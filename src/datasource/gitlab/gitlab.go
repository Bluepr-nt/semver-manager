package gitlab

import (
	"src/cmd/smgr/models"
	"src/cmd/smgr/utils"
)

type GitlabClient struct {
	config *utils.DatasourceConfig
}

func NewFetcher(config *utils.DatasourceConfig) *GitlabClient {
	return &GitlabClient{config: config}
}

func (g *GitlabClient) FetchTags() ([]models.Version, error) {
	return []models.Version{}, nil
}
