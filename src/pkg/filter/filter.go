package filter

import (
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
			return versions, nil
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
			if version.Major == major {
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
			if version.Major == major && version.Minor == minor {
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
			if version.Major == major && version.Minor == minor && version.Patch == patch {
				filtered = append(filtered, version)
			}
		}

		return filtered, nil
	}
}

func PreReleaseVersionStream(major, minor, patch uint64, prerelease []string) FilterFunc {
	return func(versions []models.Version) ([]models.Version, error) {
		var filtered []models.Version

		for _, version := range versions {
			if version.Major == major && version.Minor == minor && version.Patch == patch {
				matchingPreRelease := matchPreRelease(version.Prerelease, prerelease)
				if matchingPreRelease {
					filtered = append(filtered, version)
				}
			}
		}

		return filtered, nil
	}
}

func matchPreRelease(prVersion models.PRVersion, prerelease []string) bool {
	if len(prVersion.Identifiers) < len(prerelease) {
		return false
	}

	for i, prIdentifier := range prerelease {
		if prVersion.Identifiers[i].String() != prIdentifier {
			return false
		}
	}

	return true
}
