package findIDP

import ()

func init() {
	var s GoogleProber
	RegisterSearcher(s)
}

const KEY string = "AIzaSyCrsPhpPagFLHkbKCtGWZMo937gYahnjsM"

type GoogleProber struct{}

func (_ GoogleProber) Search(email *EMail, c chan SearchResult, handles Handles) {
	_, domain := email.parts()
	result := SearchResult{GoogleProbeType, []IDP{}}
	uri := "https://www.googleapis.com/rpc?apiVersion=v1&" +
		"method=identitytoolkit.relyingparty.createAuthUrl&" +
		"identifier=" + domain + "&" +
		"continueUrl=http://localhost&" +
		"key=" + KEY
	body, err := fetchJSON(uri, handles, "Google hosting prober")
	if err == nil {
		m := body.(map[string]interface{})
		if m["result"] != nil {
			result.idps = []IDP{*wellKnownIDP(Google)}
		}
	}
	c <- result
}

func (_ GoogleProber) Label() string {
	return "Host prober"
}
