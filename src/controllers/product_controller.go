package controllers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/alexbsec/MiniMarketplace/src/db/models"
)

// Handle 
func HandleProducts(w http.ResponseWriter, r *http.Request) {
    switch r.Method {
    case http.MethodPost:
        handleCreateProduct(w, r)
    case http.MethodGet:
        handleFetchProduct(w, r)
    case http.MethodPut:
        handleUpdateProduct(w, r)
    case http.MethodDelete:
        handleDeleteProduct(w, r)
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

func handleFetchProduct(w http.ResponseWriter, r *http.Request) {
    idStr := r.URL.Query().Get("id")
    id, err := strconv.Atoi(idStr)
    if err != nil {
        http.Error(w, "Invalid product ID", http.StatusBadRequest)
        return
    }

    product, err := productService.Fetch(uint(id)) 
    if err != nil {
        http.Error(w, "Failed to fetch product", http.StatusNotFound)
        return
    }

    if product == nil {
        http.Error(w, "Produto não encontrado", http.StatusNotFound)
        return
    }

    json.NewEncoder(w).Encode(product)
}


func handleUpdateProduct(w http.ResponseWriter, r *http.Request) {
    idStr := r.URL.Path[len("/products/"):]
    id, err := strconv.Atoi(idStr)
    if err != nil {
        http.Error(w, "Invalid product ID", http.StatusBadRequest)
        return
    }

    product, err := productService.Fetch(uint(id))
    if err != nil {
        http.Error(w, "Failed to find product", http.StatusNotFound)
        return
    }
    
    if product == nil {
        http.Error(w, "Produto não encontrado", http.StatusNotFound)
        return
    }

    var newProduct models.Product
    if err = json.NewDecoder(r.Body).Decode(&newProduct); err != nil {
        http.Error(w, "Invalid parameters when updating", http.StatusBadRequest)
        return
    }

    // Only update non-nil fields
    if newProduct.Name != nil {
        product.Name = newProduct.Name
    }
    if newProduct.Description != nil {
        product.Description = newProduct.Description
    }
    if newProduct.Price != nil {
        product.Price = newProduct.Price
    }
    if newProduct.Points != nil {
        product.Points = newProduct.Points
    }
    if newProduct.Category != nil {
        product.Category = newProduct.Category
    }

    if err = productService.Update(uint(id), product); err != nil {
        http.Error(w, "Failed to update product", http.StatusInternalServerError)
        return
    }

    json.NewEncoder(w).Encode(product)
}

func handleDeleteProduct(w http.ResponseWriter, r *http.Request) {
    idStr := r.URL.Path[len("/products/"):]
    id, err := strconv.Atoi(idStr)
    if err != nil {
        http.Error(w, "Invalid product ID", http.StatusBadRequest)
        return
    }

    product, err := productService.Fetch(uint(id))
    if err != nil {
        http.Error(w, "Failed to find product", http.StatusNotFound)
        return
    }
    
    if product == nil {
        http.Error(w, "Produto não encontrado", http.StatusNotFound)
        return
    }

    if err = productService.Delete(uint(id)); err != nil {
        http.Error(w, "Failed to delete product", http.StatusInternalServerError)
        return
    }

    w.WriteHeader(http.StatusOK)
}
