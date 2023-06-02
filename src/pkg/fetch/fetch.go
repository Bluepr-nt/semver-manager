package fetch

import (
	"errors"

	"src/cmd/smgr/datasource/github"
	"src/cmd/smgr/datasource/gitlab"
	"src/cmd/smgr/datasource/oci"
	"src/cmd/smgr/models"
	"src/cmd/smgr/util"
)

type Fetcher interface {
	FetchTags() ([]models.Version, error)
}

func NewFetcher(config *util.DatasourceConfig) (Fetcher, error) {
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
