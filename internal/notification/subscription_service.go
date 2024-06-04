package notification

import "errors"

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

type subscriptionService struct {
	subscriptions []Subscription
	nextID        int
}

func NewSubscriptionService() SubscriptionService {
	return &subscriptionService{
		subscriptions: []Subscription{},
		nextID:        1,
	}
}

func (s *subscriptionService) Subscribe(userID, employeeID int) (Subscription, error) {
	for _, sub := range s.subscriptions {
		if sub.UserID == userID && sub.EmployeeID == employeeID {
			return Subscription{}, errors.New("already subscribed")
		}
	}

	subscription := Subscription{
		ID:         s.nextID,
		UserID:     userID,
		EmployeeID: employeeID,
	}
	s.nextID++
	s.subscriptions = append(s.subscriptions, subscription)
	return subscription, nil
}

func (s *subscriptionService) Unsubscribe(userID, employeeID int) error {
	for i, sub := range s.subscriptions {
		if sub.UserID == userID && sub.EmployeeID == employeeID {
			s.subscriptions = append(s.subscriptions[:i], s.subscriptions[i+1:]...)
			return nil
		}
	}
	return errors.New("subscription not found")
}

func (s *subscriptionService) ListSubscriptions(userID int) ([]Subscription, error) {
	var result []Subscription
	for _, sub := range s.subscriptions {
		if sub.UserID == userID {
			result = append(result, sub)
		}
	}
	return result, nil
}
