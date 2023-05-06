package models

type Version struct {
	Major      uint64
	Minor      uint64
	Patch      uint64
	Prerelease PRVersion
	Build      BuildMetadata
}

type PRVersion struct {
	Identifiers []PRIdentifier
}

type BuildMetadata {
	Identifiers []BuildIdentifier
}

type PRIdentifier {
	identifier string
}

type BuildIdentifier struct {
	identifier string
}
