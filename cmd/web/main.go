package main

import (
	"log"
	"net/http"
)

type application struct{}

func main() {

	//set up an app confgiuration
	app := application{}

	//get application routes

	mux := app.routes()

	log.Println("Starting server port 8080... ")

	//start the server

	err := http.ListenAndServe(":8080", mux)
	if err != nil {
		log.Fatal(err)
	}
}
