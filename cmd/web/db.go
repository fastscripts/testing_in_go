package main

import (
	"database/sql"
	"log"
	"strings"

	_ "github.com/jackc/pgconn"
	_ "github.com/jackc/pgx/v4"
	_ "github.com/jackc/pgx/v4/stdlib"
	_ "github.com/mattn/go-sqlite3"
)

func openDB(dsn string) (*sql.DB, error) {

	var db *sql.DB = nil
	var err error
	if strings.Contains(dsn, "sqlite") {
		db, err = sql.Open("sqlite3", dsn)
	} else {
		db, err = sql.Open("pgx", dsn)
	}
	if db == nil {
		log.Fatal("db is nil")
	}
	if err != nil {
		return nil, err
	}

	err = db.Ping()
	if err != nil {
		return nil, err
	}

	return db, nil
}

func (app *application) connectToDB() (*sql.DB, error) {
	connection, err := openDB(app.DSN)
	if err != nil {
		return nil, err
	}

	log.Println("Connected to Postgres!")

	return connection, nil
}
