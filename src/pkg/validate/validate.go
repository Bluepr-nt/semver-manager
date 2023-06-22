package validate

import (
	"errors"
	"strings"

	"github.com/blang/semver/v4"
)

type Validator interface {
	IsValid(version string) (bool, error)
}

type LooseValidator struct{}
type StrictValidator struct{}

func NewSemverValidator(vType string) Validator {
	if vType == "loose" {
		return &LooseValidator{}
	} else if vType == "strict" || vType == "" {
		return &StrictValidator{}
	}

	return nil
}

func (v *LooseValidator) IsValid(version string) (bool, error) {
	version = strings.TrimPrefix(version, "v")
	_, err := semver.Parse(version)
	if err != nil {
		return false, errors.New("invalid semver")
	}

	return true, nil
}

func (v *StrictValidator) IsValid(version string) (bool, error) {
	_, err := semver.Parse(version)
	if err != nil {
		return false, errors.New("invalid semver")
	}

	return true, nil
}

func IsSemverValid(version string, validator Validator) (bool, error) {
	return validator.IsValid(version)
}
