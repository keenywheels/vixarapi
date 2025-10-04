package main

import (
	"log"

	"github.com/keenywheels/backend/internal/app"
)

func main() {
	if err := app.New().Run(); err != nil {
		log.Fatal(err)
	}
}
