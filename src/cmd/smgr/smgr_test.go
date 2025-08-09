package main

import (
	"bytes"

	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewRootCommand(t *testing.T) {

	t.Run("Base Tests", func(t *testing.T) {
		output := &bytes.Buffer{}
		cmd := NewRootCommand(output)
		assert.NotEmpty(t, cmd, "NewRootCommand() should not return nil command")
		assert.Equal(t, "smgr", cmd.Name(), "NewRootCommand() should return command with name 'smgr'")

		assert.Equal(t, "dry-run", cmd.PersistentFlags().Lookup("dry-run").Name, "NewRootCommand() should have 'dry-run' flag")
	})
}
