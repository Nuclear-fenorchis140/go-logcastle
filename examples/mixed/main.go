package main

import (
	"log"
	"time"

	"github.com/sirupsen/logrus"
	"github.com/yourusername/go-logcastle"
	"go.uber.org/zap"
)

func main() {
	// ONE LINE SETUP - works for ALL loggers!
	logcastle.Init(logcastle.Config{
		Format: logcastle.JSON,
		EnrichFields: map[string]interface{}{
			"service": "my-app",
			"env":     "production",
		},
	})
	defer logcastle.Close()

	// Mix different logging libraries - all get standardized!

	// 1. Standard library
	log.Println("Starting application")

	// 2. Logrus
	logrusLogger := logrus.New()
	logrusLogger.WithField("component", "auth").Info("Authentication initialized")

	// 3. Zap
	zapLogger, _ := zap.NewProduction()
	zapLogger.Info("Database connection established",
		zap.String("db", "postgres"),
	)
	zapLogger.Sync()

	// 4. More stdlib
	log.Printf("Listening on port %d", 8080)

	time.Sleep(200 * time.Millisecond)

	// All logs above are now in the SAME standardized format!
}
