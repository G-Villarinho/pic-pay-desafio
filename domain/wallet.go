package domain

import (
	"time"

	"github.com/google/uuid"
)

type WalletType uint8

const (
	COMMON WalletType = iota
	MERCHANT
)

type Wallet struct {
	UserID    uuid.UUID  `gorm:"column:userId;type:char(36);primaryKey"`
	User      User       `gorm:"foreignKey:UserID"`
	Type      WalletType `gorm:"column:type;type:tinyint;not null;index"`
	Balance   float64    `gorm:"column:Balance;type:decimal(15, 2);not null"`
	CreatedAt time.Time  `gorm:"column:createdAt;index"`
	UpdatedAt time.Time  `gorm:"column:updatedAt"`
}

func (Wallet) TableName() string {
	return "Wallet"
}

func (wt WalletType) IsValid() bool {
	switch wt {
	case COMMON, MERCHANT:
		return true
	}
	return false
}
