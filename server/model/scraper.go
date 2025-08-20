package model

type HeaderTagsCount struct {
	H1 int `json:"h1,omitempty"`
	H2 int `json:"h2,omitempty"`
	H3 int `json:"h3,omitempty"`
	H4 int `json:"h4,omitempty"`
	H5 int `json:"h5,omitempty"`
	H6 int `json:"h6,omitempty"`
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
