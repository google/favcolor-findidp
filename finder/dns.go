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
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

func init() {
	var s DnsSearcher
	RegisterSearcher(s)
}

type DnsSearcher struct{}

func (_ DnsSearcher) Label() string {
	return "MX record searcher"
}

func (_ DnsSearcher) Search(email *EMail, c chan SearchResult, handles Handles) {
	_, domain := email.parts()
	nothing := SearchResult{MXType, []IDP{}}
	resp, err := handles.client.Get("http://www.tbray.org:8001/d?d=" + domain)
	if err != nil {
		handles.logger.logError(fmt.Sprintf("MX lookup HTTP error: %s", err))
		c <- nothing
		return
	}
	if resp.StatusCode != 200 {
		msg := field(resp, "err")
		handles.logger.logError(fmt.Sprintf("MX lookup error %d: %s", resp.StatusCode, msg))
		c <- nothing
		return
	}
	msg := field(resp, "idp")
	if msg == "" {
		handles.logger.logError(fmt.Sprintf("No idp in MX search: %s", domain))
		c <- nothing
		return
	}
	key := IDPKey(msg)
	idp := wellKnownIDP(key)
	if idp == nil {
		idp = &IDP{key, msg, Unknown}
	}
	c <- SearchResult{MXType, []IDP{*idp}}
}

func field(resp *http.Response, key string) string {
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return ""
	}
	var f interface{}
	err = json.Unmarshal(body, &f)
	if err != nil {
		return ""
	}
	m := f.(map[string]interface{})
	return m[key].(string)
}

/*
import (
	"fmt"
	"net"
	"os"
	"strings"
)

func init() {
	var s DnsSearcher
	RegisterSearcher(s)
}

var gnames = [...]string{"google.com", "googlemail.com"}

type DnsSearcher struct{}

func (_ DnsSearcher) Label() string {
	return "MX record searcher"
}
func (_ DnsSearcher) Search(email *EMail, c chan SearchResult, handles Handles) {
	_, domain := email.parts()
	mxes, err := net.LookupMX(domain)
	if err != nil {
		handles.logger.logError(fmt.Sprintf("lookupMX error; %s", err))
		return
	}
	result := SearchResult{MXType, []IDP{}}
	for _, mx := range mxes {
		fmt.Fprintf(os.Stderr, " HOST %s\n", mx.Host)
		for _, gname := range gnames {
			if strings.Index(mx.Host, gname) != -1 {
				result.idps = []IDP{*wellKnownIDP(Google)}
				break
			}
		}
	}
	c <- result
}

import (
	"fmt"
	"github.com/miekg/dns"
	"net"
	"strings"
)

func init() {
	var s DnsSearcher
	RegisterSearcher(s)
}

type DnsSearcher struct{}

func (_ DnsSearcher) Search(email *EMail, c chan SearchResult, handles Handles) {

	_, domain := email.parts()
	request := makeDNSQuery(domain)

	conn, err := net.Dial("udp", "8.8.8.8:53")
	if err != nil {
		handles.logger.logError(fmt.Sprintf("Connection error: %s", err.Error()))
		return
	}

	response, _, err := new(dns.Client).ExchangeConn(request, conn)
	// response, _, err := new(dns.Client).Exchange(request, "8.8.8.8:53")
	if response == nil || response.Rcode != dns.RcodeSuccess {
		handles.logger.logError(fmt.Sprintf("DNS error: %s", err.Error()))
		return
	}

	mx := preferredDomain(response.Answer)
	result := SearchResult{MXType, []IDP{}}
	if mx != "" {
		// TODO: be a bit more sophisticated?
		if strings.Index(mx, "google.") != -1 {
			idp := wellKnownIDP(Google)
			result.idps = []IDP{*idp}
		}
	}
	c <- result
}
func (_ DnsSearcher) Label() string {
	return "MX record searcher"
}

func preferredDomain(rrs []dns.RR) (preferred string) {
	if len(rrs) == 0 {
		return
	}
	candidate := rrs[0].(*dns.MX)
	lowPref := candidate.Preference

	for _, answer := range rrs[1:] {
		thisMX := answer.(*dns.MX)
		if thisMX.Preference < lowPref {
			lowPref = thisMX.Preference
			candidate = thisMX
		}
	}

	preferred = candidate.Mx
	return
}

func makeDNSQuery(domain string) (message *dns.Msg) {
	message = new(dns.Msg)
	message.SetQuestion(domain+".", dns.TypeMX)
	message.RecursionDesired = true
	return
}
*/
