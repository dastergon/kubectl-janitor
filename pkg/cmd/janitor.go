package cmd

import (
	"github.com/spf13/cobra"
	cmdutil "k8s.io/kubectl/pkg/cmd/util"
)

var cmdExample = `# List Pods that are in a pending state (waiting to be scheduled)
kubectl janitor pods unscheduled

# List Pods in an unhealthy state.
kubectl janitor pods unhealthy

# List Pods that are currently in a running phase but not ready for some reason.
kubectl janitor pods unready

# List the current statuses of the Pods and their respective count.
kubectl janitor pods status

# List Jobs that have failed to run and have restartPolicy: Never.
kubectl janitor jobs failed

# List PesistentVolumes that are available for claim.
kubectl janitor pvs unclaimed

# List PersistentVolumeClaims in an pending state (unbound).
kubectl janitor pvcs pending
`

// NewJanitorCommand provides the base command when called without any subcommands.
func NewJanitorCommand() *cobra.Command {
	o := NewJanitorOptions()

	cmd := &cobra.Command{
		Use:          "janitor",
		Example:      cmdExample,
		Short:        "Find objects in a problematic state in your Kubernetes cluster",
		SilenceUsage: true,
	}

	cmd.PersistentFlags().Bool("no-headers", false, "Don't print headers (default print headers).")

	flags := cmd.PersistentFlags()
	o.ConfigFlags.AddFlags(flags)

	matchVersionFlags := cmdutil.NewMatchVersionFlags(o.ConfigFlags)
	matchVersionFlags.AddFlags(flags)

	f := cmdutil.NewFactory(matchVersionFlags)

	cmd.AddCommand(newJobsCommand(f, o))
	cmd.AddCommand(newPodsCommand(f, o))
	cmd.AddCommand(newPVCsCommand(f, o))
	cmd.AddCommand(newPVsCommand(f, o))

	return cmd
}
