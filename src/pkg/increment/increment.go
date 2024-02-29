package increment

import (
	"fmt"
	"src/cmd/smgr/models"
	"src/cmd/smgr/pkg/filter"
	"strconv"
)

// Get all versions (fetch)

// Promote highest prerelease to target prerelease or release stream

// Get highest release on target stream
// Get highest prerelease on target stream
// Compare existing increment to target increment, if target increment is higher, use target increment

// filter by release stream and Get highest release
// filter by prerelease stream Get highest prerelease

// Calculate increment between highest release and highest prerelease as existing increment
func CalculateIncrementTypeForNewPrerelease(highestRelease models.Version, highestPrerelease models.Version, requestedIncrement models.Increment) models.Increment {
	existingIncrement := GetIncrementType(highestRelease, highestPrerelease)
	if existingIncrement == models.None {
		return requestedIncrement
	} else if requestedIncrement.IsHigherThan(existingIncrement) {
		return requestedIncrement
	}
	return models.None
}

func GetIncrementType(highestRelease models.Version, comparedTo models.Version) models.Increment {
	if highestRelease.Release.Major < comparedTo.Release.Major {
		return models.Major
	} else if highestRelease.Release.Minor < comparedTo.Release.Minor {
		return models.Minor
	} else if highestRelease.Release.Patch < comparedTo.Release.Patch {
		return models.Patch
	}
	return models.None
}

func IncrementRelease(sourceVersion models.Version, increment models.Increment) models.Version {
	incrementedVersion := models.Version{}
	if increment == models.Major {
		incrementedVersion.Release.Major = sourceVersion.Release.Major + 1
		incrementedVersion.Release.Minor = 0
		incrementedVersion.Release.Patch = 0
	} else if increment == models.Minor {
		incrementedVersion.Release.Major = sourceVersion.Release.Major
		incrementedVersion.Release.Minor = sourceVersion.Release.Minor + 1
		incrementedVersion.Release.Patch = 0
	} else if increment == models.Patch {
		incrementedVersion.Release.Major = sourceVersion.Release.Major
		incrementedVersion.Release.Minor = sourceVersion.Release.Minor
		incrementedVersion.Release.Patch = sourceVersion.Release.Patch + 1
	}
	return incrementedVersion
}

func IncrementReleaseFromStream(sourceVersions []models.Version, streamPattern models.VersionPattern, increment models.Increment) (models.Version, error) {
	if !streamPattern.IsReleaseOnlyPattern() {
		return models.Version{}, fmt.Errorf("error: stream pattern must be release only")
	}

	sourceVersion, err := filter.GetHighestStreamVersion(sourceVersions, streamPattern)
	if err != nil {
		if _, ok := err.(*models.EmptyVersionListError); ok {
			return streamPattern.FirstVersion(), nil
		} else {
			return models.Version{}, err
		}
	}
	return IncrementRelease(sourceVersion, increment), nil
}

// func PrereleaseIncrement(sourceVersion models.Version, incrementType models.Increment) models.Version {
// 	incrementedVersion := models.Version{}
// 	if incrementType == models.Major {
// 		incrementedVersion.Release.Major = sourceVersion.Release.Major
// 		incrementedVersion.Release.Minor = sourceVersion.Release.Minor
// 		incrementedVersion.Release.Patch = sourceVersion.Release.Patch
// 		incrementedVersion.Prerelease.Major = sourceVersion.Prerelease.Major + 1
// 		incrementedVersion.Prerelease.Minor = 0
// 		incrementedVersion.Prerelease.Patch = 0
// 	} else if incrementType == models.Minor {
// 		incrementedVersion.Release.Major = sourceVersion.Release.Major
// 		incrementedVersion.Release.Minor = sourceVersion.Release.Minor
// 		incrementedVersion.Release.Patch = sourceVersion.Release.Patch
// 		incrementedVersion.Prerelease.Major = sourceVersion.Prerelease.Major
// 		incrementedVersion.Prerelease.Minor = sourceVersion.Prerelease.Minor + 1
// 		incrementedVersion.Prerelease.Patch = 0
// 	} else if incrementType == models.Patch {
// 		incrementedVersion.Release.Major = sourceVersion.Release.Major
// 		incrementedVersion.Release.Minor = sourceVersion.Release.Minor
// 		incrementedVersion.Release.Patch = sourceVersion.Release.Patch
// 		incrementedVersion.Prerelease.Major = sourceVersion.Prerelease.Major
// 		incrementedVersion.Prerelease.Minor = sourceVersion.Prerelease.Minor
// 		incrementedVersion.Prerelease.Patch = sourceVersion.Prerelease.Patch + 1
// 	}
// 	return incrementedVersion
// }

func NumericalPRIncrement(sourceIdentifier models.PRIdentifier) (models.PRIdentifier, error) {
	incrementedIdentifier := models.PRIdentifier{}
	sourceNumber, err := strconv.ParseUint(sourceIdentifier.String(), 10, 64)
	if err != nil {
		return incrementedIdentifier, err
	}
	incrementedNumber := sourceNumber + 1
	incrementedValue := strconv.FormatUint(incrementedNumber, 10)
	incrementedIdentifier.Set(incrementedValue)
	return incrementedIdentifier, nil
}

func AlphabeticalIncrement(sourceIdentifier models.PRIdentifier) (models.PRIdentifier, error) {
	incrementedIdentifier := models.PRIdentifier{}

	sourceChar := sourceIdentifier.String()
	if len(sourceChar) != 1 {
		return incrementedIdentifier, fmt.Errorf("expected a single character identifier")
	}

	sourceRunes := []rune(sourceChar)
	if (sourceRunes[0] < 'a' || sourceRunes[0] > 'z') && (sourceRunes[0] < 'A' || sourceRunes[0] > 'Z') {
		return incrementedIdentifier, fmt.Errorf("expected an alphabetical identifier")
	}

	if sourceRunes[0] == 'z' {
		sourceRunes = append(sourceRunes, 'a')
	} else if sourceRunes[0] == 'Z' {
		sourceRunes = append(sourceRunes, 'A')
	} else {
		sourceRunes[0]++
	}

	incrementedIdentifier.Set(string(sourceRunes))
	return incrementedIdentifier, nil
}

func PromoteVersion(sourceVersion models.Version, targetStream models.VersionPattern) models.Version {
	promotedVersion := models.Version{}
	// TODOS
	return promotedVersion
}
