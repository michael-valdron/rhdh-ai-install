package deployment

const (
	backstageContainerName   = "backstage-backend" // Developer Hub base container name
	locationContainerName    = "location"          // location sidecar container name
	storageRestContainerName = "storage-rest"      // storage-rest sidecar container name
	normalizerContainerName  = "rhoai-normalizer"  // normalizer sidecar container name
	containerWorkingDir      = "/opt/app-root/src" // working directory for the containers

	// Backstage CRD values

	backstageVersion  = "v1alpha4"        // Backstage CRD version
	backstageGroup    = "rhdh.redhat.com" // Backstage CRD group name
	backstageResource = "backstages"      // Backstage CRD resource name
)
