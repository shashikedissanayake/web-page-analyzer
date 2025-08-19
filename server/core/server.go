package core

import (
	"github.com/shashikedissanayake/web-page-analyzer/server/config"
	"github.com/shashikedissanayake/web-page-analyzer/server/controller"
	"github.com/shashikedissanayake/web-page-analyzer/server/service"
)

type Server struct {
	Config            *config.Configuration // Configuration
	ScraperController controller.IScraperController
	ScraperService    service.IScraperService
}

func CreateNewServer(
	config *config.Configuration,
	scraperService service.IScraperService,
	scraperController controller.IScraperController,
) *Server {
	server := &Server{
		Config:            config,
		ScraperController: scraperController,
		ScraperService:    scraperService,
	}
	return server
}
