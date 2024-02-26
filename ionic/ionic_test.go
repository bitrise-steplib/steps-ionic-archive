package ionic

import (
	"testing"

	ver "github.com/hashicorp/go-version"
	"github.com/stretchr/testify/require"
)

func Test_PackageNameFromVersion(t *testing.T) {
	const oldPackageName = "ionic"
	const newPackageName = "@ionic/cli"

	tests := []struct {
		name    string
		version string
		want    string
		wantErr bool
	}{
		{
			name:    "latest",
			version: "latest",
			want:    newPackageName,
		},
		{
			name:    "simple semver",
			version: "5.1.14",
			want:    oldPackageName,
		},
		{
			name:    "semver with multi digit numbers",
			version: "11.1.14",
			want:    newPackageName,
		},
		{
			name:    "semver with v prefix",
			version: "v5.1.14",
			want:    oldPackageName,
		},
		{
			name:    "major version only",
			version: "7",
			want:    newPackageName,
		},
		{
			name:    "major and minor version only",
			version: "5.6",
			want:    oldPackageName,
		},
		{
			name:    "major and minor version only",
			version: "5.x",
			want:    oldPackageName,
		},
		{
			name:    "equals operator",
			version: "=v5.1.14",
			want:    oldPackageName,
		},
		{
			name:    "equals operator",
			version: "=v5.1.14",
			want:    oldPackageName,
		},
		{
			name:    "more subversions",
			version: "6.1.14.6.3",
			want:    newPackageName,
		},
		{
			name:    "version suffix",
			version: "5.1.14-alpha.7",
			want:    oldPackageName,
		},
		{
			name:    "zero prefix",
			version: "05.01.14",
			want:    oldPackageName,
		},
		{
			name:    "invalid format",
			version: "beta",
			want:    newPackageName,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := PackageNameFromVersion(tt.version)
			if (err != nil) != tt.wantErr {
				t.Errorf("majorVersion() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			require.Equal(t, tt.want, got, "majorVersion() return value")
		})
	}
}

func TestFindIosTargetPathComponent(t *testing.T) {
    testCases := []struct {
        name           string
        target         string
        configuration  string
        cordovaVersion string
        want           string
    }{
        {
            name:           "Device + debug",
            target:         "device",
            configuration:  "debug",
            cordovaVersion: "7.1.0",
            want:           "Debug-iphoneos",
        },
        {
            name:           "Emulator + release",
            target:         "emulator",
            configuration:  "release",
            cordovaVersion: "99.99.99",
            want:           "Release-iphonesimulator",
        },
        {
            name:           "Device + release + old Cordova behavior",
            target:         "device",
            configuration:  "release",
            cordovaVersion: "6.5.0",
            want:           "device",
        },
        {
            name:           "Emulator + debug + old Cordova behavior",
            target:         "emulator",
            configuration:  "debug",
            cordovaVersion: "4.0.0",
            want:           "emulator",
        },
    }

    for _, tc := range testCases {
        t.Run(tc.name, func(t *testing.T) {
			cordovaVersion := ver.Must(ver.NewVersion(tc.cordovaVersion))
            got := FindIosTargetPathComponent(tc.target, tc.configuration, cordovaVersion)
            require.Equal(t, tc.want, got)
        })
    }
}
