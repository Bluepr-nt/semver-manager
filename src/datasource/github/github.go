package github

import (
	"src/cmd/smgr/models"
	"src/cmd/smgr/utils"
)

type GithubClient struct {
	config *utils.DatasourceConfig
}

func NewFetcher(config *utils.DatasourceConfig) *GithubClient {
	return &GithubClient{config: config}
}

func (g *GithubClient) FetchTags() ([]models.Version, error) {
	return []models.Version{}, nil
}
