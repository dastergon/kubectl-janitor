package cmd

import (
	"github.com/spf13/cobra"
	cmdutil "k8s.io/kubectl/pkg/cmd/util"
)

// newPVCsommand provides the base command when called without any subcommands.
func newPVCsCommand(factory cmdutil.Factory, options JanitorOptions) *cobra.Command {

	cmd := &cobra.Command{
		Use:          "pvcs",
		Short:        "Find PersistentVolumeClaims in a problematic state",
		SilenceUsage: true,
	}

	cmd.AddCommand(newPendingPVCsCommand(factory, options))

	return cmd
}
