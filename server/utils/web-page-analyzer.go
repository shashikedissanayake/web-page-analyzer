package utils

import (
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"
	"sync"

	logger "github.com/sirupsen/logrus"
	"golang.org/x/net/html"
)

type Links struct {
	Internal map[string]bool
	External map[string]bool
}

type Response struct {
	HtmlVersion string
	Title       string
	Links       Links
	IsLoginForm bool
	HeaderTags  map[string]int
}

type IWebPageAnalyzer interface {
	AnalyzeWebPage(string) (*Response, error)
}

type WebPageAnalyzer struct{}

func CreateWebPageAnalyzer() IWebPageAnalyzer {
	return &WebPageAnalyzer{}
}

func (wpa *WebPageAnalyzer) AnalyzeWebPage(url string) (*Response, error) {
	tokenizer, err := wpa.fetchWebPage(url)
	fmt.Println("err:", err)
	if err != nil {
		return nil, err
	}

	response := wpa.iterateTokenizer(tokenizer, url)

	return &response, nil
}

func (wpa *WebPageAnalyzer) fetchWebPage(url string) (*html.Tokenizer, error) {
	resp, err := http.Get(url)
	if err != nil {
		logger.Error("Failed to fetch webpage for given URL:", url)
		return nil, err
	} else if resp.StatusCode != http.StatusOK {
		logger.Error("Failed to fetch webpage for given URL:", url, " with status:", resp.StatusCode)
		return nil, errors.New("FAILED_TO_FETCH")
	}

	defer resp.Body.Close()
	webpage, err := io.ReadAll(resp.Body)
	if err != nil {
		logger.Error("Failed to read webpage for given URL:", url)
		return nil, err
	}

	tokenizer := html.NewTokenizer(strings.NewReader(string(webpage)))

	logger.Info("Successfully fetched URL:", url)

	return tokenizer, nil
}

func (wpa *WebPageAnalyzer) iterateTokenizer(
	tokenizer *html.Tokenizer, url string,
) Response {
	response := Response{
		HeaderTags: map[string]int{},
		Links: Links{
			Internal: map[string]bool{},
			External: map[string]bool{},
		},
	}

	if tokenizer == nil {
		logger.Error("Provided empty tokenizer")
		return response
	}

	isForm, containsEmailInput, containsPasswordInput := false, false, false
	wg := sync.WaitGroup{}
	mux := sync.Mutex{}

	logger.Info("Start iteration of tokenized webpage")
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
						wg.Add(1)
						go wpa.markAccessiblityOfLink(
							url,
							attrbute.Val,
							&response.Links,
							&wg,
							&mux,
						)
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
	logger.Info("Finised all worker thread")
	return response
}

func (wpa *WebPageAnalyzer) getHTMLVersion(doctypeToken *html.Token) string {
	if doctypeToken != nil && doctypeToken.Type == html.DoctypeToken {
		if strings.ToLower(doctypeToken.Data) == "html" && len(doctypeToken.Attr) == 0 {
			return "HTML5"
		}

		return fmt.Sprintf("HTML (DOCTYPE: %s)", doctypeToken.Data)
	}
	return "Unknown"
}

func (wpa *WebPageAnalyzer) markAccessiblityOfLink(
	currentUrl string,
	key string,
	links *Links,
	wg *sync.WaitGroup,
	mux *sync.Mutex,
) {
	logger.Info("Started fetching head object of key:", key)
	defer wg.Done()

	url := key
	isValidUrl := IsValidURL(url)
	if !isValidUrl {
		url = GenerateInternalUrl(currentUrl, key)
	}

	resp, err := http.Head(url)
	isAccessibleLink := err == nil && resp.StatusCode == http.StatusOK
	mux.Lock()
	if isValidUrl {
		links.External[url] = isAccessibleLink
	} else {
		links.Internal[url] = isAccessibleLink
	}
	mux.Unlock()
	logger.Info("Finished fetching head object of key:", key, " with results:", isAccessibleLink)
}
