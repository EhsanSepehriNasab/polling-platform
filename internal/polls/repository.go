package polls

import (
	"context"
	"fmt"
	"log"

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
	query := `INSERT INTO polls (title, options, tags, user_id) VALUES ($1, $2, $3, $4) RETURNING id, created_at`
	err := r.db.QueryRow(ctx, query, poll.Title, poll.Options, poll.Tags, poll.UserID).Scan(&poll.ID, &poll.CreatedAt)
	if err != nil {
		return nil, fmt.Errorf("error inserting poll: %v", err)
	}
	return poll, nil
}

// GetPollsForFeed retrieves polls for a user's feed that they have not voted on or skipped
func (r *PollRepository) GetPollsForFeed(ctx context.Context, userID int, tag string, page int, limit int) ([]models.Poll, error) {
	var polls []models.Poll

	query := `SELECT id, title, options, tags, user_id, created_at, updated_at 
              FROM polls 
              WHERE id NOT IN (
                  SELECT poll_id 
                  FROM user_votes 
                  WHERE user_id = $1
              )`

	if tag != "" {
		query += ` AND tags @> $2::text[] ORDER BY created_at DESC LIMIT $3 OFFSET $4`
	} else {
		query += ` ORDER BY created_at DESC LIMIT $2 OFFSET $3`
	}

	args := []interface{}{userID, limit, (page - 1) * limit}

	if tag != "" {
		args = append([]interface{}{userID, []string{tag}}, args[1:]...)
	}

	rows, err := r.db.Query(ctx, query, args...)
	if err != nil {
		log.Println("error retrieving polls: %v", err, query)
		return nil, fmt.Errorf("error retrieving polls: %v", err)
	}
	defer rows.Close()

	for rows.Next() {
		var poll models.Poll
		if err := rows.Scan(&poll.ID, &poll.Title, &poll.Options, &poll.Tags, &poll.UserID, &poll.CreatedAt, &poll.UpdatedAt); err != nil {
			log.Println("error scanning poll: %v", err)
			return nil, fmt.Errorf("error scanning poll: %v", err)
		}
		polls = append(polls, poll)
	}

	if err := rows.Err(); err != nil {
		log.Println("error with rows: %v", err)
		return nil, fmt.Errorf("error with rows: %v", err)
	}

	return polls, nil
}

// SaveVote stores a user's vote for a poll option
func (r *PollRepository) SaveVote(ctx context.Context, userID, pollID, optionIndex int) error {
	query := `
		INSERT INTO user_votes (user_id, poll_id, option_index)
		VALUES ($1, $2, $3)
		ON CONFLICT (user_id, poll_id) DO NOTHING
	`
	commandTag, err := r.db.Exec(ctx, query, userID, pollID, optionIndex)
	if err != nil {
		return fmt.Errorf("error saving vote: %w", err)
	}
	if commandTag.RowsAffected() == 0 {
		return fmt.Errorf("user has already voted on this poll")
	}
	return nil
}

// SkipPoll marks a poll as skipped for a user.
func (r *PollRepository) SkipPoll(ctx context.Context, userID, pollID int) error {
	query := `INSERT INTO user_votes (user_id, poll_id, skipped) 
              VALUES ($1, $2, TRUE) 
              ON CONFLICT (user_id, poll_id) 
              DO UPDATE SET skipped = TRUE`

	_, err := r.db.Exec(ctx, query, userID, pollID)
	if err != nil {
		return fmt.Errorf("error skipping poll: %v", err)
	}
	return nil
}

// HasUserVotedOrSkipped checks if the user has voted or skipped this poll.
func (r *PollRepository) HasUserVotedOrSkipped(ctx context.Context, userID, pollID int) (bool, error) {
	var count int
	query := `SELECT COUNT(*) FROM user_votes WHERE user_id = $1 AND poll_id = $2`
	err := r.db.QueryRow(ctx, query, userID, pollID).Scan(&count)
	if err != nil {
		return false, fmt.Errorf("error checking if user has voted or skipped: %v", err)
	}

	return count > 0, nil
}

// GetPollResults retrieves the results for a specific poll.
func (r *PollRepository) GetPollResults(ctx context.Context, pollID int) (*models.PollResults, error) {
	var results models.PollResults
	results.PollID = pollID

	query := `
	SELECT 
	options[i] AS option_text,
	(
		SELECT COUNT(*) 
		FROM user_votes 
		WHERE poll_id = polls.id AND option_index = i - 1
	) AS vote_count
	FROM polls,
	generate_subscripts(options, 1) AS i
	WHERE polls.id = $1;

	`

	rows, err := r.db.Query(ctx, query, pollID)
	if err != nil {
		log.Println("error retrieving poll results:", err)
		return nil, fmt.Errorf("error retrieving poll results: %v", err)
	}
	defer rows.Close()

	for rows.Next() {
		var option string
		var voteCount int
		if err := rows.Scan(&option, &voteCount); err != nil {
			log.Println("error scanning poll result row:", err)
			return nil, fmt.Errorf("error scanning poll result row: %v", err)
		}
		results.Votes = append(results.Votes, models.OptionResult{
			Option: option,
			Count:  voteCount,
		})
	}

	if err := rows.Err(); err != nil {
		log.Println("error with rows:", err)
		return nil, fmt.Errorf("error with rows: %v", err)
	}

	return &results, nil
}
