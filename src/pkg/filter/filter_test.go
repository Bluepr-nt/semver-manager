package filter

import (
	"errors"
	"src/cmd/smgr/models"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestHighest(t *testing.T) {
	tests := []struct {
		name     string
		versions []models.Version
		want     string
		wantErr  bool
	}{
		{
			name: "Simple tags list",
			versions: []models.Version{
				newVersion("0.0.0-1+054"),
				newVersion("0.1.0-1+054"),
			},
			want:    "0.1.0-1+054",
			wantErr: false,
		},
		{
			name: "Mixed release and prerelease versions",
			versions: []models.Version{
				newVersion("1.0.0"),
				newVersion("0.0.0-2"),
			},
			want:    "1.0.0",
			wantErr: false,
		},
		{
			name: "Prerelease identifiers with different lengths",
			versions: []models.Version{
				newVersion("0.1.0-alpha.2"),
				newVersion("0.1.0-alpha"),
			},
			want:    "0.1.0-alpha.2",
			wantErr: false,
		},
		{
			name: "Versions with different build metadata",
			versions: []models.Version{
				newVersion("0.1.0+100"),
				newVersion("0.1.0+200"),
			},
			want:    "0.1.0+100",
			wantErr: false,
		},
		{
			name:     "Empty version list",
			versions: []models.Version{},
			want:     "",
			wantErr:  true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Highest()(tt.versions)
			if tt.wantErr {
				assert.Error(t, err)
				assert.Equal(t, []models.Version{}, got)

			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.want, got[0].String())
			}
		})
	}
}

func newPrIdentifier(v string) models.PRIdentifier {
	i, _ := models.ParsePrIdentifier(v)
	return i
}

func newBuildIdentifier(v string) models.BuildIdentifier {
	i, _ := models.ParseBuildIdentifier(v)
	return i
}

func newPRVersion(identifiers []string) (p models.PRVersion) {
	p, _ = models.ParsePRVersion(identifiers)
	return
}

func TestApplyFilters(t *testing.T) {
	type args struct {
		versions models.VersionSlice
		filters  []FilterFunc
	}
	tests := []struct {
		name    string
		args    args
		want    models.VersionSlice
		wantErr bool
	}{
		{
			name: "Apply multiple filters",
			args: args{
				versions: models.VersionSlice{
					newVersion("1.2.3"),
					newVersion("1.2.4"),
					newVersion("2.3.4"),
				},
				filters: []FilterFunc{
					VersionPatternFilter(newVersionPattern("1.2.*")),
					Highest(),
				},
			},
			want: models.VersionSlice{
				newVersion("1.2.4"),
			},
			wantErr: false,
		},
		{
			name: "No matching version after filters",
			args: args{
				versions: models.VersionSlice{
					newVersion("1.2.3"),
					newVersion("1.3.4"),
					newVersion("2.3.4"),
				},
				filters: []FilterFunc{
					VersionPatternFilter(newVersionPattern("4.*.*")),
				},
			},
			want:    nil,
			wantErr: false,
		},
		{
			name: "Filter returns an error",
			args: args{
				versions: models.VersionSlice{
					newVersion("1.2.3"),
					newVersion("1.3.4"),
				},
				filters: []FilterFunc{
					VersionPatternFilter(newVersionPattern("1.1.1")),
					alwaysError(),
				},
			},
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ApplyFilters(tt.args.versions, tt.args.filters...)
			if (err != nil) != tt.wantErr {
				t.Errorf("ApplyFilters() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			assert.Equal(t, tt.want, got)
		})
	}
}

func alwaysError() FilterFunc {
	return func(versions []models.Version) ([]models.Version, error) {
		return nil, errors.New("Filter error")
	}
}

func TestVersionPatternFilter(t *testing.T) {
	tests := []struct {
		name     string
		pattern  models.VersionPattern
		versions []models.Version
		want     []models.Version
	}{
		{
			name:    "Exact match",
			pattern: newVersionPattern("1.2.3-alpha.beta"),
			versions: []models.Version{
				newVersion("1.2.3-alpha.beta"),
				newVersion("2.3.4-alpha.beta"),
			},
			want: []models.Version{
				newVersion("1.2.3-alpha.beta"),
			},
		},
		{
			name:    "Patch wildcard match",
			pattern: newVersionPattern("1.2.*-alpha.beta"),
			versions: []models.Version{
				newVersion("1.2.3-alpha.beta"),
				newVersion("1.2.4-alpha.beta"),
				newVersion("2.3.4-alpha.beta"),
			},
			want: []models.Version{
				newVersion("1.2.3-alpha.beta"),
				newVersion("1.2.4-alpha.beta"),
			},
		},
		{
			name:    "Prerelease wildcard match",
			pattern: newVersionPattern("1.2.3-alpha.*"),
			versions: []models.Version{
				newVersion("1.2.3-alpha.beta"),
				newVersion("1.2.3-alpha.gamma"),
				newVersion("2.3.4-alpha.beta"),
			},
			want: []models.Version{
				newVersion("1.2.3-alpha.beta"),
				newVersion("1.2.3-alpha.gamma"),
			},
		},
		{
			name:    "Multiple wildcard match",
			pattern: newVersionPattern("*.*.*-alpha.*"),
			versions: []models.Version{
				newVersion("1.2.3-alpha.beta"),
				newVersion("1.2.3-alpha.gamma"),
				newVersion("2.3.4-alpha.beta"),
				newVersion("2.3.4-no-match.beta"),
			},
			want: []models.Version{
				newVersion("1.2.3-alpha.beta"),
				newVersion("1.2.3-alpha.gamma"),
				newVersion("2.3.4-alpha.beta"),
			},
		},
		{
			name:    "Build metadata exact match",
			pattern: newVersionPattern("1.2.3-alpha.beta+20130313144700"),
			versions: []models.Version{
				newVersion("1.2.3-alpha.beta+20130313144700"),
				newVersion("1.2.3-alpha.beta+exp.sha.5114f85"),
				newVersion("2.3.4-alpha.beta+20130313144700"),
			},
			want: []models.Version{
				newVersion("1.2.3-alpha.beta+20130313144700"),
			},
		},
		{
			name:    "Build metadata empty",
			pattern: newVersionPattern("*.*.*-alpha.beta"),
			versions: []models.Version{
				newVersion("1.2.3-alpha.beta+20130313144700"),
				newVersion("1.2.3-alpha.beta+exp.sha.5114f85"),
				newVersion("2.3.4-alpha.beta+20130313144700"),
			},
			want: []models.Version{
				newVersion("1.2.3-alpha.beta+20130313144700"),
				newVersion("1.2.3-alpha.beta+exp.sha.5114f85"),
				newVersion("2.3.4-alpha.beta+20130313144700"),
			},
		},
		{
			name:    "Invalid pattern error",
			pattern: newVersionPattern("1.2.3-alpha..beta"),
			versions: []models.Version{
				newVersion("1.2.3-alpha.beta"),
				newVersion("1.2.3-alpha.gamma"),
				newVersion("2.3.4-alpha.beta"),
			},
			want: nil,
		},
		{
			name:    "Different lengths",
			pattern: newVersionPattern("1.2.3-alpha.*"),
			versions: []models.Version{
				newVersion("1.2.3-alpha.beta"),
				newVersion("1.2.3-alpha"),
			},
			want: []models.Version{
				newVersion("1.2.3-alpha.beta"),
			},
		},
		{
			name:    "Release only",
			pattern: newVersionPattern("*.*.*"),
			versions: []models.Version{
				newVersion("1.2.3"),
				newVersion("1.2.4"),
				newVersion("1.2.3+exp.sha.5114f85"),
				newVersion("1.2.3-alpha.beta"),
				newVersion("1.2.3-beta"),
				newVersion("1.2.3-alpha.beta+20130313144700"),
			},
			want: []models.Version{
				newVersion("1.2.3"),
				newVersion("1.2.4"),
				newVersion("1.2.3+exp.sha.5114f85"),
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			filter := VersionPatternFilter(tt.pattern)
			got, err := filter(tt.versions)
			assert.NoError(t, err)
			assert.Equal(t, tt.want, got)
		})
	}
}

func newVersion(s string) models.Version {
	v, _ := models.ParseVersion(s)
	return v
}

func newVersionPattern(s string) models.VersionPattern {
	v, _ := models.ParseVersionPattern(s)
	return v
}
