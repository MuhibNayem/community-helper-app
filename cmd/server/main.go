package main

import (
	"log"

	"github.com/MuhibNayem/community-helper-app/internal/app"
)

func main() {
	application, err := app.New()
	if err != nil {
		log.Fatalf("initialize app: %v", err)
	}

	if err := application.Run(); err != nil {
		log.Fatalf("run app: %v", err)
	}
}
