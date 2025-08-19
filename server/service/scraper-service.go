package service

import (
	"github.com/shashikedissanayake/web-page-analyzer/server/model"
	"github.com/shashikedissanayake/web-page-analyzer/server/utils"
)

type IScraperService interface {
	ScrapeWebPage(string) (*model.ScraperResponse, error)
}

type ScraperService struct {
}

func CreateNewScraperService() IScraperService {
	return &ScraperService{}
}

func (ss *ScraperService) ScrapeWebPage(url string) (*model.ScraperResponse, error) {
	res, err := utils.ScrapeWebPage(url)
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
	for _, val := range res.Links.Internal {
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
	}, nil
}
