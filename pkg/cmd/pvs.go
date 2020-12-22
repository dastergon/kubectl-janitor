package cmd

import (
	"github.com/spf13/cobra"
	cmdutil "k8s.io/kubectl/pkg/cmd/util"
)

// newPVsommand provides the base command when called without any subcommands.
func newPVsCommand(factory cmdutil.Factory, options JanitorOptions) *cobra.Command {

	cmd := &cobra.Command{
		Use:          "pvs",
		Short:        "Find PersistentVolumes in a problematic state",
		SilenceUsage: true,
	}

	cmd.AddCommand(newUnclaimedPVsCommand(factory, options))

	return cmd
}
