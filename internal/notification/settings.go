package notification

import "errors"

type UserNotificationSettings struct {
	UserID     int    `json:"user_id"`
	NotifyTime string `json:"notify_time"` // Формат HH:MM
}

type SettingsService interface {
	UpdateNotificationSettings(userID int, notifyTime string) (UserNotificationSettings, error)
	GetNotificationSettings(userID int) (UserNotificationSettings, error)
}

type settingsService struct {
	settings map[int]UserNotificationSettings
}

func NewSettingsService() SettingsService {
	return &settingsService{
		settings: make(map[int]UserNotificationSettings),
	}
}

func (s *settingsService) UpdateNotificationSettings(userID int, notifyTime string) (UserNotificationSettings, error) {
	settings := UserNotificationSettings{
		UserID:     userID,
		NotifyTime: notifyTime,
	}
	s.settings[userID] = settings
	return settings, nil
}

func (s *settingsService) GetNotificationSettings(userID int) (UserNotificationSettings, error) {
	settings, exists := s.settings[userID]
	if !exists {
		return UserNotificationSettings{}, errors.New("settings not found")
	}
	return settings, nil
}
