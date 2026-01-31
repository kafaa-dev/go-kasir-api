package main

import (
	"encoding/json"
	"kasir-api/database"
	"kasir-api/handlers"
	"kasir-api/repositories"
	"kasir-api/services"
	"log"
	"net/http"
	"os"
)

type Config struct {
	Port string
}

func main() {
	config := Config{
		Port: os.Getenv("PORT"),
	}

	http.HandleFunc("GET /health", checkHealthHandler)

	productRepository := repositories.NewProductRepository(database.Products)
	productService := services.NewProductService(productRepository)
	productHandler := handlers.NewProductHandler(productService)

	http.HandleFunc("GET /api/products", productHandler.GetAll)
	http.HandleFunc("POST /api/products", productHandler.Create)
	http.HandleFunc("GET /api/products/{id}", productHandler.GetByID)
	http.HandleFunc("PUT /api/products/{id}", productHandler.Update)
	http.HandleFunc("DELETE /api/products/{id}", productHandler.Delete)

	categoryRepository := repositories.NewCategoryRepository(database.Categories)
	categoryService := services.NewCategoryService(categoryRepository)
	categoryHandler := handlers.NewCategoryHandler(categoryService)

	http.HandleFunc("GET /api/categories", categoryHandler.GetAll)
	http.HandleFunc("POST /api/categories", categoryHandler.Create)
	http.HandleFunc("GET /api/categories/{id}", categoryHandler.GetByID)
	http.HandleFunc("PUT /api/categories/{id}", categoryHandler.Update)
	http.HandleFunc("DELETE /api/categories/{id}", categoryHandler.Delete)

	log.Println("Server listening on port " + config.Port)

	log.Fatal(http.ListenAndServe(":"+config.Port, nil))
}

func checkHealthHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	json.NewEncoder(w).Encode(map[string]string{
		"status":  "OK",
		"message": "API Running",
	})
}
