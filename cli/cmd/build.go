/*
 * GoAstra CLI - Build Command
 *
 * Handles production builds for both Go backend and Angular frontend.
 * Loads .env.production, optimizes assets, and generates deployment artifacts.
 */
package cmd

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"

	"github.com/fatih/color"
	"github.com/channdev/goastra/cli/internal/env"
	"github.com/spf13/cobra"
)

var (
	buildOutput   string
	buildPlatform string
	buildArch     string
)

/*
 * buildCmd compiles the project for production deployment.
 * Creates optimized binaries and static assets.
 */
var buildCmd = &cobra.Command{
	Use:   "build",
	Short: "Build for production",
	Long: `Builds the GoAstra project for production:
  - Compiles Go backend to optimized binary
  - Builds Angular app with AOT compilation
  - Loads .env.production environment
  - Outputs to dist/ directory`,
	RunE: runBuild,
}

func init() {
	rootCmd.AddCommand(buildCmd)

	buildCmd.Flags().StringVarP(&buildOutput, "output", "o", "dist", "output directory")
	buildCmd.Flags().StringVar(&buildPlatform, "platform", runtime.GOOS, "target platform (linux, windows, darwin)")
	buildCmd.Flags().StringVar(&buildArch, "arch", runtime.GOARCH, "target architecture (amd64, arm64)")
}

/*
 * runBuild executes the production build pipeline.
 * Coordinates backend and frontend builds with proper ordering.
 */
func runBuild(cmd *cobra.Command, args []string) error {
	projectRoot, err := findBuildProjectRoot()
	if err != nil {
		return fmt.Errorf("not in a GoAstra project: %w", err)
	}

	if err := env.Load("production"); err != nil {
		color.Yellow("Warning: Could not load production environment: %v\n", err)
	}

	color.Cyan("Building GoAstra project for production...\n")

	outputPath, err := filepath.Abs(buildOutput)
	if err != nil {
		return fmt.Errorf("invalid output path: %w", err)
	}

	if err := os.MkdirAll(outputPath, 0755); err != nil {
		return fmt.Errorf("failed to create output directory: %w", err)
	}

	color.Yellow("[1/3] Building Go backend...\n")
	if err := buildBackend(projectRoot, outputPath); err != nil {
		return fmt.Errorf("backend build failed: %w", err)
	}
	color.Green("[1/3] Backend build complete.\n")

	color.Yellow("[2/3] Building Angular frontend...\n")
	if err := buildFrontend(projectRoot, outputPath); err != nil {
		return fmt.Errorf("frontend build failed: %w", err)
	}
	color.Green("[2/3] Frontend build complete.\n")

	color.Yellow("[3/3] Copying static assets...\n")
	if err := copyAssets(projectRoot, outputPath); err != nil {
		return fmt.Errorf("asset copy failed: %w", err)
	}
	color.Green("[3/3] Assets copied.\n")

	color.Green("\nBuild complete! Output: %s\n", outputPath)

	return nil
}

func findBuildProjectRoot() (string, error) {
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

/*
 * buildBackend compiles the Go backend for the target platform.
 * Uses cross-compilation flags for different OS/arch combinations.
 */
func buildBackend(projectRoot, outputPath string) error {
	binaryName := "server"
	if buildPlatform == "windows" {
		binaryName = "server.exe"
	}

	binaryPath := filepath.Join(outputPath, binaryName)
	appDir := filepath.Join(projectRoot, "app")

	buildCmd := exec.Command("go", "build",
		"-ldflags", "-s -w",
		"-o", binaryPath,
		"./cmd/server",
	)

	buildCmd.Dir = appDir
	buildCmd.Env = append(os.Environ(),
		fmt.Sprintf("GOOS=%s", buildPlatform),
		fmt.Sprintf("GOARCH=%s", buildArch),
		"CGO_ENABLED=0",
	)

	buildCmd.Stdout = os.Stdout
	buildCmd.Stderr = os.Stderr

	return buildCmd.Run()
}

/*
 * buildFrontend runs the Angular production build.
 * Enables AOT compilation and output hashing for cache busting.
 */
func buildFrontend(projectRoot, outputPath string) error {
	webOutputPath := filepath.Join(outputPath, "public")
	webDir := filepath.Join(projectRoot, "web")

	var ngBuildCmd *exec.Cmd
	if os.PathSeparator == '\\' && os.PathListSeparator == ';' {
		ngBuildCmd = exec.Command("cmd", "/c", "npx", "ng", "build",
			"--configuration", "production",
			"--output-path", webOutputPath,
		)
	} else {
		ngBuildCmd = exec.Command("npx", "ng", "build",
			"--configuration", "production",
			"--output-path", webOutputPath,
		)
	}

	ngBuildCmd.Dir = webDir
	ngBuildCmd.Stdout = os.Stdout
	ngBuildCmd.Stderr = os.Stderr

	return ngBuildCmd.Run()
}

/*
 * copyAssets transfers static files to the output directory.
 * Includes configuration files and non-compiled resources.
 */
func copyAssets(projectRoot, outputPath string) error {
	configSrc := filepath.Join(projectRoot, "goastra.json")
	configDst := filepath.Join(outputPath, "goastra.json")

	if _, err := os.Stat(configSrc); err == nil {
		input, err := os.ReadFile(configSrc)
		if err != nil {
			return err
		}
		if err := os.WriteFile(configDst, input, 0644); err != nil {
			return err
		}
	}

	envSrc := filepath.Join(projectRoot, ".env.production")
	envDst := filepath.Join(outputPath, ".env")

	if _, err := os.Stat(envSrc); err == nil {
		input, err := os.ReadFile(envSrc)
		if err != nil {
			return err
		}
		if err := os.WriteFile(envDst, input, 0644); err != nil {
			return err
		}
	}

	return nil
}
