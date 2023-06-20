package fetch

import (
	"src/cmd/smgr/cmd/utils"
	datasourceUtils "src/cmd/smgr/datasource/utils"

	"strings"

	"github.com/spf13/cobra"
	"k8s.io/klog"
)

// fetchCmd represents the fetch command
// var fetchCmd = &cobra.Command{
// 	Use:   "fetch",
// 	Short: "Fetch is a CLI tool for fetching versions",
// 	Long:  `Fetch is a CLI tool for fetching versions from some source.`,
// 	RunE: func(cmd *cobra.Command, args []string) error {
// 		fmt.Println("Fetching versions...")
// 		fetcher, err := fetch.NewFetcher(&util.DatasourceConfig{})
// 		if err != nil {
// 			return err
// 		}
// 		versions, err := fetcher.FetchTags()
// 		if err != nil {
// 			log.Fatal(err)
// 		}

// 		filters := []filter.FilterFunc{}

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
// 		if highest {
// 			filters = append(filters, filter.Highest())
// 		}

// 		filteredVersions, err := filter.ApplyFilters(versions, filters...)
// 		if err != nil {
// 			log.Fatal(err)
// 		}
// 		cmd.OutOrStdout().Write([]byte(filteredVersions.String()))
// 		return nil
// 	},
// }

// func init() {
// 	fetchCmd.Flags().BoolVar(&highest, "highest", false, "Select the highest versions")
// }

type config struct {
	Token      string `san:"trim"`
	Repository string `san:"trim"`
	Owner      string `san:"trim"`
	Platform   string `san:"trim"`
	dryRun     bool
}

func NewFetchCommand(filterCmd *cobra.Command) *cobra.Command {
	config := &config{}
	var fetchCmd = &cobra.Command{
		Use:   "fetch",
		Short: "Fetch semver tags from a registry or repository.",
		Long: `Fetch semver tags from a registry or repository and
    sorted by highest version first. Includes tags starting 'v' as in 'v0.0.0`,

		PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
			return utils.InitializeConfig(cmd)
		},

		RunE: func(cmd *cobra.Command, args []string) error {
			dryRun, _ := cmd.PersistentFlags().GetBool("dry-run")
			config.dryRun = dryRun

			return RunFetchSemverTags(config, cmd, filterCmd)
		},
	}
	fetchCmd.Flags().StringVarP(&config.Owner, "owner", "o", "", "The owner of the registry or repository")
	fetchCmd.Flags().StringVarP(&config.Repository, "repo", "r", "", "The repository or registry to fetch the Semver tags from")
	fetchCmd.Flags().StringVarP(&config.Token, "token", "t", "", "The token to access the repository")
	fetchCmd.Flags().StringVarP(&config.Platform, "platform", "p", "github", "The platform to fetch the Semver from, options: github")
	if filterCmd != nil {
		fetchCmd.Flags().AddFlagSet(filterCmd.Flags())
	}
	return fetchCmd
}

func RunFetchSemverTags(config *config, cmd *cobra.Command, filterCmd *cobra.Command) error {

	datasource := newDatasource(config.dryRun, config.Platform, config.Token)

	klog.V(1).Info("Fetching tags...")
	semverTags, err := datasource.FetchSemverTags(config.Owner, config.Repository)
	if err != nil {
		return err
	}

	if filterCmd != nil {
		klog.V(1).Info("Filtering tags...")
		err = filterCmd.Flags().Set("versions", strings.Join(semverTags, " "))
		if err != nil {
			return err
		}

		err = filterCmd.Execute()
		if err != nil {
			return err
		}
	} else {
		cmd.Println(strings.Join(semverTags, " "))
	}

	return nil
}

func newDatasource(dryRun bool, platform, token string) datasourceUtils.Datasource {
	if len(platform) == 0 {
		platform = "github"
	}
	if dryRun {
		platform = "dry-run"
	}
	semverSvc := datasourceUtils.NewSemverSvc(platform, token)

	return semverSvc
}
