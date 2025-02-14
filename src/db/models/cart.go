package models

import (
	"fmt"
	"log/slog"

	"github.com/alexbsec/MiniMarketplace/src/db/config"
	"github.com/alexbsec/MiniMarketplace/src/db/models/utils"
	"github.com/alexbsec/MiniMarketplace/src/logging"
	"gorm.io/gorm"
)

type Cart struct {
	ID        uint      `gorm:"primaryKey"`
	Items     *string   `json:"items"`
	Total     *float64  `json:"total"`
	UserID    uint      `gorm:"not null" json:"user_id"`
	User      User      `gorm:"foreignKey:UserID"`
}

type CartService struct {
	Service *config.Service
}

func (cs *CartService) Create(cart *Cart) error {
	if !cs.isServiceRunning() {
		return fmt.Errorf("Cannot proceed because service is offline")
	}

	return models_utils.DoTransaction(cs.Service, models_utils.CREATE, func(tx *gorm.DB) error {
		if err := tx.Create(cart).Error; err != nil {
			return fmt.Errorf("failed to create cart: %w", err)
		}

		return nil
	})
}

func (cs *CartService) FetchByUser(userID uint) (*Cart, error) {
	if !cs.isServiceRunning() {
		return nil, fmt.Errorf("Cannot proceed because service is offline")
	}

	dbGorm, err := cs.Service.Db()
	if err != nil {
		return nil, err
	}

	var cart Cart
	res := dbGorm.Where("user_id = ?", userID).First(&cart)
	if res.Error != nil {
		if res.Error == gorm.ErrRecordNotFound {
			logging.Log.Error("Carrinho não encontrado", slog.String("error", res.Error.Error()))
			return nil, res.Error
		}

		logging.Log.Error("Error while searching for cart", slog.String("error", res.Error.Error()))
		return nil, res.Error
	}

	return &cart, nil
}

func (cs *CartService) UpdateByUser(userID uint, newCart *Cart) error {
	if !cs.isServiceRunning() {
		return fmt.Errorf("Cannot proceed because service is offline")
	}

	return models_utils.DoTransaction(cs.Service, models_utils.UPDATE, func(tx *gorm.DB) error {
		var cart Cart
		if err := tx.Where("user_id = ?", userID).First(&cart).Error; err != nil {
			return fmt.Errorf("cart does not exist")
		}

		if err := tx.Model(&cart).Updates(newCart).Error; err != nil {
			return fmt.Errorf("Failed to update user's cart")
		}

		return nil
	})
}

func (cs *CartService) DeleteByUser(userID uint) error {
	if !cs.isServiceRunning() {
		return fmt.Errorf("Cannot proceed because service is offline")

	}
	return models_utils.DoTransaction(cs.Service, models_utils.DELETE, func(tx *gorm.DB) error {
		if err := tx.Where("user_id = ?", userID).Delete(&Cart{}); err != nil {
			return fmt.Errorf("cart not found")
		}

		return nil
	})
}

func (cs *CartService) isServiceRunning() bool {
	if cs.Service == nil {
		logging.Log.Error("Cart Service is not initialized! Aborting")
	}

	return cs.Service != nil
}
