package config

// Function Path type
type FuncPath string

const (
	ModelCatalogCheck     FuncPath = "modelcatalog.check"     // Function path string of the model catalog check function
	ModelCatalogInstall   FuncPath = "modelcatalog.install"   // Function path string of the model catalog install function
	ModelCatalogUninstall FuncPath = "modelcatalog.uninstall" // Function path string of the model catalog uninstall function
)
