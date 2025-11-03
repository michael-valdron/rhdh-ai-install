package deployment

import (
	"reflect"
	"testing"

	"github.com/redhat-ai-dev/rhdh-ai-install/pkg/modelcatalog/serviceaccount"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/utils/ptr"
)

func makeMockDeployment() appsv1.Deployment {
	return appsv1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "test-deploy",
			Namespace: "test-deploy",
		},
		Spec: appsv1.DeploymentSpec{
			Template: corev1.PodTemplateSpec{
				Spec: corev1.PodSpec{
					Volumes: []corev1.Volume{
						{
							Name: "dynamic-plugins-root",
							VolumeSource: corev1.VolumeSource{
								Ephemeral: &corev1.EphemeralVolumeSource{
									VolumeClaimTemplate: &corev1.PersistentVolumeClaimTemplate{
										Spec: corev1.PersistentVolumeClaimSpec{
											AccessModes: []corev1.PersistentVolumeAccessMode{
												corev1.ReadWriteOnce,
											},
											VolumeMode: ptr.To(corev1.PersistentVolumeFilesystem),
										},
									},
								},
							},
						},
						{
							Name: "dynamic-plugins",
							VolumeSource: corev1.VolumeSource{
								ConfigMap: &corev1.ConfigMapVolumeSource{
									LocalObjectReference: corev1.LocalObjectReference{
										Name: "redhat-developer-hub-dynamic-plugins",
									},
									DefaultMode: ptr.To(int32(420)),
									Optional:    ptr.To(true),
								},
							},
						},
					},
					Containers: []corev1.Container{
						{
							Name: "backstage-backend",
							VolumeMounts: []corev1.VolumeMount{
								{
									Name:      "backstage-app-config",
									MountPath: "/opt/app-root/src/app-config-from-configmap.yaml",
								},
							},
						},
					},
				},
			},
		},
	}
}

func makeMockPatchedDeployment() appsv1.Deployment {
	serviceAccountTokenName := serviceaccount.GetServiceAccountTokenName(serviceaccount.ServiceAccountName)

	deploySpec := makeMockDeployment()

	deploySpec.Spec.Template.Spec.Volumes = append(deploySpec.Spec.Template.Spec.Volumes, corev1.Volume{
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
	})

	deploySpec.Spec.Template.Spec.Containers[0].VolumeMounts = append(deploySpec.Spec.Template.Spec.Containers[0].VolumeMounts, corev1.VolumeMount{
		MountPath: containerWorkingDir,
		Name:      serviceAccountTokenName,
	})

	deploySpec.Spec.Template.Spec.Containers = append(deploySpec.Spec.Template.Spec.Containers,
		getLocationContainer(),
		getStorageRestContainer(),
		getNormalizerContainer(),
	)

	return deploySpec
}

func TestPatchDeploymentSpec(t *testing.T) {
	tests := []struct {
		name       string
		deployment appsv1.Deployment
		expected   appsv1.Deployment
	}{
		{
			name:       "Test deployment patch",
			deployment: makeMockDeployment(),
			expected:   makeMockPatchedDeployment(),
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(tt *testing.T) {
			err := PatchDeploymentSpec(&test.deployment)

			if err != nil {
				tt.Fatalf("unexpected error: %v", err)
			} else if !reflect.DeepEqual(test.deployment, test.expected) {
				tt.Errorf("expected patch result to be equal to expected patch.\nexpected: %v\nactual: %v\n", test.expected, test.deployment)
			}
		})
	}
}

func TestUnpatchDeploymentSpec(t *testing.T) {
	tests := []struct {
		name       string
		deployment appsv1.Deployment
		expected   appsv1.Deployment
	}{
		{
			name:       "Test deployment patch",
			deployment: makeMockPatchedDeployment(),
			expected:   makeMockDeployment(),
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(tt *testing.T) {
			err := UnpatchDeploymentSpec(&test.deployment)

			if err != nil {
				tt.Fatalf("unexpected error: %v", err)
			} else if !reflect.DeepEqual(test.deployment, test.expected) {
				tt.Errorf("expected patch result to be equal to expected patch.\nexpected: %v\nactual: %v\n", test.expected, test.deployment)
			}
		})
	}
}
