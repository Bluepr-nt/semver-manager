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
		{
			name: "Simple Patch increment",
			args: args{
				highestRelease:     testutils.NewVersion("1.0.0"),
				highestPrerelease:  testutils.NewVersion("1.0.0-alpha"),
				requestedIncrement: models.Patch,
			},
			want: models.Patch,
		},
		{
			name: "No increment",
			args: args{
				highestRelease:     testutils.NewVersion("1.0.0"),
				highestPrerelease:  testutils.NewVersion("1.0.1-alpha"),
				requestedIncrement: models.Patch,
			},
			want: models.None,
		},
		{
			name: "Minor increment",
			args: args{
				highestRelease:     testutils.NewVersion("1.0.0"),
				highestPrerelease:  testutils.NewVersion("1.0.1-alpha"),
				requestedIncrement: models.Minor,
			},
			want: models.Minor,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := CalculateIncrementTypeForNewPrerelease(tt.args.highestRelease, tt.args.highestPrerelease, tt.args.requestedIncrement)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestIncrementRelease(t *testing.T) {
	tests := []struct {
		name          string
		sourceVersion models.Version
		increment     models.Increment
		want          models.Version
	}{
		{
			name:          "Major increment",
			sourceVersion: testutils.NewVersion("1.0.0"),
			increment:     models.Major,
			want:          testutils.NewVersion("2.0.0"),
		},
		{
			name:          "Minor increment",
			sourceVersion: testutils.NewVersion("1.0.0"),
			increment:     models.Minor,
			want:          testutils.NewVersion("1.1.0"),
		},
		{
			name:          "Patch increment",
			sourceVersion: testutils.NewVersion("1.0.0"),
			increment:     models.Patch,
			want:          testutils.NewVersion("1.0.1"),
		},
		{
			name:          "None increment",
			sourceVersion: testutils.NewVersion("1.0.0"),
			increment:     models.None,
			want:          testutils.NewVersion("1.0.0"),
		},
		{
			name:          "First version",
			sourceVersion: models.Version{},
			increment:     models.None,
			want:          testutils.NewVersion("0.0.0"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := IncrementRelease(tt.sourceVersion, tt.increment)
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
		{
			name: "Increment on release stream",
			args: args{
				sourceVersions: []models.Version{
					testutils.NewVersion("1.0.0"),
					testutils.NewVersion("0.1.0"),
					testutils.NewVersion("2.0.0"),
				},
				streamPattern: testutils.NewVersionPattern("1.*.*"),
				increment:     models.Minor,
			},
			want:    testutils.NewVersion("1.1.0"),
			wantErr: false,
		},
		{
			name: "No matching source version, new stream",
			args: args{
				sourceVersions: []models.Version{
					testutils.NewVersion("1.0.0"),
				},
				streamPattern: testutils.NewVersionPattern("2.*.*"),
				increment:     models.Minor,
			},
			want:    testutils.NewVersion("2.0.0"),
			wantErr: false,
		},
		{
			name: "Increment on alpha stream",
			args: args{
				sourceVersions: []models.Version{
					testutils.NewVersion("1.0.0"),
					testutils.NewVersion("1.1.0-alpha.1"),
					testutils.NewVersion("1.0.0-alpha.1"),
				},
				streamPattern: testutils.NewVersionPattern("1.*.*-alpha.*"),
				increment:     models.Minor,
			},
			want:    testutils.NewVersion("0.0.0"),
			wantErr: true,
		},
		{
			name: "Increment major version",
			args: args{
				sourceVersions: []models.Version{
					testutils.NewVersion("1.0.0"),
				},
				streamPattern: testutils.NewVersionPattern("1.*.*"),
				increment:     models.Major,
			},
			want:    testutils.NewVersion("2.0.0"),
			wantErr: false,
		},
		{
			name: "Increment patch version",
			args: args{
				sourceVersions: []models.Version{
					testutils.NewVersion("1.0.0"),
				},
				streamPattern: testutils.NewVersionPattern("1.*.*"),
				increment:     models.Patch,
			},
			want:    testutils.NewVersion("1.0.1"),
			wantErr: false,
		},
		{
			name: "Weird version",
			args: args{
				sourceVersions: []models.Version{
					testutils.NewVersion("2.1.000000000000000000000000"),
					testutils.NewVersion("2.0.0"),
				},
				streamPattern: testutils.NewVersionPattern("2.*.*"),
				increment:     models.Patch,
			},
			want:    testutils.NewVersion("2.0.1"),
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := IncrementReleaseToStream(tt.args.sourceVersions, tt.args.streamPattern, tt.args.increment)
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

func TestNumericalPRIncrement(t *testing.T) {

	tests := []struct {
		name             string
		sourceIdentifier models.PRIdentifier
		want             models.PRIdentifier
		wantErr          bool
	}{
		{
			name:             "Simple increment",
			sourceIdentifier: testutils.NewPRIdentifier("0"),
			want:             testutils.NewPRIdentifier("1"),
			wantErr:          false,
		},
		{
			name:             "Source Identifier too big",
			sourceIdentifier: testutils.NewPRIdentifier("18446744073709551616"),
			want:             models.PRIdentifier{},
			wantErr:          true,
		},
		{
			name:             "Source Identifier is not a number",
			sourceIdentifier: testutils.NewPRIdentifier("a"),
			want:             models.PRIdentifier{},
			wantErr:          true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NumericalPRIncrement(tt.sourceIdentifier)
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

	tests := []struct {
		name             string
		sourceIdentifier models.PRIdentifier
		want             models.PRIdentifier
		wantErr          bool
	}{
		{
			name:             "Simple Increment",
			sourceIdentifier: testutils.NewPRIdentifier("a"),
			want:             testutils.NewPRIdentifier("b"),
			wantErr:          false,
		},
		{
			name:             "Simple Increment with Caps",
			sourceIdentifier: testutils.NewPRIdentifier("A"),
			want:             testutils.NewPRIdentifier("B"),
			wantErr:          false,
		},
		{
			name:             "Error expected a single character identifier",
			sourceIdentifier: testutils.NewPRIdentifier("AA"),
			want:             models.PRIdentifier{},
			wantErr:          true,
		},
		{
			name:             "Error expected an alphabetical identifier",
			sourceIdentifier: testutils.NewPRIdentifier("0"),
			want:             models.PRIdentifier{},
			wantErr:          true,
		},
		{
			name:             "Increment from z to za",
			sourceIdentifier: testutils.NewPRIdentifier("z"),
			want:             testutils.NewPRIdentifier("za"),
			wantErr:          false,
		},
		{
			name:             "Increment from Z to ZA",
			sourceIdentifier: testutils.NewPRIdentifier("Z"),
			want:             testutils.NewPRIdentifier("ZA"),
			wantErr:          false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := AlphabeticalIncrement(tt.sourceIdentifier)
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

func TestPromotePRVersion(t *testing.T) {

	tests := []struct {
		name          string
		sourceVersion models.Version
		targetStream  models.VersionPattern
		versions      []models.Version
		want          models.Version
		wantErr       bool
	}{
		{
			name:          "Promote from Alpha to Beta with digit",
			sourceVersion: testutils.NewVersion("1.0.0-Alpha.1"),
			targetStream:  testutils.NewVersionPattern("1.0.0-Beta.*"),
			versions: []models.Version{
				testutils.NewVersion("1.0.0-Alpha.0"),
				testutils.NewVersion("1.0.0-Beta.0"),
			},
			want: testutils.NewVersion("1.0.0-Beta.1"),
		},
		{
			name:          "Error pattern is a release pattern",
			sourceVersion: testutils.NewVersion("1.0.1-Alpha"),
			targetStream:  testutils.NewVersionPattern("1.0.*"),
			versions: []models.Version{
				testutils.NewVersion("1.0.0-Alpha"),
				testutils.NewVersion("1.0.0-Beta"),
				testutils.NewVersion("1.0.0"),
			},
			want:    models.Version{},
			wantErr: true,
		},
		{
			name:          "Promote from Alpha to Beta with digit first version on stream",
			sourceVersion: testutils.NewVersion("1.0.0-Alpha.1"),
			targetStream:  testutils.NewVersionPattern("1.0.0-Beta.*"),
			versions: []models.Version{
				testutils.NewVersion("1.0.0-Alpha.0"),
			},
			want: testutils.NewVersion("1.0.0-Beta.0"),
		},
		{
			name:          "Promote from Alpha to Beta -- no digit",
			sourceVersion: testutils.NewVersion("1.0.0-Alpha"),
			targetStream:  testutils.NewVersionPattern("1.0.0-Beta"),
			versions: []models.Version{
				testutils.NewVersion("0.1.0-Alpha"),
				testutils.NewVersion("0.1.0-Beta"),
			},
			want: testutils.NewVersion("1.0.0-Beta"),
		},
		{
			name:          "Promote from Alpha to Beta -- no digit and existing version",
			sourceVersion: testutils.NewVersion("1.0.0-Alpha"),
			targetStream:  testutils.NewVersionPattern("1.0.0-Beta"),
			versions: []models.Version{
				testutils.NewVersion("0.1.0-Alpha"),
				testutils.NewVersion("0.1.0-Beta"),
				testutils.NewVersion("1.0.0-Alpha"),
				testutils.NewVersion("1.0.0-Beta"),
			},
			want: testutils.NewVersion("1.0.0-Beta.0"),
		},
		{
			name:          "Promote from Alpha to Longer Beta-teta release",
			sourceVersion: testutils.NewVersion("1.0.1-Alpha"),
			targetStream:  testutils.NewVersionPattern("1.0.0-Beta.teta.*"),
			versions: []models.Version{
				testutils.NewVersion("1.0.0-Alpha"),
				testutils.NewVersion("1.0.0-Beta"),
				testutils.NewVersion("1.0.0"),
				testutils.NewVersion("1.0.0-Beta.teta.0"),
			},
			want: testutils.NewVersion("1.0.0-Beta.teta.1"),
		},
		{
			name:          "Promote to stream with longer Prerelease than stream",
			sourceVersion: testutils.NewVersion("1.0.0-Alpha.1"),
			targetStream:  testutils.NewVersionPattern("1.0.0-Beta.*"),
			versions: []models.Version{
				testutils.NewVersion("1.0.0-Beta.Alpha.0"),
			},
			want: testutils.NewVersion("1.0.0-Beta.0"),
		},

		{
			name:          "Promote to loose prerelease stream with some PreRelease identifiers set by source version",
			sourceVersion: testutils.NewVersion("1.0.0-Alpha.Beta.0"),
			targetStream:  testutils.NewVersionPattern("1.0.0-*.Beta.*"),
			versions: []models.Version{
				testutils.NewVersion("1.0.0-Alpha.0"),
				testutils.NewVersion("1.0.0-Beta.0"),
			},
			want: testutils.NewVersion("1.0.0-Alpha.Beta.1"),
		},
		{
			name:          "Promote to loose prerelease stream with matching prerelease versions",
			sourceVersion: testutils.NewVersion("1.0.0-Alpha"),
			targetStream:  testutils.NewVersionPattern("1.0.0-*.Beta.*"),
			versions: []models.Version{
				testutils.NewVersion("1.0.0-Alpha.Beta.0"),
				testutils.NewVersion("1.0.0-Teta.Beta.0"),
			},
			want: testutils.NewVersion("1.0.0-Teta.Beta.1"),
		},
		{
			name:          "Promote to multi-wildcard prerelease stream with matching prerelease versions with more identifiers than pattern",
			sourceVersion: testutils.NewVersion("1.0.0-Alpha"),
			targetStream:  testutils.NewVersionPattern("1.0.0-*.Beta.*"),
			versions: []models.Version{
				testutils.NewVersion("1.0.0-Alpha.Beta.0"),
				testutils.NewVersion("1.0.0-Teta.Beta.0.0"),
			},
			want: testutils.NewVersion("1.0.0-Alpha.Beta.1"),
		},
		{
			name:          "Promote to loose prerelease stream",
			sourceVersion: testutils.NewVersion("1.0.0-Alpha"),
			targetStream:  testutils.NewVersionPattern("1.0.0-*.*"),
			versions: []models.Version{
				testutils.NewVersion("1.0.0-Alpha.0"),
				testutils.NewVersion("1.0.0-Beta.0"),
				testutils.NewVersion("1.0.0-Beta.0.1"),
			},
			want: testutils.NewVersion("1.0.0-Beta.1"),
		},
		{
			name:          "Promote to very loose prerelease stream",
			sourceVersion: testutils.NewVersion("1.0.0-Alpha"),
			targetStream:  testutils.NewVersionPattern("1.0.0-*.*.*"),
			versions: []models.Version{
				testutils.NewVersion("1.0.0-Alpha.0"),
				testutils.NewVersion("1.0.0-Beta.0"),
			},
			want: testutils.NewVersion("1.0.0-0.0.0"),
		},
		{
			name:          "Promote to very loose prerelease stream",
			sourceVersion: testutils.NewVersion("1.0.0-Alpha"),
			targetStream:  testutils.NewVersionPattern("1.0.0-*.*.*"),
			versions: []models.Version{
				testutils.NewVersion("1.0.0-Alpha.0"),
				testutils.NewVersion("1.0.0-Beta.0.1"),
			},
			want: testutils.NewVersion("1.0.0-Beta.0.2"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			versions := append(tt.versions, tt.sourceVersion)
			got, err := IncrementPReleaseToStream(versions, tt.targetStream, models.None)
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

func TestIncrementVersion(t *testing.T) {

	tests := []struct {
		name                   string
		sourceVersions         []models.Version
		streamPattern          models.VersionPattern
		increment              models.Increment
		wantIncrementedVersion models.Version
		wantErr                bool
	}{
		{
			name: "Increment patch version",
			sourceVersions: []models.Version{
				testutils.NewVersion("1.0.0"),
				testutils.NewVersion("1.0.1"),
				testutils.NewVersion("1.0.2"),
			},
			streamPattern:          testutils.NewVersionPattern("1.0.*"),
			increment:              models.Patch,
			wantIncrementedVersion: testutils.NewVersion("1.0.3"),
			wantErr:                false,
		},
		{
			name: "Increment minor version",
			sourceVersions: []models.Version{
				testutils.NewVersion("1.0.0"),
				testutils.NewVersion("1.1.0"),
				testutils.NewVersion("1.2.0"),
			},
			streamPattern:          testutils.NewVersionPattern("1.*.*"),
			increment:              models.Minor,
			wantIncrementedVersion: testutils.NewVersion("1.3.0"),
			wantErr:                false,
		},
		{
			name: "Increment major version",
			sourceVersions: []models.Version{
				testutils.NewVersion("1.0.0"),
				testutils.NewVersion("2.0.0"),
				testutils.NewVersion("3.0.0"),
			},
			streamPattern:          testutils.NewVersionPattern("*.0.0"),
			increment:              models.Major,
			wantIncrementedVersion: testutils.NewVersion("4.0.0"),
			wantErr:                false,
		},
		{
			name:                   "Increment with empty source versions",
			sourceVersions:         []models.Version{},
			streamPattern:          testutils.NewVersionPattern("1.0.*"),
			increment:              models.Patch,
			wantIncrementedVersion: testutils.NewVersion("1.0.0"),
			wantErr:                false,
		},
		{
			name: "Increment with no matching stream pattern",
			sourceVersions: []models.Version{
				testutils.NewVersion("1.0.0"),
				testutils.NewVersion("1.0.1"),
				testutils.NewVersion("1.0.2"),
			},
			streamPattern:          testutils.NewVersionPattern("2.0.*"),
			increment:              models.Patch,
			wantIncrementedVersion: testutils.NewVersion("2.0.0"),
			wantErr:                false,
		},
		{
			name:                   "Increment with no source versions and no stream pattern",
			sourceVersions:         []models.Version{},
			streamPattern:          models.VersionPattern{},
			increment:              models.Patch,
			wantIncrementedVersion: testutils.NewVersion("0.0.0"),
			wantErr:                false,
		},
		{
			name: "Increment with PreRelease stream pattern",
			sourceVersions: []models.Version{
				testutils.NewVersion("1.0.0-alpha.1"),
				testutils.NewVersion("1.0.0-alpha.2"),
				testutils.NewVersion("1.0.0-alpha.3"),
			},
			streamPattern:          testutils.NewVersionPattern("1.0.*-alpha.*"),
			increment:              models.None,
			wantIncrementedVersion: testutils.NewVersion("1.0.0-alpha.4"),
			wantErr:                false,
		},
		{
			name:                   "Increment with PreRelease stream pattern and no source versions",
			sourceVersions:         []models.Version{},
			streamPattern:          testutils.NewVersionPattern("1.0.*-alpha.*"),
			increment:              models.None,
			wantIncrementedVersion: testutils.NewVersion("1.0.0-alpha.0"),
			wantErr:                false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotIncrementedVersion, err := IncrementVersion(tt.sourceVersions, tt.streamPattern, tt.increment)
			if tt.wantErr {
				assert.Error(t, err)
				assert.Equal(t, models.Version{}, gotIncrementedVersion)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.wantIncrementedVersion, gotIncrementedVersion)
			}
		})
	}
}
