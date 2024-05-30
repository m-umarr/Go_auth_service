package repository

import (
	model "github.com/m-umarr/Go_auth_service/auth_service/internal/model"

	"database/sql"
	"errors"
)

type UserRepository struct {
	db *sql.DB
}

func NewUserRepository(db *sql.DB) *UserRepository {
	return &UserRepository{db: db}
}

func (r *UserRepository) CreateUser(phoneNumber string) error {
	var exists bool
	err := r.db.QueryRow("SELECT EXISTS(SELECT 1 FROM users WHERE phone_number = $1)", phoneNumber).Scan(&exists)
	if err != nil {
		return err
	}
	if exists {
		return errors.New("user already exists")
	}

	_, err = r.db.Exec("INSERT INTO users (phone_number) VALUES ($1)", phoneNumber)
	return err
}

func (r *UserRepository) VerifyUser(phoneNumber string) error {
	_, err := r.db.Exec("UPDATE users SET verified = TRUE WHERE phone_number = $1", phoneNumber)
	return err
}

func (r *UserRepository) GetUserProfile(phoneNumber string) (*model.User, error) {
	var profile model.User
	err := r.db.QueryRow("SELECT phone_number, profile_data FROM users WHERE phone_number = $1", phoneNumber).
		Scan(&profile.PhoneNumber, &profile.ProfileData)
	if err != nil {
		return nil, err
	}
	return &profile, nil
}

func (r *UserRepository) LogEvent(phoneNumber string, event string) error {
	_, err := r.db.Exec("INSERT INTO events (phone_number, event, event_time) VALUES ($1, $2, NOW())", phoneNumber, event)
	return err
}
