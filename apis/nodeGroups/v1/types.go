package v1

// NodeGroup - represents a Node Group
type NodeGroup struct {
	Name                   string                 `json:"name"`
	Metadata               map[string]interface{} `json:"metadata"`
	KubeProvider           string                 `json:"kubeprovider"`
	InfrastructureProvider string                 `json:"infrastructureprovider"`
}

// NodeGroupList - a list of Node Groups
type NodeGroupList struct {
	Items []NodeGroup `json:"items"`
}
