package domain

//go:generate mockgen -source=wallet.go -destination=../mocks/wallet_mock.go -package=mocks

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

var (
	ErrGetWallet                 = errors.New("error when trying to obtain wallet")
	ErrPayeeWalletNotFound       = errors.New("payee's wallet not found")
	ErrPayerWalletNotFound       = errors.New("payer's wallet not found")
	ErrSelfTransactionNotAllowed = errors.New("payer cannot perform transfers to themselves")
	ErrWalletAlredyRegister      = errors.New("the user already has a wallet")
	ErrDebitWallet               = errors.New("failed to debit the wallet")
	ErrCreditWallet              = errors.New("failed to credit the wallet")
)

type Wallet struct {
	UserID    uuid.UUID      `gorm:"column:userId;type:char(36);primaryKey"`
	User      User           `gorm:"foreignKey:UserID"`
	Type      WalletType     `gorm:"column:type;type:tinyint;not null;index"`
	Balance   float64        `gorm:"column:balance;type:decimal(15, 2);not null"`
	CreatedAt time.Time      `gorm:"column:createdAt;not null"`
	UpdatedAt time.Time      `gorm:"column:updatedAt;default:NULL"`
	DeletedAt gorm.DeletedAt `gorm:"column:deletedAt;index"`
}

func (Wallet) TableName() string {
	return "Wallet"
}

type WalletPayload struct {
	Type WalletType `json:"type" validate:"required,wallettype"`
}

type WalletHandler interface {
	Create(echo.Context) error
}

type WalletService interface {
	Create(ctx context.Context, payload *WalletPayload) error
}

type WalletRepository interface {
	Create(ctx context.Context, wallet *Wallet) error
	GetByUserID(ctx context.Context, userID uuid.UUID) (*Wallet, error)
}

func (w *WalletPayload) Validate() map[string]string {
	return ValidateStruct(w)
}

func (w *WalletPayload) ToWallet(userID uuid.UUID) *Wallet {
	return &Wallet{
		UserID:    userID,
		Type:      w.Type,
		CreatedAt: time.Now().UTC(),
	}
}
