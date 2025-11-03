package infra

const (
	defaultRHOAIOperatorNamespace = "redhat-ods-operator"
	rhoaiOperatorControllerName   = "rhods-operator"

	defaultNFDOperatorNamespace = "openshift-nfd"
	nfdOperatorControllerName   = "nfd-controller-manager"

	defaultNvidiaGPUOperatorNamespace = "nvidia-gpu-operator"
	nvidiaGPUOperatorControllerName   = "gpu-operator"

	defaultKMMOperatorNamespace = "openshift-kmm"
	kmmOperatorControllerName   = "kmm-operator-controller"
	kmmOperatorWebhookName      = "kmm-operator-webhook"

	defaultAuthorinoOperatorNamespace = "openshift-operators"
	authorinoOperatorControllerName   = "authorino-operator"

	dataScienceClusterGroup = "datasciencecluster.opendatahub.io"
	nfdGroup                = "nfd.openshift.io"
	clusterPolicyGroup      = "nvidia.com"

	dataScienceClusterKind = "DataScienceCluster"
	nfdKind                = "NodeFeatureDiscovery"
	clusterPolicyKind      = "ClusterPolicy"
)
