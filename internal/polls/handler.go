package polls

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/EhsanSepehriNasab/polling-platform/internal/cache"
	"github.com/EhsanSepehriNasab/polling-platform/internal/metrics"
	"github.com/EhsanSepehriNasab/polling-platform/internal/models"
	"github.com/go-chi/chi/v5"
	"github.com/go-redis/redis/v8"
)

// PollHandler defines methods to handle HTTP requests related to polls
type PollHandler struct {
	service *PollService
}

// NewPollHandler creates a new instance of PollHandler
func NewPollHandler(service *PollService) *PollHandler {
	return &PollHandler{service: service}
}

// CreatePoll godoc
// @Summary Create a new poll
// @Description This endpoint creates a new poll with a title, options, tags, and an associated user.
// @Tags polls
// @Accept json
// @Produce json
// @Param userId header int true "User ID of the creator"
// @Param poll body models.PollRequest true "Poll object"
// @Success 201 {object} models.Poll
// @Failure 400 {object} string "Invalid request"
// @Failure 500 {object} string "Internal server error"
// @Router /polls [post]
func (h *PollHandler) CreatePoll(w http.ResponseWriter, r *http.Request) {
	userID := r.Header.Get("userId")

	var poll models.Poll
	if err := json.NewDecoder(r.Body).Decode(&poll); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Convert userID to int and associate it with the poll
	userIDInt, err := strconv.Atoi(userID)
	if err != nil {
		http.Error(w, "Invalid user ID", http.StatusBadRequest)
		return
	}
	poll.UserID = userIDInt

	// Create the poll in the database using the service layer
	createdPoll, err := h.service.CreatePoll(r.Context(), &poll)
	if err != nil {
		http.Error(w, "Failed to create poll", http.StatusInternalServerError)
		return
	}

	for _, tag := range poll.Tags {
		cacheKeyPattern := fmt.Sprintf("user:*:feed:%s:*:*", tag)
		if err := cache.DeleteKeysByPattern(cacheKeyPattern); err != nil {
			log.Fatalf("Error deleting keys for tag %s: %v", tag, err)
		}
	}

	if err := cache.DeleteKeysByPattern("user:*:feed::*:*"); err != nil {
		log.Fatalf("Error deleting keys for tag: %v", err)
	}

	// Respond with the created poll data
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(createdPoll)
}

// PollFeedHandler godoc
// @Summary Retrieve a list of polls for the user's feed
// @Description This endpoint retrieves a list of polls that the user has not voted on or skipped, filtered by optional tag, with pagination support.
// @Tags polls
// @Accept json
// @Produce json
// @Param userId header int true "User ID of the creator"
// @Param tag query string false "Filter by tag"
// @Param page query int false "Page number"
// @Param limit query int false "Number of items per page"
// @Success 200 {array} models.Poll
// @Failure 400 {object} string "Invalid request"
// @Failure 500 {object} string "Internal server error"
// @Router /polls [get]
func (h *PollHandler) PollFeedHandler(w http.ResponseWriter, r *http.Request) {
	start := time.Now() // For PollFeedDuration
	defer func() {
		duration := time.Since(start).Seconds()
		metrics.PollFeedDuration.Observe(duration)
	}()

	// Get query parameters
	userIDStr := r.Header.Get("userId")
	tag := r.URL.Query().Get("tag")
	pageStr := r.URL.Query().Get("page")
	limitStr := r.URL.Query().Get("limit")

	// Validate and convert parameters
	userID, err := strconv.Atoi(userIDStr)
	if err != nil {
		http.Error(w, "Invalid userId", http.StatusBadRequest)
		return
	}

	page, err := strconv.Atoi(pageStr)
	if err != nil || page < 1 {
		page = 1
	}

	limit, err := strconv.Atoi(limitStr)
	if err != nil || limit < 1 {
		limit = 10
	}

	// Check Redis cache for the feed
	cacheKey := fmt.Sprintf("user:%d:feed:%s:%d:%d", userID, tag, page, limit)
	cachedFeed, err := cache.GetCacheClient().Get(context.Background(), cacheKey).Result()

	if err == nil {
		// Cache hit
		metrics.CacheHits.Inc()

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(cachedFeed))
		return
	} else {
		// Cache miss
		metrics.CacheMisses.Inc()
	}

	// Time the DB query
	dbStart := time.Now()

	polls, err := h.service.GetPollsForFeed(r.Context(), userID, tag, page, limit)

	dbDuration := time.Since(dbStart).Seconds()
	metrics.DBQueryDuration.Observe(dbDuration)

	if err != nil {
		http.Error(w, "Failed to retrieve polls", http.StatusInternalServerError)
		return
	}

	// Cache the response
	cachedFeedJSON, err := json.Marshal(polls)
	if err != nil {
		http.Error(w, "Failed to encode polls", http.StatusInternalServerError)
		return
	}
	cache.GetCacheClient().Set(context.Background(), cacheKey, cachedFeedJSON, time.Minute*10)

	// Respond
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(cachedFeedJSON)
}

// VoteRequest represents the request body for voting
type VoteRequest struct {
	OptionIndex int `json:"optionIndex"`
}

// VotePollHandler godoc
// @Summary Vote for a poll option
// @Description Allows a user to vote for a specific poll option.
// @Tags polls
// @Accept json
// @Produce json
// @Param userId header int true "User ID"
// @Param pollID path int true "Poll ID"
// @Param voteRequest body VoteRequest true "Vote Request"
// @Success 200 {string} string "Vote successful"
// @Failure 400 {string} string "Bad request"
// @Failure 409 {string} string "Already voted"
// @Failure 500 {string} string "Internal server error"
// @Router /polls/{pollID}/vote [post]
func (h *PollHandler) VotePollHandler(w http.ResponseWriter, r *http.Request) {
	userIDStr := r.Header.Get("userId")
	if userIDStr == "" {
		log.Printf("Missing userId header")
		http.Error(w, "Missing userId header", http.StatusBadRequest)
		return
	}
	userID, err := strconv.Atoi(userIDStr)
	if err != nil {
		log.Printf("Invalid userId header", err)
		http.Error(w, "Invalid userId header", http.StatusBadRequest)
		return
	}

	// Check Redis for the rate limit
	redisKey := fmt.Sprintf("user:%d:votes:today", userID)
	votesToday, err := cache.GetCacheClient().Get(context.Background(), redisKey).Int()
	if err != nil && err != redis.Nil {
		log.Printf("Internal server error", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	// Allow a user to vote only 100 times a day
	if votesToday >= 100 {
		log.Printf("Rate limit exceeded, you can vote up to 100 polls per day")
		http.Error(w, "Rate limit exceeded, you can vote up to 100 polls per day", http.StatusTooManyRequests)
		return
	}

	// Proceed with the vote logic
	pollIDStr := chi.URLParam(r, "pollID")
	pollID, err := strconv.Atoi(pollIDStr)
	if err != nil {
		log.Printf("Invalid poll ID", err)
		http.Error(w, "Invalid poll ID", http.StatusBadRequest)
		return
	}

	var payload struct {
		OptionIndex int `json:"optionIndex"`
	}
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		log.Printf("Invalid request body", err)
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if payload.OptionIndex < 0 {
		log.Printf("Option index must be non-negative")
		http.Error(w, "Option index must be non-negative", http.StatusBadRequest)
		return
	}

	// Voting logic here
	err = h.service.VotePoll(r.Context(), userID, pollID, payload.OptionIndex)
	if err != nil {
		if err.Error() == "user has already voted on this poll" {
			log.Printf("User already voted on this poll", err)
			http.Error(w, "User already voted on this poll", http.StatusConflict)
		} else {
			log.Printf("Failed to vote", err)
			http.Error(w, "Failed to vote", http.StatusInternalServerError)
		}
		return
	}

	// Increment the user's vote count in Redis (expire at midnight)
	cache.GetCacheClient().Incr(context.Background(), redisKey)
	cache.GetCacheClient().ExpireAt(context.Background(), redisKey, time.Now().Truncate(24*time.Hour).Add(24*time.Hour))

	//TODO: We can also check what is the tags of pollID and delete those tags from user cache.
	// Delete feed cache (already done)
	cacheKeyPattern := fmt.Sprintf("user:%d:feed:*:*:*", userID)
	if err := cache.DeleteKeysByPattern(cacheKeyPattern); err != nil {
		log.Printf("Error deleting feed cache keys: %v", err)
	}

	// Delete poll result cache
	pollResultCacheKey := fmt.Sprintf("poll:%d:results", pollID)
	if err := cache.GetCacheClient().Del(context.Background(), pollResultCacheKey).Err(); err != nil {
		log.Printf("Error deleting poll result cache: %v", err)
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Vote successful"))
}

// SkipPollHandler godoc
// @Summary Skip a poll
// @Description Allows a user to skip a poll without voting.
// @Tags polls
// @Accept json
// @Produce json
// @Param userId header int true "User ID"
// @Param pollID path int true "Poll ID"
// @Success 200 {string} string "Poll skipped successfully"
// @Failure 400 {string} string "Bad request"
// @Failure 409 {string} string "User already skipped or voted on this poll"
// @Failure 500 {string} string "Internal server error"
// @Router /polls/{pollID}/skip [post]
func (h *PollHandler) SkipPollHandler(w http.ResponseWriter, r *http.Request) {
	// Extract userID from the header
	userIDStr := r.Header.Get("userId")
	if userIDStr == "" {
		log.Printf("Missing userId header")
		http.Error(w, "Missing userId header", http.StatusBadRequest)
		return
	}
	userID, err := strconv.Atoi(userIDStr)
	if err != nil {
		log.Printf("Invalid userId header")
		http.Error(w, "Invalid userId header", http.StatusBadRequest)
		return
	}

	// Extract pollID from the URL parameters
	pollIDStr := chi.URLParam(r, "pollID")
	pollID, err := strconv.Atoi(pollIDStr)
	if err != nil {
		log.Printf("Invalid poll ID")
		http.Error(w, "Invalid poll ID", http.StatusBadRequest)
		return
	}

	// Call the service to skip the poll
	err = h.service.SkipPoll(r.Context(), userID, pollID)
	if err != nil {
		if err.Error() == "user has already voted or skipped this poll" {
			log.Printf("User already skipped or vote on this poll", err)
			http.Error(w, "User already voted or skipped this poll", http.StatusConflict)
		} else {
			log.Printf("Failed to skip poll", err)
			http.Error(w, "Failed to skip poll", http.StatusInternalServerError)
		}
		return
	}

	cacheKeyPattern := fmt.Sprintf("user:%d:feed:*:*:*", userID)
	if err := cache.DeleteKeysByPattern(cacheKeyPattern); err != nil {
		log.Fatalf("Error deleting keys for : %v", err)
	}

	// Respond with success
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Poll skipped successfully"))
}

// GetPollResultsHandler godoc
// @Summary Get results for a poll
// @Description This endpoint retrieves the results of a specific poll, including the number of votes for each option.
// @Tags polls
// @Accept json
// @Produce json
// @Param pollID path int true "Poll ID"
// @Success 200 {object} models.PollResults "Poll results with option votes"
// @Failure 400 {string} string "Bad request"
// @Failure 500 {string} string "Internal server error"
// @Router /polls/{pollID}/stats [get]
func (h *PollHandler) GetPollResultsHandler(w http.ResponseWriter, r *http.Request) {
	// Extract pollID from URL parameters
	pollIDStr := chi.URLParam(r, "pollID")
	pollID, err := strconv.Atoi(pollIDStr)
	if err != nil {
		http.Error(w, "Invalid poll ID", http.StatusBadRequest)
		return
	}

	// Try to get results from Redis
	cacheKey := fmt.Sprintf("poll:%d:results", pollID)
	cachedResults, err := cache.GetCacheClient().Get(context.Background(), cacheKey).Result()
	if err == nil {
		// Cache hit: return cached results
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(cachedResults))
		return
	}

	// Cache miss: Fetch from DB
	results, err := h.service.GetPollResults(r.Context(), pollID)
	if err != nil {
		http.Error(w, "Failed to get poll results", http.StatusInternalServerError)
		return
	}

	// Serialize results to JSON
	resultJSON, err := json.Marshal(results)
	if err != nil {
		http.Error(w, "Error encoding response", http.StatusInternalServerError)
		return
	}

	// Store in Redis for 1 minute
	cache.GetCacheClient().Set(context.Background(), cacheKey, resultJSON, time.Minute)

	// Return response
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(resultJSON)
}
