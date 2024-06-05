package notification

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestUpdateNotificationSettings(t *testing.T) {
	service := NewSettingsService()

	settings, err := service.UpdateNotificationSettings(1, "09:00")
	assert.NoError(t, err)
	assert.Equal(t, "09:00", settings.NotifyTime)

	settings, err = service.GetNotificationSettings(1)
	assert.NoError(t, err)
	assert.Equal(t, "09:00", settings.NotifyTime)
}
