package main

import (
	"WebPanel/internal/handlers"
	"WebPanel/internal/middleware"
	"context"
	"io"
	"log/slog"
	"net/http"
	"os"
	"time"

	"github.com/MicahParks/keyfunc/v3"
	"github.com/jackc/pgx/v5/pgxpool"
)

func Conectdatabase() *pgxpool.Pool {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	conn, err := pgxpool.New(ctx, os.Getenv("CONECT_DATABASE_URL"))
	if err != nil {
		slog.Error("Failed conect to Database", "error", err)
	}
	return conn
}

func GetJWKS() keyfunc.Keyfunc {
	project_URL := os.Getenv("PROJECT_URL") + "/auth/v1/.well-known/jwks.json"
	res, err := http.Get(project_URL)
	if err != nil {
		slog.Info("Failed get public key Supabase", "error", err)
	}
	body, err := io.ReadAll(res.Body)
	if err != nil {
		slog.Info("Failed read body", "error", err)
	}
	defer func() { _ = res.Body.Close() }()
	k, err := keyfunc.NewJWKSetJSON(body)
	if err != nil {
		slog.Info("Failed to create a keyfunc.Keyfunc", "error", err)
	}
	return k
}

func main() {
	mux := http.NewServeMux()
	Mykeyfunc := GetJWKS()
	conect := Conectdatabase()

	mux.HandleFunc("/healthz", handlers.Healthz)
	mux.HandleFunc("/login", handlers.Login)
	mux.HandleFunc("/auth/callback", handlers.Callback(conect))
	mux.HandleFunc("/auth", handlers.GenerateVerefiToken)
	mux.Handle("/profile", middleware.Authmiddleware(Mykeyfunc, http.HandlerFunc(handlers.Profile)))
	mux.Handle("/", middleware.Authmiddleware(Mykeyfunc, http.HandlerFunc(handlers.NotFound)))

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	serv := &http.Server{
		Addr:         ":" + port,
		Handler:      middleware.Logging(mux),
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
	}
	err := serv.ListenAndServe()
	if err != nil {
		slog.Info("start server", "error", err)
	}
}
