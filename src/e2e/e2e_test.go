package e2e

import (
	"fmt"
	"os"
	"os/exec"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMain(m *testing.M) {
	err := buildBinary()
	if err != nil {
		fmt.Println("Failed to build binary")
		fmt.Println(err)
		os.Exit(1)
	}

	exitVal := m.Run()

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
			args:        []string{"filter", "--versions", "1.0.0 2.0.0", "--highest"},
			expectedOut: "2.0.0\n",
		},
		{
			name:        "No matching version",
			args:        []string{"filter", "--versions", "1.0.0 2.0.0", "--stream", "3.*.*"},
			expectedOut: "\n",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cmd := exec.Command("./smgr", tt.args...)

			out, err := cmd.CombinedOutput()
			outStr := string(out)

			if tt.expectedErr != "" {
				assert.Error(t, err)
				assert.Contains(t, outStr, tt.expectedErr)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expectedOut, outStr)
			}
		})
	}
}

func TestFetchCommand(t *testing.T) {
	token := os.Getenv("GITHUB_TOKEN")
	if token == "" {
		t.Skip("Skipping TestFetchCommand: GITHUB_TOKEN environment variable not set")
		return
	}

	tests := []struct {
		name        string
		args        []string
		expectedOut string
		expectedErr string
	}{
		{
			name:        "Fetch specific version",
			args:        []string{"fetch", "-o", "Bluepr-nt", "-r", "semver-manager", "-t", token},
			expectedOut: "0.1.0 0.1.1 0.1.2 0.1.3 0.1.4\n",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cmd := exec.Command("./smgr", tt.args...)

			out, err := cmd.CombinedOutput()
			outStr := string(out)

			if tt.expectedErr != "" {
				assert.Error(t, err)
				assert.Contains(t, outStr, tt.expectedErr)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expectedOut, outStr)
			}
		})
	}
}
