package github

import "src/pkg/fetch/models"

type GithubClient struct {
	config *models.DatasourceConfig
}

func NewFetcher(config *models.DatasourceConfig) *GithubClient {
	return &GithubClient{config: config}
}

func (g *GithubClient) FetchTags() ([]models.Version, error) {
	return []models.Version{}, nil
}
