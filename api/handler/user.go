package handler

import (
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
		apiError := domain.NewAPIError(http.StatusUnprocessableEntity, "Invalid Request", "Failed to process the payload").
			WithErrors(map[string]string{"binding": "Failed to bind payload"})
		return ctx.JSON(apiError.Status, apiError)
	}

	errors := payload.Validate()
	if errors != nil {
		log.Warn("Validation failed", slog.Any("errors", errors))
		apiError := domain.NewAPIError(http.StatusBadRequest, "Validation Failed", "One or more fields failed validation").
			WithErrors(errors)
		return ctx.JSON(apiError.Status, apiError)
	}

	return ctx.JSON(http.StatusCreated, "User created successfully")
}
