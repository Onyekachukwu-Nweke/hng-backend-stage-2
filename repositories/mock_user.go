package repositories

import "github.com/Onyekachukwu-Nweke/hng-backend-stage-2/models"

type MockUserRepository struct {
	users map[string]*models.User
}

func NewMockUserRepository() *MockUserRepository {
	return &MockUserRepository{users: make(map[string]*models.User)}
}

func (r *MockUserRepository) Create(user *models.User) error {
	r.users[user.Email] = user
	return nil
}

func (r *MockUserRepository) FindByEmail(email string) (*models.User, error) {
	user, exists := r.users[email]
	if !exists {
		return nil, nil
	}
	return user, nil
}

func (r *MockUserRepository) FindByID(id string) (*models.User, error) {
	for _, user := range r.users {
		if user.UserID == id {
			return user, nil
		}
	}
	return nil, nil
}
