package database

import (
	"context"
	"log"
	"os"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

func Connect() *pgxpool.Pool {
	dbURL := os.Getenv("DATABASE_URL")
	if dbURL == "" {
		dbURL = "postgres://admin:admin@db:5432/rinha"
	}
	config, err := pgxpool.ParseConfig(dbURL)
	if err != nil {
		log.Fatal("Erro parse config:", err)
	}
	config.MaxConns = 20
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	pool, err := pgxpool.NewWithConfig(ctx, config)
	if err != nil {
		log.Fatal("Erro conexão pool:", err)
	}
	return pool
}