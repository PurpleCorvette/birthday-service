package auth

import (
	"birthday-service/internal/db"
	"context"
	"errors"
	"golang.org/x/crypto/bcrypt"

	"github.com/jackc/pgx/v4"
)

type User struct {
	ID         int    `json:"id"`
	Username   string `json:"username"`
	Password   string `json:"password"`
	TelegramID int64  `json:"telegram_id"`
}

type AuthService interface {
	Register(username, password string, telegramID int64) (User, error)
	Login(username, password string) (User, error)
	GetUser(userID int) (User, error)
	GetUsers() ([]User, error)
	GetUserByTelegramID(telegramID int64) (User, error)
}

type authService struct {
	db db.DB
}

func NewAuthService(db db.DB) AuthService {
	return &authService{db: db}
}

func (s *authService) Register(username, password string, telegramID int64) (User, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return User{}, err
	}

	var user User
	err = s.db.QueryRow(context.Background(),
		"INSERT INTO users (username, password, telegram_id) VALUES ($1, $2, $3) RETURNING id, username, telegram_id",
		username, string(hashedPassword), telegramID).Scan(&user.ID, &user.Username, &user.TelegramID)
	if err != nil {
		return User{}, err
	}

	return user, nil
}

func (s *authService) Login(username, password string) (User, error) {
	var user User
	err := s.db.QueryRow(context.Background(),
		"SELECT id, username, password, telegram_id FROM users WHERE username=$1",
		username).Scan(&user.ID, &user.Username, &user.Password, &user.TelegramID)
	if err == pgx.ErrNoRows {
		return User{}, errors.New("invalid username or password")
	}
	if err != nil {
		return User{}, err
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		return User{}, errors.New("invalid username or password")
	}

	return user, nil
}

func (s *authService) GetUser(userID int) (User, error) {
	var user User
	err := s.db.QueryRow(context.Background(),
		"SELECT id, username, telegram_id FROM users WHERE id=$1", userID).Scan(&user.ID, &user.Username, &user.TelegramID)
	if err == pgx.ErrNoRows {
		return User{}, errors.New("user not found")
	}
	if err != nil {
		return User{}, err
	}

	return user, nil
}

func (s *authService) GetUsers() ([]User, error) {
	rows, err := s.db.Query(context.Background(), "SELECT id, username, telegram_id FROM users")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []User
	for rows.Next() {
		var user User
		err := rows.Scan(&user.ID, &user.Username, &user.TelegramID)
		if err != nil {
			return nil, err
		}
		users = append(users, user)
	}

	return users, nil
}

func (s *authService) GetUserByTelegramID(telegramID int64) (User, error) {
	var user User
	err := s.db.QueryRow(context.Background(),
		"SELECT id, username, telegram_id FROM users WHERE telegram_id=$1", telegramID).Scan(&user.ID, &user.Username, &user.TelegramID)
	if err == pgx.ErrNoRows {
		return User{}, errors.New("user not found")
	}
	if err != nil {
		return User{}, err
	}

	return user, nil
}
