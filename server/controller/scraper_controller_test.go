package controller

import (
	"testing"

	"github.com/shashikedissanayake/web-page-analyzer/server/service"
	"github.com/shashikedissanayake/web-page-analyzer/server/utils"
	"github.com/stretchr/testify/suite"
	gomock "go.uber.org/mock/gomock"
)

type ScraperControllerSuite struct {
	suite.Suite
	controller         IScraperController
	mockService        *service.MockIScraperService
	mockResponseWriter *utils.MockIResponseWriter
}

func TestScraperControllerSuite(t *testing.T) {
	suite.Run(t, &ScraperControllerSuite{})
}

func (scs *ScraperControllerSuite) SetupTest() {
	scs.mockResponseWriter = utils.NewMockIResponseWriter(gomock.NewController(scs.T()))
	scs.mockService = service.NewMockIScraperService(gomock.NewController(scs.T()))
	scs.controller = CreateNewScraperController(scs.mockService, scs.mockResponseWriter)
}

func (scs *ScraperControllerSuite) TestScrapeWebPage() {
	// type RetFunctionCallArgumants struct {
	// 	count   int
	// 	code    int
	// 	message string
	// }

	// type RetScrapeWebPage struct {
	// 	count    int
	// 	input    string
	// 	response *model.ScraperResponse
	// 	err      error
	// }

	// var testCases = []struct {
	// 	name                      string
	// 	resWriter                 http.ResponseWriter
	// 	reqReader                 http.Request
	// 	invokeSendErrorResponse   RetFunctionCallArgumants
	// 	invokeSendSuccessResponse RetFunctionCallArgumants
	// 	invokeScrapeWebPage       RetScrapeWebPage
	// }{
	// 	{
	// 		"Call with invalid JSON",
	// 		http.ResponseWriter{},
	// 		http.Request{Body: "{test}"},
	// 		RetFunctionCallArgumants{1, http.StatusUnprocessableEntity, "Failed to decode json"},
	// 		RetFunctionCallArgumants{0, 0, ""},
	// 		RetScrapeWebPage{0, "", nil, nil},
	// 	},
	// }

	// for _, test := range testCases {
	// 	scs.T().Run(test.name, func(t *testing.T) {
	// 		scs.mockResponseWriter.EXPECT().SendErrorResponse(
	// 			gomock.Any(),
	// 			test.invokeSendErrorResponse.code,
	// 			test.invokeSendErrorResponse.message,
	// 			gomock.Any(),
	// 		).Times(test.invokeSendErrorResponse.count)

	// 		scs.mockResponseWriter.EXPECT().SendSuccessResponse(
	// 			gomock.Any(),
	// 			test.invokeSendSuccessResponse.code,
	// 			test.invokeSendSuccessResponse.message,
	// 			gomock.Any(),
	// 		).Times(test.invokeSendSuccessResponse.count)

	// 		scs.mockService.EXPECT().ScrapeWebPage(
	// 			test.invokeScrapeWebPage.input,
	// 		).Return(
	// 			test.invokeScrapeWebPage.response,
	// 			test.invokeScrapeWebPage.err,
	// 		).Times(test.invokeScrapeWebPage.count)

	// 		scs.controller.ScrapeWebPage(test.resWriter, &test.reqReader)

	// 	})
	// }
}
