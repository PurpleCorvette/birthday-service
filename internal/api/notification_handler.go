package api

import (
	"birthday-service/internal/notification"
	"log"
	"net/http"
)

type NotificationHandler struct {
	notifyService notification.NotifyService
}

func NewNotificationHandler(notifyService notification.NotifyService) *NotificationHandler {
	return &NotificationHandler{notifyService: notifyService}
}

func (h *NotificationHandler) TriggerNotifications(w http.ResponseWriter, r *http.Request) {
	err := h.notifyService.ScheduleDailyNotifications()
	if err != nil {
		log.Printf("Error triggering notifications: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}
