package notification

import "errors"

type Notification struct {
	ID         int    `json:"id"`
	UserID     int    `json:"user_id"`
	EmployeeID int    `json:"employee_id"`
	Message    string `json:"message"`
}

type NotificationService interface {
	CreateNotification(userID, employeeID int, message string) (Notification, error)
	GetNotification(id int) (Notification, error)
	DeleteNotification(id int) error
	ListNotifications() ([]Notification, error)
}

type notificationService struct {
	notifications []Notification
	nextID        int
}

func NewNotificationService() NotificationService {
	return &notificationService{
		notifications: []Notification{},
		nextID:        1,
	}
}

func (s *notificationService) CreateNotification(userID, employeeID int, message string) (Notification, error) {
	notification := Notification{
		ID:         s.nextID,
		UserID:     userID,
		EmployeeID: employeeID,
		Message:    message,
	}
	s.nextID++
	s.notifications = append(s.notifications, notification)
	return notification, nil
}

func (s *notificationService) GetNotification(id int) (Notification, error) {
	for _, not := range s.notifications {
		if not.ID == id {
			return not, nil
		}
	}
	return Notification{}, errors.New("notification not found")
}

func (s *notificationService) DeleteNotification(id int) error {
	for i, not := range s.notifications {
		if not.ID == id {
			s.notifications = append(s.notifications[:i], s.notifications[i+1:]...)
			return nil
		}
	}
	return errors.New("notification not found")
}

func (s *notificationService) ListNotifications() ([]Notification, error) {
	return s.notifications, nil
}
