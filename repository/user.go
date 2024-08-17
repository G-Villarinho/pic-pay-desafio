package repository

import (
	"context"
	"errors"
	"log/slog"

	"github.com/GSVillas/pic-pay-desafio/domain"
	"github.com/go-redis/redis/v8"
	"github.com/google/uuid"
	"github.com/samber/do"
	"gorm.io/gorm"
)

type userRepository struct {
	i           *do.Injector
	db          *gorm.DB
	redisCLient *redis.Client
}

func NewUserRepository(i *do.Injector) (domain.UserRepository, error) {
	db, err := do.Invoke[*gorm.DB](i)
	if err != nil {
		return nil, err
	}

	redisClient, err := do.Invoke[*redis.Client](i)
	if err != nil {
		return nil, err
	}

	return &userRepository{
		i:           i,
		db:          db,
		redisCLient: redisClient,
	}, nil
}

func (u *userRepository) Create(ctx context.Context, user *domain.User) error {
	log := slog.With(
		slog.String("repository", "user"),
		slog.String("func", "Create"),
	)

	log.Info("Initializing user creation process")
	if err := u.db.WithContext(ctx).Create(&user).Error; err != nil {
		log.Error("Failed to create user", slog.String("error", err.Error()))
		return err
	}

	log.Info("Create user process executed successfully")
	return nil
}

func (u *userRepository) GetByEmail(ctx context.Context, email string) (*domain.User, error) {
	log := slog.With(
		slog.String("repository", "user"),
		slog.String("func", "GetByEmail"),
	)

	log.Info("Initializing process of obtaining user by email")

	var user *domain.User
	if err := u.db.WithContext(ctx).Where("email = ?", email).First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			log.Warn("User not found")
			return nil, nil
		}

		log.Error("Failed to get user by email", slog.String("error", err.Error()))
		return nil, err
	}

	log.Info("Process of obtaining user by email executed successfully")
	return user, nil
}

func (u *userRepository) GetByCPF(ctx context.Context, CPF string) (*domain.User, error) {
	log := slog.With(
		slog.String("repository", "user"),
		slog.String("func", "GetByEmail"),
	)

	log.Info("Initializing process of obtaining user by CPF")

	var user *domain.User
	if err := u.db.WithContext(ctx).Where("cpf = ?", CPF).First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			log.Warn("User not found")
			return nil, nil
		}

		log.Error("Failed to get user by email", slog.String("error", err.Error()))
		return nil, err
	}

	log.Info("Process of obtaining user by cpf executed successfully")
	return user, nil
}

func (u *userRepository) GetByID(ctx context.Context, ID uuid.UUID) (*domain.User, error) {
	log := slog.With(
		slog.String("repository", "user"),
		slog.String("func", "GetByID"),
	)

	log.Info("Initializing process of obtaining user by ID")

	var user *domain.User
	if err := u.db.WithContext(ctx).Where("id = ?", ID).First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			log.Warn("User not found")
			return nil, nil
		}

		log.Error("Failed to get user by id", slog.String("error", err.Error()))
		return nil, err
	}

	log.Info("Process of obtaining user by id executed successfully")
	return user, nil
}
