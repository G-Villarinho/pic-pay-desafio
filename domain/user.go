package domain

import (
	"context"
	"errors"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/klassmann/cpfcnpj"
	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

var (
	ErrUserNotFound          = errors.New("user not found")
	ErrEmailAlreadyRegister  = errors.New("email already exists")
	ErrCPFAlreadyRegister    = errors.New("cpf already exists")
	ErrHashingPassword       = errors.New("failed to hash password")
	ErrInvalidPassword       = errors.New("invalid password")
	ErrUserNotFoundInContext = errors.New("user not found in context")
)

type User struct {
	ID           uuid.UUID      `gorm:"column:id;type:char(36);primaryKey"`
	Name         string         `gorm:"column:name;type:varchar(255);not null;index"`
	CPF          string         `gorm:"column:cpf;type:char(11);uniqueIndex;not null"`
	Email        string         `gorm:"column:email;type:varchar(255);uniqueIndex;not null"`
	PasswordHash string         `gorm:"column:passwordHash;type:varchar(255);not null"`
	CreatedAt    time.Time      `gorm:"column:createdAt;index"`
	UpdatedAt    time.Time      `gorm:"column:updatedAt"`
	DeletedAt    gorm.DeletedAt `gorm:"column:deletedAt;index"`
}

func (User) TableName() string {
	return "User"
}

type UserPayload struct {
	Name            string `json:"name" validate:"required,min=1,max=75"`
	CPF             string `json:"cpf" validate:"required,cpf"`
	Email           string `json:"email" validate:"required,email"`
	ConfirmEmail    string `json:"confirmEmail" validate:"required,eqfield=Email"`
	Password        string `json:"password,omitempty" validate:"required,max=255,strongpassword"`
	ConfirmPassword string `json:"confirmPassword" validate:"required,eqfield=Password"`
}

type SignInPayload struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password,omitempty" validate:"required"`
}

type SignInResponse struct {
	Token string `json:"token"`
}

type UserHandler interface {
	Create(ctx echo.Context) error
	SignIn(ctx echo.Context) error
}

type UserService interface {
	Create(ctx context.Context, payload *UserPayload) error
	SignIn(ctx context.Context, payload *SignInPayload) (*SignInResponse, error)
}

type UserRepository interface {
	Create(ctx context.Context, user *User) error
	GetByID(ctx context.Context, ID uuid.UUID) (*User, error)
	GetByEmail(ctx context.Context, email string) (*User, error)
	GetByCPF(ctx context.Context, CPF string) (*User, error)
}

func (u *UserPayload) trim() {
	u.Name = strings.TrimSpace(u.Name)
	u.CPF = strings.TrimSpace(u.CPF)
	u.Email = strings.TrimSpace(strings.ToLower(u.Email))
	u.ConfirmEmail = strings.TrimSpace(strings.ToLower(u.ConfirmEmail))
}

func (s *SignInPayload) trim() {
	s.Email = strings.TrimSpace(strings.ToLower(s.Email))
}

func (u *UserPayload) Validate() map[string]string {
	u.trim()
	return ValidateStruct(u)
}

func (s *SignInPayload) Validate() map[string]string {
	s.trim()
	return ValidateStruct(s)
}

func (u *UserPayload) ToUser(passwordHash string) *User {
	return &User{
		ID:           uuid.New(),
		Name:         u.Name,
		CPF:          string(cpfcnpj.NewCPF(u.CPF)),
		Email:        u.Email,
		PasswordHash: passwordHash,
		CreatedAt:    time.Now().UTC(),
	}
}
