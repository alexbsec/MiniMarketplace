package controllers

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/alexbsec/MiniMarketplace/src/db/config"
	"github.com/alexbsec/MiniMarketplace/src/db/models"
)

var productService *models.ProductService

func init() {
    service, err := config.InitService()
    if err != nil {
        panic(fmt.Sprintf("Failed to initialize database service: %v", err))
    }
    productService = &models.ProductService{Service: service}
}

// Handle 
func HandleProducts(w http.ResponseWriter, r *http.Request) {
    switch r.Method {
    case http.MethodPost:
        handleCreateProduct(w, r)
    default:
        http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
    }
}


func handleCreateProduct(w http.ResponseWriter, r *http.Request) {
    var product models.Product
    if err := json.NewDecoder(r.Body).Decode(&product); err != nil {
        http.Error(w, "Invalid request payload", http.StatusBadRequest)
        return
    }

    if err := productService.Create(&product); err != nil {
        http.Error(w, "Failed to create product", http.StatusInternalServerError)
        return
    }
    
    w.WriteHeader(http.StatusCreated)
    json.NewEncoder(w).Encode(product)
}
