package models

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewBuildIdentifier(t *testing.T) {
	type args struct {
		v string
	}
	tests := []struct {
		name    string
		args    args
		want    BuildIdentifier
		wantErr bool
	}{
		{
			name: "Valid identifier",
			args: args{
				v: "valid",
			},
			want:    BuildIdentifier{identifier: "valid"},
			wantErr: false,
		},
		{
			name: "Invalid identifier",
			args: args{
				v: "@v",
			},
			want:    BuildIdentifier{identifier: ""},
			wantErr: true,
		},
		{
			name: "Empty identifier",
			args: args{
				v: "",
			},
			want:    BuildIdentifier{},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ParseBuildIdentifier(tt.args.v)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestNewPrIdentifier(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		want    PRIdentifier
		wantErr bool
	}{
		{
			name:    "Valid PRIdentifier",
			input:   "beta",
			want:    PRIdentifier{identifier: "beta"},
			wantErr: false,
		},
		{
			name:    "Empty PRIdentifier",
			input:   "",
			want:    PRIdentifier{},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ParsePrIdentifier(tt.input)

			assert.Equal(t, tt.wantErr, err != nil, "NewPrIdentifier() error = %v, wantErr %v", err, tt.wantErr)
			assert.Equal(t, tt.want, got, "NewPrIdentifier() = %v, want %v", got, tt.want)
		})
	}
}

func TestNewVersion(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		want    Version
		wantErr bool
	}{
		{
			name:  "Valid Version",
			input: "1.0.0-beta",
			want: Version{
				Release: Release{
					Major: 1,
					Minor: 0,
					Patch: 0,
				},
				Prerelease:    PRVersion{Identifiers: []PRIdentifier{{identifier: "beta"}}},
				BuildMetadata: BuildMetadata{},
			},
			wantErr: false,
		},
		{
			name:  "Valid Version with build metadata and prerelease",
			input: "1.0.0-beta+0000044ttt",
			want: Version{
				Release: Release{
					Major: 1,
					Minor: 0,
					Patch: 0,
				},
				Prerelease:    PRVersion{Identifiers: []PRIdentifier{{identifier: "beta"}}},
				BuildMetadata: BuildMetadata{Identifiers: []BuildIdentifier{{identifier: "0000044ttt"}}},
			},
			wantErr: false,
		},
		{
			name:  "Valid Version with build metadata",
			input: "1.0.0+0000044ttt",
			want: Version{
				Release: Release{
					Major: 1,
					Minor: 0,
					Patch: 0,
				},
				Prerelease:    PRVersion{},
				BuildMetadata: BuildMetadata{Identifiers: []BuildIdentifier{{identifier: "0000044ttt"}}},
			},
			wantErr: false,
		},
		{
			name:  "Valid Version with release only",
			input: "1.0.0",
			want: Version{
				Release: Release{
					Major: 1,
					Minor: 0,
					Patch: 0,
				},
				Prerelease:    PRVersion{},
				BuildMetadata: BuildMetadata{},
			},
			wantErr: false,
		},
		{
			name:    "Invalid patch version",
			input:   "1.0.a",
			want:    Version{},
			wantErr: true,
		},
		{
			name:    "Invalid prerelease version",
			input:   "1.0.0-@",
			want:    Version{},
			wantErr: true,
		},
		{
			name:    "Invalid build metadata",
			input:   "1.0.0+@",
			want:    Version{},
			wantErr: true,
		},
		{
			name:    "Invalid minor version",
			input:   "1.a.0",
			want:    Version{},
			wantErr: true,
		},
		{
			name:    "Invalid major version too big",
			input:   "18446744073709551616.0.0",
			want:    Version{},
			wantErr: true,
		},
		{
			name:    "Invalid minor version too big",
			input:   "1.18446744073709551616.0",
			want:    Version{},
			wantErr: true,
		},
		{
			name:    "Invalid patch version too big",
			input:   "1.0.18446744073709551616",
			want:    Version{},
			wantErr: true,
		},
		{
			name:    "Invalid major version has leading zeroes",
			input:   "01.0.0",
			want:    Version{},
			wantErr: true,
		},
		{
			name:    "Invalid minor version has leading zeroes",
			input:   "1.01.0",
			want:    Version{},
			wantErr: true,
		},
		{
			name:    "Invalid patch version has leading zeroes",
			input:   "1.0.01",
			want:    Version{},
			wantErr: true,
		},
		{
			name:    "Invalid prerelease version has leading zeroes",
			input:   "1.0.0-01",
			want:    Version{},
			wantErr: true,
		},
		{
			name:  "Valid Weird Version",
			input: "1.0.0-beta-.1.-",
			want: Version{
				Release: Release{
					Major: 1,
					Minor: 0,
					Patch: 0,
				},
				Prerelease: PRVersion{Identifiers: []PRIdentifier{
					{identifier: "beta-"},
					{identifier: "1"},
					{identifier: "-"},
				}},
				BuildMetadata: BuildMetadata{},
			},
			wantErr: false,
		},
		{
			name:    "Invalid Version",
			input:   "not a version",
			want:    Version{},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ParseVersion(tt.input)

			assert.Equal(t, tt.wantErr, err != nil, "NewVersion() error = %v, wantErr %v", err, tt.wantErr)
			assert.Equal(t, tt.want, got, "NewVersion() = %v, want %v", got, tt.want)
		})
	}
}

func TestPRIdentifier_String(t *testing.T) {
	tests := []struct {
		name string
		i    PRIdentifier
		want string
	}{
		{
			name: "Valid PRIdentifier",
			i:    PRIdentifier{identifier: "beta"},
			want: "beta",
		},
		{
			name: "Empty PRIdentifier",
			i:    PRIdentifier{},
			want: "",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.i.String()
			assert.Equal(t, tt.want, got, "PRIdentifier.String() = %v, want %v", got, tt.want)
		})
	}
}

func TestVersion_String(t *testing.T) {
	tests := []struct {
		name string
		v    Version
		want string
	}{
		{
			name: "Full version",
			v: Version{
				Release:       Release{Major: 1, Minor: 2, Patch: 3},
				Prerelease:    PRVersion{Identifiers: []PRIdentifier{{identifier: "alpha"}}},
				BuildMetadata: BuildMetadata{Identifiers: []BuildIdentifier{{identifier: "001"}}},
			},
			want: "1.2.3-alpha+001",
		},
		{
			name: "Release only",
			v: Version{
				Release: Release{Major: 1, Minor: 0, Patch: 0},
			},
			want: "1.0.0",
		},
		{
			name: "Release and prerelease",
			v: Version{
				Release:    Release{Major: 2, Minor: 0, Patch: 0},
				Prerelease: PRVersion{Identifiers: []PRIdentifier{{identifier: "beta"}}},
			},
			want: "2.0.0-beta",
		},
		{
			name: "Release and build metadata",
			v: Version{
				Release:       Release{Major: 3, Minor: 1, Patch: 4},
				BuildMetadata: BuildMetadata{Identifiers: []BuildIdentifier{{identifier: "123456"}}},
			},
			want: "3.1.4+123456",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.v.String()
			assert.Equal(t, tt.want, got, "Version.String() = %v, want %v", got, tt.want)
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
		// add more test cases here
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

func TestVersion_IsHigherThan(t *testing.T) {
	type args struct {
		versionB string
	}
	tests := []struct {
		name string
		v    string
		args args
		want bool
	}{
		{
			name: "Major version higher",
			v:    "2.2.2",
			args: args{
				versionB: "1.2.2",
			},
			want: true,
		},
		{
			name: "Minor version higher",
			v:    "1.2.2",
			args: args{
				versionB: "1.1.2",
			},
			want: true,
		},
		{
			name: "Patch version higher",
			v:    "1.0.2",
			args: args{
				versionB: "1.0.1",
			},
			want: true,
		},
		{
			name: "Same version not higher",
			v:    "1.0.0",
			args: args{
				versionB: "1.0.0",
			},
			want: false,
		},
		{
			name: "Prerelease version higher",
			v:    "1.0.0-beta.2",
			args: args{
				versionB: "1.0.0-beta.1",
			},
			want: true,
		},
		{
			name: "Prerelease version lower",
			v:    "1.0.0-alpha.1",
			args: args{
				versionB: "1.0.0-beta.1",
			},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			v, _ := ParseVersion(tt.v)
			versionB, _ := ParseVersion(tt.args.versionB)
			assert.Equal(t, tt.want, v.IsHigherThan(versionB), "Version.IsHigherThan() result not as expected")
		})
	}
}
