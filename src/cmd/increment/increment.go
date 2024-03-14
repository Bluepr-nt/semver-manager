package increment

import (
	"src/cmd/smgr/cmd/utils"
	"src/cmd/smgr/models"

	"github.com/spf13/cobra"
)

type config struct {
	dryRun         bool
	incrementType  string
	sourceVersion  string
	sourceVersions string
	repository     string
	sourceStream   string
	targetStream   string
}

func NewIncrementCommand() *cobra.Command {
	config := &config{}
	incrementCmd := &cobra.Command{
		Use:   "increment",
		Short: "Increment a version",
		Long:  `Increment a version according to the provided increment type.`,
		PersistentPreRun: func(cmd *cobra.Command, args []string) {
		},

		PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
			return utils.InitializeConfig(cmd)
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			dryRun, _ := cmd.PersistentFlags().GetBool("dry-run")
			config.dryRun = dryRun
			return RunIncrement(config, cmd)
		},
	}

	incrementCmd.Flags().StringVarP(&config.incrementType, "level", "l", string(models.Patch), "The level of increment to perform, options: major, minor, patch")
	incrementCmd.Flags().StringVarP(&config.sourceVersion, "source-version", "V", "0.0.0", "The source version to increment from")
	incrementCmd.Flags().StringVarP(&config.sourceStream, "source-stream", "s", "", "The source stream to increment from")
	incrementCmd.Flags().StringVarP(&config.targetStream, "target-stream", "t", "", "The target stream to increment to")
	incrementCmd.Flags().StringVarP(&config.sourceVersions, "source-versions", "u", "", "The source versions to increment from")
	incrementCmd.Flags().StringVarP(&config.repository, "repository", "r", "", "The repository to increment the version of")
	return incrementCmd
}

func RunIncrement(config *config, cmd *cobra.Command) error {

	return nil
}
