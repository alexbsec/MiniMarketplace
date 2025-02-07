package models_utils

import (
	"log/slog"

	"github.com/alexbsec/MiniMarketplace/src/db/config"
	"github.com/alexbsec/MiniMarketplace/src/logging"
	"gorm.io/gorm"
)

func DoTransaction(
    service *config.Service,
	event TransactionEvent,
	txFunc func(*gorm.DB) error) error {
	dbGorm, err := service.Db()
	if err != nil {
		logging.Log.Error(
			"Erro ao obter conexão com banco de dados",
			slog.String("error", err.Error()),
		)
		return err
	}

	// Execute transaction dynamically
	return dbGorm.Transaction(func(tx *gorm.DB) error {
		logging.Log.Info("Iniciando transação", slog.String("event", string(event)))
		if err := txFunc(tx); err != nil {
			logging.Log.Error("Erro durante a transação", slog.String("error", err.Error()))
			return err
		}

		logging.Log.Info("Transação concluída com sucesso", slog.String("event", string(event)))
		return nil
	})
}
