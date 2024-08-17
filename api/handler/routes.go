package handler

import (
	"github.com/GSVillas/pic-pay-desafio/domain"
	"github.com/GSVillas/pic-pay-desafio/middleware"
	"github.com/labstack/echo/v4"
	"github.com/samber/do"
)

func SetupRoutes(e *echo.Echo, i *do.Injector) {
	setupUserRoutes(e, i)
	setupWalletRoutes(e, i)
}

func setupUserRoutes(e *echo.Echo, i *do.Injector) {
	userHandler, err := do.Invoke[domain.UserHandler](i)
	if err != nil {
		panic(err)
	}

	group := e.Group("/v1/users")
	group.POST("", userHandler.Create)
	group.POST("/sign-in", userHandler.SignIn)
}

func setupWalletRoutes(e *echo.Echo, i *do.Injector) {
	walletHandler, err := do.Invoke[domain.WalletHandler](i)
	if err != nil {
		panic(err)
	}

	group := e.Group("v1/wallets", middleware.CheckLoggedIn(i))
	group.POST("", walletHandler.Create)
}
