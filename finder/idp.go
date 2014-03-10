package findIDP

import ()

type IDPKey string

const (
	Google    IDPKey = "Google"
	Microsoft IDPKey = "Microsoft"
	Facebook  IDPKey = "Facebook"
	Yahoo     IDPKey = "Yahoo"
)

type Protocol int

const (
	Unknown Protocol = iota
	OIDC
	OpenID2
	Persona
	SAML
)

func (p Protocol) name() string {
	return protocolNames[p]
}

var protocolNames = [...]string{
	"unknown",
	"oidc",
	"openid2",
	"persona",
	"saml",
}

var baseURIs = map[IDPKey]string{
	Google:    "https://accounts.google.com",
	Microsoft: "https://login.live.com",
	Facebook:  "https://www.facebook.com",
	Yahoo:     "https://login.yahoo.com",
}

func wellKnownIDP(key IDPKey) (idpp *IDP) {
	switch key {
	case Google:
		idpp = &IDP{Google, baseURIs[Google], OIDC}
	case Microsoft:
		idpp = &IDP{Microsoft, baseURIs[Microsoft], Unknown}
	case Facebook:
		idpp = &IDP{Facebook, baseURIs[Facebook], Unknown}
	case Yahoo:
		idpp = &IDP{Yahoo, baseURIs[Yahoo], OpenID2}
	default:
		idpp = nil
	}
	return
}

type reported struct {
	IDP      string
	Protocol string
	Verified bool
}

func reportOne(idp IDP, verified bool) reported {
	return reported{idp.URI, idp.Protocol.name(), verified}
}

type IDP struct {
	Key      IDPKey
	URI      string
	Protocol Protocol
}
