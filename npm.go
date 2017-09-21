package main

import (
	"fmt"

	"github.com/bitrise-io/go-utils/command"
	"github.com/bitrise-io/go-utils/log"
)

func npmRemove(isGlobal bool, pkg ...string) error {
	args := []string{"remove"}
	if isGlobal {
		args = append(args, "-g")
	}
	args = append(args, pkg...)
	cmd := command.New("npm", args...)

	log.Donef("$ %s", cmd.PrintableCommandArgs())

	if out, err := cmd.RunAndReturnTrimmedCombinedOutput(); err != nil {
		return fmt.Errorf("command failed, output: %s, error: %s", out, err)
	}
	return nil
}

func npmInstall(isGlobal bool, pkg ...string) error {
	args := []string{"install"}
	if isGlobal {
		args = append(args, "-g")
	}
	args = append(args, pkg...)
	cmd := command.New("npm", args...)

	log.Donef("$ %s", cmd.PrintableCommandArgs())

	if out, err := cmd.RunAndReturnTrimmedCombinedOutput(); err != nil {
		return fmt.Errorf("command failed, output: %s, error: %s", out, err)
	}
	return nil
}
