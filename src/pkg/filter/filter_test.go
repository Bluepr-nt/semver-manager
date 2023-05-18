package filter

import (
	"errors"
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

func TestPreReleaseVersionStream(t *testing.T) {
	type args struct {
		release    models.Release
		prerelease models.PRVersion
	}
	tests := []struct {
		name     string
		args     args
		versions []models.Version
		want     []models.Version
		wantErr  bool
	}{
		{
			name: "Matching prerelease versions",
			args: args{
				release: models.Release{
					Major: 1,
					Minor: 2,
					Patch: 3,
				},
				prerelease: newPRVersion([]string{"alpha", "beta"}),
			},
			versions: []models.Version{
				{
					Release: models.Release{
						Major: 1,
						Minor: 2,
						Patch: 3,
					},
					Prerelease: newPRVersion([]string{"alpha", "beta"}),
				},
				{
					Release: models.Release{
						Major: 1,
						Minor: 2,
						Patch: 3,
					},
					Prerelease: newPRVersion([]string{"alpha"}),
				},
			},
			want: []models.Version{
				{
					Release: models.Release{
						Major: 1,
						Minor: 2,
						Patch: 3,
					},
					Prerelease: newPRVersion([]string{"alpha", "beta"}),
				},
			},
			wantErr: false,
		},
		{
			name: "No matching prerelease versions",
			args: args{
				release: models.Release{
					Major: 1,
					Minor: 2,
					Patch: 3,
				},
				prerelease: newPRVersion([]string{"alpha", "beta"}),
			},
			versions: []models.Version{
				{
					Release: models.Release{
						Major: 1,
						Minor: 2,
						Patch: 3,
					},
					Prerelease: newPRVersion([]string{"gamma"}),
				},
			},
			want:    nil,
			wantErr: false,
		},
		{
			name: "No matching prerelease versions with same amount of identifiers",
			args: args{
				release: models.Release{
					Major: 1,
					Minor: 2,
					Patch: 3,
				},
				prerelease: newPRVersion([]string{"alpha", "beta"}),
			},
			versions: []models.Version{
				{
					Release: models.Release{
						Major: 1,
						Minor: 2,
						Patch: 3,
					},
					Prerelease: newPRVersion([]string{"gamma", "alpha"}),
				},
			},
			want:    nil,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			filter := PreReleaseVersionStream(tt.args.release, tt.args.prerelease)
			got, err := filter(tt.versions)
			if (err != nil) != tt.wantErr {
				t.Errorf("PreReleaseVersionStream() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			assert.Equal(t, tt.want, got)
		})
	}
}

func newPRVersion(identifiers []string) (p models.PRVersion) {
	p, _ = models.NewPRVersion(identifiers)
	return
}

func TestPatchVersionStream(t *testing.T) {
	type args struct {
		major uint64
		minor uint64
		patch uint64
	}
	tests := []struct {
		name     string
		args     args
		versions []models.Version
		want     []models.Version
		wantErr  bool
	}{
		{
			name: "Matching patch versions",
			args: args{
				major: 1,
				minor: 2,
				patch: 3,
			},
			versions: []models.Version{
				{
					Release: models.Release{
						Major: 1,
						Minor: 2,
						Patch: 3,
					},
				},
				{
					Release: models.Release{
						Major: 1,
						Minor: 2,
						Patch: 4,
					},
				},
			},
			want: []models.Version{
				{
					Release: models.Release{
						Major: 1,
						Minor: 2,
						Patch: 3,
					},
				},
			},
			wantErr: false,
		},
		{
			name: "No matching patch versions",
			args: args{
				major: 1,
				minor: 2,
				patch: 3,
			},
			versions: []models.Version{
				{
					Release: models.Release{
						Major: 1,
						Minor: 2,
						Patch: 4,
					},
				},
			},
			want:    nil,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			filter := PatchVersionStream(tt.args.major, tt.args.minor, tt.args.patch)
			got, err := filter(tt.versions)
			if (err != nil) != tt.wantErr {
				t.Errorf("PatchVersionStream() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestMinorVersionStream(t *testing.T) {
	type args struct {
		major uint64
		minor uint64
	}
	tests := []struct {
		name     string
		args     args
		versions []models.Version
		want     []models.Version
		wantErr  bool
	}{
		{
			name: "Matching minor versions",
			args: args{
				major: 1,
				minor: 2,
			},
			versions: []models.Version{
				{
					Release: models.Release{
						Major: 1,
						Minor: 2,
						Patch: 3,
					},
				},
				{
					Release: models.Release{
						Major: 1,
						Minor: 3,
						Patch: 4,
					},
				},
			},
			want: []models.Version{
				{
					Release: models.Release{
						Major: 1,
						Minor: 2,
						Patch: 3,
					},
				},
			},
			wantErr: false,
		},
		{
			name: "No matching minor versions",
			args: args{
				major: 1,
				minor: 2,
			},
			versions: []models.Version{
				{
					Release: models.Release{
						Major: 1,
						Minor: 3,
						Patch: 4,
					},
				},
			},
			want:    nil,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			filter := MinorVersionStream(tt.args.major, tt.args.minor)
			got, err := filter(tt.versions)
			if (err != nil) != tt.wantErr {
				t.Errorf("MinorVersionStream() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestMajorVersionStream(t *testing.T) {
	type args struct {
		major uint64
	}
	tests := []struct {
		name     string
		args     args
		versions []models.Version
		want     []models.Version
		wantErr  bool
	}{
		{
			name: "Matching major versions",
			args: args{
				major: 1,
			},
			versions: []models.Version{
				{
					Release: models.Release{
						Major: 1,
						Minor: 2,
						Patch: 3,
					},
				},
				{
					Release: models.Release{
						Major: 2,
						Minor: 3,
						Patch: 4,
					},
				},
			},
			want: []models.Version{
				{
					Release: models.Release{
						Major: 1,
						Minor: 2,
						Patch: 3,
					},
				},
			},
			wantErr: false,
		},
		{
			name: "No matching major versions",
			args: args{
				major: 1,
			},
			versions: []models.Version{
				{
					Release: models.Release{
						Major: 2,
						Minor: 3,
						Patch: 4,
					},
				},
			},
			want:    nil,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			filter := MajorVersionStream(tt.args.major)
			got, err := filter(tt.versions)
			if (err != nil) != tt.wantErr {
				t.Errorf("MajorVersionStream() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestApplyFilters(t *testing.T) {
	type args struct {
		versions []models.Version
		filters  []FilterFunc
	}
	tests := []struct {
		name    string
		args    args
		want    []models.Version
		wantErr bool
	}{
		{
			name: "Apply multiple filters",
			args: args{
				versions: []models.Version{
					{
						Release: models.Release{
							Major: 1,
							Minor: 2,
							Patch: 3,
						},
					},
					{
						Release: models.Release{
							Major: 1,
							Minor: 3,
							Patch: 4,
						},
					},
					{
						Release: models.Release{
							Major: 2,
							Minor: 3,
							Patch: 4,
						},
					},
				},
				filters: []FilterFunc{
					MajorVersionStream(1),
					MinorVersionStream(1, 2),
				},
			},
			want: []models.Version{
				{
					Release: models.Release{
						Major: 1,
						Minor: 2,
						Patch: 3,
					},
				},
			},
			wantErr: false,
		},
		{
			name: "No matching version after filters",
			args: args{
				versions: []models.Version{
					{
						Release: models.Release{
							Major: 1,
							Minor: 2,
							Patch: 3,
						},
					},
					{
						Release: models.Release{
							Major: 1,
							Minor: 3,
							Patch: 4,
						},
					},
					{
						Release: models.Release{
							Major: 2,
							Minor: 3,
							Patch: 4,
						},
					},
				},
				filters: []FilterFunc{
					MajorVersionStream(3),
				},
			},
			want:    nil,
			wantErr: false,
		},
		{
			name: "Filter returns an error",
			args: args{
				versions: []models.Version{
					{
						Release: models.Release{
							Major: 1,
							Minor: 2,
							Patch: 3,
						},
					},
					{
						Release: models.Release{
							Major: 1,
							Minor: 3,
							Patch: 4,
						},
					},
				},
				filters: []FilterFunc{
					MajorVersionStream(1),
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
