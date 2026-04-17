package main

import (
	"context"
	"log/slog"
	"os"

	utils "github.com/HadeedTariq/market-place-api-go/internal"
	"github.com/jackc/pgx/v5"
	"github.com/joho/godotenv"
)

func main() {
	ctx := context.Background()
	err := godotenv.Load()

	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))
	slog.SetDefault(logger)
	if err != nil {
		logger.Error("env not loaded")
	}

	cfg := config{
		addr: ":8080",
		dbConfig: dbConfig{
			dsn: utils.GetEnv("GOOSE_DBSTRING", "host=localhost user=postgres password=postgres dbname=marketplace sslmode=disable"),
		}}

	conn, err := pgx.Connect(ctx, cfg.dbConfig.dsn)

	if err != nil {
		panic(err)
	}

	// commandTag, err := conn.Exec(ctx, "delete from users where role = 'user'")

	// if err != nil {
	// 	log.Printf("Update failed: %v", err)
	// }
	// fmt.Printf("Rows affected: %d\n", commandTag.RowsAffected())

	defer conn.Close(ctx)

	logger.Info("connected to database", "dsn", cfg.dbConfig.dsn)

	api := application{
		config: cfg,
		db:     conn,
	}

	if err := api.run(api.mount()); err != nil {
		slog.Error("Error starting server", "error", err)
		os.Exit(1)
	}
}
