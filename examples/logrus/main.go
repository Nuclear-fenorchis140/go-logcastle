package main

import (
	"time"

	"github.com/sirupsen/logrus"
	"github.com/yourusername/go-logcastle"
)

func main() {
	// ONE LINE SETUP - that's it!
	logcastle.Init(logcastle.Config{
		Format: logcastle.JSON,
	})
	defer logcastle.Close()

	// Use logrus as normal - logs will be intercepted and standardized
	log := logrus.New()
	log.SetFormatter(&logrus.JSONFormatter{})

	log.WithFields(logrus.Fields{
		"user_id": "123",
		"action":  "login",
	}).Info("User logged in")

	log.WithFields(logrus.Fields{
		"path":   "/api/users",
		"method": "GET",
		"status": 200,
	}).Info("API request")

	log.Error("Something went wrong")

	time.Sleep(200 * time.Millisecond)
}
