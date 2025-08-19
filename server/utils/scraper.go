package utils

import (
	"fmt"
	"net/http"
	"strings"

	"sync"

	"github.com/gocolly/colly"
	logger "github.com/sirupsen/logrus"
	"golang.org/x/net/html"
)

type Response struct {
	HtmlVersion string `json:"htmlVersion"`
	Title       string `json:"title"`
	Links       struct {
		Internal map[string]bool
		External map[string]bool
	}
	IsLoginForm bool `json:"isLoginForm"`
}

var response Response
var m = sync.Mutex{}
var wg = sync.WaitGroup{}

func ScrapeWebPage(url string) (Response, error) {
	var errors error
	isForm, containsEmailInput, containsPasswordInput := false, false, false

	c := colly.NewCollector()

	c.OnRequest(func(r *colly.Request) {
		logger.Info("Visiting page with URL: ", r.URL)
	})

	c.OnError(func(r *colly.Response, err error) {
		logger.Error("Failed with error: ", err.Error())
		errors = err
	})

	c.OnResponse(func(r *colly.Response) {
		docString, err := html.Parse(strings.NewReader(string(r.Body)))
		if err != nil {
			logger.Error("Error parsing HTML from string: ", err)
			response.HtmlVersion = "Unknown"
		} else {
			response.HtmlVersion = getHTMLVersion(docString.FirstChild)
		}
	})

	c.OnHTML("head title", func(h *colly.HTMLElement) {
		response.Title = h.Text
	})

	c.OnHTML("a[href]", func(h *colly.HTMLElement) {
		url = h.Attr("href")
		if IsValidURL(url) {
			if response.Links.External == nil {
				response.Links.External = map[string]bool{}
			}
			response.Links.External[url] = false
		} else {
			if response.Links.Internal == nil {
				response.Links.Internal = map[string]bool{}
			}
			response.Links.Internal[url] = false
		}
	})

	c.OnHTML("input[type=\"email\"],input[type=\"text\"],input[type=\"password\"],form", func(h *colly.HTMLElement) {
		switch CleanFields(h.Attr("type")) {
		case "email":
			containsEmailInput = true
		case "password":
			containsPasswordInput = true
		default:
			isForm = true
		}
	})

	c.OnScraped(func(r *colly.Response) {
		if isForm && containsEmailInput && containsPasswordInput {
			response.IsLoginForm = true
		}

		for url := range response.Links.External {
			wg.Add(1)
			go markAccessibleLinks(url)
		}
		wg.Wait()
	})

	c.Visit(url)

	return response, errors
}

func markAccessibleLinks(url string) {
	res, err := http.Head(url)
	if err == nil && res.StatusCode == 200 {
		m.Lock()
		response.Links.External[url] = true
		m.Unlock()
	}
	wg.Done()
}

// Extracts the HTML version based on the DOCTYPE.
func getHTMLVersion(doctypeNode *html.Node) string {
	if doctypeNode != nil && doctypeNode.Type == html.DoctypeNode {
		// HTML5 DOCTYPE is simply <!DOCTYPE html>
		if strings.ToLower(doctypeNode.Data) == "html" && len(doctypeNode.Attr) == 0 {
			return "HTML5"
		}
		// For older HTML versions, the DOCTYPE will contain more information.
		return fmt.Sprintf("HTML (DOCTYPE: %s)", doctypeNode.Data)
	}
	return "Unknown (No DOCTYPE found or invalid)"
}
