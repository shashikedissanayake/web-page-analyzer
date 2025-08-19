package model

type HeaderTagsCount struct {
	H1TagCount int `json:"h1TagCount,omitempty"`
	H2TagCount int `json:"h2TagCount,omitempty"`
	H3TagCount int `json:"h3TagCount,omitempty"`
	H4TagCount int `json:"h4TagCount,omitempty"`
	H5TagCount int `json:"h5TagCount,omitempty"`
	H6TagCount int `json:"h6TagCount,omitempty"`
}

type ScraperResponse struct {
	Url                           string          `json:"url"`
	HtmlVersion                   string          `json:"htmlVersion"`
	PageTitle                     string          `json:"pageTitle"`
	HeaderTagCount                HeaderTagsCount `json:"headerTagCount"`
	InternalLinkCount             int             `json:"internalLinkCount"`
	InternalInaccessibleLinkCount int             `json:"internalInaccessibleLinkCount"`
	ExternalLinkCount             int             `json:"externalLinkCount"`
	ExternalInaccessibleLinkCount int             `json:"externalInaccessibleLinkCount"`
	IsLoginForm                   bool            `json:"isLoginForm"`
}

type ScraperRequest struct {
	Url string `json:"url"`
}
