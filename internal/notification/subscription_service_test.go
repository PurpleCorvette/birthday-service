package notification

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSubscribe(t *testing.T) {
	service := NewSubscriptionService()

	sub, err := service.Subscribe(1, 1)
	assert.NoError(t, err)
	assert.Equal(t, 1, sub.UserID)
	assert.Equal(t, 1, sub.EmployeeID)

	_, err = service.Subscribe(1, 1)
	assert.Error(t, err)
}

func TestUnsubscribe(t *testing.T) {
	service := NewSubscriptionService()
	service.Subscribe(1, 1)

	err := service.Unsubscribe(1, 1)
	assert.NoError(t, err)

	err = service.Unsubscribe(1, 1)
	assert.Error(t, err)
}
