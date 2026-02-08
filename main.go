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
	"strings"

	"github.com/spf13/viper"
)

type Config struct {
	Port        string `mapstructure:"PORT"`
	DatabaseUrl string `mapstructure:"DATABASE_URL"`
}

func main() {
	viper.AutomaticEnv()
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	if _, err := os.Stat(".env"); err == nil {
		viper.SetConfigFile(".env")
		_ = viper.ReadInConfig()
	}

	config := Config{
		Port:        viper.GetString("PORT"),
		DatabaseUrl: viper.GetString("DATABASE_URL"),
	}

	db, err := database.InitDB(config.DatabaseUrl)
	if err != nil {
		log.Fatal("Failed to initialize database: ", err)
	}
	defer db.Close()

	http.HandleFunc("GET /health", checkHealthHandler)

	productRepository := repositories.NewProductRepository(db)
	productService := services.NewProductService(productRepository)
	productHandler := handlers.NewProductHandler(productService)

	http.HandleFunc("GET /api/products", productHandler.GetAll)
	http.HandleFunc("POST /api/products", productHandler.Create)
	http.HandleFunc("GET /api/products/{id}", productHandler.GetByID)
	http.HandleFunc("PUT /api/products/{id}", productHandler.Update)
	http.HandleFunc("DELETE /api/products/{id}", productHandler.Delete)

	categoryRepository := repositories.NewCategoryRepository(db)
	categoryService := services.NewCategoryService(categoryRepository)
	categoryHandler := handlers.NewCategoryHandler(categoryService)

	http.HandleFunc("GET /api/categories", categoryHandler.GetAll)
	http.HandleFunc("POST /api/categories", categoryHandler.Create)
	http.HandleFunc("GET /api/categories/{id}", categoryHandler.GetByID)
	http.HandleFunc("PUT /api/categories/{id}", categoryHandler.Update)
	http.HandleFunc("DELETE /api/categories/{id}", categoryHandler.Delete)

	transactionRepository := repositories.NewTransactionRepository(db)
	transactionService := services.NewTransactionService(transactionRepository)
	transactionHandler := handlers.NewTransactionHandler(transactionService)

	http.HandleFunc("POST /api/checkout", transactionHandler.Checkout)
	http.HandleFunc("GET /api/report/today", transactionHandler.ReportToday)

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
