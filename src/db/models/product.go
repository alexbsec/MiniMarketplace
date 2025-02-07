package models

import (
	"fmt"
	"log/slog"

	"github.com/alexbsec/MiniMarketplace/src/db/config"
	"github.com/alexbsec/MiniMarketplace/src/db/models/utils"
	"github.com/alexbsec/MiniMarketplace/src/logging"
	"gorm.io/gorm"
)

type Product struct {
	ID          uint     `gorm:"primaryKey"`
	Name        *string  `json:"name"`
	Description *string  `json:"description"`
	Price       *float64 `json:"price"`
	Points      *uint    `json:"points"`
	Category    *string  `json:"category"`
	Stock       *uint    `json:"stock"`
}

type ProductService struct {
	Service *config.Service
}

func (ps *ProductService) Create(product *Product) error {
    if !ps.isServiceRunning() {
        return fmt.Errorf("Cannot proceed because service is offline") 
    }

	return models_utils.DoTransaction(ps.Service, models_utils.CREATE, func(tx *gorm.DB) error {
		if err := tx.Create(product).Error; err != nil {
			return fmt.Errorf("failed to create product: %w", err)
		}

		return nil
	})
}

func (ps *ProductService) Fetch(id uint) (*Product, error) {
    if !ps.isServiceRunning() {
        return nil, fmt.Errorf("Cannot proceed because service is offline") 
    }

	dbGorm, err := ps.Service.Db()
	if err != nil {
		return nil, err
	}

	var product Product

	res := dbGorm.First(&product, id)
	if res.Error != nil {
		if res.Error == gorm.ErrRecordNotFound {
			logging.Log.Info("Produto não encontrado", slog.String("error", res.Error.Error()))
			return nil, res.Error
		}

		logging.Log.Error("Error while searching the product", slog.String("error", res.Error.Error()))
		return nil, res.Error
	}

	return &product, nil
}

func (ps *ProductService) Update(id uint, newProduct *Product) error {
    if !ps.isServiceRunning() {
        return fmt.Errorf("Cannot proceed because service is offline") 
    }

	return models_utils.DoTransaction(ps.Service, models_utils.UPDATE, func(tx *gorm.DB) error {
		var product Product
		if err := tx.First(&product, id).Error; err != nil {
			return fmt.Errorf("product with id %d not found: %w", id, err)
		}

		if err := tx.Model(&product).Updates(newProduct).Error; err != nil {
			return fmt.Errorf("failed to update product with id %d: %w", id, err)
		}

		return nil
	})
}

func (ps *ProductService) Delete(id uint) error {
    if !ps.isServiceRunning() {
        return fmt.Errorf("Cannot proceed because service is offline") 
    }
	return models_utils.DoTransaction(ps.Service, models_utils.DELETE, func(tx *gorm.DB) error {
		if err := tx.Delete(&Product{}, id).Error; err != nil {
			return fmt.Errorf("failed to delete product: %w", err)
		}

		return nil
	})
}


func (ps *ProductService) isServiceRunning() bool {
    if ps.Service == nil {
        logging.Log.Error("Product Service is not initialized! Aborting")
    }

    return ps.Service != nil
}
