package cmd

import (
	"bytes"
	"context"
	"fmt"
	"strconv"

	"github.com/spf13/cobra"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	cmdutil "k8s.io/kubectl/pkg/cmd/util"
)

// StatusPodsOptions embeds JanitorOptions struct.
type StatusPodsOptions struct {
	JanitorOptions
}

// newStatusPodsOptions creates an instance of StatusPodsOptions.
func newStatusPodsOptions(options JanitorOptions) *StatusPodsOptions {
	return &StatusPodsOptions{
		JanitorOptions: options,
	}
}

// newStatusPodsCommand returns a cobra command wrapping StatusPodsOptions.
func newStatusPodsCommand(factory cmdutil.Factory, options JanitorOptions) *cobra.Command {
	o := newStatusPodsOptions(options)

	cmd := &cobra.Command{
		Use:          "status",
		Short:        "List the current statuses of the Pods and their respective count",
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

// Run lists statuses of the Pods.
func (o *StatusPodsOptions) Run(ctx context.Context, noHeader bool) error {
	client, err := o.GetClient()
	if err != nil {
		return err
	}

	pods, err := client.CoreV1().Pods(o.namespace).List(ctx, metav1.ListOptions{})
	if err != nil {
		return err
	}

	counter := make(map[string]map[string]int)
	for _, pod := range pods.Items {
		status := getPodStatus(pod)
		if counter[pod.Namespace] != nil {
			counter[pod.Namespace][status]++
		} else {
			counter[pod.Namespace] = make(map[string]int)
			counter[pod.Namespace][status] = 1
		}
	}

	var matrix [][]string

	for ns, status := range counter {
		for st, count := range status {
			row := []string{st, strconv.Itoa(count)}
			if o.allNamespaces {
				row = append([]string{ns}, row...)
			}
			matrix = append(matrix, row)
		}
	}

	headers := []string{"STATUS", "COUNT"}

	buf := bytes.NewBuffer(nil)
	writeResults(buf, headers, matrix, o.namespace, noHeader)
	fmt.Printf("%s", buf.String())

	return nil
}
