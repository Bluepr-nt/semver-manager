package models

import (
	"fmt"
	"strconv"
	"strings"
)

const (
	numbers  string = "0123456789"
	alphas          = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ-"
	alphanum        = alphas + numbers
)

type DatasourceConfig struct {
	Owner      string
	Repository string
	Token      string
	Platform   string
}
type Version struct {
	Major      uint64
	Minor      uint64
	Patch      uint64
	Prerelease PRVersion
	Build      BuildMetadata
}

func (v *Version) String() string {
	return ""
}

func NewVersion(v string) (Version, error) {
	tokens := strings.SplitN(v, ".", 1)
	major := tokens[0]
	if !containsOnly(tokens[0], numbers) {
		return Version{}, fmt.Errorf("major version MUST be non-negative integer, got: %s", tokens[0])
	}

	if len(tokens[0]) > 1 && tokens[0][0] == '0' {
		return Version{}, fmt.Errorf("major version MUST NOT contain leading zeroes, got: %s", tokens[0])
	}

	majorUint, err := strconv.ParseUint(tokens[0], 10, 64)
	if err != nil {
		return Version{}, fmt.Errorf("error converting major version number to uint64: %w", err)
	}

	tokens = strings.SplitN(tokens[1], ".", 1)
	if !containsOnly(tokens[0], numbers) {
		return Version{}, fmt.Errorf("minor version MUST be non-negative integer, got: %s", tokens[0])
	}
	minorUint, err := strconv.ParseUint(tokens[0], 10, 64)
	if len(tokens[0]) > 1 && tokens[0][0] == '0' {
		return Version{}, fmt.Errorf("minor version MUST NOT contain leading zeroes, got: %s", tokens[0])
	}

	tokens = strings.SplitN(tokens[1], "-", 1)
	if !containsOnly(tokens[0], numbers) {
		return Version{}, fmt.Errorf("patch version MUST be non-negative integer, got: %s", tokens[0])
	}

	return Version{
			Major: majorUint,
			Minor: minorUint,
		},
		nil
}

type PRVersion struct {
	Identifiers []PRIdentifier
}

type BuildMetadata struct {
	Identifiers []BuildIdentifier
}

type PRIdentifier struct {
	identifier string
}

func (i *PRIdentifier) String() string {
	return i.identifier
}

func NewPrIdentifier(v string) (PRIdentifier, error) {
	i := PRIdentifier{}
	if err := i.Set(v); err != nil {
		return PRIdentifier{}, err
	}
	return i, nil
}

func (i *PRIdentifier) Set(v string) error {
	if len(v) < 1 {
		return fmt.Errorf("prerelease identifiers MUST NOT be empty, got: %s", v)
	}

	if containsOnly(v, numbers) {
		if v[0] == '0' {
			return fmt.Errorf("prerelease numeric identifiers MUST NOT include leading zeros, got: %s", v)
		}
	}

	if containsOnly(v, alphanum) {
		return fmt.Errorf("prerelease identifiers MUST comprise only ASCII alphanumerics and hyphens [0-9-Za-z-], got: %s", v)
	}

	i.identifier = v
	return nil
}

func containsOnly(s string, set string) bool {
	return strings.IndexFunc(s, func(r rune) bool {
		return !strings.ContainsRune(set, r)
	}) == -1
}

type BuildIdentifier struct {
	identifier string
}

func NewBuildIdentifier(v string) (BuildIdentifier, error) {
	i := BuildIdentifier{}
	if err := i.Set(v); err != nil {
		return BuildIdentifier{}, err
	}
	return i, nil
}

func (i *BuildIdentifier) Set(v string) error {
	if len(v) < 1 {
		return fmt.Errorf("prerelease identifiers MUST NOT be empty, got: %s", v)
	}

	if containsOnly(v, alphanum) {
		return fmt.Errorf("build identifiers MUST comprise only ASCII alphanumerics and hyphens [0-9-Za-z-], got: %s", v)
	}
	return nil
}
