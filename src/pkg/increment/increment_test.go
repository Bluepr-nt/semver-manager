package increment

import (
	"reflect"
	"src/cmd/smgr/models"
	"testing"
)

func TestGetReleaseToPreReleaseIncrementType(t *testing.T) {
	type args struct {
		highestRelease    models.Version
		highestPrerelease models.Version
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
			if got := GetReleaseToPreReleaseIncrementType(tt.args.highestRelease, tt.args.highestPrerelease); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetReleaseToPreReleaseIncrementType() = %v, want %v", got, tt.want)
			}
		})
	}
}
