package response

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/kevinssheva/go-backend-template/internal/errs"
)

type APIError struct {
	Code    string      `json:"code"`
	Message string      `json:"message"`
	Details interface{} `json:"details,omitempty"`
}

type APIResponse struct {
	Success bool        `json:"success"`
	Message string      `json:"message,omitempty"`
	Data    interface{} `json:"data,omitempty"`
	Error   *APIError   `json:"error,omitempty"`
	Meta    interface{} `json:"meta,omitempty"`
}

type PaginationMeta struct {
	Page       int `json:"page"`
	Size       int `json:"size"`
	TotalItems int `json:"total_items"`
	TotalPages int `json:"total_pages"`
}

func writeJSON(w http.ResponseWriter, status int, payload any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	if err := json.NewEncoder(w).Encode(payload); err != nil {
		http.Error(w, fmt.Sprintf("failed to encode response: %v", err), http.StatusInternalServerError)
	}
}

func Success(w http.ResponseWriter, message string, data interface{}) {
	writeJSON(w, http.StatusOK, APIResponse{
		Success: true,
		Message: message,
		Data:    data,
	})
}

func Error(w http.ResponseWriter, err error) {
	se := errs.AsServiceError(err)

	writeJSON(w, se.Status, APIResponse{
		Success: false,
		Error: &APIError{
			Code:    se.Code,
			Message: se.Message,
			Details: se.Details,
		},
	})
}

func Pagination(w http.ResponseWriter, message string, data interface{}, meta PaginationMeta) {
	writeJSON(w, http.StatusOK, APIResponse{
		Success: true,
		Message: message,
		Data:    data,
		Meta:    meta,
	})
}
