package router

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/kevinssheva/go-backend-template/internal/handler"
)

func registerHealthRoutes(r *mux.Router, h *handler.HealthHandler) {
	r.HandleFunc("/ping", h.Ping).Methods(http.MethodPost)
}
