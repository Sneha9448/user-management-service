package database

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

var DB *pgxpool.Pool

func ConnectDB(databaseURL string) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	config, err := pgxpool.ParseConfig(databaseURL)
	if err != nil {
		log.Fatalf("Unable to parse database URL: %v", err)
	}

	// Optimization: Increase max connections for higher throughput
	config.MaxConns = 25
	config.MinConns = 5
	config.MaxConnLifetime = time.Hour
	config.MaxConnIdleTime = 10 * time.Minute

	// Create the connection pool
	DB, err = pgxpool.NewWithConfig(context.Background(), config)
	if err != nil {
		log.Fatalf("Unable to create connection pool: %v", err)
	}

	// Verify connection with a ping using the timeout context
	if err := DB.Ping(ctx); err != nil {
		log.Fatalf("Unable to connect to database (ping failed): %v\nCheck your DATABASE_URL and ensure PostgreSQL is running.", err)
	}

	fmt.Println("Connected to PostgreSQL successfully")
}

func CloseDB() {
	if DB != nil {
		DB.Close()
		fmt.Println("Database connection closed")
	}
}
