package main

import (
	"net/http"
	"os"
	"strconv"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	"github.com/robfig/cron/v3"

	"birthday-service/internal/api"
	"birthday-service/internal/auth"
	"birthday-service/internal/db"
	"birthday-service/internal/employee"
	"birthday-service/internal/notification"
	"birthday-service/pkg/logging"
)

func main() {
	log := logging.GetLogger()

	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading .env file")
	}

	// Connect to the database
	dbURL := os.Getenv("DATABASE_URL")
	database, err := db.ConnectDatabase(dbURL)
	if err != nil {
		log.Fatalf("Could not connect to database: %s", err)
	}
	defer database.Close()

	// Create Telegram bot
	botToken := os.Getenv("TELEGRAM_BOT_TOKEN")
	bot, err := tgbotapi.NewBotAPI(botToken)
	if err != nil {
		log.Fatalf("Error creating Telegram bot: %s", err)
	}

	groupChatIDStr := os.Getenv("TELEGRAM_CHAT_ID")
	groupChatID, err := strconv.ParseInt(groupChatIDStr, 10, 64)
	if err != nil {
		log.Fatalf("Error parsing GROUP_CHAT_ID: %v", err)
	}

	r := mux.NewRouter()

	authService := auth.NewAuthService(database)
	employeeService := employee.NewEmployeeService(database)
	subscriptionService := notification.NewSubscriptionService(database)
	settingsService := notification.NewSettingsService(database)
	notifyService := notification.NewNotifyService(subscriptionService, employeeService, authService, bot, groupChatID)
	authHandler := auth.NewAuthHandler(authService)
	employeeHandler := employee.NewEmployeeHandler(employeeService)
	subscriptionHandler := api.NewSubscriptionHandler(subscriptionService, settingsService)
	notificationHandler := api.NewNotificationHandler(*notifyService)

	r.Handle("/auth", authHandler).Methods("POST", "GET")
	r.HandleFunc("/employee/{id:[0-9]+}", employeeHandler.GetEmployee).Methods("GET")
	r.HandleFunc("/employee", employeeHandler.AddEmployee).Methods("POST")
	r.HandleFunc("/employees", employeeHandler.GetAllEmployees).Methods("GET")
	r.HandleFunc("/employee/{id:[0-9]+}", employeeHandler.UpdateEmployee).Methods("PUT")
	r.HandleFunc("/employee/{id:[0-9]+}", employeeHandler.DeleteEmployee).Methods("DELETE")
	r.HandleFunc("/subscription/{userID:[0-9]+}/{employeeID:[0-9]+}", subscriptionHandler.DeleteSubscription).Methods("DELETE")
	r.HandleFunc("/subscription", subscriptionHandler.CreateSubscription).Methods("POST")
	r.HandleFunc("/subscriptions/{userID:[0-9]+}", subscriptionHandler.GetSubscriptions).Methods("GET")
	r.HandleFunc("/api/notifications/settings", subscriptionHandler.UpdateNotificationSettings).Methods("POST")
	r.HandleFunc("/trigger-notifications", notificationHandler.TriggerNotifications).Methods("POST")

	c := cron.New()
	c.AddFunc("@hourly", func() {
		err := notifyService.ScheduleDailyNotifications()
		if err != nil {
			log.Printf("Error scheduling daily notifications: %v", err)
		}
	})
	c.Start()

	log.Infoln("Starting server on :8080")
	err = http.ListenAndServe(":8080", r)
	if err != nil {
		log.Fatalf("Could not start server: %s", err.Error())
	}
}
