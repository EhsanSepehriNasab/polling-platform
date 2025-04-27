package users

import (
	"encoding/json"
	"net/http"

	"github.com/EhsanSepehriNasab/polling-platform/internal/models"
)

// UserHandler defines methods to handle HTTP requests related to users
type UserHandler struct {
	service *UserService
}

// NewUserHandler creates a new instance of UserHandler
func NewUserHandler(service *UserService) *UserHandler {
	return &UserHandler{service: service}
}

// RegisterUser godoc
// @Summary Register a new user
// @Description This endpoint registers a new user with a name and email.
// @Tags users
// @Accept json
// @Produce json
// @Param user body models.UserSignUp true "User object"
// @Success 201 {object} models.User
// @Failure 400 {object} string "Invalid request"
// @Failure 500 {object} string "Internal server error"
// @Router /users/register [post]
func (h *UserHandler) RegisterUser(w http.ResponseWriter, r *http.Request) {
	var user models.User
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	createdUser, err := h.service.RegisterUser(r.Context(), &user)
	if err != nil {
		http.Error(w, "Failed to register user", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(createdUser)
}
