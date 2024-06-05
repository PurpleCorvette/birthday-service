package db

import (
	"context"

	"github.com/jackc/pgx/v4/pgxpool"

	"birthday-service/pkg/logging"
)

var Conn *pgxpool.Pool

func ConnectDatabase(url string, log *logging.Logger) {
	var err error
	Conn, err = pgxpool.Connect(context.Background(), url)
	if err != nil {
		log.Fatalf("Unable to connect to database: %v\n", err)
	}

	log.Infoln("Connected to th database!")
}

func CloseDatabase(log *logging.Logger) {
	Conn.Close()
	log.Infoln("Database connection closed")
}
