/*
 * GoAstra CLI - Project Generator
 *
 * DEPRECATED: Use github.com/channdev/goastra/cli/internal/scaffold instead.
 * This file is kept for backwards compatibility but delegates to scaffold.
 */
package project

import (
	"github.com/channdev/goastra/cli/internal/scaffold"
)

type Config struct {
	Name        string
	Path        string
	SkipAngular bool
	SkipBackend bool
	UseGraphQL  bool
	Template    string
	DBDriver    string
}

type Generator struct {
	config *Config
}

func NewGenerator(cfg *Config) *Generator {
	return &Generator{config: cfg}
}

func (g *Generator) Generate() error {
	template := g.config.Template
	if template == "" {
		template = "default"
	}
	db := g.config.DBDriver
	if db == "" {
		db = "postgres"
	}

	return scaffold.CreateProject(scaffold.Options{
		ProjectName:  g.config.Name,
		ProjectPath:  g.config.Path,
		Template:     template,
		DBDriver:     db,
		SkipBackend:  g.config.SkipBackend,
		SkipFrontend: g.config.SkipAngular,
	})
}
