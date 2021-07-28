package ionic

import (
	"bufio"
	"fmt"
	"regexp"
	"strconv"
	"strings"

	"github.com/bitrise-io/go-utils/command"
	ver "github.com/hashicorp/go-version"
	"github.com/pkg/errors"
)

// Version returns ionic version
func Version() (*ver.Version, error) {
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

// CordovaVersion returns cordova version
func CordovaVersion() (*ver.Version, error) {
	cmd := command.New("cordova", "-v")
	out, err := cmd.RunAndReturnTrimmedCombinedOutput()
	if err != nil {
		return nil, err
	}
	re := regexp.MustCompile("(?:.*[^0-9])?([0-9]+\\.[0-9]+\\.[0-9]+)(?:[^0-9].*)?")
	matches := re.FindStringSubmatch(out)
	if len(matches) < 2 {
		return nil, fmt.Errorf("failed to parse version: %s", out)
	}
	out = matches[1]

	version, err := ver.NewVersion(out)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to parse version: %s", out)
	}
	return version, nil
}

// LoginCommand returns ionic login comand model
func LoginCommand(username string, password string) *command.Model {
	cmdArgs := []string{"ionic", "login", username, password}
	return command.New(cmdArgs[0], cmdArgs[1:]...)
}

// PrepareCommand returns ionic cordova prepare command model
func PrepareCommand(ionicMajorVersion int) *command.Model {
	cmdArgs := []string{"ionic"}
	if ionicMajorVersion > 2 {
		cmdArgs = append(cmdArgs, "cordova")
	}
	cmdArgs = append(cmdArgs, "prepare", "--no-build")
	return command.New(cmdArgs[0], cmdArgs[1:]...)
}

// PackageNameFromVersion returns either "ionic" or "@ionic/cli" based on the required version
// "ionic" is deprecated: https://www.npmjs.com/package/ionic
// "@ionic/cli" starts from version 6.0.0: https://www.npmjs.com/package/@ionic/cli
func PackageNameFromVersion(version string) (string, error) {
	newPackageName := "@ionic/cli"
	if version == "latest" {
		return newPackageName, nil
	}

	majorVersion, err := parseMajorVersion(version)
	if err != nil {
		return newPackageName, err
	}

	if majorVersion < 6 {
		return "ionic", nil
	}

	return newPackageName, nil
}

func parseMajorVersion(version string) (uint64, error) {
	matcher := regexp.MustCompile(`(\d+)(.\d+)*`)
	matches := matcher.FindStringSubmatch(version)
	if matches == nil {
		return 0, fmt.Errorf("failed to parse ionic major version (%s): no match", version)
	}

	if len(matches) < 2 {
		return 0, fmt.Errorf("failed to parse ionic major version (%s): unexpected match", version)
	}

	majorVersion, err := strconv.ParseUint(matches[1], 10, 64)
	if err != nil {
		return 0, fmt.Errorf("failed to parse ionic major version (%s): %s", version, err)
	}

	return majorVersion, nil
}
