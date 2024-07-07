package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/Onyekachukwu-Nweke/hng-backend-stage-2/models"
	"github.com/Onyekachukwu-Nweke/hng-backend-stage-2/services"
	"github.com/Onyekachukwu-Nweke/hng-backend-stage-2/utils"
	"github.com/go-playground/validator/v10"
	// "github.com/gorilla/mux"
)

type AuthHandler struct {
	userService services.UserService
	validate    *validator.Validate
}

func NewAuthHandler(userService services.UserService) *AuthHandler {
	return &AuthHandler{userService: userService, validate: validator.New()}
}

func (h *AuthHandler) Register(w http.ResponseWriter, r *http.Request) {
	var user models.User
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	// Validate user
	if err := h.validate.Struct(user); err != nil {
		utils.RespondWithValidationError(w, err)
		return
	}

	// Register user
	createdUser, err := h.userService.Register(&user)
	if err != nil {
		http.Error(w, "Registration unsuccessful", http.StatusInternalServerError)
		return
	}

	// Generate JWT token
	token, err := utils.GenerateJWT(createdUser.Email, createdUser.UserID)
	if err != nil {
		http.Error(w, "Registration unsuccessful", http.StatusInternalServerError)
		return
	}

	response := map[string]interface{}{
		"status":  "success",
		"message": "Registration successful",
		"data": map[string]interface{}{
			"accessToken": token,
			"user":        createdUser,
		},
	}
	utils.RespondWithJSON(w, http.StatusCreated, response)
}

func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Email    string `json:"email" validate:"required,email"`
		Password string `json:"password" validate:"required"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	// Validate request
	if err := h.validate.Struct(req); err != nil {
		utils.RespondWithValidationError(w, err)
		return
	}

	// Authenticate user
	user, err := h.userService.Login(req.Email, req.Password)
	if err != nil {
		http.Error(w, "Authentication failed", http.StatusUnauthorized)
		return
	}

	// Generate JWT token
	token, err := utils.GenerateJWT(user.Email, user.UserID)
	if err != nil {
		http.Error(w, "Authentication failed", http.StatusInternalServerError)
		return
	}

	response := map[string]interface{}{
		"status":  "success",
		"message": "Login successful",
		"data": map[string]interface{}{
			"accessToken": token,
			"user":        user,
		},
	}
	utils.RespondWithJSON(w, http.StatusOK, response)
}