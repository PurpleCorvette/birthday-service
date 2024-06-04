package main

import (
	"birthday-service/internal/api"
	"net/http"

	"github.com/gorilla/mux"

	"birthday-service/internal/auth"
	"birthday-service/internal/employee"
	"birthday-service/internal/notification"
	"birthday-service/pkg/logging"
	"github.com/joho/godotenv"
	"github.com/robfig/cron/v3"
)

func main() {
	log := logging.GetLogger()

	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading .env file")
	}

	r := mux.NewRouter()

	authService := auth.NewAuthService()
	authHandler := auth.NewAuthHandler(authService)

	employeeService := employee.NewEmployeeService()
	employeeHandler := employee.NewEmployeeHandler(employeeService)

	notificationService := notification.NewNotificationService()
	notificationHandler := notification.NewNotificationHandler(notificationService)

	subscriptionService := notification.NewSubscriptionService()
	settingsService := notification.NewSettingsService()
	subscriptionHandler := api.NewSubscriptionHandler(subscriptionService, settingsService)

	notifyService := notification.NewNotifyService(subscriptionService, employeeService, authService)

	r.Handle("/auth", authHandler).Methods("POST", "GET")
	r.Handle("/employee/{id:[0-9]+}", employeeHandler).Methods("GET", "PUT", "DELETE")
	r.Handle("/employee", employeeHandler).Methods("POST")
	r.Handle("/notification/{id:[0-9]+}", notificationHandler).Methods("GET", "DELETE")
	r.Handle("/notification", notificationHandler).Methods("POST")
	r.HandleFunc("/subscription/{userID:[0-9]+}/{employeeID:[0-9]+}", subscriptionHandler.DeleteSubscription).Methods("DELETE")
	r.HandleFunc("/subscription", subscriptionHandler.CreateSubscription).Methods("POST")
	r.HandleFunc("/subscriptions/{userID:[0-9]+}", subscriptionHandler.GetSubscriptions).Methods("GET")
	r.HandleFunc("/api/notifications/settings", subscriptionHandler.UpdateNotificationSettings).Methods("POST")

	c := cron.New()
	c.AddFunc("@daily", func() {
		err := notifyService.ScheduleDailyNotifications()
		if err != nil {
			log.Errorf("Error scheduling daily notifications: %v", err)
		}
	})
	c.Start()

	log.Infoln("Starting server on :8080")
	err = http.ListenAndServe(":8080", r)
	if err != nil {
		log.Fatalf("Could not start server: %s", err.Error())
	}
}
