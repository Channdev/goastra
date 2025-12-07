/*
 * GoAstra CLI - Start Command
 *
 * Runs the production build server.
 * Executes the compiled binary from the dist directory.
 */
package cmd

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"

	"github.com/fatih/color"
	"github.com/spf13/cobra"
)

var (
	startPort   int
	startDir    string
)

var startCmd = &cobra.Command{
	Use:   "start",
	Short: "Start production server",
	Long: `Starts the GoAstra production server:
  - Runs the compiled Go backend binary
  - Serves the Angular frontend (API + static files in one server)
  - Uses .env for environment configuration

Run 'goastra build' first to create the production build.`,
	RunE: runStart,
}

func init() {
	rootCmd.AddCommand(startCmd)
	startCmd.Flags().IntVarP(&startPort, "port", "p", 8080, "server port")
	startCmd.Flags().StringVarP(&startDir, "dir", "d", "dist", "build directory")
}

func runStart(cmd *cobra.Command, args []string) error {
	projectRoot, err := findStartProjectRoot()
	if err != nil {
		return fmt.Errorf("not in a GoAstra project: %w", err)
	}

	distPath := filepath.Join(projectRoot, startDir)

	if _, err := os.Stat(distPath); os.IsNotExist(err) {
		return fmt.Errorf("build directory not found: %s\nRun 'goastra build' first", distPath)
	}

	binaryName := "server"
	if runtime.GOOS == "windows" {
		binaryName = "server.exe"
	}

	binaryPath := filepath.Join(distPath, binaryName)

	if _, err := os.Stat(binaryPath); os.IsNotExist(err) {
		return fmt.Errorf("server binary not found: %s\nRun 'goastra build' first", binaryPath)
	}

	color.Cyan("Starting GoAstra production server...\n")
	color.Green("Server running at http://localhost:%d\n\n", startPort)

	serverCmd := exec.Command(binaryPath)
	serverCmd.Dir = distPath
	serverCmd.Env = append(os.Environ(),
		fmt.Sprintf("PORT=%d", startPort),
		"APP_ENV=production",
	)
	serverCmd.Stdout = os.Stdout
	serverCmd.Stderr = os.Stderr

	return serverCmd.Run()
}

func findStartProjectRoot() (string, error) {
	dir, err := os.Getwd()
	if err != nil {
		return "", err
	}
	for {
		configPath := filepath.Join(dir, "goastra.json")
		if _, err := os.Stat(configPath); err == nil {
			return dir, nil
		}
		parent := filepath.Dir(dir)
		if parent == dir {
			break
		}
		dir = parent
	}
	return "", fmt.Errorf("goastra.json not found")
}
