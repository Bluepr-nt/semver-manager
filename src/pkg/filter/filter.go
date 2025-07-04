package filter

import (
	"src/cmd/smgr/models"
	"strconv"

	"github.com/blang/semver/v4"
)

type FilterFunc func(versions []models.Version) ([]models.Version, error)

func ApplyFilters(versions []models.Version, filters ...FilterFunc) (models.VersionSlice, error) {
	var err error
	filtered := versions

	for _, filter := range filters {
		filtered, err = filter(filtered)
		if err != nil {
			return nil, err
		}
	}

	return filtered, nil
}

func Highest() FilterFunc {
	return func(versions []models.Version) ([]models.Version, error) {
		if len(versions) == 0 {
			return versions, &models.EmptyVersionListError{}
		}

		highest := versions[0]
		for _, version := range versions[1:] {
			semverVersion, _ := semver.ParseTolerant(version.String())
			highestSemver, _ := semver.ParseTolerant(highest.String())
			semver.Parse("")
			if semverVersion.GT(highestSemver) {
				highest = version
			}
		}

		return []models.Version{highest}, nil
	}
}

func GetHighestStreamVersion(versions []models.Version, streamPattern models.VersionPattern) (models.Version, error) {
	var err error
	streamFilter := VersionPatternFilter(streamPattern)
	highestFilter := Highest()
	sourceVersion, err := ApplyFilters(versions, streamFilter, highestFilter)
	if err != nil {
		return models.Version{}, err
	}
	return sourceVersion[0], nil
}

func GetHighestStreamVersionWithReleases(versions []models.Version, streamPattern models.VersionPattern) (models.Version, error) {

	streamFilter := VersionPatternFilter(streamPattern)
	sourceVersions, err := ApplyFilters(versions, streamFilter)
	if err != nil {
		return models.Version{}, err
	}

	if !streamPattern.Release.IsStrict() {
		releaseOnlyFilter := VersionPatternFilter(models.VersionPattern{Release: streamPattern.Release, Prerelease: models.PRVersionPattern{}, Build: streamPattern.Build})
		sourceReleaseVersion, err := ApplyFilters(versions, releaseOnlyFilter)
		if err != nil {
			return models.Version{}, err
		}

		sourceVersions = append(sourceVersions, sourceReleaseVersion...)
	}

	highestFilter := Highest()
	highestVersion, err := ApplyFilters(sourceVersions, highestFilter)
	if err != nil {
		return models.Version{}, err
	}

	return highestVersion[0], nil
}

// VersionPatternFilter returns a filter function that
// filters versions based on a VersionPattern
func VersionPatternFilter(pattern models.VersionPattern) FilterFunc {

	return func(versions []models.Version) ([]models.Version, error) {

		var filtered []models.Version
		releaseFilter := ReleasePatternFilter(pattern.Release)
		prereleaseFilter := PrereleasePatternFilter(pattern.Prerelease)
		buildMetadataFilter := BuildMetadataFilter(pattern.Build)

		filtered, err := ApplyFilters(versions, releaseFilter, prereleaseFilter, buildMetadataFilter)
		if err != nil {
			return nil, err
		}

		return filtered, nil
	}
}

// ReleasePatternFilter returns a filter function that
// filters versions based on a ReleasePattern
// all release and prelease versions are returned
func ReleasePatternFilter(pattern models.ReleasePattern) FilterFunc {
	return func(versions []models.Version) ([]models.Version, error) {
		var filtered []models.Version

		for _, version := range versions {
			if pattern.Major.Value() != "*" && version.Release.Major.String() != pattern.Major.Value() {
				continue
			}
			if pattern.Minor.Value() != "*" && version.Release.Minor.String() != pattern.Minor.Value() {
				continue
			}
			if pattern.Patch.Value() != "*" && version.Release.Patch.String() != pattern.Patch.Value() {
				continue
			}

			filtered = append(filtered, version)
		}

		return filtered, nil
	}
}

// PrereleasePatternFilter returns a filter function that
// filters versions based on a PrereleasePattern
// prerelease versions are returned
// all release versions are also returned
func PrereleasePatternFilter(pattern models.PRVersionPattern) FilterFunc {
	return func(versions []models.Version) ([]models.Version, error) {

		var filtered []models.Version

		for _, version := range versions {
			if !matchPrerelease(pattern.Identifiers, version.Prerelease) {
				continue
			}

			filtered = append(filtered, version)
		}

		return filtered, nil
	}
}

// GetValidVersions returns a list of valid versions from one or more string lists of versions
// The versions are split by comma or space and then parsed into a Version struct
func GetValidVersions(stringVersionsList ...string) []models.Version {
	var versions []models.Version
	for _, stringVersions := range stringVersionsList {
		for _, stringVersion := range models.SplitVersions(stringVersions) {
			version, _ := models.ParseVersion(stringVersion)
			versions = append(versions, version)
		}
	}
	return versions
}

func matchPrerelease(prIdentifiersPattern []models.PRIdentifierPattern, prerelease models.PRVersion) bool {
	if len(prIdentifiersPattern) != len(prerelease.Identifiers) {
		return false
	}

	for i, prIdentifierPattern := range prIdentifiersPattern {
		if prIdentifierPattern.Value() != "*" && prIdentifierPattern.Value() != prerelease.Identifiers[i].Value() {
			return false
		}
	}

	return true
}

func toUint(s string) uint64 {
	v, _ := strconv.ParseUint(s, 10, 64)
	return v
}

func BuildMetadataFilter(pattern models.BuildMetadataPattern) FilterFunc {
	return func(versions []models.Version) ([]models.Version, error) {
		var filtered []models.Version

		for _, version := range versions {
			if len(pattern.Identifiers) == 0 {
				filtered = append(filtered, version)
				continue
			}

			if !matchBuildMetadata(pattern.Identifiers, version.BuildMetadata) {
				continue
			}

			filtered = append(filtered, version)
		}

		return filtered, nil
	}
}

func matchBuildMetadata(buildIdentifiersPattern []models.BuildIdentifierPattern, buildMetadata models.BuildMetadata) bool {
	if len(buildIdentifiersPattern) != len(buildMetadata.Identifiers) {
		return false
	}

	for i, buildIdentifierPattern := range buildIdentifiersPattern {
		if buildIdentifierPattern.Value() != "*" && buildIdentifierPattern.Value() != buildMetadata.Identifiers[i].String() {
			return false
		}
	}

	return true
}
