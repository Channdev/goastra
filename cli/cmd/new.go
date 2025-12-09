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
	templateName string
	dbDriver     string
	apiType      string
	ormType      string
)

var newCmd = &cobra.Command{
	Use:   "new <project-name>",
	Short: "Create a new GoAstra project",
	Long: `Creates a new GoAstra project with Go backend and Angular frontend.

Templates:
  default   - Full-featured template with auth, dashboard, and beautiful landing page
  minimal   - Minimal starter template with basic structure

API Types:
  rest      - REST API with Gin framework (default)
  graphql   - GraphQL API with gqlgen
  trpc      - tRPC API with Connect-Go

ORM Options:
  sqlx      - SQLx for raw SQL with type safety (default)
  ent       - Ent ORM by Facebook (Prisma-like experience)

Database:
  postgres  - PostgreSQL (default)
  mysql     - MySQL/MariaDB

Examples:
  goastra new my-app                              # REST + SQLx + PostgreSQL
  goastra new my-app --api graphql --orm ent      # GraphQL + Ent + PostgreSQL
  goastra new my-app --api trpc --db mysql        # tRPC + SQLx + MySQL
  goastra new my-app --api rest --orm ent -t minimal  # REST + Ent + minimal template`,
	Args: cobra.ExactArgs(1),
	RunE: runNew,
}

func init() {
	rootCmd.AddCommand(newCmd)
	newCmd.Flags().BoolVar(&skipAngular, "skip-angular", false, "skip Angular frontend generation")
	newCmd.Flags().BoolVar(&skipBackend, "skip-backend", false, "skip Go backend generation")
	newCmd.Flags().StringVarP(&templateName, "template", "t", "default", "project template (default, minimal)")
	newCmd.Flags().StringVar(&dbDriver, "db", "postgres", "database driver (postgres, mysql)")
	newCmd.Flags().StringVar(&apiType, "api", "rest", "API type (rest, graphql, trpc)")
	newCmd.Flags().StringVar(&ormType, "orm", "sqlx", "ORM type (sqlx, ent)")
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

	if apiType != "rest" && apiType != "graphql" && apiType != "trpc" {
		return fmt.Errorf("invalid API type: %s (use 'rest', 'graphql', or 'trpc')", apiType)
	}

	if ormType != "sqlx" && ormType != "ent" {
		return fmt.Errorf("invalid ORM type: %s (use 'sqlx' or 'ent')", ormType)
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
		APIType:      apiType,
		ORMType:      ormType,
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
