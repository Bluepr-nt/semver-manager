package oci

import "src/pkg/fetch/models"

type OciCLient struct {
	config *models.DatasourceConfig
}

func NewFetcher(config *models.DatasourceConfig) *OciCLient {
	return &OciCLient{config: config}
}

func (g *OciCLient) FetchTags() ([]models.Version, error) {
	return []models.Version{}, nil
}
