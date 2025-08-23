package main

import (
	"fmt"
	"net/http"
	"os"

	"github.com/shashikedissanayake/web-page-analyzer/server/config"
	"github.com/shashikedissanayake/web-page-analyzer/server/controller"
	"github.com/shashikedissanayake/web-page-analyzer/server/router"
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

	// Services
	webpageAnalyzerService := service.CreateWebPageAnalyzerService()
	scraperService := service.CreateNewScraperService(webpageAnalyzerService)

	// Controllers
	scraperController := controller.CreateNewScraperController(
		scraperService, responseWriterUtil,
	)
	healthCheckController := controller.CreateNewHealthCheckController(
		responseWriterUtil,
	)

	server := router.CreateNewServer(
		config,
		scraperController,
		healthCheckController,
	)

	router := router.CreateNewRouter(server)

	port := config.Port
	logger.Info("Server started on port: ", port)
	http.ListenAndServe(fmt.Sprintf(":%v", port), router)
}
