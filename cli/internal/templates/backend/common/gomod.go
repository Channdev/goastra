/*
 * GoAstra CLI - Go Module Template
 *
 * Generates go.mod with dependencies based on API type and ORM choice.
 * Supports REST (Gin), GraphQL (gqlgen), tRPC (Connect-Go), SQLx, and Ent.
 */
package common

import "fmt"

// GoModOptions configures go.mod generation.
type GoModOptions struct {
	ProjectName string
	DBDriver    string // postgres, mysql
	APIType     string // rest, graphql, trpc
	ORMType     string // sqlx, ent
}

// GoMod returns the go.mod template with conditional dependencies.
func GoMod(opts GoModOptions) string {
	// Base dependencies (always included)
	deps := `	github.com/joho/godotenv v1.5.1
	go.uber.org/zap v1.26.0
	github.com/golang-jwt/jwt/v5 v5.2.0
	golang.org/x/crypto v0.16.0
	github.com/go-playground/validator/v10 v10.16.0`

	// Database driver
	dbDep := "\n\tgithub.com/lib/pq v1.10.9"
	if opts.DBDriver == "mysql" {
		dbDep = "\n\tgithub.com/go-sql-driver/mysql v1.7.1"
	}

	// API-specific dependencies
	apiDeps := ""
	switch opts.APIType {
	case "graphql":
		apiDeps = `
	github.com/99designs/gqlgen v0.17.43
	github.com/vektah/gqlparser/v2 v2.5.10
	github.com/gin-gonic/gin v1.9.1`
	case "trpc":
		apiDeps = `
	connectrpc.com/connect v1.14.0
	google.golang.org/protobuf v1.32.0
	golang.org/x/net v0.19.0`
	default: // rest
		apiDeps = `
	github.com/gin-gonic/gin v1.9.1`
	}

	// ORM-specific dependencies
	ormDeps := ""
	switch opts.ORMType {
	case "ent":
		ormDeps = `
	entgo.io/ent v0.12.5
	ariga.io/atlas v0.19.0`
	default: // sqlx
		ormDeps = `
	github.com/jmoiron/sqlx v1.3.5`
	}

	return fmt.Sprintf(`module github.com/%s/app

go 1.21

require (
%s%s%s%s
)
`, opts.ProjectName, deps, dbDep, apiDeps, ormDeps)
}

// GoModREST returns go.mod for REST API projects (backward compatibility).
func GoModREST(projectName, dbDriver string) string {
	return GoMod(GoModOptions{
		ProjectName: projectName,
		DBDriver:    dbDriver,
		APIType:     "rest",
		ORMType:     "sqlx",
	})
}

// GoModGraphQL returns go.mod for GraphQL projects.
func GoModGraphQL(projectName, dbDriver, ormType string) string {
	return GoMod(GoModOptions{
		ProjectName: projectName,
		DBDriver:    dbDriver,
		APIType:     "graphql",
		ORMType:     ormType,
	})
}

// GoModTRPC returns go.mod for tRPC projects.
func GoModTRPC(projectName, dbDriver, ormType string) string {
	return GoMod(GoModOptions{
		ProjectName: projectName,
		DBDriver:    dbDriver,
		APIType:     "trpc",
		ORMType:     ormType,
	})
}
