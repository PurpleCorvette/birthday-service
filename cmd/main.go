package main

import (
	"net/http"
	"os"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	"github.com/robfig/cron/v3"

	"birthday-service/internal/api"
	"birthday-service/internal/auth"
	"birthday-service/internal/db"
	"birthday-service/internal/employee"
	"birthday-service/internal/notification"
	"birthday-service/internal/teltegram"
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

	r := mux.NewRouter()

	authService := auth.NewAuthService(database)
	authHandler := auth.NewAuthHandler(authService)

	employeeService := employee.NewEmployeeService(database)
	employeeHandler := employee.NewEmployeeHandler(employeeService)

	subscriptionService := notification.NewSubscriptionService(database)
	settingsService := notification.NewSettingsService(database)
	subscriptionHandler := api.NewSubscriptionHandler(subscriptionService, settingsService)

	notifyService := notification.NewNotifyService(subscriptionService, employeeService, authService, bot)

	telegramBot, err := teltegram.NewBot(authService, employeeService, subscriptionService)
	if err != nil {
		log.Fatalf("Error creating Telegram bot: %s", err)
	}
	go telegramBot.Start()

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

	c := cron.New()
	c.AddFunc("@daily", func() {
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
