package middleware

import (
	"log/slog"
	"net/http"
	"time"

	"github.com/MicahParks/keyfunc/v3"
	"github.com/golang-jwt/jwt/v5"
)

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
		token, err := jwt.Parse(cookie.Value, k.Keyfunc)
		if err != nil {
			slog.Info("Failed parse JWT token", "error", err)
			http.Redirect(w, r, "/login", http.StatusSeeOther)
			return
		}
		if !token.Valid {
			slog.Info("The token is not valid")
			http.Redirect(w, r, "/login", http.StatusSeeOther)
			return
		}
		next.ServeHTTP(w, r)
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
