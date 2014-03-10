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
	"errors"
	"fmt"
	"net/http"
)

func init() {
	http.HandleFunc("/find", find)
	http.HandleFunc("/report", report)
	http.HandleFunc("/", main)
	http.HandleFunc("/set", set)
}

type Handles struct {
	db     Database
	logger Logger
	client *http.Client
}

func checkRequest(r *http.Request) (email *EMail, handles Handles, err error) {
	m := r.URL.Query()["m"]
	if len(m) == 0 {
		err = errors.New("missing m= argument")
		return
	}
	address := m[0]

	handles = Handles{getDatabase(r), getLogger(r), getClient(r)}
	email, err = parseEMail(address)
	return
}

func htmlError(w http.ResponseWriter, status int, message string) {
	w.WriteHeader(status)
	html := "<html><head><title>%s</title></head>" +
		"<body><h2>Error</h2><p>%s</p></body>"
	fmt.Fprintf(w, html, message, message)
}
