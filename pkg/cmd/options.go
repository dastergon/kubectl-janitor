package cmd

import (
	"os"

	"github.com/spf13/cobra"
	"k8s.io/cli-runtime/pkg/genericclioptions"
	"k8s.io/client-go/kubernetes"
	cmdutil "k8s.io/kubectl/pkg/cmd/util"
)

// JanitorOptions holds user configuration options.
// setting up the CLI and the connection to the cluster.
type JanitorOptions struct {
	Streams              genericclioptions.IOStreams
	ConfigFlags          *genericclioptions.ConfigFlags
	ResourceBuilderFlags *genericclioptions.ResourceBuilderFlags
	namespace            string
	allNamespaces        bool
}

// NewJanitorOptions provides an instance of JanitorOptions with default values.
func NewJanitorOptions() JanitorOptions {
	rbFlags := &genericclioptions.ResourceBuilderFlags{}
	rbFlags.WithAllNamespaces(false)

	return JanitorOptions{
		ConfigFlags:          genericclioptions.NewConfigFlags(true),
		ResourceBuilderFlags: rbFlags,
		Streams: genericclioptions.IOStreams{
			In:     os.Stdin,
			Out:    os.Stdout,
			ErrOut: os.Stderr,
		},
	}
}

func (o *JanitorOptions) GetClient() (*kubernetes.Clientset, error) {
	restConfig, err := o.ConfigFlags.ToRESTConfig()
	if err != nil {
		return nil, err
	}

	clientSet, err := kubernetes.NewForConfig(restConfig)
	if err != nil {
		return nil, err
	}
	return clientSet, nil

}

// Complete sets all information required for working with Kubernetes.
func (o *JanitorOptions) Complete(factory cmdutil.Factory, cmd *cobra.Command) error {
	var err error
	o.namespace, _, err = factory.ToRawKubeConfigLoader().Namespace()
	if err != nil {
		return err
	}

	if cmd.Flag("all-namespaces").Changed {
		o.allNamespaces = *o.ResourceBuilderFlags.AllNamespaces
		o.namespace = ""
	}

	return nil
}
