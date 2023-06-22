package filter

import (
	"src/cmd/smgr/models"
	"src/cmd/smgr/pkg/filter"

	"github.com/spf13/cobra"
)

func NewFilterCommand() *cobra.Command {
	filterArgs := &FilterArgs{}
	var filterCmd = &cobra.Command{
		Use:   "filter",
		Short: "Filter is a CLI tool for filtering versions",
		Long:  `Filter is a CLI tool for filtering versions using various criteria.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			versions := filter.GetVersions(filterArgs.Versions)

			filters := []filter.FilterFunc{}
			if filterArgs.StreamFilter != "" {
				pattern, err := models.ParseVersionPattern(filterArgs.StreamFilter)
				if err != nil {
					return err
				}
				filters = append(filters, filter.VersionPatternFilter(pattern))
			}

			if filterArgs.Highest {
				filters = append(filters, filter.Highest())
			}

			semverTags, err := filter.ApplyFilters(versions, filters...)
			if err != nil {
				return err
			}
			cmd.Println(semverTags.String())
			return nil
		},
	}

	filterCmd.Flags().StringVarP(&filterArgs.Versions, "versions", "V", "", "Version list to filter")
	filterCmd.Flags().StringVarP(&filterArgs.StreamFilter, "stream", "s", "", "Filter by major, minor, patch, prerelease version and build metadata streams")
	filterCmd.Flags().BoolVarP(&filterArgs.Highest, "highest", "H", false, "Filter by highest version")
	return filterCmd
}
