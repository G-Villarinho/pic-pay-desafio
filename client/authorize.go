package client

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net/http"

	"github.com/GSVillas/pic-pay-desafio/config"
	jsoniter "github.com/json-iterator/go"
	"github.com/samber/do"
)

var (
	ErrUnauthorized         = errors.New("authorization failed: user is not authorized")
	ErrAPIConnection        = errors.New("failed to connect to the authorization API")
	ErrCheckAuthorization   = errors.New("erro to check transactio authorization")
	ErrUnexpectedStatusCode = func(statusCode int) error {
		return fmt.Errorf("unexpected status code from authorization API: %d", statusCode)
	}
)

type AuthorizationResponse struct {
	Status string            `json:"status"`
	Data   AuthorizationData `json:"data"`
}

type AuthorizationData struct {
	Authorization bool `json:"authorization"`
}

type AuthorizationService interface {
	CheckAuthorization(ctx context.Context) (*AuthorizationResponse, error)
}

type authorizationService struct {
	i          *do.Injector
	httpClient *http.Client
}

func NewAuthorizationService(i *do.Injector) (AuthorizationService, error) {
	httpClient, err := do.Invoke[*http.Client](i)
	if err != nil {
		return nil, err
	}

	return &authorizationService{
		i:          i,
		httpClient: httpClient,
	}, nil
}

func (a *authorizationService) CheckAuthorization(ctx context.Context) (*AuthorizationResponse, error) {
	log := slog.With(
		slog.String("service", "authorization"),
		slog.String("func", "CheckAuthorization"),
	)

	log.Info("Initializing check authorization process")

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, config.Env.AuthorizationURL, nil)
	if err != nil {
		log.Error("Failed to create request", slog.String("error", err.Error()))
		return nil, ErrAPIConnection
	}

	resp, err := a.httpClient.Do(req)
	if err != nil {
		log.Error("Failed to perform HTTP request", slog.String("error", err.Error()))
		return nil, ErrAPIConnection
	}
	defer func() {
		if err := resp.Body.Close(); err != nil {
			log.Error("Failed to close response body", slog.String("error", err.Error()))
		}
	}()

	log.Info("HTTP request completed", slog.Int("statusCode", resp.StatusCode))

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusForbidden {
		log.Warn("Unexpected status code received", slog.Int("statusCode", resp.StatusCode))
		return nil, ErrUnexpectedStatusCode(resp.StatusCode)
	}

	var authorizationResponse *AuthorizationResponse
	if err := jsoniter.NewDecoder(resp.Body).Decode(&authorizationResponse); err != nil {
		log.Error("Failed to decode response body", slog.String("error", err.Error()))
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	log.Info("Authorization response decoded successfully", slog.Bool("authorized", authorizationResponse.Data.Authorization))

	return authorizationResponse, nil
}
