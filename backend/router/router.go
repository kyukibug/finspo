package router

import (
	"com.fukubox/handlers"
	"github.com/go-chi/chi"
)

func SetupRoutes(r *chi.Mux) {
	r.Route("/categories", func(r chi.Router) {
		r.Get("/", handlers.GetCategories)
		r.Get("/{category_id}", handlers.GetCategoriesById)
	})
}
