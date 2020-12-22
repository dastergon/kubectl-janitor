package cmd

import (
	"bytes"
	"context"
	"fmt"

	"github.com/spf13/cobra"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	cmdutil "k8s.io/kubectl/pkg/cmd/util"
)

// FailedJobsOptions embeds JanitorOptions struct.
type FailedJobsOptions struct {
	JanitorOptions
}

// newFailedJobsOptions create a instance of FailedJobsOptions.
func newFailedJobsOptions(options JanitorOptions) *FailedJobsOptions {
	return &FailedJobsOptions{
		JanitorOptions: options,
	}
}

// newFailedJobsCommand returns a cobra command wrapping BlockedJobsOptions.
func newFailedJobsCommand(factory cmdutil.Factory, options JanitorOptions) *cobra.Command {
	o := newFailedJobsOptions(options)

	cmd := &cobra.Command{
		Use:          "failed",
		Short:        "List Jobs that have failed to run and have restartPolicy: Never",
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

// Run lists stuck Jobs that cannot restart.
func (o *FailedJobsOptions) Run(ctx context.Context, noHeader bool) error {
	client, err := o.GetClient()
	if err != nil {
		return err
	}

	jobs, err := client.BatchV1().Jobs(o.namespace).List(ctx, metav1.ListOptions{})
	if err != nil {
		return err
	}

	var matrix [][]string

	for _, job := range jobs.Items {
		if job.Spec.Template.Spec.RestartPolicy == "Never" {
			for _, c := range job.Status.Conditions {
				if c.Reason == "BackoffLimitExceeded" || c.Reason == "DeadlineExceeded" {
					age := getAge(job.CreationTimestamp)
					row := []string{job.Name, c.Reason, c.Message, age}
					if o.allNamespaces {
						row = append([]string{job.Namespace}, row...)
					}
					matrix = append(matrix, row)
				}
			}
		}
	}

	headers := []string{"NAME", "REASON", "MESSAGE", "AGE"}

	buf := bytes.NewBuffer(nil)
	writeResults(buf, headers, matrix, o.namespace, noHeader)
	fmt.Printf("%s", buf.String())

	return nil
}
