package service

import (
	"fmt"
	"io"
	"net/http"
	"regexp"
	"strings"
	"sync"

	"github.com/shashikedissanayake/web-page-analyzer/server/model"
	"github.com/shashikedissanayake/web-page-analyzer/server/utils"
	logger "github.com/sirupsen/logrus"
	"golang.org/x/net/html"
)

//go:generate mockgen -source=webpage_analyzer_service.go -destination=webpage_analyzer_service_mock.go -package=service
type IWebPageAnalyzerService interface {
	AnalyzeWebPage(string) (*model.AnalyzerResults, error)
}

type WebPageAnalyzerService struct{}

func CreateWebPageAnalyzerService() IWebPageAnalyzerService {
	return &WebPageAnalyzerService{}
}

func (wpa *WebPageAnalyzerService) AnalyzeWebPage(url string) (*model.AnalyzerResults, error) {
	tokenizer, err := wpa.fetchWebPage(url)
	if err != nil {
		return nil, err
	}

	response := wpa.iterateTokenizer(tokenizer, url)

	return &response, nil
}

func (wpa *WebPageAnalyzerService) fetchWebPage(url string) (*html.Tokenizer, error) {
	resp, err := http.Get(url)
	if err != nil {
		logger.Error("Failed to fetch webpage for given URL:", url, "with error:", err.Error())
		return nil, fmt.Errorf("failed to get url with error: %s", err.Error())
	} else if resp.StatusCode != http.StatusOK {
		logger.Error("Failed to fetch webpage for given URL:", url, " with status:", resp.StatusCode)
		return nil, fmt.Errorf("failed to fetch page with status code: %d", resp.StatusCode)
	}

	defer resp.Body.Close()
	webpage, err := io.ReadAll(resp.Body)
	if err != nil {
		logger.Error("Failed to read webpage for given URL:", url, " with error:", err.Error())
		return nil, fmt.Errorf("failed to read response with error:%s", err.Error())
	}

	tokenizer := html.NewTokenizer(strings.NewReader(string(webpage)))

	logger.Info("Successfully fetched URL:", url)

	return tokenizer, nil
}

func (wpa *WebPageAnalyzerService) iterateTokenizer(
	tokenizer *html.Tokenizer, url string,
) model.AnalyzerResults {
	response := model.AnalyzerResults{
		HeaderTags: map[string]int{},
		Links:      map[string]*model.Links{},
	}

	if tokenizer == nil {
		logger.Error("Provided empty tokenizer")
		return response
	}

	isForm, containsEmailInput, containsPasswordInput := false, false, false
	wg := sync.WaitGroup{}
	linkDetailsChannel := make(chan model.LinkDetails)

	logger.Info("Start iteration of tokenized webpage")

	go wpa.resultsReader(response.Links, linkDetailsChannel)
loop:
	for {
		tokenType := tokenizer.Next()

		switch tokenType {
		case html.ErrorToken:
			break loop
		case html.DoctypeToken:
			token := tokenizer.Token()
			response.HtmlVersion = wpa.getHTMLVersion(&token)
		case html.StartTagToken, html.SelfClosingTagToken:
			token := tokenizer.Token()
			switch token.Data {
			case "title":
				tokenizer.Next()
				token := tokenizer.Token()
				response.Title = token.Data
			case "h1", "h2", "h3", "h4", "h5", "h6":
				count, exists := response.HeaderTags[token.Data]
				if exists {
					response.HeaderTags[token.Data] = count + 1
				} else {
					response.HeaderTags[token.Data] = 1
				}

			case "a":
				for _, attrbute := range token.Attr {
					if attrbute.Key == "href" {
						record, isFound := response.Links[attrbute.Val]
						if !isFound {
							response.Links[attrbute.Val] = &model.Links{
								Count: 1,
							}

							wg.Add(1)
							go wpa.markAccessiblityOfLink(
								url,
								attrbute.Val,
								&wg,
								linkDetailsChannel,
							)
						} else {
							record.Count++
						}
					}
				}
			case "form":
				isForm = true
			case "input":
				for _, attrbute := range token.Attr {
					if attrbute.Key == "type" {
						switch attrbute.Val {
						case "email", "text":
							containsEmailInput = true
						case "password":
							containsPasswordInput = true
						}
					}
				}
			}
		}
	}
	response.IsLoginForm = isForm && containsEmailInput && containsPasswordInput
	logger.Info("Finised iteration of tokenized webpage")
	wg.Wait()
	close(linkDetailsChannel)
	logger.Info("Finised all worker thread")
	return response
}

func (wpa *WebPageAnalyzerService) getHTMLVersion(doctypeToken *html.Token) string {
	olderHtmlversionIdentifierRegEx := regexp.MustCompile(` ".+" `)
	if doctypeToken != nil && doctypeToken.Type == html.DoctypeToken {
		versionDetails := utils.CleanFields(olderHtmlversionIdentifierRegEx.FindString(doctypeToken.Data))
		if strings.ToLower(doctypeToken.Data) == "html" {
			return "HTML 5"
		} else if versionDetails != "" {
			versionDetails := strings.ReplaceAll(versionDetails, "\"", "")
			return strings.Join(strings.Split(versionDetails, " ")[1:], " ")
		}
	}
	return "Unknown"
}

func (wpa *WebPageAnalyzerService) markAccessiblityOfLink(
	currentUrl string,
	key string,
	wg *sync.WaitGroup,
	linkDetailsChan chan<- model.LinkDetails,
) {
	logger.Info("Started fetching head object of key:", key)
	defer wg.Done()

	url := key
	isValidUrl := utils.IsValidURL(url)
	if !isValidUrl {
		url = utils.GenerateInternalUrl(currentUrl, key)
	}

	resp, err := http.Head(url)
	isAccessibleLink := err == nil && resp.StatusCode == http.StatusOK

	linkDetails := model.LinkDetails{
		URL:          key,
		IsAccessible: isAccessibleLink,
	}

	if isValidUrl {
		linkDetails.LinkType = model.EXTERNAL
	} else {
		linkDetails.LinkType = model.INTERNAL
	}
	linkDetailsChan <- linkDetails

	logger.Info("Finished fetching head object of key:", key, " with results:", isAccessibleLink)
}

func (wpa *WebPageAnalyzerService) resultsReader(links map[string]*model.Links, linkDetailsChan <-chan model.LinkDetails) {
	for linkDetail := range linkDetailsChan {
		logger.Info("Successfully consumed message with url:", linkDetail.URL)

		record, ok := links[linkDetail.URL]
		if !ok {
			logger.Error("Record not found for key:", linkDetail.URL)
			return
		}

		switch linkDetail.LinkType {
		case model.INTERNAL:
			record.LinkType = model.INTERNAL
			record.IsAccessible = linkDetail.IsAccessible
		case model.EXTERNAL:
			record.LinkType = model.EXTERNAL
			record.IsAccessible = linkDetail.IsAccessible
		}
	}
	logger.Info("Successfully exited from channel")
}
