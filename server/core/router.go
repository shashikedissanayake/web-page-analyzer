package core

import (
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
)

func CreateNewRouter(s *Server) *chi.Mux {
	router := chi.NewRouter()

	// Middlewares
	router.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"https://*", "http://*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowCredentials: true,
	}))
	router.Use(middleware.Logger)

	// Routes
	router.Post("/analyze", s.ScraperController.ScrapeWebPage)

	return router
}
