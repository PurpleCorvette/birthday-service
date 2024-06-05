package notification

import (
	"birthday-service/internal/db"
	"context"
)

type NotificationSettings struct {
	UserID   int    `json:"user_id"`
	Time     string `json:"time"`
	Interval string `json:"interval"`
}

type SettingsService interface {
	UpdateSettings(settings NotificationSettings) error
	GetSettings(userID int) (NotificationSettings, error)
}

type settingsService struct {
	db db.DB
}

func NewSettingsService(db db.DB) SettingsService {
	return &settingsService{db: db}
}

func (s *settingsService) UpdateSettings(settings NotificationSettings) error {
	_, err := s.db.Exec(context.Background(),
		"UPDATE notification_settings SET time=$1, interval=$2 WHERE user_id=$3",
		settings.Time, settings.Interval, settings.UserID)
	return err
}

func (s *settingsService) GetSettings(userID int) (NotificationSettings, error) {
	var settings NotificationSettings
	err := s.db.QueryRow(context.Background(),
		"SELECT user_id, time, interval FROM notification_settings WHERE user_id=$1",
		userID).Scan(&settings.UserID, &settings.Time, &settings.Interval)
	if err != nil {
		return NotificationSettings{}, err
	}

	return settings, nil
}
