package cmd

import (
	"bytes"
	"context"
	"fmt"

	"github.com/spf13/cobra"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	cmdutil "k8s.io/kubectl/pkg/cmd/util"
)

// UnreadyPodsOptions embeds JanitorOptions struct.
type UnreadyPodsOptions struct {
	JanitorOptions
}

// newUnreadyPodsOptions creates an instance of UnreadyPodsOptions.
func newUnreadyPodsOptions(options JanitorOptions) *UnreadyPodsOptions {
	return &UnreadyPodsOptions{
		JanitorOptions: options,
	}
}

// newUnreadyPodsCommand returns a cobra command wrapping UnreadyPodsOptions.
func newUnreadyPodsCommand(factory cmdutil.Factory, options JanitorOptions) *cobra.Command {
	o := newUnreadyPodsOptions(options)

	cmd := &cobra.Command{
		Use:          "unready",
		Short:        "List Pods that are currently in a running phase but not ready for some reason",
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

// Run finds pods in a not ready mode.
func (o *UnreadyPodsOptions) Run(ctx context.Context, noHeader bool) error {
	client, err := o.GetClient()
	if err != nil {
		return err
	}

	options := metav1.ListOptions{}
	options.FieldSelector = "status.phase=Running"

	pods, err := client.CoreV1().Pods(o.namespace).List(ctx, options)
	if err != nil {
		return err
	}

	var matrix [][]string

	for _, pod := range pods.Items {
		if !isPodReady(pod) {
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
