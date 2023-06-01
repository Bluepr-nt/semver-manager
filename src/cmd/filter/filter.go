package filtercmd

import (
	"fmt"
	"src/pkg/fetch/models"
	"src/pkg/fetch/pkg/filter"

	"github.com/spf13/cobra"
)

func NewFilterCommand() *cobra.Command {
	filterArgs := &FilterArgs{}
	var filterCmd = &cobra.Command{
		Use:   "filter",
		Short: "Filter is a CLI tool for filtering versions",
		Long:  `Filter is a CLI tool for filtering versions using various criteria.`,
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println("Command-line arguments:", args)
			versions := []models.Version{}
			var err error
			for i, versionStr := range filterArgs.Versions {
				versions[i], err = models.ParseVersion(versionStr)
				if err != nil {
					panic(err)
				}
			}

			filters := []filter.FilterFunc{}
			if filterArgs.StreamFilter != "" {
				filters = append(filters, filter.StreamFilter(filterArgs.StreamFilter))
				filter.ApplyFilters(versions)
				// ... Rest of your filter command implementation ...
			}
		},
	}
	filterCmd.Flags().StringArrayVarP(&filterArgs.Versions, "versions", "v", []string{}, "Version list to filter")
	filterCmd.MarkFlagRequired("versions")
	filterCmd.Flags().StringVarP(&filterArgs.StreamFilter, "stream", "s", "", "Filter by major, minor, patch, prerelease version and build metadata streams")
	filterCmd.Flags().BoolVarP(&filterArgs.Highest, "highest", "H", false, "Filter by highest version")
	return filterCmd
}
