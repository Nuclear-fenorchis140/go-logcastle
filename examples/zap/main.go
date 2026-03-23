package main

import (
	"time"

	"github.com/yourusername/go-logcastle"
	"go.uber.org/zap"
)

func main() {
	// ONE LINE SETUP
	logcastle.Init(logcastle.Config{
		Format: logcastle.JSON,
	})
	defer logcastle.Close()

	// Use zap as normal
	logger, _ := zap.NewProduction()
	defer logger.Sync()

	logger.Info("Application started",
		zap.String("version", "1.0.0"),
		zap.Int("port", 8080),
	)

	logger.Info("Processing request",
		zap.String("user", "alice"),
		zap.Duration("latency", 23*time.Millisecond),
	)

	logger.Error("Failed to connect to database",
		zap.String("host", "localhost"),
		zap.Int("port", 5432),
	)

	time.Sleep(200 * time.Millisecond)
}
