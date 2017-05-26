package ionic

import "github.com/bitrise-io/go-utils/command"

// Model ...
type Model struct {
	platforms     []string
	configuration string
	target        string
	buildConfig   string
	customOptions []string
}

// New ...
func New() Model {
	return Model{}
}

// SetPlatforms ...
func (builder *Model) SetPlatforms(platforms ...string) *Model {
	builder.platforms = platforms
	return builder
}

// SetConfiguration ...
func (builder *Model) SetConfiguration(configuration string) *Model {
	builder.configuration = configuration
	return builder
}

// SetTarget ...
func (builder *Model) SetTarget(target string) *Model {
	builder.target = target
	return builder
}

// SetBuildConfig ...
func (builder *Model) SetBuildConfig(buildConfig string) *Model {
	builder.buildConfig = buildConfig
	return builder
}

// SetCustomOptions ...
func (builder *Model) SetCustomOptions(customOptions ...string) *Model {
	builder.customOptions = customOptions
	return builder
}

func (builder *Model) defaultIonicCommandSlice(cmd ...string) []string {
	return []string{"ionic", "--no-interactive", "--confirm"}
}

func (builder *Model) defaultIonicCordovaCommandSlice(cmd ...string) []string {
	return []string{"ionic", "cordova", "--no-interactive", "--confirm"}
}

func (builder *Model) cordovaCommandSlice(cmd ...string) []string {
	cmdSlice := builder.defaultIonicCordovaCommandSlice()
	cmdSlice = append(cmdSlice, cmd...)

	if len(cmd) == 1 && cmd[0] == "build" {
		if builder.configuration != "" {
			cmdSlice = append(cmdSlice, "--"+builder.configuration)
		}
		if builder.target != "" {
			cmdSlice = append(cmdSlice, "--"+builder.target)
		}
	}

	if len(builder.platforms) > 0 {
		cmdSlice = append(cmdSlice, builder.platforms...)
	}

	if len(cmd) == 1 && cmd[0] == "build" {
		if builder.buildConfig != "" {
			cmdSlice = append(cmdSlice, "--buildConfig", builder.buildConfig)
		}
	}

	return cmdSlice
}

// VersionCommand ...
func (builder *Model) VersionCommand() *command.Model {
	cmdSlice := builder.defaultIonicCommandSlice()
	cmdSlice = append(cmdSlice, "-v")
	return command.New(cmdSlice[0], cmdSlice[1:]...)
}

// PlatformCommand ...
func (builder *Model) PlatformCommand(cmd string) *command.Model {
	cmdSlice := builder.cordovaCommandSlice("platform", cmd)
	return command.New(cmdSlice[0], cmdSlice[1:]...)
}

// BuildCommand ...
func (builder *Model) BuildCommand() *command.Model {
	cmdSlice := builder.cordovaCommandSlice("build")
	return command.New(cmdSlice[0], cmdSlice[1:]...)
}
