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

// NoServiceINGsOptions embeds JanitorOptions struct.
type NoServiceINGsOptions struct {
	JanitorOptions
}

// newNoServiceINGsOptions create a instance of NoServiceINGsOptions.
func newNoServiceINGsOptions(options JanitorOptions) *NoServiceINGsOptions {
	return &NoServiceINGsOptions{
		JanitorOptions: options,
	}
}

// newNoServiceINGsCommand returns a cobra command wrapping NoServiceINGsOptions.
func newNoServiceINGsCommand(factory cmdutil.Factory, options JanitorOptions) *cobra.Command {
	o := newNoServiceINGsOptions(options)

	cmd := &cobra.Command{
		Use:          "no-service",
		Short:        "List Ingresses without service",
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

// Run finds Ingresses without Serivices.
func (o *NoServiceINGsOptions) Run(ctx context.Context, noHeader bool) error {
	client, err := o.GetClient()
	if err != nil {
		return err
	}

	ings, err := client.NetworkingV1().Ingresses(o.namespace).List(ctx, metav1.ListOptions{})
	if err != nil {
		return err
	}

	svcs, err := client.CoreV1().Services(o.namespace).List(ctx, metav1.ListOptions{})
	if err != nil {
		return err
	}

	serviceMap := make(map[string]corev1.Service, len(svcs.Items))
	for _, end := range svcs.Items {
		serviceMap[end.Name] = end
	}

	var matrix [][]string

	for _, ing := range ings.Items {
		for _, r := range ing.Spec.Rules {
			for _, p := range r.HTTP.Paths {
				_, ok := serviceMap[p.Backend.Service.Name]
				if !ok {
					row := []string{ing.Name, r.Host, p.Path}
					if o.allNamespaces {
						row = append([]string{ing.Namespace}, row...)
					}
					matrix = append(matrix, row)
					continue
				}
			}
		}
	}

	headers := []string{"NAME", "HOST", "PATH"}

	buf := bytes.NewBuffer(nil)
	writeResults(buf, headers, matrix, o.namespace, noHeader)
	fmt.Printf("%s", buf.String())

	return nil
}
