package main

import (
	"log"
	"net/http"

	"github.com/alexedwards/scs/v2"
)

type application struct {
	Session *scs.SessionManager
}

func main() {

	//set up an app confgiuration
	app := application{}

	// get a session manager
	app.Session = getSession()

	// get a session manager
	app.Session = getSession()

	//get application routes
	mux := app.routes()

	log.Println("Starting server port 8080... ")

	//start the server

	err := http.ListenAndServe(":8080", mux)
	if err != nil {
		log.Fatal(err)
	}
}
