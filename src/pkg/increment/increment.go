package increment

import (
	"fmt"
	"src/cmd/smgr/models"
	"src/cmd/smgr/pkg/filter"
	"src/cmd/smgr/utils"
	"strconv"
)

// TODO write test
func IncrementVersion(sourceVersions []models.Version, streamPattern models.VersionPattern, increment models.Increment) (incrementedVersion models.Version, err error) {
	if streamPattern.IsEmpty() {
		streamPattern, _ = models.ParseVersionPattern("*.*.*")
	}

	if streamPattern.IsPRPattern() {
		incrementedVersion, err = IncrementPReleaseToStream(sourceVersions, streamPattern, increment)

	} else if streamPattern.IsReleaseOnlyPattern() {
		incrementedVersion, err = IncrementReleaseToStream(sourceVersions, streamPattern, increment)

	} else {
		err = fmt.Errorf("error: stream pattern or increment level is required") // TODO default to patch increment

	}

	return incrementedVersion, err
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

func IncrementPReleaseToStream(sourceVersions []models.Version, streamPattern models.VersionPattern, increment models.Increment) (models.Version, error) {
	newVersion := models.Version{}
	if streamPattern.IsReleaseOnlyPattern() {
		return models.Version{}, fmt.Errorf("error: stream pattern must be prerelease only")
	}

	highestStreamVersion, err := filter.GetHighestStreamVersionWithReleases(sourceVersions, streamPattern)
	if err != nil {
		if isStreamEmpty(err) {
			return streamPattern.FirstVersion(), nil
		}
		return models.Version{}, err
	}

	streamVersion := streamPattern.FirstVersion()

	if !streamVersion.IsHigherThan(highestStreamVersion) {
		if !streamVersion.Release.IsHigherThan(highestStreamVersion.Release) && increment == models.None && highestStreamVersion.IsRelease() {
			increment = models.Patch
		}
		newVersion = IncrementRelease(highestStreamVersion, increment)

		if !streamVersion.Prerelease.IsHigherThan(highestStreamVersion.Prerelease) {
			newVersion.Prerelease = PrereleaseIncrement(highestStreamVersion.Prerelease)
		} else {
			newVersion.Prerelease = streamVersion.Prerelease
		}

	} else {
		newVersion = streamVersion
	}

	return newVersion, nil
}

func isStreamEmpty(err error) bool {
	_, ok := err.(*models.EmptyVersionListError)
	return ok
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
