package handler

import (
	"bytes"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/GSVillas/pic-pay-desafio/domain"
	"github.com/GSVillas/pic-pay-desafio/mocks"
	"github.com/GSVillas/pic-pay-desafio/utils"
	"github.com/golang/mock/gomock"
	jsoniter "github.com/json-iterator/go"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
)

func TestUserHandler_Create_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	userServiceMock := mocks.NewMockUserService(ctrl)

	handler := &userHandler{
		userService: userServiceMock,
	}

	payload := domain.UserPayload{
		Name:            "Test User",
		CPF:             "077.351.310-89",
		Email:           "test@example.com",
		ConfirmEmail:    "test@example.com",
		Password:        utils.Password,
		ConfirmPassword: utils.Password,
	}

	jsonPayload, _ := jsoniter.Marshal(payload)

	req := httptest.NewRequest(http.MethodPost, "/v1/users", bytes.NewReader(jsonPayload))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()

	e := echo.New()
	ctx := e.NewContext(req, rec)

	userServiceMock.EXPECT().Create(gomock.Any(), &payload).Return(nil)

	err := handler.Create(ctx)

	assert.NoError(t, err)
	assert.Equal(t, http.StatusCreated, rec.Code)
}

func TestUserHandler_Create_WhenValidationFails_ShouldReturnBadRequest(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	userServiceMock := mocks.NewMockUserService(ctrl)

	handler := &userHandler{
		userService: userServiceMock,
	}

	payload := domain.UserPayload{
		Name:            "",
		CPF:             "077.351.310-89",
		Email:           "test@example.com",
		ConfirmEmail:    "test@example.com",
		Password:        utils.Password,
		ConfirmPassword: utils.Password,
	}

	PayloadJSON, _ := jsoniter.Marshal(payload)

	req := httptest.NewRequest(http.MethodPost, "/v1/users", bytes.NewReader(PayloadJSON))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()

	e := echo.New()
	ctx := e.NewContext(req, rec)

	err := handler.Create(ctx)

	assert.NoError(t, err)
	assert.Equal(t, http.StatusBadRequest, rec.Code)

	expectedResponse := map[string]any{
		"status": 400,
		"title":  "Validation Failed",
		"detail": "One or more fields failed validation",
		"errors": map[string]any{
			"name": "This field is required",
		},
	}

	var actualResponse map[string]any
	err = jsoniter.NewDecoder(rec.Body).Decode(&actualResponse)
	assert.NoError(t, err)

	if actualStatus, ok := actualResponse["status"].(float64); ok {
		actualResponse["status"] = int(actualStatus)
	}

	assert.Equal(t, expectedResponse, actualResponse)
}

func TestUserHandler_Create_WhenServiceFails_ShouldReturnInternalServerError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	userServiceMock := mocks.NewMockUserService(ctrl)

	handler := &userHandler{
		userService: userServiceMock,
	}

	payload := domain.UserPayload{
		Name:            "Test User",
		CPF:             "077.351.310-89",
		Email:           "test@example.com",
		ConfirmEmail:    "test@example.com",
		Password:        utils.Password,
		ConfirmPassword: utils.Password,
	}

	jsonPayload, _ := jsoniter.Marshal(payload)

	req := httptest.NewRequest(http.MethodPost, "/v1/users", bytes.NewReader(jsonPayload))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	e := echo.New()
	ctx := e.NewContext(req, rec)

	userServiceMock.EXPECT().Create(gomock.Any(), &payload).Return(errors.New("service error"))

	err := handler.Create(ctx)

	assert.NoError(t, err)
	assert.Equal(t, http.StatusInternalServerError, rec.Code)

	expectedResponse := map[string]any{
		"status": 500,
		"title":  "Internal Server Error",
		"detail": "Failed to process the payload",
	}

	var actualResponse map[string]any
	err = jsoniter.NewDecoder(rec.Body).Decode(&actualResponse)
	assert.NoError(t, err)

	if actualStatus, ok := actualResponse["status"].(float64); ok {
		actualResponse["status"] = int(actualStatus)
	}

	assert.Equal(t, expectedResponse, actualResponse)
}
