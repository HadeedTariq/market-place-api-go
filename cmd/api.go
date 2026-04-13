package main

import (
	"log"
	"net/http"
	"time"

	repo "github.com/HadeedTariq/market-place-api-go/internal/adapters/postgresql/sqlc"
	"github.com/HadeedTariq/market-place-api-go/internal/auth"
	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/chi/v5"
	"github.com/jackc/pgx/v5"
)

// ~ so over there the application mount function exist which will setup the chi router
func (app *application) mount() http.Handler {
	r := chi.NewRouter()

	r.Use(middleware.RequestID) // important for rate limiting
	r.Use(middleware.RealIP)    // import for rate limiting and analytics and tracing
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer) // recover from crashes

	r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("all good"))
	})

	authValidator := auth.InitValidator()
	authService := auth.NewService(repo.New(app.db))
	authHandler := auth.NewHandler(authService, authValidator)
	r.Post("/auth/register-user", authHandler.RegisterUser)
	return r
}

func (app *application) run(h http.Handler) error {
	server := &http.Server{
		Addr:         app.config.addr,
		Handler:      h,
		WriteTimeout: time.Second * 30,
		ReadTimeout:  time.Second * 10,
		IdleTimeout:  time.Minute,
	}

	log.Printf("server has started at addr %s", app.config.addr)
	return server.ListenAndServe()
}

type config struct {
	addr     string
	dbConfig dbConfig
}

type application struct {
	config config
	db     *pgx.Conn
}

type dbConfig struct {
	dsn string
}
