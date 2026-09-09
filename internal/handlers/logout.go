package handlers

import (
	"net/http"
)

func Logout(w http.ResponseWriter, r *http.Request) {
	swp := http.Cookie{
		Name:     "access_token",
		Value:    "",
		Path:     "/",
		MaxAge:   -1,
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteLaxMode,
	}
	http.SetCookie(w, &swp)

	http.Redirect(w, r, "/login", http.StatusSeeOther)
}
