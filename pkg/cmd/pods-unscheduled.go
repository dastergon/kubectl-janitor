package cmd

import (
	"bytes"
	"context"
	"fmt"

	"github.com/spf13/cobra"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	cmdutil "k8s.io/kubectl/pkg/cmd/util"
)

// UnscheduledPodsOptions embeds JanitorOptions struct.
type UnscheduledPodsOptions struct {
	JanitorOptions
}

// newUnscheduledPodsOptions creates an instance of PendingPodsOptions.
func newUnscheduledPodsOptions(options JanitorOptions) *UnscheduledPodsOptions {
	return &UnscheduledPodsOptions{
		JanitorOptions: options,
	}
}

// newUnscheduledPodsCommand returns a cobra command wrapping PendingPodsOptions.
func newUnscheduledPodsCommand(factory cmdutil.Factory, options JanitorOptions) *cobra.Command {
	o := newUnscheduledPodsOptions(options)

	cmd := &cobra.Command{
		Use:          "unscheduled",
		Short:        "List Pods that are in a pending state (waiting to be scheduled)",
		SilenceUsage: true,
		RunE: func(c *cobra.Command, args []string) error {
			if err := o.Complete(factory, c); err != nil {
				return err
			}

			ctx := context.Background()
			noHeader := c.Flag("no-headers").Changed
			if err := o.Run(ctx, noHeader); err != nil {
				fmt.Fprintln(options.Streams.ErrOut, err.Error())
				return nil
			}
			return nil
		},
	}

	o.ResourceBuilderFlags.AddFlags(cmd.Flags())

	return cmd
}

// Run lists Pods waiting to be scheduled.
func (o *UnscheduledPodsOptions) Run(ctx context.Context, noHeader bool) error {
	client, err := o.GetClient()
	if err != nil {
		return err
	}

	options := metav1.ListOptions{}
	options.FieldSelector = "status.phase=Pending"

	pods, err := client.CoreV1().Pods(o.namespace).List(ctx, options)
	if err != nil {
		return err
	}

	var matrix [][]string

	for _, pod := range pods.Items {
		for _, c := range pod.Status.Conditions {
			if c.Type == "PodScheduled" && c.Status == "False" {
				age := getAge(pod.CreationTimestamp)
				row := []string{pod.Name, c.Reason, c.Message, age}
				if o.allNamespaces {
					row = append([]string{pod.Namespace}, row...)
				}
				matrix = append(matrix, row)
			}
		}
	}

	headers := []string{"NAME", "REASON", "MESSAGE", "AGE"}

	buf := bytes.NewBuffer(nil)
	writeResults(buf, headers, matrix, o.namespace, noHeader)
	fmt.Printf("%s", buf.String())

	return nil
}
