package middleware

import (
	"context"
	"log/slog"
	"net/http"
	"time"

	"github.com/MicahParks/keyfunc/v3"
	"github.com/golang-jwt/jwt/v5"
)

type ClaimsUserData struct {
	Profile UserProfile `json:"user_metadata"`
	jwt.RegisteredClaims
}
type UserProfile struct {
	FullName  string `json:"full_name"`
	AvatarURL string `json:"avatar_url"`
}

type loggingResponseWriter struct {
	http.ResponseWriter
	statusCode int
}

func (lrw *loggingResponseWriter) WriteHeader(code int) {
	lrw.statusCode = code
	lrw.ResponseWriter.WriteHeader(code)
}

func Authmiddleware(k keyfunc.Keyfunc, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cookie, err := r.Cookie("access_token")
		if err != nil {
			slog.Info("Failed to retrieve cookie", "error", err)
			http.Redirect(w, r, "/login", http.StatusSeeOther)
			return
		}
		token, err := jwt.ParseWithClaims(cookie.Value, &ClaimsUserData{}, k.Keyfunc)
		if err != nil {
			slog.Info("Failed parse JWT token", "error", err)
			http.Redirect(w, r, "/login", http.StatusSeeOther)
			return
		}
		claims, ok := token.Claims.(*ClaimsUserData)
		if !ok {
			slog.Info("unknown claims type, cannot proceed")
		}
		subject, err := token.Claims.GetSubject()
		if err != nil {
			slog.Info("Failed get Subject", "error", err)
			http.Redirect(w, r, "/login", http.StatusSeeOther)
			return
		}
		ctx := context.WithValue(r.Context(), "Subject", subject)
		ctx = context.WithValue(ctx, "FullName", claims.Profile.FullName)
		ctx = context.WithValue(ctx, "AvatarURL", claims.Profile.AvatarURL)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func Logging(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		lrw := &loggingResponseWriter{ResponseWriter: w, statusCode: http.StatusOK}
		next.ServeHTTP(lrw, r)
		slog.Info("http request", "method", r.Method, "path", r.URL.Path, "status", lrw.statusCode, "time", time.Since(start))
	})
}
