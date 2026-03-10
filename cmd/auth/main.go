package main

import (
	"log"
	"net/http"
)

func main() {
	log.Print("starting server....")
	if err := http.ListenAndServe("localhost:8080", nil); err != nil {
		panic(err)
	}
}
