package cmd

import (
	"bytes"
	"context"
	"fmt"

	"github.com/spf13/cobra"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	cmdutil "k8s.io/kubectl/pkg/cmd/util"
)

// NoEndpointsSVCsOptions embeds JanitorOptions struct.
type NoEndpointsSVCsOptions struct {
	JanitorOptions
}

// newNoEndpointsSVCsOptions create a instance of NoEndpointsSVCsOptions.
func newNoEndpointsSVCsOptions(options JanitorOptions) *NoEndpointsSVCsOptions {
	return &NoEndpointsSVCsOptions{
		JanitorOptions: options,
	}
}

// newNoEndpointsSVCsCommand returns a cobra command wrapping NoEndpointsSVCsCommand.
func newNoEndpointsSVCsCommand(factory cmdutil.Factory, options JanitorOptions) *cobra.Command {
	o := newNoEndpointsSVCsOptions(options)

	cmd := &cobra.Command{
		Use:          "no-endpoints",
		Short:        "List Services without endpoints",
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

// Run finds Services without Endpoints.
func (o *NoEndpointsSVCsOptions) Run(ctx context.Context, noHeader bool) error {
	client, err := o.GetClient()
	if err != nil {
		return err
	}

	svcs, err := client.CoreV1().Services(o.namespace).List(ctx, metav1.ListOptions{})
	if err != nil {
		return err
	}

	endpoints, err := client.CoreV1().Endpoints(o.namespace).List(ctx, metav1.ListOptions{})
	if err != nil {
		return err
	}

	endpointsMap := make(map[string]corev1.Endpoints, len(endpoints.Items))
	for _, end := range endpoints.Items {
		endpointsMap[end.Name] = end
	}

	var matrix [][]string

	for _, svc := range svcs.Items {
		end, ok := endpointsMap[svc.Name]
		if !ok {
			row := []string{svc.Name}
			if o.allNamespaces {
				row = append([]string{svc.Namespace}, row...)
			}
			matrix = append(matrix, row)
			continue
		}

		if len(end.Subsets) == 0 {
			row := []string{svc.Name}
			if o.allNamespaces {
				row = append([]string{svc.Namespace}, row...)
			}
			matrix = append(matrix, row)
		}
	}

	headers := []string{"NAME"}

	buf := bytes.NewBuffer(nil)
	writeResults(buf, headers, matrix, o.namespace, noHeader)
	fmt.Printf("%s", buf.String())

	return nil
}
