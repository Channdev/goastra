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
	"github.com/channdev/goastra/cli/internal/generator"
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
  goastra generate crud <name>     Generate full CRUD stack
  goastra generate graphql <name>  Generate GraphQL schema and resolvers
  goastra generate trpc <name>     Generate tRPC proto and service
  goastra generate ent <name>      Generate Ent ORM schema`,
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

/*
 * generateGraphQLCmd creates GraphQL schema and resolvers.
 * Generates schema definitions and resolver implementations.
 */
var generateGraphQLCmd = &cobra.Command{
	Use:   "graphql <name>",
	Short: "Generate GraphQL schema and resolvers",
	Long: `Generates GraphQL artifacts for a resource:
  - Schema file with types, queries, and mutations
  - Resolver implementations with CRUD operations

After generation, run 'go generate ./...' to regenerate gqlgen code.`,
	Args: cobra.ExactArgs(1),
	RunE: runGenerateGraphQL,
}

/*
 * generateTRPCCmd creates tRPC proto and service.
 * Generates Protocol Buffer definitions and Connect-Go service.
 */
var generateTRPCCmd = &cobra.Command{
	Use:   "trpc <name>",
	Short: "Generate tRPC proto and service",
	Long: `Generates tRPC artifacts for a resource:
  - Protocol Buffer definition with service and messages
  - Connect-Go service implementation

After generation, run 'buf generate' to regenerate proto code.`,
	Args: cobra.ExactArgs(1),
	RunE: runGenerateTRPC,
}

/*
 * generateEntCmd creates an Ent ORM schema.
 * Generates entity schema with fields, edges, and indexes.
 */
var generateEntCmd = &cobra.Command{
	Use:   "ent <name>",
	Short: "Generate Ent ORM schema",
	Long: `Generates an Ent schema for an entity:
  - Schema file with fields, edges, and indexes
  - Timestamps and soft delete mixin support

After generation, run 'go generate ./ent' to regenerate Ent code.`,
	Args: cobra.ExactArgs(1),
	RunE: runGenerateEnt,
}

func init() {
	rootCmd.AddCommand(generateCmd)
	generateCmd.AddCommand(generateAPICmd)
	generateCmd.AddCommand(generateModuleCmd)
	generateCmd.AddCommand(generateCRUDCmd)
	generateCmd.AddCommand(generateGraphQLCmd)
	generateCmd.AddCommand(generateTRPCCmd)
	generateCmd.AddCommand(generateEntCmd)
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
 * runGenerateGraphQL executes GraphQL generation workflow.
 * Creates schema and resolver files.
 */
func runGenerateGraphQL(cmd *cobra.Command, args []string) error {
	name := args[0]
	normalizedName := normalizeResourceName(name)

	color.Cyan("Generating GraphQL schema and resolvers: %s\n", normalizedName)

	gen := generator.NewGraphQLGenerator(normalizedName)

	if err := gen.GenerateSchema(); err != nil {
		return fmt.Errorf("failed to generate schema: %w", err)
	}

	if err := gen.GenerateResolver(); err != nil {
		return fmt.Errorf("failed to generate resolver: %w", err)
	}

	color.Green("GraphQL artifacts generated successfully!\n")
	fmt.Printf("\nGenerated files:\n")
	fmt.Printf("  app/graph/%s.graphqls\n", normalizedName)
	fmt.Printf("  app/graph/%s.resolvers.go\n", normalizedName)
	fmt.Printf("\nNext steps:\n")
	fmt.Printf("  1. Add fields to the schema\n")
	fmt.Printf("  2. Run: go generate ./...\n")
	fmt.Printf("  3. Implement resolver methods\n")

	return nil
}

/*
 * runGenerateTRPC executes tRPC generation workflow.
 * Creates proto and service files.
 */
func runGenerateTRPC(cmd *cobra.Command, args []string) error {
	name := args[0]
	normalizedName := normalizeResourceName(name)

	color.Cyan("Generating tRPC proto and service: %s\n", normalizedName)

	gen := generator.NewTRPCGenerator(normalizedName)

	if err := gen.GenerateProto(); err != nil {
		return fmt.Errorf("failed to generate proto: %w", err)
	}

	if err := gen.GenerateService(); err != nil {
		return fmt.Errorf("failed to generate service: %w", err)
	}

	color.Green("tRPC artifacts generated successfully!\n")
	fmt.Printf("\nGenerated files:\n")
	fmt.Printf("  app/proto/v1/%s.proto\n", strings.ReplaceAll(normalizedName, "-", "_"))
	fmt.Printf("  app/internal/rpc/%s_service.go\n", strings.ReplaceAll(normalizedName, "-", "_"))
	fmt.Printf("\nNext steps:\n")
	fmt.Printf("  1. Add fields to the proto\n")
	fmt.Printf("  2. Run: buf generate\n")
	fmt.Printf("  3. Register service in main.go\n")
	fmt.Printf("  4. Implement service methods\n")

	return nil
}

/*
 * runGenerateEnt executes Ent schema generation workflow.
 * Creates entity schema file.
 */
func runGenerateEnt(cmd *cobra.Command, args []string) error {
	name := args[0]
	normalizedName := normalizeResourceName(name)

	color.Cyan("Generating Ent schema: %s\n", normalizedName)

	gen := generator.NewEntGenerator(normalizedName)

	if err := gen.GenerateSchema(); err != nil {
		return fmt.Errorf("failed to generate schema: %w", err)
	}

	color.Green("Ent schema generated successfully!\n")
	fmt.Printf("\nGenerated files:\n")
	fmt.Printf("  app/ent/schema/%s.go\n", strings.ReplaceAll(normalizedName, "-", "_"))
	fmt.Printf("\nNext steps:\n")
	fmt.Printf("  1. Add fields to the schema\n")
	fmt.Printf("  2. Run: go generate ./ent\n")
	fmt.Printf("  3. Use generated client in your code\n")

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
