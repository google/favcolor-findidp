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
	"html/template"
	"net/http"
)

type MainPage struct{}

func main(w http.ResponseWriter, r *http.Request) {
	t, err := template.ParseFiles("templates/Main.html")
	if err != nil {
		htmlError(w, 500, err.Error())
		return
	}
	page := MainPage{}
	err = t.Execute(w, page)
	if err != nil {
		htmlError(w, 500, err.Error())
	}
}
