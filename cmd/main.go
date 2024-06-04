package main

import (
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
	subscriptionHandler := notification.NewSubscriptionHandler(subscriptionService)

	notyfyService := notification.NewNotifyService(subscriptionService, employeeService, authService)

	r.Handle("/auth", authHandler).Methods("POST", "GET")
	r.Handle("/employee/{id:[0-9]+}", employeeHandler).Methods("GET", "PUT", "DELETE")
	r.Handle("/employee", employeeHandler).Methods("POST")
	r.Handle("/notification/{id:[0-9]+}", notificationHandler).Methods("GET", "DELETE")
	r.Handle("/notification", notificationHandler).Methods("POST")
	r.Handle("/subscription/{userID:[0-9]+}/{employeeID:[0-9]+}", subscriptionHandler).Methods("DELETE")
	r.Handle("/subscription", subscriptionHandler).Methods("POST")
	r.Handle("/subscriptions/{userID:[0-9]+}", subscriptionHandler).Methods("GET")

	c := cron.New()
	c.AddFunc("@daily", func() {
		err := notyfyService.ScheduleDailyNotifications()
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
