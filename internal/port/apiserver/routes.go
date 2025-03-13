package apiserver

import (
	"context"
	"final_course/internal/cases"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"net/http"
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

func (h *RESTHandler) getMaxRate(w http.ResponseWriter, r *http.Request) {
	titles := []string{"title1", "title2"} // Замените на реальные данные
	coins, err := h.service.GetMaxRate(context.Background(), titles)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	// Вернуть данные пользователю
}

func (h *RESTHandler) getMinRate(w http.ResponseWriter, r *http.Request) {
	titles := []string{"title1", "title2"} // Замените на реальные данные
	coins, err := h.service.GetMinRate(context.Background(), titles)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	// Вернуть данные пользователю
}

func (h *RESTHandler) getAvgRate(w http.ResponseWriter, r *http.Request) {
	titles := []string{"title1", "title2"} // Замените на реальные данные
	coins, err := h.service.GetAvgRate(context.Background(), titles)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	// Вернуть данные пользователю
}
