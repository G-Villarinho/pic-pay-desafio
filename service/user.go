package service

import (
	"context"
	"log/slog"

	"github.com/GSVillas/pic-pay-desafio/domain"
	"github.com/GSVillas/pic-pay-desafio/secure"
	"github.com/klassmann/cpfcnpj"
	"github.com/samber/do"
)

type userService struct {
	i              *do.Injector
	userRepository domain.UserRepository
}

func NewUserService(i *do.Injector) (domain.UserService, error) {
	userRepository, err := do.Invoke[domain.UserRepository](i)
	if err != nil {
		return nil, err
	}

	return &userService{
		i:              i,
		userRepository: userRepository,
	}, nil
}

func (u *userService) Create(ctx context.Context, payload *domain.UserPayload) error {
	log := slog.With(
		slog.String("service", "user"),
		slog.String("func", "Create"),
	)

	log.Info("Initializing user creation process")

	user, err := u.userRepository.GetByEmail(ctx, payload.Email)
	if err != nil {
		log.Error("Failed to get user by email", slog.String("error", err.Error()))
		return err
	}

	if user != nil {
		log.Warn("There is already a user with this ", slog.String("email:", payload.Email))
		return domain.ErrEmailAlreadyRegister
	}

	user, err = u.userRepository.GetByCPF(ctx, string(cpfcnpj.NewCPF(payload.CPF)))
	if err != nil {
		log.Error("Failed to get user by cpf", slog.String("error", err.Error()))
		return err
	}

	if user != nil {
		log.Warn("There is already a user with this cpf")
		return domain.ErrCPFAlreadyRegister
	}

	passwordHash, err := secure.Hash(payload.Password)
	if err != nil {
		log.Error("Failed to hash password", slog.String("error", err.Error()))
		return domain.ErrHashingPassword
	}

	user = payload.ToUser(string(passwordHash))

	if err := u.userRepository.Create(ctx, user); err != nil {
		log.Error("Failed to create user", slog.String("error", err.Error()))
		return err
	}

	log.Info("user creation process executed successfully")
	return nil
}
