package testutils

import "src/cmd/smgr/models"

func NewVersion(s string) models.Version {
	v, _ := models.ParseVersion(s)
	return v
}

func NewVersionPattern(s string) models.VersionPattern {
	v, _ := models.ParseVersionPattern(s)
	return v
}
