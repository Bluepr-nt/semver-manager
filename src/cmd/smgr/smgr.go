package main

import (
	"flag"
	"fmt"
	"io"
	"os"

	"github.com/spf13/cobra"
	"k8s.io/klog/v2"

	"src/cmd/smgr/cmd/fetch"
	"src/cmd/smgr/cmd/filter"
	"src/cmd/smgr/cmd/increment"
	"src/cmd/smgr/cmd/utils"
)

type config struct {
	dryRun bool
}

func NewRootCommand(output io.Writer) *cobra.Command {
	config := config{
		dryRun: false,
	}
	cmd := &cobra.Command{
		Use:   "smgr",
		Short: "Manage Semantic Versioning compliant versions.",
		Long:  `Manage Semantic Versioning compliant versions and integrate with popular repository and registry platform to facilitate the task.`,
		PersistentPreRun: func(cmd *cobra.Command, args []string) {
			utils.InitializeConfig(cmd)
		},
		Run: func(cmd *cobra.Command, args []string) {
			cmd.Help()
		},
	}

	cmd.PersistentFlags().BoolVar(&config.dryRun, "dry-run", false, "Execute the command in dry-run mode")
	cmd.PersistentFlags().AddGoFlagSet(flag.CommandLine)
	filterCmd := filter.NewFilterCommand()
	fetchCmd := fetch.NewFetchCommand(filterCmd)
	incrementCmd := increment.NewIncrementCommand()
	cmd.AddCommand(filterCmd, fetchCmd, incrementCmd)

	return cmd
}

func main() {
	klog.InitFlags(nil)
	flag.Set("logtostderr", "false")
	flag.Set("alsologtostderr", "false")
	flag.Parse()
	rootCmd := NewRootCommand(nil)

	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
