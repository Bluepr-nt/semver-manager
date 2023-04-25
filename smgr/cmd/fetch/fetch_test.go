package fetch

import (
	"testing"

	"github.com/joho/godotenv"
)

func LoadEnvFromFile() error {
	err := godotenv.Load("ccs.yaml")
	if err != nil {
		return err
	}
	return nil
}

func TestNewFetchCommand(t *testing.T) {
	t.Run("Returns non-nil command", func(t *testing.T) {
		cmd := NewFetchCommand(nil)
		if cmd == nil {
			t.Error("NewFetchCommand returned nil command.")
		}
	})

	t.Run("Command name is 'fetch'", func(t *testing.T) {
		cmd := NewFetchCommand(nil)
		if cmd.Name() != "fetch" {
			t.Errorf("Command name is '%s', expected 'fetch'", cmd.Name())
		}
	})

	t.Run("Command has expected flags", func(t *testing.T) {
		cmd := NewFetchCommand(nil)
		flags := cmd.Flags()
		expectedFlags := []string{"owner", "repo", "token", "platform", "highest"}
		for _, expectedFlag := range expectedFlags {
			if flags.Lookup(expectedFlag) == nil {
				t.Errorf("Command does not have expected flag '%s'", expectedFlag)
			}
		}
	})

	t.Run("End-to-end command with valid inputs and validate outputs", func(t *testing.T) {
		// TODO
	})
}
