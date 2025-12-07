package cmd

import (
	"context"
	"fmt"
	"net"
	"os"
	"os/exec"
	"os/signal"
	"path/filepath"
	"sync"
	"syscall"
	"time"

	"github.com/fatih/color"
	"github.com/goastra/cli/internal/env"
	"github.com/spf13/cobra"
)

var (
	devPort         int
	devFrontendPort int
	devBackendOnly  bool
	devFrontendOnly bool
)

var devCmd = &cobra.Command{
	Use:   "dev",
	Short: "Start development servers",
	Long:  "Starts Go backend and Angular frontend development servers",
	RunE:  runDev,
}

func init() {
	rootCmd.AddCommand(devCmd)
	devCmd.Flags().IntVarP(&devPort, "port", "p", 8080, "backend server port")
	devCmd.Flags().IntVar(&devFrontendPort, "frontend-port", 4200, "frontend server port")
	devCmd.Flags().BoolVar(&devBackendOnly, "backend", false, "run backend only")
	devCmd.Flags().BoolVar(&devFrontendOnly, "frontend", false, "run frontend only")
}

func runDev(cmd *cobra.Command, args []string) error {
	projectRoot, err := findProjectRoot()
	if err != nil {
		return fmt.Errorf("not in a GoAstra project: %w", err)
	}

	env.Load("development")

	color.Cyan("Starting GoAstra development servers...\n")

	if !devFrontendOnly {
		if isPortInUse(devPort) {
			color.Yellow("[Backend] Port %d is in use, trying %d...\n", devPort, devPort+1)
			devPort++
		}
	}

	if !devBackendOnly {
		if isPortInUse(devFrontendPort) {
			color.Yellow("[Frontend] Port %d is in use, trying %d...\n", devFrontendPort, devFrontendPort+1)
			devFrontendPort++
		}
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	var wg sync.WaitGroup
	errChan := make(chan error, 2)

	if !devFrontendOnly {
		wg.Add(1)
		go func() {
			defer wg.Done()
			color.Green("[Backend] Starting on port %d...\n", devPort)
			if err := runBackendServer(ctx, projectRoot, devPort); err != nil {
				select {
				case <-ctx.Done():
				default:
					errChan <- fmt.Errorf("backend error: %w", err)
				}
			}
		}()
	}

	if !devBackendOnly {
		wg.Add(1)
		go func() {
			defer wg.Done()
			time.Sleep(500 * time.Millisecond)
			color.Blue("[Frontend] Starting on port %d...\n", devFrontendPort)
			if err := runFrontendServer(ctx, projectRoot, devFrontendPort); err != nil {
				select {
				case <-ctx.Done():
				default:
					errChan <- fmt.Errorf("frontend error: %w", err)
				}
			}
		}()
	}

	select {
	case <-sigChan:
		color.Yellow("\nShutting down development servers...\n")
		cancel()
	case err := <-errChan:
		color.Red("Error: %v\n", err)
		cancel()
	}

	wg.Wait()
	color.Green("Development servers stopped.\n")
	return nil
}

func runBackendServer(ctx context.Context, projectRoot string, port int) error {
	appDir := filepath.Join(projectRoot, "app")
	cmd := exec.CommandContext(ctx, "go", "run", "./cmd/server")
	cmd.Dir = appDir
	cmd.Env = append(os.Environ(),
		fmt.Sprintf("PORT=%d", port),
		"APP_ENV=development",
	)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

func runFrontendServer(ctx context.Context, projectRoot string, port int) error {
	webDir := filepath.Join(projectRoot, "web")
	var cmd *exec.Cmd
	if isWindows() {
		cmd = exec.CommandContext(ctx, "cmd", "/c", "npx", "ng", "serve", "--port", fmt.Sprintf("%d", port), "--proxy-config", "proxy.conf.json")
	} else {
		cmd = exec.CommandContext(ctx, "npx", "ng", "serve", "--port", fmt.Sprintf("%d", port), "--proxy-config", "proxy.conf.json")
	}
	cmd.Dir = webDir
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

func findProjectRoot() (string, error) {
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

func isWindows() bool {
	return os.PathSeparator == '\\' && os.PathListSeparator == ';'
}

func isPortInUse(port int) bool {
	conn, err := net.DialTimeout("tcp", fmt.Sprintf("127.0.0.1:%d", port), time.Second)
	if err != nil {
		return false
	}
	conn.Close()
	return true
}
