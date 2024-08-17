package handler

import (
	"errors"
	"log/slog"
	"net/http"

	"github.com/GSVillas/pic-pay-desafio/domain"
	"github.com/labstack/echo/v4"
	"github.com/samber/do"
)

type walletHandler struct {
	i             *do.Injector
	walletService domain.WalletService
}

func NewWalletHandler(i *do.Injector) (domain.WalletHandler, error) {
	walletService, err := do.Invoke[domain.WalletService](i)
	if err != nil {
		return nil, err
	}

	return &walletHandler{
		i:             i,
		walletService: walletService,
	}, nil
}

func (w *walletHandler) Create(ctx echo.Context) error {
	log := slog.With(
		slog.String("handler", "wallet"),
		slog.String("func", "Create"),
	)

	log.Info("Initializing create wallet process")

	var payload domain.WalletPayload
	if err := ctx.Bind(&payload); err != nil {
		log.Warn("Failed to bind payload", slog.String("error", err.Error()))
		return ctx.JSON(http.StatusUnprocessableEntity, domain.CannotBindPayloadAPIError)
	}

	validationErrors := payload.Validate()
	if validationErrors != nil {
		log.Warn("Validation failed", slog.Any("errors", validationErrors))
		apiError := domain.NewAPIError(http.StatusBadRequest, "Validation Failed", "One or more fields failed validation").
			WithErrors(validationErrors)
		return ctx.JSON(apiError.Status, apiError)
	}

	if err := w.walletService.Create(ctx.Request().Context(), &payload); err != nil {

		if errors.Is(err, domain.ErrSessionNotFound) {
			log.Warn("Unauthorized attempt to create wallet", slog.String("error", err.Error()))
			return ctx.JSON(http.StatusForbidden, domain.SessionNotFoundAPIError)
		}

		if errors.Is(err, domain.ErrWalletAlredyRegister) {
			log.Warn("Fail to create wallet", slog.String("error", err.Error()))
			apiError := domain.NewAPIError(http.StatusConflict, "conflict", "The user already has a wallet")
			return ctx.JSON(http.StatusConflict, apiError)
		}

		log.Error("Fail to create user wallet", slog.Any("error", err))
		return ctx.JSON(http.StatusInternalServerError, domain.InternalServerAPIError)
	}

	log.Info("Create wallet process executed succefully")
	return ctx.NoContent(http.StatusCreated)
}
