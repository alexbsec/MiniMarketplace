package models

import (
	"fmt"
	"log/slog"

	"github.com/alexbsec/MiniMarketplace/src/db/config"
	"github.com/alexbsec/MiniMarketplace/src/db/models/utils"
	"github.com/alexbsec/MiniMarketplace/src/logging"
	"gorm.io/gorm"
)

type Wallet struct {
	ID     uint     `gorm:"primaryKey"`
	Name   *string  `gorm:"not null" json:"name"`
	Amount *float64 `gorm:"not null" json:"amount"`
	Points *float64 `gorm:"not null" json:"points"`
	UserID uint     `gorm:"not null" json:"user_id"`
	User   User     `gorm:"foreignKey:UserID"`
}

type WalletService struct {
	Service *config.Service
}

func (ws *WalletService) Create(wallet *Wallet) error {
	if !ws.isServiceRunning() {
		return fmt.Errorf("Cannot proceed because service is offline")
	}

	return models_utils.DoTransaction(ws.Service, models_utils.CREATE, func(tx *gorm.DB) error {
		if err := tx.Create(wallet).Error; err != nil {
			return fmt.Errorf("failed to create wallet: %w", err)
		}

		return nil
	})
}

func (ws *WalletService) Fetch(id uint) (*Wallet, error) {
	if !ws.isServiceRunning() {
		return nil, fmt.Errorf("Cannot proceed because service is offline")
	}

	dbGorm, err := ws.Service.Db()
	if err != nil {
		return nil, err
	}

	var wallet Wallet
	res := dbGorm.First(&wallet, id)
	if res.Error != nil {
		if res.Error == gorm.ErrRecordNotFound {
			logging.Log.Error("Carteira não encontrada", slog.String("error", res.Error.Error()))
			return nil, res.Error
		}

		logging.Log.Error("Error while searching for wallet", slog.String("error", res.Error.Error()))
		return nil, res.Error
	}

	return &wallet, nil
}

func (ws *WalletService) Update(id uint, newWallet *Wallet) error {
    if !ws.isServiceRunning() {
		return fmt.Errorf("Cannot proceed because service is offline")
    }

    return models_utils.DoTransaction(ws.Service, models_utils.UPDATE, func(tx *gorm.DB) error {
        var wallet Wallet
        if err := tx.First(&wallet, id).Error; err != nil {
            return fmt.Errorf("wallet with id %d not found: %w", id, err)
        }

        if err := tx.Model(&wallet).Updates(newWallet).Error; err != nil {
            return fmt.Errorf("failed to update wallet with id %d: %w", id, err)
        }

        return nil
    })
}

func (ws *WalletService) Delete(id uint) error {
    if !ws.isServiceRunning() {
		return fmt.Errorf("Cannot proceed because service is offline")
    }

    return models_utils.DoTransaction(ws.Service, models_utils.DELETE, func(tx *gorm.DB) error {
        if err := tx.Delete(&Wallet{}, id).Error; err != nil {
            return fmt.Errorf("failed to delete wallet: %w", err)
        }

        return nil
    })
}

func (ws *WalletService) isServiceRunning() bool {
	if ws.Service == nil {
		logging.Log.Error("Wallet Service is not initialized! Aborting")
	}

	return ws.Service != nil
}

