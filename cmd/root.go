package cmd

import (
	"os"

	"github.com/redhat-ai-dev/rhdh-ai-install/pkg/constants"
	"github.com/redhat-ai-dev/rhdh-ai-install/pkg/flags"
	"github.com/spf13/cobra"
)

type RootCmd struct {
	cmd   *cobra.Command
	flags *flags.RootFlags
}

func runRootCmd(cmd *cobra.Command, f *flags.RootFlags, args []string) error {
	return nil
}

func (r *RootCmd) runE(cmd *cobra.Command, args []string) error {
	return runRootCmd(cmd, r.flags, args)
}

func NewRootCmd() *RootCmd {
	rootCmd := &RootCmd{
		cmd: &cobra.Command{
			Use:   constants.AppName,
			Short: "CLI tool for managing Red Hat Developer Hub (RHDH) AI integrations",
			Long:  "CLI tool for managing AI integrations under given Red Hat Developer Hub (RHDH) instance",
		},
		flags: flags.NewRootFlags(),
	}
	rootCmd.flags.SetupFlags(rootCmd.cmd.Flags())
	rootCmd.cmd.RunE = rootCmd.runE
	return rootCmd
}

func (r *RootCmd) Cmd() *cobra.Command {
	r.flags.GetLogger(os.Stdout)

	return r.cmd
}
