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
	Major           = "major"
	Minor           = "minor"
	Patch           = "patch"
)

type Release struct {
	Major uint64
	Minor uint64
	Patch uint64
}

func (r *Release) String() string {
	return fmt.Sprintf("%d.%d.%d", r.Major, r.Minor, r.Patch)
}

type Version struct {
	Release       Release
	Prerelease    PRVersion
	BuildMetadata BuildMetadata
}

func (v *Version) String() string {
	version := fmt.Sprintf("%s%s%s", v.Release.String(), v.Prerelease.String(), v.BuildMetadata.String())
	return version
}

func ParseVersion(v string) (Version, error) {

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
			Release:       release,
			Prerelease:    prVersion,
			BuildMetadata: buildMetadata,
		},
		nil
}

type VersionSlice []Version

func (vs VersionSlice) String() string {
	var builder strings.Builder

	for i, version := range vs {
		builder.WriteString(version.String())
		if i < len(vs)-1 {
			builder.WriteString(", ")
		}
	}

	return builder.String()
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
		buildIdentifier, err := ParseBuildIdentifier(identifier)
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
	release := strings.SplitN(v, "-", 2)[0]
	release = strings.SplitN(release, "+", 2)[0]

	majorUint, err := parseMajor(release)
	if err != nil {
		return Release{}, err
	}

	minorUint, err := parseMinor(release)
	if err != nil {
		return Release{}, err
	}

	patchUint, err := parsePatch(release)
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
		prIdentifier, err := ParsePrIdentifier(identifier)
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
	if err := versionDigitsCompliance(patch, Patch); err != nil {
		return 0, err
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

	if err := versionDigitsCompliance(minor, Minor); err != nil {
		return 0, err
	}

	minorUint, err := strconv.ParseUint(minor, 10, 64)
	if err != nil {
		return 0, err
	}
	return minorUint, nil
}

func parseMajor(v string) (uint64, error) {
	tokens := strings.SplitN(v, ".", 2)
	major := tokens[0]
	if err := versionDigitsCompliance(major, Major); err != nil {
		return 0, err
	}

	majorUint, err := strconv.ParseUint(major, 10, 64)
	if err != nil {
		return 0, err
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

func ParsePRVersion(identifiers []string) (PRVersion, error) {
	prVersion := PRVersion{}
	for _, identifier := range identifiers {
		prIdentifier, err := ParsePrIdentifier(identifier)
		if err != nil {
			return PRVersion{}, err
		}
		prVersion.Identifiers = append(prVersion.Identifiers, prIdentifier)
	}
	return prVersion, nil
}

func StringToVersionList(stringVersions string) ([]Version, error) {
	versionList := []string{}

	stringVersions = strings.TrimSpace(stringVersions)

	if strings.Contains(stringVersions, ",") {
		stringVersions = strings.ReplaceAll(stringVersions, " ", "")
		versionList = strings.Split(stringVersions, ",")
	} else {
		versionList = strings.Split(stringVersions, " ")
	}

	versions := make([]Version, len(versionList))
	var err error
	for i, versionStr := range versionList {
		versions[i], err = ParseVersion(versionStr)
		if err != nil {
			return nil, err
		}
	}
	return versions, nil
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

func ParsePrIdentifier(v string) (PRIdentifier, error) {
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
		} else {
			i.identifier = v
			return nil
		}
	}

	if containsOnly(v, alphanum) {
		i.identifier = v
		return nil
	}

	return fmt.Errorf("prerelease identifiers MUST comprise only ASCII alphanumerics and hyphens [0-9-Za-z-], got: %s", v)
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

func (i *BuildIdentifier) String() string {
	return i.identifier
}
func ParseBuildIdentifier(v string) (BuildIdentifier, error) {
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

func versionDigitsCompliance(version, increment string) error {
	if increment != "major" && increment != "minor" && increment != "patch" {
		return fmt.Errorf("increment MUST be one of 'major', 'minor', or 'patch', got: %s", increment)
	}
	if len(version) < 1 {
		return fmt.Errorf("%s MUST NOT be empty, got: %s", increment, version)
	}

	if !containsOnly(version, numbers) {
		return fmt.Errorf("%s MUST comprise only ASCII numerics [0-9], got: %s", increment, version)
	}

	if len(version) > 1 && version[0] == '0' {
		return fmt.Errorf("%s MUST NOT contain leading zeroes, got: %s", increment, version)
	}

	return nil
}
