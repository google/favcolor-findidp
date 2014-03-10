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
