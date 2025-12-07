package scaffold

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/channdev/goastra/cli/internal/templates/backend"
	"github.com/channdev/goastra/cli/internal/templates/config"
	"github.com/channdev/goastra/cli/internal/templates/frontend"
	"github.com/fatih/color"
)

type Options struct {
	ProjectName  string
	ProjectPath  string
	Template     string
	DBDriver     string
	SkipBackend  bool
	SkipFrontend bool
}

func CreateProject(opts Options) error {
	color.Cyan("Creating new GoAstra project: %s\n", opts.ProjectName)
	color.Cyan("Template: %s | Database: %s\n\n", opts.Template, opts.DBDriver)

	color.Yellow("[1/7] Creating project structure...\n")
	if err := createDirectories(opts.ProjectPath, opts.Template); err != nil {
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
		if err := generateBackend(opts.ProjectPath, opts.ProjectName, opts.DBDriver); err != nil {
			return err
		}
	}

	if !opts.SkipFrontend {
		color.Yellow("[5/7] Generating Angular frontend...\n")
		if err := generateFrontend(opts.ProjectPath, opts.ProjectName, opts.Template); err != nil {
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

func createDirectories(projectPath, template string) error {
	dirs := []string{
		"app/cmd/server",
		"app/internal/config",
		"app/internal/middleware",
		"app/internal/handlers",
		"app/internal/models",
		"app/internal/repository",
		"app/internal/services",
		"app/internal/router",
		"app/internal/auth",
		"app/internal/logger",
		"app/internal/database",
		"app/internal/validator",
		"app/migrations",
		"web/src/environments",
		"web/src/assets",
		"schema/types",
	}

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

func generateBackend(projectPath, projectName, db string) error {
	files := map[string]string{
		"app/go.mod":                             backend.GoMod(projectName, db),
		"app/cmd/server/main.go":                 backend.MainGo(),
		"app/internal/config/config.go":         backend.ConfigGo(),
		"app/internal/logger/logger.go":         backend.LoggerGo(),
		"app/internal/database/database.go":     backend.DatabaseGo(db),
		"app/internal/auth/auth.go":             backend.AuthGo(),
		"app/internal/middleware/middleware.go": backend.MiddlewareGo(),
		"app/internal/models/models.go":         backend.ModelsGo(),
		"app/internal/handlers/handlers.go":     backend.HandlersGo(),
		"app/internal/repository/repository.go": backend.RepositoryGo(),
		"app/internal/services/services.go":     backend.ServicesGo(),
		"app/internal/router/router.go":         backend.RouterGo(),
		"app/internal/validator/validator.go":   backend.ValidatorGo(),
	}

	for path, content := range files {
		if err := writeFile(projectPath, path, content); err != nil {
			return err
		}
	}
	return nil
}

func generateFrontend(projectPath, projectName, template string) error {
	files := map[string]string{
		"web/package.json":                         frontend.PackageJSON(projectName),
		"web/angular.json":                         frontend.AngularJSON(projectName),
		"web/tsconfig.app.json":                    frontend.TSConfigApp(),
		"web/proxy.conf.json":                      frontend.ProxyConf(),
		"web/src/index.html":                       frontend.IndexHTML(projectName),
		"web/src/main.ts":                          frontend.MainTS(),
		"web/src/app/app.component.ts":             frontend.AppComponent(),
		"web/src/app/app.config.ts":                frontend.AppConfig(),
		"web/src/environments/environment.ts":      frontend.EnvDev(),
		"web/src/environments/environment.prod.ts": frontend.EnvProd(),
	}

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
