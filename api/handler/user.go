package handler

import (
	"errors"
	"log/slog"
	"net/http"

	"github.com/GSVillas/pic-pay-desafio/domain"
	jsoniter "github.com/json-iterator/go"
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

	if err := u.userService.Create(ctx.Request().Context(), &payload); err != nil {

		if errors.Is(err, domain.ErrEmailAlreadyRegister) {
			log.Warn("Fail to create user", slog.Any("error", err))
			apiError := domain.NewAPIError(http.StatusConflict, "conflict", "The email already registered. Please try again with a different email.")
			return ctx.JSON(http.StatusConflict, apiError)
		}

		if errors.Is(err, domain.ErrCPFAlreadyRegister) {
			log.Warn("Fail to create user", slog.Any("error", err))
			apiError := domain.NewAPIError(http.StatusConflict, "conflict", "The cpf already registered. Please try again with a different cpf.")
			return ctx.JSON(http.StatusConflict, apiError)
		}

		log.Error("Fail to create user", slog.Any("error", err))
		return ctx.JSON(http.StatusInternalServerError, domain.InternalServerAPIError)
	}

	log.Info("User created successfully")
	return ctx.NoContent(http.StatusCreated)
}

func (u *userHandler) SignIn(ctx echo.Context) error {
	log := slog.With(
		slog.String("handler", "user"),
		slog.String("func", "SignIn"),
	)

	log.Info("Initializing user sign in process")

	var payload domain.SignInPayload
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

	response, err := u.userService.SignIn(ctx.Request().Context(), &payload)
	if err != nil {

		if errors.Is(err, domain.ErrUserNotFound) || errors.Is(err, domain.ErrInvalidPassword) {
			log.Warn("Fail to excute user sign in", slog.Any("error", err))
			apiError := domain.NewAPIError(http.StatusUnauthorized, "Unauthorized credentials", "Unauthorized credentials. Review the data sent")
			return ctx.JSON(http.StatusUnauthorized, apiError)
		}

		log.Error("Fail to create user", slog.Any("error", err))
		return ctx.JSON(http.StatusInternalServerError, domain.InternalServerAPIError)
	}

	log.Info("user sign in executed succefully")
	return ctx.JSON(http.StatusOK, response)
}
