package plugins

import (
	"context"
	"fmt"
	"slices"

	"github.com/redhat-ai-dev/rhdh-ai-install/pkg/client"
	"github.com/redhat-ai-dev/rhdh-ai-install/pkg/util"
	"go.yaml.in/yaml/v3"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
)

type DynamicPlugin struct {
	Disabled bool   `json:"disabled,omitempty"`
	Package  string `json:"package,omitempty"`
}

type DynamicPluginList []DynamicPlugin

type DynamicPluginConfig struct {
	Includes []string          `json:"includes"`
	Plugins  DynamicPluginList `json:"plugins"`
}

func (dpl DynamicPluginList) Count() int {
	return len(dpl)
}

func (dpl DynamicPluginList) Add(al DynamicPluginList) DynamicPluginList {
	for _, a := range al {
		doesContain := slices.ContainsFunc(dpl, func(di DynamicPlugin) bool { return a.Package == di.Package })

		if !doesContain {
			dpl = append(dpl, a)
		}
	}

	return dpl
}

func (dpl DynamicPluginList) Remove(rl DynamicPluginList) DynamicPluginList {
	var result []DynamicPlugin = dpl

	for _, r := range rl {
		result, _ = util.RemoveFromSliceByFunc(result, getDynamicPluginByPackageFunc(r.Package))
	}

	return result
}

func (c *DynamicPluginConfig) MarshalConfig(configMap *corev1.ConfigMap) error {
	var (
		data []byte
		err  error
	)

	if data, err = yaml.Marshal(c); err != nil {
		return err
	}

	configMap.Data[DefaultDynamicPluginsKey] = string(data)

	return nil
}

func getDynamicPluginByPackageFunc(value string) func(DynamicPlugin) bool {
	return func(dp DynamicPlugin) bool { return dp.Package == value }
}

func UnmarshalConfig(configMap *corev1.ConfigMap) (*DynamicPluginConfig, error) {
	data := []byte(configMap.Data[DefaultDynamicPluginsKey])
	var config DynamicPluginConfig

	if err := yaml.Unmarshal(data, &config); err != nil {
		return nil, err
	}

	return &config, nil
}

func GetDynamicPluginList() DynamicPluginList {
	return DynamicPluginList{
		{
			Disabled: false,
			Package:  DynamicPluginModelCatalog,
		},
		{
			Disabled: false,
			Package:  DynamicPluginTechdocUrlReader,
		},
	}
}

func PatchDynamicPlugins(ctx context.Context, c *client.Client, patchedConfigMap *corev1.ConfigMap, opts metav1.PatchOptions) error {
	var (
		patchData []byte
		err       error
	)

	if patchData, err = patchedConfigMap.Marshal(); err != nil {
		return err
	}

	_, err = c.ConfigMaps(client.ClientParams{Namespace: patchedConfigMap.Namespace}).Patch(ctx, patchedConfigMap.Name, types.StrategicMergePatchType, patchData, opts)

	return err
}

func Check(ctx context.Context, c *client.Client, configMapName string) error {
	errs := []any{}

	if _, err := c.ConfigMaps(client.ClientParams{}).Get(ctx, configMapName, metav1.GetOptions{}); err != nil {
		errs = append(errs, err)
	}

	if len(errs) != 0 {
		if c.Options.Verbose {
			return fmt.Errorf("model catalog plugin resources are not ready, see the following errors (note: remember uninstall first if there is a previous install):\n%s", util.Join(errs, "\n"))
		} else {
			return fmt.Errorf("model catalog plugin resources are not ready")
		}
	} else {
		return nil
	}
}
