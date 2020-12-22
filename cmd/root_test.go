package main

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMainUnknownCommand(t *testing.T) {
	defer func(args []string) {
		os.Args = args
	}(os.Args)

	os.Args = []string{"janitor", "etoomanycookies"}

	err := Execute()

	assert.Error(t, err)
}
