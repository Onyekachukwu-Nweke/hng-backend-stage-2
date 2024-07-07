package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/Onyekachukwu-Nweke/hng-backend-stage-2/models"
	"github.com/Onyekachukwu-Nweke/hng-backend-stage-2/services"
	"github.com/Onyekachukwu-Nweke/hng-backend-stage-2/utils"
	"github.com/go-playground/validator/v10"
	"github.com/gorilla/mux"
	// "github.com/gorilla/mux"
)

type OrganisationHandler struct {
	orgService services.OrganisationService
	validate   *validator.Validate
}

func NewOrganisationHandler(orgService services.OrganisationService) *OrganisationHandler {
	return &OrganisationHandler{orgService: orgService, validate: validator.New()}
}

func (h *OrganisationHandler) Create(w http.ResponseWriter, r *http.Request) {
	var org models.Organisation
	if err := json.NewDecoder(r.Body).Decode(&org); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	// Validate organisation
	if err := h.validate.Struct(org); err != nil {
		utils.RespondWithValidationError(w, err)
		return
	}

	// Create organisation
	err := h.orgService.Create(&org)
	if err != nil {
		http.Error(w, "Client error", http.StatusInternalServerError)
		return
	}

	response := map[string]interface{}{
		"status":  "success",
		"message": "Organisation created successfully",
		"data":    org,
	}
	utils.RespondWithJSON(w, http.StatusCreated, response)
}

func (h *OrganisationHandler) GetUser(w http.ResponseWriter, r *http.Request) {
	userID := mux.Vars(r)["id"]
	user, err := h.orgService.GetUser(userID)
	if err != nil {
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}

	response := map[string]interface{}{
		"status":  "success",
		"message": "User retrieved successfully",
		"data":    user,
	}
	utils.RespondWithJSON(w, http.StatusOK, response)
}


func (h *OrganisationHandler) GetOrganisations(w http.ResponseWriter, r *http.Request) {
	claims, err := utils.ValidateJWTFromRequest(r)
	if err != nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	orgs, err := h.orgService.GetByUser(claims.UserID)
	if err != nil {
		http.Error(w, "Organisations not found", http.StatusNotFound)
		return
	}

	response := map[string]interface{}{
		"status":  "success",
		"message": "Organisations retrieved successfully",
		"data":    orgs,
	}
	utils.RespondWithJSON(w, http.StatusOK, response)
}


func (h *OrganisationHandler) GetOrganisation(w http.ResponseWriter, r *http.Request) {
	orgID := mux.Vars(r)["orgId"]
	org, err := h.orgService.GetByID(orgID)
	if err != nil {
		http.Error(w, "Organisation not found", http.StatusNotFound)
		return
	}

	response := map[string]interface{}{
		"status":  "success",
		"message": "Organisation retrieved successfully",
		"data":    org,
	}
	utils.RespondWithJSON(w, http.StatusOK, response)
}

func (h *OrganisationHandler) AddUserToOrganisation(w http.ResponseWriter, r *http.Request) {
	orgID := mux.Vars(r)["orgId"]
	var req struct {
		UserID string `json:"userId" validate:"required"`
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

	err := h.orgService.AddUserToOrganisation(orgID, req.UserID)
	if err != nil {
		http.Error(w, "Failed to add user to organisation", http.StatusInternalServerError)
		return
	}

	response := map[string]interface{}{
		"status":  "success",
		"message": "User added to organisation successfully",
	}
	utils.RespondWithJSON(w, http.StatusOK, response)
}
