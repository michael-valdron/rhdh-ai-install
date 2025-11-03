package modelcatalog

import (
	"context"
	"fmt"

	"github.com/redhat-ai-dev/rhdh-ai-install/pkg/client"
	"github.com/redhat-ai-dev/rhdh-ai-install/pkg/config"
	"github.com/redhat-ai-dev/rhdh-ai-install/pkg/modelcatalog/deployment"
	"github.com/redhat-ai-dev/rhdh-ai-install/pkg/modelcatalog/infra"
	"github.com/redhat-ai-dev/rhdh-ai-install/pkg/modelcatalog/plugins"
	"github.com/redhat-ai-dev/rhdh-ai-install/pkg/modelcatalog/serviceaccount"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
)

func patchBackstageCR(ctx context.Context, c *client.Client, name string) error {
	gvr := deployment.NewBackstageGVR()
	deploymentPatch, err := c.GetUnstructured(ctx, gvr, name, client.ClientParams{}, metav1.GetOptions{})
	if err != nil {
		return err
	}

	if err = deployment.PatchBackstageSpec(deploymentPatch); err != nil {
		return err
	}

	return c.PatchUnstructured(ctx, gvr, deploymentPatch, metav1.PatchOptions{})
}

func patchDeployment(ctx context.Context, c *client.Client, name string) error {
	deploymentPatch, err := c.Deployments(client.ClientParams{}).Get(ctx, name, metav1.GetOptions{})
	if err != nil {
		return err
	}

	if err = deployment.PatchDeploymentSpec(deploymentPatch); err != nil {
		return err
	}

	return deployment.PatchDeployment(ctx, c, deploymentPatch, metav1.PatchOptions{})
}

func CreateServiceAccount(ctx context.Context, c *client.Client, rhdhNamespace string, rhoaiNamespace string) error {
	content := serviceaccount.New(rhdhNamespace, rhoaiNamespace)

	return content.Create(ctx, c, metav1.CreateOptions{})
}

func DeleteServiceAccount(ctx context.Context, c *client.Client, rhdhNamespace string, rhoaiNamespace string) error {
	content := serviceaccount.New(rhdhNamespace, rhoaiNamespace)

	return content.Delete(ctx, c, metav1.DeleteOptions{})
}

func PatchDynamicPlugins(ctx context.Context, c *client.Client, configMapName string) error {
	dynamicPluginPatchList := plugins.GetDynamicPluginList()
	var (
		err                     error
		dynamicPluginsConfigMap *corev1.ConfigMap
		config                  *plugins.DynamicPluginConfig
	)

	if c.Options.Verbose {
		fmt.Println("Running: fetching current plugins list")
	}

	if dynamicPluginsConfigMap, err = c.ConfigMaps(client.ClientParams{}).Get(ctx, configMapName, metav1.GetOptions{}); err != nil {
		return err
	}

	if c.Options.Verbose {
		fmt.Println("Running: Preparing dynamic plugins list patch")
		fmt.Println("Decoding dynamic-plugins.yaml value")
	}

	if config, err = plugins.UnmarshalConfig(dynamicPluginsConfigMap); err != nil {
		return err
	}

	if c.Options.Verbose {
		fmt.Println("Adding model catalog dynamic plugins")
	}

	initialPluginsCount := config.Plugins.Count()
	config.Plugins = config.Plugins.Add(dynamicPluginPatchList)

	if c.Options.Verbose {
		fmt.Println("Checking if changes can be applied")
	}
	if initialPluginsCount != config.Plugins.Count() {
		if c.Options.Verbose {
			fmt.Println("Changes found")
			fmt.Println("Encoding changed dynamic-plugins.yaml value")
		}

		if err = config.MarshalConfig(dynamicPluginsConfigMap); err != nil {
			return err
		}

		if c.Options.Verbose {
			fmt.Println("Patch is ready")
			fmt.Println("Running: Running dynamic plugins list patch")
		}

		return plugins.PatchDynamicPlugins(ctx, c, dynamicPluginsConfigMap, metav1.PatchOptions{})
	} else {
		if c.Options.Verbose {
			fmt.Println("No changes found, skipping patch")
		}
		return nil
	}
}

func UnpatchDynamicPlugins(ctx context.Context, c *client.Client, configMapName string) error {
	dynamicPluginPatchList := plugins.GetDynamicPluginList()
	var (
		err                     error
		dynamicPluginsConfigMap *corev1.ConfigMap
		config                  *plugins.DynamicPluginConfig
	)

	if c.Options.Verbose {
		fmt.Println("Running: fetching current plugins list")
	}

	if dynamicPluginsConfigMap, err = c.ConfigMaps(client.ClientParams{}).Get(ctx, configMapName, metav1.GetOptions{}); err != nil {
		return err
	}

	if c.Options.Verbose {
		fmt.Println("Running: Preparing dynamic plugins list patch")
		fmt.Println("Decoding dynamic-plugins.yaml value")
	}

	if config, err = plugins.UnmarshalConfig(dynamicPluginsConfigMap); err != nil {
		return err
	}

	if c.Options.Verbose {
		fmt.Println("Removing model catalog dynamic plugins")
	}

	initialPluginsCount := config.Plugins.Count()
	config.Plugins = config.Plugins.Remove(dynamicPluginPatchList)

	if c.Options.Verbose {
		fmt.Println("Checking if changes can be applied")
	}
	if initialPluginsCount != config.Plugins.Count() {
		if c.Options.Verbose {
			fmt.Println("Changes found")
			fmt.Println("Encoding changed dynamic-plugins.yaml value")
		}

		if err = config.MarshalConfig(dynamicPluginsConfigMap); err != nil {
			return err
		}

		if c.Options.Verbose {
			fmt.Println("Patch is ready")
			fmt.Println("Running: Running dynamic plugins list patch")
		}

		return plugins.PatchDynamicPlugins(ctx, c, dynamicPluginsConfigMap, metav1.PatchOptions{})
	} else {
		if c.Options.Verbose {
			fmt.Println("No changes found, skipping patch")
		}
		return nil
	}
}

func PatchDeployment(ctx context.Context, c *client.Client, name string) error {
	if c.Options.RhdhOperator {
		return patchBackstageCR(ctx, c, name)
	} else {
		return patchDeployment(ctx, c, name)
	}
}

func Check(ctx context.Context, c *client.Client, cfg config.Config) error {
	var (
		modelRegistriesNamespace string                              = config.GetModelRegistriesNamespace(&cfg)
		serviceAccountResources  *serviceaccount.ServiceAccountSuite = serviceaccount.New(c.Namespace, modelRegistriesNamespace)
		deploymentPatch          *appsv1.Deployment
		backstageCRPatch         *unstructured.Unstructured
		err                      error
	)

	if err = cfg.Check(config.ModelCatalogCheck); err != nil {
		return err
	}

	if err = infra.Check(ctx, c); err != nil {
		return err
	}

	if err = serviceAccountResources.Check(ctx, c, metav1.GetOptions{}); err != nil {
		return err
	}

	if err = plugins.Check(ctx, c, cfg.Global.DeveloperHub.PluginsName); err != nil {
		return err
	}

	if c.Options.RhdhOperator {
		if backstageCRPatch, err = c.GetUnstructured(ctx, deployment.NewBackstageGVR(), cfg.Global.DeveloperHub.DeployName, client.ClientParams{}, metav1.GetOptions{}); err != nil {
			return err
		}

		if err = deployment.CheckBackstage(ctx, c, backstageCRPatch, metav1.GetOptions{}); err != nil {
			return err
		}
	} else {
		if deploymentPatch, err = c.Deployments(client.ClientParams{}).Get(ctx, cfg.Global.DeveloperHub.DeployName, metav1.GetOptions{}); err != nil {
			return err
		}

		if err = deployment.CheckDeployment(ctx, c, deploymentPatch, metav1.GetOptions{}); err != nil {
			return err
		}
	}

	return nil
}

func Install(ctx context.Context, c *client.Client, cfg config.Config) error {
	panic("unimplemented")
}

func Uninstall(ctx context.Context, c *client.Client, cfg config.Config) error {
	panic("unimplemented")
}
