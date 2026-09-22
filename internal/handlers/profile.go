package handlers

import (
	"html/template"
	"log/slog"
	"net/http"
)

type ProfileTemplate struct {
	UserSubject interface{}
	UserName    interface{}
	URL         string
	Title       string
	// AvatarURL string dont need right now
}

func Profile(tmpl *template.Template) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		subject := r.Context().Value("Subject")
		name := r.Context().Value("FullName")
		//	avatarURL := r.Context().Value("AvatarURL") dont need right now
		data := ProfileTemplate{
			UserSubject: subject,
			UserName:    name,
			URL:         "/logout",
			Title:       "Logout",
		}
		err := tmpl.ExecuteTemplate(w, "profile.html", data)
		if err != nil {
			slog.Error("Isert data to template", "error", err)
		}
	}
}
