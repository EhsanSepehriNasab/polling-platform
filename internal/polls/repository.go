package polls

import (
	"context"
	"fmt"

	"github.com/EhsanSepehriNasab/polling-platform/internal/models"
	"github.com/jackc/pgx/v5/pgxpool"
)

// PollRepository defines methods for interacting with the database
type PollRepository struct {
	db *pgxpool.Pool
}

// NewPollRepository creates a new instance of PollRepository
func NewPollRepository(db *pgxpool.Pool) *PollRepository {
	return &PollRepository{db: db}
}

// CreatePoll inserts a new poll into the database
func (r *PollRepository) CreatePoll(ctx context.Context, poll *models.Poll) (*models.Poll, error) {
	query := `INSERT INTO polls (title, options, tags) VALUES ($1, $2, $3) RETURNING id, created_at`
	err := r.db.QueryRow(ctx, query, poll.Title, poll.Options, poll.Tags).Scan(&poll.ID, &poll.CreatedAt)
	if err != nil {
		return nil, fmt.Errorf("error inserting poll: %v", err)
	}
	return poll, nil
}
