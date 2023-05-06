package fetch

import (
	"errors"

	"src/datasource/github"
	"src/datasource/gitlab"
	"src/datasource/oci"
)

type Fetcher interface {
	FetchTags() ([]string, error)
}

type Config struct {
	Owner      string
	Repository string
	Token      string
	Platform   string
}

func NewFetcher(config *Config) (Fetcher, error) {
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
