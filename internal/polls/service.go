package polls

import (
	"context"

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
