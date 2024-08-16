package domain

import (
	"context"
	"strings"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
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

type UserPayload struct {
	Name            string `json:"name" validate:"required,min=1,max=75"`
	CPF             string `json:"cpf" validate:"required,cpf"`
	Email           string `json:"email" validate:"required,email"`
	ConfirmEmail    string `json:"confirmEmail" validate:"required,email,eqfield=Email"`
	Password        string `json:"password,omitempty" validate:"required,max=255,strongpassword"`
	ConfirmPassword string `json:"confirmPassword" validate:"required,eqfield=Password"`
}

type UserHandler interface {
	Create(ctx echo.Context) error
}

type UserService interface {
	Create(ctx context.Context, payload *UserPayload)
}

type UserRepository interface {
	Create(ctx context.Context, user *User) error
}

func (u *UserPayload) trim() {
	u.Name = strings.TrimSpace(u.Name)
	u.CPF = strings.TrimSpace(u.CPF)
	u.Email = strings.TrimSpace(u.Email)
	u.ConfirmEmail = strings.TrimSpace(u.ConfirmEmail)
}

func (u *UserPayload) Validate() map[string]string {
	u.trim()
	validate := validator.New()
	err := validate.Struct(u)

	validationErrors := make(map[string]string)

	if err != nil {
		for _, err := range err.(validator.ValidationErrors) {
			fieldName := strings.ToLower(err.Field())
			validationErrors[fieldName] = getErrorMessage(err)
		}
	}

	if len(validationErrors) > 0 {
		return validationErrors
	}

	return nil
}

func getErrorMessage(err validator.FieldError) string {
	switch err.Tag() {
	case "required":
		return "This field is required"
	case "email":
		return "Invalid email format"
	case "min":
		return "Value is too short"
	case "max":
		return "Value is too long"
	case "eqfield":
		return "Fields do not match"
	default:
		return "Invalid value"
	}
}

func (User) TableName() string {
	return "User"
}
