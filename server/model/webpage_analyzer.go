package model

type LinkType string

const (
	INTERNAL LinkType = "INTERNAL"
	EXTERNAL LinkType = "EXTERNAL"
)

type LinkDetails struct {
	LinkType     LinkType
	URL          string
	IsAccessible bool
}

type Links struct {
	LinkType     LinkType
	Count        int
	IsAccessible bool
}

type AnalyzerResults struct {
	HtmlVersion string
	Title       string
	Links       map[string]*Links
	IsLoginForm bool
	HeaderTags  map[string]int
}
