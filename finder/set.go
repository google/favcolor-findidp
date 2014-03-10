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
	"net/http"
)

func set(w http.ResponseWriter, r *http.Request) {

	if r.Method != "POST" {
		htmlError(w, 400, "Must POST to set IDP")
		return
	}
	err := r.ParseForm()
	if err != nil {
		htmlError(w, 400, err.Error())
		return
	}
	fields := r.Form
	if fields["e"] == nil {
		htmlError(w, 400, "Missing e= argument")
		return
	}
	email, err := parseEMail(fields["e"][0])
	if err != nil {
		htmlError(w, 400, "Invalid e= argument")
		return
	}
	var idp IDP
	f := fields["i"][0]
	if f != "-" {
		switch f {
		case "F":
			idp = *wellKnownIDP(Facebook)
		case "G":
			idp = *wellKnownIDP(Google)
		case "M":
			idp = *wellKnownIDP(Microsoft)
		case "Y":
			idp = *wellKnownIDP(Yahoo)
		default:
			htmlError(w, 400, "Bad value for i= argument")
			return
		}
	} else {
		u := fields["u"][0]
		if u == "-" {
			htmlError(w, 400, "Must provide either i= or u=")
			return
		}
		var protocol Protocol
		f = fields["p"][0]
		if f == "?" {
			protocol = Unknown
		} else {
			switch f {
			case "oidc":
				protocol = OIDC
			case "oid2":
				protocol = OpenID2
			case "persona":
				protocol = Persona
			case "saml":
				protocol = SAML
			default:
				htmlError(w, 400, "Bad value for p= argument")
				return
			}
		}
		idp = IDP{IDPKey(u), u, protocol}
	}

	email.IDPs = []IDP{idp}
	err = getDatabase(r).storeEMail(email)
	if err != nil {
		htmlError(w, 500, "Failed to store email/IDP mapping: "+err.Error())
		return
	}
}
