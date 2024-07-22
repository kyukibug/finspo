package database

import (
	"context"
	"errors"
	"fmt"
	"os"

	"github.com/jackc/pgx/v5/pgxpool"
)

var dbpool *pgxpool.Pool

func GetDB() *pgxpool.Pool {
	return dbpool
}

func StartDB() error {
	db_url := os.Getenv("DB_URL")
	if db_url == "" {
		return errors.New("DB_URL is not set")
	}

	var err error
	dbpool, err = pgxpool.New(context.Background(), db_url)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to connect to database: %v\n", err)
		return errors.New("can't connect to database")
	}

	err = dbpool.Ping(context.Background())
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to verify a connection: %v\n", err)
		return errors.New("can't verify a connection")
	}

	return nil
}

func CloseDB() {
	dbpool.Close()
	// maybe error check for this but .Close does not return an errro??
}
