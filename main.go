package main

import (
	"github.com/topfreegames/kaas-management-api/internal/k8s"
	"github.com/topfreegames/kaas-management-api/internal/server"

	"log"
)

func main() {

	k8sClient := k8s.CreateK8sInstance()
	err := server.InitServer(k8sClient)
	if err != nil {
		log.Fatalf("Error initializing server: %v", err)
	}
}
