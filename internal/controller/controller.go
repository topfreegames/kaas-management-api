package controller

import "github.com/topfreegames/kaas-management-api/internal/k8s"

type ControllerConfig struct {
	K8sInstance *k8s.Kubernetes
}

func ConfigureControllers(k8sInstance *k8s.Kubernetes) ControllerConfig {
	return ControllerConfig{K8sInstance: k8sInstance}
}
