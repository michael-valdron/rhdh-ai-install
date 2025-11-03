package flags

import (
	"fmt"
	"io"
	"log/slog"
	"strings"
	"time"

	"github.com/redhat-ai-dev/rhdh-ai-install/pkg/constants"
	"github.com/spf13/pflag"
)

type RootFlags struct {
	ConfigFile    string
	Debug         bool
	DryRun        bool
	LogLevel      *slog.Level
	Timeout       time.Duration
	Namespace     string
	DeployName    string
	PluginsName   string
	AppConfigName string
	IsOperator    bool
}

type CheckFlags struct {
	RootFlags
}

func NewRootFlags() *RootFlags {
	logLevel := defaultLogLevel
	return &RootFlags{
		Debug:    false,
		DryRun:   false,
		LogLevel: &logLevel,
		Timeout:  defaultTimeout,
	}
}

func NewCheckFlags() *CheckFlags {
	return &CheckFlags{
		RootFlags: *NewRootFlags(),
	}
}

func (f *RootFlags) SetupFlags(flagParams *pflag.FlagSet) {
	flagParams.StringVarP(&f.ConfigFile, "config-file", "f", "", "specify parameters config file to use with commands")
	flagParams.BoolVar(&f.Debug, "debug", f.Debug, "enable debug mode")
	flagParams.BoolVar(&f.DryRun, "dry-run", f.DryRun, "enable dry-run mode")
	flagParams.Var(NewLogLevel(f.LogLevel),
		"log-level",
		fmt.Sprintf("logging verbosity level (default %q)", strings.ToLower(f.LogLevel.String())))
	flagParams.VarP(NewDuration(&f.Timeout),
		"timeout",
		"t",
		fmt.Sprintf("client timeout duration (default %s)", f.Timeout.String()))
	flagParams.StringVarP(&f.Namespace, "rhdh-namespace", "n", "", "namespace of the target RHDH instance")
	flagParams.StringVarP(&f.DeployName, "rhdh-instance", "i", "", "target RHDH instance name/backstage CR name")
	flagParams.StringVarP(&f.PluginsName, "rhdh-dynamic-plugins", "p", "", "target dynamic plugins ConfigMap name")
	flagParams.StringVarP(&f.AppConfigName, "rhdh-appconfig", "a", "", "target AppConfig ConfigMap name")
	flagParams.BoolVar(&f.IsOperator, "rhdh-operator", constants.DefaultIsOperator, "indicates whether target RHDH instance is deployed using the RHDH operator")
}

func (f *RootFlags) GetLogger(out io.Writer) *slog.Logger {
	logOpts := &slog.HandlerOptions{Level: f.LogLevel}
	return slog.New(slog.NewTextHandler(out, logOpts))
}
