package test

import (
	clusterapikopsv1alpha1 "github.com/topfreegames/kubernetes-kops-operator/apis/infrastructure/v1alpha1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/dynamic/fake"
	"k8s.io/kops/pkg/apis/kops/v1alpha2"
	"log"
	clusterapiv1beta1 "sigs.k8s.io/cluster-api/api/v1beta1"
	clusterapiexpv1beta1 "sigs.k8s.io/cluster-api/exp/api/v1beta1"
)

type K8sRequest struct {
	ResourceName  string
	ResourceKind  string
	Cluster       string
	TestResources []runtime.Object
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

func NewK8sFakeDynamicClientWithResource(resource interface{}) *fake.FakeDynamicClient {
	generic, err := runtime.DefaultUnstructuredConverter.ToUnstructured(&resource)
	if err != nil {
		log.Fatalf("Error generating unstructured: %v", err)
	}
	Unstructured := &unstructured.Unstructured{}
	Unstructured.SetUnstructuredContent(generic)
	client := fake.NewSimpleDynamicClient(runtime.NewScheme(), Unstructured)
	return client
}

func NewTestKopsMachinePool(name string, clusterName string) *clusterapikopsv1alpha1.KopsMachinePool {
	testResource := clusterapikopsv1alpha1.KopsMachinePool{
		TypeMeta: v1.TypeMeta{
			Kind:       "KopsMachinePool",
			APIVersion: "infrastructure.cluster.x-k8s.io/v1alpha1",
		},
		ObjectMeta: v1.ObjectMeta{
			Name:        name,
			Namespace:   clusterName,
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

	testResource := clusterapiexpv1beta1.MachinePool{
		TypeMeta:   v1.TypeMeta{
			Kind:       "MachinePool",
			APIVersion: "cluster.x-k8s.io/v1beta1",
		},
		ObjectMeta: v1.ObjectMeta{
			Name:        name,
			Namespace:   clusterName,
			ClusterName: clusterName,
		},
		Spec:       clusterapiexpv1beta1.MachinePoolSpec{
			ClusterName:     clusterName,
			Replicas:        nil,
			Template:        clusterapiv1beta1.MachineTemplateSpec{
				ObjectMeta: clusterapiv1beta1.ObjectMeta{},
				Spec:       clusterapiv1beta1.MachineSpec{
					ClusterName:       clusterName,
					Bootstrap:         clusterapiv1beta1.Bootstrap{},
					InfrastructureRef: corev1.ObjectReference{
						Kind:            infrastructureKind,
						Namespace:       clusterName,
						Name:            infrastructureName,
						APIVersion:      infrastructureApiVersion,
					},
				},
			},
		},
		Status:     clusterapiexpv1beta1.MachinePoolStatus{},
	}

	return &testResource
}

func NewTestMachineDeployment(name string, clusterName string, infrastructureKind string, infrastructureName string, infrastructureApiVersion string) *clusterapiv1beta1.MachineDeployment {

	testResource := clusterapiv1beta1.MachineDeployment{
		TypeMeta:   v1.TypeMeta{
			Kind:       "MachineDeployment",
			APIVersion: "cluster.x-k8s.io/v1beta1",
		},
		ObjectMeta: v1.ObjectMeta{
			Name:        name,
			Namespace:   clusterName,
			ClusterName: clusterName,
		},
		Spec:       clusterapiv1beta1.MachineDeploymentSpec{
			ClusterName:             clusterName,
			Replicas:                nil,
			Template:                clusterapiv1beta1.MachineTemplateSpec{
				ObjectMeta: clusterapiv1beta1.ObjectMeta{},
				Spec:       clusterapiv1beta1.MachineSpec{
					ClusterName:       clusterName,
					Bootstrap:         clusterapiv1beta1.Bootstrap{
						ConfigRef:      &corev1.ObjectReference{
							Kind:            "KubeadmConfigTemplate",
							Namespace:       clusterName,
							Name:            name,
							APIVersion:      "bootstrap.cluster.x-k8s.io/v1beta1",
						},
					},
					InfrastructureRef: corev1.ObjectReference{
						Kind:            infrastructureKind,
						Namespace:       clusterName,
						Name:            infrastructureName,
						APIVersion:      infrastructureApiVersion,
					},
				},
			},
		},
		Status:     clusterapiv1beta1.MachineDeploymentStatus{},
	}

	return &testResource
}