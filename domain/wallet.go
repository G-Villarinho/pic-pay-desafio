package domain

import (
	"context"
	"errors"
	"time"

	"github.com/GSVillas/pic-pay-desafio/domain/types"
	"github.com/GSVillas/pic-pay-desafio/utils"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

var (
	ErrWalletAlredyRegister = errors.New("the user already has a wallet")
)

type Wallet struct {
	UserID    uuid.UUID        `gorm:"column:userId;type:char(36);primaryKey"`
	User      User             `gorm:"foreignKey:UserID"`
	Type      types.WalletType `gorm:"column:type;type:tinyint;not null;index"`
	Balance   float64          `gorm:"column:Balance;type:decimal(15, 2);not null"`
	CreatedAt time.Time        `gorm:"column:createdAt;index"`
	UpdatedAt time.Time        `gorm:"column:updatedAt"`
}

func (Wallet) TableName() string {
	return "Wallet"
}

type WalletPayload struct {
	Type types.WalletType `json:"type" validate:"required,wallettype"`
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
	return utils.ValidateStruct(w)
}

func (w *WalletPayload) ToWallet(userID uuid.UUID) *Wallet {
	return &Wallet{
		UserID:    userID,
		Type:      w.Type,
		CreatedAt: time.Now().UTC(),
	}
}
