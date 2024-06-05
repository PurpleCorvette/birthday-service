package teltegram

import (
	"os"
	"strconv"
	"strings"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"

	"birthday-service/internal/auth"
	"birthday-service/internal/employee"
	"birthday-service/internal/notification"
)

type Bot struct {
	bot                 *tgbotapi.BotAPI
	authService         auth.AuthService
	employeeService     employee.EmployeeService
	subscriptionService notification.SubscriptionService
}

func NewBot(authService auth.AuthService, employeeService employee.EmployeeService, subscriptionService notification.SubscriptionService) (*Bot, error) {
	token := os.Getenv("TELEGRAM_BOT_TOKEN")
	bot, err := tgbotapi.NewBotAPI(token)
	if err != nil {
		return nil, err
	}

	return &Bot{
		bot:                 bot,
		authService:         authService,
		employeeService:     employeeService,
		subscriptionService: subscriptionService,
	}, nil
}

func (b *Bot) Start() {
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates := b.bot.GetUpdatesChan(u)

	for update := range updates {
		if update.Message != nil {
			b.handleMessage(update.Message)
		}
	}
}

func (b *Bot) handleMessage(msg *tgbotapi.Message) {
	if msg.Text == "/start" {
		b.handleStartCommand(msg)
	} else if strings.HasPrefix(msg.Text, "/register") {
		b.handleRegisterCommand(msg)
	} else if strings.HasPrefix(msg.Text, "/subscribe") {
		b.handleSubscribeCommand(msg)
	}
}

func (b *Bot) handleStartCommand(msg *tgbotapi.Message) {
	reply := "Welcome! Use /register <username> <password> to register. Use /subscribe <employee_id> to subscribe to birthday notifications."
	message := tgbotapi.NewMessage(msg.Chat.ID, reply)
	b.bot.Send(message)
}

func (b *Bot) handleRegisterCommand(msg *tgbotapi.Message) {
	parts := strings.Split(msg.Text, " ")
	if len(parts) != 3 {
		reply := "Usage: /register <username> <password>"
		message := tgbotapi.NewMessage(msg.Chat.ID, reply)
		b.bot.Send(message)
		return
	}

	username := parts[1]
	password := parts[2]
	telegramID := msg.From.ID

	user, err := b.authService.Register(username, password, telegramID)
	if err != nil {
		reply := "Error registering user: " + err.Error()
		message := tgbotapi.NewMessage(msg.Chat.ID, reply)
		b.bot.Send(message)
		return
	}

	reply := "User registered successfully! Username: " + user.Username
	message := tgbotapi.NewMessage(msg.Chat.ID, reply)
	b.bot.Send(message)
}

func (b *Bot) handleSubscribeCommand(msg *tgbotapi.Message) {
	parts := strings.Split(msg.Text, " ")
	if len(parts) != 2 {
		reply := "Usage: /subscribe <employee_id>"
		message := tgbotapi.NewMessage(msg.Chat.ID, reply)
		b.bot.Send(message)
		return
	}

	employeeID, err := strconv.Atoi(parts[1])
	if err != nil {
		reply := "Invalid employee ID"
		message := tgbotapi.NewMessage(msg.Chat.ID, reply)
		b.bot.Send(message)
		return
	}

	user, err := b.authService.GetUserByTelegramID(msg.From.ID)
	if err != nil {
		reply := "User not registered"
		message := tgbotapi.NewMessage(msg.Chat.ID, reply)
		b.bot.Send(message)
		return
	}

	err = b.subscriptionService.Subscribe(user.ID, employeeID)
	if err != nil {
		reply := "Error subscribing to employee: " + err.Error()
		message := tgbotapi.NewMessage(msg.Chat.ID, reply)
		b.bot.Send(message)
		return
	}

	reply := "Subscribed to employee's birthday notifications successfully!"
	message := tgbotapi.NewMessage(msg.Chat.ID, reply)
	b.bot.Send(message)
}
