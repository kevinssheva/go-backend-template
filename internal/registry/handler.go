package registry

import (
	"github.com/kevinssheva/go-backend-template/internal/handler"
	"go.uber.org/zap"
)

type Handlers struct {
	HealthHandler *handler.HealthHandler
}

func NewHandlers(services *Services, log *zap.Logger) *Handlers {
	return &Handlers{
		HealthHandler: handler.NewHealthHandler(services.HealthService, log),
	}
}
