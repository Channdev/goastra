package scaffold

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/channdev/goastra/cli/internal/templates/backend"
	"github.com/channdev/goastra/cli/internal/templates/backend/common"
	"github.com/channdev/goastra/cli/internal/templates/backend/graphql"
	"github.com/channdev/goastra/cli/internal/templates/backend/rest"
	"github.com/channdev/goastra/cli/internal/templates/backend/trpc"
	"github.com/channdev/goastra/cli/internal/templates/config"
	"github.com/channdev/goastra/cli/internal/templates/frontend"
	"github.com/channdev/goastra/cli/internal/templates/frontend/api"
	entorm "github.com/channdev/goastra/cli/internal/templates/orm/ent"
	"github.com/channdev/goastra/cli/internal/templates/orm/sqlx"
	"github.com/fatih/color"
)

type Options struct {
	ProjectName  string
	ProjectPath  string
	Template     string
	DBDriver     string
	APIType      string // rest, graphql, trpc
	ORMType      string // sqlx, ent
	SkipBackend  bool
	SkipFrontend bool
}

func CreateProject(opts Options) error {
	// Set defaults
	if opts.APIType == "" {
		opts.APIType = "rest"
	}
	if opts.ORMType == "" {
		opts.ORMType = "sqlx"
	}

	color.Cyan("Creating new GoAstra project: %s\n", opts.ProjectName)
	color.Cyan("Template: %s | API: %s | ORM: %s | Database: %s\n\n", opts.Template, opts.APIType, opts.ORMType, opts.DBDriver)

	color.Yellow("[1/7] Creating project structure...\n")
	if err := createDirectories(opts.ProjectPath, opts.Template, opts.APIType, opts.ORMType); err != nil {
		return err
	}

	color.Yellow("[2/7] Generating configuration files...\n")
	if err := generateConfigFiles(opts.ProjectPath, opts.ProjectName); err != nil {
		return err
	}

	color.Yellow("[3/7] Generating environment files...\n")
	if err := generateEnvFiles(opts.ProjectPath, opts.DBDriver); err != nil {
		return err
	}

	if !opts.SkipBackend {
		color.Yellow("[4/7] Generating Go backend...\n")
		if err := generateBackend(opts.ProjectPath, opts.ProjectName, opts.DBDriver, opts.APIType, opts.ORMType); err != nil {
			return err
		}
	}

	if !opts.SkipFrontend {
		color.Yellow("[5/7] Generating Angular frontend...\n")
		if err := generateFrontend(opts.ProjectPath, opts.ProjectName, opts.Template, opts.APIType); err != nil {
			return err
		}
	}

	color.Yellow("[6/7] Generating schema types...\n")
	if err := generateSchema(opts.ProjectPath); err != nil {
		return err
	}

	color.Yellow("[7/7] Installing dependencies...\n")
	if err := installDependencies(opts.ProjectPath, opts.SkipBackend, opts.SkipFrontend); err != nil {
		color.Yellow("Warning: Failed to install some dependencies: %v\n", err)
		color.Yellow("You may need to run 'go mod tidy' in app/ and 'npm install' in web/ manually.\n")
	}

	color.Green("\nProject created successfully!\n\n")
	fmt.Printf("Next steps:\n")
	fmt.Printf("  cd %s\n", opts.ProjectName)
	fmt.Printf("  goastra dev\n\n")
	fmt.Printf("Your app will be available at:\n")
	fmt.Printf("  Frontend: http://localhost:4200\n")
	fmt.Printf("  Backend:  http://localhost:8080\n")

	return nil
}

func createDirectories(projectPath, template, apiType, ormType string) error {
	// Base directories for all projects
	dirs := []string{
		"app/cmd/server",
		"app/internal/config",
		"app/internal/middleware",
		"app/internal/models",
		"app/internal/auth",
		"app/internal/logger",
		"app/internal/database",
		"app/migrations",
		"web/src/environments",
		"web/src/assets",
		"schema/types",
	}

	// API-type specific directories
	switch apiType {
	case "graphql":
		dirs = append(dirs,
			"app/graph",
			"app/graph/model",
			"app/graph/generated",
		)
	case "trpc":
		dirs = append(dirs,
			"app/proto/v1",
			"app/internal/rpc",
			"app/internal/rpc/gen/proto/v1",
		)
	default: // rest
		dirs = append(dirs,
			"app/internal/handlers",
			"app/internal/repository",
			"app/internal/services",
			"app/internal/router",
			"app/internal/validator",
		)
	}

	// ORM-type specific directories
	switch ormType {
	case "ent":
		dirs = append(dirs,
			"app/ent",
			"app/ent/schema",
		)
	default: // sqlx
		dirs = append(dirs,
			"app/internal/repository",
		)
	}

	// Frontend template directories
	if template == "default" {
		dirs = append(dirs,
			"web/src/app/core/services",
			"web/src/app/core/guards",
			"web/src/app/core/interceptors",
			"web/src/app/core/models",
			"web/src/app/shared/components",
			"web/src/app/shared/directives",
			"web/src/app/shared/pipes",
			"web/src/app/features/home",
			"web/src/app/features/auth/login",
			"web/src/app/features/auth/register",
			"web/src/app/features/dashboard",
			"web/src/app/features/not-found",
		)
	} else {
		dirs = append(dirs, "web/src/app/home")
	}

	for _, dir := range dirs {
		fullPath := filepath.Join(projectPath, dir)
		if err := os.MkdirAll(fullPath, 0755); err != nil {
			return fmt.Errorf("failed to create directory %s: %w", dir, err)
		}
	}
	return nil
}

func generateConfigFiles(projectPath, projectName string) error {
	files := map[string]string{
		"goastra.json": config.GoastraJSON(projectName),
		".gitignore":   config.Gitignore(),
	}

	for name, content := range files {
		if err := writeFile(projectPath, name, content); err != nil {
			return err
		}
	}
	return nil
}

func generateEnvFiles(projectPath, db string) error {
	files := map[string]string{
		".env.development": config.EnvDevelopment(db),
		".env.production":  config.EnvProduction(db),
		".env.test":        config.EnvTest(db),
	}

	for name, content := range files {
		if err := writeFile(projectPath, name, content); err != nil {
			return err
		}
	}
	return nil
}

func generateBackend(projectPath, projectName, db, apiType, ormType string) error {
	files := make(map[string]string)

	// Common files for all API types
	files["app/internal/config/config.go"] = common.ConfigGo()
	files["app/internal/logger/logger.go"] = common.LoggerGo()
	files["app/internal/auth/auth.go"] = common.AuthGo()
	files["app/internal/middleware/middleware.go"] = common.MiddlewareGo()
	files["app/internal/models/models.go"] = common.ModelsGo()

	// Generate go.mod based on API and ORM type
	files["app/go.mod"] = common.GoMod(common.GoModOptions{
		ProjectName: projectName,
		DBDriver:    db,
		APIType:     apiType,
		ORMType:     ormType,
	})

	// ORM-specific files
	switch ormType {
	case "ent":
		files["app/internal/database/database.go"] = entorm.ClientGo(db)
		files["app/ent/generate.go"] = entorm.GenerateGo()
		files["app/ent/schema/user.go"] = entorm.UserSchemaGo()
		files["app/ent/schema/mixin.go"] = entorm.BaseMixinGo()
	default: // sqlx
		files["app/internal/database/database.go"] = sqlx.DatabaseGo(db)
		files["app/internal/repository/repository.go"] = sqlx.RepositoryGo()
	}

	// API-specific files
	switch apiType {
	case "graphql":
		files["app/cmd/server/main.go"] = graphql.MainGo()
		files["app/gqlgen.yml"] = graphql.GqlgenYML()
		files["app/graph/generate.go"] = graphql.GenerateGo()
		files["app/graph/schema.graphqls"] = graphql.SchemaGraphQL()
		files["app/graph/resolver.go"] = graphql.ResolverGo()
		files["app/graph/schema.resolvers.go"] = graphql.SchemaResolversGo()
		files["app/tools.go"] = graphql.ToolsGo()

	case "trpc":
		files["app/cmd/server/main.go"] = trpc.MainGo()
		files["app/proto/v1/service.proto"] = trpc.ServiceProto()
		files["app/buf.yaml"] = trpc.BufYAML()
		files["app/buf.gen.yaml"] = trpc.BufGenYAML()
		files["app/buf.work.yaml"] = trpc.BufWorkYAML()
		files["app/internal/rpc/service.go"] = trpc.ServiceGo()
		files["app/internal/rpc/interceptor.go"] = trpc.InterceptorGo()

	default: // rest
		files["app/cmd/server/main.go"] = rest.MainGo()
		files["app/internal/router/router.go"] = rest.RouterGo()
		files["app/internal/handlers/handlers.go"] = rest.HandlersGo()
		files["app/internal/services/services.go"] = rest.ServicesGo()
		files["app/internal/validator/validator.go"] = common.ValidatorGo()
	}

	for path, content := range files {
		if err := writeFile(projectPath, path, content); err != nil {
			return err
		}
	}
	return nil
}

func generateFrontend(projectPath, projectName, template, apiType string) error {
	files := make(map[string]string)

	// API-specific package.json and config
	switch apiType {
	case "graphql":
		files["web/package.json"] = api.GraphQLPackageJSON(projectName)
		files["web/src/app/app.config.ts"] = api.ApolloConfigTS()
		files["web/codegen.yml"] = api.CodegenYML()
		files["web/src/environments/environment.ts"] = api.GraphQLEnvTS()
		files["web/src/environments/environment.prod.ts"] = api.GraphQLEnvProdTS()
		files["web/src/app/core/services/graphql.service.ts"] = api.GraphQLServiceTS()
	case "trpc":
		files["web/package.json"] = api.TRPCPackageJSON(projectName)
		files["web/src/app/app.config.ts"] = api.TRPCAppConfigTS()
		files["web/buf.gen.yaml"] = api.BufGenYAMLWeb()
		files["web/src/environments/environment.ts"] = api.TRPCEnvTS()
		files["web/src/environments/environment.prod.ts"] = api.TRPCEnvProdTS()
		files["web/src/app/core/services/trpc.service.ts"] = api.TRPCServiceTS()
	default: // rest
		files["web/package.json"] = frontend.PackageJSON(projectName)
		files["web/src/app/app.config.ts"] = frontend.AppConfig()
		files["web/src/environments/environment.ts"] = frontend.EnvDev()
		files["web/src/environments/environment.prod.ts"] = frontend.EnvProd()
		files["web/src/app/core/services/api.service.ts"] = api.RESTServiceTS()
		files["web/src/app/core/interceptors/auth.interceptor.ts"] = api.AuthInterceptorTS()
		files["web/src/app/core/services/auth.service.ts"] = api.AuthServiceTS()
	}

	// Common frontend files
	files["web/angular.json"] = frontend.AngularJSON(projectName)
	files["web/tsconfig.app.json"] = frontend.TSConfigApp()
	files["web/proxy.conf.json"] = frontend.ProxyConf()
	files["web/src/index.html"] = frontend.IndexHTML(projectName)
	files["web/src/main.ts"] = frontend.MainTS()
	files["web/src/app/app.component.ts"] = frontend.AppComponent()

	// Template-specific files
	if template == "minimal" {
		files["web/tsconfig.json"] = frontend.TSConfigSimple()
		files["web/src/styles.css"] = frontend.MinimalStylesCSS()
		files["web/src/app/app.routes.ts"] = frontend.MinimalAppRoutes()
		files["web/src/app/home/home.component.ts"] = frontend.MinimalHomeComponent()
	} else {
		files["web/tsconfig.json"] = frontend.TSConfigWithPaths()
		files["web/src/styles.css"] = frontend.DefaultStylesCSS()
		files["web/src/app/app.routes.ts"] = frontend.DefaultAppRoutes()
		files["web/src/app/features/home/home.component.ts"] = frontend.DefaultHomeComponent()
		files["web/src/app/features/auth/login/login.component.ts"] = frontend.DefaultLoginComponent()
		files["web/src/app/features/auth/register/register.component.ts"] = frontend.DefaultRegisterComponent()
		files["web/src/app/features/dashboard/dashboard.component.ts"] = frontend.DefaultDashboardComponent()
	}

	for path, content := range files {
		if err := writeFile(projectPath, path, content); err != nil {
			return err
		}
	}
	return nil
}

func generateSchema(projectPath string) error {
	files := map[string]string{
		"schema/types/types.go": backend.SchemaTypesGo(),
		"schema/go.mod":         backend.SchemaGoMod(),
	}

	for path, content := range files {
		if err := writeFile(projectPath, path, content); err != nil {
			return err
		}
	}
	return nil
}

func installDependencies(projectPath string, skipBackend, skipFrontend bool) error {
	if !skipBackend {
		color.Blue("  Running 'go mod tidy' in app/...\n")
		appDir := filepath.Join(projectPath, "app")
		cmd := exec.Command("go", "mod", "tidy")
		cmd.Dir = appDir
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		if err := cmd.Run(); err != nil {
			return fmt.Errorf("go mod tidy failed: %w", err)
		}
	}

	if !skipFrontend {
		color.Blue("  Running 'npm install' in web/...\n")
		webDir := filepath.Join(projectPath, "web")
		var cmd *exec.Cmd
		if os.PathSeparator == '\\' && os.PathListSeparator == ';' {
			cmd = exec.Command("cmd", "/c", "npm", "install")
		} else {
			cmd = exec.Command("npm", "install")
		}
		cmd.Dir = webDir
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		if err := cmd.Run(); err != nil {
			return fmt.Errorf("npm install failed: %w", err)
		}
	}

	return nil
}

func writeFile(basePath, relativePath, content string) error {
	fullPath := filepath.Join(basePath, relativePath)
	return os.WriteFile(fullPath, []byte(content), 0644)
}
