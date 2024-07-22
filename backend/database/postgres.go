package database

import (
	"context"
	"errors"
	"fmt"
	"os"

	"github.com/jackc/pgx/v5"
)

var conn *pgx.Conn

func StartDB() error {
	db_url := os.Getenv("DB_URL")
	if db_url == "" {
		return errors.New("DB_URL is not set")
	}
	fmt.Println("DB_URL is set to: ", db_url)
	var err error
	conn, err = pgx.Connect(context.Background(), db_url)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to connect to database: %v\n", err)
		return errors.New("can't connect to database")
	}

	err = conn.Ping(context.Background())
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to verify a connection: %v\n", err)
		return errors.New("can't verify a connection")
	}

	return nil
}

func CloseDB() error {
	err := conn.Close(context.Background())
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to close the connection: %v\n", err)
		return errors.New("can't close the connection")
	}

	return nil
}
