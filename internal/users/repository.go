package users

import (
	"context"
	"fmt"
	"log"

	"github.com/EhsanSepehriNasab/polling-platform/internal/models"
	"github.com/jackc/pgx/v5/pgxpool"
)

// UserRepository defines methods for interacting with the user table
type UserRepository struct {
	db *pgxpool.Pool
}

// NewUserRepository creates a new instance of UserRepository
func NewUserRepository(db *pgxpool.Pool) *UserRepository {
	return &UserRepository{db: db}
}

// CreateUser inserts a new user into the database
func (r *UserRepository) CreateUser(ctx context.Context, user *models.User) (*models.User, error) {
	query := `INSERT INTO users (name, email) VALUES ($1, $2) RETURNING id, created_at, updated_at`
	err := r.db.QueryRow(ctx, query, user.Name, user.Email).Scan(&user.ID, &user.CreatedAt, &user.UpdatedAt)
	if err != nil {
		log.Println("error retrieving user results:", err)
		return nil, fmt.Errorf("error inserting user: %v", err)
	}
	return user, nil
}
