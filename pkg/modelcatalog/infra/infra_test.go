package infra

import (
	"context"
	"reflect"
	"testing"

	"github.com/redhat-ai-dev/rhdh-ai-install/pkg/client"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	apiextv1 "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
)

// TestFindDeploymentFromAllNamespaces tests if findDeploymentFromAllNamespaces returns the correct deployment object from all namespaces.
func TestFindDeploymentFromAllNamespaces(t *testing.T) {
	// Setup
	deploymentName := "test-deployment"
	namespace := &corev1.Namespace{
		ObjectMeta: metav1.ObjectMeta{
			Name: "test-namespace",
		},
	}
	expectedDeployment := &appsv1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name:      deploymentName,
			Namespace: namespace.Name,
		},
	}
	mockClient, err := client.NewMock(client.ClientOptions{}, []runtime.Object{namespace, expectedDeployment}, []runtime.Object{})
	if err != nil {
		t.Fatalf("unexpected error creating the mock client: %v", err)
	}

	// Test
	result, err := findDeploymentFromAllNamespaces(context.Background(), mockClient, deploymentName)
	if err != nil {
		t.Fatalf("unexpected error running findDeploymentFromAllNamespaces: %v", err)
	}

	// Assertions
	if !reflect.DeepEqual(expectedDeployment, result) {
		t.Errorf("Expected: %v\nActual: %v", expectedDeployment, result)
	}
}

// TestCheckForOpenShiftAI tests if checkForOpenShiftAI returns nil when OpenShift AI operator controller deployment exists and is ready in the default namespace.
func TestCheckForOpenShiftAI(t *testing.T) {
	// Setup
	namespace := &corev1.Namespace{
		ObjectMeta: metav1.ObjectMeta{
			Name: defaultRHOAIOperatorNamespace,
		},
	}
	expectedController := &appsv1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name:      rhoaiOperatorControllerName,
			Namespace: defaultRHOAIOperatorNamespace,
		},
		Status: appsv1.DeploymentStatus{
			ReadyReplicas: 1,
		},
	}
	dataScienceClusterCRD := &apiextv1.CustomResourceDefinition{
		ObjectMeta: metav1.ObjectMeta{
			Name: dataScienceClusterKind,
		},
	}
	var (
		mockClient *client.Client
		err        error
	)

	// Set CRD GroupVersionKind to match DataScienceCluster CRD
	dataScienceClusterCRD.SetGroupVersionKind(schema.GroupVersionKind{
		Group:   dataScienceClusterGroup,
		Version: "v1",
		Kind:    dataScienceClusterKind,
	})

	mockClient, err = client.NewMock(client.ClientOptions{}, []runtime.Object{namespace, expectedController}, []runtime.Object{dataScienceClusterCRD})
	if err != nil {
		t.Fatalf("unexpected error creating the mock client: %v", err)
	}

	// Test
	err = checkForOpenShiftAI(context.Background(), mockClient, nil)

	// Assertions
	if err != nil {
		t.Errorf("Expected no error, got: %v", err)
	}
}

// TestCheckForNvidiaGPU tests if checkForNvidiaGPU returns nil when Nvidia GPU operator controller deployment exists and is ready in the default namespace.
func TestCheckForNvidiaGPU(t *testing.T) {
	// Setup
	namespace := &corev1.Namespace{
		ObjectMeta: metav1.ObjectMeta{
			Name: defaultNvidiaGPUOperatorNamespace,
		},
	}
	expectedController := &appsv1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name:      nvidiaGPUOperatorControllerName,
			Namespace: defaultNvidiaGPUOperatorNamespace,
		},
		Status: appsv1.DeploymentStatus{
			ReadyReplicas: 1,
		},
	}
	clusterPolicyCRD := &apiextv1.CustomResourceDefinition{
		ObjectMeta: metav1.ObjectMeta{
			Name: clusterPolicyKind,
		},
	}
	var (
		mockClient *client.Client
		err        error
	)

	// Set CRD GroupVersionKind to match ClusterPolicy CRD
	clusterPolicyCRD.SetGroupVersionKind(schema.GroupVersionKind{
		Group:   clusterPolicyGroup,
		Version: "v1",
		Kind:    clusterPolicyKind,
	})

	mockClient, err = client.NewMock(client.ClientOptions{}, []runtime.Object{namespace, expectedController}, []runtime.Object{clusterPolicyCRD})
	if err != nil {
		t.Fatalf("unexpected error creating the mock client: %v", err)
	}

	// Test
	err = checkForNvidiaGPU(context.Background(), mockClient, nil)

	// Assertions
	if err != nil {
		t.Errorf("Expected no error, got: %v", err)
	}
}

// TestCheckForNodeFeatureDiscovery tests if checkForNodeFeatureDiscovery returns nil when Node Feature Discovery operator controller deployment exists and is ready in the default namespace.
func TestCheckForNodeFeatureDiscovery(t *testing.T) {
	// Setup
	namespace := &corev1.Namespace{
		ObjectMeta: metav1.ObjectMeta{
			Name: defaultNFDOperatorNamespace,
		},
	}
	expectedController := &appsv1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name:      nfdOperatorControllerName,
			Namespace: defaultNFDOperatorNamespace,
		},
		Status: appsv1.DeploymentStatus{
			ReadyReplicas: 1,
		},
	}
	nodeFeatureDiscoveryCRD := &apiextv1.CustomResourceDefinition{
		ObjectMeta: metav1.ObjectMeta{
			Name: nfdKind,
		},
	}
	var (
		mockClient *client.Client
		err        error
	)

	// Set CRD GroupVersionKind to match NodeFeatureDiscovery CRD
	nodeFeatureDiscoveryCRD.SetGroupVersionKind(schema.GroupVersionKind{
		Group:   nfdGroup,
		Version: "v1",
		Kind:    nfdKind,
	})

	mockClient, err = client.NewMock(client.ClientOptions{}, []runtime.Object{namespace, expectedController}, []runtime.Object{nodeFeatureDiscoveryCRD})
	if err != nil {
		t.Fatalf("unexpected error creating the mock client: %v", err)
	}

	// Test
	err = checkForNodeFeatureDiscovery(context.Background(), mockClient, nil)

	// Assertions
	if err != nil {
		t.Errorf("Expected no error, got: %v", err)
	}
}

// TestCheckForKMM tests if checkForKMM returns nil when Kernel Module Management (KMM) operator controller deployment exists and is ready in the default namespace.
func TestCheckForKMM(t *testing.T) {
	// Setup
	namespace := &corev1.Namespace{
		ObjectMeta: metav1.ObjectMeta{
			Name: defaultKMMOperatorNamespace,
		},
	}
	expectedController := &appsv1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name:      kmmOperatorControllerName,
			Namespace: defaultKMMOperatorNamespace,
		},
		Status: appsv1.DeploymentStatus{
			ReadyReplicas: 1,
		},
	}
	expectedWebhook := &appsv1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name:      kmmOperatorWebhookName,
			Namespace: defaultKMMOperatorNamespace,
		},
		Status: appsv1.DeploymentStatus{
			ReadyReplicas: 1,
		},
	}
	var (
		mockClient *client.Client
		err        error
	)

	mockClient, err = client.NewMock(client.ClientOptions{}, []runtime.Object{namespace, expectedController, expectedWebhook}, []runtime.Object{})
	if err != nil {
		t.Fatalf("unexpected error creating the mock client: %v", err)
	}

	// Test
	err = checkForKMM(context.Background(), mockClient)

	// Assertions
	if err != nil {
		t.Errorf("Expected no error, got: %v", err)
	}
}

// TestCheckForAuthorino tests if checkForAuthorino returns nil when Authorino operator controller deployment exists and is ready in the default namespace.
func TestCheckForAuthorino(t *testing.T) {
	// Setup
	namespace := &corev1.Namespace{
		ObjectMeta: metav1.ObjectMeta{
			Name: defaultAuthorinoOperatorNamespace,
		},
	}
	expectedDeployment := &appsv1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name:      authorinoOperatorControllerName,
			Namespace: defaultAuthorinoOperatorNamespace,
		},
		Status: appsv1.DeploymentStatus{
			ReadyReplicas: 1,
		},
	}
	var (
		mockClient *client.Client
		err        error
	)

	mockClient, err = client.NewMock(client.ClientOptions{}, []runtime.Object{namespace, expectedDeployment}, []runtime.Object{})
	if err != nil {
		t.Fatalf("unexpected error creating the mock client: %v", err)
	}

	// Test
	err = checkForAuthorino(context.Background(), mockClient)

	// Assertions
	if err != nil {
		t.Errorf("Expected no error, got: %v", err)
	}
}

// TestCheck checks if Check returns nil when all infrastructure resources are ready.
func TestCheck(t *testing.T) {
	// Setup
	namespaces := &corev1.NamespaceList{
		Items: []corev1.Namespace{
			// RHOAI operator namespace
			{
				ObjectMeta: metav1.ObjectMeta{
					Name: defaultRHOAIOperatorNamespace,
				},
			},
			// Nvidia GPU operator namespace
			{
				ObjectMeta: metav1.ObjectMeta{
					Name: defaultNvidiaGPUOperatorNamespace,
				},
			},
			// NFD operator namespace
			{
				ObjectMeta: metav1.ObjectMeta{
					Name: defaultNFDOperatorNamespace,
				},
			},
			// KMM operator namespace
			{
				ObjectMeta: metav1.ObjectMeta{
					Name: defaultKMMOperatorNamespace,
				},
			},
			// Authorino operator namespace
			{
				ObjectMeta: metav1.ObjectMeta{
					Name: defaultAuthorinoOperatorNamespace,
				},
			},
		},
	}
	deployments := &appsv1.DeploymentList{
		Items: []appsv1.Deployment{
			// RHOAI operator controller deployment
			{
				ObjectMeta: metav1.ObjectMeta{
					Name:      rhoaiOperatorControllerName,
					Namespace: defaultRHOAIOperatorNamespace,
				},
				Status: appsv1.DeploymentStatus{
					ReadyReplicas: 1,
				},
			},
			// Nvidia GPU controller deployment
			{
				ObjectMeta: metav1.ObjectMeta{
					Name:      nvidiaGPUOperatorControllerName,
					Namespace: defaultNvidiaGPUOperatorNamespace,
				},
				Status: appsv1.DeploymentStatus{
					ReadyReplicas: 1,
				},
			},
			// NFD controller deployment
			{
				ObjectMeta: metav1.ObjectMeta{
					Name:      nfdOperatorControllerName,
					Namespace: defaultNFDOperatorNamespace,
				},
				Status: appsv1.DeploymentStatus{
					ReadyReplicas: 1,
				},
			},
			// KMM controller deployment
			{
				ObjectMeta: metav1.ObjectMeta{
					Name:      kmmOperatorControllerName,
					Namespace: defaultKMMOperatorNamespace,
				},
				Status: appsv1.DeploymentStatus{
					ReadyReplicas: 1,
				},
			},
			// KMM webhooks deployment
			{
				ObjectMeta: metav1.ObjectMeta{
					Name:      kmmOperatorWebhookName,
					Namespace: defaultKMMOperatorNamespace,
				},
				Status: appsv1.DeploymentStatus{
					ReadyReplicas: 1,
				},
			},
			// Authorino controller deployment
			{
				ObjectMeta: metav1.ObjectMeta{
					Name:      authorinoOperatorControllerName,
					Namespace: defaultAuthorinoOperatorNamespace,
				},
				Status: appsv1.DeploymentStatus{
					ReadyReplicas: 1,
				},
			},
		},
	}
	crds := &apiextv1.CustomResourceDefinitionList{
		Items: []apiextv1.CustomResourceDefinition{
			// DataScienceCluster
			{
				ObjectMeta: metav1.ObjectMeta{
					Name: dataScienceClusterKind,
				},
			},
			// ClusterPolicy
			{
				ObjectMeta: metav1.ObjectMeta{
					Name: clusterPolicyKind,
				},
			},
			// NodeFeatureDiscovery
			{
				ObjectMeta: metav1.ObjectMeta{
					Name: nfdKind,
				},
			},
		},
	}

	crds.Items[0].SetGroupVersionKind(schema.GroupVersionKind{
		Group:   dataScienceClusterGroup,
		Version: "v1",
		Kind:    dataScienceClusterKind,
	})
	crds.Items[1].SetGroupVersionKind(schema.GroupVersionKind{
		Group:   clusterPolicyGroup,
		Version: "v1",
		Kind:    clusterPolicyKind,
	})
	crds.Items[2].SetGroupVersionKind(schema.GroupVersionKind{
		Group:   nfdGroup,
		Version: "v1",
		Kind:    nfdKind,
	})

	mockClient, err := client.NewMock(client.ClientOptions{Verbose: true}, []runtime.Object{namespaces, deployments}, []runtime.Object{crds})
	if err != nil {
		t.Fatalf("unexpected error creating the mock client: %v", err)
	}

	// Test
	err = Check(context.Background(), mockClient)

	// Assertions
	if err != nil {
		t.Errorf("Expected no error, got: %v", err)
	}
}
