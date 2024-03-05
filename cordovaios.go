package main

import (
	"github.com/bitrise-io/go-utils/log"
	"github.com/bitrise-io/go-utils/pathutil"
)

func getIOSOutoutDirPath(candidateDirPaths []string) string {
	for _, path := range candidateDirPaths {
		exist, err := pathutil.IsDirExists(path)
		if err != nil {
			log.Warnf("Failed to check if dir (%s) exist: %s", path, err)
			continue
		}

		if exist {
			return path
		}
	}

	return ""
}
