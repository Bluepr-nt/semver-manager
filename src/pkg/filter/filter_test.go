package filter

import (
	"reflect"
	"src/pkg/fetch/models"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestHighest(t *testing.T) {
	tests := []struct {
		name     string
		versions []models.Version
		want     string
		err      error
	}{
		{
			name: "Empty tags list",
			versions: []models.Version{
				{
					Major: 0,
					Minor: 0,
					Patch: 0,
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
					Major: 0,
					Minor: 1,
					Patch: 0,
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
			want: "",
			err:  nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Highest()(tt.versions)
			assert.NoError(t, err, "Expected no error")
			assert.Equal(t, tt.want, got[0].String())
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
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ReleaseOnly(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ReleaseOnly() = %v, want %v", got, tt.want)
			}
		})
	}
}
