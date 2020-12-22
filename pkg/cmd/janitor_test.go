package cmd

import (
	"bytes"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"k8s.io/cli-runtime/pkg/genericclioptions"
)

func NewTestJanitorOptions() (JanitorOptions, *bytes.Buffer, *bytes.Buffer, *bytes.Buffer) {
	streams, in, out, errout := genericclioptions.NewTestIOStreams()
	rbFlags := &genericclioptions.ResourceBuilderFlags{}
	rbFlags.WithAllNamespaces(false)
	return JanitorOptions{
		ConfigFlags:          genericclioptions.NewConfigFlags(true),
		ResourceBuilderFlags: rbFlags,
		Streams:              streams,
	}, in, out, errout
}

func TestNewJanitorCommandHelp(t *testing.T) {
	_, _, stdout, stderr := NewTestJanitorOptions()

	defer func(args []string) {
		os.Args = args
	}(os.Args)
	os.Args = []string{"janitor", "help"}

	root := NewJanitorCommand()
	root.SetOut(stdout)
	err := root.Execute()

	assert.NoError(t, err)
	assert.Equal(t, "", stderr.String())
	assert.Contains(t, stdout.String(), "Available Commands:")
}

func TestNewJanitorCommandUnknownCommand(t *testing.T) {
	defer func(args []string) {
		os.Args = args
	}(os.Args)
	os.Args = []string{"janitor", "etoomuchcookies"}

	root := NewJanitorCommand()
	err := root.Execute()

	assert.Error(t, err)
}
