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
