package semverSvc

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_FilterHighestSemver(t *testing.T) {
	tests := []struct {
		name    string
		fields  SemverSvc
		tags    []string
		want    []string
		wantErr bool
	}{
		{
			name: "Empty tags list",
			fields: SemverSvc{
				client: nil,
			},
			tags:    []string{},
			want:    []string{},
			wantErr: true,
		},
		{
			name: "Non-SemVer compliant tags",
			fields: SemverSvc{
				client: nil,
			},

			tags: []string{"1.0.0", "1.1.0", "1.2.0", "not-a-tag", "1.2.1"},

			want:    []string{"1.2.1"},
			wantErr: false,
		},
		{
			name: "SemVer compliant tags",
			fields: SemverSvc{
				client: nil,
			},

			tags: []string{"1.0.0", "1.1.0", "1.2.0", "1.2.1"},

			want:    []string{"1.2.1"},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			svSvc := &SemverSvc{
				client: tt.fields.client,
			}
			got, err := svSvc.FilterSemverTags(tt.tags, &Filters{Highest: true})

			if len(got) == 0 && len(tt.want) == 0 {
				return
			}
			if tt.wantErr {
				assert.Error(t, err, "svSvc.FilterHighestSemver() error = %v, wantErr %v", err, tt.wantErr)
			} else {
				assert.NoError(t, err, "svSvc.FilterHighestSemver() error = %v, wantErr %v", err, tt.wantErr)
			}

			assert.Equal(t, tt.want, got, "svSvc.FilterHighestSemver() = %v, want %v", got, tt.want)
		})
	}
}

func Test_IsSemver(t *testing.T) {
	type args struct {
		version string
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{name: "Valid Semantic Version: 0.0.4", args: args{version: "0.0.4"}, want: true},
		{name: "Valid Semantic Version: 1.2.3", args: args{version: "1.2.3"}, want: true},
		{name: "Valid Semantic Version: 10.20.30", args: args{version: "10.20.30"}, want: true},
		{name: "Valid Semantic Version: 1.1.2-prerelease+meta", args: args{version: "1.1.2-prerelease+meta"}, want: true},
		{name: "Valid Semantic Version: 1.1.2+meta", args: args{version: "1.1.2+meta"}, want: true},
		{name: "Valid Semantic Version: 1.1.2+meta-valid", args: args{version: "1.1.2+meta-valid"}, want: true},
		{name: "Valid Semantic Version: 1.0.0-alpha", args: args{version: "1.0.0-alpha"}, want: true},
		{name: "Valid Semantic Version: 1.0.0-beta", args: args{version: "1.0.0-beta"}, want: true},
		{name: "Valid Semantic Version: 1.0.0-alpha.beta", args: args{version: "1.0.0-alpha.beta"}, want: true},
		{name: "Valid Semantic Version: 1.0.0-alpha.beta.1", args: args{version: "1.0.0-alpha.beta.1"}, want: true},
		{name: "Valid Semantic Version: 1.0.0-alpha.1", args: args{version: "1.0.0-alpha.1"}, want: true},
		{name: "Valid Semantic Version: 1.0.0-alpha0.valid", args: args{version: "1.0.0-alpha0.valid"}, want: true},
		{name: "Valid Semantic Version: 1.0.0-alpha.0valid", args: args{version: "1.0.0-alpha.0valid"}, want: true},
		{name: "Valid Semantic Version: 1.0.0-alpha-a.b-c-somethinglong+build.1-aef.1-its-okay", args: args{version: "1.0.0-alpha-a.b-c-somethinglong+build.1-aef.1-its-okay"}, want: true},
		{name: "Valid Semantic Version: 1.0.0-rc.1+build.1", args: args{version: "1.0.0-rc.1+build.1"}, want: true},
		{name: "Valid Semantic Version: 2.0.0-rc.1+build.123", args: args{version: "2.0.0-rc.1+build.123"}, want: true},
		{name: "Valid Semantic Version: 1.2.3-beta", args: args{version: "1.2.3-beta"}, want: true},
		{name: "Valid Semantic Version: 10.2.3-DEV-SNAPSHOT", args: args{version: "10.2.3-DEV-SNAPSHOT"}, want: true},
		{name: "Valid Semantic Version: 1.2.3-SNAPSHOT-123", args: args{version: "1.2.3-SNAPSHOT-123"}, want: true},
		{name: "Valid Semantic Version: 1.0.0", args: args{version: "1.0.0"}, want: true},
		{name: "Valid Semantic Version: 2.0.0", args: args{version: "2.0.0"}, want: true},
		{name: "Valid Semantic Version: 1.1.7", args: args{version: "1.1.7"}, want: true},
		{name: "Valid Semantic Version: 2.0.0+build.1848", args: args{version: "2.0.0+build.1848"}, want: true},
		{name: "Valid Semantic Version: 2.0.1-alpha.1227", args: args{version: "2.0.1-alpha.1227"}, want: true},
		{name: "Valid Semantic Version: 1.0.0-alpha+beta", args: args{"1.0.0-alpha+beta"}, want: true},
		{name: "Valid Semantic Version: 1.2.3----RC-SNAPSHOT.12.9.1--.12+788", args: args{"1.2.3----RC-SNAPSHOT.12.9.1--.12+788"}, want: true},
		{name: "Valid Semantic Version: 1.2.3----R-S.12.9.1--.12+meta", args: args{"1.2.3----R-S.12.9.1--.12+meta"}, want: true},
		{name: "Valid Semantic Version: 1.2.3----RC-SNAPSHOT.12.9.1--.12", args: args{"1.2.3----RC-SNAPSHOT.12.9.1--.12"}, want: true},
		{name: "Valid Semantic Version: 1.0.0+0.build.1-rc.10000aaa-kk-0.1", args: args{"1.0.0+0.build.1-rc.10000aaa-kk-0.1"}, want: true},
		{name: "Valid Semantic Version: 99999999999999999999999.999999999999999999.99999999999999999", args: args{"99999999999999999999999.999999999999999999.99999999999999999"}, want: true},
		{name: "Valid Semantic Version: 1.0.0-0A.is.legal", args: args{"1.0.0-0A.is.legal"}, want: true},
		{name: "Invalid Semantic Version: 1", args: args{"1"}, want: false},
		{name: "Invalid Semantic Version: 1.2", args: args{"1.2"}, want: false},
		{name: "Invalid Semantic Version: 1.2.3-0123", args: args{"1.2.3-0123"}, want: false},
		{name: "Invalid Semantic Version: 1.2.3-0123.0123", args: args{"1.2.3-0123.0123"}, want: false},
		{name: "Invalid Semantic Version: 1.1.2+.123", args: args{"1.1.2+.123"}, want: false},
		{name: "Invalid Semantic Version: +invalid", args: args{"+invalid"}, want: false},
		{name: "Invalid Semantic Version: -invalid", args: args{"-invalid"}, want: false},
		{name: "Invalid Semantic Version: -invalid+invalid", args: args{"-invalid+invalid"}, want: false},
		{name: "Invalid Semantic Version: -invalid.01", args: args{"-invalid.01"}, want: false},
		{name: "Invalid Semantic Version: alpha", args: args{"alpha"}, want: false},
		{name: "Invalid Semantic Version: alpha.beta", args: args{"alpha.beta"}, want: false},
		{name: "Invalid Semantic Version: alpha.beta.1", args: args{"alpha.beta.1"}, want: false},
		{name: "Invalid Semantic Version: alpha.1", args: args{"alpha.1"}, want: false},
		{name: "Invalid Semantic Version: alpha+beta", args: args{"alpha+beta"}, want: false},
		{name: "Invalid Semantic Version: alpha_beta", args: args{"alpha_beta"}, want: false},
		{name: "Invalid Semantic Version: alpha.", args: args{"alpha."}, want: false},
		{name: "Invalid Semantic Version: alpha..", args: args{"alpha.."}, want: false},
		{name: "Invalid Semantic Version: beta", args: args{"beta"}, want: false},
		{name: "Invalid Semantic Version: 1.0.0-alpha_beta", args: args{"1.0.0-alpha_beta"}, want: false},
		{name: "Invalid Semantic Version: -alpha.", args: args{"-alpha."}, want: false},
		{name: "Invalid Semantic Version: 1.0.0-alpha..", args: args{"1.0.0-alpha.."}, want: false},
		{name: "Invalid Semantic Version: 1.0.0-alpha..1", args: args{"1.0.0-alpha..1"}, want: false},
		{name: "Invalid Semantic Version: 1.0.0-alpha...1", args: args{"1.0.0-alpha...1"}, want: false},
		{name: "Invalid Semantic Version: 1.0.0-alpha....1", args: args{"1.0.0-alpha....1"}, want: false},
		{name: "Invalid Semantic Version: 1.0.0-alpha.....1", args: args{"1.0.0-alpha.....1"}, want: false},
		{name: "Invalid Semantic Version: 1.0.0-alpha......1", args: args{"1.0.0-alpha......1"}, want: false},
		{name: "Invalid Semantic Version: 1.0.0-alpha.......1", args: args{"1.0.0-alpha.......1"}, want: false},
		{name: "Invalid Semantic Version: 01.1.1", args: args{"01.1.1"}, want: false},
		{name: "Invalid Semantic Version: 1.01.1", args: args{"1.01.1"}, want: false},
		{name: "Invalid Semantic Version: 1.1.01", args: args{"1.1.01"}, want: false},
		{name: "Invalid Semantic Version: 1.2", args: args{"1.2"}, want: false},
		{name: "Invalid Semantic Version: 1.2.3.DEV", args: args{"1.2.3.DEV"}, want: false},
		{name: "Invalid Semantic Version: 1.2-SNAPSHOT", args: args{"1.2-SNAPSHOT"}, want: false},
		{name: "Invalid Semantic Version: 1.2.31.2.3----RC-SNAPSHOT.12.09.1--..12+788", args: args{"1.2.31.2.3----RC-SNAPSHOT.12.09.1--..12+788"}, want: false},
		{name: "Invalid Semantic Version: 1.2-RC-SNAPSHOT", args: args{"1.2-RC-SNAPSHOT"}, want: false},
		{name: "Invalid Semantic Version: -1.0.3-gamma+b7718", args: args{"-1.0.3-gamma+b7718"}, want: false},
		{name: "Invalid Semantic Version: +justmeta", args: args{"+justmeta"}, want: false},
		{name: "Invalid Semantic Version: 9.8.7+meta+meta", args: args{"9.8.7+meta+meta"}, want: false},
		{name: "Invalid Semantic Version: 9.8.7-whatever+meta+meta", args: args{"9.8.7-whatever+meta+meta"}, want: false},
		{name: "Invalid Semantic Version: 99999999999999999999999.999999999999999999.99999999999999999----RC-SNAPSHOT.12.09.1--------------------------------..12", args: args{"99999999999999999999999.999999999999999999.99999999999999999----RC-SNAPSHOT.12.09.1--------------------------------..12"}, want: false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := IsSemver(tt.args.version); got != tt.want {
				t.Errorf("IsSemVer() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_FilterSemverTags(t *testing.T) {
	type args struct {
		tags []string
	}
	tests := []struct {
		svSvc          SemverSvc
		name           string
		args           args
		wantSemverTags []string
	}{
		{
			name:  "No SemVer compliant tags",
			svSvc: SemverSvc{client: nil},
			args: args{
				tags: []string{"v1.0", "v2.0.0", "release-1.1"},
			},
			wantSemverTags: []string{},
		},
		{
			name: "Single SemVer compliant tag",
			args: args{
				tags: []string{"v1.0.0", "v2.0.0", "release-1.1", "1.2.3-alpha"},
			},
			wantSemverTags: []string{"1.2.3-alpha"},
		},
		{
			name: "Multiple SemVer compliant tags",
			args: args{
				tags: []string{"v1.0.0", "v2.0.0", "release-1.1", "1.2.3-alpha", "2.0.0", "v3.0.0-beta.1", "4.5.6-rc.1+build.123"},
			},
			wantSemverTags: []string{"1.2.3-alpha", "2.0.0", "4.5.6-rc.1+build.123"},
		},
		{
			name: "Invalid SemVer compliant tags",
			args: args{
				tags: []string{"v1.0.0", "1.2", "v2.0.0", "release-1.1", "1.2.3-0123", "1.2.3-0123.0123", "+invalid", "-invalid", "-invalid+invalid", "-invalid.01", "alpha", "alpha_beta", "01.1.1", "1.01.1", "1.1.01", "1.2.3.DEV", "1.2-SNAPSHOT", "1.2-RC-SNAPSHOT", "-1.0.3-gamma+b7718", "+justmeta", "9.8.7+meta+meta", "9.8.7-whatever+meta+meta", "99999999999999999999999.999999999999999999.99999999999999999----RC-SNAPSHOT.12.09.1--------------------------------..12"},
			},
			wantSemverTags: []string{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotSemverTags, _ := tt.svSvc.FilterSemverTags(tt.args.tags, nil)
			if len(gotSemverTags) == 0 && len(tt.wantSemverTags) == 0 {
				return
			}
			assert.Equal(t, tt.wantSemverTags, gotSemverTags, "FilterSemverTags() = %v, want %v", gotSemverTags, tt.wantSemverTags)
		})
	}
}

func Test_IsRelease(t *testing.T) {
	type args struct {
		version string
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "valid release version",
			args: args{version: "1.2.3"},
			want: true,
		},
		{
			name: "invalid release version with too few segments",
			args: args{version: "1.2"},
			want: false,
		},
		{
			name: "invalid release version with too many segments",
			args: args{version: "1.2.3.4"},
			want: false,
		},
		{
			name: "invalid release version with non-numeric major version",
			args: args{version: "a.2.3"},
			want: false,
		},
		{
			name: "invalid release version with non-numeric minor version",
			args: args{version: "1.b.3"},
			want: false,
		},
		{
			name: "invalid release version with non-numeric patch version",
			args: args{version: "1.2.c"},
			want: false,
		},
		{
			name: "invalid release version with non-integer major version",
			args: args{version: "1.2.3.0"},
			want: false,
		},
		{
			name: "invalid release version with negative version number",
			args: args{version: "1.2.-3"},
			want: false,
		},
		{
			name: "invalid release version with leading zeros",
			args: args{version: "01.02.003"},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := IsRelease(tt.args.version); got != tt.want {
				t.Errorf("isRelease() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestFilterSemverTagsWithRelease(t *testing.T) {
	tests := []struct {
		name         string
		tags         []string
		wantFiltered []string
		wantErr      bool
	}{
		{
			name: "Filter release versions",
			tags: []string{
				"1.0.0-alpha",
				"1.0.0",
				"1.1.0",
				"2.0.0-beta",
				"2.0.0",
			},
			wantFiltered: []string{
				"1.0.0",
				"1.1.0",
				"2.0.0",
			},
			wantErr: false,
		},
		{
			name: "Filter release versions with invalid semver",
			tags: []string{
				"1.0.0-alpha",
				"1.0.0",
				"1.1.0",
				"2.0.0-beta",
				"2.0.0",
				"invalid",
			},
			wantFiltered: []string{
				"1.0.0",
				"1.1.0",
				"2.0.0",
			},
			wantErr: false,
		},
		{
			name: "Filter release versions with only prerelease versions",
			tags: []string{
				"1.0.0-alpha",
				"1.0.0-beta",
				"2.0.0-beta",
				"2.0.0-alpha",
			},
			wantFiltered: []string{},
			wantErr:      false,
		},
		{
			name:         "Empty input tags",
			tags:         []string{},
			wantFiltered: []string{},
			wantErr:      false,
		},
		{
			name: "Filter release versions with mixed semver and non-semver tags",
			tags: []string{
				"1.0.0",
				"1.1.0",
				"2.0.0",
				"v3.0.0",
				"3.0.1",
				"3.0.2",
				"non-semver",
			},
			wantFiltered: []string{
				"1.0.0",
				"1.1.0",
				"2.0.0",
				"3.0.1",
				"3.0.2",
			},
			wantErr: false,
		},
	}

	svSvc := &SemverSvc{}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			filteredTags, err := svSvc.FilterSemverTags(tt.tags, &Filters{Release: true})
			if tt.wantErr {
				assert.Error(t, err, "FilterSemverTags() error = %v, wantErr %v", err, tt.wantErr)
			} else {
				assert.NoError(t, err, "FilterSemverTags() error = %v, wantErr %v", err, tt.wantErr)
			}
			if len(filteredTags) == 0 && len(tt.wantFiltered) == 0 {
				return
			}

			assert.Equal(t, tt.wantFiltered, filteredTags, "FilterSemverTags() = %v, want %v", filteredTags, tt.wantFiltered)
		})
	}
}
