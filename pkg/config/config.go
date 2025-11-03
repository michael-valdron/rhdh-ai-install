package config

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"time"
)

type OptionsConfig struct {
	Debug    bool          `json:"debug,omitempty"`
	DryRun   bool          `json:"dry-run,omitempty"`
	LogLevel *slog.Level   `json:"log-level,omitempty"`
	Timeout  time.Duration `json:"timeout,omitempty"`
}

type DeveloperHubConfig struct {
	Namespace     string `json:"namespace,omitempty"`
	DeployName    string `json:"name,omitempty"`
	PluginsName   string `json:"plugins,omitempty"`
	AppConfigName string `json:"appconfig,omitempty"`
	IsOperator    bool   `json:"operator,omitempty"`
}

type LocationConfig struct {
	Image string `json:"image,omitempty"`
}

type StorageRestConfig struct {
	Image string `json:"image,omitempty"`
}

type NormalizerConfig struct {
	Image string `json:"image,omitempty"`
}

type ModelCatalogConfig struct {
	ModelRegistriesNamespace string             `json:"model-registries-namespace,omitempty"`
	ServiceAccountName       string             `json:"serviceaccount,omitempty"`
	DashboardRoleBindingName string             `json:"dashboard-role-binding,omitempty"`
	Location                 *LocationConfig    `json:"location,omitempty"`
	StorageRest              *StorageRestConfig `json:"storage-rest,omitempty"`
	Normalizer               *NormalizerConfig  `json:"normalizer,omitempty"`
}

type GlobalConfig struct {
	Options      *OptionsConfig     `json:"options,omitempty"`
	DeveloperHub DeveloperHubConfig `json:"developer-hub"`
}

type Config struct {
	Global       GlobalConfig       `json:"global"`
	ModelCatalog ModelCatalogConfig `json:"model-catalog"`
}

// ensures required global fields are set
func (c *Config) globalCheck() error {
	if c.Global.DeveloperHub.DeployName != "" {
		return fmt.Errorf("name of the developer hub deployment should be specified")
	} else if c.Global.DeveloperHub.AppConfigName != "" {
		return fmt.Errorf("appconfig configmap name of the developer hub deployment should be specified")
	} else if c.Global.DeveloperHub.PluginsName != "" {
		return fmt.Errorf("plugins configmap name of the developer hub deployment should be specified")
	}

	return nil
}

// ensures all required field for given library function path are set
func (c *Config) Check(funcPath FuncPath) error {
	if err := c.globalCheck(); err != nil {
		return err
	}
	switch funcPath {
	case ModelCatalogCheck:
		return nil
	case ModelCatalogInstall:
		return nil
	case ModelCatalogUninstall:
		return nil
	default:
		return fmt.Errorf("function path undefined")
	}
}

func (c *Config) Marshal() ([]byte, error) {
	return json.Marshal(c)
}

func (c *Config) Unmarshal(data []byte) error {
	return json.Unmarshal(data, c)
}
