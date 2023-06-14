package filter

import (
	"src/cmd/smgr/models"
	"src/cmd/smgr/pkg/filter"
	"strings"

	"github.com/spf13/cobra"
)

func NewFilterCommand() *cobra.Command {
	filterArgs := &FilterArgs{}
	var filterCmd = &cobra.Command{
		Use:   "filter",
		Short: "Filter is a CLI tool for filtering versions",
		Long:  `Filter is a CLI tool for filtering versions using various criteria.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			versionList := []string{}
			// trim begining and trailing whitespaces in filterArgs.Versions
			filterArgs.Versions = strings.TrimSpace(filterArgs.Versions)
			// if filterArgs contains commas, remove all whitespaces
			if strings.Contains(filterArgs.Versions, ",") {
				filterArgs.Versions = strings.ReplaceAll(filterArgs.Versions, " ", "")
				versionList = strings.Split(filterArgs.Versions, ",")
			} else {
				versionList = strings.Split(filterArgs.Versions, " ")
			}

			filterArgs.Versions = strings.ReplaceAll(filterArgs.Versions, " ", ",")

			versions := make([]models.Version, len(versionList))
			var err error
			for i, versionStr := range versionList {
				versions[i], err = models.ParseVersion(versionStr)
				if err != nil {
					return err
				}
			}

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
	filterCmd.MarkFlagRequired("versions")
	filterCmd.Flags().StringVarP(&filterArgs.StreamFilter, "stream", "s", "", "Filter by major, minor, patch, prerelease version and build metadata streams")
	filterCmd.Flags().BoolVarP(&filterArgs.Highest, "highest", "H", false, "Filter by highest version")
	return filterCmd
}
