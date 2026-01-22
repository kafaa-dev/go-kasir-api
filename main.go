package main

import (
	"log"
	"net/http"
)

func main() {
	log.Println("Server listening on port 8080")

	log.Fatal(http.ListenAndServe(":8080", nil))
}
