package ionic

import (
	"bufio"
	"fmt"
	"os"
	"regexp"
	"strings"

	"github.com/bitrise-io/go-utils/command"
	"github.com/bitrise-io/go-utils/log"
	ver "github.com/hashicorp/go-version"
	"github.com/pkg/errors"
)

// LogIntf declares used logger functions
type LogIntf interface {
	// Printf(format string, v ...interface{})
	Infof(format string, v ...interface{})
	Donef(format string, v ...interface{})
}

// Log is the currently set logger
var Log LogIntf = log.NewDummyLogger()

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
	out = strings.Split(out, "(")[0]
	out = strings.TrimSpace(out)

	version, err := ver.NewVersion(out)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to parse version: %s", out)
	}
	return version, nil
}

// LoginCommand runs ionic login comand
func LoginCommand(username string, password string) error {
	if username != "" && password != "" {
		fmt.Println()
		Log.Infof("Ionic login")

		cmdArgs := []string{"ionic", "login", username, password}
		cmd := command.New(cmdArgs[0], cmdArgs[1:]...)
		cmd.SetStdout(os.Stdout).SetStderr(os.Stderr).SetStdin(strings.NewReader("y"))

		Log.Donef("$ ionic login *** ***")

		if err := cmd.Run(); err != nil {
			return fmt.Errorf("command failed, error: %s", err)
		}
	}
	return nil
}

// PrepareCommand runs ionic cordova prepare command
func PrepareCommand(ionicMajorVersion int) error {
	cmdArgs := []string{"ionic"}
	if ionicMajorVersion > 2 {
		cmdArgs = append(cmdArgs, "cordova")
	}

	cmdArgs = append(cmdArgs, "prepare", "--no-build")
	cmd := command.New(cmdArgs[0], cmdArgs[1:]...)
	cmd.SetStdout(os.Stdout).SetStderr(os.Stderr)
	Log.Donef("$ %s", cmd.PrintableCommandArgs())

	if err := cmd.Run(); err != nil {
		return fmt.Errorf("command %s failed, error: %s", cmd.PrintableCommandArgs(), err)
	}
	return nil
}
