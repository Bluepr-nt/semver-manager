package e2e

import (
	"fmt"
	"os"
	"os/exec"
	"testing"

	"github.com/joho/godotenv"
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
	cmd := exec.Command("go", "build", "-buildvcs=false", "-o", "smgr", "../cmd/smgr/")
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
	tests := []struct {
		name        string
		args        []string
		expectedOut string
		expectedErr string
	}{
		{
			name:        "Fetch specific version",
			args:        []string{"fetch", "-o", "13013SwagR", "-r", "semver-manager-test", "-t", GetGithubToken(), ""},
			expectedOut: "0.0.4 1.0.0 1.0.0+0.build.1-rc.10000aaa-kk-0.1 1.1.2+meta 1.1.2+meta-valid 1.1.7 1.2.3 2.0.0+build.1848 2.0.0 10.20.30\n",
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

func GetGithubToken() string {
	LoadEnvFromFile()
	token := os.Getenv("GITHUB_TOKEN")
	if token == "" {
		panic("Please set the GITHUB_TOKEN environment variable to run the tests.")
	}
	return token
}

func LoadEnvFromFile() error {
	err := godotenv.Load("fetch_test.yaml")
	if err != nil {
		return err
	}
	return nil
}
