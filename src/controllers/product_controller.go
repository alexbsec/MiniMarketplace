package controllers

import (
	"net/http"

	"github.com/alexbsec/MiniMarketplace/src/db/config"
)

type ProductController struct {
    Service *config.Service
}

func (pc *ProductController) CreateProduct(w http.ResponseWriter, r *http.Request) {
    
}

