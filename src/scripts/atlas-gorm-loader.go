package main

import (
	"io"
	"log/slog"
	"os"

	_ "ariga.io/atlas-go-sdk/recordriver"
	"ariga.io/atlas-provider-gorm/gormschema"
	"github.com/alexbsec/MiniMarketplace/src/logging"
    "github.com/alexbsec/MiniMarketplace/src/models"
)

func main() {
    loader := gormschema.New("postgres") 
    log := slog.New(logging.NewHandler(nil))

    stmts, err := loader.Load(&models.Product{})
    if err != nil {
        log.Error("Failed to load GORM schema", slog.String("error", err.Error()))
        return
    }

    io.WriteString(os.Stdout, stmts)
}
