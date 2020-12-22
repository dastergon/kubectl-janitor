package main

import (
	"os"

	"github.com/dastergon/kubectl-janitor/pkg/cmd"
	"github.com/spf13/pflag"

	_ "k8s.io/client-go/plugin/pkg/client/auth"
)

// version is populated by goreleaser
var version string

// Execute consolidates all sub-commands to the root command and sets the flags.
func Execute() error {
	root := cmd.NewJanitorCommand()
	root.Version = version
	return root.Execute()
}

func main() {
	flags := pflag.NewFlagSet("kubectl-janitor", pflag.ExitOnError)
	pflag.CommandLine = flags

	if err := Execute(); err != nil {
		os.Exit(1)
	}
}
