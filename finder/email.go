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
	"strings"
)

type EMail struct {
	Name, Domain string
	IDPs         []IDP
}

func parseEMail(in string) (e *EMail, err error) {
	at := strings.Index(in, "@")
	if at == -1 {
		err = errors.New("malformed email")
		return
	}
	e = new(EMail)
	e.Name = in[:at]
	e.Domain = in[at+1:]
	return
}

func (e EMail) parts() (name, domain string) {
	return e.Name, e.Domain
}

func (e EMail) string() string {
	return e.Name + "@" + e.Domain
}
