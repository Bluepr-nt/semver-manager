package filter

import (
	"fmt"
	"src/pkg/fetch/models"
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

// matchPrerelease checks if the prerelease identifiers match the pattern.
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

// toUint converts a string to a uint64.
func toUint(s string) uint64 {
	v, _ := strconv.ParseUint(s, 10, 64)
	return v
}

// create filter function where * is a wildcard for major, minor, patch and prerelease identifiers and exclude prerelease versions if no prerelease identifiers are specified or * inplace of prerelease identifiers
// func VersionPatternFilterWithWildcard(pattern models.VersionPattern) FilterFunc {
// 	return func(versions []models.Version) ([]models.Version, error) {
// 		var filtered []models.Version

// 		for _, version := range versions {
// 			if pattern.Release.Major.Value() != "*" && version.Release.Major != toUint(pattern.Release.Major.Value()) {
// 				continue
// 			}
// 			if pattern.Release.Minor.Value() != "*" && version.Release.Minor != toUint(pattern.Release.Minor.Value()) {
// 				continue
// 			}
// 			if pattern.Release.Patch.Value() != "*" && version.Release.Patch != toUint(pattern.Release.Patch.Value()) {
// 				continue
// 			}
// 			if len(pattern.Prerelease.Identifiers) == 0 && len(version.Prerelease.Identifiers) > 0 {
// 				// If no prerelease data in pattern, exclude versions that have prerelease identifiers
// 				continue
// 			}
// 			if len(pattern.Prerelease.Identifiers) > 0 && !matchPrerelease(pattern.Prerelease.Identifiers, version.Prerelease) {
// 				continue
// 			}

// 			filtered = append(filtered, version)
// 		}

// 		return filtered, nil
// 	}
// }

func VersionPatternFilterWithWildcard(pattern models.VersionPattern) FilterFunc {
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

func matchPrereleaseWithWildcard(prIdentifiersPattern []models.PRIdentifierPattern, prerelease models.PRVersion) bool {
	if len(prIdentifiersPattern) > len(prerelease.Identifiers) {
		return false
	}

	for i, prIdentifierPattern := range prIdentifiersPattern {
		if prIdentifierPattern.Value() != "*" && prIdentifierPattern.Value() != prerelease.Identifiers[i].String() {
			return false
		}
	}

	return true
}
