package handlers

import (
	"bytes"
	"encoding/json"
	"io"
	"log/slog"
	"net/http"
	"os"
	"time"
)

type TokenRequest struct {
	AuthCode     string `json:"auth_code"`
	CodeVerifier string `json:"code_verifier"`
}

type User struct {
	ID    string `json:"id"`
	Email string `json:"email"`
}

type tokenResponse struct {
	AccessToken string `json:"access_token"`
	User        User   `json:"user"`
}

func Callback(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie("pkce_verifier")
	if err != nil {
		slog.Error("Failed to retrieve cookie", "error", err)
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	tokenURL := os.Getenv("TOKEN_URL")
	supaAPIKey := os.Getenv("SUPABASE_API_KEY")
	takeAPIKey := TokenRequest{
		AuthCode:     r.URL.Query().Get("code"),
		CodeVerifier: cookie.Value,
	}

	b, err := json.Marshal(takeAPIKey)
	if err != nil {
		slog.Error("Failed to marshal request", "error", err)
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	req, err := http.NewRequest("POST", tokenURL, bytes.NewBuffer(b))
	if err != nil {
		slog.Error("Failed to create request", "error", err)
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("apikey", supaAPIKey)

	client := &http.Client{
		Timeout: 10 * time.Second,
	}
	resp, err := client.Do(req)
	if err != nil {
		slog.Error("Failed to send request", "error", err)
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}
	defer func() { _ = resp.Body.Close() }()
	if resp.StatusCode != http.StatusOK {
		respBody, readErr := io.ReadAll(resp.Body)
		if readErr != nil {
			slog.Error("Failed to read error response", "error", readErr)
		}
		slog.Error("Non-OK response from token endpoint", "status", resp.StatusCode, "body", respBody)
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	var tokenResp tokenResponse
	err = json.NewDecoder(resp.Body).Decode(&tokenResp)
	if err != nil {
		slog.Error("Failed to decode response", "error", err)
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	swp := http.Cookie{
		Name:     "pkce_verifier",
		Value:    "",
		Path:     "/",
		MaxAge:   -1,
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteLaxMode,
	}
	http.SetCookie(w, &swp)

	swp = http.Cookie{
		Name:     "access_token",
		Value:    tokenResp.AccessToken,
		Path:     "/",
		MaxAge:   300,
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteLaxMode,
	}
	http.SetCookie(w, &swp)
	http.Redirect(w, r, "/", http.StatusSeeOther)
}
