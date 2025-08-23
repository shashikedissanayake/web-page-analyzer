package service

import (
	"strings"
	"sync"
	"testing"

	"github.com/shashikedissanayake/web-page-analyzer/server/model"
	"github.com/stretchr/testify/suite"
	"golang.org/x/net/html"
)

type WebpageAnalyzerServiceSuite struct {
	suite.Suite
	service WebPageAnalyzerService
}

func TestWebpageAnalyzerService(t *testing.T) {
	suite.Run(t, &WebpageAnalyzerServiceSuite{})
}

func (wass *WebpageAnalyzerServiceSuite) SetupTest() {
	wass.service = WebPageAnalyzerService{}
}

func (wass *WebpageAnalyzerServiceSuite) TestResultsReader() {
	var testCases = []struct {
		name         string
		initialLinks map[string]*model.Links
		input        model.LinkDetails
		output       map[string]*model.Links
	}{
		{
			"Insert value into empty list",
			map[string]*model.Links{},
			model.LinkDetails{LinkType: model.INTERNAL, URL: "/test", IsAccessible: true},
			map[string]*model.Links{},
		},
		{
			"Push data for a existing value internal",
			map[string]*model.Links{
				"/test":              {Count: 1},
				"https://google.com": {Count: 1},
			},
			model.LinkDetails{LinkType: model.INTERNAL, URL: "/test", IsAccessible: true},
			map[string]*model.Links{
				"/test":              {Count: 1, LinkType: model.INTERNAL, IsAccessible: true},
				"https://google.com": {Count: 1},
			},
		},
		{
			"Push data for a existing value external",
			map[string]*model.Links{
				"/test":              {Count: 1},
				"https://google.com": {Count: 1},
			},
			model.LinkDetails{LinkType: model.EXTERNAL, URL: "https://google.com", IsAccessible: false},
			map[string]*model.Links{
				"/test":              {Count: 1},
				"https://google.com": {Count: 1, LinkType: model.EXTERNAL, IsAccessible: false},
			},
		},
	}

	for _, test := range testCases {
		wass.T().Run(test.name, func(t *testing.T) {
			links := test.initialLinks
			channel := make(chan model.LinkDetails)
			wg := sync.WaitGroup{}

			wg.Add(1)
			go func(wg *sync.WaitGroup) {
				wass.service.resultsReader(links, channel)
				wg.Done()
			}(&wg)

			channel <- test.input
			close(channel)
			wg.Wait()

			for key, val := range test.output {
				wass.Equal(*val, *links[key])
			}
		})
	}
}

func (wass *WebpageAnalyzerServiceSuite) TestGetHtmlVersion() {
	var testCases = []struct {
		name          string
		docTypeString string
		output        string
	}{
		{"Parsing html 5 doctype", "<!DOCTYPE html>", "HTML 5"},
		{"Parsing html 4 doctype", "<!DOCTYPE HTML PUBLIC \"-//W3C//DTD HTML 4.01 Transitional//EN\" \"http://www.w3.org/TR/html4/loose.dtd\">", "HTML 4.01 Transitional//EN"},
		{"Parsing xhtml 1 doctype", "<!DOCTYPE html PUBLIC \"-//W3C//DTD XHTML 1.1//EN\" \"http://www.w3.org/TR/xhtml11/DTD/xhtml11.dtd\">", "XHTML 1.1//EN"},
	}

	for _, test := range testCases {
		wass.T().Run(test.name, func(t *testing.T) {
			tokenizor := html.NewTokenizer(strings.NewReader(test.docTypeString))
			tokenizor.Next()
			token := tokenizor.Token()
			res := wass.service.getHTMLVersion(&token)

			wass.Equal(test.output, res)
		})
	}
}
