package main

import (
	"log"
	"net/http"
	"strconv"
	"time"

	utils "github.com/HadeedTariq/market-place-api-go/internal"
	repo "github.com/HadeedTariq/market-place-api-go/internal/adapters/postgresql/sqlc"
	"github.com/HadeedTariq/market-place-api-go/internal/auth"
	"github.com/HadeedTariq/market-place-api-go/internal/mail"
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

	senderUsername := utils.GetEnv("MAIL_SENDER_USERNAME", "apikey")          // For SendGrid, this is usually "apikey"
	senderPassword := utils.GetEnv("MAIL_SENDER_PASSWORD", "your_sg_api_key") // Your actual API Key
	senderHost := utils.GetEnv("MAIL_SENDER_SERVER", "smtp.sendgrid.net")
	senderPort := utils.GetEnv("MAIL_SENDER_PORT", "587")
	senderPortInt, _ := strconv.Atoi(senderPort)

	mailSenderService := mail.NewMailer(senderHost, senderPortInt, senderUsername, senderPassword, "computeranalog351@gmail.com")

	emailService := mail.NewEmailService(mailSenderService)
	authValidator := auth.InitValidator()
	authService := auth.NewService(repo.New(app.db))
	authHandler := auth.NewHandler(authService, authValidator, emailService)
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
