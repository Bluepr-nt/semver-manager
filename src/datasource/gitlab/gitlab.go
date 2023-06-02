package gitlab

import (
	"src/cmd/smgr/models"
	"src/cmd/smgr/util"
)

type GitlabClient struct {
	config *util.DatasourceConfig
}

func NewFetcher(config *util.DatasourceConfig) *GitlabClient {
	return &GitlabClient{config: config}
}

func (g *GitlabClient) FetchTags() ([]models.Version, error) {
	return []models.Version{}, nil
}
