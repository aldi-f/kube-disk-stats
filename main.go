package main

import (
	"fmt"
	"os"

	"github.com/aldi-f/kube-disk-stats/cmd"
)

var Version = "1.2.0"

func main() {
	cmd.Version = Version
	rootCmd := cmd.NewRootCmd()
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}
