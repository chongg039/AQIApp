package main

import (
	"log"
	"net/http"
)

func main() {
	log.Println("Start server engine...")

	router := NewRouter()

	log.Fatal(http.ListenAndServe(":8088", router))
}
