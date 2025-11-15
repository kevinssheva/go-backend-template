package router

import (
	"github.com/gorilla/mux"
	"github.com/kevinssheva/go-backend-template/internal/registry"
)

func NewRouter(handlers *registry.Handlers) *mux.Router {
	router := mux.NewRouter()

	registerHealthRoutes(router, handlers.HealthHandler)

	return router
}
