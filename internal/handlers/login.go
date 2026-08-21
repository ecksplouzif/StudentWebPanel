package handlers

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
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
	b := make([]byte, 32)
	_, err := rand.Read(b)
	if err != nil {
		log.Println(err)
		return
	}
	verifierBase64 := base64.RawURLEncoding.EncodeToString(b)
	sum := sha256.Sum256([]byte(verifierBase64))
	challenge := base64.RawURLEncoding.EncodeToString(sum[:])
	swp := http.Cookie{
		Name:     "pkce_verifier",
		Value:    verifierBase64,
		Path:     "/",
		MaxAge:   300,
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteLaxMode,
	}
	http.SetCookie(w, &swp)

	tmpl := template.Must(template.ParseFiles("web/login.html"))
	authURL := os.Getenv("PROVIDER_GOOGLE_URL") + challenge + "&code_challenge_method=S256"
	data := LoginTemplate{
		Title: "Authorization with Google",
		URL:   authURL,
	}
	err = tmpl.Execute(w, data)
	if err != nil {
		log.Println(err)
	}
}
