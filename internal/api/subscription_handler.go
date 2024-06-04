package api

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"

	"birthday-service/internal/notification"
)

type SubscriptionHandler struct {
	SubscriptionService notification.SubscriptionService
	SettingsService     notification.SettingsService
}

func NewSubscriptionHandler(subService notification.SubscriptionService, setService notification.SettingsService) *SubscriptionHandler {
	return &SubscriptionHandler{
		SubscriptionService: subService,
		SettingsService:     setService,
	}
}

func (h *SubscriptionHandler) GetSubscriptions(w http.ResponseWriter, r *http.Request) {
	userID, err := strconv.Atoi(r.URL.Query().Get("user_id"))
	if err != nil {
		http.Error(w, "Invalid user ID", http.StatusBadRequest)
		return
	}

	subscriptions, err := h.SubscriptionService.ListSubscriptions(userID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(subscriptions)
}

func (h *SubscriptionHandler) CreateSubscription(w http.ResponseWriter, r *http.Request) {
	var sub notification.Subscription
	if err := json.NewDecoder(r.Body).Decode(&sub); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	newSub, err := h.SubscriptionService.Subscribe(sub.UserID, sub.EmployeeID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(newSub)
}

func (h *SubscriptionHandler) DeleteSubscription(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	userID, err := strconv.Atoi(params["userID"])
	if err != nil {
		http.Error(w, "Invalid user ID", http.StatusBadRequest)
		return
	}
	employeeID, err := strconv.Atoi(params["employeeID"])
	if err != nil {
		http.Error(w, "Invalid employee ID", http.StatusBadRequest)
		return
	}

	err = h.SubscriptionService.Unsubscribe(userID, employeeID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (h *SubscriptionHandler) UpdateNotificationSettings(w http.ResponseWriter, r *http.Request) {
	var settings notification.UserNotificationSettings
	if err := json.NewDecoder(r.Body).Decode(&settings); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	updatedSettings, err := h.SettingsService.UpdateNotificationSettings(settings.UserID, settings.NotifyTime)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(updatedSettings)
}
