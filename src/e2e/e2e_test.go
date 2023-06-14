package e2e

import (
	"fmt"
	"os"
	"os/exec"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestMain(m *testing.M) {
	err := buildBinary()
	if err != nil {
		fmt.Println("Failed to build binary")
		fmt.Println(err)
		os.Exit(1)
	}

	exitVal := m.Run()

	// You might want to cleanup the binary
	err = os.Remove("./smgr")
	if err != nil {
		os.Exit(1)
	}

	os.Exit(exitVal)
}

func buildBinary() error {
	cmd := exec.Command("go", "build", "-o", "smgr", "../cmd/smgr/")
	err := cmd.Run()
	return err
}

func TestFilterCommand(t *testing.T) {
	tests := []struct {
		name        string
		args        []string
		expectedOut string
		expectedErr string
	}{
		{
			name:        "Filter versions with wildcard pattern",
			args:        []string{"filter", "--versions", "1.0.0 2.0.0", "--stream", "1.*.*"},
			expectedOut: "1.0.0\n",
		},
		{
			name:        "Filter versions and select the highest",
			args:        []string{"filter", "--versions", "1.0.0", "2.0.0", "--highest"},
			expectedOut: "2.0.0\n",
		},
		{
			name:        "No matching version",
			args:        []string{"filter", "--versions", "1.0.0", "2.0.0", "--stream", "3.*.*"},
			expectedOut: "\n",
		},
		{
			name:        "Invalid version",
			args:        []string{"filter", "--versions", "invalid", "2.0.0", "--stream", "1.*.*"},
			expectedErr: "error: invalid version format\n",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cmd := exec.Command("./smgr", tt.args...)

			out, err := cmd.CombinedOutput()
			outStr := string(out)

			if tt.expectedErr != "" {
				require.Error(t, err)
				require.Contains(t, outStr, tt.expectedErr)
			} else {
				require.NoError(t, err)
				require.Equal(t, tt.expectedOut, outStr)
			}
		})
	}
}
