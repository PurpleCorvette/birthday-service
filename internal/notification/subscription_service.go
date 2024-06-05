package notification

import (
	"context"

	"birthday-service/internal/db"
)

type Subscription struct {
	UserID     int `json:"user_id"`
	EmployeeID int `json:"employee_id"`
}

type SubscriptionService interface {
	Subscribe(userID, employeeID int) error
	Unsubscribe(userID, employeeID int) error
	GetSubscriptions(userID int) ([]Subscription, error)
}

type subscriptionService struct {
	db db.DB
}

func NewSubscriptionService(db db.DB) SubscriptionService {
	return &subscriptionService{db: db}
}

func (s *subscriptionService) Subscribe(userID, employeeID int) error {
	_, err := s.db.Exec(context.Background(),
		"INSERT INTO subscriptions (user_id, employee_id) VALUES ($1, $2)",
		userID, employeeID)
	return err
}

func (s *subscriptionService) Unsubscribe(userID, employeeID int) error {
	_, err := s.db.Exec(context.Background(),
		"DELETE FROM subscriptions WHERE user_id=$1 AND employee_id=$2",
		userID, employeeID)
	return err
}

func (s *subscriptionService) GetSubscriptions(userID int) ([]Subscription, error) {
	rows, err := s.db.Query(context.Background(),
		"SELECT user_id, employee_id FROM subscriptions WHERE user_id=$1",
		userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var subscriptions []Subscription
	for rows.Next() {
		var subscription Subscription
		err := rows.Scan(&subscription.UserID, &subscription.EmployeeID)
		if err != nil {
			return nil, err
		}
		subscriptions = append(subscriptions, subscription)
	}

	return subscriptions, nil
}
