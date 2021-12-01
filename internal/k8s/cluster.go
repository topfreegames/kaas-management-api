package k8s

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/topfreegames/kaas-management-api/util/clientError"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	clusterapiv1beta1 "sigs.k8s.io/cluster-api/api/v1beta1"
)

// GetCluster gets a cluster-API cluster CR by name from the Kubernetes API. We follow the standard of one cluster per namespace.
func (k Kubernetes) GetCluster(clusterName string) (*clusterapiv1beta1.Cluster, error) {
	client := k.K8sAuth.DynamicClient

	resource := client.Resource(ClusterResourceSchemaV1beta1)

	clustersRaw, err := resource.Namespace(clusterName).Get(context.TODO(), clusterName, metav1.GetOptions{})
	if err != nil {
		if errors.IsNotFound(err) {
			return nil, clientError.NewClientError(err, clientError.ResourceNotFound, fmt.Sprintf("The requested cluster %s was not found in namespace %s!", clusterName, clusterName))
		} else if statusError, isStatus := err.(*errors.StatusError); isStatus {
			return nil, fmt.Errorf("Error getting Cluster: %v\n", statusError.ErrStatus.Message)
		}
		return nil, fmt.Errorf("Internal server clientError: %v\n", err)
	}

	var cluster clusterapiv1beta1.Cluster
	clustersRawsJson, err := clustersRaw.MarshalJSON()
	if err != nil {
		return nil, fmt.Errorf("could not Marshal Clusters response: %v", err)
	}

	err = json.Unmarshal(clustersRawsJson, &cluster)
	if err != nil {
		return nil, fmt.Errorf("could not Unmarshal Clusters JSON into clusterAPI list: %v", err)
	}

	return &cluster, nil
}

// ListClusters list all cluster-api clusters CR in the kubernetes API and returns as cluster-api struct. We follow the standard of one cluster per namespace
func (k Kubernetes) ListClusters() (*clusterapiv1beta1.ClusterList, error) {
	client := k.K8sAuth.DynamicClient

	clustersRaw, err := client.Resource(ClusterResourceSchemaV1beta1).List(context.TODO(), metav1.ListOptions{})

	if errors.IsNotFound(err) {
		return nil, clientError.NewClientError(err, clientError.ResourceNotFound, "Could not find any cluster in the Kubernetes API!")
	} else if statusError, isStatus := err.(*errors.StatusError); isStatus {
		return nil, fmt.Errorf("Error getting Cluster %v\n", statusError.ErrStatus.Message)
	} else if err != nil {
		return nil, fmt.Errorf("Internal server clientError\n")
	}

	var clusters clusterapiv1beta1.ClusterList

	clustersRawsJson, err := clustersRaw.MarshalJSON()
	if err != nil {
		return nil, fmt.Errorf("could not Marshal Clusters response: %v", err)
	}

	err = json.Unmarshal(clustersRawsJson, &clusters)
	if err != nil {
		return nil, fmt.Errorf("could not Unmarshal Clusters JSON into clusterAPI list: %v", err)
	}

	if len(clusters.Items) == 0 {
		return nil, clientError.NewClientError(err, clientError.EmptyResponse, "no Clusters were found!")
	}

	return &clusters, nil
}
