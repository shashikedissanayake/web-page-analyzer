package service

import (
	"github.com/shashikedissanayake/web-page-analyzer/server/model"
	"github.com/shashikedissanayake/web-page-analyzer/server/utils"
)

type IScraperService interface {
	ScrapeWebPage(string) (*model.ScraperResponse, error)
}

type ScraperService struct {
	scraper utils.IWebPageAnalyzer
}

func CreateNewScraperService(scraper utils.IWebPageAnalyzer) IScraperService {
	return &ScraperService{
		scraper,
	}
}

func (ss *ScraperService) ScrapeWebPage(url string) (*model.ScraperResponse, error) {
	res, err := ss.scraper.AnalyzeWebPage(url)
	if err != nil {
		return nil, err
	}

	internalLinkCount, internalInaccessibleLinkCount := 0, 0
	for _, val := range res.Links.Internal {
		internalLinkCount++
		if !val {
			internalInaccessibleLinkCount++
		}
	}

	externalLinkCount, externalInaccessibleLinkCount := 0, 0
	for _, val := range res.Links.External {
		externalLinkCount++
		if !val {
			externalInaccessibleLinkCount++
		}
	}

	return &model.ScraperResponse{
		Url:                           url,
		HtmlVersion:                   res.HtmlVersion,
		PageTitle:                     res.Title,
		InternalLinkCount:             internalLinkCount,
		InternalInaccessibleLinkCount: internalInaccessibleLinkCount,
		ExternalLinkCount:             externalLinkCount,
		ExternalInaccessibleLinkCount: externalInaccessibleLinkCount,
		IsLoginForm:                   res.IsLoginForm,
		HeaderTagCount: model.HeaderTagsCount{
			H1: res.HeaderTags["h1"],
			H2: res.HeaderTags["h2"],
			H3: res.HeaderTags["h3"],
			H4: res.HeaderTags["h4"],
			H5: res.HeaderTags["h5"],
			H6: res.HeaderTags["h6"],
		},
	}, nil
}
