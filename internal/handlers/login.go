package handlers

import (
	"html/template"
	"log/slog"
	"net/http"
)

type LoginTemplate struct {
	Title string
	URL   string
}

func Login(tmpl *template.Template) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		data := LoginTemplate{
			Title: "Authorization with Google",
			URL:   "/auth",
		}
		err := tmpl.ExecuteTemplate(w, "login.html", data)
		if err != nil {
			slog.Error("Isert data to template", "error", err)
		}
	}
}
