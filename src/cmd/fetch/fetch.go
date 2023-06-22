package fetch

import (
	"src/cmd/smgr/cmd/utils"
	datasourceUtils "src/cmd/smgr/datasource/utils"
	"src/cmd/smgr/models"

	"strings"

	"github.com/spf13/cobra"
	"k8s.io/klog"
)

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
		Short: "Fetch semver tags from a repository.",
		Long: `Fetch semver tags from a repository and
sorted by highest version first. Fetch also supports all
the filters from the filter command. If the --versions
flag is set, the versions passed will be merged with the
fetched versions.`,

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
		cliTags, err := filterCmd.Flags().GetString("versions")
		if err != nil {
			return err
		}
		if len(cliTags) > 0 {
			klog.V(2).Infof("Merging tags: %s", cliTags)
			semverTags = append(semverTags, models.SplitVersions(cliTags)...)
		}
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
