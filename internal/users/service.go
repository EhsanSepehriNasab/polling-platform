package users

import (
	"context"

	"github.com/EhsanSepehriNasab/polling-platform/internal/models"
)

// UserService defines methods for user-related business logic
type UserService struct {
	repo *UserRepository
}

// NewUserService creates a new instance of UserService
func NewUserService(repo *UserRepository) *UserService {
	return &UserService{repo: repo}
}

// RegisterUser registers a new user in the system
func (s *UserService) RegisterUser(ctx context.Context, user *models.User) (*models.User, error) {
	return s.repo.CreateUser(ctx, user)
}
