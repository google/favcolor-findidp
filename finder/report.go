package findIDP

import (
	"fmt"
	"html/template"
	"net/http"
	"sort"
)

type PageCrumb struct {
	Method   string
	IDPs     string
	Strength string
	Verified string
}
type ReportPage struct {
	EMail    string
	IDP      string
	Verified string
	Crumbs   []PageCrumb
}

func report(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "text/html")

	email, handles, err := checkRequest(r)
	if err != nil {
		htmlError(w, 400, err.Error())
		return
	}

	crumbs := new(Crumbs)
	bestResult, verified := scanForIDPs(email, handles, crumbs)
	t, err := template.ParseFiles("templates/Report.html")
	if err != nil {
		htmlError(w, 500, err.Error())
		return
	}
	var yn string
	if verified {
		yn = "(verified)"
	} else {
		yn = "(not verified)"
	}

	sort.Sort(crumbs)
	pageCrumbs := make([]PageCrumb, len(crumbTrail(crumbs)))
	for i, result := range crumbTrail(crumbs) {
		pageCrumbs[i].Method = resultLabel(result.rtype)
		if len(result.idps) == 0 {
			pageCrumbs[i].IDPs = "-none found-"
			pageCrumbs[i].Verified = ""
			pageCrumbs[i].Strength = ""
		} else {
			pageCrumbs[i].IDPs = result.idps[0].URI
			strength := ResultStrengths[result.rtype]
			pageCrumbs[i].Strength = fmt.Sprintf("%d", strength.strength)
			if strength.verified {
				pageCrumbs[i].Verified = "Y"
			} else {
				pageCrumbs[i].Verified = "N"
			}
		}
	}
	var idpURI string
	if len(bestResult.idps) == 0 {
		idpURI = "<none found>"
		yn = ""
	} else {
		idpURI = bestResult.idps[0].URI
	}

	page := ReportPage{email.string(), idpURI, yn, pageCrumbs}

	err = t.Execute(w, page)
	if err != nil {
		htmlError(w, 500, err.Error())
	}
}
