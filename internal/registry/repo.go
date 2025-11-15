package registry

import (
	"database/sql"

	"github.com/kevinssheva/go-backend-template/internal/repository"
	"go.uber.org/zap"
)

type Repos struct {
	HealthRepo repository.HealthRepository
}

func NewRepos(db *sql.DB, log *zap.Logger) *Repos {
	return &Repos{
		HealthRepo: repository.NewHealthRepository(db),
	}
}
