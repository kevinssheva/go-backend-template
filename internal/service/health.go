package service

import (
	"context"

	"github.com/kevinssheva/go-backend-template/internal/domain"
	"github.com/kevinssheva/go-backend-template/internal/repository"
	"go.uber.org/zap"
)

type HealthService interface {
	Ping(ctx context.Context, includeDB bool) (*domain.HealthStatus, error)
}

type healthService struct {
	healthRepo repository.HealthRepository
	logger     *zap.Logger
}

func NewHealthService(healthRepo repository.HealthRepository, logger *zap.Logger) HealthService {
	return &healthService{
		healthRepo: healthRepo,
		logger:     logger,
	}
}

func (s *healthService) Ping(ctx context.Context, includeDB bool) (*domain.HealthStatus, error) {
	if includeDB {
		if err := s.healthRepo.CheckDB(ctx); err != nil {
			s.logger.Error("database health check failed", zap.Error(err))
			return nil, err
		}
	}

	return &domain.HealthStatus{
		Message: "pong",
	}, nil
}
