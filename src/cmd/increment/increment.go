package increment

import (
	"errors"
	"fmt"

	"src/cmd/smgr/cmd/utils"
	"src/cmd/smgr/models"
	"src/cmd/smgr/pkg/filter"
	"src/cmd/smgr/pkg/increment"

	"github.com/spf13/cobra"
)

type config struct {
	dryRun        bool
	incrementType string
	// sourceVersion  string
	sourceVersions string
	// repository     string
	// sourceStream   string
	targetStream string
}

func NewIncrementCommand() *cobra.Command {
	config := &config{}
	incrementCmd := &cobra.Command{
		Use:   "increment",
		Short: "Increment a version",
		Long: `
Increment a version according to one of the required flags --level or --target-stream, 
and any combination of the optional flags:

- Use --level to specify the increment level (major, minor, patch).
- Define the source with --repository, --source-stream, --source-version, or --source-versions.

Increment a version according to the provided:
  - Increment level (major, minor, patch)
  - The source, any of: repository, source-stream, source-version, source-versions
  - 
  `,
		PersistentPreRun: func(cmd *cobra.Command, args []string) {
		},

		PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
			return utils.InitializeConfig(cmd)
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			dryRun, _ := cmd.PersistentFlags().GetBool("dry-run")
			config.dryRun = dryRun
			level, _ := cmd.Flags().GetString("level")
			targetStream, _ := cmd.Flags().GetString("target-stream")

			tmp := fmt.Sprintf("targetstream %s", targetStream)
			fmt.Println(tmp)
			if level == "" {
				fmt.Println("level is empty")
			}
			if targetStream == "" {
				fmt.Println("targetstream is empty")
			}
			if level == "" && targetStream == "" {
				cmd.Usage() // Prints usage information
				return errors.New("error: either --level or --target-stream must be specified")
			}

			return RunIncrement(config, cmd)
		},
	}

	// TODO require level or target stream

	incrementCmd.Flags().StringVarP(&config.incrementType, "level", "l", string(models.Patch), "The level of increment to perform, options: major, minor, patch (required if --target-stream not specified)")
	// incrementCmd.Flags().StringVarP(&config.sourceVersion, "source-version", "V", "0.0.0", "The source version to increment from")
	// incrementCmd.Flags().StringVarP(&config.sourceStream, "source-stream", "s", "", "The source stream to increment from ")
	incrementCmd.Flags().StringVarP(&config.targetStream, "target-stream", "t", "", "The target stream to increment to e.g. 1.2.* (required if --level not specified)")
	incrementCmd.Flags().StringVarP(&config.sourceVersions, "source-versions", "u", "", "The source versions to increment from e.g. \"0.0.0,1.0.0,1.1.0\" (optional)")
	// incrementCmd.Flags().StringVarP(&config.repository, "repository", "r", "", "The repository to increment the version of e.g. https://github.com/<user|org>/<repo> (optional)")

	return incrementCmd
}

func RunIncrement(config *config, cmd *cobra.Command) error {
	sourceVersions := filter.GetValidVersions(config.sourceVersions)

	var err error
	var targetStream models.VersionPattern
	if config.targetStream != "" {
		targetStream, err = models.ParseVersionPattern(config.targetStream)
		if err != nil {
			return err
		}
	}

	newVersion, err := increment.IncrementReleaseToStream(sourceVersions, targetStream, models.Increment(config.incrementType))
	if err != nil {
		return err
	}
	cmd.Print(newVersion.String())
	return nil
}
