package github

import (
	"src/pkg/fetch/models"
	"src/pkg/fetch/util"
)

type GithubClient struct {
	config *util.DatasourceConfig
}

func NewFetcher(config *util.DatasourceConfig) *GithubClient {
	return &GithubClient{config: config}
}

func (g *GithubClient) FetchTags() ([]models.Version, error) {
	return []models.Version{}, nil
}
