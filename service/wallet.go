package service

import (
	"context"
	"log/slog"

	"github.com/GSVillas/pic-pay-desafio/domain"
	"github.com/GSVillas/pic-pay-desafio/domain/types"
	"github.com/samber/do"
)

type walletService struct {
	i                *do.Injector
	walletRepository domain.WalletRepository
}

func NewWalletService(i *do.Injector) (domain.WalletService, error) {
	walletRepository, err := do.Invoke[domain.WalletRepository](i)
	if err != nil {
		return nil, err
	}

	return &walletService{
		i:                i,
		walletRepository: walletRepository,
	}, nil
}

func (w *walletService) Create(ctx context.Context, payload *domain.WalletPayload) error {
	log := slog.With(
		slog.String("service", "wallet"),
		slog.String("func", "Create"),
	)

	log.Info("Initializing create wallet process")

	session, ok := ctx.Value(types.SessionKey).(*domain.Session)
	if !ok || session == nil {
		return domain.ErrSessionNotFound
	}

	wallet, err := w.walletRepository.GetByUserID(ctx, session.UserID)
	if err != nil {
		log.Error("Failed to get wallet by ", slog.String("userId", session.UserID.String()))
		return err
	}

	if wallet != nil {
		log.Warn("There is already a wallet for this user ", slog.String("userId", session.UserID.String()))
		return domain.ErrWalletAlredyRegister
	}

	wallet = payload.ToWallet(session.UserID)
	if err := w.walletRepository.Create(ctx, wallet); err != nil {
		log.Error("Failed to create wallet", slog.String("error", err.Error()))
		return err
	}

	log.Info("Wallet creation process executed successfully")
	return nil
}
