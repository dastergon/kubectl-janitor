package cmd

import (
	"github.com/spf13/cobra"
	cmdutil "k8s.io/kubectl/pkg/cmd/util"
)

// NewJobsCommand provides the base command when called without any subcommands.
func newJobsCommand(factory cmdutil.Factory, options JanitorOptions) *cobra.Command {

	cmd := &cobra.Command{
		Use:          "jobs",
		Short:        "Find Jobs in a problematic state",
		SilenceUsage: true,
	}

	cmd.AddCommand(newFailedJobsCommand(factory, options))

	return cmd
}
