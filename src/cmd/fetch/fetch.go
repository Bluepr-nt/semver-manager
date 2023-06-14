package fetchcmd

import (
	"fmt"
	"log"
	"src/cmd/smgr/pkg/fetch"
	"src/cmd/smgr/pkg/filter"
	"src/cmd/smgr/util"

	"github.com/spf13/cobra"
)

var (
	major       uint64
	minor       uint64
	patch       uint64
	prerelease  []string
	releaseOnly bool
	highest     bool
)

// fetchCmd represents the fetch command
var fetchCmd = &cobra.Command{
	Use:   "fetch",
	Short: "Fetch is a CLI tool for fetching versions",
	Long:  `Fetch is a CLI tool for fetching versions from some source.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		fmt.Println("Fetching versions...")
		fetcher, err := fetch.NewFetcher(&util.DatasourceConfig{})
		if err != nil {
			return err
		}
		versions, err := fetcher.FetchTags()
		if err != nil {
			log.Fatal(err)
		}

		filters := []filter.FilterFunc{}

		// if cmd.Flags().Changed("major") {
		// 	filters = append(filters, filter.MajorVersionStream(major))
		// }
		// if cmd.Flags().Changed("minor") {
		// 	filters = append(filters, filter.MinorVersionStream(major, minor))
		// }
		// if cmd.Flags().Changed("patch") {
		// 	filters = append(filters, filter.PatchVersionStream(major, minor, patch))
		// }
		// if cmd.Flags().Changed("prerelease") {
		// 	filters = append(filters, filter.PreReleaseVersionStream(models.Release{}, models.PRVersion{}))
		// }
		if highest {
			filters = append(filters, filter.Highest())
		}

		filteredVersions, err := filter.ApplyFilters(versions, filters...)
		if err != nil {
			log.Fatal(err)
		}
		cmd.OutOrStdout().Write([]byte(filteredVersions.String()))
		return nil
	},
}

func init() {
	fetchCmd.Flags().BoolVar(&highest, "highest", false, "Select the highest versions")
}
