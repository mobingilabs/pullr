package main

import (
	"github.com/mobingilabs/pullr/cmd/eventsrv/app"
)

func main() {
	app.Execute()

	// r := mux.NewRouter()
	// api := r.PathPrefix("/api").Subrouter()
	// api.HandleFunc("/github", LogRequest("githubHandler", githubHandler)).Methods("POST")
	// api.HandleFunc("/version", LogRequest("version", version))
}
