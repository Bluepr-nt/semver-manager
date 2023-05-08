package models

type DatasourceConfig struct {
	Owner      string
	Repository string
	Token      string
	Platform   string
}
type Version struct {
	Major      uint64
	Minor      uint64
	Patch      uint64
	Prerelease PRVersion
	Build      BuildMetadata
}

func (v *Version) String() string {
	return ""
}

type PRVersion struct {
	Identifiers []PRIdentifier
}

type BuildMetadata struct {
	Identifiers []BuildIdentifier
}

type PRIdentifier struct {
	identifier string
}

func (i *PRIdentifier) String() string {
	return i.identifier
}

type BuildIdentifier struct {
	identifier string
}
