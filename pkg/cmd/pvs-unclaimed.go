package cmd

import (
	"bytes"
	"context"
	"fmt"

	"github.com/spf13/cobra"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	cmdutil "k8s.io/kubectl/pkg/cmd/util"
)

// UnclaimedPVsOptions embeds JanitorOptions struct.
type UnclaimedPVsOptions struct {
	JanitorOptions
}

// newUnclaimedPVsOptions create a instance of UmclaimedPVOptions.
func newUnclaimedPVsOptions(options JanitorOptions) *UnclaimedPVsOptions {
	return &UnclaimedPVsOptions{
		JanitorOptions: options,
	}
}

// newUnclaimedPVsCommand returns a cobra command wrapping UnclaimedPVsOptions.
func newUnclaimedPVsCommand(factory cmdutil.Factory, options JanitorOptions) *cobra.Command {
	o := newUnclaimedPVsOptions(options)

	cmd := &cobra.Command{
		Use:          "unclaimed",
		Short:        "List PersistentVolumes that are available for claim",
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

// Run finds unclaimed PersistentVolumes.
func (o *UnclaimedPVsOptions) Run(ctx context.Context, noHeader bool) error {
	client, err := o.GetClient()
	if err != nil {
		return err
	}

	pvs, err := client.CoreV1().PersistentVolumes().List(ctx, metav1.ListOptions{})
	if err != nil {
		return err
	}

	var matrix [][]string

	for _, pv := range pvs.Items {
		if pv.Status.Phase == "Available" {
			age := getAge(pv.CreationTimestamp)
			row := []string{pv.Name, string(pv.Spec.PersistentVolumeReclaimPolicy), pv.Spec.StorageClassName, age}
			if o.allNamespaces {
				row = append([]string{pv.Namespace}, row...)
			}
			matrix = append(matrix, row)
		}
	}

	headers := []string{"NAME", "RECLAIM POLICY", "STORAGECLASS", "AGE"}

	buf := bytes.NewBuffer(nil)
	writeResults(buf, headers, matrix, o.namespace, noHeader)
	fmt.Printf("%s", buf.String())

	return nil
}
