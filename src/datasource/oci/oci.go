package oci

import (
	"src/cmd/smgr/models"
	"src/cmd/smgr/util"
)

type OciCLient struct {
	config *util.DatasourceConfig
}

func NewFetcher(config *util.DatasourceConfig) *OciCLient {
	return &OciCLient{config: config}
}

func (g *OciCLient) FetchTags() ([]models.Version, error) {
	return []models.Version{}, nil
}
