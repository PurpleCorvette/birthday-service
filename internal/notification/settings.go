package notification

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v4"

	"birthday-service/internal/db"
)

type UserNotificationSettings struct {
	UserID     int    `json:"user_id"`
	NotifyTime string `json:"notify_time"` // Формат HH:MM
}

type SettingsService interface {
	UpdateNotificationSettings(userID int, notifyTime string) (UserNotificationSettings, error)
	GetNotificationSettings(userID int) (UserNotificationSettings, error)
}

type settingsService struct{}

func NewSettingsService() SettingsService {
	return &settingsService{}
}

func (s *settingsService) UpdateNotificationSettings(userID int, notifyTime string) (UserNotificationSettings, error) {
	var settings UserNotificationSettings
	err := db.Conn.QueryRow(context.Background(),
		"INSERT INTO settings (user_id, notify_time) VALUES ($1, $2) ON CONFLICT (user_id) DO UPDATE SET notify_time = $2 RETURNING user_id, notify_time",
		userID, notifyTime).Scan(&settings.UserID, &settings.NotifyTime)
	if err != nil {
		return UserNotificationSettings{}, err
	}

	return settings, nil
}

func (s *settingsService) GetNotificationSettings(userID int) (UserNotificationSettings, error) {
	var settings UserNotificationSettings
	err := db.Conn.QueryRow(context.Background(),
		"SELECT user_id, notify_time FROM settings WHERE user_id=$1", userID).Scan(&settings.UserID, &settings.NotifyTime)
	if err == pgx.ErrNoRows {
		return UserNotificationSettings{}, errors.New("settings not found")
	}
	if err != nil {
		return UserNotificationSettings{}, err
	}

	return settings, nil
}
