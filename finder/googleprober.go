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
