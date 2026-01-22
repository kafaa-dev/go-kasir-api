package main

import (
	"encoding/json"
	"log"
	"net/http"
)

type Product struct {
	ID    int    `json:"id"`
	Name  string `json:"name"`
	Price int    `json:"price"`
	Stock int    `json:"stock"`
}

var products = []Product{
	{ID: 1, Name: "Mie rebus", Price: 3500, Stock: 10},
	{ID: 2, Name: "Air minum 800ml", Price: 3000, Stock: 40},
	{ID: 3, Name: "Kecap", Price: 12000, Stock: 20},
}

func main() {
	http.HandleFunc("GET /health", checkHealthHandler)

	http.HandleFunc("GET /api/products", listProductsHandler)

	log.Println("Server listening on port 8080")

	log.Fatal(http.ListenAndServe(":8080", nil))
}

func checkHealthHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	json.NewEncoder(w).Encode(map[string]string{
		"status":  "OK",
		"message": "API Running",
	})
}

func listProductsHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	json.NewEncoder(w).Encode(products)
}
