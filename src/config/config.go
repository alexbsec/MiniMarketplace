package config

import (
	"encore.dev/storage/sqldb"
	"github.com/alexbsec/MiniMarketplace/src/logging"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
    "log/slog"
)

type Service struct {
    db *gorm.DB
}

var mktPlaceDB = sqldb.NewDatabase("marketplace_db", sqldb.DatabaseConfig{
    Migrations: "src/migrations",
}) 

// initService initialize the service
// It is automatically called by Encore on service startup
func initService() (*Service, error) {
    opts := &slog.HandlerOptions{
        Level: slog.LevelDebug,
        AddSource: true,
    }
    log := slog.New(logging.NewHandler(opts))
    db, err := gorm.Open(postgres.New(postgres.Config{
        Conn: mktPlaceDB.Stdlib(),
    }))
    if err != nil {
        log.Error(
            "Failed to initialize service",
            slog.String("error", err.Error()),
        )
        return nil, err 
    }

    return &Service{db: db}, nil
}
