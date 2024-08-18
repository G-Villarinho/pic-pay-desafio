package repository

import (
	"context"
	"errors"
	"log/slog"

	"github.com/GSVillas/pic-pay-desafio/domain"
	"github.com/go-redis/redis/v8"
	"github.com/google/uuid"
	"github.com/samber/do"
	"gorm.io/gorm"
)

type walletRepository struct {
	i           *do.Injector
	db          *gorm.DB
	redisClient *redis.Client
}

func NewWalletRepository(i *do.Injector) (domain.WalletRepository, error) {
	db, err := do.Invoke[*gorm.DB](i)
	if err != nil {
		return nil, err
	}

	redisClient, err := do.Invoke[*redis.Client](i)
	if err != nil {
		return nil, err
	}

	return &walletRepository{
		i:           i,
		db:          db,
		redisClient: redisClient,
	}, nil
}

func (w *walletRepository) Create(ctx context.Context, wallet *domain.Wallet) error {
	log := slog.With(
		slog.String("repository", "wallet"),
		slog.String("func", "Create"),
	)

	log.Info("Initializing create wallet process")
	if err := w.db.WithContext(ctx).Create(&wallet).Error; err != nil {
		log.Error("Failed to create wallet", slog.String("error", err.Error()))
		return err
	}

	log.Info("Create user wallet process executed successfully")
	return nil
}

func (w *walletRepository) GetByUserID(ctx context.Context, userID uuid.UUID) (*domain.Wallet, error) {
	log := slog.With(
		slog.String("repository", "wallet"),
		slog.String("func", "GetByUserID"),
	)

	log.Info("Initializing get wallet by userId process")

	var wallet *domain.Wallet
	if err := w.db.WithContext(ctx).Where("userId = ?", userID.String()).First(&wallet).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			log.Warn("Wallet not found")
			return nil, nil
		}

		log.Error("Failed to get user by email", slog.String("error", err.Error()))
		return nil, err
	}

	log.Info("Process of obtaining wallet by userID executed successfully")
	return wallet, nil
}

func (w *walletRepository) Credit(ctx context.Context, userID uuid.UUID, value float64) error {
	log := slog.With(
		slog.String("repository", "wallet"),
		slog.String("func", "Credit"),
	)

	log.Info("Starting to credit value to user's wallet", slog.String("userID", userID.String()), slog.Float64("value", value))

	if err := w.db.WithContext(ctx).Model(&domain.Wallet{}).Where("userId = ?", userID).UpdateColumn("balance", gorm.Expr("balance + ?", value)).Error; err != nil {
		log.Error("Failed to credit value to wallet", slog.String("error", err.Error()))
		return err
	}

	log.Info("Successfully credited value to user's wallet", slog.String("userID", userID.String()), slog.Float64("value", value))
	return nil
}

func (w *walletRepository) Debit(ctx context.Context, userID uuid.UUID, value float64) error {
	log := slog.With(
		slog.String("repository", "wallet"),
		slog.String("func", "Debit"),
	)

	log.Info("Starting to debit value from user's wallet", slog.String("userID", userID.String()), slog.Float64("value", value))

	if err := w.db.WithContext(ctx).Model(&domain.Wallet{}).Where("userId = ?", userID).UpdateColumn("balance", gorm.Expr("balance - ?", value)).Error; err != nil {
		log.Error("Failed to debit value from wallet", slog.String("error", err.Error()))
		return err
	}

	log.Info("Successfully debited value from user's wallet", slog.String("userID", userID.String()), slog.Float64("value", value))
	return nil
}
