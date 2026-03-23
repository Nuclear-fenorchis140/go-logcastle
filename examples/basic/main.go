package main

import (
	"fmt"
	"log"
	"time"

	"github.com/bhaskarblur/go-logcastle"
)

func main() {
	// Initialize logcastle with JSON output
	// This is the ONLY setup needed!
	err := logcastle.Init(logcastle.Config{
		Format: logcastle.JSON,
		Level:  logcastle.LevelInfo,
	})
	if err != nil {
		log.Fatal("Failed to initialize logcastle:", err)
	}
	defer logcastle.Close()

	// Wait for initialization to complete
	logcastle.WaitReady()

	fmt.Println("=== go-logcastle Basic Example ===")
	fmt.Println("")

	// From here on, ALL logs are intercepted and standardized!

	// Standard library log
	log.Println("This is a stdlib log message")

	// fmt.Println also gets captured
	fmt.Println("This is a fmt.Println message")

	// Direct println
	println("This is a plain println")

	// More structured logging
	log.Printf("User %s logged in from %s", "alice", "192.168.1.1")

	// Give time for logs to flush
	time.Sleep(200 * time.Millisecond)

	fmt.Println("")
	fmt.Println("=== All logs above were standardized to JSON! ===")
}
