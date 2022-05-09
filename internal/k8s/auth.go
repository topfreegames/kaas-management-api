package k8s

import (
	"fmt"
	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	"log"
	"os"
)

type Auth struct {
	AuthConfig    *rest.Config
	DynamicClient dynamic.Interface
}

func Authenticate() *Auth {
	config, err := rest.InClusterConfig()
	if err != nil {
		log.Printf("Could not retrieve pod service Account configuration: %v", err)
		log.Print("Trying local authentication")
		return LocalAuthenticate()
	}

	client, err := dynamic.NewForConfig(config)
	if err != nil {
		log.Fatalf("Could not create client as a pod: %v", err)
	}

	log.Print("Using local authentication")
	return &Auth{
		AuthConfig:    config,
		DynamicClient: client,
	}
}

func LocalAuthenticate() *Auth {
	var kubeConfigPath string
	kubeConfigPath = os.Getenv("KUBECONFIG")
	if kubeConfigPath == "" {
		kubeConfigPath = fmt.Sprintf("%s/.kube/config", os.Getenv("HOME"))
		log.Printf("KUBECONFIG env is unset, using default location in $HOME: %s", kubeConfigPath)
	}

	config, err := clientcmd.BuildConfigFromFlags("", kubeConfigPath)
	if err != nil {
		log.Fatalf("Could not retrieve Kubeconfig configuration: %v", err)
	}

	client, err := dynamic.NewForConfig(config)
	if err != nil {
		log.Fatalf("Could not create client using Kubeconfig: %v", err)
	}

	return &Auth{
		AuthConfig:    config,
		DynamicClient: client,
	}
}
