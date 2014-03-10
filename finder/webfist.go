package findIDP

import (
	"fmt"
	"net/url"
)

func init() {
	var s WebFistSearcher
	RegisterSearcher(s)
}

type WebFistSearcher struct{}

func (_ WebFistSearcher) Search(email *EMail, c chan SearchResult, handles Handles) {
	query := "resource=acct:" + url.QueryEscape(email.string())
	uri := "http://webfist.org/.well-known/webfinger?" + query
	c <- procWebFist(uri, handles)
}

func (_ WebFistSearcher) Label() string {
	return "WebFist searcher"
}

func procWebFist(uri string, handles Handles) SearchResult {
	failure := SearchResult{WebFistType, []IDP{}}
	f, err := fetchJSON(uri, handles, "WebFist")
	if err == nil {
		webfingerURI := linkFromJRD(f, "http://webfist.org/spec/rel")
		if webfingerURI == "" {
			msg := fmt.Sprintf("No Webfinger link found in WebFist for %s", uri)
			handles.logger.logError(msg)
		} else {
			return procWebfingerURI(webfingerURI, handles, WebFistType)
		}
	}
	return failure
}
