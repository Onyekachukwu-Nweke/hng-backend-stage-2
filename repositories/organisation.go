package repositories

import (
	"github.com/Onyekachukwu-Nweke/hng-backend-stage-2/models"
	"gorm.io/gorm"
)


type OrganisationRepository interface {
	Create(org *models.Organisation) error
	FindByID(id string) (*models.Organisation, error)
	FindByUser(userID string) ([]models.Organisation, error)
}

type organisationRepository struct {
	db *gorm.DB
}

func NewOrganisationRepository(db *gorm.DB) OrganisationRepository {
	return &organisationRepository{db}
}

func (r *organisationRepository) Create(org *models.Organisation) error {
	return r.db.Create(org).Error
}

func (r *organisationRepository) FindByID(id string) (*models.Organisation, error) {
	var org models.Organisation
	err := r.db.Where("org_id = ?", id).First(&org).Error
	return &org, err
}

func (r *organisationRepository) FindByUser(userID string) ([]models.Organisation, error) {
	var orgs []models.Organisation
	err := r.db.Model(&models.User{UserID: userID}).Association("Organisations").Find(&orgs)
	return orgs, err
}