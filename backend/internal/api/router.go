package api

import (
	"github.com/go-chi/chi/v5"
)

func NewRouter() *chi.Mux {
	r := chi.NewRouter()
	
	r.Get("/api/health", HandleHealth)
	r.Get("/api/scene/example", HandleGetExample)
	r.Post("/api/scene/preview", HandlePreview)
	
	return r
}
