package models

import (
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestVersionPattern_FirstRelease(t *testing.T) {
	type fields struct {
		Release    ReleasePattern
		Prerelease PRVersionPattern
		Build      BuildMetadataPattern
	}
	tests := []struct {
		name             string
		fields           fields
		wantFirstRelease Release
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			v := VersionPattern{
				Release:    tt.fields.Release,
				Prerelease: tt.fields.Prerelease,
				Build:      tt.fields.Build,
			}
			if gotFirstRelease := v.FirstRelease(); !reflect.DeepEqual(gotFirstRelease, tt.wantFirstRelease) {
				t.Errorf("VersionPattern.FirstRelease() = %v, want %v", gotFirstRelease, tt.wantFirstRelease)
			}
		})
	}
}

func TestVersionPattern_FirstPrerelease(t *testing.T) {
	type fields struct {
		Release    ReleasePattern
		Prerelease PRVersionPattern
		Build      BuildMetadataPattern
	}
	tests := []struct {
		name   string
		fields fields
		want   PRVersion
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			v := VersionPattern{
				Release:    tt.fields.Release,
				Prerelease: tt.fields.Prerelease,
				Build:      tt.fields.Build,
			}
			if got := v.FirstPrerelease(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("VersionPattern.FirstPrerelease() = %v, want %v", got, tt.want)
			}
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
