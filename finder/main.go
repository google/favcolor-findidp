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
