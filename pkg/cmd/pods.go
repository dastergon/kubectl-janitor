package cmd

import (
	"github.com/spf13/cobra"
	cmdutil "k8s.io/kubectl/pkg/cmd/util"
)

// newPodsCommand provides the base command when called without any subcommands.
func newPodsCommand(factory cmdutil.Factory, options JanitorOptions) *cobra.Command {

	cmd := &cobra.Command{
		Use:          "pods",
		Short:        "Find Pods in a problematic state",
		SilenceUsage: true,
	}

	cmd.AddCommand(newUnhealthyPodsCommand(factory, options))
	cmd.AddCommand(newUnreadyPodsCommand(factory, options))
	cmd.AddCommand(newStatusPodsCommand(factory, options))
	cmd.AddCommand(newUnscheduledPodsCommand(factory, options))

	return cmd
}
