package main

import (
	"encoding/gob"
	"flag"
	"log"
	"net/http"

	"github.com/alexedwards/scs/v2"
	"github.com/fastscripts/testing_in_go/data"
	"github.com/fastscripts/testing_in_go/db"
)

type application struct {
	DSN     string
	DB      db.PostgresConn
	Session *scs.SessionManager
}

func main() {

	// register the data type for the session manager
	gob.Register(data.User{})

	//set up an app confgiuration
	app := application{}

	//flag.StringVar(&app.DSN, "dsn", "host=localhost port=5432 user=postgres password=postgres dbname=users sslmode=disable timezone=UTC connect_timeout=5", "Posgtres connection")
	flag.StringVar(&app.DSN, "dsn", "test.sqlite", "SQLite connection")

	flag.Parse()

	conn, err := app.connectToDB()
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	app.DB = db.PostgresConn{DB: conn}

	// get a session manager
	app.Session = getSession()

	// get a session manager
	app.Session = getSession()

	//get application routes
	mux := app.routes()

	log.Println("Starting server port 8080... ")

	//start the server

	err = http.ListenAndServe(":8080", mux)
	if err != nil {
		log.Fatal(err)
	}
}
