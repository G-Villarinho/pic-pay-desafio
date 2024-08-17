package middleware

import (
	"context"
	"net/http"
	"strings"

	"github.com/GSVillas/pic-pay-desafio/domain"
	"github.com/labstack/echo/v4"
	"github.com/samber/do"
)

func CheckLoggedIn(i *do.Injector) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(ctx echo.Context) error {
			sessionService, err := do.Invoke[domain.SessionService](i)
			if err != nil {
				return ctx.JSON(http.StatusInternalServerError, domain.InternalServerAPIError)
			}

			authorizationHeader := ctx.Request().Header.Get("Authorization")
			if authorizationHeader == "" {
				apiError := domain.NewAPIError(http.StatusUnauthorized, "Access Denied", "You need to be logged in to access this resource.")
				return ctx.JSON(http.StatusUnauthorized, apiError)
			}

			content := strings.Split(authorizationHeader, " ")
			if len(content) != 2 {
				apiError := domain.NewAPIError(http.StatusUnauthorized, "Access Denied", "You need to be logged in to access this resource.")
				return ctx.JSON(http.StatusUnauthorized, apiError)
			}

			token := content[1]

			session, err := sessionService.GetSession(ctx.Request().Context(), token)
			if err != nil {
				if err == domain.ErrTokenInvalid || err == domain.ErrSessionMismatch || err == domain.ErrSessionNotFound {
					apiError := domain.NewAPIError(http.StatusUnauthorized, "Access Denied", "Invalid or expired session token.")
					return ctx.JSON(http.StatusUnauthorized, apiError)
				}

				return ctx.JSON(http.StatusInternalServerError, domain.InternalServerAPIError)
			}

			newCtx := context.WithValue(ctx.Request().Context(), domain.SessionKey, session)
			ctx.SetRequest(ctx.Request().WithContext(newCtx))

			return next(ctx)
		}
	}
}
