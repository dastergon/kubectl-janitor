package cmd

import (
	"github.com/spf13/cobra"
	cmdutil "k8s.io/kubectl/pkg/cmd/util"
)

// newSVCsCommand provides the base command when called without any subcommands.
func newSVCsCommand(factory cmdutil.Factory, options JanitorOptions) *cobra.Command {

	cmd := &cobra.Command{
		Use:          "svcs",
		Short:        "Find Services in a problematic state",
		SilenceUsage: true,
	}

	cmd.AddCommand(newNoEndpointsSVCsCommand(factory, options))

	return cmd
}
