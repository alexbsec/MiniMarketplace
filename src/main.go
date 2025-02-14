package main

import (
	"github.com/alexbsec/MiniMarketplace/src/core"
	"github.com/alexbsec/MiniMarketplace/src/db/config"
)

func main() {
    config.Connect()
    app := app.App{}
    app.Run(":7676")
}
