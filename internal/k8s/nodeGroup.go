package k8s

func (k Kubernetes) GetNodeGroup(clusterName string, nodeGroupName string) {
	client := k.K8sAuth.DynamicClient

	resource := client.Resource(clusterResource)
}
