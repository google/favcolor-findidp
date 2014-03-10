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
	"regexp"
)

var reMap = make(map[*regexp.Regexp][]IDP)

func init() {

	googleRE, _ := regexp.Compile("(google|gmail|googlemail)\\.com$")
	microsoftRE, _ := regexp.Compile("(live|outlook|hotmail)\\.com$")
	yahooRE, _ := regexp.Compile("yahoo(\\.[a-z]+)?\\.[a-z]+$")

	reMap[googleRE] = []IDP{*wellKnownIDP(Google)}
	reMap[microsoftRE] = []IDP{*wellKnownIDP(Microsoft)}
	reMap[yahooRE] = []IDP{*wellKnownIDP(Yahoo)}

	var s DomainMatcher
	RegisterSearcher(s)
}

type DomainMatcher struct{}

func (_ DomainMatcher) Search(email *EMail, c chan SearchResult, _ Handles) {
	_, domain := email.parts()
	result := SearchResult{DomainMatchType, []IDP{}}

	for re, idps := range reMap {
		if re.MatchString(domain) {
			result.idps = idps
		}
	}
	c <- result
}

func (_ DomainMatcher) Label() string {
	return "Domain name matcher"
}
