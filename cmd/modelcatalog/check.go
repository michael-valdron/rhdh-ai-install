package modelcatalog

import (
	"os"

	"github.com/redhat-ai-dev/rhdh-ai-install/pkg/flags"
	"github.com/spf13/cobra"
)

type CheckCmd struct {
	cmd   *cobra.Command
	flags *flags.CheckFlags
}

func runCheckCmd(cmd *cobra.Command, f *flags.CheckFlags, args []string) error {
	return nil
}

func (c *CheckCmd) runE(cmd *cobra.Command, args []string) error {
	return runCheckCmd(cmd, c.flags, args)
}

func NewCheckCmd() *CheckCmd {
	checkCmd := &CheckCmd{
		cmd: &cobra.Command{
			Use:   "check",
			Short: "Checks readiness for the model catalog bridge RHDH integration",
			Long:  "Checks readiness of OpenShift cluster to ensure all component and resource requires are available for the model catalog bridge RHDH integration",
		},
		flags: flags.NewCheckFlags(),
	}
	checkCmd.flags.SetupFlags(checkCmd.cmd.Flags())
	checkCmd.cmd.RunE = checkCmd.runE
	return checkCmd
}

func (c *CheckCmd) Cmd() *cobra.Command {
	c.flags.GetLogger(os.Stdout)

	return c.cmd
}
