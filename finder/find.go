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
