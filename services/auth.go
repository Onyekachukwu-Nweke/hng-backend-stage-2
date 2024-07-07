package services

import (
	"errors"

	"github.com/Onyekachukwu-Nweke/hng-backend-stage-2/models"
	"github.com/Onyekachukwu-Nweke/hng-backend-stage-2/repositories"
	// "github.com/Onyekachukwu-Nweke/hng-backend-stage-2/utils"
	"golang.org/x/crypto/bcrypt"
)


type UserService interface {
	Register(user *models.User) (*models.User, error)
	Login(email, password string) (*models.User, error)
}

type userService struct {
	userRepo repositories.UserRepository
	orgRepo  repositories.OrganisationRepository
}

func NewUserService(userRepo repositories.UserRepository, orgRepo repositories.OrganisationRepository) UserService {
	return &userService{userRepo, orgRepo}
}

func (s *userService) Register(user *models.User) (*models.User, error) {
	// Hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}
	user.Password = string(hashedPassword)

	// Create user
	err = s.userRepo.Create(user)
	if err != nil {
		return nil, err
	}

	// Create default organisation
	org := &models.Organisation{
		Name: user.FirstName + "'s Organisation",
	}
	err = s.orgRepo.Create(org)
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (s *userService) Login(email, password string) (*models.User, error) {
	user, err := s.userRepo.FindByEmail(email)
	if err != nil {
		return nil, errors.New("invalid email or password")
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		return nil, errors.New("invalid email or password")
	}

	return user, nil
}