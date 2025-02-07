package config

import (
	"fmt"
	"log/slog"
	"os"

	"github.com/alexbsec/MiniMarketplace/src/logging"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type Service struct {
	db *gorm.DB
}

func (s *Service) Db() (*gorm.DB, error) {
    if s.db == nil {
        return nil, fmt.Errorf("database connection is not initialized")
    }

    return s.db, nil
}

func InitMockService(mockDB *gorm.DB) *Service {
	return &Service{db: mockDB}
}

// InitService initialize the service
func InitService() (*Service, error) {
    host := os.Getenv("DB_HOST")
    user := os.Getenv("DB_USER")
    pass := os.Getenv("DB_PASSWORD")
    dbname := os.Getenv("DB_NAME")
    port := 5432

	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%d sslmode=disable TimeZone=UTC",
        host, user, pass, dbname, port)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{PrepareStmt: false})
	if err != nil {
		logging.Log.Error(
			"Failed to initialize service",
			slog.String("error", err.Error()),
		)
		return nil, err
	}

	return &Service{db: db}, nil
}
