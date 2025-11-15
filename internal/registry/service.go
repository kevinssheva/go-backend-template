package registry

import (
	"github.com/kevinssheva/go-backend-template/internal/service"
	"go.uber.org/zap"
)

type Services struct {
	HealthService service.HealthService
}

func NewServices(repos *Repos, log *zap.Logger) *Services {
	return &Services{
		HealthService: service.NewHealthService(repos.HealthRepo, log),
	}
}
