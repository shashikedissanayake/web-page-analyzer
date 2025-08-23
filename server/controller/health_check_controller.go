package controller

import (
	"net/http"

	"github.com/shashikedissanayake/web-page-analyzer/server/utils"
)

type IHealthCheckController interface {
	HealthCheck(http.ResponseWriter, *http.Request)
}

type HealthCheckController struct {
	responseWriter utils.IResponseWriter
}

func CreateNewHealthCheckController(
	responseWriter utils.IResponseWriter,
) IHealthCheckController {
	return &HealthCheckController{
		responseWriter,
	}
}

func (hc *HealthCheckController) HealthCheck(w http.ResponseWriter, r *http.Request) {
	hc.responseWriter.SendSuccessResponse(
		w, http.StatusOK, "Success", nil,
	)
}
