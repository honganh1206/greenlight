package vcs

import (
	"fmt"
	"runtime/debug"
)

func Version() string {
	var (
		revision string
		modified bool
		time     string
	)

	// Same like when we run go version -m
	bi, ok := debug.ReadBuildInfo()
	if ok {
		for _, s := range bi.Settings {
			switch s.Key {
			// Commit time
			case "vcs.time":
				time = s.Value
			// Hash for the latest Git commit
			case "vcs.revision":
				revision = s.Value
			// Whether the code tracked by Git has been modified since the commit was made
			case "vcs.modified":
				if s.Value == "true" {
					modified = true
				}
			}
		}
	}

	if modified {
		return fmt.Sprintf("%s-%s-dirty", time, revision)
	}

	return fmt.Sprintf("%s-%s", time, revision)
}
