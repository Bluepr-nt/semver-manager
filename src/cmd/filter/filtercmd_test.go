package filtercmd

import (
	"bytes"
	"strings"
	"testing"

	"github.com/spf13/cobra"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewFilterCommand(t *testing.T) {
	tests := []struct {
		name          string
		inputArgs     []string
		expectedError bool
		expectedOut   string
	}{
		{
			name:          "Missing required flag",
			inputArgs:     []string{},
			expectedError: true,
			expectedOut:   "Error: required flag(s) \"versions\" not set\nUsage:\n  filter [flags]\n\nFlags:\n  -h, --help                   help for filter\n  -H, --highest                Filter by highest version\n  -s, --stream string          Filter by major, minor, patch, prerelease version and build metadata streams\n  -v, --versions stringArray   Version list to filter",
		},
		{
			name:        "Provided version",
			inputArgs:   []string{"--versions", "1.2.3", "--highest"},
			expectedOut: "Command-line arguments: []\n",
		},
		// ... More test cases ...
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			filtercmd := NewFilterCommand()
			output, err := executeCommand(filtercmd, test.inputArgs...)

			if test.expectedError {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}

			assert.NoError(t, err)
			assert.Equal(t, test.expectedOut, output)
		})
	}
}

func executeCommand(cmd *cobra.Command, args ...string) (string, error) {
	buf := new(bytes.Buffer)
	cmd.SetOut(buf)
	cmd.SetErr(buf)
	cmd.SetArgs(args)
	err := cmd.Execute()
	return strings.TrimSpace(buf.String()), err
}
