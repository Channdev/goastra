/*
 * GoAstra CLI - TypeSync Command
 *
 * Synchronizes Go struct definitions to TypeScript interfaces.
 * Parses schema package and generates corresponding TypeScript models
 * and Angular service stubs for type-safe frontend development.
 */
package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/fatih/color"
	"github.com/goastra/cli/internal/codegen"
	"github.com/spf13/cobra"
)

var (
	typesyncOutput  string
	typesyncWatch   bool
	typesyncService bool
)

/*
 * typesyncCmd generates TypeScript from Go type definitions.
 * Maintains type safety across the full stack.
 */
var typesyncCmd = &cobra.Command{
	Use:   "typesync",
	Short: "Sync Go types to TypeScript",
	Long: `Generates TypeScript interfaces from Go struct definitions:
  - Parses schema/types/*.go files
  - Generates TypeScript interfaces in web/src/app/core/models/
  - Optionally generates Angular services for API calls
  - Supports watch mode for continuous synchronization`,
	RunE: runTypesync,
}

func init() {
	rootCmd.AddCommand(typesyncCmd)

	typesyncCmd.Flags().StringVarP(&typesyncOutput, "output", "o", "web/src/app/core/models", "output directory for TypeScript files")
	typesyncCmd.Flags().BoolVarP(&typesyncWatch, "watch", "w", false, "watch for changes and regenerate")
	typesyncCmd.Flags().BoolVar(&typesyncService, "services", false, "also generate Angular services")
}

/*
 * runTypesync executes the type synchronization process.
 * Parses Go files and generates corresponding TypeScript.
 */
func runTypesync(cmd *cobra.Command, args []string) error {
	schemaPath := "schema/types"
	outputPath, err := filepath.Abs(typesyncOutput)
	if err != nil {
		return fmt.Errorf("invalid output path: %w", err)
	}

	if _, err := os.Stat(schemaPath); os.IsNotExist(err) {
		return fmt.Errorf("schema directory not found: %s", schemaPath)
	}

	if err := os.MkdirAll(outputPath, 0755); err != nil {
		return fmt.Errorf("failed to create output directory: %w", err)
	}

	color.Cyan("Syncing Go types to TypeScript...\n")

	parser := codegen.NewGoParser(schemaPath)
	types, err := parser.Parse()
	if err != nil {
		return fmt.Errorf("failed to parse Go types: %w", err)
	}

	color.Yellow("Found %d type definitions\n", len(types))

	tsGenerator := codegen.NewTypeScriptGenerator(outputPath)
	if err := tsGenerator.Generate(types); err != nil {
		return fmt.Errorf("failed to generate TypeScript: %w", err)
	}

	if typesyncService {
		color.Yellow("Generating Angular services...\n")
		serviceGen := codegen.NewServiceGenerator(outputPath)
		if err := serviceGen.Generate(types); err != nil {
			return fmt.Errorf("failed to generate services: %w", err)
		}
	}

	color.Green("TypeScript generation complete!\n")
	fmt.Printf("Output: %s\n", outputPath)

	if typesyncWatch {
		color.Cyan("Watching for changes... (Ctrl+C to stop)\n")
		return watchForChanges(schemaPath, outputPath)
	}

	return nil
}

/*
 * watchForChanges monitors the schema directory for modifications.
 * Triggers regeneration when Go files are updated.
 */
func watchForChanges(schemaPath, outputPath string) error {
	/* TODO: Implement fsnotify watcher for continuous sync */
	select {}
}
