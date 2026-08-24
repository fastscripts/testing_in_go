package main

import (
	"log"
	"os"
	"testing"

	"github.com/fastscripts/testing_in_go/db"
)

var app application

func TestMain(m *testing.M) {

	pathToTemplates = "./../../templates/"
	app.Session = getSession()
	app.DSN = "./../../test.sqlite"

	conn, err := app.connectToDB()
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	app.DB = db.PostgresConn{DB: conn}

	os.Exit(m.Run())
}
