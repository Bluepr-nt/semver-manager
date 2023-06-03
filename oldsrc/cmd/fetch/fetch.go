package fetch

import (
	"strings"

	"github.com/spf13/cobra"
	"k8s.io/klog/v2"
	"src/cmd/smgr/internal/utils"
	smvr "src/cmd/smgr/pkg/semver"
)

type config struct {
	Token      string `san:"trim"`
	Repository string `san:"trim"`
	Owner      string `san:"trim"`
	Platform   string `san:"trim"`
	dryRun     bool
	Filters    smvr.Filters
}

func NewFetchCommand() *cobra.Command {
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
			if err := utils.SanitizeInputs(); err != nil {
				klog.Errorf("CLI argument error: %w", err)
				panic(config)
			}

			dryRun, _ := cmd.PersistentFlags().GetBool("dry-run")
			config.dryRun = dryRun

			return RunFetchSemverTags(config, cmd)
		},
	}
	fetchCmd.Flags().StringVarP(&config.Owner, "owner", "o", "", "The owner of the registry or repository")
	fetchCmd.Flags().StringVarP(&config.Repository, "repo", "r", "", "The repository or registry to fetch the Semver tags from")
	fetchCmd.Flags().StringVarP(&config.Token, "token", "t", "", "The token to access the repository")
	fetchCmd.Flags().StringVarP(&config.Platform, "platform", "p", "github", "The platform to fetch the Semver from, options: github")
	fetchCmd.Flags().BoolVarP(&config.Filters.Highest, "highest", "H", false, "Fetches only the highest Semver tag")
	fetchCmd.Flags().BoolVarP(&config.Filters.Release, "release", "R", false, "Fetches only Release Semver tag (x.x.x)")

	return fetchCmd
}

func RunFetchSemverTags(config *config, cmd *cobra.Command) error {

	semverSvc := newSemverSvc(config.dryRun, config.Platform, config.Token)

	klog.V(1).Info("Fetching tags...")
	semverTags, err := semverSvc.FetchSemverTags(config.Owner, config.Repository, &config.Filters)
	if err != nil {
		return err
	}

	cmd.Println(strings.Join(semverTags, " "))
	return nil
}

func newSemverSvc(dryRun bool, platform, token string) smvr.Semver {
	if len(platform) == 0 {
		platform = "github"
	}
	if dryRun {
		platform = "dry-run"
	}
	semverSvc := smvr.NewSemverSvc(platform, token)

	return semverSvc
}
