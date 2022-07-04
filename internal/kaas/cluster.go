package kaas

import (
	"fmt"
	"github.com/topfreegames/kaas-management-api/internal/k8s"
	"github.com/topfreegames/kaas-management-api/util/clientError"
	"log"
	"sigs.k8s.io/cluster-api/api/v1beta1"
	clusterapiv1beta1 "sigs.k8s.io/cluster-api/api/v1beta1"
)

type Cluster struct {
	Name                     string
	ApiEndpoint              string
	ControlPlaneEndpointHost string
	ControlPlaneEndpointPort int32
	Region                   string
	ClusterGroup             string
	Environment              string
	CIDR                     []string
	ControlPlane             *ClusterControlPlane
	Infrastructure           *ClusterInfrastructure
}

func GetCluster(k *k8s.Kubernetes, name string) (*Cluster, error) {

	clusterAPICR, err := k.GetCluster(name)
	if err != nil {
		clientErr, ok := err.(*clientError.ClientError)
		if !ok {
			return nil, clientError.NewClientError(clientErr, clientError.UnexpectedError, fmt.Sprintf("Something went wrong while getting cluster %s", name))
		} else {
			if clientErr.ErrorMessage == clientError.ResourceNotFound {
				return nil, clientError.NewClientError(clientErr, clientError.ResourceNotFound, fmt.Sprintf("Could not find cluster %s", name))
			} else {
				return nil, clientError.NewClientError(clientErr, clientError.UnexpectedError, fmt.Sprintf("Error getting cluster %s", name))
			}
		}
	}

	err = ValidateClusterComponents(clusterAPICR)
	if err != nil {
		return nil, clientError.NewClientError(err, clientError.InvalidConfiguration, fmt.Sprintf("Cluster %s have an invalid configuration", name))
	}

	cluster := &Cluster{}
	err = cluster.GetClusterProperties(clusterAPICR)
	if err != nil {
		clientErr, ok := err.(*clientError.ClientError)
		if !ok {
			return nil, clientError.NewClientError(err, clientError.UnexpectedError, fmt.Sprintf("An Unexpected error happened while reading cluster %s properties", name))
		} else {
			if clientErr.ErrorMessage == clientError.InvalidConfiguration {
				return nil, clientError.NewClientError(clientErr, clientError.InvalidConfiguration, fmt.Sprintf("Cluster %s is invalid due to missing or invalid labels", name))
			}
		}
	}

	return cluster, nil
}

func ListClusters(k *k8s.Kubernetes) ([]*Cluster, error) {

	var clusterList []*Cluster

	clusterListAPICR, err := k.ListClusters()
	if err != nil {
		clientErr, ok := err.(*clientError.ClientError)
		if !ok {
			return nil, clientError.NewClientError(clientErr, clientError.UnexpectedError, "Error listing clusters")
		} else {
			if clientErr.ErrorMessage == clientError.ResourceNotFound {
				return nil, clientErr
			} else if clientErr.ErrorMessage == clientError.EmptyResponse {
				return nil, clientErr
			} else {
				return nil, clientError.NewClientError(clientErr, clientError.UnexpectedError, "Something went wrong when listing clusters")
			}
		}
	}

	for _, clusterAPICR := range clusterListAPICR.Items {
		cluster := &Cluster{}
		err = ValidateClusterComponents(&clusterAPICR)
		if err != nil {
			log.Printf("Skipping cluster %s because of invalid configuration: %s", cluster.Name, err.Error())
			continue
		}
		err = cluster.GetClusterProperties(&clusterAPICR)
		if err != nil {
			clientErr, ok := err.(*clientError.ClientError)
			if !ok {
				log.Printf("Skipping cluster %s: An Unexpected error happened while reading the cluster properties: %s", clusterAPICR.Name, err.Error())
			} else {
				if clientErr.ErrorMessage == clientError.InvalidConfiguration {
					log.Printf("Skipping cluster %s: Cluster is invalid due to missing or invalid labels: %s", clusterAPICR.Name, err.Error())
				}
			}
			continue
		}
		clusterList = append(clusterList, cluster)
	}

	if len(clusterList) == 0 {
		return nil, clientError.NewClientError(nil, clientError.EmptyResponse, "No valid clusters were found, some clusters have invalid configuration")
	}

	return clusterList, nil
}

// TODO do the validation on each Get method from each component
func ValidateClusterComponents(cluster *clusterapiv1beta1.Cluster) error {
	if cluster.Spec.InfrastructureRef == nil {
		return clientError.NewClientError(nil, clientError.InvalidConfiguration, "Cluster doesn't have an infrastructure Reference")
	}

	if cluster.Spec.ControlPlaneRef == nil {
		return clientError.NewClientError(nil, clientError.InvalidConfiguration, "Cluster doesn't have a ControlPlane Reference")
	}

	if !cluster.Spec.ControlPlaneEndpoint.IsValid() {
		return clientError.NewClientError(nil, clientError.InvalidConfiguration, "Cluster doesn't have a valid ControlPlane endpoint")
	}
	return nil
}

func (c *Cluster) GetClusterProperties(clusterAPICR *v1beta1.Cluster) error {
	c.Name = clusterAPICR.Name
	c.ControlPlaneEndpointHost = clusterAPICR.Spec.ControlPlaneEndpoint.Host
	c.ControlPlaneEndpointPort = clusterAPICR.Spec.ControlPlaneEndpoint.Port
	c.ApiEndpoint = fmt.Sprintf("https://%s:%d", c.ControlPlaneEndpointHost, c.ControlPlaneEndpointPort)
	// TODO: turn customizable
	c.ClusterGroup = clusterAPICR.Labels["clusterGroup"]
	c.Region = clusterAPICR.Labels["region"]
	c.Environment = clusterAPICR.Labels["environment"]
	c.CIDR = clusterAPICR.Spec.ClusterNetwork.Services.CIDRBlocks

	cp, err := GetControlPlane(clusterAPICR.Spec.ControlPlaneRef.Kind)
	if err != nil {
		clientErr, ok := err.(*clientError.ClientError)
		if !ok {
			return clientError.NewClientError(err, clientError.UnexpectedError, fmt.Sprintf("An Unexpected error heppened while reading cluster control plane resource for cluster %s", c.Name))
		} else {
			if clientErr.ErrorMessage == clientError.KindNotFound {
				return clientError.NewClientError(clientErr, clientError.InvalidConfiguration, fmt.Sprintf("Could not get cluster %s controlplane property", c.Name))
			} else {
				return clientError.NewClientError(err, clientError.UnexpectedError, fmt.Sprintf("An Unexpected error heppened while reading cluster control plane resource for cluster %s", c.Name))
			}
		}
	}
	c.ControlPlane = cp

	c.Infrastructure, err = GetClusterInfrastructure(clusterAPICR.Spec.InfrastructureRef.Kind)
	if err != nil {
		clientErr, ok := err.(*clientError.ClientError)
		if !ok {
			return clientError.NewClientError(err, clientError.UnexpectedError, fmt.Sprintf("An Unexpected error heppened while reading cluster Infrastrucutre resource for cluster %s", c.Name))
		} else {
			if clientErr.ErrorMessage == clientError.KindNotFound {
				return clientError.NewClientError(clientErr, clientError.InvalidConfiguration, fmt.Sprintf("Could not get cluster %s infrastructure property", c.Name))
			} else {
				return clientError.NewClientError(err, clientError.UnexpectedError, fmt.Sprintf("An Unexpected error heppened while reading cluster Infrastrucutre resource for cluster %s", c.Name))
			}
		}
	}
	return nil
}
