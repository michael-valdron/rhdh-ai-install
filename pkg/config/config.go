package config

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"time"
)

// OptionsConfig holds general library configuration options (e.g., such as setting debugging mode, dry-run mode, log level, and timeout for operations).
type OptionsConfig struct {
	Debug    bool          `json:"debug,omitempty"`     // Enable debug mode if true; otherwise, false by default.
	DryRun   bool          `json:"dry-run,omitempty"`   // Run in dry-run mode if true; otherwise, normal execution mode by default.
	LogLevel *slog.Level   `json:"log-level,omitempty"` // Set the log level for logging messages (e.g., DEBUG, INFO, WARN, ERROR).
	Timeout  time.Duration `json:"timeout,omitempty"`   // Timeout duration for operations; if not specified, uses default timeout.
}

// DeveloperHubConfig contains configuration details for the target developer hub deployment.
type DeveloperHubConfig struct {
	Namespace     string `json:"namespace,omitempty"` // Namespace where the developer hub resources are deployed.
	DeployName    string `json:"name,omitempty"`      // Name of the developer hub deployment (Backstage CR name if IsOperator is true).
	PluginsName   string `json:"plugins,omitempty"`   // Name of the plugins configmap for the developer hub deployment.
	AppConfigName string `json:"appconfig,omitempty"` // Name of the appconfig configmap for the developer hub deployment.
	IsOperator    bool   `json:"operator,omitempty"`  // Indicates if the deployment is an operator-managed one; defaults to false.
}

// LocationConfig encapsulates configuration parameters for the location sidecar container, including its image.
// If not set in ModelCatalogConfig, it assumes default values; otherwise, the Image field overrides the default image to use.
type LocationConfig struct {
	Image      string `json:"image,omitempty"`       // Custom image for the location sidecar container; if omitted, defaults are used.
	Address    string `json:"address,omitempty"`     // The port the location service listens on; if omitted, defaults are used.
	StorageUrl string `json:"storage-url,omitempty"` // The REST endpoint of the storage container; if omitted, defaults are used.
	// Data format for the normalizer, can either be `JsonArrayFormat` for our new format from the schema folder, or the legacy
	// `CatalogInfoYamlFormat`; if omitted, defaults are used (`JsonArrayFormat`).
	NormalizerFormat string `json:"normalizer-format,omitempty"`
}

// StorageRestConfig encapsulates configuration parameters for the storage-rest sidecar container, including its image.
// If not set in ModelCatalogConfig, it assumes default values; otherwise, the Image field overrides the default image to use.
type StorageRestConfig struct {
	Image       string `json:"image,omitempty"`        // Custom image for the storage-rest sidecar container; if omitted, defaults are used.
	Address     string `json:"address,omitempty"`      // The port the storage service listens on; if omitted, defaults are used.
	LocationUrl string `json:"location-url,omitempty"` // The REST endpoint of our location container; if omitted, defaults are used.
	StorageType string `json:"storage-type,omitempty"` // The storage type used, for now only the development mode `ConfigMap` is supported; we'll add `GitHub` soon.
	// Data format for the normalizer, can either be `JsonArrayFormat` for our new format from the schema folder, or the legacy
	// `CatalogInfoYamlFormat`; if omitted, defaults are used (`JsonArrayFormat`).
	NormalizerFormat string `json:"normalizer-format,omitempty"`
}

// NormalizerConfig encapsulates configuration parameters for the normalizer sidecar container, including its image.
// If not set in ModelCatalogConfig, it assumes default values; otherwise, the Image field overrides the default image to use.
type NormalizerConfig struct {
	Image           string         `json:"image,omitempty"`            // Custom image for the normalizer sidecar container; if omitted, defaults are used.
	PprofAddress    string         `json:"pprof-address,omitempty"`    // The address the pprof endpoint binds to; if omitted, defaults are used.
	StorageUrl      string         `json:"storage-url,omitempty"`      // The REST endpoint of the storage container; if omitted, defaults are used.
	PollingInterval *time.Duration `json:"polling-interval,omitempty"` // The interval at which the RHOAI Model Registry REST endpoint is polled for updates; if omitted, defaults are used.
	// Data format for the normalizer, can either be `JsonArrayFormat` for our new format from the schema folder, or the legacy
	// `CatalogInfoYamlFormat`; if omitted, defaults are used (`JsonArrayFormat`).
	NormalizerFormat string `json:"normalizer-format,omitempty"`
}

// ModelCatalogConfig encapsulates all configuration parameters specific to the model catalog integration.
type ModelCatalogConfig struct {
	ModelRegistriesNamespace string             `json:"model-registries-namespace,omitempty"`
	ServiceAccountName       string             `json:"serviceaccount,omitempty"`
	DashboardRoleBindingName string             `json:"dashboard-role-binding,omitempty"`
	Location                 *LocationConfig    `json:"location,omitempty"`
	StorageRest              *StorageRestConfig `json:"storage-rest,omitempty"`
	Normalizer               *NormalizerConfig  `json:"normalizer,omitempty"`
}

// GlobalConfig aggregates all global configuration settings used across various integrations.
type GlobalConfig struct {
	Options      *OptionsConfig     `json:"options,omitempty"`
	DeveloperHub DeveloperHubConfig `json:"developer-hub"`
}

// Config is the main configuration structure that serves as a container for global and integration-specific configurations.
type Config struct {
	Global       GlobalConfig       `json:"global"`
	ModelCatalog ModelCatalogConfig `json:"model-catalog"`
}

// Ensures required global fields for a developer hub deployment are set.
func (c *Config) globalCheck() error {
	if c.Global.DeveloperHub.DeployName != "" {
		return fmt.Errorf("name of the developer hub deployment should be specified")
	}

	if c.Global.DeveloperHub.AppConfigName != "" {
		return fmt.Errorf("appconfig configmap name of the developer hub deployment should be specified")
	}

	if c.Global.DeveloperHub.PluginsName != "" {
		return fmt.Errorf("plugins configmap name of the developer hub deployment should be specified")
	}

	return nil
}

// Ensures all required fields for a given library function path are set.
func (c *Config) Check(funcPath FuncPath) error {
	if err := c.globalCheck(); err != nil {
		return err
	}

	switch funcPath {
	case ModelCatalogCheck:
		return nil // No additional checks required for this function path.
	case ModelCatalogInstall:
		return nil // No additional checks required for this function path.
	case ModelCatalogUninstall:
		return nil // No additional checks required for this function path.
	default:
		return fmt.Errorf("function path undefined")
	}
}

// Marshal converts the Config structure into JSON bytes.
func (c *Config) Marshal() ([]byte, error) {
	return json.Marshal(c)
}

// Unmarshal deserializes JSON bytes into the Config structure.
func (c *Config) Unmarshal(data []byte) error {
	return json.Unmarshal(data, c)
}
