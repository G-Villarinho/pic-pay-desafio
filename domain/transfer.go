package domain

//go:generate mockgen -source=transfer.go -destination=../mocks/transfer_mock.go -package=mocks

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

var (
	ErrInsufficientBalance             = errors.New("insufficient balance for the transfer")
	ErrTransferNotAuthorized           = errors.New("authorization service not authorized this transfer")
	ErrTransferNotAllowedForWalletType = errors.New("this wallet type is not allowed to transfer")
	ErrCreateTransfer                  = errors.New("fail to create transfer")
)

type Transfer struct {
	ID        uuid.UUID      `gorm:"column:id;type:char(36);primaryKey"`
	PayerID   uuid.UUID      `gorm:"column:payerId;type:char(36);not null;index"`
	PayeeID   uuid.UUID      `gorm:"column:payeeId;type:char(36);not null;index"`
	Payer     User           `gorm:"foreignKey:PayerID"`
	Payee     User           `gorm:"foreignKey:PayeeID"`
	Value     float64        `gorm:"column:value;type:decimal(15, 2);not null"`
	CreatedAt time.Time      `gorm:"column:createdAt;not null"`
	UpdatedAt time.Time      `gorm:"column:updatedAt;default:NULL"`
	DeletedAt gorm.DeletedAt `gorm:"column:deletedAt;index"`
}

func (Transfer) TableName() string {
	return "Transfer"
}

func (t *Transfer) BeforeUpdate(tx *gorm.DB) (err error) {
	t.UpdatedAt = time.Now().UTC()
	return nil
}

type TransferPayload struct {
	PayeeID uuid.UUID `json:"payeeId" validate:"required,uuid"`
	Value   float64   `json:"value" validate:"required,gt=0"`
}

type TransferHandler interface {
	Transfer(ctx echo.Context) error
}

type TransferService interface {
	Transfer(ctx context.Context, payload *TransferPayload) error
}

type TransferRepository interface {
	Transfer(ctx context.Context, transfer *Transfer) error
}

func (t *TransferPayload) Validate() map[string]string {
	return ValidateStruct(t)
}

func (t *TransferPayload) ToTansaction(payerID uuid.UUID) *Transfer {
	return &Transfer{
		ID:        uuid.New(),
		PayerID:   payerID,
		PayeeID:   t.PayeeID,
		Value:     t.Value,
		CreatedAt: time.Now().UTC(),
	}
}
