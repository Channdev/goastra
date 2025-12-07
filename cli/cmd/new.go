package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/channdev/goastra/cli/internal/scaffold"
	"github.com/spf13/cobra"
)

var (
	skipAngular  bool
	skipBackend  bool
	useGraphQL   bool
	templateName string
	dbDriver     string
)

var newCmd = &cobra.Command{
	Use:   "new <project-name>",
	Short: "Create a new GoAstra project",
	Long: `Creates a new GoAstra project with Go backend and Angular frontend.

Templates:
  default   - Full-featured template with auth, dashboard, and beautiful landing page
  minimal   - Minimal starter template with basic structure

Database:
  postgres  - PostgreSQL (default)
  mysql     - MySQL/MariaDB`,
	Args: cobra.ExactArgs(1),
	RunE: runNew,
}

func init() {
	rootCmd.AddCommand(newCmd)
	newCmd.Flags().BoolVar(&skipAngular, "skip-angular", false, "skip Angular frontend generation")
	newCmd.Flags().BoolVar(&skipBackend, "skip-backend", false, "skip Go backend generation")
	newCmd.Flags().BoolVar(&useGraphQL, "graphql", false, "use GraphQL instead of REST")
	newCmd.Flags().StringVarP(&templateName, "template", "t", "default", "project template (default, minimal)")
	newCmd.Flags().StringVar(&dbDriver, "db", "postgres", "database driver (postgres, mysql)")
}

func runNew(cmd *cobra.Command, args []string) error {
	projectName := args[0]

	if err := validateProjectName(projectName); err != nil {
		return err
	}

	if templateName != "default" && templateName != "minimal" {
		return fmt.Errorf("invalid template: %s (use 'default' or 'minimal')", templateName)
	}

	if dbDriver != "postgres" && dbDriver != "mysql" {
		return fmt.Errorf("invalid database driver: %s (use 'postgres' or 'mysql')", dbDriver)
	}

	projectPath, err := filepath.Abs(projectName)
	if err != nil {
		return fmt.Errorf("failed to resolve project path: %w", err)
	}

	if _, err := os.Stat(projectPath); !os.IsNotExist(err) {
		return fmt.Errorf("directory already exists: %s", projectPath)
	}

	return scaffold.CreateProject(scaffold.Options{
		ProjectName:  projectName,
		ProjectPath:  projectPath,
		Template:     templateName,
		DBDriver:     dbDriver,
		SkipBackend:  skipBackend,
		SkipFrontend: skipAngular,
	})
}

func validateProjectName(name string) error {
	if len(name) == 0 {
		return fmt.Errorf("project name cannot be empty")
	}
	if len(name) > 100 {
		return fmt.Errorf("project name too long (max 100 characters)")
	}
	for _, c := range name {
		if !isValidNameChar(c) {
			return fmt.Errorf("invalid character in project name: %c", c)
		}
	}
	return nil
}

func isValidNameChar(c rune) bool {
	return (c >= 'a' && c <= 'z') || (c >= 'A' && c <= 'Z') || (c >= '0' && c <= '9') || c == '-' || c == '_'
}
