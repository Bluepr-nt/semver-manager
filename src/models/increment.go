package models

import "errors"

type Increment string

const (
	Major Increment = "major"
	Minor Increment = "minor"
	Patch Increment = "patch"
	None  Increment = "none"
)

func (i Increment) ValidateIncrement() error {
	if i != Major && i != Minor && i != Patch {
		return errors.New("invalid increment type")
	}
	return nil
}

func (i Increment) IsHigherThan(comparedTo Increment) bool {
	if i == Major {
		return comparedTo == Major
	} else if i == Minor {
		return comparedTo == Patch
	} else if i == Patch {
		return false
	}
	return false
}

func IncrementVersion(versionList []Version, i Increment, targetStream VersionPattern) error {
	if err := i.ValidateIncrement(); err != nil {
		return err
	}

	var highest Version
	for _, version := range versionList {
		if version.IsHigherThan(highest) {
			highest = version
		}
	}
	// get highest release
	return nil
}
