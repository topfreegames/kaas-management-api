package k8s

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/topfreegames/kaas-management-api/util"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime/schema"
	clusterapiv1beta1 "sigs.k8s.io/cluster-api/api/v1beta1"
)

// TODO check a better way to define this resource
var clusterResource = schema.GroupVersionResource{Group: "cluster.x-k8s.io", Version: "v1beta1", Resource: "clusters"}

// GetCluster gets a cluster-API cluster CR by name from the Kubernetes API. We follow the standard of one cluster per namespace.
func (k Kubernetes) GetCluster(clusterName string) (clusterapiv1beta1.Cluster, error) {
	client := k.K8sAuth.DynamicClient

	resource := client.Resource(clusterResource)

	clustersRaw, err := resource.Namespace(clusterName).Get(context.TODO(), clusterName, metav1.GetOptions{})
	if err != nil {
		if errors.IsNotFound(err) {
			return clusterapiv1beta1.Cluster{}, util.NewClientError(err, fmt.Sprintf("The requested cluster %s was not found in namespace %s!", clusterName, clusterName))
		} else if statusError, isStatus := err.(*errors.StatusError); isStatus {
			return clusterapiv1beta1.Cluster{}, fmt.Errorf("Error getting Cluster: %v\n", statusError.ErrStatus.Message)
		}
		return clusterapiv1beta1.Cluster{}, fmt.Errorf("Internal server error: %v\n", err)
	}

	var cluster clusterapiv1beta1.Cluster
	clustersRawsJson, err := clustersRaw.MarshalJSON()
	if err != nil {
		return clusterapiv1beta1.Cluster{}, fmt.Errorf("could not Marshal Clusters response: %v", err)
	}

	err = json.Unmarshal(clustersRawsJson, &cluster)
	if err != nil {
		return clusterapiv1beta1.Cluster{}, fmt.Errorf("could not Unmarshal Clusters JSON into clusterAPI list: %v", err)
	}

	return cluster, nil
}

// ListClusters list all cluster-api clusters CR in the kubernetes API and returns as cluster-api struct. We follow the standard of one cluster per namespace
func (k Kubernetes) ListClusters() (clusterapiv1beta1.ClusterList, error) {
	client := k.K8sAuth.DynamicClient

	clustersRaw, err := client.Resource(clusterResource).List(context.TODO(), metav1.ListOptions{})

	if errors.IsNotFound(err) {
		return clusterapiv1beta1.ClusterList{}, util.NewClientError(err, "Could not find any cluster in the cluster!")
	} else if statusError, isStatus := err.(*errors.StatusError); isStatus {
		return clusterapiv1beta1.ClusterList{}, fmt.Errorf("Error getting Cluster %v\n", statusError.ErrStatus.Message)
	} else if err != nil {
		return clusterapiv1beta1.ClusterList{}, fmt.Errorf("Internal server error\n")
	}

	var clusters clusterapiv1beta1.ClusterList

	clustersRawsJson, err := clustersRaw.MarshalJSON()
	if err != nil {
		return clusterapiv1beta1.ClusterList{}, fmt.Errorf("could not Marshal Clusters response: %v", err)
	}

	err = json.Unmarshal(clustersRawsJson, &clusters)
	if err != nil {
		return clusterapiv1beta1.ClusterList{}, fmt.Errorf("could not Unmarshal Clusters JSON into clusterAPI list: %v", err)
	}

	return clusters, nil
}
