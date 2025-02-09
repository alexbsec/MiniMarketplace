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
	Amount *float64 `gorm:"not null" json:"amount"`
	Points *float64 `gorm:"not null" json:"points"`
	User   User     `gorm:"foreignKey:UserRefer" json:"user"`
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

func (ws *WalletService) isServiceRunning() bool {
	if ws.Service == nil {
		logging.Log.Error("Wallet Service is not initialized! Aborting")
	}

	return ws.Service != nil
}
