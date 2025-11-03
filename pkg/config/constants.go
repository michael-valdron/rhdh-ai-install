package config

// Function Path type
type FuncPath string

const (
	ModelCatalogCheck     FuncPath = "modelcatalog.check"
	ModelCatalogInstall   FuncPath = "modelcatalog.install"
	ModelCatalogUninstall FuncPath = "modelcatalog.uninstall"
)
