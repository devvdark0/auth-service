package main

import (
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
)

func main() {
	r := chi.NewRouter()
	r.Get("/home", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Hello, World!"))
	})
	log.Println("server started on port :81...")
	if err := http.ListenAndServe(":81", r); err != nil {
		log.Fatal(err)
	}

}
