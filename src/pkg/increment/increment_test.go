package increment

import (
	"src/cmd/smgr/models"
	"src/cmd/smgr/testutils"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetIncrementType(t *testing.T) {
	type args struct {
		highestRelease    models.Version
		highestPrerelease models.Version
	}
	tests := []struct {
		name string
		args args
		want models.Increment
	}{
		{
			name: "Both versions are the same",
			args: args{
				highestRelease:    testutils.NewVersion("1.0.0"),
				highestPrerelease: testutils.NewVersion("1.0.0-alpha"),
			},
			want: models.None,
		},
		{
			name: "Prerelease version is lower",
			args: args{
				highestRelease:    testutils.NewVersion("1.1.0"),
				highestPrerelease: testutils.NewVersion("1.0.0-alpha"),
			},
			want: models.None,
		},
		{
			name: "Prerelease version has higher patch",
			args: args{
				highestRelease:    testutils.NewVersion("1.0.0"),
				highestPrerelease: testutils.NewVersion("1.0.1-alpha"),
			},
			want: models.Patch,
		},
		{
			name: "Prerelease version has lower patch",
			args: args{
				highestRelease:    testutils.NewVersion("1.0.1"),
				highestPrerelease: testutils.NewVersion("1.0.0-alpha"),
			},
			want: models.None,
		},
		{
			name: "Prerelease version has higher minor",
			args: args{
				highestRelease:    testutils.NewVersion("1.0.0"),
				highestPrerelease: testutils.NewVersion("1.1.0-alpha"),
			},
			want: models.Minor,
		},
		{
			name: "Prerelease version has lower minor",
			args: args{
				highestRelease:    testutils.NewVersion("1.1.0"),
				highestPrerelease: testutils.NewVersion("1.0.0-alpha"),
			},
			want: models.None,
		},
		{
			name: "Prerelease version has higher major",
			args: args{
				highestRelease:    testutils.NewVersion("1.0.0"),
				highestPrerelease: testutils.NewVersion("2.0.0-alpha"),
			},
			want: models.Major,
		},
		{
			name: "Prerelease version has lower major",
			args: args{
				highestRelease:    testutils.NewVersion("2.0.0"),
				highestPrerelease: testutils.NewVersion("1.0.0-alpha"),
			},
			want: models.None,
		},
		{
			name: "Prerelease version has higher major and minor",
			args: args{
				highestRelease:    testutils.NewVersion("1.0.0"),
				highestPrerelease: testutils.NewVersion("2.1.0-alpha"),
			},
			want: models.Major,
		},
		{
			name: "Prerelease version has higher major, minor and patch",
			args: args{
				highestRelease:    testutils.NewVersion("1.0.0"),
				highestPrerelease: testutils.NewVersion("2.1.1-alpha"),
			},
			want: models.Major,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := GetIncrementType(tt.args.highestRelease, tt.args.highestPrerelease)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestCalculateIncrementTypeForNewPrerelease(t *testing.T) {
	type args struct {
		highestRelease     models.Version
		highestPrerelease  models.Version
		requestedIncrement models.Increment
	}
	tests := []struct {
		name string
		args args
		want models.Increment
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := CalculateIncrementTypeForNewPrerelease(tt.args.highestRelease, tt.args.highestPrerelease, tt.args.requestedIncrement)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestIncrementRelease(t *testing.T) {
	type args struct {
		sourceVersion models.Version
		increment     models.Increment
	}
	tests := []struct {
		name string
		args args
		want models.Version
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := IncrementRelease(tt.args.sourceVersion, tt.args.increment)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestIncrementReleaseFromStream(t *testing.T) {
	type args struct {
		sourceVersions []models.Version
		streamPattern  models.VersionPattern
		increment      models.Increment
	}
	tests := []struct {
		name    string
		args    args
		want    models.Version
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := IncrementReleaseFromStream(tt.args.sourceVersions, tt.args.streamPattern, tt.args.increment)
			if tt.wantErr {
				assert.Error(t, err)
				assert.Equal(t, tt.want, got)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.want, got)
			}
		})
	}
}

func TestNumericalIncrement(t *testing.T) {
	type args struct {
		sourceIdentifier models.PRIdentifier
	}
	tests := []struct {
		name    string
		args    args
		want    models.PRIdentifier
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NumericalIncrement(tt.args.sourceIdentifier)
			if tt.wantErr {
				assert.Error(t, err)
				assert.Equal(t, tt.want, got)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.want, got)
			}
		})
	}
}

func TestAlphabeticalIncrement(t *testing.T) {
	type args struct {
		sourceIdentifier models.PRIdentifier
	}
	tests := []struct {
		name    string
		args    args
		want    models.PRIdentifier
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := AlphabeticalIncrement(tt.args.sourceIdentifier)
			if tt.wantErr {
				assert.Error(t, err)
				assert.Equal(t, tt.want, got)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.want, got)
			}
		})
	}
}

func TestPromoteVersion(t *testing.T) {
	type args struct {
		sourceVersion models.Version
		targetStream  models.VersionPattern
	}
	tests := []struct {
		name string
		args args
		want models.Version
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := PromoteVersion(tt.args.sourceVersion, tt.args.targetStream)
			assert.Equal(t, tt.want, got)
		})
	}
}
