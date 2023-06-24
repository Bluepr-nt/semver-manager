package increment

import (
	"fmt"
	"src/cmd/smgr/models"
	"src/cmd/smgr/pkg/filter"
	"strconv"
)

// Get all versions (fetch)

// Calculate increment between highest release and highest prerelease as existing increment

// Promote highest prerelease to target prerelease or release stream
// Get highest release on target stream
// Get highest prerelease on target stream
// Compare existing increment to target increment, if target increment is higher, use target increment

// filter by release stream and Get highest release
// filter by prerelease stream Get highest prerelease
func GetHighestStreamVersion(versions []models.Version, streamPattern models.VersionPattern) (models.Version, error) {
	var err error
	streamFilter := filter.VersionPatternFilter(streamPattern)
	highestFilter := filter.Highest()
	sourceVersion, err := filter.ApplyFilters(versions, streamFilter, highestFilter)
	if err != nil {
		return models.Version{}, err
	}
	return sourceVersion[0], nil
}

func GetReleaseToPreReleaseIncrementType(highestRelease models.Version, highestPrerelease models.Version) string {
	if highestRelease.Release.Major < highestPrerelease.Release.Major {
		return models.Major
	} else if highestRelease.Release.Minor < highestPrerelease.Release.Minor {
		return models.Minor
	} else if highestRelease.Release.Patch < highestPrerelease.Release.Patch {
		return models.Patch
	}
	return ""
}

func ReleaseIncrement(sourceVersion models.Version, incrementType string) models.Version {
	incrementedVersion := models.Version{}
	if incrementType == models.Major {
		incrementedVersion.Release.Major = sourceVersion.Release.Major + 1
		incrementedVersion.Release.Minor = 0
		incrementedVersion.Release.Patch = 0
	} else if incrementType == models.Minor {
		incrementedVersion.Release.Major = sourceVersion.Release.Major
		incrementedVersion.Release.Minor = sourceVersion.Release.Minor + 1
		incrementedVersion.Release.Patch = 0
	} else if incrementType == models.Patch {
		incrementedVersion.Release.Major = sourceVersion.Release.Major
		incrementedVersion.Release.Minor = sourceVersion.Release.Minor
		incrementedVersion.Release.Patch = sourceVersion.Release.Patch + 1
	}
	return incrementedVersion
}

// func PrereleaseIncrement(sourceVersion models.Version, incrementType string) models.Version {

func NumericalIncrement(sourceIdentifier models.PRIdentifier) (models.PRIdentifier, error) {
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
		return incrementedIdentifier, fmt.Errorf("Expected a single character identifier")
	}

	sourceRunes := []rune(sourceChar)
	if sourceRunes[0] < 'a' || sourceRunes[0] > 'z' {
		return incrementedIdentifier, fmt.Errorf("Expected an alphabetical identifier")
	}

	if sourceRunes[0] == 'z' {
		sourceRunes[0] = 'a'
	} else {
		sourceRunes[0]++
	}

	incrementedIdentifier.Set(string(sourceRunes))
	return incrementedIdentifier, nil
}

func PromoteVersion(sourceVersion models.Version, targetStream models.VersionPattern) models.Version {
	promotedVersion := models.Version{}

	return promotedVersion
}
