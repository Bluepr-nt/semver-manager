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
			args:        []string{"filter", "--versions", "1.0.0 2.0.0", "--highest"},
			expectedOut: "2.0.0\n",
		},
		{
			name:        "No matching version",
			args:        []string{"filter", "--versions", "1.0.0 2.0.0", "--stream", "3.*.*"},
			expectedOut: "\n",
		},
		{
			name: "Invalid version",
			args: []string{"filter", "--versions", "invalid 2.0.0", "--stream", "1.*.*"},
			expectedErr: `Error: major MUST comprise only ASCII numerics [0-9], got: invalid
Usage:
  smgr filter [flags]

Flags:
  -h, --help              help for filter
  -H, --highest           Filter by highest version
  -s, --stream string     Filter by major, minor, patch, prerelease version and build metadata streams
  -V, --versions string   Version list to filter

Global Flags:
      --add_dir_header                   If true, adds the file directory to the header of the log messages
      --alsologtostderr                  log to standard error as well as files (no effect when -logtostderr=true)
      --dry-run                          Execute the command in dry-run mode
      --log_backtrace_at traceLocation   when logging hits line file:N, emit a stack trace (default :0)
      --log_dir string                   If non-empty, write log files in this directory (no effect when -logtostderr=true)
      --log_file string                  If non-empty, use this log file (no effect when -logtostderr=true)
      --log_file_max_size uint           Defines the maximum size a log file can grow to (no effect when -logtostderr=true). Unit is megabytes. If the value is 0, the maximum file size is unlimited. (default 1800)
      --logtostderr                      log to standard error instead of files
      --one_output                       If true, only write logs to their native severity level (vs also writing to each lower severity level; no effect when -logtostderr=true)
      --skip_headers                     If true, avoid header prefixes in the log messages
      --skip_log_headers                 If true, avoid headers when opening log files (no effect when -logtostderr=true)
      --stderrthreshold severity         logs at or above this threshold go to stderr when writing to files and stderr (no effect when -logtostderr=true or -alsologtostderr=false) (default 2)
  -v, --v Level                          number for the log level verbosity
      --vmodule moduleSpec               comma-separated list of pattern=N settings for file-filtered logging

major MUST comprise only ASCII numerics [0-9], got: invalid`,
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
