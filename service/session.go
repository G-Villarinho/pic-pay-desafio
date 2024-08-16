package service

import (
	"context"
	"log/slog"

	"github.com/GSVillas/pic-pay-desafio/config"
	"github.com/GSVillas/pic-pay-desafio/domain"
	"github.com/golang-jwt/jwt"
	jsoniter "github.com/json-iterator/go"
	"github.com/samber/do"
)

type sessionService struct {
	i                 *do.Injector
	sessionRepository domain.SessionRepository
}

func NewSessionService(i *do.Injector) (domain.SessionService, error) {
	sessionRepository, err := do.Invoke[domain.SessionRepository](i)
	if err != nil {
		return nil, err
	}

	return &sessionService{
		i:                 i,
		sessionRepository: sessionRepository,
	}, nil
}

func (s *sessionService) Create(ctx context.Context, user *domain.User) (string, error) {
	log := slog.With(
		slog.String("service", "session"),
		slog.String("func", "Create"),
	)

	log.Info("Initializing create user session process")

	token, err := s.createToken(user)
	if err != nil {
		log.Error("Failed to create token", slog.String("error", err.Error()))
		return "", err
	}

	session := &domain.Session{
		Token:  token,
		Name:   user.Name,
		UserID: user.ID,
		Email:  user.Email,
	}

	if err := s.sessionRepository.Create(ctx, session); err != nil {
		log.Error("Failed to save user session", slog.String("error", err.Error()))
		return "", err
	}

	log.Info("session creation process excuted succefully")
	return token, err
}

func (s *sessionService) GetSession(ctx context.Context, token string) (*domain.Session, error) {
	log := slog.With(
		slog.String("service", "session"),
		slog.String("func", "GetSession"),
		slog.String("token", token),
	)

	log.Info("Starting session retrieval process")

	sessionToken, err := s.extractSessionFromToken(token)
	if err != nil {
		log.Error("Failed to extract session from token", slog.String("error", err.Error()))
		return nil, err
	}

	session, err := s.sessionRepository.GetSession(ctx, sessionToken.UserID)
	if err != nil {
		log.Error("Failed to retrieve session from repository", slog.Any("userID", sessionToken.UserID), slog.String("error", err.Error()))
		return nil, err
	}

	if session == nil {
		log.Warn("Session not found in repository", slog.Any("userID", sessionToken.UserID))
		return nil, domain.ErrSessionNotFound
	}

	if token != session.Token {
		log.Warn("Session token mismatch", slog.Any("userID", sessionToken.UserID))
		return nil, domain.ErrSessionMismatch
	}

	log.Info("Session retrieved successfully", slog.Any("userID", sessionToken.UserID))
	return session, nil
}

func (s *sessionService) createToken(user *domain.User) (string, error) {
	log := slog.With(
		slog.String("service", "session"),
		slog.String("func", "createToken"),
	)

	log.Info("Initializing create token process")

	token := jwt.NewWithClaims(jwt.SigningMethodES256, jwt.MapClaims{
		"picPayId": user.ID,
		"name":     user.Name,
		"email":    user.Email,
	})

	tokenString, err := token.SignedString(config.Env.PrivateKey)
	if err != nil {
		log.Error("Error to signed token string", slog.Any("error:", err.Error()))
		return "", err
	}

	log.Info("Create token process executed successfully")
	return tokenString, nil
}

func (s *sessionService) extractSessionFromToken(tokenString string) (*domain.Session, error) {
	log := slog.With(
		slog.String("service", "session"),
		slog.String("func", "extractSessionFromToken"),
	)

	log.Info("Initializing extract session from token process")

	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (any, error) {
		if _, ok := token.Method.(*jwt.SigningMethodECDSA); !ok {
			log.Error("Unexpected signing method", slog.String("expected", "ECDSA"), slog.String("actual", token.Method.Alg()))
			return nil, domain.ErrorUnexpectedMethod
		}
		log.Info("Token signing method validated", slog.String("method", token.Method.Alg()))
		return config.Env.PublicKey, nil
	})

	if err != nil {
		log.Error("Failed to parse token", slog.Any("error: ", err.Error()))
		return nil, err
	}

	if !token.Valid {
		log.Warn("Invalid token", slog.String("token", tokenString))
		return nil, domain.ErrTokenInvalid
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		log.Error("Failed to assert token claims as MapClaims")
		return nil, domain.ErrTokenInvalid
	}

	sessionJSON, err := jsoniter.Marshal(claims)
	if err != nil {
		log.Error("Failed to marshal claims to JSON", slog.Any("error: ", err.Error()))
		return nil, err
	}

	log.Info("Claims marshaled to JSON successfully", slog.String("json", string(sessionJSON)))

	var session domain.Session
	err = jsoniter.Unmarshal(sessionJSON, &session)
	if err != nil {
		log.Error("Failed to unmarshal JSON to session struct", slog.Any("error: ", err.Error()))
		return nil, err
	}

	log.Info("Session successfully extracted from token", slog.Any("session", session))
	return &session, nil
}
