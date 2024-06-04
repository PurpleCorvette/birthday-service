package notification

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

type SubscriptionHandler struct {
	service SubscriptionService
}

func NewSubscriptionHandler(service SubscriptionService) *SubscriptionHandler {
	return &SubscriptionHandler{service: service}
}

func (h *SubscriptionHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "POST":
		var sub Subscription
		if err := json.NewDecoder(r.Body).Decode(&sub); err != nil {
			http.Error(w, "invalid request payload", http.StatusBadRequest)
			return
		}

		newSub, err := h.service.Subscribe(sub.UserID, sub.EmployeeID)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(newSub)
	case "DELETE":
		params := mux.Vars(r)
		userID, err := strconv.Atoi(params["userID"])
		if err != nil {
			http.Error(w, "invalid user ID", http.StatusBadRequest)
			return
		}

		employeeID, err := strconv.Atoi(params["employeeID"])
		if err != nil {
			http.Error(w, "invalid employee ID", http.StatusBadRequest)
			return
		}

		if err := h.service.Unsubscribe(userID, employeeID); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusNoContent)
	case "GET":
		params := mux.Vars(r)
		userID, err := strconv.Atoi(params["userID"])
		if err != nil {
			http.Error(w, "invalid user ID", http.StatusBadRequest)
			return
		}

		subs, err := h.service.ListSubscriptions(userID)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(subs)
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}
