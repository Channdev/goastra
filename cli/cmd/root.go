/*
 * GoAstra CLI - Root Command
 *
 * Defines the root command for the GoAstra CLI framework.
 * All subcommands attach to this root and inherit global configuration.
 */
package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var (
	cfgFile string
	verbose bool
)

/*
 * rootCmd represents the base command when called without any subcommands.
 * It provides framework information and usage instructions.
 */
var rootCmd = &cobra.Command{
	Use:   "goastra",
	Short: "GoAstra - Full-stack Go + Angular framework",
	Long: `GoAstra is a production-grade full-stack framework that combines
Go backend with Angular frontend. It provides code generation,
environment management, and unified development tooling.

Usage:
  goastra new <project>     Create a new GoAstra project
  goastra dev               Start development servers
  goastra build             Build for production
  goastra generate          Generate code artifacts
  goastra typesync          Sync Go types to TypeScript
  goastra test              Run test suites`,
	Version: "1.0.0",
}

/*
 * Execute initializes and runs the root command.
 * This is called by main.main() and handles all command routing.
 */
func Execute() error {
	return rootCmd.Execute()
}

func init() {
	cobra.OnInitialize(initConfig)

	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is ./goastra.json)")
	rootCmd.PersistentFlags().BoolVarP(&verbose, "verbose", "v", false, "enable verbose output")
}

/*
 * initConfig reads in config file and ENV variables if set.
 * Called before any command execution to ensure proper configuration.
 */
func initConfig() {
	if cfgFile != "" {
		if _, err := os.Stat(cfgFile); os.IsNotExist(err) {
			fmt.Fprintf(os.Stderr, "Config file not found: %s\n", cfgFile)
		}
	}
}
