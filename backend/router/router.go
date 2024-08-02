package router

import (
	"com.fukubox/handlers"
	"github.com/go-chi/chi"
)

func SetupRoutes(r *chi.Mux) {
	r.Route("/clothes", func(r chi.Router) {
		r.Get("/", handlers.GetClothes)
		r.Get("/{id}", handlers.GetClothesById)
		r.Post("/", handlers.CreateClothes)
		r.Patch("/{id}", handlers.UpdateClothes)
		r.Delete("/{id}", handlers.DeleteClothes)
	})

	r.Route("/categories", func(r chi.Router) {
		r.Get("/", handlers.GetCategories)
		r.Get("/{id}", handlers.GetCategoriesById)
		r.Post("/", handlers.CreateCategory)
		r.Patch("/{id}", handlers.UpdateCategory)
		r.Delete("/{id}", handlers.DeleteCategory)
	})

	r.Route("/tags", func(r chi.Router) {
		r.Get("/", handlers.GetTags)
		r.Get("/{id}", handlers.GetClothesById)
		r.Post("/", handlers.CreateTag)
		r.Patch("/{id}", handlers.UpdateTag)
		r.Delete("/{id}", handlers.DeleteTag)
	})
}
