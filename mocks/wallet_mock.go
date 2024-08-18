// Code generated by MockGen. DO NOT EDIT.
// Source: wallet.go

// Package mocks is a generated GoMock package.
package mocks

import (
	context "context"
	reflect "reflect"

	domain "github.com/GSVillas/pic-pay-desafio/domain"
	gomock "github.com/golang/mock/gomock"
	uuid "github.com/google/uuid"
	echo "github.com/labstack/echo/v4"
)

// MockWalletHandler is a mock of WalletHandler interface.
type MockWalletHandler struct {
	ctrl     *gomock.Controller
	recorder *MockWalletHandlerMockRecorder
}

// MockWalletHandlerMockRecorder is the mock recorder for MockWalletHandler.
type MockWalletHandlerMockRecorder struct {
	mock *MockWalletHandler
}

// NewMockWalletHandler creates a new mock instance.
func NewMockWalletHandler(ctrl *gomock.Controller) *MockWalletHandler {
	mock := &MockWalletHandler{ctrl: ctrl}
	mock.recorder = &MockWalletHandlerMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockWalletHandler) EXPECT() *MockWalletHandlerMockRecorder {
	return m.recorder
}

// Create mocks base method.
func (m *MockWalletHandler) Create(arg0 echo.Context) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Create", arg0)
	ret0, _ := ret[0].(error)
	return ret0
}

// Create indicates an expected call of Create.
func (mr *MockWalletHandlerMockRecorder) Create(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Create", reflect.TypeOf((*MockWalletHandler)(nil).Create), arg0)
}

// MockWalletService is a mock of WalletService interface.
type MockWalletService struct {
	ctrl     *gomock.Controller
	recorder *MockWalletServiceMockRecorder
}

// MockWalletServiceMockRecorder is the mock recorder for MockWalletService.
type MockWalletServiceMockRecorder struct {
	mock *MockWalletService
}

// NewMockWalletService creates a new mock instance.
func NewMockWalletService(ctrl *gomock.Controller) *MockWalletService {
	mock := &MockWalletService{ctrl: ctrl}
	mock.recorder = &MockWalletServiceMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockWalletService) EXPECT() *MockWalletServiceMockRecorder {
	return m.recorder
}

// Create mocks base method.
func (m *MockWalletService) Create(ctx context.Context, payload *domain.WalletPayload) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Create", ctx, payload)
	ret0, _ := ret[0].(error)
	return ret0
}

// Create indicates an expected call of Create.
func (mr *MockWalletServiceMockRecorder) Create(ctx, payload interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Create", reflect.TypeOf((*MockWalletService)(nil).Create), ctx, payload)
}

// MockWalletRepository is a mock of WalletRepository interface.
type MockWalletRepository struct {
	ctrl     *gomock.Controller
	recorder *MockWalletRepositoryMockRecorder
}

// MockWalletRepositoryMockRecorder is the mock recorder for MockWalletRepository.
type MockWalletRepositoryMockRecorder struct {
	mock *MockWalletRepository
}

// NewMockWalletRepository creates a new mock instance.
func NewMockWalletRepository(ctrl *gomock.Controller) *MockWalletRepository {
	mock := &MockWalletRepository{ctrl: ctrl}
	mock.recorder = &MockWalletRepositoryMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockWalletRepository) EXPECT() *MockWalletRepositoryMockRecorder {
	return m.recorder
}

// Create mocks base method.
func (m *MockWalletRepository) Create(ctx context.Context, wallet *domain.Wallet) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Create", ctx, wallet)
	ret0, _ := ret[0].(error)
	return ret0
}

// Create indicates an expected call of Create.
func (mr *MockWalletRepositoryMockRecorder) Create(ctx, wallet interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Create", reflect.TypeOf((*MockWalletRepository)(nil).Create), ctx, wallet)
}

// GetByUserID mocks base method.
func (m *MockWalletRepository) GetByUserID(ctx context.Context, userID uuid.UUID) (*domain.Wallet, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetByUserID", ctx, userID)
	ret0, _ := ret[0].(*domain.Wallet)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetByUserID indicates an expected call of GetByUserID.
func (mr *MockWalletRepositoryMockRecorder) GetByUserID(ctx, userID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetByUserID", reflect.TypeOf((*MockWalletRepository)(nil).GetByUserID), ctx, userID)
}