package handlers

import (
	"html/template"
	"log"
	"net/http"
)

type LoginTemplate struct {
	Title string
	URL   string
}

func Login(w http.ResponseWriter, r *http.Request) {
	tmpl := template.Must(template.ParseFiles("web/login.html"))
	data := LoginTemplate{
		Title: "Authorization with Google",
		URL:   "/auth",
	}
	err := tmpl.Execute(w, data)
	if err != nil {
		log.Println(err)
	}
}
