package notification

import (
	"context"

	"birthday-service/internal/db"
)

type Subscription struct {
	ID         int `json:"id"`
	UserID     int `json:"user_id"`
	EmployeeID int `json:"employee_id"`
}

type SubscriptionService interface {
	Subscribe(userID, employeeID int) (Subscription, error)
	Unsubscribe(userID, employeeID int) error
	ListSubscriptions(userID int) ([]Subscription, error)
}

type subscriptionService struct{}

func NewSubscriptionService() SubscriptionService {
	return &subscriptionService{}
}

func (s *subscriptionService) Subscribe(userID, employeeID int) (Subscription, error) {
	var sub Subscription
	err := db.Conn.QueryRow(context.Background(),
		"INSERT INTO subscriptions (user_id, employee_id) VALUES ($1, $2) RETURNING id, user_id, employee_id",
		userID, employeeID).Scan(&sub.ID, &sub.UserID, &sub.EmployeeID)
	if err != nil {
		return Subscription{}, err
	}

	return sub, nil
}

func (s *subscriptionService) Unsubscribe(userID, employeeID int) error {
	_, err := db.Conn.Exec(context.Background(), "DELETE FROM subscriptions WHERE user_id=$1 AND employee_id=$2",
		userID, employeeID)
	if err != nil {
		return err
	}

	return nil
}

func (s *subscriptionService) ListSubscriptions(userID int) ([]Subscription, error) {
	rows, err := db.Conn.Query(context.Background(), "SELECT id, user_id, employee_id FROM subscriptions WHERE user_id=$1", userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var subs []Subscription
	for rows.Next() {
		var sub Subscription
		err := rows.Scan(&sub.ID, &sub.UserID, &sub.EmployeeID)
		if err != nil {
			return nil, err
		}
		subs = append(subs, sub)
	}

	return subs, nil
}
