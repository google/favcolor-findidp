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
