package web

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestExtractUID(t *testing.T) {

	// test : expected
	paths := map[string]string{
		"/":                                     "",
		"/1.5/abc":                              "",
		"/somepath/123":                         "",
		"/somepath/123/storage/collectionname/": "",
		"123": "",
		"":    "",

		"/1.5/123":                         "123",
		"/1.5/123/info":                    "123",
		"/1.5/123/info/collections":        "123",
		"/1.5/123/storage/collectionname":  "123",
		"/1.5/123/storage/collectionname/": "123",
	}

	for path, expected := range paths {
		assert.Equal(t, expected, extractUID(path))
	}
}

func TestAcceptHeaderOk(t *testing.T) {

	// test headers are acceptable
	acceptable := []string{
		"application/json",
		"application/newlines",
	}

	for _, contentType := range acceptable {
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/", nil)
		req.Header.Set("Accept", contentType)
		assert.True(t, AcceptHeaderOk(w, req), contentType)
	}

	// test default rewrites
	defaults := []string{
		"",
		"*/*",
		"application/*",
		"*/json",

		// https://github.com/mostlygeek/go-syncstorage/issues/85
		"text/html,application/xhtml+xml,application/xml;q=0.9,*/*;q=0.8",
	}

	for _, contentType := range defaults {
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/", nil)
		req.Header.Set("Accept", contentType)
		assert.True(t, AcceptHeaderOk(w, req), contentType)
		assert.Equal(t, req.Header.Get("Accept"), "application/json")
	}

	// test 406 StatusNotAcceptable
	notAcceptable := []string{
		"x/yy",
		"text/html",
		"application/xhtml+xml",
		"application/xml",
		"text/html,application/xhtml+xml,application/xml;q=0.9",
	}

	for _, contentType := range notAcceptable {
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/", nil)
		req.Header.Set("Accept", contentType)
		assert.False(t, AcceptHeaderOk(w, req), contentType)
		assert.Equal(t, http.StatusNotAcceptable, w.Code)
	}

}