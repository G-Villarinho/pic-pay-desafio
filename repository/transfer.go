package repository

import (
	"context"
	"log/slog"

	"github.com/GSVillas/pic-pay-desafio/domain"
	"github.com/go-redis/redis/v8"
	"github.com/google/uuid"
	"github.com/samber/do"
	"gorm.io/gorm"
)

type transferRepository struct {
	i           *do.Injector
	db          *gorm.DB
	redisClient *redis.Client
}

func NewTransferRepository(i *do.Injector) (domain.TransferRepository, error) {
	db, err := do.Invoke[*gorm.DB](i)
	if err != nil {
		return nil, err
	}

	redisClient, err := do.Invoke[*redis.Client](i)
	if err != nil {
		return nil, err
	}

	return &transferRepository{
		i:           i,
		db:          db,
		redisClient: redisClient,
	}, nil
}

func (t *transferRepository) Transfer(ctx context.Context, transfer *domain.Transfer) error {
	log := slog.With(
		slog.String("repository", "transfer"),
		slog.String("func", "Transfer"),
	)

	tx := t.db.Begin()
	if err := tx.Error; err != nil {
		log.Error("Failed to begin transaction", slog.String("error", err.Error()))
		return err
	}

	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
			log.Error("Panic occurred, transaction rolled back", slog.Any("recover", r))
			panic(r)
		}
	}()

	log.Info("Starting to process transfer", slog.String("payerID", transfer.PayerID.String()), slog.String("payeeID", transfer.PayeeID.String()), slog.Float64("value", transfer.Value))

	if err := t.debit(ctx, transfer.PayerID, transfer.Value); err != nil {
		tx.Rollback()
		log.Error("Failed to debit payer's wallet, transaction rolled back", slog.String("payerID", transfer.PayerID.String()), slog.Float64("value", transfer.Value), slog.String("error", err.Error()))
		return err
	}

	if err := t.credit(ctx, transfer.PayeeID, transfer.Value); err != nil {
		tx.Rollback()
		log.Error("Failed to credit payee's wallet, transaction rolled back", slog.String("payeeID", transfer.PayeeID.String()), slog.Float64("value", transfer.Value), slog.String("error", err.Error()))
		return err
	}

	if err := tx.WithContext(ctx).Create(transfer).Error; err != nil {
		tx.Rollback()
		log.Error("Failed to record transfer, transaction rolled back", slog.String("transferID", transfer.ID.String()), slog.String("error", err.Error()))
		return err
	}

	if err := tx.Commit().Error; err != nil {
		log.Error("Failed to commit transaction", slog.String("error", err.Error()))
		return err
	}

	log.Info("Transfer completed successfully", slog.String("payerID", transfer.PayerID.String()), slog.String("payeeID", transfer.PayeeID.String()), slog.Float64("value", transfer.Value))
	return nil
}

func (t *transferRepository) credit(ctx context.Context, userID uuid.UUID, value float64) error {
	log := slog.With(
		slog.String("repository", "wallet"),
		slog.String("func", "Credit"),
	)

	log.Info("Starting to credit value to user's wallet", slog.String("userID", userID.String()), slog.Float64("value", value))

	if err := t.db.WithContext(ctx).Model(&domain.Wallet{}).Where("userId = ?", userID).UpdateColumn("balance", gorm.Expr("balance + ?", value)).Error; err != nil {
		log.Error("Failed to credit value to wallet", slog.String("error", err.Error()))
		return err
	}

	log.Info("Successfully credited value to user's wallet", slog.String("userID", userID.String()), slog.Float64("value", value))
	return nil
}

func (t *transferRepository) debit(ctx context.Context, userID uuid.UUID, value float64) error {
	log := slog.With(
		slog.String("repository", "wallet"),
		slog.String("func", "Debit"),
	)

	log.Info("Starting to debit value from user's wallet", slog.String("userID", userID.String()), slog.Float64("value", value))

	if err := t.db.WithContext(ctx).Model(&domain.Wallet{}).Where("userId = ?", userID).UpdateColumn("balance", gorm.Expr("balance - ?", value)).Error; err != nil {
		log.Error("Failed to debit value from wallet", slog.String("error", err.Error()))
		return err
	}

	log.Info("Successfully debited value from user's wallet", slog.String("userID", userID.String()), slog.Float64("value", value))
	return nil
}
