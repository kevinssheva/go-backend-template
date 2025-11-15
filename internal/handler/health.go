package handler

import (
	"encoding/json"
	"net/http"

	"github.com/kevinssheva/go-backend-template/internal/errs"
	"github.com/kevinssheva/go-backend-template/internal/handler/request"
	"github.com/kevinssheva/go-backend-template/internal/handler/response"
	"github.com/kevinssheva/go-backend-template/internal/service"
	"github.com/kevinssheva/go-backend-template/internal/validation"
	"go.uber.org/zap"
)

type HealthHandler struct {
	healthService service.HealthService
	logger        *zap.Logger
}

func NewHealthHandler(healthService service.HealthService, logger *zap.Logger) *HealthHandler {
	return &HealthHandler{
		healthService: healthService,
		logger:        logger,
	}
}

func (h *HealthHandler) Ping(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var req request.HealthCheckRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.logger.Error("failed to decode request", zap.Error(err))
		response.Error(w, errs.ErrInvalidJSON)
		return
	}

	if err := validation.ValidateStruct(req); err != nil {
		h.logger.Error("validation failed", zap.Error(err))
		response.Error(w, err)
		return
	}

	status, err := h.healthService.Ping(ctx, req.IncludeDB)
	if err != nil {
		h.logger.Error("ping failed", zap.Error(err))
		response.Error(w, err)
		return
	}

	res := response.HealthResponse{
		Message: status.Message,
	}

	response.Success(w, "pong", res)
}
