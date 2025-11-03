package deployment

import (
	"context"
	"fmt"
	"slices"

	"github.com/redhat-ai-dev/rhdh-ai-install/pkg/client"
	"github.com/redhat-ai-dev/rhdh-ai-install/pkg/modelcatalog/serviceaccount"
	"github.com/redhat-ai-dev/rhdh-ai-install/pkg/util"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
)

var (
	normalizerEnvVar = corev1.EnvVar{
		Name:  "NORMALIZER_FORMAT",
		Value: "JsonArrayFormat",
	}
	storageRestEnvVars = []corev1.EnvVar{
		{
			Name:  "STORAGE_TYPE",
			Value: "ConfigMap",
		},
		{
			Name:  "BRIDGE_URL",
			Value: "http://localhost:9090",
		},
	}
	pluginsVolumeMount = corev1.VolumeMount{
		MountPath: "/opt/app-root/src/dynamic-plugins-root",
		Name:      "dynamic-plugins-root",
	}
)

func getPodEnvVars() []corev1.EnvVar {
	return []corev1.EnvVar{
		{
			Name: "POD_IP",
			ValueFrom: &corev1.EnvVarSource{
				FieldRef: &corev1.ObjectFieldSelector{
					FieldPath: "status.podIP",
				},
			},
		},
		{
			Name: "POD_NAMESPACE",
			ValueFrom: &corev1.EnvVarSource{
				FieldRef: &corev1.ObjectFieldSelector{
					FieldPath: "metadata.namespace",
				},
			},
		},
	}
}

func getServiceAccountTokenRef() corev1.EnvFromSource {
	serviceAccountTokenName := serviceaccount.GetServiceAccountTokenName(serviceaccount.ServiceAccountName)
	return corev1.EnvFromSource{
		SecretRef: &corev1.SecretEnvSource{
			LocalObjectReference: corev1.LocalObjectReference{
				Name: serviceAccountTokenName,
			},
		},
	}
}

func getServiceAccountVolume() corev1.Volume {
	serviceAccountTokenName := serviceaccount.GetServiceAccountTokenName(serviceaccount.ServiceAccountName)
	return corev1.Volume{
		Name: serviceAccountTokenName,
		VolumeSource: corev1.VolumeSource{
			Secret: &corev1.SecretVolumeSource{
				SecretName: serviceAccountTokenName,
				Items: []corev1.KeyToPath{
					{
						Key: "token",
					},
				},
			},
		},
	}
}

func getLocationContainer() corev1.Container {
	return corev1.Container{
		Name:            locationContainerName,
		Image:           "quay.io/redhat-ai-dev/model-catalog-location-service:latest",
		ImagePullPolicy: corev1.PullAlways,
		Env:             append(getPodEnvVars(), normalizerEnvVar),
		EnvFrom:         []corev1.EnvFromSource{getServiceAccountTokenRef()},
		Ports: []corev1.ContainerPort{
			{
				ContainerPort: 9090,
				Name:          "location",
				Protocol:      corev1.ProtocolTCP,
			},
		},
		VolumeMounts: []corev1.VolumeMount{pluginsVolumeMount},
		WorkingDir:   containerWorkingDir,
	}
}

func getStorageRestContainer() corev1.Container {
	return corev1.Container{
		Name:            storageRestContainerName,
		Image:           "quay.io/redhat-ai-dev/model-catalog-storage-rest:latest",
		ImagePullPolicy: corev1.PullAlways,
		Env:             append(getPodEnvVars(), append(storageRestEnvVars, normalizerEnvVar)...),
		EnvFrom:         []corev1.EnvFromSource{getServiceAccountTokenRef()},
		VolumeMounts:    []corev1.VolumeMount{pluginsVolumeMount},
		WorkingDir:      containerWorkingDir,
	}
}

func getNormalizerContainer() corev1.Container {
	return corev1.Container{
		Name:            normalizerContainerName,
		Image:           "quay.io/redhat-ai-dev/model-catalog-rhoai-normalizer:latest",
		ImagePullPolicy: corev1.PullAlways,
		Env:             append(getPodEnvVars(), normalizerEnvVar),
		EnvFrom:         []corev1.EnvFromSource{getServiceAccountTokenRef()},
		VolumeMounts:    []corev1.VolumeMount{pluginsVolumeMount},
		WorkingDir:      containerWorkingDir,
	}
}

func getContainerByNameFunc(name string) func(corev1.Container) bool {
	return func(container corev1.Container) bool { return container.Name == name }
}

func getVolumeMountByNameFunc(name string) func(corev1.VolumeMount) bool {
	return func(vm corev1.VolumeMount) bool { return vm.Name == name }
}

func getVolumeByNameFunc(name string) func(corev1.Volume) bool {
	return func(v corev1.Volume) bool { return v.Name == name }
}

func PatchDeploymentSpec(deploymentSpec *appsv1.Deployment) error {
	serviceAccountTokenName := serviceaccount.GetServiceAccountTokenName(serviceaccount.ServiceAccountName)

	// Patch volume mapping to service account token secret
	deploymentSpec.Spec.Template.Spec.Volumes = append(deploymentSpec.Spec.Template.Spec.Volumes, getServiceAccountVolume())

	// Patch volume mount to backstage container for service account token volume
	util.ModifySliceElementByFunc(deploymentSpec.Spec.Template.Spec.Containers, getContainerByNameFunc(backstageContainerName), func(container *corev1.Container) {
		container.VolumeMounts = append(container.VolumeMounts, corev1.VolumeMount{
			MountPath: containerWorkingDir,
			Name:      serviceAccountTokenName,
		})
	})

	// Patch sidecar containers for model catalog
	deploymentSpec.Spec.Template.Spec.Containers = append(deploymentSpec.Spec.Template.Spec.Containers,
		getLocationContainer(),
		getStorageRestContainer(),
		getNormalizerContainer(),
	)
	return nil
}

func UnpatchDeploymentSpec(deploymentSpec *appsv1.Deployment) error {
	serviceAccountTokenName := serviceaccount.GetServiceAccountTokenName(serviceaccount.ServiceAccountName)
	patchContainers := []string{locationContainerName, storageRestContainerName, normalizerContainerName}
	var err error

	// Remove model catalog sidecar containers
	for _, containerName := range patchContainers {
		deploymentSpec.Spec.Template.Spec.Containers, err = util.RemoveFromSliceByFunc(deploymentSpec.Spec.Template.Spec.Containers, getContainerByNameFunc(containerName))
		if err != nil {
			return err
		}
	}

	// Remove service account token volume mount
	util.ModifySliceElementByFunc(deploymentSpec.Spec.Template.Spec.Containers, getContainerByNameFunc("backstage-backend"), func(container *corev1.Container) {
		container.VolumeMounts, err = util.RemoveFromSliceByFunc(container.VolumeMounts, getVolumeMountByNameFunc(serviceAccountTokenName))
	})
	if err != nil {
		return err
	}

	// Remove service account token volume
	deploymentSpec.Spec.Template.Spec.Volumes, err = util.RemoveFromSliceByFunc(deploymentSpec.Spec.Template.Spec.Volumes, getVolumeByNameFunc(serviceAccountTokenName))

	return err
}

func PatchDeployment(ctx context.Context, c *client.Client, deploymentPatch *appsv1.Deployment, opts metav1.PatchOptions) error {
	var (
		patchData []byte
		err       error
	)

	if patchData, err = deploymentPatch.Marshal(); err != nil {
		return err
	}

	_, err = c.Deployments(client.ClientParams{Namespace: deploymentPatch.Namespace}).Patch(ctx, deploymentPatch.Name, types.StrategicMergePatchType, patchData, opts)

	return err
}

func CheckDeployment(ctx context.Context, c *client.Client, deployment *appsv1.Deployment, opts metav1.GetOptions) error {
	serviceAccountTokenName := serviceaccount.GetServiceAccountTokenName(serviceaccount.ServiceAccountName)
	errs := []any{}

	if slices.ContainsFunc(deployment.Spec.Template.Spec.Volumes, getVolumeByNameFunc(serviceAccountTokenName)) {
		errs = append(errs, fmt.Errorf("%s volume is already present in the deployment", serviceAccountTokenName))
	}

	if container, err := util.GetElementFromSliceByFunc(deployment.Spec.Template.Spec.Containers, getContainerByNameFunc(backstageContainerName)); err != nil {
		errs = append(errs, fmt.Errorf("error getting backstage container %s: %v", backstageContainerName, err))
	} else if slices.ContainsFunc(container.VolumeMounts, getVolumeMountByNameFunc(serviceAccountTokenName)) {
		errs = append(errs, fmt.Errorf("%s volume mount is already present in the backstage (\"%s\") container", serviceAccountTokenName, backstageContainerName))
	}

	for _, containerName := range []string{locationContainerName, storageRestContainerName, normalizerContainerName} {
		if slices.ContainsFunc(deployment.Spec.Template.Spec.Containers, getContainerByNameFunc(containerName)) {
			errs = append(errs, fmt.Errorf("%s container is already present in the deployment", containerName))
		}
	}

	if len(errs) != 0 {
		if c.Options.Verbose {
			return fmt.Errorf("model catalog deployment resources are not ready, see the following errors (note: remember uninstall first if there is a previous install):\n%s", util.Join(errs, "\n"))
		} else {
			return fmt.Errorf("model catalog deployment resources are not ready")
		}
	} else {
		return nil
	}
}
