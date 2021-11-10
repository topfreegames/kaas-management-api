package k8s

import (
	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/rest"
	"log"
)

type Auth struct {
	AuthConfig    *rest.Config
	DynamicClient dynamic.Interface
}

func Authenticate() *Auth {
	config, err := rest.InClusterConfig()
	if err != nil {
		log.Fatalf("Could not retrieve pod service Account configuration: %v", err)
	}

	client, err := dynamic.NewForConfig(config)
	if err != nil {
		log.Fatalf("Could not authenticate to cluster as a pod: %v", err)
	}

	return &Auth{
		AuthConfig:    config,
		DynamicClient: client,
	}
}
