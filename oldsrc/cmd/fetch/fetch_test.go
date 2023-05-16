package fetch

import (
	"bytes"
	"fmt"
	"os"
	"testing"

	"github.com/joho/godotenv"
	"github.com/stretchr/testify/assert"
)

type testFlag struct {
	name  string
	value string
}

func LoadEnvFromFile() error {
	err := godotenv.Load("fetch_test.yaml")
	if err != nil {
		return err
	}
	return nil
}

func TestNewFetchCommand(t *testing.T) {
	t.Run("Returns non-nil command", func(t *testing.T) {
		cmd := NewFetchCommand()
		if cmd == nil {
			t.Error("NewFetchCommand returned nil command.")
		}
	})

	t.Run("Command name is 'fetch'", func(t *testing.T) {
		cmd := NewFetchCommand()
		if cmd.Name() != "fetch" {
			t.Errorf("Command name is '%s', expected 'fetch'", cmd.Name())
		}
	})

	t.Run("Command has expected flags", func(t *testing.T) {
		cmd := NewFetchCommand()
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

func TestNewFetchCommandRealRepo(t *testing.T) {
	LoadEnvFromFile()
	tests := []struct {
		name         string
		flags        []testFlag
		expectedTags string
	}{
		{
			name: "fetch all semver tags",
			flags: []testFlag{
				{"owner", os.Getenv("OWNER")},
				{"repo", os.Getenv("REPO")},
				{"token", os.Getenv("TOKEN")},
				{"platform", "github"},
				{"highest", "false"},
			},
			expectedTags: "0.0.4 1.0.0-0A.is.legal 1.0.0-alpha 1.0.0-alpha+beta 1.0.0-alpha.1 1.0.0-alpha.0valid 1.0.0-alpha.beta 1.0.0-alpha.beta.1 1.0.0-alpha-a.b-c-somethinglong+build.1-aef.1-its-okay 1.0.0-alpha0.valid 1.0.0-beta 1.0.0-rc.1+build.1 1.0.0 1.0.0+0.build.1-rc.10000aaa-kk-0.1 1.1.2-prerelease+meta 1.1.2+meta 1.1.2+meta-valid 1.1.7 1.2.3----R-S.12.9.1--.12+meta 1.2.3----RC-SNAPSHOT.12.9.1--.12+788 1.2.3----RC-SNAPSHOT.12.9.1--.12 1.2.3-SNAPSHOT-123 1.2.3-beta 1.2.3 2.0.0-rc.1+build.123 2.0.0+build.1848 2.0.0 2.0.1-alpha.1227 10.2.3-DEV-SNAPSHOT 10.20.30\n",
		},
		{
			name: "fetch only the highest semver tag",
			flags: []testFlag{
				{"owner", os.Getenv("OWNER")},
				{"repo", os.Getenv("REPO")},
				{"token", os.Getenv("TOKEN")},
				{"platform", "github"},
				{"highest", "true"},
			},
			expectedTags: "10.20.30\n",
		},
		{
			name: "fetch only release semver tags",
			flags: []testFlag{
				{"owner", os.Getenv("OWNER")},
				{"repo", os.Getenv("REPO")},
				{"token", os.Getenv("TOKEN")},
				{"platform", "github"},
				{"highest", "false"},
				{"release", "true"},
			},
			expectedTags: "0.0.4 1.0.0 1.0.0+0.build.1-rc.10000aaa-kk-0.1 1.1.2+meta 1.1.2+meta-valid 1.1.7 1.2.3 2.0.0+build.1848 2.0.0 10.20.30\n",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			cmd := NewFetchCommand()

			output := new(bytes.Buffer)
			cmd.SetOut(output)
			cmd.SetErr(output)

			var args []string
			for _, flag := range tc.flags {
				err := cmd.Flags().Set(flag.name, flag.value)
				args = append(args, fmt.Sprintf("--%s", flag.name))
				args = append(args, flag.value)
				if err != nil {
					t.Fatalf("unexpected error setting flag '%s': %v", flag.name, err)
				}
			}

			err := cmd.Execute()
			assert.NoError(t, err, "Expected no error")

			assert.Equal(t, tc.expectedTags, output.String())
		})
	}
}
