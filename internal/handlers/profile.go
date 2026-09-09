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
	// AvatarURL string dont need rigt now
}

func Profile(w http.ResponseWriter, r *http.Request) {
	subject := r.Context().Value("Subject")
	name := r.Context().Value("FullName")
	//	avatarURL := r.Context().Value("AvatarURL") dont need rigt now
	tmpl := template.Must(template.ParseFiles("web/profile.html"))
	data := ProfileTemplate{
		UserSubject: subject,
		UserName:    name,
		URL:         "/logout",
		Title:       "Logout",
	}
	err := tmpl.Execute(w, data)
	if err != nil {
		slog.Info("Failed execute data", "error", err)
	}
}
