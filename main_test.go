package main

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func Test_parseMajorVersion(t *testing.T) {
	tests := []struct {
		name    string
		version string
		want    uint64
		wantErr bool
	}{
		{
			name:    "latest",
			version: "latest",
			wantErr: true,
		},
		{
			name:    "simple semver",
			version: "5.1.14",
			want:    5,
		},
		{
			name:    "semver with multi digit numbers",
			version: "11.1.14",
			want:    11,
		},
		{
			name:    "semver with v prefix",
			version: "v5.1.14",
			want:    5,
		},
		{
			name:    "major version only",
			version: "5",
			want:    5,
		},
		{
			name:    "major and minor version only",
			version: "5.6",
			want:    5,
		},
		{
			name:    "major and minor version only",
			version: "5.x",
			want:    5,
		},
		{
			name:    "equals operator",
			version: "=v5.1.14",
			want:    5,
		},
		{
			name:    "equals operator",
			version: "=v5.1.14",
			want:    5,
		},
		{
			name:    "more subversions",
			version: "5.1.14.6.3",
			want:    5,
		},
		{
			name:    "version suffix",
			version: "5.1.14-alpha.7",
			want:    5,
		},
		{
			name:    "zero prefix",
			version: "05.01.14",
			want:    5,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := parseMajorVersion(tt.version)
			if (err != nil) != tt.wantErr {
				t.Errorf("majorVersion() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			require.Equal(t, tt.want, got, "majorVersion() return value")
		})
	}
}
