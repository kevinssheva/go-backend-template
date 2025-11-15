package handler

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/kevinssheva/go-backend-template/internal/domain"
	"github.com/kevinssheva/go-backend-template/internal/errs"
	"github.com/kevinssheva/go-backend-template/internal/handler/response"
	"github.com/kevinssheva/go-backend-template/internal/service/mocks"
	"github.com/stretchr/testify/mock"
	"go.uber.org/zap"
)

func TestHealthHandler_Ping(t *testing.T) {
	logger, _ := zap.NewDevelopment()

	tests := []struct {
		name           string
		requestBody    string
		setupMock      func(*mocks.HealthService)
		expectedStatus int
		checkResponse  func(t *testing.T, resp response.APIResponse)
	}{
		{
			name:        "successful ping without db check",
			requestBody: `{"include_db": false}`,
			setupMock: func(m *mocks.HealthService) {
				m.EXPECT().
					Ping(mock.Anything, false).
					Return(&domain.HealthStatus{Message: "pong"}, nil).
					Once()
			},
			expectedStatus: http.StatusOK,
			checkResponse: func(t *testing.T, resp response.APIResponse) {
				if !resp.Success {
					t.Error("expected success to be true")
				}
				if resp.Message != "pong" {
					t.Errorf("expected message 'pong', got %q", resp.Message)
				}
				data, ok := resp.Data.(map[string]interface{})
				if !ok {
					t.Fatal("expected data to be a map")
				}
				if data["message"] != "pong" {
					t.Errorf("expected data.message 'pong', got %q", data["message"])
				}
			},
		},
		{
			name:        "successful ping with db check",
			requestBody: `{"include_db": true}`,
			setupMock: func(m *mocks.HealthService) {
				m.EXPECT().
					Ping(mock.Anything, true).
					Return(&domain.HealthStatus{Message: "pong"}, nil).
					Once()
			},
			expectedStatus: http.StatusOK,
			checkResponse: func(t *testing.T, resp response.APIResponse) {
				if !resp.Success {
					t.Error("expected success to be true")
				}
			},
		},
		{
			name:        "database check failure",
			requestBody: `{"include_db": true}`,
			setupMock: func(m *mocks.HealthService) {
				m.EXPECT().
					Ping(mock.Anything, true).
					Return(nil, errs.New(
						"database_unavailable",
						503,
						"Database is unavailable",
						errs.WithError(errors.New("connection refused")),
					)).
					Once()
			},
			expectedStatus: http.StatusServiceUnavailable,
			checkResponse: func(t *testing.T, resp response.APIResponse) {
				if resp.Success {
					t.Error("expected success to be false")
				}
				if resp.Error == nil {
					t.Fatal("expected error to be present")
				}
				if resp.Error.Code != "database_unavailable" {
					t.Errorf("expected error code 'database_unavailable', got %q", resp.Error.Code)
				}
			},
		},
		{
			name:           "invalid request body",
			requestBody:    `{invalid json}`,
			setupMock:      func(m *mocks.HealthService) {},
			expectedStatus: http.StatusBadRequest,
			checkResponse: func(t *testing.T, resp response.APIResponse) {
				if resp.Success {
					t.Error("expected success to be false")
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockService := mocks.NewHealthService(t)
			tt.setupMock(mockService)

			handler := NewHealthHandler(mockService, logger)

			req := httptest.NewRequest(http.MethodPost, "/ping", bytes.NewBufferString(tt.requestBody))
			req.Header.Set("Content-Type", "application/json")
			w := httptest.NewRecorder()

			handler.Ping(w, req)

			if w.Code != tt.expectedStatus {
				t.Errorf("expected status %d, got %d", tt.expectedStatus, w.Code)
			}

			var resp response.APIResponse
			if err := json.NewDecoder(w.Body).Decode(&resp); err != nil {
				t.Fatalf("failed to decode response: %v", err)
			}

			if tt.checkResponse != nil {
				tt.checkResponse(t, resp)
			}
		})
	}
}
