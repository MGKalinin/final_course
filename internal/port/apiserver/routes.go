package apiserver

import (
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

// SetupRoutes регистрирует маршруты API
func (h *RESTHandler) SetupRoutes(r chi.Router) {
	r.Group(func(r chi.Router) {
		r.Use(
			middleware.Logger,
			middleware.Recoverer,
			middleware.RequestID,
		)

		r.Get("/coins", h.getCoins)
		r.Get("/coins/max", h.getMaxRate)
		r.Get("/coins/min", h.getMinRate)
		r.Get("/coins/avg", h.getAvgRate)
	})
}
