package cli

import (
	_ "embed"
	"strings"
)

//go:embed version.txt
var versionFile string

// Version is the CLI package version (synced from sdk/cli/VERSION at generate time,
// or the embedded version.txt checked into the module).
var Version = strings.TrimSpace(versionFile)
