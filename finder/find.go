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
	"net/http"
)

func find(w http.ResponseWriter, r *http.Request) {

	email, handles, err := checkRequest(r)

	w.Header().Add("Content-Type", "application/json")
	if err != nil {
		w.WriteHeader(400)
		fmt.Fprintf(w, "{\"error\":\"%s\"}", err.Error())
		return
	}
	crumbs := new(Crumbs)
	bestResult, verified := scanForIDPs(email, handles, crumbs)
	w.Write(listIDPsAsJSON(bestResult.idps, verified))
}

func listIDPsAsJSON(idps []IDP, verified bool) []byte {
	top := make(map[string][]reported)
	list := make([]reported, len(idps))
	for i, idp := range idps {
		list[i] = reportOne(idp, verified)
	}
	top["idps"] = list
	json, _ := json.Marshal(top) // we made this, no errors
	return json
}
