package config

import (
	"github.com/redhat-ai-dev/rhdh-ai-install/pkg/constants"
)

func GetModelRegistriesNamespace(cfg *Config) string {
	if cfg.ModelCatalog.ModelRegistriesNamespace == "" {
		return constants.DefaultRHOAINamespace
	}

	return cfg.ModelCatalog.ModelRegistriesNamespace
}
