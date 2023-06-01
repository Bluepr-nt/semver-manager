package filtercmd

import (
	"fmt"
	"src/pkg/fetch/models"
	"src/pkg/fetch/pkg/filter"

	"github.com/spf13/cobra"
)

type FilterArgs struct {
	VersionStream string
	Filters       []string
	StreamFilter  string // Major, Minor, Patch, Release, Pre-release
	Highest       string
}

func NewFilterCommand() *cobra.Command {
	config := &FilterArgs{}
	var filterCmd = &cobra.Command{
		Use:   "filter",
		Short: "Filter is a CLI tool for filtering versions",
		Long:  `Filter is a CLI tool for filtering versions using various criteria.`,
		Run: func(cmd *cobra.Command, args []string) {
			// args contains the command-line arguments
			fmt.Println("Command-line arguments:", args)

			// Assuming args contain the versions and you have a ParseVersion function to parse the string to a version.
			// Make sure to handle the error from ParseVersion properly.
			versions := make([]models.Version, len(args))
			for i, versionStr := range args {
				versions[i], _ = models.ParseVersion(versionStr)
			}

			filter.ApplyFilters(versions)
			// ... Rest of your filter command implementation ...
		},
	}
	filterCmd.Flags().StringVarP(&config.StreamFilter, "stream", "s", "", "Filter by major, minor, or patch version stream")
	filterCmd.Flags().StringVarP(&config.Highest, "highest", "H", "", "Filter by highest version")

	return filterCmd
}
