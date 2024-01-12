package models

import "errors"

type Increment string

const (
	Major Increment = "major"
	Minor Increment = "minor"
	Patch Increment = "patch"
	None  Increment = "none"
)

func (i Increment) Validate() error {
	if i != Major && i != Minor && i != Patch {
		return errors.New("invalid increment type")
	}
	return nil
}

func (i Increment) IsHigherThan(comparedTo Increment) bool {
	if i == Major && comparedTo == Major {
		return false
	} else if i == Minor && comparedTo != Major {
		return true
	} else if i == Patch && comparedTo == Patch {
		return true
	}
	return false
}
