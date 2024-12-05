package increment

import (
	"fmt"
	"src/cmd/smgr/models"
	"src/cmd/smgr/pkg/filter"
	"src/cmd/smgr/utils"
	"strconv"
)

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

func IncrementReleaseToStream(sourceVersions []models.Version, streamPattern models.VersionPattern, increment models.Increment) (models.Version, error) {
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

func IncrementPReleaseToStream(sourceVersion models.Version, targetStream models.VersionPattern, versionList []models.Version) (models.Version, error) {
	if targetStream.IsReleaseOnlyPattern() {
		return models.Version{}, fmt.Errorf("error: stream pattern must be prerelease only")
	}

	newVersion := promoteToTargetStream(targetStream, sourceVersion)
	versionList = append(versionList, sourceVersion)
	highestStreamVersion, err := filter.GetHighestStreamVersion(versionList, targetStream)

	if highestStreamVersion.IsHigherThan(newVersion) && isStreamEmpty(err) {

		highestStreamVersion.Prerelease = PrereleaseIncrement(highestStreamVersion.Prerelease)
		newVersion = highestStreamVersion
	} else if highestStreamVersion.IsEqualTo(newVersion) {

		newVersion.Prerelease = PrereleaseIncrement(newVersion.Prerelease)
	}
	return newVersion, nil
}

func isStreamEmpty(err error) bool {
	_, ok := err.(*models.EmptyVersionListError)
	return !ok
}

func PrereleaseIncrement(prVersion models.PRVersion) models.PRVersion {

	incrementedIdentifier, err := PRIdentifierIncrement(prVersion.LastID())
	if err != nil {

		newId, _ := models.ParsePrIdentifier("0")
		prVersion.Identifiers = append(prVersion.Identifiers, newId)
	} else {

		prVersion.Identifiers[len(prVersion.Identifiers)-1] = incrementedIdentifier
	}

	return prVersion
}

func PRIdentifierIncrement(sourceId models.PRIdentifier) (newId models.PRIdentifier, err error) {
	if utils.IsNumerical(sourceId.Value()) {
		newId, err = NumericalPRIncrement(sourceId)
		if err != nil {
			return models.PRIdentifier{}, err
		}
	} else {
		newId, err = AlphabeticalIncrement(sourceId)
		if err != nil {
			return models.PRIdentifier{}, err
		}

	}
	return newId, nil
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

func promoteToTargetStream(targetStream models.VersionPattern, sourceVersion models.Version) (promotedVersion models.Version) {

	promotedVersion.Release = promoteToTargetStreamRelease(targetStream.Release, sourceVersion.Release)

	if !targetStream.IsReleaseOnlyPattern() {
		promotedVersion.Prerelease = promoteToTargetPrereleaseStream(targetStream.Prerelease.Identifiers, sourceVersion, promotedVersion)
	}

	return promotedVersion
}

func promoteToTargetPrereleaseStream(targetStream []models.PRIdentifierPattern, sourceVersion models.Version, promotedVersion models.Version) models.PRVersion {

	for _, targetId := range targetStream {

		if targetId.Value() == models.Wildcard {
			newId := models.PRIdentifier{}
			newId.Set("0")
			promotedVersion.Prerelease.Identifiers = append(promotedVersion.Prerelease.Identifiers, newId)

		} else {
			newId := models.PRIdentifier{}
			newId.Set(targetId.Value())
			promotedVersion.Prerelease.Identifiers = append(promotedVersion.Prerelease.Identifiers, newId)
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
