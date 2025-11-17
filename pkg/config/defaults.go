package config

import (
	"github.com/redhat-ai-dev/rhdh-ai-install/pkg/constants"
)

// GetModelRegistriesNamespace retrieves the ModelRegistriesNamespace from the provided Config.
// If not set, it defaults to the RHOAI namespace defined in constants.DefaultRHOAINamespace.
func GetModelRegistriesNamespace(cfg *Config) string {
	if cfg.ModelCatalog.ModelRegistriesNamespace == "" {
		return constants.DefaultRHOAINamespace
	}

	return cfg.ModelCatalog.ModelRegistriesNamespace
}
