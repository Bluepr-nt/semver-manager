package models

import (
	"fmt"
	"strings"
)

const Wildcard = "*"

type VersionPattern struct {
	Release    ReleasePattern
	Prerelease PRVersionPattern
	Build      BuildMetadataPattern
}

func (v VersionPattern) IsReleaseOnlyPattern() bool {
	if v.IsEmpty() || len(v.Prerelease.Identifiers) > 0 {
		return false
	}
	return true
}

func (v VersionPattern) IsEmpty() bool {
	if v.Release.Major.Value() == "" &&
		v.Release.Minor.Value() == "" &&
		v.Release.Patch.Value() == "" &&
		len(v.Prerelease.Identifiers) < 1 &&
		len(v.Build.Identifiers) < 1 {
		return true
	}
	return false
}
func (v VersionPattern) FirstVersion() (firstVersion Version) {
	firstVersion.Release = v.FirstRelease()
	firstVersion.Prerelease = v.FirstPrerelease()
	firstVersion.BuildMetadata = v.FirstBuildMetadata()

	return firstVersion
}

func (v VersionPattern) FirstRelease() (FirstRelease Release) {
	major := getAbsoluteValue(v.Release.Major.pattern)
	minor := getAbsoluteValue(v.Release.Minor.pattern)
	patch := getAbsoluteValue(v.Release.Patch.pattern)

	release, _ := parseRelease(fmt.Sprintf("%s.%s.%s", major, minor, patch))
	return release
}

func (v VersionPattern) FirstPrerelease() PRVersion {
	var rawPRVersion string
	for i, identifier := range v.Prerelease.Identifiers {
		rawId := getAbsoluteValue(identifier.pattern)
		if i == 0 {
			rawPRVersion = fmt.Sprintf("%s-%s", rawPRVersion, rawId)
		} else {
			rawPRVersion = fmt.Sprintf("%s.%s", rawPRVersion, rawId)
		}
	}
	prVersion, _ := parsePrerelease(rawPRVersion)
	return prVersion
}

func (v VersionPattern) FirstBuildMetadata() BuildMetadata {
	var rawBuildMetadata string
	for i, buildId := range v.Build.Identifiers {
		rawId := getAbsoluteValue(buildId.pattern)
		if i == 0 {
			rawBuildMetadata = fmt.Sprintf("%s+%s", rawBuildMetadata, rawId)
		} else {
			rawBuildMetadata = fmt.Sprintf("%s.%s", rawBuildMetadata, rawId)
		}
	}
	buildMetadata, _ := parseBuildMetadata(rawBuildMetadata)
	return buildMetadata
}

func getAbsoluteValue(patten Pattern) string {
	if patten.value == Wildcard {
		return "0"
	}
	return patten.value
}

// func getAbsoluteIdentifiers
func ParseVersionPattern(pattern string) (VersionPattern, error) {
	release, err := parseReleasePattern(pattern)
	if err != nil {
		return VersionPattern{}, err
	}
	prerelease, err := parsePrereleasePattern(pattern)
	if err != nil {
		return VersionPattern{}, err
	}

	buildMetadata, err := parseBuildMetadataPattern(pattern)
	if err != nil {
		return VersionPattern{}, err
	}
	return VersionPattern{
		Release:    release,
		Prerelease: prerelease,
		Build:      buildMetadata,
	}, nil
}

func parsePrereleasePattern(pattern string) (PRVersionPattern, error) {
	tokens := strings.SplitN(pattern, "+", 2)
	tokens = strings.SplitN(tokens[0], "-", 2)
	if len(tokens) < 2 {
		return PRVersionPattern{}, nil
	}
	prereleasePattern := tokens[1]
	identifiersPattern := strings.Split(prereleasePattern, ".")
	var prIdentifiersPattern []PRIdentifierPattern
	for _, identifierPattern := range identifiersPattern {
		p, err := parsePrIdentifierPattern(identifierPattern)
		if err != nil {
			return PRVersionPattern{}, err
		}
		prIdentifiersPattern = append(prIdentifiersPattern, p)
	}
	return PRVersionPattern{Identifiers: prIdentifiersPattern}, nil
}

func parsePrIdentifierPattern(pattern string) (PRIdentifierPattern, error) {
	identifierPattern := PRIdentifierPattern{}
	if err := identifierPattern.Set(pattern); err != nil {
		return PRIdentifierPattern{}, err
	}
	return identifierPattern, nil
}

func (i *PRIdentifierPattern) Set(pattern string) error {
	if len(pattern) < 1 {
		return fmt.Errorf("prerelease identifiers MUST NOT be empty, got: %s", pattern)
	}
	if containsOnly(pattern, numbers) {
		if len(pattern) > 1 && pattern[0] == '0' {
			return fmt.Errorf("prerelease numeric identifiers MUST NOT include leading zeros, got: %s", pattern)
		}
		i.pattern = Pattern{value: pattern}
		return nil
	}
	if containsOnly(pattern, alphanum) {
		i.pattern = Pattern{value: pattern}
		return nil
	}

	if pattern == Wildcard {
		i.pattern = Pattern{value: pattern}
		return nil
	}
	return fmt.Errorf("prerelease identifiers MUST contain only alphanumerics and hyphens, got: %s", pattern)
}

func parseReleasePattern(pattern string) (ReleasePattern, error) {
	release := strings.SplitN(pattern, "-", 2)[0]
	release = strings.SplitN(release, "+", 2)[0]
	major, err := parseMajorPattern(release)
	if err != nil {
		return ReleasePattern{}, err
	}

	minor, err := parseMinorPattern(release)
	if err != nil {
		return ReleasePattern{}, err
	}

	patch, err := parsePatchPattern(release)
	if err != nil {
		return ReleasePattern{}, err
	}

	return ReleasePattern{
		Major: major,
		Minor: minor,
		Patch: patch,
	}, nil
}

func parseMajorPattern(pattern string) (ReleaseDigitPattern, error) {
	tokens := strings.SplitN(pattern, ".", 2)
	majorPattern := tokens[0]
	p, err := parseDigitsPattern(majorPattern, Major)
	if err != nil {
		return ReleaseDigitPattern{}, err
	}
	return ReleaseDigitPattern{pattern: p}, nil
}

func parseMinorPattern(pattern string) (ReleaseDigitPattern, error) {
	tokens := strings.SplitN(pattern, ".", 3)
	minor := tokens[1]
	p, err := parseDigitsPattern(minor, Minor)
	if err != nil {
		return ReleaseDigitPattern{}, err
	}
	return ReleaseDigitPattern{pattern: p}, nil
}

func parsePatchPattern(pattern string) (ReleaseDigitPattern, error) {
	tokens := strings.SplitN(pattern, ".", 3)
	patch := tokens[2]
	p, err := parseDigitsPattern(patch, Patch)
	if err != nil {
		return ReleaseDigitPattern{}, err
	}
	return ReleaseDigitPattern{pattern: p}, nil
}

func parseDigitsPattern(pattern string, increment Increment) (Pattern, error) {
	if err := versionDigitsCompliance(pattern, increment); err != nil {
		if pattern == Wildcard {
			return Pattern{value: pattern}, nil
		} else {
			return Pattern{}, err
		}
	}
	return Pattern{value: pattern}, nil
}

func parseBuildMetadataPattern(pattern string) (BuildMetadataPattern, error) {
	tokens := strings.SplitN(pattern, "+", 2)
	if len(tokens) < 2 {
		return BuildMetadataPattern{}, nil
	}
	buildMetadataPattern := tokens[1]
	identifiersPattern := strings.Split(buildMetadataPattern, ".")
	var buildIdentifiersPattern []BuildIdentifierPattern
	for _, identifierPattern := range identifiersPattern {
		p, err := parseBuildIdentifierPattern(identifierPattern)
		if err != nil {
			return BuildMetadataPattern{}, err
		}
		buildIdentifiersPattern = append(buildIdentifiersPattern, p)
	}
	return BuildMetadataPattern{Identifiers: buildIdentifiersPattern}, nil
}

func parseBuildIdentifierPattern(pattern string) (BuildIdentifierPattern, error) {
	identifierPattern := BuildIdentifierPattern{}
	if err := identifierPattern.Set(pattern); err != nil {
		return BuildIdentifierPattern{}, err
	}
	return identifierPattern, nil
}

func (i *BuildIdentifierPattern) Set(pattern string) error {
	if len(pattern) < 1 {
		return fmt.Errorf("build identifiers MUST NOT be empty, got: %s", pattern)
	}
	if containsOnly(pattern, alphanum) {
		i.pattern = Pattern{value: pattern}
		return nil
	}
	if pattern == Wildcard {
		i.pattern = Pattern{value: pattern}
		return nil
	}
	return fmt.Errorf("build identifiers MUST contain only alphanumerics and hyphens, got: %s", pattern)
}

type ReleasePattern struct {
	Major ReleaseDigitPattern
	Minor ReleaseDigitPattern
	Patch ReleaseDigitPattern
}

type ReleaseDigitPattern struct {
	pattern Pattern
}

func (m ReleaseDigitPattern) Value() string {
	return m.pattern.value
}

type PRVersionPattern struct {
	Identifiers []PRIdentifierPattern
}

type PRIdentifierPattern struct {
	pattern Pattern
}

func (p PRIdentifierPattern) Value() string {
	return p.pattern.value
}

type BuildMetadataPattern struct {
	Identifiers []BuildIdentifierPattern
}

type BuildIdentifierPattern struct {
	pattern Pattern
}

func (p BuildIdentifierPattern) Value() string {
	return p.pattern.value
}

type Pattern struct {
	value string
}
