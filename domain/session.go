package domain

import (
	"context"
	"errors"

	"github.com/google/uuid"
)

var (
	ErrTokenInvalid           = errors.New("invalid token")
	ErrSessionNotFound        = errors.New("token not found")
	ErrorUnexpectedMethod     = errors.New("unexpected signing method")
	ErrTokenNotFoundInContext = errors.New("token not found in context")
	ErrSessionMismatch        = errors.New("session icompatible for user requested")
)

type Session struct {
	Token  string    `json:"token"`
	Name   string    `json:"name"`
	UserID uuid.UUID `json:"picPayId"`
	Email  string    `json:"email"`
}

type SessionService interface {
	Create(ctx context.Context, user *User) (string, error)
	GetSession(ctx context.Context, token string) (*Session, error)
}

type SessionRepository interface {
	Create(ctx context.Context, session *Session) error
	GetSession(ctx context.Context, userID uuid.UUID) (*Session, error)
}
