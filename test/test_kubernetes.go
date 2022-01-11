package test

import (
	"fmt"
	clusterapikopsv1alpha1 "github.com/topfreegames/kubernetes-kops-operator/apis/infrastructure/v1alpha1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/dynamic/fake"
	"k8s.io/kops/pkg/apis/kops/v1alpha2"
	"log"
	clusterapiv1beta1 "sigs.k8s.io/cluster-api/api/v1beta1"
	clusterapiexpv1beta1 "sigs.k8s.io/cluster-api/exp/api/v1beta1"
	"strings"
)

type K8sRequest struct {
	ResourceName string
	ResourceKind string
	Cluster      string
}

// GetK8sRequest returns the request of the test as an instance of the struct *K8sRequest
func (t TestCase) GetK8sRequest() *K8sRequest {
	request, ok := t.Request.(*K8sRequest)
	if !ok {
		log.Fatalf("Could not convert TestCase %s Request to k8sRequest", t.Name)
	}
	return request
}

func NewK8sFakeDynamicClient() *fake.FakeDynamicClient {
	client := fake.NewSimpleDynamicClient(runtime.NewScheme())
	return client
}

func NewK8sFakeDynamicClientWithResources(resources ...runtime.Object) *fake.FakeDynamicClient {

	client := fake.NewSimpleDynamicClient(runtime.NewScheme(), resources...)

	return client
}

func GetTestClusterNamespace(clusterName string) string {
	prefix := "kubernetes"
	clusterNamespace := strings.ReplaceAll(clusterName, ".", "-")
	namespace := fmt.Sprintf("%s-%s", prefix, clusterNamespace)
	return namespace
}

func NewTestCluster(name string, controlPlaneName string, controPlaneKind string, controlPlaneApiVersion string, infrastructureName string, infrastructureKind string, infrastructureApiVersion string) *clusterapiv1beta1.Cluster {

	var (
		controlplaneRef   *corev1.ObjectReference
		infrastructureRef *corev1.ObjectReference
	)

	if controlPlaneName != "" {
		controlplaneRef = &corev1.ObjectReference{
			Kind:       controPlaneKind,
			Namespace:  name,
			Name:       controlPlaneName,
			APIVersion: controlPlaneApiVersion,
		}
	}

	if infrastructureName != "" {
		infrastructureRef = &corev1.ObjectReference{
			Kind:       infrastructureKind,
			Namespace:  name,
			Name:       infrastructureName,
			APIVersion: infrastructureApiVersion,
		}
	}

	namespace := GetTestClusterNamespace(name)
	testResource := clusterapiv1beta1.Cluster{
		TypeMeta: metav1.TypeMeta{
			Kind:       "Cluster",
			APIVersion: "cluster.x-k8s.io/v1beta1",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      name,
			Namespace: namespace,
			Labels: map[string]string{
				"region":       "us-east-1",
				"environment":  "test",
				"clusterGroup": "test-clusters",
			},
			ClusterName: name,
		},
		Spec: clusterapiv1beta1.ClusterSpec{
			Paused: false,
			ClusterNetwork: &clusterapiv1beta1.ClusterNetwork{
				APIServerPort: nil,
				Services:      &clusterapiv1beta1.NetworkRanges{CIDRBlocks: []string{"192.168.0.0/24"}},
				Pods:          &clusterapiv1beta1.NetworkRanges{CIDRBlocks: []string{"10.0.0.0/16"}},
				ServiceDomain: "cluster.local",
			},
			ControlPlaneEndpoint: clusterapiv1beta1.APIEndpoint{
				Host: "api." + name + ".cluster.example.com",
				Port: 443,
			},
			ControlPlaneRef:   controlplaneRef,
			InfrastructureRef: infrastructureRef,
		},
		Status: clusterapiv1beta1.ClusterStatus{},
	}

	return &testResource
}

func NewTestKopsMachinePool(name string, clusterName string) *clusterapikopsv1alpha1.KopsMachinePool {
	namespace := GetTestClusterNamespace(clusterName)
	testResource := clusterapikopsv1alpha1.KopsMachinePool{
		TypeMeta: metav1.TypeMeta{
			Kind:       "KopsMachinePool",
			APIVersion: "infrastructure.cluster.x-k8s.io/v1alpha1",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:        name,
			Namespace:   namespace,
			ClusterName: clusterName,
		},
		Spec: clusterapikopsv1alpha1.KopsMachinePoolSpec{
			KopsInstanceGroupSpec: v1alpha2.InstanceGroupSpec{
				MinSize:     nil,
				MaxSize:     nil,
				MachineType: "m5.xlarge",
				Subnets:     []string{"us-east-1a"},
			}},
		Status: clusterapikopsv1alpha1.KopsMachinePoolStatus{},
	}

	return &testResource
}

func NewTestMachinePool(name string, clusterName string, infrastructureKind string, infrastructureName string, infrastructureApiVersion string) *clusterapiexpv1beta1.MachinePool {
	var infrastructureRef corev1.ObjectReference

	if infrastructureName != "" {
		infrastructureRef = corev1.ObjectReference{
			Kind:       infrastructureKind,
			Namespace:  clusterName,
			Name:       infrastructureName,
			APIVersion: infrastructureApiVersion,
		}
	}

	namespace := GetTestClusterNamespace(clusterName)
	testResource := clusterapiexpv1beta1.MachinePool{
		TypeMeta: metav1.TypeMeta{
			Kind:       "MachinePool",
			APIVersion: "cluster.x-k8s.io/v1beta1",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:        name,
			Namespace:   namespace,
			ClusterName: clusterName,
		},
		Spec: clusterapiexpv1beta1.MachinePoolSpec{
			ClusterName: clusterName,
			Replicas:    nil,
			Template: clusterapiv1beta1.MachineTemplateSpec{
				ObjectMeta: clusterapiv1beta1.ObjectMeta{},
				Spec: clusterapiv1beta1.MachineSpec{
					ClusterName:       clusterName,
					Bootstrap:         clusterapiv1beta1.Bootstrap{},
					InfrastructureRef: infrastructureRef,
				},
			},
		},
		Status: clusterapiexpv1beta1.MachinePoolStatus{},
	}

	return &testResource
}

func NewTestMachineDeployment(name string, clusterName string, infrastructureKind string, infrastructureName string, infrastructureApiVersion string) *clusterapiv1beta1.MachineDeployment {
	namespace := GetTestClusterNamespace(clusterName)
	testResource := clusterapiv1beta1.MachineDeployment{
		TypeMeta: metav1.TypeMeta{
			Kind:       "MachineDeployment",
			APIVersion: "cluster.x-k8s.io/v1beta1",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:        name,
			Namespace:   namespace,
			ClusterName: clusterName,
		},
		Spec: clusterapiv1beta1.MachineDeploymentSpec{
			ClusterName: clusterName,
			Replicas:    nil,
			Template: clusterapiv1beta1.MachineTemplateSpec{
				ObjectMeta: clusterapiv1beta1.ObjectMeta{},
				Spec: clusterapiv1beta1.MachineSpec{
					ClusterName: clusterName,
					Bootstrap: clusterapiv1beta1.Bootstrap{
						ConfigRef: &corev1.ObjectReference{
							Kind:       "KubeadmConfigTemplate",
							Namespace:  clusterName,
							Name:       name,
							APIVersion: "bootstrap.cluster.x-k8s.io/v1beta1",
						},
					},
					InfrastructureRef: corev1.ObjectReference{
						Kind:       infrastructureKind,
						Namespace:  clusterName,
						Name:       infrastructureName,
						APIVersion: infrastructureApiVersion,
					},
				},
			},
		},
		Status: clusterapiv1beta1.MachineDeploymentStatus{},
	}

	return &testResource
}
