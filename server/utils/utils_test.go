package utils

import (
	"testing"

	"github.com/stretchr/testify/suite"
)

type UtilsSuite struct {
	suite.Suite
}

func TestUtilsSuite(t *testing.T) {
	suite.Run(t, &UtilsSuite{})
}

func (us *UtilsSuite) TestCleanFields() {
	var testCases = []struct {
		name   string
		input  string
		output string
	}{
		{"Passing empty string", "", ""},
		{"Passing a valid string", "test string", "test string"},
		{"Passing a valid string with spaces in front", "   test string", "test string"},
		{"Passing a valid string with spaces in back", "test string  ", "test string"},
		{"Passing a valid string with spaces in front and back", "  test string  ", "test string"},
	}

	for _, test := range testCases {
		us.T().Run(test.name, func(t *testing.T) {
			out := CleanFields(test.input)

			us.Equal(out, test.output)
		})
	}
}

func (us *UtilsSuite) TestIsValidURL() {
	var testCases = []struct {
		name   string
		input  string
		output bool
	}{
		{"Passing valid URL", "https://google.com", true},
		{"Passing valid URL with http", "http://google.com", true},
		{"Passing valid URL with path", "http://google.com/test", true},
		{"Passing valid URL with page section", "http://google.com/test#test", true},
		{"Passing string", "test string", false},
		{"Passing url path", "/test", false},
		{"Passing url with page section", "#test", false},
		{"Passing url without protocol", "google.com/test", false},
	}

	for _, test := range testCases {
		us.T().Run(test.name, func(t *testing.T) {
			out := IsValidURL(test.input)

			us.Equal(out, test.output)
		})
	}
}

func (us *UtilsSuite) TestGenerateInternalUrl() {
	var testCases = []struct {
		name   string
		url    string
		href   string
		output string
	}{
		{"Passing valid URL with path in href", "https://google.com/test1", "/test", "https://google.com/test"},
		{"Passing valid URL with section in href", "https://google.com/test1", "#test", "https://google.com/test1#test"},
	}

	for _, test := range testCases {
		us.T().Run(test.name, func(t *testing.T) {
			out := GenerateInternalUrl(test.url, test.href)

			us.Equal(test.output, out)
		})
	}
}
