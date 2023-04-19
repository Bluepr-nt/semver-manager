package main

import (
	"flag"
	"io"
	"os"

	"github.com/spf13/cobra"

	"k8s.io/klog/v2"
	fetch "semver-manager.thephoenixhomelab.com/cmd"
	"semver-manager.thephoenixhomelab.com/internal/utils"
)

// parse flags

// Review password length
type fetchFlags struct {
	Token      string `san:"trim"`
	Repository string `san:"trim"`
	Owner      string `san:"trim"`
	Platform   string `san:"trim"`
}

type config struct {
	dryRun bool
}

func main() {
	klog.InitFlags(nil)
	flag.Set("logtostderr", "false")
	flag.Set("alsologtostderr", "false")
	flag.Parse()

	cmd := NewRootCommand(os.Stdout)
	if err := cmd.Execute(); err != nil {
		os.Exit(1)
	}
}

func NewRootCommand(output io.Writer) *cobra.Command {
	klog.SetOutput(output)
	config := config{
		dryRun: false,
	}

	rootCmd := &cobra.Command{
		Use:   "smgr",
		Short: "Manage Semantic Versioning compliant versions.",
		Long:  `Manage Semantic Versioning compliant versions and integrate with popular or registry platform to facilitate the task.`,
		PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
			return utils.InitializeConfig(cmd)
		},
		Run: func(cmd *cobra.Command, args []string) {
			cmd.Help()
		},
	}

	rootCmd.PersistentFlags().BoolVar(&config.dryRun, "dry-run", false, "Execute the command in dry-run mode")
	rootCmd.PersistentFlags().AddGoFlagSet(flag.CommandLine)
	rootCmd.AddCommand(fetch.NewFetchCommand(output))

	return rootCmd
}

// Bind each cobra flag to its associated viper configuration (config file and environment variable)
