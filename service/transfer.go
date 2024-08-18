package service

import (
	"context"
	"log/slog"

	"github.com/GSVillas/pic-pay-desafio/client"
	"github.com/GSVillas/pic-pay-desafio/domain"
	"github.com/samber/do"
)

type transactionService struct {
	i                    *do.Injector
	transferRepository   domain.TransferRepository
	walletRepository     domain.WalletRepository
	authorizationService client.AuthorizationService
}

func NewTransferService(i *do.Injector) (domain.TransferService, error) {
	transactionRepository, err := do.Invoke[domain.TransferRepository](i)
	if err != nil {
		return nil, err
	}

	walletRepository, err := do.Invoke[domain.WalletRepository](i)
	if err != nil {
		return nil, err
	}

	authorizationService, err := do.Invoke[client.AuthorizationService](i)
	if err != nil {
		return nil, err
	}

	return &transactionService{
		i:                    i,
		transferRepository:   transactionRepository,
		walletRepository:     walletRepository,
		authorizationService: authorizationService,
	}, nil
}

func (t *transactionService) Transfer(ctx context.Context, payload *domain.TransferPayload) error {
	log := slog.With(
		slog.String("service", "transaction"),
		slog.String("func", "Transfer"),
	)

	log.Info("Initializing create transaction process")

	session, ok := ctx.Value(domain.SessionKey).(*domain.Session)
	if !ok || session == nil {
		return domain.ErrSessionNotFound
	}

	if session.UserID == payload.PayeeID {
		log.Warn("Attempted self-transfer detected", slog.String("userID", session.UserID.String()), slog.String("action", "transaction to self"))
		return domain.ErrSelfTransactionNotAllowed
	}

	payer, err := t.walletRepository.GetByUserID(ctx, session.UserID)
	if err != nil {
		log.Error("Failed to get wallet by userID ", slog.String("Error: ", err.Error()))
		return domain.ErrGetWallet
	}

	if payer == nil {
		log.Warn("No wallets were found for this user", slog.String("userId: ", session.UserID.String()))
		return domain.ErrPayerWalletNotFound
	}

	payee, err := t.walletRepository.GetByUserID(ctx, payload.PayeeID)
	if err != nil {
		log.Error("Failed to get wallet by userID ", slog.String("Error", err.Error()))
		return domain.ErrGetWallet
	}

	if payee == nil {
		log.Warn("No wallets were found for this user", slog.String("userId", payload.PayeeID.String()))
		return domain.ErrPayeeWalletNotFound
	}

	if err := t.validateTransfer(ctx, payload, payer); err != nil {
		log.Warn("Transfer validation failed", slog.String("error", err.Error()))
		return err
	}

	transaction := payload.ToTansaction(payer.UserID)
	if err := t.transferRepository.Transfer(ctx, transaction); err != nil {
		log.Error("Failed to create transaction the user's wallet", slog.String("error", err.Error()))
		return domain.ErrCreateTransfer
	}

	log.Info("Session retrieved successfully")
	return nil
}

func (t *transactionService) validateTransfer(ctx context.Context, payload *domain.TransferPayload, payer *domain.Wallet) error {
	log := slog.With(
		slog.String("service", "transaction"),
		slog.String("func", "validateTransfer"),
	)

	log.Info("Initializing validation transfer process")

	if payer.Type == domain.WalletTypeMERCHANT {
		log.Warn("Transfer not allowed for merchant wallet", slog.String("walletType", "MERCHANT"))
		return domain.ErrTransferNotAllowedForWalletType
	}

	if payer.Balance < payload.Value {
		log.Warn("Insufficient balance for transaction")
		return domain.ErrInsufficientBalance
	}

	authorizationData, err := t.authorizationService.CheckAuthorization(ctx)
	if err != nil {
		log.Error("Error to check user authorization", slog.String("Error: ", err.Error()))
		return client.ErrCheckAuthorization
	}

	if !authorizationData.Data.Authorization {
		log.Warn("Transfer authorization failed")
		return domain.ErrTransferNotAuthorized
	}

	log.Info("Validation transfer process successfully")
	return nil
}
