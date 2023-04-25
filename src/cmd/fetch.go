package fetch

import (
	"fmt"
	"io"

	"github.com/spf13/cobra"
	"k8s.io/klog/v2"
	"github.com/13013SwagR/semver-manager/src/internal/utils"
	services "github.com/13013SwagR/semver-manager/src/pkg"
)

type config struct {
	Token      string `san:"trim"`
	Repository string `san:"trim"`
	Owner      string `san:"trim"`
	Platform   string `san:"trim"`
	dryRun     bool
}

func NewFetchCommand(output io.Writer) *cobra.Command {
	config := config{}
	var fetchCmd = &cobra.Command{
		Use:   "fetch",
		Short: "Fetch semver tags from a registry or repository.",
		Long: `Fetch semver tags from a registry or repository and
    sorted by highest version first. Includes tags starting 'v' as in 'v0.0.0`,
		PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
			return utils.InitializeConfig(cmd)
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := utils.SanitizeInputs(&config); err != nil {
				klog.Errorf("CLI argument error: %w", err)
				panic(config)
			}
			dryRun, _ := cmd.PersistentFlags().GetBool("dry-run")
			config.dryRun = dryRun

			if enabled, _ := cmd.Flags().GetBool("highest"); enabled {
				err := RunFetchHighestSemver(config)
				if err != nil {
					return fmt.Errorf("unable to fetch highest semver error: %w", err)
				}
				return nil
			}

			return RunFetchSemverTags(config)
		},
	}
	fetchCmd.Flags().StringVarP(&config.Owner, "owner", "o", "", "The owner of the registry or repository")
	fetchCmd.Flags().StringVarP(&config.Repository, "repo", "r", "", "The repository or registry to fetch the Semver tags from")
	fetchCmd.Flags().StringVarP(&config.Token, "token", "t", "", "The token to access the repository")
	fetchCmd.Flags().StringVarP(&config.Platform, "platform", "p", "github", "The platform to fetch the Semver from, options: github")
	fetchCmd.Flags().BoolP("highest", "H", false, "Fetches only the highest Semver tag")

	return fetchCmd
}

func RunFetchSemverTags(config config) error {
	semverSvc := newSemverSvc(config.dryRun, config.Platform, config.Token)
	klog.V(1).Info("Fetching tags...")
	semverTags, err := semverSvc.FetchSemverTags(config.Owner, config.Repository)

	if err != nil {
		return err
	}
	fmt.Println(semverTags)
	return nil
}

func RunFetchHighestSemver(config config) error {
	semverSvc := newSemverSvc(config.dryRun, config.Platform, config.Token)
	klog.V(1).Info("Fetching highest tag...")
	semver, err := semverSvc.FetchHighestSemver(config.Owner, config.Repository)
	if err != nil {
		return err
	}
	fmt.Println(semver)
	return nil
}
func newSemverSvc(dryRun bool, platform, token string) services.SemverSvc {
	if len(platform) == 0 {
		platform = "github"
	}
	if dryRun {
		platform = "dry-run"
	}
	semverSvc := services.NewSemverSvc(platform, token)

	return semverSvc
}
