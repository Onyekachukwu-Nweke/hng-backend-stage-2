package tests

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/Onyekachukwu-Nweke/hng-backend-stage-2/handlers"
	"github.com/Onyekachukwu-Nweke/hng-backend-stage-2/models"
	"github.com/Onyekachukwu-Nweke/hng-backend-stage-2/repositories"
	"github.com/Onyekachukwu-Nweke/hng-backend-stage-2/services"
	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
)

func setupRouter() *mux.Router {
	userRepo := repositories.NewMockUserRepository()
	orgRepo := repositories.NewMockOrganisationRepository()
	userService := services.NewUserService(userRepo, orgRepo)
	authHandler := handlers.NewAuthHandler(userService)

	router := mux.NewRouter()
	router.HandleFunc("/auth/register", authHandler.Register).Methods("POST")
	return router
}

func TestRegisterUserSuccessfully(t *testing.T) {
	router := setupRouter()

	payload := models.User{
		FirstName: "Onyekachukwu",
		LastName:  "Nweke",
		Email:     "werey@naster.com",
		Password:  "sky_walker",
	}

	body, _ := json.Marshal(payload)
	req, _ := http.NewRequest("POST", "/auth/register", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")

	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusCreated, rr.Code)

	var response map[string]interface{}
	json.Unmarshal(rr.Body.Bytes(), &response)

	assert.Equal(t, "success", response["status"])
	assert.Equal(t, "Registration successful", response["message"])
	assert.NotEmpty(t, response["data"].(map[string]interface{})["accessToken"])
	assert.Equal(t, "Onyekachukwu", response["data"].(map[string]interface{})["user"].(map[string]interface{})["firstName"])
	assert.Equal(t, "Nweke", response["data"].(map[string]interface{})["user"].(map[string]interface{})["lastName"])
	assert.Equal(t, "werey@naster.com", response["data"].(map[string]interface{})["user"].(map[string]interface{})["email"])
}

func TestRegisterUserValidationErrors(t *testing.T) {
	router := setupRouter()

	tests := []struct {
		payload     models.User
		expectedErr string
	}{
		{models.User{LastName: "Nweke", Email: "werey@naster.com", Password: "sky_walker"}, "FirstName is required"},
		{models.User{FirstName: "Onyekachukwu", Email: "werey@naster.com", Password: "sky_walker"}, "LastName is required"},
		{models.User{FirstName: "Onyekachukwu", LastName: "Nweke", Password: "sky_walker"}, "Email is required"},
		{models.User{FirstName: "Onyekachukwu", LastName: "Nweke", Email: "werey@naster.com"}, "Password is required"},
	}

	for _, tt := range tests {
		body, _ := json.Marshal(tt.payload)
		req, _ := http.NewRequest("POST", "/auth/register", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")

		rr := httptest.NewRecorder()
		router.ServeHTTP(rr, req)

		assert.Equal(t, http.StatusUnprocessableEntity, rr.Code)

		var response map[string]interface{}
		json.Unmarshal(rr.Body.Bytes(), &response)

		assert.Equal(t, "error", response["status"])
		assert.Contains(t, response["message"], tt.expectedErr)
	}
}

func TestRegisterDuplicateEmail(t *testing.T) {
	router := setupRouter()

	// Register the first user
	payload := models.User{
		FirstName: "Onyekachukwu",
		LastName:  "Nweke",
		Email:     "werey@naster.com",
		Password:  "sky_walker",
	}
	body, _ := json.Marshal(payload)
	req, _ := http.NewRequest("POST", "/auth/register", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)
	assert.Equal(t, http.StatusCreated, rr.Code)

	// Attempt to register the same user again
	rr = httptest.NewRecorder()
	router.ServeHTTP(rr, req)
	assert.Equal(t, http.StatusBadRequest, rr.Code)

	var response map[string]interface{}
	json.Unmarshal(rr.Body.Bytes(), &response)

	assert.Equal(t, "Bad request", response["status"])
	assert.Equal(t, "Email already exists", response["message"])
}