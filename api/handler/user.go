package handler

import (
	"errors"
	"log/slog"
	"net/http"

	"github.com/GSVillas/pic-pay-desafio/domain"
	"github.com/labstack/echo/v4"
	"github.com/samber/do"
)

type userHandler struct {
	i           *do.Injector
	userService domain.UserService
}

func NewUserHandler(i *do.Injector) (domain.UserHandler, error) {
	userService, err := do.Invoke[domain.UserService](i)
	if err != nil {
		return nil, err
	}

	return &userHandler{
		i:           i,
		userService: userService,
	}, nil
}

func (u *userHandler) Create(ctx echo.Context) error {
	log := slog.With(
		slog.String("handler", "user"),
		slog.String("func", "Create"),
	)

	log.Info("Initializing user creation process")

	var payload domain.UserPayload
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

	if err := u.userService.Create(ctx.Request().Context(), &payload); err != nil {
		log.Warn("Fail to create user", slog.Any("error", err))

		if errors.Is(err, domain.ErrUserAlreadyExists) {
			apiError := domain.NewAPIError(http.StatusConflict, "conflict", "The user already exists. Please try again with a different email.")
			return ctx.JSON(http.StatusConflict, apiError)
		}

		return ctx.JSON(http.StatusInternalServerError, domain.InternalServerAPIError)
	}

	log.Info("User created successfully")
	return ctx.NoContent(http.StatusCreated)
}
