/*
Copyright [2014] Google, Inc.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
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
