package models

import (
	"log/slog"

	"github.com/alexbsec/MiniMarketplace/src/db/config"
	"github.com/alexbsec/MiniMarketplace/src/logging"
)

type Product struct {
	ID          uint `gorm:"primaryKey"`
	Name        string
	Description string
	Price       float64
    Points      uint
	Category    string
}

type ProductView struct {
    Product *Product
    service *config.Service
}

func (pv *ProductView) Create() error {
    log := slog.New(logging.NewHandler(nil))
    dbGorm, err := pv.service.Db()
    if err != nil {
        return err
    }
    
    if err := dbGorm.Create(pv.Product).Error; err != nil {
        log.Error("Failed to retrieve DB from GORM", slog.String("error", err.Error()))
        return err
    } 

    log.Info("Product created successfully", slog.String("product_name", pv.Product.Name))
    return nil
}
