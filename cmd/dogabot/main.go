package main

import (
	"os"

	"github.com/dogabot/dogabot-cli/internal/cli"
)

func main() {
	if err := cli.Execute(); err != nil {
		os.Exit(1)
	}
}
