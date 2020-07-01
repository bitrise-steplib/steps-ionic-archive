package jsdependency

import (
	"testing"

	"github.com/bitrise-io/go-utils/command"
	"github.com/stretchr/testify/require"
)

func TestInstallGlobalDependencyCommand(t *testing.T) {
	tests := []struct {
		name           string
		packageManager Tool
		dependency     string
		version        string
		want           []InstallCommand
		wantErr        bool
	}{
		{
			name:           "Install latest ionic with yarn",
			packageManager: Yarn,
			dependency:     "ionic",
			version:        "latest",
			want: []InstallCommand{
				{
					Slice:       command.New("yarn", "remove", "ionic"),
					IgnoreError: true,
				},
				{
					Slice:       command.New("yarn", "global", "remove", "@ionic/cli"),
					IgnoreError: true,
				},
				{
					Slice:       command.New("yarn", "global", "add", "ionic@latest"),
					IgnoreError: false,
				},
			},
		},
		{
			name:           "Install latest @ionic/cli with yarn",
			packageManager: Yarn,
			dependency:     "@ionic/cli",
			version:        "latest",
			want: []InstallCommand{
				{
					Slice:       command.New("yarn", "remove", "@ionic/cli"),
					IgnoreError: true,
				},
				{
					Slice:       command.New("yarn", "global", "remove", "ionic"),
					IgnoreError: true,
				},
				{
					Slice:       command.New("yarn", "global", "add", "@ionic/cli@latest"),
					IgnoreError: false,
				},
			},
		},
		{
			name:           "Install latest corodva with yarn",
			packageManager: Yarn,
			dependency:     "cordova",
			version:        "latest",
			want: []InstallCommand{
				{
					Slice:       command.New("yarn", "remove", "cordova"),
					IgnoreError: true,
				},
				{
					Slice:       command.New("yarn", "global", "add", "cordova@latest"),
					IgnoreError: false,
				},
			},
		},
		{
			name:           "Install latest @ionic/cli with npm",
			packageManager: Npm,
			dependency:     "@ionic/cli",
			version:        "latest",
			want: []InstallCommand{
				{
					Slice:       command.New("npm", "remove", "@ionic/cli", "--force"),
					IgnoreError: false,
				},
				{
					Slice:       command.New("npm", "install", "-g", "@ionic/cli@latest", "--force"),
					IgnoreError: false,
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := InstallGlobalDependencyCommand(tt.packageManager, tt.dependency, tt.version)
			if (err != nil) != tt.wantErr {
				t.Errorf("InstallGlobalDependencyCommand() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			require.Equal(t, tt.want, got, "InstallGlobalDependencyCommand() return value")
		})
	}
}
