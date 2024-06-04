package notification

import (
	"errors"
	"fmt"
	"os"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"

	"birthday-service/internal/auth"
	"birthday-service/internal/employee"
)

type NotifyService interface {
	Notify(userID, employeeID int, message string) error
	ScheduleDailyNotifications() error
}

type notifyService struct {
	subscriptionService SubscriptionService
	employeeService     employee.EmployeeService
	userService         auth.AuthService
}

func NewNotifyService(subServ SubscriptionService, empServ employee.EmployeeService, userServ auth.AuthService) NotifyService {
	return &notifyService{
		subscriptionService: subServ,
		employeeService:     empServ,
		userService:         userServ,
	}
}

func (n *notifyService) Notify(userID, employeeID int, message string) error {
	botToken := os.Getenv("TELEGRAM_BOT_TOKEN")
	bot, err := tgbotapi.NewBotAPI(botToken)
	if err != nil {
		return err
	}

	user, err := n.userService.GetUser(userID)
	if err != nil {
		return err
	}

	employee, err := n.employeeService.GetEmployee(employeeID)
	if err != nil {
		return err
	}

	chatID := os.Getenv("TELEGRAM_CHAT_ID")
	if chatID == "" {
		return errors.New("TELEGRAM_CHAT_ID is not set")
	}

	personalizedMessage := fmt.Sprintf("Привет, %s! Сегодня день рождения у %s. Не забудьте поздравить!", user.Username, employee.Name)
	msg := tgbotapi.NewMessageToChannel(chatID, personalizedMessage)
	_, err = bot.Send(msg)
	if err != nil {
		return err
	}

	return nil
}

func (n *notifyService) ScheduleDailyNotifications() error {
	users, err := n.userService.GetUsers()
	if err != nil {
		return err
	}

	for _, user := range users {
		subscriptions, err := n.subscriptionService.ListSubscriptions(user.ID)
		if err != nil {
			return err
		}

		for _, sub := range subscriptions {
			empl, err := n.employeeService.GetEmployee(sub.EmployeeID)
			if err != nil {
				return err
			}

			// Check if today is the empl's birthday
			today := time.Now().Format("2006-01-02")
			if empl.DOB == today {
				message := fmt.Sprintf("Сегодня день рождения у %s. Не забудьте поздравить!", empl.Name)
				err = n.Notify(user.ID, empl.ID, message)
				if err != nil {
					return err
				}
			}
		}
	}
	return nil
}
