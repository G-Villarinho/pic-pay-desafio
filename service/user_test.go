package service

import (
	"context"
	"errors"
	"testing"

	"github.com/GSVillas/pic-pay-desafio/domain"
	"github.com/GSVillas/pic-pay-desafio/mocks"
	"github.com/GSVillas/pic-pay-desafio/utils"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func TestUserService_Create_WhenUserAlreadyExistsByEmail_ShouldReturnErrEmailAlreadyRegister(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	userRepositoryMock := mocks.NewMockUserRepository(ctrl)
	sessionServiceMock := mocks.NewMockSessionService(ctrl)

	userService := &userService{
		userRepository: userRepositoryMock,
		sessionService: sessionServiceMock,
	}

	payload := &domain.UserPayload{
		Name:            "Test User",
		CPF:             "12345678901",
		Email:           "test@example.com",
		ConfirmEmail:    "test@example.com",
		Password:        "password123",
		ConfirmPassword: "password123",
	}

	existingUser := &domain.User{}
	userRepositoryMock.EXPECT().GetByEmail(gomock.Any(), payload.Email).Return(existingUser, nil)

	err := userService.Create(context.Background(), payload)

	assert.ErrorIs(t, err, domain.ErrEmailAlreadyRegister)
}

func TestUserService_Create_WhenUserAlreadyExistsByCPF_ShouldReturnErrCPFAlreadyRegister(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	userRepositoryMock := mocks.NewMockUserRepository(ctrl)
	sessionServiceMock := mocks.NewMockSessionService(ctrl)

	userService := &userService{
		userRepository: userRepositoryMock,
		sessionService: sessionServiceMock,
	}

	payload := &domain.UserPayload{
		Name:            "Test User",
		CPF:             "12345678901",
		Email:           "test@example.com",
		ConfirmEmail:    "test@example.com",
		Password:        "password123",
		ConfirmPassword: "password123",
	}

	existingUser := &domain.User{}
	userRepositoryMock.EXPECT().GetByEmail(gomock.Any(), payload.Email).Return(nil, nil)
	userRepositoryMock.EXPECT().GetByCPF(gomock.Any(), payload.CPF).Return(existingUser, nil)

	err := userService.Create(context.Background(), payload)

	assert.ErrorIs(t, err, domain.ErrCPFAlreadyRegister)
}

func TestUserService_Create_WhenCreateUserFails_ShouldReturnError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	userRepositoryMock := mocks.NewMockUserRepository(ctrl)
	sessionServiceMock := mocks.NewMockSessionService(ctrl)

	userService := &userService{
		userRepository: userRepositoryMock,
		sessionService: sessionServiceMock,
	}

	payload := &domain.UserPayload{
		Name:            "Test User",
		CPF:             "12345678901",
		Email:           "test@example.com",
		ConfirmEmail:    "test@example.com",
		Password:        "password123",
		ConfirmPassword: "password123",
	}

	userRepositoryMock.EXPECT().GetByEmail(gomock.Any(), payload.Email).Return(nil, nil)
	userRepositoryMock.EXPECT().GetByCPF(gomock.Any(), payload.CPF).Return(nil, nil)

	userRepositoryMock.EXPECT().Create(gomock.Any(), gomock.Any()).Return(errors.New("db error"))

	err := userService.Create(context.Background(), payload)

	assert.ErrorIs(t, err, domain.ErrCreateUser)
}

func TestUserService_Create_WhenHashingPasswordFails_ShouldReturnErrHashingPassword(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	userRepositoryMock := mocks.NewMockUserRepository(ctrl)
	sessionServiceMock := mocks.NewMockSessionService(ctrl)

	userService := &userService{
		userRepository: userRepositoryMock,
		sessionService: sessionServiceMock,
	}

	payload := &domain.UserPayload{
		Name:            "Test User",
		CPF:             "12345678901",
		Email:           "test@example.com",
		ConfirmEmail:    "test@example.com",
		Password:        utils.LargeString,
		ConfirmPassword: utils.LargeString,
	}

	userRepositoryMock.EXPECT().GetByEmail(gomock.Any(), payload.Email).Return(nil, nil)
	userRepositoryMock.EXPECT().GetByCPF(gomock.Any(), payload.CPF).Return(nil, nil)

	err := userService.Create(context.Background(), payload)

	assert.ErrorIs(t, err, domain.ErrHashingPassword)
}

func TestUserService_Create_WhenSuccess_ShouldReturnNil(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	userRepositoryMock := mocks.NewMockUserRepository(ctrl)
	sessionServiceMock := mocks.NewMockSessionService(ctrl)

	userService := &userService{
		userRepository: userRepositoryMock,
		sessionService: sessionServiceMock,
	}

	payload := &domain.UserPayload{
		Name:            "Test User",
		CPF:             "12345678901",
		Email:           "test@example.com",
		ConfirmEmail:    "test@example.com",
		Password:        "password123",
		ConfirmPassword: "password123",
	}

	userRepositoryMock.EXPECT().GetByEmail(gomock.Any(), payload.Email).Return(nil, nil)
	userRepositoryMock.EXPECT().GetByCPF(gomock.Any(), payload.CPF).Return(nil, nil)

	userRepositoryMock.EXPECT().Create(gomock.Any(), gomock.Any()).Return(nil)

	err := userService.Create(context.Background(), payload)

	assert.NoError(t, err)
}

func TestUserService_SignIn_WhenUserNotFound_ShouldReturnErrUserNotFound(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	userRepositoryMock := mocks.NewMockUserRepository(ctrl)
	sessionServiceMock := mocks.NewMockSessionService(ctrl)

	userService := &userService{
		userRepository: userRepositoryMock,
		sessionService: sessionServiceMock,
	}

	payload := &domain.SignInPayload{
		Email:    "test@example.com",
		Password: "password123",
	}

	userRepositoryMock.EXPECT().GetByEmail(gomock.Any(), payload.Email).Return(nil, nil)

	_, err := userService.SignIn(context.Background(), payload)

	assert.ErrorIs(t, err, domain.ErrUserNotFound)
}

func TestUserService_SignIn_WhenPasswordIsInvalid_ShouldReturnErrInvalidPassword(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	userRepositoryMock := mocks.NewMockUserRepository(ctrl)
	sessionServiceMock := mocks.NewMockSessionService(ctrl)

	userService := &userService{
		userRepository: userRepositoryMock,
		sessionService: sessionServiceMock,
	}

	payload := &domain.SignInPayload{
		Email:    "test@example.com",
		Password: "Teste@123",
	}

	user := &domain.User{
		Email:        payload.Email,
		PasswordHash: "wrong_password",
	}

	userRepositoryMock.EXPECT().GetByEmail(gomock.Any(), payload.Email).Return(user, nil)

	_, err := userService.SignIn(context.Background(), payload)

	assert.ErrorIs(t, err, domain.ErrInvalidPassword)
}

func TestUserService_SignIn_WhenSuccess_ShouldReturnSignInResponseAndNil(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	userRepositoryMock := mocks.NewMockUserRepository(ctrl)
	sessionServiceMock := mocks.NewMockSessionService(ctrl)

	userService := &userService{
		userRepository: userRepositoryMock,
		sessionService: sessionServiceMock,
	}

	payload := &domain.SignInPayload{
		Email:    "test@example.com",
		Password: utils.Password,
	}

	user := &domain.User{
		Email:        payload.Email,
		PasswordHash: utils.PasswordHash,
	}

	userRepositoryMock.EXPECT().GetByEmail(gomock.Any(), payload.Email).Return(user, nil)
	sessionServiceMock.EXPECT().Create(gomock.Any(), user).Return("validtoken", nil)

	response, err := userService.SignIn(context.Background(), payload)

	assert.NoError(t, err)
	assert.NotNil(t, response)
	assert.Equal(t, "validtoken", response.Token)
}

func TestUserService_SignIn_WhenCreateSessionFails_ShouldReturnErrCreateSession(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	userRepositoryMock := mocks.NewMockUserRepository(ctrl)
	sessionServiceMock := mocks.NewMockSessionService(ctrl)

	userService := &userService{
		userRepository: userRepositoryMock,
		sessionService: sessionServiceMock,
	}

	payload := &domain.SignInPayload{
		Email:    "test@example.com",
		Password: utils.Password,
	}

	user := &domain.User{
		Email:        payload.Email,
		PasswordHash: utils.PasswordHash,
	}

	userRepositoryMock.EXPECT().GetByEmail(gomock.Any(), payload.Email).Return(user, nil)

	sessionServiceMock.EXPECT().Create(gomock.Any(), user).Return("", errors.New("session error"))

	_, err := userService.SignIn(context.Background(), payload)

	assert.ErrorIs(t, err, domain.ErrCreateSession)
}
