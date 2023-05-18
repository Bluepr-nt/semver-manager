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

type Release struct {
	Major uint64
	Minor uint64
	Patch uint64
}

func (r *Release) String() string {
	return fmt.Sprintf("%d.%d.%d", r.Major, r.Minor, r.Patch)
}

type Version struct {
	Release    Release
	Prerelease PRVersion
	Build      BuildMetadata
}

func (v *Version) String() string {
	version := fmt.Sprintf("%s%s%s", v.Release.String(), v.Prerelease.String(), v.Build.String())
	return version
}

func NewVersion(v string) (Version, error) {

	release, err := parseRelease(v)
	if err != nil {
		return Version{}, err
	}

	prVersion, err := parsePrerelease(v)
	if err != nil {
		return Version{}, err
	}

	buildMetadata, err := parseBuildMetadata(v)
	if err != nil {
		return Version{}, err
	}

	return Version{
			Release:    release,
			Prerelease: prVersion,
			Build:      buildMetadata,
		},
		nil
}

func parseBuildMetadata(v string) (BuildMetadata, error) {
	tokens := strings.SplitN(v, "+", 2)
	if len(tokens) < 2 {
		return BuildMetadata{}, nil
	}
	metadata := tokens[1]
	identifiers := strings.Split(metadata, ".")
	buildIdentifiers := []BuildIdentifier{}
	for _, identifier := range identifiers {
		buildIdentifier, err := NewBuildIdentifier(identifier)
		if err != nil {
			return BuildMetadata{}, err
		}
		buildIdentifiers = append(buildIdentifiers, buildIdentifier)
	}
	return BuildMetadata{
		Identifiers: buildIdentifiers,
	}, nil
}

func parseRelease(v string) (Release, error) {
	r := strings.SplitN(v, "-", 2)[0]
	r = strings.SplitN(r, "+", 2)[0]

	majorUint, err := parseMajor(r)
	if err != nil {
		return Release{}, err
	}

	minorUint, err := parseMinor(r)
	if err != nil {
		return Release{}, err
	}

	patchUint, err := parsePatch(r)
	if err != nil {
		return Release{}, err
	}
	return Release{
		Major: majorUint,
		Minor: minorUint,
		Patch: patchUint,
	}, nil
}

func parsePrerelease(v string) (PRVersion, error) {
	tokens := strings.SplitN(v, "+", 2)
	tokens = strings.SplitN(tokens[0], "-", 2)
	if len(tokens) < 2 {
		return PRVersion{}, nil
	}
	pr := tokens[1]
	identifiers := strings.Split(pr, ".")
	prIdentifiers := []PRIdentifier{}
	for _, identifier := range identifiers {
		prIdentifier, err := NewPrIdentifier(identifier)
		if err != nil {
			return PRVersion{}, err
		}
		prIdentifiers = append(prIdentifiers, prIdentifier)

	}

	return PRVersion{Identifiers: prIdentifiers}, nil
}

func parsePatch(v string) (uint64, error) {
	tokens := strings.SplitN(v, ".", 3)
	patch := tokens[2]
	if !containsOnly(patch, numbers) {
		return 0, fmt.Errorf("patch version MUST be non-negative integer, got: %s", patch)
	}
	if len(patch) > 1 && patch[0] == '0' {
		return 0, fmt.Errorf("patch version MUST NOT contain leading zeroes, got: %s", patch)
	}
	patchUint, err := strconv.ParseUint(patch, 10, 64)
	if err != nil {
		return 0, err
	}

	return patchUint, nil
}

func parseMinor(v string) (uint64, error) {
	tokens := strings.SplitN(v, ".", 3)
	minor := tokens[1]
	if !containsOnly(minor, numbers) {
		return 0, fmt.Errorf("minor version MUST be non-negative integer, got: %s", minor)
	}
	minorUint, err := strconv.ParseUint(minor, 10, 64)
	if err != nil {
		return 0, err
	}
	if len(minor) > 1 && minor[0] == '0' {
		return 0, fmt.Errorf("minor version MUST NOT contain leading zeroes, got: %s", minor)
	}
	return minorUint, nil
}

func parseMajor(v string) (uint64, error) {
	tokens := strings.SplitN(v, ".", 2)
	major := tokens[0]
	if !containsOnly(major, numbers) {
		return 0, fmt.Errorf("major version MUST be non-negative integer, got: %s", major)
	}

	if len(major) > 1 && major[0] == '0' {
		return 0, fmt.Errorf("major version MUST NOT contain leading zeroes, got: %s", major)
	}

	majorUint, err := strconv.ParseUint(major, 10, 64)
	if err != nil {
		return 0, fmt.Errorf("error converting major version number to uint64: %w", err)
	}
	return majorUint, nil
}

type PRVersion struct {
	Identifiers []PRIdentifier
}

func (pr *PRVersion) String() string {
	identifierList := []string{}
	for _, identifier := range pr.Identifiers {
		identifierList = append(identifierList, identifier.identifier)
	}
	prerelease := strings.Join(identifierList, ".")
	if len(prerelease) > 0 {
		return fmt.Sprintf("-%s", prerelease)
	}
	return ""
}

func NewPRVersion(identifiers []string) (PRVersion, error) {
	prVersion := PRVersion{}
	for _, identifier := range identifiers {
		prIdentifier, err := NewPrIdentifier(identifier)
		if err != nil {
			return PRVersion{}, err
		}
		prVersion.Identifiers = append(prVersion.Identifiers, prIdentifier)
	}
	return prVersion, nil
}

type BuildMetadata struct {
	Identifiers []BuildIdentifier
}

func (bm *BuildMetadata) String() string {
	identifierList := []string{}
	for _, identifier := range bm.Identifiers {
		identifierList = append(identifierList, identifier.identifier)
	}
	buildMetadata := strings.Join(identifierList, ".")
	if len(buildMetadata) > 0 {
		return fmt.Sprintf("+%s", buildMetadata)
	}
	return ""
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
		if len(v) > 1 && v[0] == '0' {
			return fmt.Errorf("prerelease numeric identifiers MUST NOT include leading zeros, got: %s", v)
		}
	}

	if !containsOnly(v, alphanum) {
		return fmt.Errorf("prerelease identifiers MUST comprise only ASCII alphanumerics and hyphens [0-9-Za-z-], got: %s", v)
	}

	i.identifier = v
	return nil
}

// Function containsOnly checks if all characters in the input string 's' are present in the set of valid characters 'set'
func containsOnly(s string, set string) bool {

	// A helper function that checks if a single character is in the set of valid characters
	characterIsInSet := func(r rune) bool {
		return strings.ContainsRune(set, r)
	}

	// The strings.IndexFunc function will return the index of the first character in 's' that does not satisfy the helper function.
	// If all characters satisfy the helper function, it returns -1.
	// Therefore, if the returned index is -1, all characters in 's' are in the 'set' of valid characters.
	return strings.IndexFunc(s, func(r rune) bool {
		return !characterIsInSet(r)
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

	if !containsOnly(v, alphanum) {
		return fmt.Errorf("build identifiers MUST comprise only ASCII alphanumerics and hyphens [0-9-Za-z-], got: %s", v)
	}
	i.identifier = v
	return nil
}
