package main

import (
	"errors"
	"fmt"
	"log/slog"

	"github.com/alexbsec/MiniMarketplace/src/outlogger"
)

func main() {
	logger := slog.New(
		outlogger.NewHandler(nil).WithAttrs(
			[]slog.Attr{
				slog.Group("status",
					slog.String("key1", "test"),
					slog.String("key2", "test2"),
				),
			},
		),
	)

	// Log an error with additional context
	logger.Error("Ops, test failed successfully!",
		slog.String("error", errors.New("bad luck").Error()),
	)
	fmt.Println("Servidor iniciado!")
}
