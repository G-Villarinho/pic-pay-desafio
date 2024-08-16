package domain

import (
	"strings"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
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

type UserPayLoad struct {
	Name            string `json:"name" validate:"required,min=1,max=75"`
	CPF             string `json:"cpf" validate:"required,cpf"`
	Email           string `json:"email" validate:"required,email"`
	ConfirmEmail    string `json:"confirmEmail" validate:"required,email,eqfield=Email"`
	Password        string `json:"password,omitempty" validate:"required,max=255,strongpassword"`
	ConfirmPassword string `json:"confirmPassword" validate:"required,eqfield=Password"`
}

func (u *UserPayLoad) trim() {
	u.Name = strings.TrimSpace(u.Name)
	u.CPF = strings.TrimSpace(u.CPF)
	u.Email = strings.TrimSpace(u.Email)
	u.ConfirmEmail = strings.TrimSpace(u.ConfirmEmail)
}

func (u *UserPayLoad) Validate() error {
	u.trim()
	validate := validator.New()
	return validate.Struct(u)
}

func (User) TableName() string {
	return "User"
}
