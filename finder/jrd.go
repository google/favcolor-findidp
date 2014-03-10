package findIDP

import ()

func idpFromJRD(f interface{}) *IDP {
	idpURI := linkFromJRD(f, "http://openid.net/specs/connect/1.0/issuer")
	if idpURI == "" {
		return nil
	}
	return &IDP{IDPKey("WebFinger"), idpURI, OIDC}
}

func linkFromJRD(f interface{}, relName string) string {
	if f == nil {
		return ""
	}
	m := f.(map[string]interface{})
	f = m["links"]
	links := f.([]interface{})
	for _, linkI := range links {
		link := linkI.(map[string]interface{})
		rel := link["rel"].(string)
		if rel == relName {
			return link["href"].(string)
		}
	}
	return ""
}
