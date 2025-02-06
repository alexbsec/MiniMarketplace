package config

import (
	"fmt"
	"log/slog"
	"os"

	"github.com/alexbsec/MiniMarketplace/src/logging"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var (
    opts *slog.HandlerOptions = &slog.HandlerOptions{
        Level:     slog.LevelDebug,
	    AddSource: true,
    }

    log *slog.Logger = slog.New(logging.NewHandler(opts))
)

type Service struct {
	db *gorm.DB
}

func (s *Service) Db() (*gorm.DB, error) {
    service, err := InitService()
    if err != nil {
        return nil, err
    }

    s.db = service.db
    return s.db, nil
}

func InitMockService(mockDB *gorm.DB) *Service {
	return &Service{db: mockDB}
}

// InitService initialize the service
func InitService() (*Service, error) {
	log := slog.New(logging.NewHandler(opts))
    host := os.Getenv("DB_HOST")
    user := os.Getenv("DB_USER")
    pass := os.Getenv("DB_PASSWORD")
    dbname := os.Getenv("DB_NAME")
    port := 5432

	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%d sslmode=disable TimeZone=UTC",
        host, user, pass, dbname, port)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{PrepareStmt: false})
	if err != nil {
		log.Error(
			"Failed to initialize service",
			slog.String("error", err.Error()),
		)
		return nil, err
	}

	return &Service{db: db}, nil
}
