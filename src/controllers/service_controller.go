package controllers

import (
	"fmt"

	"github.com/alexbsec/MiniMarketplace/src/db/config"
	"github.com/alexbsec/MiniMarketplace/src/db/models"
)

var (
    productService  *models.ProductService
    userService     *models.UserService
    walletService   *models.WalletService
)

func init() {
	service, err := config.InitService()
	if err != nil {
		panic(fmt.Sprintf("Failed to initialize database service: %v", err))
	}
    productService = &models.ProductService{Service: service}
	userService = &models.UserService{Service: service}
    walletService = &models.WalletService{Service: service}
}
