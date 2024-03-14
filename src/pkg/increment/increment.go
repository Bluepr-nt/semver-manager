package increment

import (
	"fmt"
	"src/cmd/smgr/models"
	"src/cmd/smgr/pkg/filter"
	"src/cmd/smgr/utils"
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
	if highestRelease.Release.Major.LT(comparedTo.Release.Major) {
		return models.Major
	} else if highestRelease.Release.Minor.LT(comparedTo.Release.Minor) {
		return models.Minor
	} else if highestRelease.Release.Patch.LT(comparedTo.Release.Patch) {
		return models.Patch
	}
	return models.None
}

func IncrementRelease(sourceVersion models.Version, increment models.Increment) models.Version {
	incrementedVersion := sourceVersion
	if increment == models.Major {
		incrementedVersion.Release.Major.Set(sourceVersion.Release.Major.Value() + 1)
		incrementedVersion.Release.Minor.Set(0)
		incrementedVersion.Release.Patch.Set(0)
	} else if increment == models.Minor {
		incrementedVersion.Release.Minor.Set(sourceVersion.Release.Minor.Value() + 1)
		incrementedVersion.Release.Patch.Set(0)
	} else if increment == models.Patch {
		incrementedVersion.Release.Patch.Set(sourceVersion.Release.Patch.Value() + 1)
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

func PRIdentifierIncrement(sourceId models.PRIdentifier) (newId models.PRIdentifier) {
	if utils.IsNumerical(sourceId.Value()) {
		newId, _ = NumericalPRIncrement(sourceId)
	} else {
		newId, _ = AlphabeticalIncrement(sourceId)
	}
	return newId
}

func NumericalPRIncrement(sourceIdentifier models.PRIdentifier) (models.PRIdentifier, error) {
	incrementedIdentifier := models.PRIdentifier{}
	sourceNumber, err := strconv.ParseUint(sourceIdentifier.Value(), 10, 64)
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

	sourceChar := sourceIdentifier.Value()
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

func PromotePRVersion(sourceVersion models.Version, targetStream models.VersionPattern, versionList []models.Version) (promotedVersion models.Version) {
	// sourceVersion needs to be a Prerelease version
	// TODO validate

	// Calculate increment of current version vs highest release on target stream
	// increment := GetIncrementType(highestStreamVersion, sourceVersion)

	// Translate version to target stream
	promotedToStream := promoteToTargetStream(targetStream, sourceVersion)

	highestStreamVersion, err := filter.GetHighestStreamVersion(versionList, targetStream)
	if _, ok := err.(*models.EmptyVersionListError); ok {
		if utils.IsNumerical(promotedToStream.Prerelease.LastID().Value()) {
			promotedToStream.Prerelease.Identifiers[len(promotedToStream.Prerelease.Identifiers)-1].Set("0")
			promotedVersion = promotedToStream
		}
	} else {
		promotedVersion = promoteVersionAbove(highestStreamVersion, promotedToStream)
	}

	return promotedVersion
}

func promoteVersionAbove(highestStreamVersion models.Version, promotedToStream models.Version) models.Version {
	for i, identifier := range highestStreamVersion.Prerelease.Identifiers {
		// if current index is bigger than promoted stream ids
		if (len(promotedToStream.Prerelease.Identifiers) - 1) < i {
			promotedToStream.Prerelease.Identifiers = append(promotedToStream.Prerelease.Identifiers, highestStreamVersion.Prerelease.Identifiers[i])

			// if current id is bigger than promoted version
		} else if identifier.IsHigherThan(promotedToStream.Prerelease.Identifiers[i]) {
			promotedToStream.Prerelease.Identifiers[i] = identifier
		}

		//  if current id is the last
		if (len(highestStreamVersion.Prerelease.Identifiers) - 1) == i {
			if !promotedToStream.IsHigherThan(highestStreamVersion) {
				newId := PRIdentifierIncrement(identifier)
				promotedToStream.Prerelease.Identifiers[i] = newId
			}
		}
	}

	return promotedToStream
}

func promoteToTargetStream(targetStream models.VersionPattern, sourceVersion models.Version) (promotedVersion models.Version) {

	promotedVersion.Release = promoteToTargetStreamRelease(targetStream.Release, sourceVersion.Release)

	if !targetStream.IsReleaseOnlyPattern() {
		promotedVersion.Prerelease = promoteToTargetPrereleaseStream(targetStream.Prerelease.Identifiers, sourceVersion, promotedVersion)
	}

	return promotedVersion
}

func promoteToTargetPrereleaseStream(targetStream []models.PRIdentifierPattern, sourceVersion models.Version, promotedVersion models.Version) models.PRVersion {

	for i, targetId := range targetStream {
		if targetId.Value() != models.Wildcard {
			newId := models.PRIdentifier{}
			newId.Set(targetId.Value())
			promotedVersion.Prerelease.Identifiers = append(promotedVersion.Prerelease.Identifiers, newId)

		} else if (len(sourceVersion.Prerelease.Identifiers) - 1) < i {
			newId := models.PRIdentifier{}
			newId.Set("0")
			promotedVersion.Prerelease.Identifiers = append(promotedVersion.Prerelease.Identifiers, newId)

		} else if (len(targetStream) - 1) >= i {
			promotedVersion.Prerelease.Identifiers = append(promotedVersion.Prerelease.Identifiers, sourceVersion.Prerelease.Identifiers[i])

		}
	}
	return promotedVersion.Prerelease
}

func promoteToTargetStreamRelease(targetStreamRelease models.ReleasePattern, sourceRelease models.Release) (newRelease models.Release) {

	newRelease.Major = promotoToStreamReleaseDigit(targetStreamRelease.Major, sourceRelease.Major)
	newRelease.Minor = promotoToStreamReleaseDigit(targetStreamRelease.Minor, sourceRelease.Minor)
	newRelease.Patch = promotoToStreamReleaseDigit(targetStreamRelease.Patch, sourceRelease.Patch)

	return newRelease
}

func promotoToStreamReleaseDigit(targetStream models.ReleaseDigitPattern, sourceDigit models.ReleaseDigit) (newDigit models.ReleaseDigit) {
	if targetStream.Value() == models.Wildcard {
		newDigit.Set(sourceDigit.Value())
	} else {
		var rawDigit uint64
		rawDigit, _ = strconv.ParseUint(targetStream.Value(), 10, 64)
		newDigit.Set(rawDigit)
	}
	return newDigit
}
