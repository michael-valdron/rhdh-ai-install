package infra

import (
	"context"
	"fmt"
	"slices"

	"github.com/redhat-ai-dev/rhdh-ai-install/pkg/client"
	"github.com/redhat-ai-dev/rhdh-ai-install/pkg/util"
	appsv1 "k8s.io/api/apps/v1"
	apiextv1 "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func getCRDs(ctx context.Context, c *client.Client) (*apiextv1.CustomResourceDefinitionList, error) {
	if crds, err := c.CustomResourceDefinitions().List(ctx, metav1.ListOptions{}); err != nil {
		return nil, fmt.Errorf("error occurred trying to fetch the CRDs: %v", err)
	} else {
		return crds, nil
	}
}

func findCRDFunc(kind string, group string) func(apiextv1.CustomResourceDefinition) bool {
	return func(crd apiextv1.CustomResourceDefinition) bool {
		return crd.Kind == kind && crd.GroupVersionKind().Group == group
	}
}

func checkForOpenShiftAI(ctx context.Context, c *client.Client, crds *apiextv1.CustomResourceDefinitionList) error {
	var (
		err        error
		deployment *appsv1.Deployment
	)

	if deployment, err = c.Deployments(client.ClientParams{Namespace: defaultRHOAIOperatorNamespace}).Get(ctx, rhoaiOperatorControllerName, metav1.GetOptions{}); err != nil {
		if errors.IsNotFound(err) {
			return fmt.Errorf("the OpenShift AI component appears to be not fully installed, missing controller: %v", err)
		} else {
			return fmt.Errorf("unexpected error trying to fetch the OpenShift AI operator controller deployment")
		}
	}

	if deployment.Status.ReadyReplicas == 0 {
		return fmt.Errorf("the OpenShift AI operator controller is not ready")
	}

	if crds == nil {
		if crds, err = getCRDs(ctx, c); err != nil {
			return err
		}
	}

	if !slices.ContainsFunc(crds.Items, findCRDFunc(dataScienceClusterKind, dataScienceClusterGroup)) {
		return fmt.Errorf("missing component required %s CRD, please finish installing or reinstall OpenShift AI required component", dataScienceClusterKind)
	}

	return nil
}

func checkForNvidiaGPU(ctx context.Context, c *client.Client, crds *apiextv1.CustomResourceDefinitionList) error {
	var (
		err        error
		deployment *appsv1.Deployment
	)

	if deployment, err = c.Deployments(client.ClientParams{Namespace: defaultNvidiaGPUOperatorNamespace}).Get(ctx, nvidiaGPUOperatorControllerName, metav1.GetOptions{}); err != nil {
		if errors.IsNotFound(err) {
			return fmt.Errorf("the Nvidia GPU component appears to be not fully installed, missing controller: %v", err)
		} else {
			return fmt.Errorf("unexpected error trying to fetch the Nvidia GPU operator controller deployment")
		}
	}

	if deployment.Status.ReadyReplicas == 0 {
		return fmt.Errorf("the Nvidia GPU operator controller is not ready")
	}

	if crds == nil {
		if crds, err = getCRDs(ctx, c); err != nil {
			return err
		}
	}

	if !slices.ContainsFunc(crds.Items, findCRDFunc(clusterPolicyKind, clusterPolicyGroup)) {
		return fmt.Errorf("missing model catalog required %s CRD, please finish installing or reinstall Nvidia GPU required component", clusterPolicyKind)
	}

	return nil
}

func checkForNodeFeatureDiscovery(ctx context.Context, c *client.Client, crds *apiextv1.CustomResourceDefinitionList) error {
	var (
		err        error
		deployment *appsv1.Deployment
	)

	if deployment, err = c.Deployments(client.ClientParams{Namespace: defaultNFDOperatorNamespace}).Get(ctx, nfdOperatorControllerName, metav1.GetOptions{}); err != nil {
		if errors.IsNotFound(err) {
			return fmt.Errorf("the Node Feature Discovery component appears to be not fully installed, missing controller: %v", err)
		} else {
			return fmt.Errorf("unexpected error trying to fetch the Node Feature Discovery operator controller deployment")
		}
	}

	if deployment.Status.ReadyReplicas == 0 {
		return fmt.Errorf("the Node Feature Discovery operator controller is not ready")
	}

	if crds == nil {
		if crds, err = getCRDs(ctx, c); err != nil {
			return err
		}
	}

	if !slices.ContainsFunc(crds.Items, findCRDFunc(nfdKind, nfdGroup)) {
		return fmt.Errorf("missing model catalog required %s CRD, please finish installing or reinstall Node Feature Discovery required component", nfdKind)
	}

	return nil
}

func checkForKMM(ctx context.Context, c *client.Client) error {
	if deployment, err := c.Deployments(client.ClientParams{Namespace: defaultKMMOperatorNamespace}).Get(ctx, kmmOperatorControllerName, metav1.GetOptions{}); err != nil {
		if errors.IsNotFound(err) {
			return fmt.Errorf("the Kernel Module Management (KMM) component appears to be not fully installed, missing controller: %v", err)
		} else {
			return fmt.Errorf("unexpected error trying to fetch the Kernel Module Management (KMM) operator controller deployment")
		}
	} else if deployment.Status.ReadyReplicas == 0 {
		return fmt.Errorf("the Kernel Module Management (KMM) operator controller is not ready")
	}

	if deployment, err := c.Deployments(client.ClientParams{Namespace: defaultKMMOperatorNamespace}).Get(ctx, kmmOperatorWebhookName, metav1.GetOptions{}); err != nil {
		if errors.IsNotFound(err) {
			return fmt.Errorf("the Kernel Module Management (KMM) component appears to be not fully installed, missing webhook: %v", err)
		} else {
			return fmt.Errorf("unexpected error trying to fetch the Kernel Module Management (KMM) operator webhook deployment")
		}
	} else if deployment.Status.ReadyReplicas == 0 {
		return fmt.Errorf("the Kernel Module Management (KMM) operator webhook is not ready")
	}

	return nil
}

func checkForAuthorino(ctx context.Context, c *client.Client) error {
	if deployment, err := c.Deployments(client.ClientParams{Namespace: defaultAuthorinoOperatorNamespace}).Get(ctx, authorinoOperatorControllerName, metav1.GetOptions{}); err != nil {
		if errors.IsNotFound(err) {
			return fmt.Errorf("the Authorino component appears to be not fully installed, missing controller: %v", err)
		} else {
			return fmt.Errorf("unexpected error trying to fetch the Authorino operator controller deployment")
		}
	} else if deployment.Status.ReadyReplicas == 0 {
		return fmt.Errorf("the Authorino operator controller is not ready")
	}

	return nil
}

func Check(ctx context.Context, c *client.Client) error {
	errs := []any{}
	var (
		crds *apiextv1.CustomResourceDefinitionList
		err  error
	)

	// Get CRD list
	if crds, err = getCRDs(ctx, c); err != nil {
		return err
	}

	// Check for OpenShift AI
	errs = append(errs, checkForOpenShiftAI(ctx, c, crds))

	// Check for Node Feature Discovery
	errs = append(errs, checkForNodeFeatureDiscovery(ctx, c, crds))

	// Check for Nvidia GPU
	errs = append(errs, checkForNvidiaGPU(ctx, c, crds))

	// Check for KMM
	errs = append(errs, checkForKMM(ctx, c))

	// Check for Authorino
	errs = append(errs, checkForAuthorino(ctx, c))

	if len(errs) != 0 {
		if c.Options.Verbose {
			return fmt.Errorf("model catalog infrastructure resources are not ready, see the following errors (note: remember uninstall first if there is a previous install):\n%s", util.Join(errs, "\n"))
		} else {
			return fmt.Errorf("model catalog infrastructure resources are not ready")
		}
	} else {
		return nil
	}
}
