package service

import (
	"github.com/shashikedissanayake/web-page-analyzer/server/model"
	"github.com/shashikedissanayake/web-page-analyzer/server/utils"
)

//go:generate mockgen -source=scraper-service.go -destination=scraper-service_mock.go -package=service
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

	totalLinks, internalLinks, internalInaccessibleLinks, externalLinks, externalInaccessibleLinks := 0, 0, 0, 0, 0
	for _, link := range res.Links {
		totalLinks += link.Count
		switch link.LinkType {
		case utils.INTERNAL:
			internalLinks++
			if !link.IsAccessible {
				internalInaccessibleLinks++
			}
		case utils.EXTERNAL:
			externalLinks++
			if !link.IsAccessible {
				externalInaccessibleLinks++
			}
		}
	}

	return &model.ScraperResponse{
		Url:                           url,
		HtmlVersion:                   res.HtmlVersion,
		PageTitle:                     res.Title,
		TotalLinkCount:                totalLinks,
		InternalLinkCount:             internalLinks,
		InternalInaccessibleLinkCount: internalInaccessibleLinks,
		ExternalLinkCount:             externalLinks,
		ExternalInaccessibleLinkCount: externalInaccessibleLinks,
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
