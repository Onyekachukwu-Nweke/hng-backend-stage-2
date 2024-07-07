package repositories

import "github.com/Onyekachukwu-Nweke/hng-backend-stage-2/models"

type MockOrganisationRepository struct {
	organisations map[string]*models.Organisation
	userOrgs      map[string][]string
}

// FindUserByID implements OrganisationRepository.
func (r *MockOrganisationRepository) FindUserByID(userID string) (*models.User, error) {
	panic("unimplemented")
}

func NewMockOrganisationRepository() *MockOrganisationRepository {
	return &MockOrganisationRepository{
		organisations: make(map[string]*models.Organisation),
		userOrgs:      make(map[string][]string),
	}
}

func (r *MockOrganisationRepository) Create(org *models.Organisation) error {
	r.organisations[org.OrgID] = org
	return nil
}

func (r *MockOrganisationRepository) FindByID(id string) (*models.Organisation, error) {
	org, exists := r.organisations[id]
	if !exists {
		return nil, nil
	}
	return org, nil
}

func (r *MockOrganisationRepository) FindByUser(userID string) ([]models.Organisation, error) {
	var orgs []models.Organisation
	for _, orgID := range r.userOrgs[userID] {
		org, exists := r.organisations[orgID]
		if exists {
			orgs = append(orgs, *org)
		}
	}
	return orgs, nil
}

func (r *MockOrganisationRepository) AddUserToOrganisation(orgID, userID string) error {
	r.userOrgs[userID] = append(r.userOrgs[userID], orgID)
	return nil
}
