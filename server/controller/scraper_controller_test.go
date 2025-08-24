package controller

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/shashikedissanayake/web-page-analyzer/server/model"
	"github.com/shashikedissanayake/web-page-analyzer/server/service"
	"github.com/shashikedissanayake/web-page-analyzer/server/utils"
	"github.com/stretchr/testify/suite"
	gomock "go.uber.org/mock/gomock"
)

type ScraperControllerSuite struct {
	suite.Suite
	controller     IScraperController
	responseWriter utils.IResponseWriter
	mockService    *service.MockIScraperService
}

func TestScraperControllerSuite(t *testing.T) {
	suite.Run(t, &ScraperControllerSuite{})
}

func (scs *ScraperControllerSuite) SetupTest() {
	scs.responseWriter = utils.CreateNewResponseWriter()
	scs.mockService = service.NewMockIScraperService(gomock.NewController(scs.T()))
	scs.controller = CreateNewScraperController(scs.mockService, scs.responseWriter)
}

func (scs *ScraperControllerSuite) TestScrapeWebPage() {
	testCases := []struct {
		name            string
		requestBody     string
		mockInput       string
		mockError       error
		mockOutput      *model.ScraperResponse
		responseCode    int
		responseMessage string
	}{
		{
			"Return error when passing invaild json payload",
			"{nil}",
			"",
			nil,
			nil,
			http.StatusBadRequest,
			"Invalid payload",
		},
		{
			"Return error when Passing invaild url with payload",
			"{\"url\": \"test\"}",
			"",
			nil,
			nil,
			http.StatusBadRequest,
			"Invalid url",
		},
		{
			"Return error when ScrapeWebPage failed",
			"{\"url\": \"https://google.com\"}",
			"https://google.com",
			errors.New("ScrapeWebPage failed"),
			nil,
			http.StatusUnprocessableEntity,
			"Failed to analyze web page",
		},
		{
			"Return success when ScrapeWebPage retuned success",
			"{\"url\": \"https://google.com\"}",
			"https://google.com",
			nil,
			&model.ScraperResponse{},
			http.StatusOK,
			"Success",
		},
	}

	for _, test := range testCases {
		scs.T().Run(test.name, func(t *testing.T) {
			if test.mockInput != "" {
				scs.mockService.EXPECT().ScrapeWebPage(test.mockInput).Return(
					test.mockOutput,
					test.mockError,
				).Times(1)
			}

			server := httptest.NewServer(http.HandlerFunc(scs.controller.ScrapeWebPage))

			resp, err := http.Post(server.URL, "application/json", strings.NewReader(test.requestBody))
			if err != nil {
				scs.T().Error(err)
			}
			defer resp.Body.Close()
			var response model.BaseResponse
			json.NewDecoder(resp.Body).Decode(&response)

			scs.Equal(test.responseCode, response.StatusCode)
			scs.Equal(test.responseMessage, response.Message)
		})
	}
}
