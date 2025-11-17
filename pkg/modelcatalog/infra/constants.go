package infra

const (
	// Red Hat OpenShift AI operator values

	defaultRHOAIOperatorNamespace = "redhat-ods-operator" // default RHOAI operator namespace
	rhoaiOperatorControllerName   = "rhods-operator"      // RHOAI operator controller container name

	// Node Feature Discovery operator values

	defaultNFDOperatorNamespace = "openshift-nfd"          // default Node Feature Discovery operator namespace
	nfdOperatorControllerName   = "nfd-controller-manager" // Node Feature Discovery operator controller container name

	// Nvidia GPU operator values

	defaultNvidiaGPUOperatorNamespace = "nvidia-gpu-operator" // default Nvidia GPU operator namespace
	nvidiaGPUOperatorControllerName   = "gpu-operator"        // Nvidia GPU operator controller container name

	// Kernel Module Manager operator values

	defaultKMMOperatorNamespace = "openshift-kmm"           // default KMM operator namespace
	kmmOperatorControllerName   = "kmm-operator-controller" // KMM operator controller container name
	kmmOperatorWebhookName      = "kmm-operator-webhook"    // KMM operator webhook container name

	// Authorino operator values

	defaultAuthorinoOperatorNamespace = "openshift-operators" // default Authorino operator namespace
	authorinoOperatorControllerName   = "authorino-operator"  // Authorino operator controller container name

	// cluster resource groups

	dataScienceClusterGroup = "datasciencecluster.opendatahub.io" // DataScienceCluster CRD group name
	nfdGroup                = "nfd.openshift.io"                  // NodeFeatureDiscovery CRD group name
	clusterPolicyGroup      = "nvidia.com"                        // ClusterPolicy CRD group name

	// cluster resource kinds

	dataScienceClusterKind = "DataScienceCluster"   // DataScienceCluster CRD kind
	nfdKind                = "NodeFeatureDiscovery" // NodeFeatureDiscovery CRD kind
	clusterPolicyKind      = "ClusterPolicy"        // ClusterPolicy CRD kind
)
