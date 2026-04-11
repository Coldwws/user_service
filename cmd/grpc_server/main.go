package main

import (
	"context"
	"flag"
	"log"
	"user_service/internal/app"
)

var logLevel string

func init() {
	flag.StringVar(&logLevel, "l", "info", "log level")
}

func main() {
	ctx := context.Background()

	a, err := app.NewApp(ctx, logLevel)
	if err != nil {
		log.Fatalf("Failed to init app: %v", err)
	}

	err = a.Run()
	if err != nil {
		log.Fatalf("Failed to run app: %v", err)
	}
}
