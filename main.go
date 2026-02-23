package main

import (
	"github.com/vineethkrishnan/dockit/cmd"
)

var (
	version = "dev"
	commit  = "none"
	date    = "unknown"
)

func main() {
	// Tell the command package about the version variables injected by GoReleaser
	cmd.SetVersion(version, commit, date)
	cmd.Execute()
}
