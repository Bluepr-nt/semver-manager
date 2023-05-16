package filter

import (
	"src/pkg/fetch/models"
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
				{
					Release: models.Release{
						Major: 0,
						Minor: 0,
						Patch: 0,
					},
					Prerelease: models.PRVersion{
						Identifiers: []models.PRIdentifier{
							newPrIdentifier("1"),
						},
					},
					Build: models.BuildMetadata{
						Identifiers: []models.BuildIdentifier{
							newBuildIdentifier("054"),
						},
					},
				},
				{
					Release: models.Release{
						Major: 0,
						Minor: 1,
						Patch: 0,
					},
					Prerelease: models.PRVersion{
						Identifiers: []models.PRIdentifier{
							newPrIdentifier("1"),
						},
					},
					Build: models.BuildMetadata{
						Identifiers: []models.BuildIdentifier{
							newBuildIdentifier("054"),
						},
					},
				},
			},
			want:    "0.1.0-1+054",
			wantErr: false,
		},
		{
			name: "Mixed release and prerelease versions",
			versions: []models.Version{
				{
					Release: models.Release{
						Major: 1,
						Minor: 0,
						Patch: 0,
					},
				},
				{
					Release: models.Release{
						Major: 0,
						Minor: 0,
						Patch: 0,
					},
					Prerelease: models.PRVersion{
						Identifiers: []models.PRIdentifier{
							newPrIdentifier("2"),
						},
					},
				},
			},
			want:    "1.0.0",
			wantErr: false,
		},
		{
			name: "Prerelease identifiers with different lengths",
			versions: []models.Version{
				{
					Release: models.Release{
						Major: 0,
						Minor: 1,
						Patch: 0,
					},
					Prerelease: models.PRVersion{
						Identifiers: []models.PRIdentifier{
							newPrIdentifier("alpha"),
							newPrIdentifier("2"),
						},
					},
				},
				{
					Release: models.Release{
						Major: 0,
						Minor: 1,
						Patch: 0,
					},
					Prerelease: models.PRVersion{
						Identifiers: []models.PRIdentifier{
							newPrIdentifier("alpha"),
						},
					},
				},
			},
			want:    "0.1.0-alpha.2",
			wantErr: false,
		},
		{
			name: "Versions with different build metadata",
			versions: []models.Version{
				{
					Release: models.Release{
						Major: 0,
						Minor: 1,
						Patch: 0,
					},
					Build: models.BuildMetadata{
						Identifiers: []models.BuildIdentifier{
							newBuildIdentifier("100"),
						},
					},
				},
				{
					Release: models.Release{
						Major: 0,
						Minor: 1,
						Patch: 0,
					},
					Build: models.BuildMetadata{
						Identifiers: []models.BuildIdentifier{
							newBuildIdentifier("200"),
						},
					},
				},
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
	i, _ := models.NewPrIdentifier(v)
	return i
}

func newBuildIdentifier(v string) models.BuildIdentifier {
	i, _ := models.NewBuildIdentifier(v)
	return i
}

func TestReleaseOnly(t *testing.T) {
	tests := []struct {
		name     string
		versions []models.Version
		want     []models.Version
		wantErr  bool
	}{
		{
			name: "Mixed release and prerelease versions",
			versions: []models.Version{
				{
					Release: models.Release{
						Major: 1,
						Minor: 0,
						Patch: 0,
					},
				},
				{
					Release: models.Release{
						Major: 1,
						Minor: 0,
						Patch: 1,
					},
					Prerelease: models.PRVersion{
						Identifiers: []models.PRIdentifier{
							newPrIdentifier("alpha"),
						},
					},
				},
			},
			want: []models.Version{
				{
					Release: models.Release{
						Major: 1,
						Minor: 0,
						Patch: 0,
					},
				},
			},
			wantErr: false,
		},
		{
			name: "Only release versions",
			versions: []models.Version{
				{
					Release: models.Release{
						Major: 0,
						Minor: 1,
						Patch: 0,
					},
				},
				{
					Release: models.Release{
						Major: 0,
						Minor: 2,
						Patch: 0,
					},
				},
			},
			want: []models.Version{
				{
					Release: models.Release{
						Major: 0,
						Minor: 1,
						Patch: 0,
					},
				},
				{
					Release: models.Release{
						Major: 0,
						Minor: 2,
						Patch: 0,
					},
				},
			},
			wantErr: false,
		},
		{
			name:     "No versions",
			versions: []models.Version{},
			want:     nil,
			wantErr:  false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ReleaseOnly()(tt.versions)
			if (err != nil) != tt.wantErr {
				t.Errorf("ReleaseOnly() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			assert.Equal(t, tt.want, got)
		})
	}
}
