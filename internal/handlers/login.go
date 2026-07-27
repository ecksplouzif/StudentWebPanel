package handlers

import (
	"html/template"
	"log"
	"net/http"
	"os"
)

type LoginTemplate struct {
	Title string
	URL   string
}

func Login(w http.ResponseWriter, r *http.Request) {
	tmpl := template.Must(template.ParseFiles("web/login.html"))
	authURL := os.Getenv("PROVIDER_GOOGLE_URL")
	data := LoginTemplate{
		Title: "Authorization with Google",
		URL:   authURL,
	}
	err := tmpl.Execute(w, data)
	if err != nil {
		log.Println(err)
	}
}
