package models

import (
	"fmt"
	"log/slog"

	"github.com/alexbsec/MiniMarketplace/src/db/config"
	"github.com/alexbsec/MiniMarketplace/src/logging"
	"gorm.io/gorm"
)

type TransactionEvent string

const (
	CREATE TransactionEvent = "CREATE"
	UPDATE TransactionEvent = "UPDATE"
	DELETE TransactionEvent = "DELETE"
)

type Product struct {
	ID          uint `gorm:"primaryKey"`
	Name        string
	Description string
	Price       float64
	Points      uint
	Category    string
}

type ProductService struct {
	service *config.Service
}

func (ps *ProductService) Create(product *Product) error {
    return ps.productTransaction(CREATE, func(tx *gorm.DB) error {
        if err := tx.Create(product).Error; err != nil {
            return fmt.Errorf("failed to create product: %w", err)
        }

        return nil
    })
}

func (ps *ProductService) Fetch(id uint) (*Product, error) {
	dbGorm, err := ps.service.Db()
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
    return ps.productTransaction(UPDATE, func(tx *gorm.DB) error {
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
    return ps.productTransaction(DELETE, func(tx *gorm.DB) error {
        if err := tx.Delete(&Product{}, id).Error; err != nil {
            return fmt.Errorf("failed to delete product: %w", err) 
        }

        return nil
    })
}

func (ps *ProductService) productTransaction(
    event TransactionEvent,
    txFunc func(*gorm.DB) error) error {
    dbGorm, err := ps.service.Db()
    if err != nil {
        logging.Log.Error(
            "Erro ao obter conexão com banco de dados",
            slog.String("error", err.Error()),
        )
        return err
    }

    // Execute transaction dynamically
    return dbGorm.Transaction(func (tx *gorm.DB) error {
        logging.Log.Info("Iniciando transação", slog.String("event", string(event)))
        if err := txFunc(tx); err != nil {
            logging.Log.Error("Erro durante a transação", slog.String("error", err.Error()))
            return err
        }

        logging.Log.Info("Transação concluída com sucesso", slog.String("event", string(event)))
        return nil
    })
}
