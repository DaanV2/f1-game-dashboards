package users

import (
	"context"
	"errors"

	"github.com/charmbracelet/log"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type User struct {
	Id       string `json:"id"`
	Email    string `json:"email"`
	Password string `json:"password"` // Hashed
	Admin    bool   `json:"admin,omitempty"`
	Guest    bool   `json:"guest,omitempty"`
}

type UserStorage interface {
	GetByEmail(email string) (*User, error)
	Set(value *User) error
}

type UserManagement struct {
	db UserStorage
}

func NewUserManagement(db UserStorage) *UserManagement {
	return &UserManagement{
		db: db,
	}
}

// Authenticate checks if the user is authenticated
func (um *UserManagement) Authenticate(ctx context.Context, email, password string) (*User, error) {
	logger := log.FromContext(ctx).With("email", email)
	logger.Debug("checking if user is authenticated")

	user, err := um.db.GetByEmail(email)
	if err != nil {
		logger.Warn("user not found")
		return user, errors.New("invalid password / username")
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		logger.Warn("invalid password")
		return user, errors.New("invalid password")
	}

	logger.Debug("user authenticated")
	return user, nil
}

// UpdatePassword updates the password of the user
func (um *UserManagement) UpdatePassword(ctx context.Context, email, password string) error {
	logger := log.FromContext(ctx).With("email", email)
	logger.Debug("user authenticated")

	user, err := um.db.GetByEmail(email)
	if err != nil {
		return err
	}

	user.Password = password
	if err := hashPassword(user); err != nil {
		return err
	}

	return um.db.Set(user)
}

// Create creates a new user
func (um *UserManagement) Create(email, password string, admin bool) error {
	if _, err := um.db.GetByEmail(email); err != nil {
		return errors.New("user already exists")
	}

	user := &User{
		Id:       uuid.New().String(),
		Email:    email,
		Password: password,
		Admin:    admin,
	}
	if err := hashPassword(user); err != nil {
		return err
	}

	return um.db.Set(user)
}

// GetByEmail returns a user by email
func (um *UserManagement) GetByEmail(email string) (*User, error) {
	return um.db.GetByEmail(email)
}

func hashPassword(user *User) error {
	data, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	user.Password = string(data)
	return nil
}
