/*
 * GoAstra CLI - Generate Command
 *
 * Code generation commands for creating backend APIs, frontend modules,
 * and full CRUD implementations spanning both layers.
 */
package cmd

import (
	"fmt"
	"strings"

	"github.com/fatih/color"
	"github.com/goastra/cli/internal/generator"
	"github.com/spf13/cobra"
)

/*
 * generateCmd is the parent command for all code generation operations.
 * Provides subcommands for different artifact types.
 */
var generateCmd = &cobra.Command{
	Use:     "generate",
	Aliases: []string{"g"},
	Short:   "Generate code artifacts",
	Long: `Generate code artifacts for your GoAstra project:
  goastra generate api <name>      Generate REST API endpoint
  goastra generate module <name>   Generate Angular feature module
  goastra generate crud <name>     Generate full CRUD stack`,
}

/*
 * generateAPICmd creates a new REST API endpoint.
 * Generates handler, service, repository, and routes.
 */
var generateAPICmd = &cobra.Command{
	Use:   "api <name>",
	Short: "Generate REST API endpoint",
	Long: `Generates a complete REST API endpoint:
  - Handler with CRUD operations
  - Service layer with business logic
  - Repository for data access
  - Route registration
  - Request/Response DTOs`,
	Args: cobra.ExactArgs(1),
	RunE: runGenerateAPI,
}

/*
 * generateModuleCmd creates a new Angular feature module.
 * Generates components, services, and routing configuration.
 */
var generateModuleCmd = &cobra.Command{
	Use:   "module <name>",
	Short: "Generate Angular feature module",
	Long: `Generates an Angular feature module with:
  - Module file with lazy loading support
  - Routing module
  - Container component
  - Feature service
  - State management (if configured)`,
	Args: cobra.ExactArgs(1),
	RunE: runGenerateModule,
}

/*
 * generateCRUDCmd creates full-stack CRUD implementation.
 * Generates both backend API and frontend module with all operations.
 */
var generateCRUDCmd = &cobra.Command{
	Use:   "crud <name>",
	Short: "Generate full CRUD stack",
	Long: `Generates a complete CRUD implementation:
  Backend:
    - Model definition
    - Handler with Create, Read, Update, Delete
    - Service layer
    - Repository with database operations
    - Migration file
  Frontend:
    - Feature module with routing
    - List component with pagination
    - Detail/View component
    - Create/Edit form component
    - Delete confirmation
    - API service`,
	Args: cobra.ExactArgs(1),
	RunE: runGenerateCRUD,
}

func init() {
	rootCmd.AddCommand(generateCmd)
	generateCmd.AddCommand(generateAPICmd)
	generateCmd.AddCommand(generateModuleCmd)
	generateCmd.AddCommand(generateCRUDCmd)
}

/*
 * runGenerateAPI executes the API generation workflow.
 * Creates all backend components for a REST endpoint.
 */
func runGenerateAPI(cmd *cobra.Command, args []string) error {
	name := args[0]
	normalizedName := normalizeResourceName(name)

	color.Cyan("Generating API endpoint: %s\n", normalizedName)

	gen := generator.NewAPIGenerator(normalizedName)

	if err := gen.GenerateHandler(); err != nil {
		return fmt.Errorf("failed to generate handler: %w", err)
	}

	if err := gen.GenerateService(); err != nil {
		return fmt.Errorf("failed to generate service: %w", err)
	}

	if err := gen.GenerateRepository(); err != nil {
		return fmt.Errorf("failed to generate repository: %w", err)
	}

	if err := gen.GenerateRoutes(); err != nil {
		return fmt.Errorf("failed to generate routes: %w", err)
	}

	color.Green("API endpoint generated successfully!\n")
	fmt.Printf("\nGenerated files:\n")
	fmt.Printf("  app/internal/handlers/%s_handler.go\n", normalizedName)
	fmt.Printf("  app/internal/services/%s_service.go\n", normalizedName)
	fmt.Printf("  app/internal/repository/%s_repository.go\n", normalizedName)

	return nil
}

/*
 * runGenerateModule executes the Angular module generation.
 * Creates feature module with lazy loading configuration.
 */
func runGenerateModule(cmd *cobra.Command, args []string) error {
	name := args[0]
	normalizedName := normalizeResourceName(name)

	color.Cyan("Generating Angular module: %s\n", normalizedName)

	gen := generator.NewModuleGenerator(normalizedName)

	if err := gen.GenerateModule(); err != nil {
		return fmt.Errorf("failed to generate module: %w", err)
	}

	if err := gen.GenerateRouting(); err != nil {
		return fmt.Errorf("failed to generate routing: %w", err)
	}

	if err := gen.GenerateComponent(); err != nil {
		return fmt.Errorf("failed to generate component: %w", err)
	}

	if err := gen.GenerateService(); err != nil {
		return fmt.Errorf("failed to generate service: %w", err)
	}

	color.Green("Angular module generated successfully!\n")
	fmt.Printf("\nGenerated files:\n")
	fmt.Printf("  web/src/app/features/%s/%s.module.ts\n", normalizedName, normalizedName)
	fmt.Printf("  web/src/app/features/%s/%s-routing.module.ts\n", normalizedName, normalizedName)
	fmt.Printf("  web/src/app/features/%s/%s.component.ts\n", normalizedName, normalizedName)
	fmt.Printf("  web/src/app/features/%s/%s.service.ts\n", normalizedName, normalizedName)

	return nil
}

/*
 * runGenerateCRUD executes full-stack CRUD generation.
 * Combines API and module generation with additional CRUD components.
 */
func runGenerateCRUD(cmd *cobra.Command, args []string) error {
	name := args[0]
	normalizedName := normalizeResourceName(name)

	color.Cyan("Generating CRUD stack: %s\n", normalizedName)

	gen := generator.NewCRUDGenerator(normalizedName)

	color.Yellow("[1/6] Generating model...\n")
	if err := gen.GenerateModel(); err != nil {
		return fmt.Errorf("failed to generate model: %w", err)
	}

	color.Yellow("[2/6] Generating API...\n")
	if err := gen.GenerateAPI(); err != nil {
		return fmt.Errorf("failed to generate API: %w", err)
	}

	color.Yellow("[3/6] Generating migration...\n")
	if err := gen.GenerateMigration(); err != nil {
		return fmt.Errorf("failed to generate migration: %w", err)
	}

	color.Yellow("[4/6] Generating Angular module...\n")
	if err := gen.GenerateModule(); err != nil {
		return fmt.Errorf("failed to generate module: %w", err)
	}

	color.Yellow("[5/6] Generating CRUD components...\n")
	if err := gen.GenerateComponents(); err != nil {
		return fmt.Errorf("failed to generate components: %w", err)
	}

	color.Yellow("[6/6] Updating routes...\n")
	if err := gen.UpdateRoutes(); err != nil {
		return fmt.Errorf("failed to update routes: %w", err)
	}

	color.Green("\nCRUD stack generated successfully!\n")

	return nil
}

/*
 * normalizeResourceName converts input to lowercase with hyphens.
 * Ensures consistent naming across all generated files.
 */
func normalizeResourceName(name string) string {
	name = strings.ToLower(name)
	name = strings.ReplaceAll(name, "_", "-")
	name = strings.ReplaceAll(name, " ", "-")
	return name
}
