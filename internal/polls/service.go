package polls

import (
	"context"
	"fmt"
	"log"

	"github.com/EhsanSepehriNasab/polling-platform/internal/models"
)

// PollService defines methods for business logic related to polls
type PollService struct {
	repo *PollRepository
}

// NewPollService creates a new instance of PollService
func NewPollService(repo *PollRepository) *PollService {
	return &PollService{repo: repo}
}

// CreatePoll handles the logic for creating a new poll
func (s *PollService) CreatePoll(ctx context.Context, poll *models.Poll) (*models.Poll, error) {
	return s.repo.CreatePoll(ctx, poll)
}

// GetPollsForFeed retrieves polls for a user’s feed, excluding those they’ve voted or skipped
func (s *PollService) GetPollsForFeed(ctx context.Context, userID int, tag string, page int, limit int) ([]models.Poll, error) {
	// Use the repository method to get the polls for the feed
	return s.repo.GetPollsForFeed(ctx, userID, tag, page, limit)
}

// VotePoll lets a user vote for a poll option
func (s *PollService) VotePoll(ctx context.Context, userID, pollID, optionIndex int) error {
	// Optional: Validate poll exists and optionIndex is valid
	return s.repo.SaveVote(ctx, userID, pollID, optionIndex)
}

// SkipPoll allows a user to skip a poll.
func (s *PollService) SkipPoll(ctx context.Context, userID, pollID int) error {
	// Check if the user has already voted or skipped this poll
	hasVotedOrSkipped, err := s.repo.HasUserVotedOrSkipped(ctx, userID, pollID)
	if err != nil {
		return fmt.Errorf("error checking if user has voted or skipped: %v", err)
	}

	if hasVotedOrSkipped {
		return fmt.Errorf("user has already voted or skipped this poll")
	}

	// Mark the poll as skipped for the user
	err = s.repo.SkipPoll(ctx, userID, pollID)
	if err != nil {
		return fmt.Errorf("error skipping poll: %v", err)
	}

	return nil
}

// GetPollResults retrieves the results for a specific poll.
func (s *PollService) GetPollResults(ctx context.Context, pollID int) (*models.PollResults, error) {
	// Retrieve the results from the repository
	results, err := s.repo.GetPollResults(ctx, pollID)
	if err != nil {
		log.Println("error retrieving polls: %v", err)
		return nil, fmt.Errorf("error retrieving poll results: %v", err)
	}

	return results, nil
}
