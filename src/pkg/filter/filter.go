package filter

import (
	"fmt"
	"src/pkg/fetch/models"

	"github.com/blang/semver/v4"
)

type FilterFunc func(versions []models.Version) ([]models.Version, error)

func ApplyFilters(versions []models.Version, filters ...FilterFunc) ([]models.Version, error) {
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

func ReleaseOnly() FilterFunc {
	return func(versions []models.Version) ([]models.Version, error) {
		var filtered []models.Version

		for _, version := range versions {
			if len(version.Prerelease.Identifiers) == 0 {
				filtered = append(filtered, version)
			}
		}

		return filtered, nil
	}
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

func MajorVersionStream(major uint64) FilterFunc {
	return func(versions []models.Version) ([]models.Version, error) {
		var filtered []models.Version

		for _, version := range versions {
			if version.Release.Major == major {
				filtered = append(filtered, version)
			}
		}

		return filtered, nil
	}
}

func MinorVersionStream(major, minor uint64) FilterFunc {
	return func(versions []models.Version) ([]models.Version, error) {
		var filtered []models.Version

		for _, version := range versions {
			if version.Release.Major == major && version.Release.Minor == minor {
				filtered = append(filtered, version)
			}
		}

		return filtered, nil
	}
}

func PatchVersionStream(major, minor, patch uint64) FilterFunc {
	return func(versions []models.Version) ([]models.Version, error) {
		var filtered []models.Version

		for _, version := range versions {
			if version.Release.Major == major && version.Release.Minor == minor && version.Release.Patch == patch {
				filtered = append(filtered, version)
			}
		}

		return filtered, nil
	}
}

func PreReleaseVersionStream(release models.Release, prerelease models.PRVersion) FilterFunc {
	return func(versions []models.Version) ([]models.Version, error) {
		var filtered []models.Version

		for _, version := range versions {
			if version.Release == release {
				matchingPreRelease := matchPreRelease(version.Prerelease, prerelease)
				if matchingPreRelease {
					filtered = append(filtered, version)
				}
			}
		}

		return filtered, nil
	}
}

func matchPreRelease(prVersion models.PRVersion, prerelease models.PRVersion) bool {
	if len(prVersion.Identifiers) < len(prerelease.Identifiers) {
		return false
	}

	for i, prIdentifier := range prerelease.Identifiers {
		if prVersion.Identifiers[i].String() != prIdentifier.String() {
			return false
		}
	}

	return true
}
