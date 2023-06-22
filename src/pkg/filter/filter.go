package filter

import (
	"fmt"
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
			return versions, fmt.Errorf("error version list is empty")
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

func VersionPatternFilter(pattern models.VersionPattern) FilterFunc {
	return func(versions []models.Version) ([]models.Version, error) {
		var filtered []models.Version

		for _, version := range versions {
			if pattern.Release.Major.Value() != "*" && version.Release.Major != toUint(pattern.Release.Major.Value()) {
				continue
			}
			if pattern.Release.Minor.Value() != "*" && version.Release.Minor != toUint(pattern.Release.Minor.Value()) {
				continue
			}
			if pattern.Release.Patch.Value() != "*" && version.Release.Patch != toUint(pattern.Release.Patch.Value()) {
				continue
			}
			if len(pattern.Prerelease.Identifiers) > 0 && !matchPrerelease(pattern.Prerelease.Identifiers, version.Prerelease) {
				continue
			} else if len(pattern.Prerelease.Identifiers) == 0 && len(version.Prerelease.Identifiers) > 0 {
				continue
			}
			// Checking build metadata
			if len(pattern.Build.Identifiers) > 0 && !matchBuildMetadata(pattern.Build.Identifiers, version.BuildMetadata) {
				continue
			}

			filtered = append(filtered, version)
		}

		return filtered, nil
	}
}

func GetVersions(stringVersions string) []models.Version {
	var versions []models.Version

	for _, stringVersion := range models.SplitVersions(stringVersions) {
		version, _ := models.ParseVersion(stringVersion)
		versions = append(versions, version)
	}

	return versions
}

func matchPrerelease(prIdentifiersPattern []models.PRIdentifierPattern, prerelease models.PRVersion) bool {
	if len(prIdentifiersPattern) != len(prerelease.Identifiers) {
		return false
	}

	for i, prIdentifierPattern := range prIdentifiersPattern {
		if prIdentifierPattern.Value() != "*" && prIdentifierPattern.Value() != prerelease.Identifiers[i].String() {
			return false
		}
	}

	return true
}

func toUint(s string) uint64 {
	v, _ := strconv.ParseUint(s, 10, 64)
	return v
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
