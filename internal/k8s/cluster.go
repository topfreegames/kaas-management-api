package k8s

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/topfreegames/kaas-management-api/util/clientError"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"log"
	clusterapiv1beta1 "sigs.k8s.io/cluster-api/api/v1beta1"
	"strings"
)

func GetClusterNamespace(clusterName string) string {

	prefix := "kubernetes"
	clusterNamespace := strings.ReplaceAll(clusterName, ".", "-")
	namespace := fmt.Sprintf("%s-%s", prefix, clusterNamespace)
	return namespace
}

// GetCluster gets a cluster-API cluster CR by name from the Kubernetes API. We follow the standard of one cluster per namespace.
func (k Kubernetes) GetCluster(clusterName string) (*clusterapiv1beta1.Cluster, error) {
	client := k.K8sAuth.DynamicClient

	resource := client.Resource(ClusterResourceSchemaV1beta1)

	namespace := GetClusterNamespace(clusterName)
	clustersRaw, err := resource.Namespace(namespace).Get(context.TODO(), clusterName, metav1.GetOptions{})
	if err != nil {
		if errors.IsNotFound(err) {
			return nil, clientError.NewClientError(err, clientError.ResourceNotFound, fmt.Sprintf("The requested cluster %s was not found in namespace %s!", clusterName, namespace))
		} else if statusError, isStatus := err.(*errors.StatusError); isStatus {
			return nil, fmt.Errorf("Error getting Cluster from Kubernetes API: %v\n", statusError.ErrStatus.Message)
		}
		return nil, fmt.Errorf("Kube go-client Error: %v\n", err)
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

	err = ValidateClusterComponents(&cluster)
	if err != nil {
		return nil, clientError.NewClientError(err, clientError.InvalidConfiguration, fmt.Sprintf("Cluster %s have an invalid configuration", clusterName))
	}

	return &cluster, nil
}

// ListClusters list all cluster-api clusters CR in the kubernetes API and returns as cluster-api struct. We follow the standard of one cluster per namespace
func (k Kubernetes) ListClusters() (*clusterapiv1beta1.ClusterList, error) {
	client := k.K8sAuth.DynamicClient

	clustersRaw, err := client.Resource(ClusterResourceSchemaV1beta1).List(context.TODO(), metav1.ListOptions{})

	if errors.IsNotFound(err) {
		return nil, clientError.NewClientError(err, clientError.ResourceNotFound, "Could not find any cluster in the Kubernetes API")
	} else if statusError, isStatus := err.(*errors.StatusError); isStatus {
		return nil, fmt.Errorf("Error getting Cluster from Server API %v\n", statusError.ErrStatus.Message)
	} else if err != nil {
		return nil, fmt.Errorf("Kube go-client Error: %v\n", err)
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
		return nil, clientError.NewClientError(err, clientError.EmptyResponse, "no Clusters were found")
	}

	clustersValidated := clusterapiv1beta1.ClusterList{
		TypeMeta: clusters.TypeMeta,
		ListMeta: clusters.ListMeta,
		Items:    []clusterapiv1beta1.Cluster{},
	}

	for _, cluster := range clusters.Items {
		err := ValidateClusterComponents(&cluster)
		if err != nil {
			log.Printf("Skiping cluster %s because of invalid configuration: %v", cluster.Name, err.Error())
			continue
		}
		clustersValidated.Items = append(clustersValidated.Items, cluster)
	}

	if len(clustersValidated.Items) == 0 {
		return nil, clientError.NewClientError(nil, clientError.EmptyResponse, "no valid clusters were found, some clusters have invalid configuration")
	}

	return &clustersValidated, nil
}

func ValidateClusterComponents(cluster *clusterapiv1beta1.Cluster) error {
	if cluster.Spec.InfrastructureRef == nil {
		return clientError.NewClientError(nil, clientError.InvalidConfiguration, "Cluster doesn't have a infrastructure Reference")
	}

	if cluster.Spec.ControlPlaneRef == nil {
		return clientError.NewClientError(nil, clientError.InvalidConfiguration, "Cluster doesn't have a ControlPlane Reference")
	}

	if !cluster.Spec.ControlPlaneEndpoint.IsValid() {
		return clientError.NewClientError(nil, clientError.InvalidConfiguration, "Cluster doesn't have a valid ControlPlane endpoint")
	}
	return nil
}
