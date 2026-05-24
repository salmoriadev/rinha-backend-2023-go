package main

import (
	"context"
	"log"
	"os"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

func ConnectDatabase() (*pgxpool.Pool, error) {
	databaseURL := os.Getenv("DATABASE_URL")
	if databaseURL == "" {
		databaseURL = "postgres://rinha:rinha@localhost:5432/rinha?sslmode=disable"
	}

	cfg, err := pgxpool.ParseConfig(databaseURL)
	if err != nil {
		return nil, err
	}

	cfg.MaxConns = 20
	cfg.MinConns = 4
	cfg.MaxConnIdleTime = 30 * time.Second

	var lastErr error

	for i := 1; i <= 20; i++ {
		ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)

		db, err := pgxpool.NewWithConfig(ctx, cfg)
		if err == nil {
			err = db.Ping(ctx)
		}

		cancel()

		if err == nil {
			log.Println("Connected to database")
			return db, nil
		}

		if db != nil {
			db.Close()
		}

		lastErr = err
		log.Printf("Trying to connect to database (attempt %d/20): %v", i, err)
		time.Sleep(250 * time.Millisecond)
	}

	log.Printf("Failed to connect to database after 20 attempts")
	return nil, lastErr
}
