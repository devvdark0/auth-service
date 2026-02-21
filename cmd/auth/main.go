package main

import (
	"log"
	"net/http"
)

func main() {

	log.Println("starting server...")
	if err := http.ListenAndServe("localhost:8000", nil); err != nil {
		panic(err)
	}

}
