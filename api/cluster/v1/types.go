package v1

// Cluster - represents a cluster
type Cluster struct {
	Name                   string                 `json:"name"`
	Metadata               map[string]interface{} `json:"metadata"`
	KubeProvider           string                 `json:"kubeprovider"`
	InfrastructureProvider string                 `json:"infrastructureprovider"`
}

// ClusterList - a list of Cluster
type ClusterList struct {
	Items []Cluster `json:"items"`
}
