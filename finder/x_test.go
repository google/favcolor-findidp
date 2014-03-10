package findIDP

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
)

func init() {
	fmt.Println("Testing findIDP\n")
}

func TestJRD(t *testing.T) {
	var jrd []byte
	var i *IDP
	jrd = []byte(`{
       "subject" : "acct:carol@example.com",
       "links" :
       [
         {
           "rel" : "http://openid.net/specs/connect/1.0/issuer",
           "href" : "https://openid.example.com"
         }
       ]
     }`)
	var res interface{}
	json.Unmarshal(jrd, &res)
	i = idpFromJRD(res)
	if i.URI != "https://openid.example.com" {
		t.Errorf("Wanted URI https://openid.example.com, got %s", i.URI)
	}
	if i.Protocol != OIDC {
		t.Error("Protocol list is not [OIDC]")
	}

	jrd = []byte(`  {
       "subject" : "http://blog.example.com/article/id/314",
       "aliases" :
       [
         "http://blog.example.com/cool_new_thing",
         "http://blog.example.com/steve/article/7"
       ],
       "properties" :
       {
         "http://blgx.example.net/ns/version" : 1.3,
         "http://blgx.example.net/ns/ext" : null
       },
       "links" :
       [
         {
           "rel" : "copyright",
           "href" : "http://www.example.com/copyright"
         },
         {
           "rel" : "author",
           "href" : "http://blog.example.com/author/steve",
           "titles" :
           {
             "en-us" : "The Magical World of Steve",
             "fr" : "Le Monde Magique de Steve"
           },
           "properties" :
           {
             "http://example.com/role" : "editor"
           }
         }
       ]
     }`)
	var res2 interface{}
	json.Unmarshal(jrd, &res2)
	i = idpFromJRD(res2)
	if i != nil {
		t.Error("Found a bogus IDP in JRD without one")
	}

}

func TestIDP(t *testing.T) {
	i := wellKnownIDP(Google)
	if i.Protocol.name() != "oidc" {
		t.Error("i.Protocol name should be oidc, is: " +
			i.Protocol.name())
	}
	if i.URI != "https://accounts.google.com" {
		t.Error("URI should be http://accounts.google.com, is: " + i.URI)
	}
	ip := wellKnownIDP(Google)
	if ip.Key != Google {
		t.Errorf("Wellknown key should be Google, is %s", ip.Key)
	}
	if ip.Protocol != OIDC {
		t.Errorf("Protocol should be 0 is %d", ip.Protocol)
	}
}

type receivedIDP struct {
	IDP      string
	Protocol string
	Verified bool
}
type received map[string][]receivedIDP

func TestFind(t *testing.T) {

	recorder := httptest.NewRecorder()
	request, err := http.NewRequest("GET", "http://findidp.com/find?email=foo",
		nil)
	if err != nil {
		t.Fatal("NewRequest failed")
	}

	find(recorder, request)
	if recorder.Header()["Content-Type"][0] != "application/json" {
		t.Error("application/json not there")
	}
}

func TestRequest(t *testing.T) {

	recorder := httptest.NewRecorder()
	request, err := http.NewRequest("GET", "http://findidp.com/find",
		nil)
	if err != nil {
		t.Fatal("NewRequest failed")
	}
	find(recorder, request)
	if recorder.Code != 400 {
		t.Errorf("missing arg to find produced status %d, wanted 400",
			recorder.Code)
	}
}

func TestSearcher(t *testing.T) {

	db := getHandles(nil).db
	email := EMail{"em1", "example.com", []IDP{*wellKnownIDP(Google)}}
	db.storeEMail(&email)
	recorder := httptest.NewRecorder()
	request, err := http.NewRequest("GET", "http://findidp.com/find?m=em1@example.com",
		nil)
	if err != nil {
		t.Fatal("NewRequest failed")
	}
	find(recorder, request)
	var m received
	err = json.Unmarshal(recorder.Body.Bytes(), &m)
	if err != nil {
		t.Errorf("Error parsing JSON %s", err)
	}
	if len(m) != 1 || m["idps"] == nil || len(m["idps"]) != 1 {
		t.Error("Returned JSON doesn’t have 1 'idps' member & list-len 1")
	}
	idp := m["idps"][0]

	if idp.IDP != "https://accounts.google.com" {
		t.Error("Failed to find em1@example.com in database")
	}
	if idp.Protocol != OIDC.name() {
		t.Errorf("Didn’t get one protocol with value 'oidc'")
	}
	if !idp.Verified {
		t.Error("IDP not marked as verified")
	}

	recorder = httptest.NewRecorder()
	request, err = http.NewRequest("GET", "http://findidp.com/find?m=bengalrune@gmail.com",
		nil)
	if err != nil {
		t.Fatal("NewRequest failed")
	}
	find(recorder, request) // should work with webfist
	err = json.Unmarshal(recorder.Body.Bytes(), &m)
	if err != nil {
		t.Errorf("Error parsing JSON %s", err)
	}
	if len(m) != 1 || m["idps"] == nil || len(m["idps"]) != 1 {
		t.Error("Returned JSON doesn’t have 1 'idps' member & list-len 1")
	}
	idp = m["idps"][0]

	if idp.IDP != "https://accounts.google.com" {
		t.Error("Failed to find em1@example.com in database")
	}
	if idp.Protocol != OIDC.name() {
		t.Errorf("Didn’t get one protocol with value 'oidc'")
	}
	if !idp.Verified {
		t.Error("IDP not marked as verified")
	}

	recorder = httptest.NewRecorder()
	request, err = http.NewRequest("GET", "http://findidp.com/find?m=nobody@example.com",
		nil)
	if err != nil {
		t.Fatal("NewRequest failed")
	}
	find(recorder, request)

	// This is really sort of a system not unit test, since it actually goes
	// and uses the DNS module.  Heresy!
	recorder = httptest.NewRecorder()
	request, err = http.NewRequest("GET", "http://findidp.com/find?m=tbray@textuality.com",
		nil)
	if err != nil {
		t.Fatal("NewRequest failed")
	}
	find(recorder, request)

	err = json.Unmarshal(recorder.Body.Bytes(), &m)
	if err != nil {
		t.Errorf("Error parsing JSON %s", err)
	}

	if m["idps"] == nil || len(m["idps"]) != 1 {
		t.Error("Returned JSON doesn’t have 1 'idps' member & list-len 1")
	}
	idp = m["idps"][0]

	if idp.IDP != "https://accounts.google.com" {
		t.Error("Failed to find tbray@textuality.com via MX")
	}
	if idp.Protocol != OIDC.name() {
		t.Errorf("Didn’t get one protocol with value 'oidc'")
	}
	if idp.Verified {
		t.Error("IDP incorrectly marked as verified")
	}
	//======================
	recorder = httptest.NewRecorder()
	request, err = http.NewRequest("GET", "http://findidp.com/find?m=nobody@yahoo.com.ar",
		nil)
	if err != nil {
		t.Fatal("NewRequest failed")
	}
	find(recorder, request)

	err = json.Unmarshal(recorder.Body.Bytes(), &m)
	if err != nil {
		t.Errorf("Error parsing JSON %s", err)
	}

	if m["idps"] == nil || len(m["idps"]) != 1 {
		t.Error("Returned JSON doesn’t have 1 'idps' member & list-len 1")
	}
	idp = m["idps"][0]

	if idp.IDP != "https://login.yahoo.com" {
		t.Error("Failed to find nobody@.com via domainmatch")
	}
	if idp.Protocol != OpenID2.name() {
		t.Errorf("Didn’t get one protocol with value 'openid2'")
	}
	if idp.Verified {
		t.Error("IDP incorrectly marked as verified")
	}

}

func TestEMail(t *testing.T) {

	e, err := parseEMail("x")
	if err.Error() != "malformed email" {
		t.Error("email address 'x' should be malformed, isn't")
	}
	e, err = parseEMail("a@b.com")
	if err != nil {
		t.Error("err non-null on good email address")
	}
	name, domain := e.parts()
	if name != "a" || domain != "b.com" {
		t.Error("Wanted a/b.com, got " + name + "/" + domain)
	}
}

func TestMerge(t *testing.T) {
	i1 := []IDP{*wellKnownIDP(Google), *wellKnownIDP(Microsoft)}
	i2 := []IDP{*wellKnownIDP(Microsoft), *wellKnownIDP(Facebook)}
	i3 := merge(i1, i2)
	if len(i3) != 3 {
		t.Errorf("Wanted merge size 3, got %d", len(i3))
	}
}

func TestCrumbs(t *testing.T) {
	a0 := SearchResult{WebfingerType, []IDP{}}
	a1 := SearchResult{MXType, []IDP{*wellKnownIDP(Google)}}
	c := Crumbs{[]SearchResult{a0, a1}}
	if c.Less(0, 1) {
		t.Errorf("Empty result %s ranks higher than non-empty %s", crumbTrail(&c)[0], crumbTrail(&c)[1])
	}

	a2 := SearchResult{WebfingerType, []IDP{*wellKnownIDP(Google)}}
	a3 := SearchResult{TestingType, []IDP{*wellKnownIDP(Google)}}
	c = Crumbs{[]SearchResult{a2, a3}}
	if c.Less(1, 0) {
		t.Errorf("Unverified result %s ranks higher than verified %s", crumbTrail(&c)[1], crumbTrail(&c)[0])
	}
}

/*
func TestDNS(t *testing.T) {
	message := makeDNSQuery("foo.com")
	question := message.Question[0]
	if question.Qtype != dns.TypeMX {
		t.Error("qtype %d", question.Qtype)
	}
	if question.Name != "foo.com." {
		t.Error("name is: " + question.Name + ", wanted 'foo.com.'")
	}

	r1 := new(dns.MX)
	r1.Hdr = dns.RR_Header{Name: "r1.nl.", Rrtype: dns.TypeMX, Class: dns.ClassINET, Ttl: 3600}
	r1.Preference = 1
	r1.Mx = "r1"
	r2 := new(dns.MX)
	r2.Hdr = dns.RR_Header{Name: "r2.nl.", Rrtype: dns.TypeMX, Class: dns.ClassINET, Ttl: 3600}
	r2.Preference = 2
	r2.Mx = "r2"
	r3 := new(dns.MX)
	r3.Hdr = dns.RR_Header{Name: "r3.nl.", Rrtype: dns.TypeMX, Class: dns.ClassINET, Ttl: 3600}
	r3.Preference = 3
	r3.Mx = "r3"
	a := []dns.RR{r1, r2, r3}
	s := preferredDomain(a)
	if s != "r1" {
		t.Error("preferredDomain returned " + s + ", wanted r1")
	}
	a = []dns.RR{r3, r2, r1}
	s = preferredDomain(a)
	if s != "r1" {
		t.Error("preferredDomain returned " + s + ", wanted r1")
	}
	a = []dns.RR{r2}
	s = preferredDomain(a)
	if s != "r2" {
		t.Error("preferredDomain returned " + s + ", wanted r2")
	}
}
*/

/*
	fmt.Printf("Target: %s\n", email.string())
	for _, result := range crumbTrail(crumbs) {
		fmt.Printf(" %s (%v):", resultLabel(result.rtype), ResultStrengths[result.rtype].verified)
		for _, idp := range result.idps {
			fmt.Printf(" [%s]", idp.URI)
		}
		fmt.Println("")
	}*/
