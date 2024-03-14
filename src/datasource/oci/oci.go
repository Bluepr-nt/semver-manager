package oci

import (
	"src/cmd/smgr/models"
	"src/cmd/smgr/utils"
)

type OciCLient struct {
	config *utils.DatasourceConfig
}

func NewFetcher(config *utils.DatasourceConfig) *OciCLient {
	return &OciCLient{config: config}
}

func (g *OciCLient) FetchTags() ([]models.Version, error) {
	return []models.Version{}, nil
}
