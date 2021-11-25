package k8s

import "k8s.io/apimachinery/pkg/runtime/schema"

var (
	clusterResourceSchemaV1beta1   = schema.GroupVersionResource{Group: "cluster.x-k8s.io", Version: "v1beta1", Resource: "clusters"}
	machinePoolSchemaV1beta1       = schema.GroupVersionResource{Group: "cluster.x-k8s.io", Version: "v1beta1", Resource: "machinepools"}
	machineDeploymentSchemaV1beta1 = schema.GroupVersionResource{Group: "cluster.x-k8s.io", Version: "v1beta1", Resource: "machinedeployments"}

	dockerMachineTemplateSchemaV1beta1 = schema.GroupVersionResource{Group: "infrastructure.cluster.x-k8s.io", Version: "v1beta1", Resource: "dockermachinetemplates"}
	kopsMachinePoolSchemaV1alpha1      = schema.GroupVersionResource{Group: "infrastructure.cluster.x-k8s.io", Version: "v1alpha1", Resource: "kopsmachinepools"}
)
