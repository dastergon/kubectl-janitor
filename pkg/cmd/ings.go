package cmd

import (
	"github.com/spf13/cobra"
	cmdutil "k8s.io/kubectl/pkg/cmd/util"
)

// newINGsCommand provides the base command when called without any subcommands.
func newINGsCommand(factory cmdutil.Factory, options JanitorOptions) *cobra.Command {

	cmd := &cobra.Command{
		Use:          "ings",
		Short:        "Find Ingresses in a problematic state",
		SilenceUsage: true,
	}

	cmd.AddCommand(newNoServiceINGsCommand(factory, options))

	return cmd
}
