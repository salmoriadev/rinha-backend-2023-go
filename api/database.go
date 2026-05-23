package main

import (
	"context"
	"github.com/jackc/pgx/v5/pgxpool"
	"log"
	"os"
	"time"
)

func ConnectDatabase() (*pgxpool.Pool, error) {
	databaseURL := os.Getenv("DATABASE_URL")
	if databaseURL == "" {
		databaseURL = "postgres://rinha:rinha@localhost:5432/rinha?sslmode=disable"
	}
	var db *pgxpool.Pool
	var err error
	for i := 1; i <= 20; i++ {
		ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
		db, err = pgxpool.New(ctx, databaseURL)
		if err == nil {
			db.Ping(ctx)
		}
		cancel()
		if err == nil {
			log.Println("Connected to database")
			return db, nil
		}
		log.Printf("Trying to connect to database (attempt %d/20): %v", i, err)
	}
	log.Fatalf("Failed to connect to database after 20 attempts")
	return nil, err
}
