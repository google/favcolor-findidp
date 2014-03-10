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
	var s WebfingerSearcher
	RegisterSearcher(s)
}

type WebfingerSearcher struct{}

func (_ WebfingerSearcher) Search(email *EMail, c chan SearchResult, handles Handles) {
	_, domain := email.parts()
	query := "resource=acct:" + url.QueryEscape(email.string()) +
		"&rel=" + url.QueryEscape("http://openid.net/specs/connect/1.0/issuer")
	uri := "https://" + domain + "/.well-known/webfinger?" + query
	c <- procWebfingerURI(uri, handles, WebfingerType)
}

func (_ WebfingerSearcher) Label() string {
	return "WebFinger searcher"
}

func procWebfingerURI(uri string, handles Handles, rt ResultType) SearchResult {
	idps := []IDP{}
	f, err := fetchJSON(uri, handles, "Webfinger")
	if err == nil {
		idp := idpFromJRD(f)
		if idp == nil {
			msg := fmt.Sprintf("No IDP in webfinger resp for %s", uri)
			handles.logger.logError(msg)
		} else {
			idps = []IDP{*idp}
		}
	}
	return SearchResult{rt, idps}
}
