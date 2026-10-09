package cli

import (
	"fmt"
	"io"
	"os"
	"strings"
)

// Color modes for --color (git/ls convention).
const (
	colorAuto   = "auto"
	colorAlways = "always"
	colorNever  = "never"
)

func parseColorMode(s string) (string, error) {
	m := strings.ToLower(strings.TrimSpace(s))
	if m == "" {
		return colorAuto, nil
	}
	switch m {
	case colorAuto, colorAlways, colorNever:
		return m, nil
	default:
		return "", fmt.Errorf("--color must be auto|always|never (got %q)", s)
	}
}

// noColorSet reports whether NO_COLOR is set to a non-empty value (https://no-color.org/).
func noColorSet() bool {
	v, ok := os.LookupEnv("NO_COLOR")
	return ok && v != ""
}

// useColor decides whether decorative ANSI styling is enabled.
// --json always disables chrome. NO_COLOR (non-empty) wins over --color=always.
func useColor(mode string, w io.Writer, jsonMode bool) bool {
	if jsonMode {
		return false
	}
	if noColorSet() {
		return false
	}
	m, err := parseColorMode(mode)
	if err != nil {
		m = colorAuto
	}
	switch m {
	case colorNever:
		return false
	case colorAlways:
		return true
	default: // auto
		return isTerminal(w)
	}
}
