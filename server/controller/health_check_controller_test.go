package controller

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/shashikedissanayake/web-page-analyzer/server/model"
	"github.com/shashikedissanayake/web-page-analyzer/server/utils"
	"github.com/stretchr/testify/suite"
)

type HealthControllerSuite struct {
	suite.Suite
	controller     IHealthCheckController
	responseWriter utils.IResponseWriter
}

func TestHealthControllerSuite(t *testing.T) {
	suite.Run(t, &HealthControllerSuite{})
}

func (hcs *HealthControllerSuite) SetupTest() {
	hcs.responseWriter = utils.CreateNewResponseWriter()
	hcs.controller = CreateNewHealthCheckController(hcs.responseWriter)
}

func (hcs *HealthControllerSuite) TestHealthCheck() {
	testCases := []struct {
		name            string
		responseCode    int
		responseMessage string
	}{
		{"Return success when calls", http.StatusOK, "Success"},
	}

	for _, test := range testCases {
		hcs.T().Run(test.name, func(t *testing.T) {
			server := httptest.NewServer(http.HandlerFunc(hcs.controller.HealthCheck))

			resp, err := http.Get(server.URL)
			if err != nil {
				hcs.T().Error(err)
			}
			defer resp.Body.Close()
			var response model.BaseResponse
			json.NewDecoder(resp.Body).Decode(&response)

			hcs.Equal(test.responseCode, response.StatusCode)
			hcs.Equal(test.responseMessage, response.Message)
		})
	}
}
