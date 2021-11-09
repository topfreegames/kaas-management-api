package main

import (
	"github.com/topfreegames/kaas-management-api/internal/server"

	"log"
)

func main() {
	err := server.InitServer()
	if err != nil {
		log.Fatalf("Error initializing server: %v", err)
	}
}
