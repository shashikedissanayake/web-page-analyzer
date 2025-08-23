package router

import (
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"

	custom_middlewares "github.com/shashikedissanayake/web-page-analyzer/server/middleware"
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

	// Public routes
	router.Get("/health_check", s.HealthCheckController.HealthCheck)

	// Authenticated routes
	router.Group(func(r chi.Router) {
		r.Use(custom_middlewares.AuthMiddleware)

		r.Post("/analyze", s.ScraperController.ScrapeWebPage)
	})

	return router
}
