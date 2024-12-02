package increment

import (
	"bytes"
	"testing"

	"github.com/joho/godotenv"
	"github.com/stretchr/testify/assert"
)

type testFlag struct {
	name  string
	value string
}

func LoadEnvFromFile() error {
	err := godotenv.Load("increment_test.yaml")
	if err != nil {
		return err
	}
	return nil
}

func TestNewIncrementCommand(t *testing.T) {
	t.Run("Returns non-nil command", func(t *testing.T) {
		cmd := NewIncrementCommand()
		assert.NotNil(t, cmd)
	})

	t.Run("Command name is 'increment'", func(t *testing.T) {
		cmd := NewIncrementCommand()
		assert.Equal(t, "increment", cmd.Name())
	})

	t.Run("Command has expected flags", func(t *testing.T) {
		cmd := NewIncrementCommand()
		flags := cmd.Flags()
		expectedFlags := []string{"level", "repository", "source-versions", "target-stream"}
		for _, expectedFlag := range expectedFlags {
			assert.NotNil(t, flags.Lookup(expectedFlag))

		}
	})
}

func TestNewIncrementCommandFlags(t *testing.T) {
	LoadEnvFromFile()
	tests := []struct {
		name               string
		flags              []testFlag
		expectedNewVersion string
		expectedError      error
	}{
		{
			name: "Increment major version",
			flags: []testFlag{
				{name: "level", value: "major"},
			},
			expectedNewVersion: "1.0.0",
			expectedError:      nil,
		}, {
			name: "Increment minor version",
			flags: []testFlag{
				{name: "level", value: "minor"},
			},
			expectedNewVersion: "0.1.0",
			expectedError:      nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			output := new(bytes.Buffer)

			cmd := NewIncrementCommand()
			cmd.SetOutput(output)
			cmd.SetErr(output)

			for _, flag := range tt.flags {
				cmd.Flags().Set(flag.name, flag.value)
			}

			err := cmd.Execute()
			assert.Equal(t, tt.expectedError, err)
			assert.Equal(t, tt.expectedNewVersion, output.String())
		})
	}
}
