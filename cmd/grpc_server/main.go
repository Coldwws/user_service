package main

import (
	"context"
	"log"
	"user_service/internal/app"
)

func main() {
	ctx := context.Background()

	a,err := app.NewApp(ctx)
	if err != nil{
		log.Fatalf("Failed to init app: %v",err)
	}

	err = a.Run()
	if err != nil{
		log.Fatalf("Failed to run app: %v",err)
	}
}
