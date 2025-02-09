package main

import (
	"io"
	"log/slog"
	"os"

	_ "ariga.io/atlas-go-sdk/recordriver"
	"ariga.io/atlas-provider-gorm/gormschema"
	"github.com/alexbsec/MiniMarketplace/src/logging"
    "github.com/alexbsec/MiniMarketplace/src/db/models"
)

func main() {
    loader := gormschema.New("postgres") 

    stmts, err := loader.Load(
        &models.Product{},
        &models.User{},
        &models.Wallet{},
        )

    if err != nil {
        logging.Log.Error("Failed to load GORM schema", slog.String("error", err.Error()))
        return
    }

    io.WriteString(os.Stdout, stmts)
}
