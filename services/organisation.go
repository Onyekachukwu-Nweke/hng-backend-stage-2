package services

import (
	"github.com/Onyekachukwu-Nweke/hng-backend-stage-2/models"
	"github.com/Onyekachukwu-Nweke/hng-backend-stage-2/repositories"
)

type OrganisationService interface {
	Create(org *models.Organisation) error
	GetByID(id string) (*models.Organisation, error)
	GetByUser(userID string) ([]models.Organisation, error)
}

type organisationService struct {
	orgRepo repositories.OrganisationRepository
}

func NewOrganisationService(orgRepo repositories.OrganisationRepository) OrganisationService {
	return &organisationService{orgRepo}
}

func (s *organisationService) Create(org *models.Organisation) error {
	return s.orgRepo.Create(org)
}

func (s *organisationService) GetByID(id string) (*models.Organisation, error) {
	return s.orgRepo.FindByID(id)
}

func (s *organisationService) GetByUser(userID string) ([]models.Organisation, error) {
	return s.orgRepo.FindByUser(userID)
}