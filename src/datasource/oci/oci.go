package oci

import (
	"src/pkg/fetch/models"
	"src/pkg/fetch/util"
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
