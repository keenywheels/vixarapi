package main

import (
	"log"

	app "github.com/keenywheels/backend/internal/vixarapi"
)

func main() {
	if err := app.New().Run(); err != nil {
		log.Fatal(err)
	}
}
