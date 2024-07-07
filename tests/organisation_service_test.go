package tests

import (
	"testing"

	"github.com/Onyekachukwu-Nweke/hng-backend-stage-2/models"
	"github.com/Onyekachukwu-Nweke/hng-backend-stage-2/repositories"
	"github.com/Onyekachukwu-Nweke/hng-backend-stage-2/services"
	"github.com/stretchr/testify/assert"
)

func TestOrganisationAccess(t *testing.T) {
	// Mock data and repositories
	userRepo := repositories.NewMockUserRepository()
	orgRepo := repositories.NewMockOrganisationRepository()

	// Create user and organisation services
	// userService := services.NewUserService(userRepo, orgRepo)
	orgService := services.NewOrganisationService(orgRepo)

	// Create test user
	user := &models.User{
		UserID:    "user1",
		FirstName: "John",
		LastName:  "Doe",
		Email:     "john.doe@example.com",
		Password:  "password123",
	}
	userRepo.Create(user)

	// Create test organisation
	org := &models.Organisation{
		OrgID:       "org1",
		Name:        "John's Organisation",
		Description: "Test Organisation",
	}
	orgRepo.Create(org)
	orgRepo.AddUserToOrganisation(org.OrgID, user.UserID)

	// Attempt to access organisation by another user
	anotherUser := &models.User{
		UserID:    "user2",
		FirstName: "Jane",
		LastName:  "Smith",
		Email:     "jane.smith@example.com",
		Password:  "password123",
	}
	userRepo.Create(anotherUser)

	organisations, err := orgService.GetByUser(anotherUser.UserID)
	assert.NoError(t, err)
	assert.Empty(t, organisations)
}
