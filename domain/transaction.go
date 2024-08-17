package domain

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

var (
	ErrInsufficientBalance             = errors.New("insufficient balance for the transaction")
	ErrTransferNotAuthorized           = errors.New("authorization service not authorized this transfer")
	ErrTransferNotAllowedForWalletType = errors.New("this wallet type is not allowed to transfer")
)

type Transaction struct {
	ID        uuid.UUID `gorm:"column:id;type:char(36);primaryKey"`
	PayerID   uuid.UUID `gorm:"column:payerId;type:char(36);not null;index"`
	PayeeID   uuid.UUID `gorm:"column:payeeId;type:char(36);not null;index"`
	Payer     User      `gorm:"foreignKey:PayerID"`
	Payee     User      `gorm:"foreignKey:PayeeID"`
	Value     float64   `gorm:"column:value;type:decimal(15, 2);not null"`
	CreatedAt time.Time `gorm:"column:createdAt;index"`
	UpdatedAt time.Time `gorm:"column:updatedAt"`
}

func (Transaction) TableName() string {
	return "Transaction"
}

type TransactionPayload struct {
	PayeeID uuid.UUID `json:"payeeId" validate:"required,uuid"`
	Value   float64   `json:"value" validate:"required,gt=0"`
}

type TransactionHandler interface {
	Create(echo.Context) error
}

type TransactionService interface {
	Create(ctx context.Context, payload *TransactionPayload) error
}

type TransactionRepository interface {
	Create(ctx context.Context, transaction *Transaction) error
}

func (t *TransactionPayload) Validate() map[string]string {
	return ValidateStruct(t)
}
