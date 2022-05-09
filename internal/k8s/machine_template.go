package k8s

import (
	"github.com/topfreegames/kaas-management-api/util/clientError"
	v1 "k8s.io/api/core/v1"
	clusterapiv1beta1 "sigs.k8s.io/cluster-api/api/v1beta1"
)

// ValidateMachineTemplateComponents validates if the machineTemplate, used by both machinePool and machineDeployment have all required fields
func ValidateMachineTemplateComponents(machineTemplate clusterapiv1beta1.MachineTemplateSpec) error {

	if machineTemplate.Spec.InfrastructureRef == (v1.ObjectReference{}) {
		return clientError.NewClientError(nil, clientError.InvalidConfiguration, "MachineTemplate doesn't have an infrastructure Reference")
	}

	if machineTemplate.Spec.InfrastructureRef.Name == "" {
		return clientError.NewClientError(nil, clientError.InvalidConfiguration, "MachineTemplate infrastructure reference name is empty")
	}

	if machineTemplate.Spec.InfrastructureRef.Kind == "" {
		return clientError.NewClientError(nil, clientError.InvalidConfiguration, "MachineTemplate infrastructure Kind is empty")
	}

	if machineTemplate.Spec.InfrastructureRef.APIVersion == "" {
		return clientError.NewClientError(nil, clientError.InvalidConfiguration, "MachineTemplate infrastructure APIVersion is empty")
	}
	return nil
}
