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
