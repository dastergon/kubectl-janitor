package cmd

import (
	"bytes"
	"context"
	"fmt"

	"github.com/spf13/cobra"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	cmdutil "k8s.io/kubectl/pkg/cmd/util"
)

// PendingPVCsOptions embeds JanitorOptions struct.
type PendingPVCsOptions struct {
	JanitorOptions
}

// newPendingPVCsOptions create a instance of PendingPVCsOptions.
func newPendingPVCsOptions(options JanitorOptions) *PendingPVCsOptions {
	return &PendingPVCsOptions{
		JanitorOptions: options,
	}
}

// newPendingPVCsCommand returns a cobra command wrapping PendingPVCsOptions.
func newPendingPVCsCommand(factory cmdutil.Factory, options JanitorOptions) *cobra.Command {
	o := newPendingPVCsOptions(options)

	cmd := &cobra.Command{
		Use:          "pending",
		Short:        "List PersistentVolumeClaims in a pending state (unbound)",
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

// Run lists PersistentVolumeClaims that are in a Pending state.
func (o *PendingPVCsOptions) Run(ctx context.Context, noHeader bool) error {
	client, err := o.GetClient()
	if err != nil {
		return err
	}

	pvcs, err := client.CoreV1().PersistentVolumeClaims(o.namespace).List(ctx, metav1.ListOptions{})
	if err != nil {
		return err
	}

	var matrix [][]string

	for _, pvc := range pvcs.Items {
		if pvc.Status.Phase == "Pending" {
			age := getAge(pvc.CreationTimestamp)
			row := []string{pvc.Name, age}
			if o.allNamespaces {
				row = append([]string{pvc.Namespace}, row...)
			}
			matrix = append(matrix, row)
		}
	}

	headers := []string{"NAME", "AGE"}

	buf := bytes.NewBuffer(nil)
	writeResults(buf, headers, matrix, o.namespace, noHeader)
	fmt.Printf("%s", buf.String())

	return nil
}
