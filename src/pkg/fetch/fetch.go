package fetch

import (
	"errors"

	"github.com/13013SwagR/semver-manager/src/datastores/github"
	"github.com/13013SwagR/semver-manager/src/datastores/gitlab"
	"github.com/13013SwagR/semver-manager/src/datastores/oci"
	"github.com/13013SwagR/semver-manager/src/models"
)

type Fetcher interface {
	FetchTags() ([]models.Version, error)
}

type Config struct {
	Owner      string
	Repository string
	Token      string
	Platform   string
}

func NewFetcher(config *models.Config) (Fetcher, error) {
	switch config.Platform {
	case "github":
		return github.NewFetcher(config), nil
	case "gitlab":
		return gitlab.NewFetcher(config), nil
	case "oci":
		return oci.NewFetcher(config), nil
	default:
		return nil, errors.New("unsupported platform")
	}
}
