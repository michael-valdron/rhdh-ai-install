package deployment

import (
	"context"
	"encoding/json"
	"fmt"
	"slices"

	"github.com/redhat-ai-dev/rhdh-ai-install/pkg/client"
	"github.com/redhat-ai-dev/rhdh-ai-install/pkg/modelcatalog/serviceaccount"
	"github.com/redhat-ai-dev/rhdh-ai-install/pkg/util"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime/schema"
)

var (
	backstageContainersPath = []any{"spec", "deployment", "patch", "spec", "template", "spec", "containers"}
	backstageExtraFilesPath = []any{"spec", "application", "extraFiles"}
)

func getUnstructuredContainer(container corev1.Container) (map[string]any, error) {
	var (
		data   []byte
		result map[string]any
		err    error
	)

	if data, err = container.Marshal(); err != nil {
		return result, err
	}

	if err = json.Unmarshal(data, &result); err != nil {
		return result, err
	}

	return result, err
}

func getBackstageSecretByNameFunc(name string) func(map[string]string) bool {
	return func(secret map[string]string) bool { return secret["name"] == name }
}

func getBackstageContainerByNameFunc(name string) func(map[string]any) bool {
	return func(container map[string]any) bool { return container["name"] == name }
}

func patchBackstageContainers(backstageSpec *unstructured.Unstructured, rootPath []any, subPath []any) error {
	var (
		locationContainer    map[string]any
		storageRestContainer map[string]any
		normalizerContainer  map[string]any
		err                  error
	)

	for idx := 0; idx < len(subPath); idx++ {
		util.SetValueInUnstructured(backstageSpec.Object, make(map[string]any), rootPath...)
		rootPath = append(rootPath, subPath[idx])
		subPath = subPath[idx+1:]
	}

	if locationContainer, err = getUnstructuredContainer(getLocationContainer()); err != nil {
		return err
	}

	if storageRestContainer, err = getUnstructuredContainer(getStorageRestContainer()); err != nil {
		return err
	}

	if normalizerContainer, err = getUnstructuredContainer(getNormalizerContainer()); err != nil {
		return err
	}

	if !util.ExistsInUnstructured(backstageSpec.Object, rootPath...) {
		util.SetValueInUnstructured(backstageSpec.Object, []map[string]any{
			{
				"name": "backstage-backend",
			},
			locationContainer,
			storageRestContainer,
			normalizerContainer,
		}, rootPath...)
	} else {
		currentValue := util.GetValueInUnstructured(backstageSpec.Object, rootPath...)
		containersPatch, err := util.UnstructuredToSlice[map[string]any](currentValue)
		if err != nil {
			return err
		}

		containersPatch = append(containersPatch, locationContainer, storageRestContainer, normalizerContainer)

		util.SetValueInUnstructured(backstageSpec.Object, containersPatch, rootPath...)
	}

	return nil
}

func NewBackstageGVR() *schema.GroupVersionResource {
	return &schema.GroupVersionResource{
		Group:    backstageGroup,
		Version:  backstageVersion,
		Resource: backstageResource,
	}
}

func PatchBackstageSpec(backstageSpec *unstructured.Unstructured) error {
	serviceAccountTokenName := serviceaccount.GetServiceAccountTokenName(serviceaccount.ServiceAccountName)

	// Patch volume mount to backstage container for service account token volume
	if !util.ExistsInUnstructured(backstageSpec.Object, backstageExtraFilesPath...) {
		util.SetValueInUnstructured(backstageSpec.Object, map[string]any{
			"mountPath": containerWorkingDir,
			"secrets": []map[string]string{
				{
					"key":       "token",
					"mountPath": containerWorkingDir,
					"name":      serviceAccountTokenName,
				},
			},
		}, backstageExtraFilesPath...)
	} else {
		currentValue := util.GetValueInUnstructured(backstageSpec.Object, backstageExtraFilesPath...)
		extraFilesPatch, err := util.UnstructuredToMap[any](currentValue)
		if err != nil {
			return err
		}

		currentSecrets, err := util.UnstructuredToSlice[map[string]string](extraFilesPatch["secrets"])
		if err != nil {
			return err
		}

		if !slices.ContainsFunc(currentSecrets, getBackstageSecretByNameFunc(serviceAccountTokenName)) {
			extraFilesPatch["secrets"] = append(currentSecrets, map[string]string{
				"key":       "token",
				"mountPath": containerWorkingDir,
				"name":      serviceAccountTokenName,
			})

			util.SetValueInUnstructured(backstageSpec.Object, extraFilesPatch, backstageExtraFilesPath...)
		}
	}

	// Patch sidecar containers for model catalog
	rootPath := backstageContainersPath[:2]
	subPath := backstageContainersPath[2:]
	for len(subPath) != 0 {
		if util.ExistsInUnstructured(backstageSpec.Object, rootPath...) {
			rootPath = append(rootPath, subPath[0])
			subPath = subPath[1:]
		} else {
			break
		}
	}
	return patchBackstageContainers(backstageSpec, rootPath, subPath)
}

func UnpatchBackstageSpec(backstageSpec *unstructured.Unstructured) error {
	serviceAccountTokenName := serviceaccount.GetServiceAccountTokenName(serviceaccount.ServiceAccountName)
	patchContainers := []string{locationContainerName, storageRestContainerName, normalizerContainerName}
	secretsPath := append(backstageExtraFilesPath, "secrets")
	var (
		backstageContainersPatch []map[string]any
		err                      error
	)

	// Remove model catalog sidecar containers
	if backstageContainersPatch, err = util.UnstructuredToSlice[map[string]any](util.GetValueInUnstructured(backstageSpec.Object, backstageContainersPath...)); err != nil {
		return err
	}
	for _, containerName := range patchContainers {
		if backstageContainersPatch, err = util.RemoveFromSliceByFunc(backstageContainersPatch, getBackstageContainerByNameFunc(containerName)); err != nil {
			return err
		}
	}
	util.SetValueInUnstructured(backstageSpec.Object, backstageContainersPatch, backstageContainersPath...)

	// Remove service account token volume
	if !util.ExistsInUnstructured(backstageSpec.Object, secretsPath...) {
		return fmt.Errorf("expected %s to be set", util.Join(secretsPath, "."))
	}

	if currentSecrets, err := util.UnstructuredToSlice[map[string]string](util.GetValueInUnstructured(backstageSpec.Object, secretsPath...)); err != nil {
		return err
	} else {
		if currentSecrets, err = util.RemoveFromSliceByFunc(currentSecrets, getBackstageSecretByNameFunc(serviceAccountTokenName)); err != nil {
			return err
		}

		util.SetValueInUnstructured(backstageSpec.Object, currentSecrets, secretsPath...)
	}

	return err
}

func CheckBackstage(ctx context.Context, c *client.Client, backstage *unstructured.Unstructured, opts metav1.GetOptions) error {
	panic("unimplemented")
}
