package config

import (
	"github.com/alexbsec/MiniMarketplace/src/logging"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"log/slog"
)

type Service struct {
	db *gorm.DB
}

// InitService initialize the service
func InitService() (*Service, error) {
	opts := &slog.HandlerOptions{
		Level:     slog.LevelDebug,
		AddSource: true,
	}
	log := slog.New(logging.NewHandler(opts))

	dsn := "host=db user=user password=password dbname=marketplace_db port=5432 sslmode=disable TimeZone=UTC"

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Error(
			"Failed to initialize service",
			slog.String("error", err.Error()),
		)
		return nil, err
	}

	return &Service{db: db}, nil
}
