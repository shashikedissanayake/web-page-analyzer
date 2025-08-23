package router

import (
	"github.com/shashikedissanayake/web-page-analyzer/server/config"
	"github.com/shashikedissanayake/web-page-analyzer/server/controller"
)

type Server struct {
	Config                *config.Configuration // Configuration
	HealthCheckController controller.IHealthCheckController
	ScraperController     controller.IScraperController
}

func CreateNewServer(
	config *config.Configuration,
	scraperController controller.IScraperController,
	healthCheckController controller.IHealthCheckController,
) *Server {
	server := &Server{
		Config:                config,
		HealthCheckController: healthCheckController,
		ScraperController:     scraperController,
	}
	return server
}
