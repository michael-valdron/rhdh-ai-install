package main

import (
	"os"

	"github.com/redhat-ai-dev/rhdh-ai-install/cmd"
)

func main() {
	rootCmd := cmd.NewRootCmd()
	if err := rootCmd.Cmd().Execute(); err != nil {
		os.Exit(1)
	}
}
