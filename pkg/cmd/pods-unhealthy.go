package cmd

import (
	"bytes"
	"context"
	"fmt"

	"github.com/spf13/cobra"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	cmdutil "k8s.io/kubectl/pkg/cmd/util"
)

// UnhealthyPodsOptions embeds JanitorOptions struct.
type UnhealthyPodsOptions struct {
	JanitorOptions
}

// newUnhealthyPodsOptions creates an instance of UnhealthyPodsOptions.
func newUnhealthyPodsOptions(options JanitorOptions) *UnhealthyPodsOptions {
	return &UnhealthyPodsOptions{
		JanitorOptions: options,
	}
}

// newUnreadyPodsCommand returns a cobra command wrapping UnhealthyPodsOptions.
func newUnhealthyPodsCommand(factory cmdutil.Factory, options JanitorOptions) *cobra.Command {
	o := newUnhealthyPodsOptions(options)

	cmd := &cobra.Command{
		Use:          "unhealthy",
		Short:        "List Pods in an unhealthy state",
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

// Run finds pods that are unhealthy.
func (o *UnhealthyPodsOptions) Run(ctx context.Context, noHeader bool) error {
	client, err := o.GetClient()
	if err != nil {
		return err
	}

	pods, err := client.CoreV1().Pods(o.namespace).List(ctx, metav1.ListOptions{})
	if err != nil {
		return err
	}

	var matrix [][]string

	for _, pod := range pods.Items {
		if !isPodHealthy(pod) {
			age := getAge(pod.CreationTimestamp)
			podStatus := getPodStatus(pod)
			row := []string{pod.Name, podStatus, age}
			if o.allNamespaces {
				row = append([]string{pod.Namespace}, row...)
			}
			matrix = append(matrix, row)
		}
	}

	headers := []string{"NAME", "STATUS", "AGE"}

	buf := bytes.NewBuffer(nil)
	writeResults(buf, headers, matrix, o.namespace, noHeader)
	fmt.Printf("%s", buf.String())

	return nil
}
