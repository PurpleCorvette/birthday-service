package auth

import (
	"errors"

	"golang.org/x/crypto/bcrypt"
)

type User struct {
	ID       int    `json:"id"`
	Username string `json:"username"`
	Password string `json:"password"`
}

type AuthService interface {
	Register(username, password string) (User, error)
	Login(username, password string) (User, error)
	GetUsers() ([]User, error)
	GetUser(userID int) (User, error)
}

type authService struct {
	users  []User
	nextID int
}

func NewAuthService() AuthService {
	return &authService{
		users:  []User{},
		nextID: 1,
	}
}

// Register хеширует пароль
func (s *authService) Register(username, password string) (User, error) {
	for _, user := range s.users {
		if user.Username == username {
			return User{}, errors.New("user already exists")
		}
	}

	hashedPass, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return User{}, nil
	}

	user := User{
		ID:       s.nextID,
		Username: username,
		Password: string(hashedPass),
	}
	s.nextID++
	s.users = append(s.users, user)
	return user, nil
}

func (s *authService) Login(username, password string) (User, error) {
	for _, user := range s.users {
		if user.Username == username {
			err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
			if err != nil {
				return User{}, nil
			}
			return user, nil
		}
	}
	return User{}, errors.New("invalid username or password")
}

func (s *authService) GetUsers() ([]User, error) {
	return s.users, nil
}

func (s *authService) GetUser(userID int) (User, error) {
	for _, user := range s.users {
		if user.ID == userID {
			return user, nil
		}
	}
	return User{}, errors.New("user not found")
}
