/*
 * GoAstra CLI - Test Command
 *
 * Unified test runner for both Go backend and Angular frontend.
 * Loads .env.test environment and provides consolidated reporting.
 */
package cmd

import (
	"fmt"
	"os"
	"os/exec"

	"github.com/fatih/color"
	"github.com/channdev/goastra/cli/internal/env"
	"github.com/spf13/cobra"
)

var (
	testCoverage    bool
	testVerbose     bool
	testBackendOnly bool
	testFrontendOnly bool
	testWatch       bool
)

/*
 * testCmd runs the test suites for both backend and frontend.
 * Provides unified output and coverage reporting.
 */
var testCmd = &cobra.Command{
	Use:   "test",
	Short: "Run test suites",
	Long: `Runs tests for the GoAstra project:
  - Go backend tests with go test
  - Angular frontend tests with ng test
  - Loads .env.test environment
  - Supports coverage reporting
  - Optional watch mode for TDD`,
	RunE: runTest,
}

func init() {
	rootCmd.AddCommand(testCmd)

	testCmd.Flags().BoolVarP(&testCoverage, "coverage", "c", false, "generate coverage report")
	testCmd.Flags().BoolVarP(&testVerbose, "verbose", "V", false, "verbose test output")
	testCmd.Flags().BoolVar(&testBackendOnly, "backend", false, "run backend tests only")
	testCmd.Flags().BoolVar(&testFrontendOnly, "frontend", false, "run frontend tests only")
	testCmd.Flags().BoolVarP(&testWatch, "watch", "w", false, "watch mode for continuous testing")
}

/*
 * runTest executes the test suites based on flags.
 * Coordinates Go and Angular test execution.
 */
func runTest(cmd *cobra.Command, args []string) error {
	if err := env.Load("test"); err != nil {
		return fmt.Errorf("failed to load test environment: %w", err)
	}

	color.Cyan("Running GoAstra tests...\n")

	hasErrors := false

	if !testFrontendOnly {
		color.Yellow("\n[Backend Tests]\n")
		if err := runBackendTests(); err != nil {
			color.Red("Backend tests failed: %v\n", err)
			hasErrors = true
		} else {
			color.Green("Backend tests passed!\n")
		}
	}

	if !testBackendOnly {
		color.Yellow("\n[Frontend Tests]\n")
		if err := runFrontendTests(); err != nil {
			color.Red("Frontend tests failed: %v\n", err)
			hasErrors = true
		} else {
			color.Green("Frontend tests passed!\n")
		}
	}

	if hasErrors {
		return fmt.Errorf("some tests failed")
	}

	color.Green("\nAll tests passed!\n")
	return nil
}

/*
 * runBackendTests executes Go test suite.
 * Supports coverage and verbose modes.
 */
func runBackendTests() error {
	args := []string{"test"}

	if testVerbose {
		args = append(args, "-v")
	}

	if testCoverage {
		args = append(args, "-cover", "-coverprofile=coverage.out")
	}

	args = append(args, "./...")

	testCmd := exec.Command("go", args...)
	testCmd.Dir = "app"
	testCmd.Stdout = os.Stdout
	testCmd.Stderr = os.Stderr
	testCmd.Env = os.Environ()

	return testCmd.Run()
}

/*
 * runFrontendTests executes Angular test suite.
 * Uses Karma runner with Chrome headless by default.
 */
func runFrontendTests() error {
	args := []string{"test"}

	if !testWatch {
		args = append(args, "--watch=false")
	}

	if testCoverage {
		args = append(args, "--code-coverage")
	}

	args = append(args, "--browsers=ChromeHeadless")

	var testCmd *exec.Cmd
	if os.PathSeparator == '\\' {
		// Windows: use npx to run ng
		cmdArgs := append([]string{"/c", "npx", "ng"}, args...)
		testCmd = exec.Command("cmd", cmdArgs...)
	} else {
		testCmd = exec.Command("ng", args...)
	}
	testCmd.Dir = "web"
	testCmd.Stdout = os.Stdout
	testCmd.Stderr = os.Stderr
	testCmd.Env = os.Environ()

	return testCmd.Run()
}
