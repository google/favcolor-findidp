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
