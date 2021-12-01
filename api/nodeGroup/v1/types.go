package v1

// NodeGroup - represents a Node Group
type NodeGroup struct {
	Name                   string    `json:"name"`
	Metadata               *Metadata `json:"metadata"`
	KubeProvider           string    `json:"kubeprovider,omitempty"`
	InfrastructureProvider string    `json:"infrastructureprovider"`
}

// NodeGroupList - a list of Node Groups
type NodeGroupList struct {
	Items []NodeGroup `json:"items"`
}

// TODO: We want this to be customizable in the future
type Metadata struct {
	Cluster     string   `json:"cluster"`
	Replicas    *int32   `json:"replicas"`
	MachineType string   `json:"machinetype"`
	Zones       []string `json:"zones"`
	Environment string   `json:"environment"`
	Region      string   `json:"region"`
	Min         *int32   `json:"min,omitempty"`
	Max         *int32   `json:"max,omitempty"`
}
