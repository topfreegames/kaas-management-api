package k8s

import (
	"context"
	"encoding/json"
	"fmt"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime/schema"
	clusterapiv1beta1 "sigs.k8s.io/cluster-api/api/v1beta1"
)

var clusterResource = schema.GroupVersionResource{Group: "cluster.x-k8s.io", Version: "v1beta1", Resource: "clusters"}

func (k Kubernetes) GetCluster(clusterName string, namespace string) (clusterapiv1beta1.Cluster, error) {
	client := k.K8sAuth.DynamicClient

	resource := client.Resource(clusterResource)

	clustersRaw, err := resource.Namespace(namespace).Get(context.TODO(), clusterName, metav1.GetOptions{})
	if err != nil {
		if errors.IsNotFound(err) {
			return clusterapiv1beta1.Cluster{}, fmt.Errorf("Cluster %s not found: %v\n", clusterName, err)
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

func (k Kubernetes) ListClusters(namespace string) (clusterapiv1beta1.ClusterList, error) {
	client := k.K8sAuth.DynamicClient

	clustersRaw, err := client.Resource(clusterResource).Namespace(namespace).List(context.TODO(), metav1.ListOptions{})

	if errors.IsNotFound(err) {
		return clusterapiv1beta1.ClusterList{}, fmt.Errorf("Clusters not found in namespace %s\n", namespace)
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
