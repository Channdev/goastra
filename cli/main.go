/*
 * GoAstra CLI - Main Entry Point
 *
 * This file serves as the bootstrap for the GoAstra command-line interface.
 * It initializes the root command and delegates execution to the cmd package.
 */
package main

import (
	"os"

	"github.com/goastra/cli/cmd"
)

func main() {
	if err := cmd.Execute(); err != nil {
		os.Exit(1)
	}
}
