package constants

import (
	"log/slog"
	"time"
)

const (
	AppName               = "rhdh-ai-install"        // Name of the application
	DefaultLogLevel       = slog.LevelInfo           // Default log level
	DefaultTimeout        = 10 * time.Minute         // Default timeout of executions
	DefaultIsOperator     = false                    // Default value if RHDH deployment is from the operator
	DefaultRHOAINamespace = "rhoai-model-registries" // Default namespace for RHOAI resources
)
