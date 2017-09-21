package main

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/bitrise-io/go-utils/colorstring"
	"github.com/bitrise-io/go-utils/command"
	"github.com/bitrise-io/go-utils/log"
	"github.com/bitrise-io/go-utils/pathutil"
	"github.com/bitrise-io/go-utils/ziputil"
	"github.com/bitrise-tools/go-steputils/input"
	"github.com/bitrise-tools/go-steputils/tools"
	ver "github.com/hashicorp/go-version"
	shellquote "github.com/kballard/go-shellquote"
	"github.com/pkg/errors"
)

const (
	ipaPathEnvKey = "BITRISE_IPA_PATH"

	appZipPathEnvKey = "BITRISE_APP_PATH"
	appDirPathEnvKey = "BITRISE_APP_DIR_PATH"

	dsymDirPathEnvKey = "BITRISE_DSYM_DIR_PATH"
	dsymZipPathEnvKey = "BITRISE_DSYM_PATH"

	apkPathEnvKey = "BITRISE_APK_PATH"
)

// ConfigsModel ...
type ConfigsModel struct {
	Platform      string
	Configuration string
	Target        string
	BuildConfig   string
	Options       string

	Username string
	Password string

	CordovaVersion string
	IonicVersion   string

	WorkDir   string
	DeployDir string
}

func createConfigsModelFromEnvs() ConfigsModel {
	return ConfigsModel{
		Platform:      os.Getenv("platform"),
		Configuration: os.Getenv("configuration"),
		Target:        os.Getenv("target"),
		BuildConfig:   os.Getenv("build_config"),
		Options:       os.Getenv("options"),

		Username: os.Getenv("ionic_username"),
		Password: os.Getenv("ionic_password"),

		CordovaVersion: os.Getenv("cordova_version"),
		IonicVersion:   os.Getenv("ionic_version"),

		WorkDir:   os.Getenv("workdir"),
		DeployDir: os.Getenv("BITRISE_DEPLOY_DIR"),
	}
}

func (configs ConfigsModel) print() {
	log.Infof("Configs:")
	log.Printf("- Platform: %s", configs.Platform)
	log.Printf("- Configuration: %s", configs.Configuration)
	log.Printf("- Target: %s", configs.Target)
	log.Printf("- BuildConfig: %s", configs.BuildConfig)
	log.Printf("- Options: %s", configs.Options)

	log.Printf("- Username: %s", input.SecureInput(configs.Username))
	log.Printf("- Username: %s", input.SecureInput(configs.Password))

	log.Printf("- CordovaVersion: %s", configs.CordovaVersion)
	log.Printf("- IonicVersion: %s", configs.IonicVersion)

	log.Printf("- WorkDir: %s", configs.WorkDir)
	log.Printf("- DeployDir: %s", configs.DeployDir)
}

func (configs ConfigsModel) validate() error {
	if err := input.ValidateIfDirExists(configs.WorkDir); err != nil {
		return fmt.Errorf("WorkDir: %s", err)
	}

	if err := input.ValidateWithOptions(configs.Platform, "ios,android", "ios", "android"); err != nil {
		return fmt.Errorf("Platform: %s", err)
	}

	if err := input.ValidateIfNotEmpty(configs.Configuration); err != nil {
		return fmt.Errorf("Configuration: %s", err)
	}

	if err := input.ValidateIfNotEmpty(configs.Target); err != nil {
		return fmt.Errorf("Target: %s", err)
	}

	return nil
}

func moveAndExportOutputs(outputs []string, deployDir, envKey string) (string, error) {
	outputToExport := ""
	for _, output := range outputs {
		info, err := os.Lstat(output)
		if err != nil {
			return "", err
		}

		if info.Mode()&os.ModeSymlink != 0 {
			resolvedPth, err := os.Readlink(output)
			if err != nil {
				return "", err
			}

			log.Warnf("Output: %s is a symlink to: %s", output, resolvedPth)

			if exist, err := pathutil.IsPathExists(resolvedPth); err != nil {
				return "", err
			} else if !exist {
				return "", fmt.Errorf("resolved path: %s does not exist", resolvedPth)
			}

			resolvedInfo, err := os.Lstat(resolvedPth)
			if err != nil {
				return "", err
			}

			if resolvedInfo.Mode()&os.ModeSymlink != 0 {
				return "", fmt.Errorf("resolved path: %s is still symlink", resolvedPth)
			}

			output = resolvedPth
			info = resolvedInfo
		}

		fileName := filepath.Base(output)
		destinationPth := filepath.Join(deployDir, fileName)

		if info.IsDir() {
			if err := command.CopyDir(output, destinationPth, false); err != nil {
				return "", err
			}
		} else {
			if err := command.CopyFile(output, destinationPth); err != nil {
				return "", err
			}
		}

		outputToExport = destinationPth
	}

	if outputToExport == "" {
		return "", nil
	}

	if err := tools.ExportEnvironmentWithEnvman(envKey, outputToExport); err != nil {
		return "", err
	}

	return outputToExport, nil
}

func ionicVersion() (*ver.Version, error) {
	cmd := command.New("ionic", "-v")
	cmd.SetStdin(strings.NewReader("Y"))
	out, err := cmd.RunAndReturnTrimmedCombinedOutput()
	if err != nil {
		return nil, err
	}

	// fix for ionic-cli intercative version output: `[1000D[K3.2.0`
	pattern := `(?P<version>\d+\.\d+\.\d+)`

	reader := strings.NewReader(out)
	scanner := bufio.NewScanner(reader)
	for scanner.Scan() {
		line := scanner.Text()
		if match := regexp.MustCompile(pattern).FindStringSubmatch(line); len(match) == 2 {
			versionStr := match[1]
			version, err := ver.NewVersion(versionStr)
			if err != nil {
				return nil, errors.Wrapf(err, "failed to parse version: %s", versionStr)
			}
			return version, nil
		}
	}
	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return nil, fmt.Errorf("output: %s", out)
}

func cordovaVersion() (*ver.Version, error) {
	cmd := command.New("cordova", "-v")
	out, err := cmd.RunAndReturnTrimmedCombinedOutput()
	if err != nil {
		return nil, err
	}
	version, err := ver.NewVersion(out)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to parse version: %s", out)
	}
	return version, nil
}

func fail(format string, v ...interface{}) {
	log.Errorf(format, v...)
	os.Exit(1)
}

func main() {
	configs := createConfigsModelFromEnvs()

	fmt.Println()
	configs.print()

	if err := configs.validate(); err != nil {
		fail("Issue with input: %s", err)
	}

	// Change dir to working directory
	workDir, err := pathutil.AbsPath(configs.WorkDir)
	if err != nil {
		fail("Failed to expand WorkDir (%s), error: %s", configs.WorkDir, err)
	}

	currentDir, err := pathutil.CurrentWorkingDirectoryAbsolutePath()
	if err != nil {
		fail("Failed to get current directory, error: %s", err)
	}

	if workDir != currentDir {
		fmt.Println()
		log.Infof("Switch working directory to: %s", workDir)

		revokeFunc, err := pathutil.RevokableChangeDir(workDir)
		if err != nil {
			fail("Failed to change working directory, error: %s", err)
		}
		defer func() {
			fmt.Println()
			log.Infof("Reset working directory")
			if err := revokeFunc(); err != nil {
				fail("Failed to reset working directory, error: %s", err)
			}
		}()
	}

	// Update cordova and ionic version
	if configs.CordovaVersion != "" {
		fmt.Println()
		log.Infof("Updating cordova version to: %s", configs.CordovaVersion)

		if err := npmRemove(false, "cordova"); err != nil {
			fail("Failed to remove cordova, error: %s", err)
		}

		if err := npmInstall(true, "cordova@"+configs.CordovaVersion); err != nil {
			fail("Failed to install cordova, error: %s", err)
		}
	}

	if configs.IonicVersion != "" {
		fmt.Println()
		log.Infof("Updating ionic version to: %s", configs.IonicVersion)

		if err := npmRemove(false, "ionic"); err != nil {
			fail("Failed to remove ionic, error: %s", err)
		}

		if err := npmInstall(true, "ionic@"+configs.IonicVersion); err != nil {
			fail("Failed to install ionic, error: %s", err)
		}
	}

	fmt.Println()
	log.Infof("Installing cordova and angular plugins")
	if err := npmInstall(false, "@ionic/cli-plugin-ionic-angular@latest", "@ionic/cli-plugin-cordova@latest"); err != nil {
		fail("command failed, error: %s", err)
	}

	// Print cordova and ionic version
	cordovaVer, err := cordovaVersion()
	if err != nil {
		fail("Failed to get cordova version, error: %s", err)
	}

	fmt.Println()
	log.Printf("cordova version: %s", colorstring.Green(cordovaVer.String()))

	ionicVer, err := ionicVersion()
	if err != nil {
		fail("Failed to get ionic version, error: %s", err)
	}

	log.Printf("ionic version: %s", colorstring.Green(ionicVer.String()))

	//
	// ionic login
	if configs.Username != "" && configs.Password != "" {
		fmt.Println()
		log.Infof("Ionic login")

		cmdArgs := []string{"ionic", "login", configs.Username, configs.Password}
		cmd := command.New(cmdArgs[0], cmdArgs[1:]...)
		cmd.SetStdout(os.Stdout).SetStderr(os.Stderr).SetStdin(strings.NewReader("y"))

		log.Donef("$ ionic login *** ***")

		if err := cmd.Run(); err != nil {
			fail("command failed, error: %s", err)
		}
	}

	ionicMajorVersion := ionicVer.Segments()[0]

	platforms := []string{}
	if configs.Platform != "" {
		platformsSplit := strings.Split(configs.Platform, ",")
		for _, platform := range platformsSplit {
			platforms = append(platforms, strings.TrimSpace(platform))
		}
	}

	// ionic prepare
	fmt.Println()
	log.Infof("Building project")

	{
		// platform rm
		for _, platform := range platforms {
			cmdArgs := []string{"ionic"}
			if ionicMajorVersion > 2 {
				cmdArgs = append(cmdArgs, "cordova")
			}

			cmdArgs = append(cmdArgs, "platform", "rm")

			cmdArgs = append(cmdArgs, platform)

			cmd := command.New(cmdArgs[0], cmdArgs[1:]...)
			cmd.SetStdout(os.Stdout).SetStderr(os.Stderr).SetStdin(strings.NewReader("y"))

			log.Donef("$ %s", cmd.PrintableCommandArgs())

			if err := cmd.Run(); err != nil {
				fail("command failed, error: %s", err)
			}
		}
	}

	{
		// platform add
		for _, platform := range platforms {
			cmdArgs := []string{"ionic"}
			if ionicMajorVersion > 2 {
				cmdArgs = append(cmdArgs, "cordova")
			}

			cmdArgs = append(cmdArgs, "platform", "add")

			cmdArgs = append(cmdArgs, platform)

			cmd := command.New(cmdArgs[0], cmdArgs[1:]...)
			cmd.SetStdout(os.Stdout).SetStderr(os.Stderr).SetStdin(strings.NewReader("y"))

			log.Donef("$ %s", cmd.PrintableCommandArgs())

			if err := cmd.Run(); err != nil {
				fail("command failed, error: %s", err)
			}
		}
	}

	{
		// build
		options := []string{}
		if configs.Options != "" {
			opts, err := shellquote.Split(configs.Options)
			if err != nil {
				fail("Failed to shell split Options (%s), error: %s", configs.Options, err)
			}
			options = opts
		}

		for _, platform := range platforms {
			cmdArgs := []string{"ionic"}
			if ionicMajorVersion > 2 {
				cmdArgs = append(cmdArgs, "cordova")
			}

			cmdArgs = append(cmdArgs, "build")

			if configs.Configuration != "" {
				cmdArgs = append(cmdArgs, "--"+configs.Configuration)
			}

			if configs.Target != "" {
				cmdArgs = append(cmdArgs, "--"+configs.Target)
			}

			cmdArgs = append(cmdArgs, platform)

			if configs.BuildConfig != "" {
				cmdArgs = append(cmdArgs, "--buildConfig", configs.BuildConfig)
			}

			cmdArgs = append(cmdArgs, options...)
			cmdArgs = append(cmdArgs)

			cmd := command.New(cmdArgs[0], cmdArgs[1:]...)
			cmd.SetStdout(os.Stdout).SetStderr(os.Stderr).SetStdin(strings.NewReader("y"))

			log.Donef("$ %s", cmd.PrintableCommandArgs())

			if err := cmd.Run(); err != nil {
				fail("command failed, error: %s", err)
			}
		}
	}

	// collect outputs

	iosOutputDirExist := false
	iosOutputDir := filepath.Join(workDir, "platforms", "ios", "build", configs.Target)
	if exist, err := pathutil.IsDirExists(iosOutputDir); err != nil {
		fail("Failed to check if dir (%s) exist, error: %s", iosOutputDir, err)
	} else if exist {
		iosOutputDirExist = true

		fmt.Println()
		log.Infof("Collecting ios outputs")

		// ipa
		ipaPattern := filepath.Join(iosOutputDir, "*.ipa")
		ipas, err := filepath.Glob(ipaPattern)
		if err != nil {
			fail("Failed to find ipas, with pattern (%s), error: %s", ipaPattern, err)
		}

		if len(ipas) > 0 {
			if exportedPth, err := moveAndExportOutputs(ipas, configs.DeployDir, ipaPathEnvKey); err != nil {
				fail("Failed to export ipas, error: %s", err)
			} else if exportedPth != "" {
				log.Donef("The ipa path is now available in the Environment Variable: %s (value: %s)", ipaPathEnvKey, exportedPth)
			}
		}
		// ---

		// dsym
		dsymPattern := filepath.Join(iosOutputDir, "*.dSYM")
		dsyms, err := filepath.Glob(dsymPattern)
		if err != nil {
			fail("Failed to find dSYMs, with pattern (%s), error: %s", dsymPattern, err)
		}

		if len(dsyms) > 0 {
			if exportedPth, err := moveAndExportOutputs(dsyms, configs.DeployDir, dsymDirPathEnvKey); err != nil {
				fail("Failed to export dsyms, error: %s", err)
			} else if exportedPth != "" {
				log.Donef("The dsym dir path is now available in the Environment Variable: %s (value: %s)", dsymDirPathEnvKey, exportedPth)

				zippedExportedPth := exportedPth + ".zip"
				if err := ziputil.ZipDir(exportedPth, zippedExportedPth, false); err != nil {
					fail("Failed to zip dsym dir (%s), error: %s", exportedPth, err)
				}

				if err := tools.ExportEnvironmentWithEnvman(dsymZipPathEnvKey, zippedExportedPth); err != nil {
					fail("Failed to export dsym.zip (%s), error: %s", zippedExportedPth, err)
				}

				log.Donef("The dsym.zip path is now available in the Environment Variable: %s (value: %s)", dsymZipPathEnvKey, zippedExportedPth)
			}
		}
		// --

		// app
		appPattern := filepath.Join(iosOutputDir, "*.app")
		apps, err := filepath.Glob(appPattern)
		if err != nil {
			fail("Failed to find apps, with pattern (%s), error: %s", appPattern, err)
		}

		if len(apps) > 0 {
			if exportedPth, err := moveAndExportOutputs(apps, configs.DeployDir, appDirPathEnvKey); err != nil {
				log.Warnf("Failed to export apps, error: %s", err)
			} else if exportedPth != "" {
				log.Donef("The app dir path is now available in the Environment Variable: %s (value: %s)", appDirPathEnvKey, exportedPth)

				zippedExportedPth := exportedPth + ".zip"
				if err := ziputil.ZipDir(exportedPth, zippedExportedPth, false); err != nil {
					fail("Failed to zip app dir (%s), error: %s", exportedPth, err)
				}

				if err := tools.ExportEnvironmentWithEnvman(appZipPathEnvKey, zippedExportedPth); err != nil {
					fail("Failed to export app.zip (%s), error: %s", zippedExportedPth, err)
				}

				log.Donef("The app.zip path is now available in the Environment Variable: %s (value: %s)", appZipPathEnvKey, zippedExportedPth)
			}
		}
		// ---
	}

	androidOutputDirExist := false
	androidOutputDir := filepath.Join(workDir, "platforms", "android", "build", "outputs", "apk")
	if exist, err := pathutil.IsDirExists(androidOutputDir); err != nil {
		fail("Failed to check if dir (%s) exist, error: %s", androidOutputDir, err)
	} else if exist {
		androidOutputDirExist = true

		fmt.Println()
		log.Infof("Collecting android outputs")

		pattern := filepath.Join(androidOutputDir, "*.apk")
		apks, err := filepath.Glob(pattern)
		if err != nil {
			fail("Failed to find apks, with pattern (%s), error: %s", pattern, err)
		}

		if len(apks) > 0 {
			if exportedPth, err := moveAndExportOutputs(apks, configs.DeployDir, apkPathEnvKey); err != nil {
				fail("Failed to export apks, error: %s", err)
			} else if exportedPth != "" {
				log.Donef("The apk path is now available in the Environment Variable: %s (value: %s)", apkPathEnvKey, exportedPth)
			}
		}
	}

	if !iosOutputDirExist && !androidOutputDirExist {
		log.Warnf("No ios nor android platform's output dir exist")
		fail("No output generated")
	}
}
