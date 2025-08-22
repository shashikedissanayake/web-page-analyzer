package controller

import (
	"encoding/json"
	"net/http"

	"github.com/shashikedissanayake/web-page-analyzer/server/model"
	"github.com/shashikedissanayake/web-page-analyzer/server/service"
	"github.com/shashikedissanayake/web-page-analyzer/server/utils"
	logger "github.com/sirupsen/logrus"
)

//go:generate mockgen -source=scraper-controller.go -destination=scraper-controller_mock.go -package=controller
type IScraperController interface {
	ScrapeWebPage(http.ResponseWriter, *http.Request)
}

type ScraperController struct {
	service        service.IScraperService
	responseWriter utils.IResponseWriter
}

func CreateNewScraperController(
	service service.IScraperService,
	responseWriter utils.IResponseWriter,
) IScraperController {
	return &ScraperController{
		service,
		responseWriter,
	}
}

func (sc *ScraperController) ScrapeWebPage(w http.ResponseWriter, r *http.Request) {
	var request model.ScraperRequest
	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		sc.responseWriter.SendErrorResponse(w, http.StatusUnprocessableEntity, "Failed to decode json", err.Error())
		return
	}
	logger.Info("Recieved request to scrape webpage with payload:", request)

	if !utils.IsValidURL(request.Url) {
		sc.responseWriter.SendErrorResponse(w, http.StatusBadRequest, "Invalid url", "")
		return
	}

	res, err := sc.service.ScrapeWebPage(request.Url)
	if err != nil {
		sc.responseWriter.SendErrorResponse(w, http.StatusUnprocessableEntity, "Failed to scrape web page", err.Error())
		return
	}
	sc.responseWriter.SendSuccessResponse(w, http.StatusOK, "Success", res)
}
