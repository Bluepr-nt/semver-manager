package gitlab

import "src/pkg/fetch/models"

type GitlabClient struct {
	config *models.DatasourceConfig
}

func NewFetcher(config *models.DatasourceConfig) *GitlabClient {
	return &GitlabClient{config: config}
}

func (g *GitlabClient) FetchTags() ([]models.Version, error) {
	return []models.Version{}, nil
}
