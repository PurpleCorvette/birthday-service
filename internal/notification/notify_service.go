package notification

import (
	"fmt"
	"log"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"

	"birthday-service/internal/auth"
	"birthday-service/internal/employee"
)

type NotifyService struct {
	subscriptionService SubscriptionService
	employeeService     employee.EmployeeService
	authService         auth.AuthService
	bot                 *tgbotapi.BotAPI
	groupChatID         int64 // Добавьте идентификатор группы, если необходимо
}

func NewNotifyService(subscriptionService SubscriptionService, employeeService employee.EmployeeService, authService auth.AuthService, bot *tgbotapi.BotAPI, groupChatID int64) *NotifyService {
	return &NotifyService{
		subscriptionService: subscriptionService,
		employeeService:     employeeService,
		authService:         authService,
		bot:                 bot,
		groupChatID:         groupChatID,
	}
}

func (n *NotifyService) Notify(userID, employeeID int, message string) error {
	user, err := n.authService.GetUser(userID)
	if err != nil {
		return err
	}

	chatID := user.TelegramID
	if chatID == 0 {
		return fmt.Errorf("user %d does not have a Telegram ID", userID)
	}

	msg := tgbotapi.NewMessage(chatID, message)
	_, err = n.bot.Send(msg)
	if err != nil {
		log.Printf("failed to send message to user %d: %v", userID, err)
		return err
	}

	return nil
}

func (n *NotifyService) NotifyGroup(employeeName string) error {
	message := fmt.Sprintf("Сегодня день рождения у %s!", employeeName)
	msg := tgbotapi.NewMessage(n.groupChatID, message)
	_, err := n.bot.Send(msg)
	if err != nil {
		log.Printf("failed to send message to group: %v", err)
		return err
	}

	return nil
}

func (n *NotifyService) ScheduleDailyNotifications() error {
	employees, err := n.employeeService.GetAllEmployees()
	if err != nil {
		return err
	}

	today := time.Now().Format("2006-01-02")
	for _, employee := range employees {
		if employee.Birthday.Format("2006-01-02") == today {
			subscriptions, err := n.subscriptionService.GetSubscriptions(employee.ID)
			if err != nil {
				log.Printf("failed to get subscriptions for employee %d: %v", employee.ID, err)
				continue
			}

			message := fmt.Sprintf("Сегодня день рождения у %s!", employee.Name)
			for _, subscription := range subscriptions {
				err = n.Notify(subscription.UserID, employee.ID, message)
				if err != nil {
					log.Printf("failed to notify user %d: %v", subscription.UserID, err)
				}
			}

			// Отправка уведомления в группу
			err = n.NotifyGroup(employee.Name)
			if err != nil {
				log.Printf("failed to notify group: %v", err)
			}
		}
	}

	return nil
}
