package k8s

import "k8s.io/apimachinery/pkg/runtime/schema"

var (
	ClusterResourceSchemaV1beta1   = schema.GroupVersionResource{Group: "cluster.x-k8s.io", Version: "v1beta1", Resource: "clusters"}
	MachinePoolSchemaV1beta1       = schema.GroupVersionResource{Group: "cluster.x-k8s.io", Version: "v1beta1", Resource: "machinepools"}
	MachineDeploymentSchemaV1beta1 = schema.GroupVersionResource{Group: "cluster.x-k8s.io", Version: "v1beta1", Resource: "machinedeployments"}

	DockerMachineTemplateSchemaV1beta1 = schema.GroupVersionResource{Group: "infrastructure.cluster.x-k8s.io", Version: "v1beta1", Resource: "dockermachinetemplates"}
	KopsMachinePoolSchemaV1alpha1      = schema.GroupVersionResource{Group: "infrastructure.cluster.x-k8s.io", Version: "v1alpha1", Resource: "kopsmachinepools"}
)
