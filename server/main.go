package main

import (
	"fmt"
	"net/http"
	"os"

	"github.com/shashikedissanayake/web-page-analyzer/server/config"
	"github.com/shashikedissanayake/web-page-analyzer/server/controller"
	"github.com/shashikedissanayake/web-page-analyzer/server/core"
	"github.com/shashikedissanayake/web-page-analyzer/server/service"
	"github.com/shashikedissanayake/web-page-analyzer/server/utils"
	logger "github.com/sirupsen/logrus"
)

func init() {
	logger.SetFormatter(&logger.JSONFormatter{})
	logger.SetOutput(os.Stdout)
	logger.SetLevel(logger.InfoLevel)
}

func main() {
	config, err := config.NewConfig(".env")
	if err != nil {
		panic(err)
	}

	// Utils
	responseWriterUtil := utils.CreateNewResponseWriter()
	scraper := utils.CreateWebPageAnalyzer()

	// Services
	scraperService := service.CreateNewScraperService(scraper)

	// Controllers
	scraperController := controller.CreateNewScraperController(
		scraperService, responseWriterUtil,
	)
	healthCheckController := controller.CreateNewHealthCheckController(
		responseWriterUtil,
	)

	server := core.CreateNewServer(
		config,
		scraperService,
		scraperController,
		healthCheckController,
	)

	router := core.CreateNewRouter(server)

	port := config.Port
	logger.Info("Server started on port: ", port)
	http.ListenAndServe(fmt.Sprintf(":%v", port), router)
}
