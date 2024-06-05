package auth

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v4"
	"golang.org/x/crypto/bcrypt"

	"birthday-service/internal/db"
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

type authService struct{}

func NewAuthService() AuthService {
	return &authService{}
}

func (s *authService) Register(username, password string) (User, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return User{}, err
	}

	var user User
	err = db.Conn.QueryRow(context.Background(),
		"INSERT INTO users (username, password) VALUES ($1, $2) RETURNING id, username",
		username, string(hashedPassword)).Scan(&user.ID, &user.Username)
	if err != nil {
		return User{}, err
	}

	return user, nil
}

func (s *authService) Login(username, password string) (User, error) {
	var user User
	err := db.Conn.QueryRow(context.Background(),
		"SELECT id, username, password FROM users WHERE username=$1", username).Scan(&user.ID, &user.Username, &user.Password)
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
	err := db.Conn.QueryRow(context.Background(),
		"SELECT id, username FROM users WHERE id=$1", user).Scan(&user.ID, &user.Username)
	if err == pgx.ErrNoRows {
		return User{}, errors.New("user not found")
	}
	if err != nil {
		return User{}, err
	}

	return user, nil
}

func (s *authService) GetUsers() ([]User, error) {
	rows, err := db.Conn.Query(context.Background(),
		"SELECT id, username FROM users")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []User
	for rows.Next() {
		var user User
		err := rows.Scan(&user.ID, &user.Username)
		if err != nil {
			return nil, err
		}
		users = append(users, user)
	}

	return users, nil
}
