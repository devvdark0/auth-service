package main

import "net/http"

func main() {
	//TODO: init config

	//TODO: init logger

	//TODO: init storage

	if err := http.ListenAndServe("localhost:8000", nil); err != nil {
		panic(err)
	}

}
