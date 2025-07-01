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

type Release struct {
	Major ReleaseDigit
	Minor ReleaseDigit
	Patch ReleaseDigit
}

// IsEqualTo compares two Releases and returns true if they are equal
func (r Release) IsEqualTo(releaseB Release) bool {
	return r.Major == releaseB.Major && r.Minor == releaseB.Minor && r.Patch == releaseB.Patch
}

func (r Release) IsHigherThan(rB Release) bool {
	if r.Major.GT(rB.Major) {
		return true
	}

	if r.Major == rB.Major && r.Minor.GT(rB.Minor) {
		return true
	}

	if r.Major == rB.Major && r.Minor == rB.Minor && r.Patch.GT(rB.Patch) {
		return true
	}

	return false
}

func (r *Release) String() string {
	return fmt.Sprintf("%d.%d.%d", r.Major.value, r.Minor.value, r.Patch.value)
}

type ReleaseDigit struct {
	value uint64
}

func (r *ReleaseDigit) Increment() {
	r.value = r.value + 1
}

func (r ReleaseDigit) GT(cmp ReleaseDigit) bool {
	return r.value > cmp.value
}

func (r ReleaseDigit) LT(cmp ReleaseDigit) bool {
	return r.value < cmp.value
}

func (r *ReleaseDigit) Set(d uint64) {
	r.value = d
}

func (r ReleaseDigit) String() string {
	stringValue := fmt.Sprintf("%d", r.value)
	return stringValue
}

func (r ReleaseDigit) Value() uint64 {
	return r.value
}

type Version struct {
	Release       Release
	Prerelease    PRVersion
	BuildMetadata BuildMetadata
}

func (v Version) IsRelease() bool {
	return len(v.Prerelease.Identifiers) < 1
}

// IsEqualTo compares two Versions and returns true if they are equal
// Build metadata is ignored in the comparison as per the Semver specification
func (v *Version) IsEqualTo(versionB Version) bool {
	if v.Release.IsEqualTo(versionB.Release) && v.Prerelease.IsEqualTo(versionB.Prerelease) {
		return true
	}
	return false
}

// IsHigherThan compares two Versions and returns true if the first version is higher than the second
// according to the Semver specification
func (v Version) IsHigherThan(versionB Version) bool {

	if v.Release.IsHigherThan(versionB.Release) {
		return true
	}

	if v.IsRelease() && !versionB.IsRelease() {
		return true
	}

	if versionB.IsRelease() && !v.IsRelease() {
		return false
	}

	return v.Prerelease.IsHigherThan(versionB.Prerelease)
}

func (v *Version) String() string {
	version := fmt.Sprintf("%s%s%s", v.Release.String(), v.Prerelease.String(), v.BuildMetadata.String())
	return version
}

func ParseVersion(v string) (Version, error) {

	rRelease, rPrerelease, rBuildMetadata := GetVersionComponents(v)
	release, err := ParseRelease(rRelease)
	if err != nil {
		return Version{}, err
	}

	var prVersion PRVersion
	if rPrerelease == "" {
		prVersion = PRVersion{}
	} else {
		prVersion, err = ParsePRVersion(rPrerelease)
		if err != nil {
			return Version{}, err
		}
	}

	var buildMetadata BuildMetadata
	if rBuildMetadata == "" {
		buildMetadata = BuildMetadata{}

	} else {
		buildMetadata, err = ParseBuildMetadata(rBuildMetadata)
		if err != nil {
			return Version{}, err
		}
	}

	return Version{
			Release:       release,
			Prerelease:    prVersion,
			BuildMetadata: buildMetadata,
		},
		nil
}

func GetVersionComponents(v string) (release, prerelease, buildMetadata string) {
	release = GetRelease(v)
	prerelease = GetPrerelease(v)
	buildMetadata = GetBuildMetadata(v)
	return release, prerelease, buildMetadata

}

func GetRelease(v string) string {
	tokens := strings.SplitN(v, "+", 2)
	tokens = strings.SplitN(tokens[0], "-", 2)

	return tokens[0]
}

func GetPrerelease(v string) string {
	tokens := strings.SplitN(v, "+", 2)
	tokens = strings.SplitN(tokens[0], "-", 2)
	if len(tokens) < 2 {
		return ""
	}

	return tokens[1]
}

func GetBuildMetadata(v string) string {
	tokens := strings.SplitN(v, "+", 2)
	if len(tokens) < 2 {
		return ""
	}
	return tokens[1]
}

type VersionSlice []Version

func (vs VersionSlice) String() string {
	var builder strings.Builder

	for i, version := range vs {
		builder.WriteString(version.String())
		if i < len(vs)-1 {
			builder.WriteString(" ")
		}
	}

	return builder.String()
}

func ParseVersions(vList string) (VersionSlice, error) {

	versionStrings := SplitVersions(vList)
	var versions VersionSlice

	for _, rawVersion := range versionStrings {
		version, err := ParseVersion(rawVersion)
		if err != nil {
			return nil, err
		}

		versions = append(versions, version)

	}
	return versions, nil
}

func ParseBuildMetadata(metadata string) (BuildMetadata, error) {

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

func ParseRelease(v string) (Release, error) {
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
		Major: ReleaseDigit{value: majorUint},
		Minor: ReleaseDigit{minorUint},
		Patch: ReleaseDigit{patchUint},
	}, nil
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

func (pr PRVersion) IsEqualTo(prB PRVersion) bool {
	if len(pr.Identifiers) != len(prB.Identifiers) {
		return false
	}

	for index, identifier := range pr.Identifiers {
		if !identifier.IsEqualTo(prB.Identifiers[index]) {
			return false
		}
	}
	return true
}

func (pr *PRVersion) IsHigherThan(prB PRVersion) bool {
	for index, identifier := range pr.Identifiers {
		if len(prB.Identifiers) <= index {
			return true
		} else if identifier.IsEqualTo(prB.Identifiers[index]) {
			continue
		} else {
			return identifier.IsHigherThan(prB.Identifiers[index])
		}
	}
	return false
}

func (pr PRVersion) LastID() PRIdentifier {
	if len(pr.Identifiers) > 0 {
		return pr.Identifiers[len(pr.Identifiers)-1]
	}
	return PRIdentifier{}
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

func ParsePRVersion(rawPrerelease string) (PRVersion, error) {
	identifiers := strings.Split(rawPrerelease, ".")
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

func StringToVersions(stringVersions string) ([]Version, error) {
	versionStringSlice := SplitVersions(stringVersions)

	versions := make([]Version, len(versionStringSlice))
	var err error
	for i, versionStr := range versionStringSlice {
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

func (i PRIdentifier) Value() string {
	return i.identifier
}

// IsEqualTo compares two PRIdentifiers and returns true if they are equal
func (i PRIdentifier) IsEqualTo(identifierB PRIdentifier) bool {
	return i.identifier == identifierB.identifier
}

// IsHigherThan compares two PRIdentifiers and returns true if the first identifier is higher than the second
// according to the Semver specification
func (i PRIdentifier) IsHigherThan(identifierB PRIdentifier) bool {
	if containsOnly(i.identifier, numbers) && containsOnly(identifierB.identifier, numbers) {
		iValue, _ := strconv.ParseUint(i.identifier, 10, 64)
		identifierBValue, _ := strconv.ParseUint(identifierB.identifier, 10, 64)
		return iValue > identifierBValue
	}
	return i.identifier > identifierB.identifier
}

// ParsePrIdentifier parses a string into a PRIdentifier
func ParsePrIdentifier(v string) (PRIdentifier, error) {
	i := PRIdentifier{}
	if err := i.Set(v); err != nil {
		return PRIdentifier{}, err
	}
	return i, nil
}

// Set sets the PRIdentifier value
// It returns an error if the value is empty, contains leading zeros,
// or contains characters other than alphanumerics and hyphens
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

func containsOnly(s string, set string) bool {

	characterIsInSet := func(r rune) bool {
		return strings.ContainsRune(set, r)
	}

	return strings.IndexFunc(s, func(r rune) bool {
		return !characterIsInSet(r)
	}) == -1
}

type BuildIdentifier struct {
	identifier string
}

// String returns the BuildIdentifier value
func (i *BuildIdentifier) String() string {
	return i.identifier
}

// ParseBuildIdentifier parses a string into a BuildIdentifier
func ParseBuildIdentifier(v string) (BuildIdentifier, error) {
	i := BuildIdentifier{}
	if err := i.Set(v); err != nil {
		return BuildIdentifier{}, err
	}
	return i, nil
}

// Set sets the BuildIdentifier value
// It returns an error if the value is empty or contains characters other than alphanumerics and hyphens
func (i *BuildIdentifier) Set(v string) error {
	if len(v) < 1 {
		return fmt.Errorf("build identifiers MUST NOT be empty, got: %s", v)
	}

	if !containsOnly(v, alphanum) {
		return fmt.Errorf("build identifiers MUST comprise only ASCII alphanumerics and hyphens [0-9-Za-z-], got: %s", v)
	}
	i.identifier = v
	return nil
}

func versionDigitsCompliance(version string, increment Increment) error {
	if increment != Major && increment != Minor && increment != Patch {
		return fmt.Errorf("increment MUST be one of %s, %s, or %s, got: %s", Major, Minor, Patch, increment)
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

// SplitVersions splits a string of versions into a slice of strings
// separated by commas or spaces
func SplitVersions(stringVersions string) (versionStringSlice []string) {
	stringVersions = strings.TrimSpace(stringVersions)

	if strings.Contains(stringVersions, ",") {
		stringVersions = strings.ReplaceAll(stringVersions, " ", "")
		versionStringSlice = strings.Split(stringVersions, ",")
	} else {
		versionStringSlice = strings.Split(stringVersions, " ")
	}
	return versionStringSlice
}
