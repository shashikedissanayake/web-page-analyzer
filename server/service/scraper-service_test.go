package service

import (
	"errors"
	"testing"

	"github.com/shashikedissanayake/web-page-analyzer/server/model"
	"github.com/shashikedissanayake/web-page-analyzer/server/utils"
	"github.com/stretchr/testify/suite"
	"go.uber.org/mock/gomock"
)

type ScraperServiceSuite struct {
	suite.Suite
	service     IScraperService
	mockScraper *utils.MockIWebPageAnalyzer
}

func TestScraperServiceSuite(t *testing.T) {
	suite.Run(t, &ScraperServiceSuite{})
}

func (sss *ScraperServiceSuite) SetupTest() {
	sss.mockScraper = utils.NewMockIWebPageAnalyzer(gomock.NewController(sss.T()))
	sss.service = CreateNewScraperService(sss.mockScraper)
}

func (sss *ScraperServiceSuite) TestScrapeWebPage() {
	var testCases = []struct {
		name           string
		input          string
		mockOutput     *utils.Response
		mockError      error
		expectedOutput *model.ScraperResponse
		expectedError  error
	}{
		{
			"ScrapeWebPage function failed with an error",
			"test",
			nil,
			errors.New("ScrapeWebPage failed"),
			nil,
			errors.New("ScrapeWebPage failed"),
		},
		{
			"ScrapeWebPage return valid results",
			"test",
			&utils.Response{
				HtmlVersion: "Html5",
				Title:       "test title",
				Links: map[string]*utils.Links{
					"/test":                  {LinkType: utils.INTERNAL, Count: 2, IsAccessible: false},
					"/login":                 {LinkType: utils.INTERNAL, Count: 1, IsAccessible: true},
					"#test1":                 {LinkType: utils.INTERNAL, Count: 2, IsAccessible: false},
					"http://google.com":      {LinkType: utils.EXTERNAL, Count: 1, IsAccessible: true},
					"http://google.com/test": {LinkType: utils.EXTERNAL, Count: 1, IsAccessible: false},
				},
				IsLoginForm: true,
				HeaderTags:  map[string]int{"h1": 4, "h2": 10},
			},
			nil,
			&model.ScraperResponse{
				Url:                           "test",
				HtmlVersion:                   "Html5",
				PageTitle:                     "test title",
				HeaderTagCount:                model.HeaderTagsCount{H1: 4, H2: 10},
				TotalLinkCount:                7,
				InternalLinkCount:             3,
				InternalInaccessibleLinkCount: 2,
				ExternalLinkCount:             2,
				ExternalInaccessibleLinkCount: 1,
				IsLoginForm:                   true,
			},
			nil,
		},
		{
			"ScrapeWebPage return valid results with empty internal map",
			"test",
			&utils.Response{
				HtmlVersion: "Html5",
				Title:       "test title",
				Links: map[string]*utils.Links{
					"http://google.com": {LinkType: utils.EXTERNAL, Count: 1, IsAccessible: true},
				},
				IsLoginForm: false,
				HeaderTags:  map[string]int{"h1": 4, "h2": 10},
			},
			nil,
			&model.ScraperResponse{
				Url:                           "test",
				HtmlVersion:                   "Html5",
				PageTitle:                     "test title",
				HeaderTagCount:                model.HeaderTagsCount{H1: 4, H2: 10},
				TotalLinkCount:                1,
				InternalLinkCount:             0,
				InternalInaccessibleLinkCount: 0,
				ExternalLinkCount:             1,
				ExternalInaccessibleLinkCount: 0,
				IsLoginForm:                   false,
			},
			nil,
		},
		{
			"ScrapeWebPage return valid results with empty external map",
			"test",
			&utils.Response{
				HtmlVersion: "Html5",
				Title:       "test title",
				Links: map[string]*utils.Links{
					"/test": {LinkType: utils.INTERNAL, Count: 1, IsAccessible: false},
				},
				IsLoginForm: false,
				HeaderTags:  map[string]int{"h1": 4, "h2": 10},
			},
			nil,
			&model.ScraperResponse{
				Url:                           "test",
				HtmlVersion:                   "Html5",
				PageTitle:                     "test title",
				HeaderTagCount:                model.HeaderTagsCount{H1: 4, H2: 10},
				TotalLinkCount:                1,
				InternalLinkCount:             1,
				InternalInaccessibleLinkCount: 1,
				ExternalLinkCount:             0,
				ExternalInaccessibleLinkCount: 0,
				IsLoginForm:                   false,
			},
			nil,
		},
		{
			"ScrapeWebPage return valid results with empty header list",
			"test",
			&utils.Response{
				HtmlVersion: "Html5",
				Title:       "test title",
				Links: map[string]*utils.Links{
					"/test":             {LinkType: utils.INTERNAL, Count: 2, IsAccessible: false},
					"http://google.com": {LinkType: utils.EXTERNAL, Count: 1, IsAccessible: true},
				},
				IsLoginForm: false,
				HeaderTags:  map[string]int{},
			},
			nil,
			&model.ScraperResponse{
				Url:                           "test",
				HtmlVersion:                   "Html5",
				PageTitle:                     "test title",
				HeaderTagCount:                model.HeaderTagsCount{},
				TotalLinkCount:                3,
				InternalLinkCount:             1,
				InternalInaccessibleLinkCount: 1,
				ExternalLinkCount:             1,
				ExternalInaccessibleLinkCount: 0,
				IsLoginForm:                   false,
			},
			nil,
		},
	}

	for _, test := range testCases {
		sss.T().Run(test.name, func(t *testing.T) {
			sss.mockScraper.EXPECT().AnalyzeWebPage(test.input).Times(1).Return(
				test.mockOutput,
				test.mockError,
			)

			res, err := sss.service.ScrapeWebPage(test.input)

			sss.Equal(res, test.expectedOutput)
			sss.Equal(err, test.expectedError)
		})
	}
}
