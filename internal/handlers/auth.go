package handlers

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"log"
	"net/http"
	"os"
)

func GenerateVerefiToken(w http.ResponseWriter, r *http.Request) {
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
	authURL := os.Getenv("PROVIDER_GOOGLE_URL") + challenge + "&code_challenge_method=S256"
	http.Redirect(w, r, authURL, http.StatusSeeOther)
}
