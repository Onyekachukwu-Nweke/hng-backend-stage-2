package tests

import (
	"testing"
	"time"

	"github.com/Onyekachukwu-Nweke/hng-backend-stage-2/utils"
	"github.com/stretchr/testify/assert"
)

func TestGenerateJWT(t *testing.T) {
	email := "test@example.com"
	userID := "12345"

	token, err := utils.GenerateJWT(email, userID)
	assert.NoError(t, err)
	assert.NotEmpty(t, token)

	claims, err := utils.ValidateJWT(token)
	assert.NoError(t, err)
	assert.Equal(t, email, claims.Email)
	assert.Equal(t, userID, claims.UserID)

	expirationTime := claims.RegisteredClaims.ExpiresAt
	assert.NotNil(t, expirationTime)
	assert.True(t, time.Unix(expirationTime.Unix(), 0).After(time.Now()))
	// assert.True(t, time.Unix(expirationTime, 0).After(time.Now()))
}
