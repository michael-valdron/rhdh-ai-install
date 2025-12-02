package config

import "time"

// Function Path type
type FuncPath string

const (
	// ModelCatalog

	DefaultLocationImage    string = "quay.io/redhat-ai-dev/model-catalog-location-service:latest" // Default location container image
	DefaultStorageRestImage string = "quay.io/redhat-ai-dev/model-catalog-storage-rest:latest"     // Default storage REST container image
	DefaultNormalizerImage  string = "quay.io/redhat-ai-dev/model-catalog-rhoai-normalizer:latest" // Default normalizer container image

	DefaultNormalizerFormat string        = "JsonArrayFormat"       // Default normalizer data format
	DefaultStorageUrl       string        = "localhost:7070"        // Default storage container URL
	DefaultPollingInterval  time.Duration = 2 * time.Minute         // Default polling interval
	DefaultLocationUrl      string        = "http://localhost:9090" // Default location container URL
	DefaultStorageType      string        = "ConfigMap"             // Default storage type

	// FunPath

	ModelCatalogCheck     FuncPath = "modelcatalog.check"     // Function path string of the model catalog check function
	ModelCatalogInstall   FuncPath = "modelcatalog.install"   // Function path string of the model catalog install function
	ModelCatalogUninstall FuncPath = "modelcatalog.uninstall" // Function path string of the model catalog uninstall function
)
