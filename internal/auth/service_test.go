package auth

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRegister(t *testing.T) {
	service := NewAuthService()

	user, err := service.Register("testuser", "password")
	assert.NoError(t, err)
	assert.Equal(t, "testuser", user.Username)

	_, err = service.Register("testuser", "password")
	assert.Error(t, err)
}

func TestLogin(t *testing.T) {
	service := NewAuthService()
	service.Register("testuser", "password")

	user, err := service.Login("testuser", "password")
	assert.NoError(t, err)
	assert.Equal(t, "testuser", user.Username)

	_, err = service.Login("testuser", "wrongpassword")
	assert.Error(t, err)
}
