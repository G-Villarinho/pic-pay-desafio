package config

import (
	"log/slog"
	"os"

	"github.com/GSVillas/pic-pay-desafio/config/models"
	"github.com/Netflix/go-env"
	"github.com/joho/godotenv"
)

var Env models.Environment

func LoadEnvironments() {
	err := godotenv.Load()
	if err != nil {
		panic(err)
	}

	_, err = env.UnmarshalFromEnviron(&Env)
	if err != nil {
		panic(err)
	}

}

func ConfigureLogger() {
	handler := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		AddSource: false,
	}))
	slog.SetDefault(handler)
}
