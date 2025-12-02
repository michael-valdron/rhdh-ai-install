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

func patchBackstageCR(ctx context.Context, c *client.Client, cfg config.Config) error {
	gvr := deployment.NewBackstageGVR()
	deploymentPatch, err := c.GetUnstructured(ctx, gvr, cfg.Global.DeveloperHub.DeployName, client.ClientParams{}, metav1.GetOptions{})
	if err != nil {
		return err
	}

	if err = deployment.PatchBackstageSpec(cfg.ModelCatalog, deploymentPatch); err != nil {
		return err
	}

	return c.PatchUnstructured(ctx, gvr, deploymentPatch, metav1.PatchOptions{})
}

func unpatchBackstageCR(ctx context.Context, c *client.Client, cfg config.Config) error {
	gvr := deployment.NewBackstageGVR()
	deploymentPatch, err := c.GetUnstructured(ctx, gvr, cfg.Global.DeveloperHub.DeployName, client.ClientParams{}, metav1.GetOptions{})
	if err != nil {
		return err
	}

	if err = deployment.UnpatchBackstageSpec(deploymentPatch); err != nil {
		return err
	}

	return c.PatchUnstructured(ctx, gvr, deploymentPatch, metav1.PatchOptions{})
}

func patchDeployment(ctx context.Context, c *client.Client, cfg config.Config) error {
	deploymentPatch, err := c.Deployments(client.ClientParams{}).Get(ctx, cfg.Global.DeveloperHub.DeployName, metav1.GetOptions{})
	if err != nil {
		return err
	}

	if err = deployment.PatchDeploymentSpec(cfg.ModelCatalog, deploymentPatch); err != nil {
		return err
	}

	return deployment.PatchDeployment(ctx, c, deploymentPatch, metav1.PatchOptions{})
}

func unpatchDeployment(ctx context.Context, c *client.Client, cfg config.Config) error {
	deploymentPatch, err := c.Deployments(client.ClientParams{}).Get(ctx, cfg.Global.DeveloperHub.DeployName, metav1.GetOptions{})
	if err != nil {
		return err
	}

	if err = deployment.UnpatchDeploymentSpec(deploymentPatch); err != nil {
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

// Patches given RHDH deployment with model catalog integration sidecars.
// Patches Backstage CR if targeting an RHDH operator install else patches the
// deployment spec instead.
func PatchDeployment(ctx context.Context, c *client.Client, cfg config.Config) error {
	if c.Options.RhdhOperator {
		return patchBackstageCR(ctx, c, cfg)
	} else {
		return patchDeployment(ctx, c, cfg)
	}
}

// Removes patches from the given RHDH deployment with model catalog integration sidecars.
// Removes patches from Backstage CR if targeting an RHDH operator install else removes from the
// deployment spec instead.
func UnpatchDeployment(ctx context.Context, c *client.Client, cfg config.Config) error {
	if c.Options.RhdhOperator {
		return unpatchBackstageCR(ctx, c, cfg)
	} else {
		return unpatchDeployment(ctx, c, cfg)
	}
}

// Checks if model catalog integration resources, plugins, and sidecars can be installed to a given RHDH deployment.
// Check ensures all required input fields are set, OpenShift cluster has the required operators installed, and does not
// already contain an installation or conflicts to the installation of the model catalog integration.
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

// Installs model catalog integration resources, plugins, and sidecars from a given RHDH deployment.
// Runs `Check` before installing.
func Install(ctx context.Context, c *client.Client, cfg config.Config) error {
	var rhdhNamespace string
	if err := Check(ctx, c, cfg); err != nil {
		return err
	}

	if cfg.Global.DeveloperHub.Namespace != "" {
		rhdhNamespace = cfg.Global.DeveloperHub.Namespace
	} else {
		rhdhNamespace = c.Namespace
	}

	if err := CreateServiceAccount(ctx, c, rhdhNamespace, config.GetModelRegistriesNamespace(&cfg)); err != nil {
		return err
	}

	if err := PatchDynamicPlugins(ctx, c, cfg.Global.DeveloperHub.PluginsName); err != nil {
		return err
	}

	if err := PatchDeployment(ctx, c, cfg); err != nil {
		return err
	}

	return nil
}

// Uninstalls model catalog integration resources, plugins, and sidecars from a given RHDH deployment.
func Uninstall(ctx context.Context, c *client.Client, cfg config.Config) error {
	var rhdhNamespace string

	if cfg.Global.DeveloperHub.Namespace != "" {
		rhdhNamespace = cfg.Global.DeveloperHub.Namespace
	} else {
		rhdhNamespace = c.Namespace
	}

	if err := UnpatchDeployment(ctx, c, cfg); err != nil {
		return err
	}

	if err := UnpatchDynamicPlugins(ctx, c, cfg.Global.DeveloperHub.PluginsName); err != nil {
		return err
	}

	if err := DeleteServiceAccount(ctx, c, rhdhNamespace, config.GetModelRegistriesNamespace(&cfg)); err != nil {
		return err
	}

	return nil
}
