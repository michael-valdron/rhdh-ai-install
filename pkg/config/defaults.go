package config

import (
	"github.com/redhat-ai-dev/rhdh-ai-install/pkg/constants"
	corev1 "k8s.io/api/core/v1"
)

// GetModelRegistriesNamespace retrieves the ModelRegistriesNamespace from the provided Config.
// If not set, it defaults to the RHOAI namespace defined in constants.DefaultRHOAINamespace.
func GetModelRegistriesNamespace(cfg *Config) string {
	if cfg.ModelCatalog.ModelRegistriesNamespace == "" {
		return constants.DefaultRHOAINamespace
	}

	return cfg.ModelCatalog.ModelRegistriesNamespace
}

// GetLocationImage retrieves the model catalog location container image from the provided Config.
// If not set, it defaults to the default location image defined in config.DefaultLocationImage.
func GetLocationImage(cfg *LocationConfig) string {
	if cfg == nil || cfg.Image == "" {
		return DefaultLocationImage
	}

	return cfg.Image
}

// GetStorageRestImage retrieves the model catalog storage rest container image from the provided Config.
// If not set, it defaults to the default storage rest image defined in config.DefaultStorageRestImage.
func GetStorageRestImage(cfg *StorageRestConfig) string {
	if cfg == nil || cfg.Image == "" {
		return DefaultStorageRestImage
	}

	return cfg.Image
}

// GetNormalizerImage retrieves the model catalog normalizer container image from the provided Config.
// If not set, it defaults to the default normalizer image defined in config.DefaultNormalizerImage.
func GetNormalizerImage(cfg *NormalizerConfig) string {
	if cfg == nil || cfg.Image == "" {
		return DefaultNormalizerImage
	}

	return cfg.Image
}

func GetLocationEnvVars(cfg *LocationConfig) []corev1.EnvVar {
	envVars := []corev1.EnvVar{}

	if cfg != nil && cfg.NormalizerFormat != "" {
		envVars = append(envVars, corev1.EnvVar{
			Name:  constants.NormalizerFormatVarName,
			Value: cfg.NormalizerFormat,
		})
	} else {
		envVars = append(envVars, corev1.EnvVar{
			Name:  constants.NormalizerFormatVarName,
			Value: DefaultNormalizerFormat,
		})
	}

	if cfg != nil && cfg.StorageUrl != "" {
		envVars = append(envVars, corev1.EnvVar{
			Name:  constants.StorageUrlVarName,
			Value: cfg.StorageUrl,
		})
	} else {
		envVars = append(envVars, corev1.EnvVar{
			Name:  constants.StorageUrlVarName,
			Value: DefaultStorageUrl,
		})
	}

	return envVars
}

func GetStorageRestEnvVars(cfg *StorageRestConfig) []corev1.EnvVar {
	envVars := []corev1.EnvVar{}

	if cfg != nil && cfg.NormalizerFormat != "" {
		envVars = append(envVars, corev1.EnvVar{
			Name:  constants.NormalizerFormatVarName,
			Value: cfg.NormalizerFormat,
		})
	} else {
		envVars = append(envVars, corev1.EnvVar{
			Name:  constants.NormalizerFormatVarName,
			Value: DefaultNormalizerFormat,
		})
	}

	if cfg != nil && cfg.LocationUrl != "" {
		envVars = append(envVars, corev1.EnvVar{
			Name:  constants.LocationUrlVarName,
			Value: cfg.LocationUrl,
		})
	} else {
		envVars = append(envVars, corev1.EnvVar{
			Name:  constants.LocationUrlVarName,
			Value: DefaultLocationUrl,
		})
	}

	if cfg != nil && cfg.StorageType != "" {
		envVars = append(envVars, corev1.EnvVar{
			Name:  constants.StorageTypeVarName,
			Value: cfg.StorageType,
		})
	} else {
		envVars = append(envVars, corev1.EnvVar{
			Name:  constants.StorageTypeVarName,
			Value: DefaultStorageType,
		})
	}

	return envVars
}

func GetNormalizerEnvVars(cfg *NormalizerConfig) []corev1.EnvVar {
	envVars := []corev1.EnvVar{}

	if cfg != nil && cfg.NormalizerFormat != "" {
		envVars = append(envVars, corev1.EnvVar{
			Name:  constants.NormalizerFormatVarName,
			Value: cfg.NormalizerFormat,
		})
	} else {
		envVars = append(envVars, corev1.EnvVar{
			Name:  constants.NormalizerFormatVarName,
			Value: DefaultNormalizerFormat,
		})
	}

	if cfg != nil && cfg.StorageUrl != "" {
		envVars = append(envVars, corev1.EnvVar{
			Name:  constants.StorageUrlVarName,
			Value: cfg.StorageUrl,
		})
	} else {
		envVars = append(envVars, corev1.EnvVar{
			Name:  constants.StorageUrlVarName,
			Value: DefaultStorageUrl,
		})
	}

	if cfg != nil && cfg.PollingInterval != nil {
		envVars = append(envVars, corev1.EnvVar{
			Name:  constants.PollingIntervalVarName,
			Value: cfg.PollingInterval.String(),
		})
	} else {
		envVars = append(envVars, corev1.EnvVar{
			Name:  constants.PollingIntervalVarName,
			Value: DefaultPollingInterval.String(),
		})
	}

	return envVars
}
