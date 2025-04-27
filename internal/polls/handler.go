package polls

import (
	"encoding/json"
	"net/http"

	"github.com/EhsanSepehriNasab/polling-platform/internal/models"
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
// @Description This endpoint creates a new poll with a title, options, and tags.
// @Tags polls
// @Accept json
// @Produce json
// @Param poll body models.Poll true "Poll object"
// @Success 201 {object} models.Poll
// @Failure 400 {object} string "Invalid request"
// @Failure 500 {object} string "Internal server error"
// @Router /polls [post]
func (h *PollHandler) CreatePoll(w http.ResponseWriter, r *http.Request) {
	var poll models.Poll
	if err := json.NewDecoder(r.Body).Decode(&poll); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	createdPoll, err := h.service.CreatePoll(r.Context(), &poll)
	if err != nil {
		http.Error(w, "Failed to create poll", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(createdPoll)
}
