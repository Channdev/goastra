/*
 * GoAstra CLI - Ent Generate Template
 *
 * Generates Ent ORM code generation directives.
 * Provides go:generate comments for ent generate command.
 */
package ent

// GenerateGo returns the generate.go template with go:generate directive.
func GenerateGo() string {
	return `package ent

//go:generate go run -mod=mod entgo.io/ent/cmd/ent generate ./schema
`
}

// EntrcGo returns the .entrc configuration file template.
func EntrcGo() string {
	return `{
  "schema": "./ent/schema",
  "target": "./ent",
  "package": "app/ent",
  "features": [
    "sql/modifier",
    "sql/upsert",
    "namedges"
  ]
}
`
}

// SchemaIndexGo returns the schema package index file.
func SchemaIndexGo() string {
	return `// Package schema contains the Ent schema definitions.
// Add new entity schemas to this directory.
//
// Example:
//
//	type Product struct {
//	    ent.Schema
//	}
//
//	func (Product) Fields() []ent.Field {
//	    return []ent.Field{
//	        field.String("name").NotEmpty(),
//	        field.Float("price").Positive(),
//	    }
//	}
//
// Run 'go generate ./ent' to regenerate code after changes.
package schema
`
}
