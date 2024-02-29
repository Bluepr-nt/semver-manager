package models

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestVersionPattern_FirstRelease(t *testing.T) {
	tests := []struct {
		name    string
		pattern VersionPattern
		want    Release
	}{
		{
			name:    "Simple pattern",
			pattern: newVersionPattern("1.0.*-test.*"),
			want:    newRelease("1.0.0-test.0"),
		},
		{
			name:    "Simple pattern",
			pattern: newVersionPattern("1.*.0"),
			want:    newRelease("1.0.0"),
		},
		{
			name:    "Simple pattern",
			pattern: newVersionPattern("1.*.1"),
			want:    newRelease("1.0.1"),
		},
		{
			name:    "Simple pattern",
			pattern: newVersionPattern("*.*.1"),
			want:    newRelease("0.0.1"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.pattern.FirstRelease()
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestVersionPattern_FirstPrerelease(t *testing.T) {

	tests := []struct {
		name    string
		pattern VersionPattern
		want    PRVersion
	}{
		{
			name:    "Simple pattern",
			pattern: newVersionPattern("1.0.0-test.*"),
			want:    newPRVersion("1.0.0-test.0"),
		},
		{
			name:    "Wildcard pattern",
			pattern: newVersionPattern("1.0.0-*"),
			want:    newPRVersion("1.0.0-0"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.pattern.FirstPrerelease()
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestVersionPattern_FirstBuildMetadata(t *testing.T) {
	tests := []struct {
		name    string
		pattern VersionPattern
		want    BuildMetadata
	}{
		{
			name:    "Simple pattern",
			pattern: newVersionPattern("1.0.0+test.*"),
			want:    newBuildMetadata("1.0.0+test.0"),
		},
		{
			name:    "Wildcard pattern",
			pattern: newVersionPattern("1.0.0+*"),
			want:    newBuildMetadata("1.0.0+0"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.pattern.FirstBuildMetadata()
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestVersionPattern_IsReleaseOnlyPattern(t *testing.T) {
	tests := []struct {
		name    string
		pattern VersionPattern
		want    bool
	}{
		{
			name:    "Is a ReleaseOnlyPattern",
			pattern: newVersionPattern("1.*.*"),
			want:    true,
		},
		{
			name:    "Is a ReleaseOnlyPattern",
			pattern: newVersionPattern("1.1.1"),
			want:    true,
		},
		{
			name:    "Is NOT a ReleaseOnlyPattern",
			pattern: newVersionPattern("1.*.*-Alpha"),
			want:    false,
		},
		{
			name:    "Empty ReleaseOnlyPattern",
			pattern: VersionPattern{},
			want:    false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			got := tt.pattern.IsReleaseOnlyPattern()
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestParseVersionPattern(t *testing.T) {
	tests := []struct {
		name               string
		pattern            string
		expectedRelease    ReleasePattern
		expectedPrerelease PRVersionPattern
		expectErr          bool
	}{
		{
			name:    "Valid Version Pattern",
			pattern: "1.2.3-alpha.1",
			expectedRelease: ReleasePattern{
				Major: MajorPattern{Pattern{"1"}},
				Minor: MinorPattern{Pattern{"2"}},
				Patch: PatchPattern{Pattern{"3"}},
			},
			expectedPrerelease: PRVersionPattern{
				Identifiers: []PRIdentifierPattern{
					{pattern: Pattern{"alpha"}},
					{pattern: Pattern{"1"}},
				},
			},
			expectErr: false,
		},
		{
			name:               "Invalid Version Pattern",
			pattern:            "1.2.3-alpha.1.bad+",
			expectedRelease:    ReleasePattern{},
			expectedPrerelease: PRVersionPattern{},
			expectErr:          true,
		},
		{
			name:               "Invalid Release Major Pattern",
			pattern:            "A.2.3",
			expectedRelease:    ReleasePattern{},
			expectedPrerelease: PRVersionPattern{},
			expectErr:          true,
		},
		{
			name:               "Invalid Release Minor Pattern",
			pattern:            "1.A.3",
			expectedRelease:    ReleasePattern{},
			expectedPrerelease: PRVersionPattern{},
			expectErr:          true,
		},
		{
			name:               "Invalid Release Patch Pattern",
			pattern:            "1.2.A",
			expectedRelease:    ReleasePattern{},
			expectedPrerelease: PRVersionPattern{},
			expectErr:          true,
		},
		{
			name:               "Invalid Prerelease Pattern",
			pattern:            "1.2.3-!",
			expectedRelease:    ReleasePattern{},
			expectedPrerelease: PRVersionPattern{},
			expectErr:          true,
		},
		{
			name:               "Invalid Empty Prerelease Pattern",
			pattern:            "1.2.3-",
			expectedRelease:    ReleasePattern{},
			expectedPrerelease: PRVersionPattern{},
			expectErr:          true,
		},
		{
			name:               "Invalid leading zero Prerelease Pattern",
			pattern:            "1.2.3-00",
			expectedRelease:    ReleasePattern{},
			expectedPrerelease: PRVersionPattern{},
			expectErr:          true,
		},
		{
			name:               "Invalid BuildMetadata Pattern",
			pattern:            "1.2.3+!",
			expectedRelease:    ReleasePattern{},
			expectedPrerelease: PRVersionPattern{},
			expectErr:          true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			pattern, err := ParseVersionPattern(tt.pattern)
			if tt.expectErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expectedRelease, pattern.Release)
				assert.Equal(t, tt.expectedPrerelease, pattern.Prerelease)
			}
		})
	}
}

func newVersionPattern(s string) VersionPattern {
	v, _ := ParseVersionPattern(s)
	return v
}

func newBuildMetadata(s string) BuildMetadata {
	buildMetadata, _ := parseBuildMetadata(s)
	return buildMetadata
}

func newPRVersion(s string) PRVersion {
	prerelease, _ := parsePrerelease(s)
	return prerelease
}

func newRelease(s string) Release {
	release, _ := parseRelease(s)
	return release
}
