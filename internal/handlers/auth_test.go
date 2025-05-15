package handlers

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"

	"github.com/SanExpett/diploma/internal/domain"
	"github.com/SanExpett/diploma/internal/handlers/mocks"
	"github.com/SanExpett/diploma/internal/metrics"
	session "github.com/SanExpett/diploma/internal/session/proto"
)

func TestAuthPageHandlers_Login(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUsersClient := mocks.NewMockUsersClient(ctrl)
	mockSessionsClient := mocks.NewMockSessionsClient(ctrl)

	var usersClient session.UsersClient = mockUsersClient
	var sessionsClient session.SessionsClient = mockSessionsClient

	logger := zap.NewNop().Sugar()
	metrics := metrics.NewHttpMetrics()

	handler := NewAuthPageHandlers(&usersClient, &sessionsClient, metrics, logger)

	tests := []struct {
		name         string
		input        domain.UserSignUp
		setupMocks   func()
		expectedCode int
	}{
		{
			name: "Успешный вход",
			input: domain.UserSignUp{
				Email:    "test@test.com",
				Password: "password123",
			},
			setupMocks: func() {
				mockUsersClient.EXPECT().HasUser(gomock.Any(), &session.HasUserRequest{
					Login:    "test@test.com",
					Password: "password123",
				}).Return(&session.HasUserResponse{Has: false}, nil)

				mockUsersClient.EXPECT().GetUser(gomock.Any(), &session.GetUserRequest{
					Login: "test@test.com",
				}).Return(&session.GetUserResponse{
					User: &session.User{
						Email:   "test@test.com",
						IsAdmin: false,
						Version: 1,
						Uuid:    "test-uuid",
					},
				}, nil)

				mockSessionsClient.EXPECT().Add(gomock.Any(), gomock.Any()).Return(&session.AddResponse{}, nil)
			},
			expectedCode: http.StatusOK,
		},
		{
			name: "Некорректный email",
			input: domain.UserSignUp{
				Email:    "invalid-email",
				Password: "password123",
			},
			setupMocks:   func() {},
			expectedCode: http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setupMocks()

			body, _ := json.Marshal(tt.input)
			req := httptest.NewRequest(http.MethodPost, "/api/auth/login", bytes.NewReader(body))
			w := httptest.NewRecorder()

			handler.Login(w, req)

			assert.Equal(t, tt.expectedCode, w.Code)
		})
	}
}

func TestAuthPageHandlers_Logout(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUsersClient := mocks.NewMockUsersClient(ctrl)
	mockSessionsClient := mocks.NewMockSessionsClient(ctrl)

	var usersClient session.UsersClient = mockUsersClient
	var sessionsClient session.SessionsClient = mockSessionsClient

	logger := zap.NewNop().Sugar()
	metrics := metrics.NewHttpMetrics()

	handler := NewAuthPageHandlers(&usersClient, &sessionsClient, metrics, logger)

	tests := []struct {
		name         string
		setupCookie  func(req *http.Request)
		setupMocks   func()
		expectedCode int
	}{
		{
			name: "Успешный выход",
			setupCookie: func(req *http.Request) {
				cookie := &http.Cookie{
					Name:  "access",
					Value: "valid-token",
				}
				req.AddCookie(cookie)
			},
			setupMocks: func() {
				mockSessionsClient.EXPECT().DeleteSession(gomock.Any(), gomock.Any()).Return(&session.DeleteSessionResponse{}, nil)
				mockSessionsClient.EXPECT().GetVersion(gomock.Any(), gomock.Any()).Return(&session.GetVersionResponse{}, nil)
			},
			expectedCode: http.StatusOK,
		},
		{
			name:         "Отсутствует токен",
			setupCookie:  func(req *http.Request) {},
			setupMocks:   func() {},
			expectedCode: http.StatusUnauthorized,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setupMocks()

			req := httptest.NewRequest(http.MethodPost, "/api/auth/logout", nil)
			tt.setupCookie(req)
			w := httptest.NewRecorder()

			handler.Logout(w, req)

			assert.Equal(t, tt.expectedCode, w.Code)
		})
	}
}
