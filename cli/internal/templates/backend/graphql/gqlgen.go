/*
 * GoAstra CLI - gqlgen Config Template
 *
 * Generates the gqlgen.yml configuration file.
 * Configures code generation for GraphQL server.
 */
package graphql

// GqlgenYML returns the gqlgen.yml configuration template.
func GqlgenYML() string {
	return `# gqlgen Configuration
# See https://gqlgen.com/config/ for documentation

# Where to find the GraphQL schema files
schema:
  - graph/*.graphqls

# Where to write generated code
exec:
  filename: graph/generated/generated.go
  package: generated

# Where to write the models
model:
  filename: graph/model/models_gen.go
  package: model

# Where to write the resolver interface
resolver:
  layout: follow-schema
  dir: graph
  package: graph
  filename_template: "{name}.resolvers.go"

# Optional: autobind Go types to GraphQL types
autobind:
  - "app/graph/model"

# Model mappings (map GraphQL types to Go types)
models:
  ID:
    model:
      - github.com/99designs/gqlgen/graphql.ID
      - github.com/99designs/gqlgen/graphql.Int
      - github.com/99designs/gqlgen/graphql.Int64
      - github.com/99designs/gqlgen/graphql.Int32
  Int:
    model:
      - github.com/99designs/gqlgen/graphql.Int
      - github.com/99designs/gqlgen/graphql.Int64
      - github.com/99designs/gqlgen/graphql.Int32
`
}

// GenerateGo returns the generate.go directive for gqlgen.
func GenerateGo() string {
	return `package graph

//go:generate go run github.com/99designs/gqlgen generate
`
}

// ToolsGo returns the tools.go file for managing tool dependencies.
func ToolsGo() string {
	return `//go:build tools
// +build tools

package tools

import (
	_ "github.com/99designs/gqlgen"
	_ "github.com/99designs/gqlgen/graphql/introspection"
)
`
}
