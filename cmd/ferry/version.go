package main

import (
	"fmt"
	"runtime/debug"
)

// Set at build time via -ldflags "-X main.version=... -X main.commit=... -X main.date=..."
// (see Makefile and .goreleaser.yaml).
var (
	version = "dev"
	commit  = "none"
	date    = "unknown"
)

func printVersion() {
	v := version
	// Binaries built with `go install ...@vX.Y.Z` have no ldflags,
	// but the module version is embedded in the build info.
	if v == "dev" {
		if bi, ok := debug.ReadBuildInfo(); ok && bi.Main.Version != "" && bi.Main.Version != "(devel)" {
			v = bi.Main.Version
		}
	}
	fmt.Printf("ferry %s (commit %s, built %s)\n", v, commit, date)
}
