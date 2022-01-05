package k8s

type Kubernetes struct {
	K8sAuth *Auth
}

func CreateK8sInstance() *Kubernetes {
	auth := Authenticate()
	return &Kubernetes{K8sAuth: auth}
}
