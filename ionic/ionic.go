package ionic

import "github.com/bitrise-io/go-utils/command"

// Model ...
type Model struct {
	ionicMajorVersion int
	platforms         []string
	configuration     string
	target            string
	buildConfig       string
	customOptions     []string
}

// New ...
func New(ionicMajorVersion int) Model {
	return Model{
		ionicMajorVersion: ionicMajorVersion,
	}
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

// PlatformCommand ...
func (builder *Model) PlatformCommand(cmd string) *command.Model {
	cmdSlice := []string{}

	if builder.ionicMajorVersion > 2 {
		cmdSlice = []string{"ionic", "cordova"}
	} else {
		cmdSlice = []string{"ionic"}
	}

	cmdSlice = append(cmdSlice, "platform", cmd)
	cmdSlice = append(cmdSlice, builder.platforms...)
	return command.New(cmdSlice[0], cmdSlice[1:]...)
}

// BuildCommand ...
func (builder *Model) BuildCommand() *command.Model {
	cmdSlice := []string{}

	if builder.ionicMajorVersion > 2 {
		cmdSlice = []string{"ionic", "cordova"}
	} else {
		cmdSlice = []string{"ionic"}
	}

	cmdSlice = append(cmdSlice, "build")

	if builder.configuration != "" {
		cmdSlice = append(cmdSlice, "--"+builder.configuration)
	}
	if builder.target != "" {
		cmdSlice = append(cmdSlice, "--"+builder.target)
	}

	cmdSlice = append(cmdSlice, builder.platforms...)

	if builder.buildConfig != "" {
		cmdSlice = append(cmdSlice, "--buildConfig", builder.buildConfig)
	}

	return command.New(cmdSlice[0], cmdSlice[1:]...)
}
