package fetch

import (
	"errors"

	"src/pkg/fetch/datasource/github"
	"src/pkg/fetch/datasource/gitlab"
	"src/pkg/fetch/datasource/oci"
	"src/pkg/fetch/models"
)

type Fetcher interface {
	FetchTags() ([]models.Version, error)
}

func NewFetcher(config *models.DatasourceConfig) (Fetcher, error) {
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
