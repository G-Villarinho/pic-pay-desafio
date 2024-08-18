package handler

import (
	"errors"
	"log/slog"
	"net/http"

	"github.com/GSVillas/pic-pay-desafio/client"
	"github.com/GSVillas/pic-pay-desafio/domain"
	jsoniter "github.com/json-iterator/go"
	"github.com/labstack/echo/v4"
	"github.com/samber/do"
)

type transferHandler struct {
	i               *do.Injector
	transferService domain.TransferService
}

func NewTransferHandler(i *do.Injector) (domain.TransferHandler, error) {
	transferService, err := do.Invoke[domain.TransferService](i)
	if err != nil {
		return nil, err
	}

	return &transferHandler{
		i:               i,
		transferService: transferService,
	}, nil
}

func (t *transferHandler) Transfer(ctx echo.Context) error {
	log := slog.With(
		slog.String("handler", "transfer"),
		slog.String("func", "Transfer"),
	)

	log.Info("Initializing transfer process")

	var payload domain.TransferPayload
	if err := jsoniter.NewDecoder(ctx.Request().Body).Decode(&payload); err != nil {
		log.Warn("Failed to decode JSON payload", slog.String("error", err.Error()))
		return ctx.JSON(http.StatusUnprocessableEntity, domain.CannotBindPayloadAPIError)
	}

	validationErrors := payload.Validate()
	if validationErrors != nil {
		log.Warn("Validation failed", slog.Any("errors", validationErrors))
		apiError := domain.NewAPIError(http.StatusBadRequest, "Validation Failed", "One or more fields failed validation").
			WithErrors(validationErrors)
		return ctx.JSON(apiError.Status, apiError)
	}

	if err := t.transferService.Transfer(ctx.Request().Context(), &payload); err != nil {

		if errors.Is(err, domain.ErrSelfTransactionNotAllowed) {
			log.Warn("Transfer failed due to self-transfer attempt", slog.String("error", err.Error()))
			apiError := domain.NewAPIError(http.StatusForbidden, "Forbidden", "You cannot transfer money to yourself.")
			return ctx.JSON(http.StatusForbidden, apiError)
		}

		if errors.Is(err, domain.ErrPayerWalletNotFound) {
			log.Warn("Transfer failed due to missing payer wallet", slog.String("error", err.Error()))
			apiError := domain.NewAPIError(http.StatusNotFound, "Not Found", "Payer wallet not found.")
			return ctx.JSON(http.StatusNotFound, apiError)
		}

		if errors.Is(err, domain.ErrPayeeWalletNotFound) {
			log.Warn("Transfer failed due to missing payee wallet", slog.String("error", err.Error()))
			apiError := domain.NewAPIError(http.StatusNotFound, "Not Found", "Payee wallet not found.")
			return ctx.JSON(http.StatusNotFound, apiError)
		}

		if errors.Is(err, domain.ErrTransferNotAllowedForWalletType) {
			log.Warn("Transfer failed due to wallet type restriction", slog.String("error", err.Error()))
			apiError := domain.NewAPIError(http.StatusForbidden, "Forbidden", "Transfers are not allowed for this wallet type.")
			return ctx.JSON(http.StatusForbidden, apiError)
		}

		if errors.Is(err, domain.ErrInsufficientBalance) {
			log.Warn("Transfer failed due to insufficient balance", slog.String("error", err.Error()))
			apiError := domain.NewAPIError(http.StatusBadRequest, "Bad Request", "Insufficient balance for the transaction.")
			return ctx.JSON(http.StatusBadRequest, apiError)
		}

		if errors.Is(err, domain.ErrTransferNotAuthorized) || errors.Is(err, client.ErrCheckAuthorization) {
			log.Warn("Transfer failed due to authorization error", slog.String("error", err.Error()))
			apiError := domain.NewAPIError(http.StatusUnauthorized, "Unauthorized", "Transfer not authorized.")
			return ctx.JSON(http.StatusUnauthorized, apiError)
		}

		log.Error("Failed to process transfer", slog.String("error", err.Error()))
		return ctx.JSON(http.StatusInternalServerError, domain.InternalServerAPIError)
	}

	log.Info("Transfer completed successfully")
	return ctx.NoContent(http.StatusCreated)
}
