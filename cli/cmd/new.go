package cmd

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/fatih/color"
	"github.com/spf13/cobra"
)

var (
	skipAngular bool
	skipBackend bool
	useGraphQL  bool
)

var newCmd = &cobra.Command{
	Use:   "new <project-name>",
	Short: "Create a new GoAstra project",
	Long:  "Creates a new GoAstra project with Go backend and Angular frontend",
	Args:  cobra.ExactArgs(1),
	RunE:  runNew,
}

func init() {
	rootCmd.AddCommand(newCmd)
	newCmd.Flags().BoolVar(&skipAngular, "skip-angular", false, "skip Angular frontend generation")
	newCmd.Flags().BoolVar(&skipBackend, "skip-backend", false, "skip Go backend generation")
	newCmd.Flags().BoolVar(&useGraphQL, "graphql", false, "use GraphQL instead of REST")
}

func runNew(cmd *cobra.Command, args []string) error {
	projectName := args[0]

	if err := validateProjectName(projectName); err != nil {
		return err
	}

	projectPath, err := filepath.Abs(projectName)
	if err != nil {
		return fmt.Errorf("failed to resolve project path: %w", err)
	}

	if _, err := os.Stat(projectPath); !os.IsNotExist(err) {
		return fmt.Errorf("directory already exists: %s", projectPath)
	}

	color.Cyan("Creating new GoAstra project: %s\n\n", projectName)

	color.Yellow("[1/7] Creating project structure...\n")
	if err := createDirectories(projectPath); err != nil {
		return err
	}

	color.Yellow("[2/7] Generating configuration files...\n")
	if err := generateConfigFiles(projectPath, projectName); err != nil {
		return err
	}

	color.Yellow("[3/7] Generating environment files...\n")
	if err := generateEnvFiles(projectPath); err != nil {
		return err
	}

	if !skipBackend {
		color.Yellow("[4/7] Generating Go backend...\n")
		if err := generateBackend(projectPath, projectName); err != nil {
			return err
		}
	}

	if !skipAngular {
		color.Yellow("[5/7] Generating Angular frontend...\n")
		if err := generateFrontend(projectPath, projectName); err != nil {
			return err
		}
	}

	color.Yellow("[6/7] Generating schema types...\n")
	if err := generateSchema(projectPath); err != nil {
		return err
	}

	color.Yellow("[7/7] Installing dependencies...\n")
	if err := installDependencies(projectPath, skipBackend, skipAngular); err != nil {
		color.Yellow("Warning: Failed to install some dependencies: %v\n", err)
		color.Yellow("You may need to run 'go mod tidy' in app/ and 'npm install' in web/ manually.\n")
	}

	color.Green("\nProject created successfully!\n\n")
	fmt.Printf("Next steps:\n")
	fmt.Printf("  cd %s\n", projectName)
	fmt.Printf("  goastra dev\n\n")
	fmt.Printf("Your app will be available at:\n")
	fmt.Printf("  Frontend: http://localhost:4200\n")
	fmt.Printf("  Backend:  http://localhost:8080\n")

	return nil
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

func createDirectories(projectPath string) error {
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
		"web/src/environments",
		"web/src/assets",
		"schema/types",
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
	goastraJSON := fmt.Sprintf(`{
  "name": "%s",
  "version": "1.0.0",
  "api": {
    "type": "rest",
    "prefix": "/api/v1"
  },
  "backend": {
    "port": 8080,
    "module": "github.com/%s/app"
  },
  "frontend": {
    "port": 4200,
    "proxy": "/api"
  },
  "codegen": {
    "schemaPath": "schema/types",
    "outputPath": "web/src/app/core/models"
  },
  "database": {
    "driver": "postgres",
    "migrationsPath": "app/migrations"
  }
}`, projectName, projectName)

	gitignore := `dist/
bin/
node_modules/
vendor/
.env
.env.local
.env.*.local
.idea/
.vscode/
*.exe
*.dll
*.so
*.dylib
web/dist/
web/.angular/
coverage/
*.log
tmp/
`

	if err := os.WriteFile(filepath.Join(projectPath, "goastra.json"), []byte(goastraJSON), 0644); err != nil {
		return err
	}
	return os.WriteFile(filepath.Join(projectPath, ".gitignore"), []byte(gitignore), 0644)
}

func generateEnvFiles(projectPath string) error {
	envDev := `APP_ENV=development
API_URL=http://localhost:8080
PORT=8080
LOG_LEVEL=debug
DB_URL=postgres://user:password@localhost:5432/goastra_dev?sslmode=disable
JWT_SECRET=dev-secret-change-in-production-32chars
JWT_EXPIRY=24h
CORS_ALLOWED_ORIGINS=http://localhost:4200
`
	envProd := `APP_ENV=production
API_URL=https://api.example.com
PORT=8080
LOG_LEVEL=info
DB_URL=
JWT_SECRET=
JWT_EXPIRY=24h
CORS_ALLOWED_ORIGINS=https://example.com
`
	envTest := `APP_ENV=test
API_URL=http://localhost:8081
PORT=8081
LOG_LEVEL=error
DB_URL=postgres://user:password@localhost:5432/goastra_test?sslmode=disable
JWT_SECRET=test-secret-32-characters-long!!
JWT_EXPIRY=1h
CORS_ALLOWED_ORIGINS=*
`
	if err := os.WriteFile(filepath.Join(projectPath, ".env.development"), []byte(envDev), 0644); err != nil {
		return err
	}
	if err := os.WriteFile(filepath.Join(projectPath, ".env.production"), []byte(envProd), 0644); err != nil {
		return err
	}
	return os.WriteFile(filepath.Join(projectPath, ".env.test"), []byte(envTest), 0644)
}

func generateBackend(projectPath, projectName string) error {
	goMod := fmt.Sprintf(`module github.com/%s/app

go 1.21

require (
	github.com/gin-gonic/gin v1.9.1
	github.com/go-playground/validator/v10 v10.16.0
	github.com/golang-jwt/jwt/v5 v5.2.0
	github.com/jmoiron/sqlx v1.3.5
	github.com/joho/godotenv v1.5.1
	github.com/lib/pq v1.10.9
	go.uber.org/zap v1.26.0
	golang.org/x/crypto v0.16.0
)
`, projectName)

	mainGo := `package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	env := os.Getenv("APP_ENV")
	if env == "" {
		env = "development"
	}

	godotenv.Load("../../.env." + env)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	if os.Getenv("APP_ENV") == "production" {
		gin.SetMode(gin.ReleaseMode)
	}

	r := gin.Default()

	r.Use(corsMiddleware())

	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "healthy", "version": "1.0.0"})
	})

	v1 := r.Group("/api/v1")
	{
		auth := v1.Group("/auth")
		{
			auth.POST("/login", handleLogin)
			auth.POST("/register", handleRegister)
			auth.POST("/refresh", handleRefresh)
			auth.POST("/logout", handleLogout)
		}

		users := v1.Group("/users")
		{
			users.GET("", handleListUsers)
			users.GET("/:id", handleGetUser)
			users.PUT("/:id", handleUpdateUser)
			users.DELETE("/:id", handleDeleteUser)
		}
	}

	srv := &http.Server{
		Addr:         ":" + port,
		Handler:      r,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
	}

	go func() {
		log.Printf("Server starting on port %s", port)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Server failed: %v", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("Shutting down server...")
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	srv.Shutdown(ctx)
}

func corsMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}
		c.Next()
	}
}

func handleLogin(c *gin.Context) {
	c.JSON(200, gin.H{"message": "Login endpoint - implement me"})
}

func handleRegister(c *gin.Context) {
	c.JSON(200, gin.H{"message": "Register endpoint - implement me"})
}

func handleRefresh(c *gin.Context) {
	c.JSON(200, gin.H{"message": "Refresh endpoint - implement me"})
}

func handleLogout(c *gin.Context) {
	c.JSON(200, gin.H{"message": "Logout endpoint - implement me"})
}

func handleListUsers(c *gin.Context) {
	c.JSON(200, gin.H{"data": []interface{}{}, "total": 0})
}

func handleGetUser(c *gin.Context) {
	c.JSON(200, gin.H{"id": c.Param("id")})
}

func handleUpdateUser(c *gin.Context) {
	c.JSON(200, gin.H{"id": c.Param("id"), "updated": true})
}

func handleDeleteUser(c *gin.Context) {
	c.JSON(200, gin.H{"id": c.Param("id"), "deleted": true})
}
`

	if err := os.WriteFile(filepath.Join(projectPath, "app/go.mod"), []byte(goMod), 0644); err != nil {
		return err
	}
	return os.WriteFile(filepath.Join(projectPath, "app/cmd/server/main.go"), []byte(mainGo), 0644)
}

func generateFrontend(projectPath, projectName string) error {
	packageJSON := fmt.Sprintf(`{
  "name": "%s-web",
  "version": "1.0.0",
  "scripts": {
    "ng": "ng",
    "start": "ng serve --proxy-config proxy.conf.json",
    "build": "ng build",
    "test": "ng test"
  },
  "dependencies": {
    "@angular/animations": "^17.0.0",
    "@angular/common": "^17.0.0",
    "@angular/compiler": "^17.0.0",
    "@angular/core": "^17.0.0",
    "@angular/forms": "^17.0.0",
    "@angular/platform-browser": "^17.0.0",
    "@angular/platform-browser-dynamic": "^17.0.0",
    "@angular/router": "^17.0.0",
    "rxjs": "~7.8.0",
    "tslib": "^2.6.0",
    "zone.js": "~0.14.0"
  },
  "devDependencies": {
    "@angular-devkit/build-angular": "^17.0.0",
    "@angular/cli": "^17.0.0",
    "@angular/compiler-cli": "^17.0.0",
    "typescript": "~5.2.0"
  }
}`, projectName)

	angularJSON := fmt.Sprintf(`{
  "$schema": "./node_modules/@angular/cli/lib/config/schema.json",
  "version": 1,
  "cli": { "analytics": false },
  "newProjectRoot": "projects",
  "projects": {
    "%s": {
      "projectType": "application",
      "root": "",
      "sourceRoot": "src",
      "prefix": "app",
      "architect": {
        "build": {
          "builder": "@angular-devkit/build-angular:application",
          "options": {
            "outputPath": "dist",
            "index": "src/index.html",
            "browser": "src/main.ts",
            "polyfills": ["zone.js"],
            "tsConfig": "tsconfig.app.json",
            "assets": ["src/assets"],
            "styles": ["src/styles.css"]
          },
          "configurations": {
            "production": {
              "outputHashing": "all",
              "fileReplacements": [{
                "replace": "src/environments/environment.ts",
                "with": "src/environments/environment.prod.ts"
              }]
            },
            "development": {
              "optimization": false,
              "sourceMap": true
            }
          },
          "defaultConfiguration": "production"
        },
        "serve": {
          "builder": "@angular-devkit/build-angular:dev-server",
          "configurations": {
            "production": { "buildTarget": "%s:build:production" },
            "development": { "buildTarget": "%s:build:development" }
          },
          "defaultConfiguration": "development"
        }
      }
    }
  }
}`, projectName, projectName, projectName)

	tsconfig := `{
  "compileOnSave": false,
  "compilerOptions": {
    "outDir": "./dist/out-tsc",
    "strict": true,
    "noImplicitOverride": true,
    "noPropertyAccessFromIndexSignature": true,
    "noImplicitReturns": true,
    "noFallthroughCasesInSwitch": true,
    "skipLibCheck": true,
    "esModuleInterop": true,
    "sourceMap": true,
    "declaration": false,
    "experimentalDecorators": true,
    "moduleResolution": "bundler",
    "importHelpers": true,
    "target": "ES2022",
    "module": "ES2022",
    "lib": ["ES2022", "dom"],
    "baseUrl": "./src",
    "paths": {
      "@core/*": ["app/core/*"],
      "@shared/*": ["app/shared/*"],
      "@features/*": ["app/features/*"],
      "@env/*": ["environments/*"]
    }
  },
  "angularCompilerOptions": {
    "enableI18nLegacyMessageIdFormat": false,
    "strictInjectionParameters": true,
    "strictInputAccessModifiers": true,
    "strictTemplates": true
  }
}`

	tsconfigApp := `{
  "extends": "./tsconfig.json",
  "compilerOptions": { "outDir": "./out-tsc/app" },
  "files": ["src/main.ts"],
  "include": ["src/**/*.d.ts"]
}`

	proxyConf := `{
  "/api": {
    "target": "http://localhost:8080",
    "secure": false,
    "changeOrigin": true
  }
}`

	indexHTML := fmt.Sprintf(`<!DOCTYPE html>
<html lang="en">
<head>
  <meta charset="utf-8">
  <title>%s</title>
  <base href="/">
  <meta name="viewport" content="width=device-width, initial-scale=1">
</head>
<body>
  <app-root></app-root>
</body>
</html>`, projectName)

	mainTS := `import { bootstrapApplication } from '@angular/platform-browser';
import { AppComponent } from './app/app.component';
import { appConfig } from './app/app.config';

bootstrapApplication(AppComponent, appConfig).catch((err) => console.error(err));
`

	stylesCSS := `:root {
  --color-primary: #3b82f6;
  --color-background: #0f172a;
  --color-surface: #1e293b;
  --color-text: #f8fafc;
  --color-text-muted: #94a3b8;
  --color-border: #334155;
}

* { box-sizing: border-box; margin: 0; padding: 0; }

body {
  font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, sans-serif;
  background: var(--color-background);
  color: var(--color-text);
}

a { color: var(--color-primary); text-decoration: none; }
`

	appComponent := `import { Component } from '@angular/core';
import { RouterOutlet } from '@angular/router';

@Component({
  selector: 'app-root',
  standalone: true,
  imports: [RouterOutlet],
  template: '<router-outlet></router-outlet>'
})
export class AppComponent {}
`

	appConfig := `import { ApplicationConfig } from '@angular/core';
import { provideRouter } from '@angular/router';
import { provideHttpClient } from '@angular/common/http';
import { routes } from './app.routes';

export const appConfig: ApplicationConfig = {
  providers: [provideRouter(routes), provideHttpClient()]
};
`

	appRoutes := `import { Routes } from '@angular/router';

export const routes: Routes = [
  { path: '', redirectTo: 'home', pathMatch: 'full' },
  { path: 'home', loadComponent: () => import('@features/home/home.component').then(m => m.HomeComponent) },
  { path: '**', redirectTo: 'home' }
];
`

	homeComponent := "import { Component } from '@angular/core';\n\n@Component({\n  selector: 'app-home',\n  standalone: true,\n  template: `\n    <div class=\"container\">\n      <h1>Welcome to GoAstra</h1>\n      <p>Your full-stack Go + Angular application is ready!</p>\n      <div class=\"links\">\n        <a href=\"http://localhost:8080/health\" target=\"_blank\">Backend Health</a>\n        <a href=\"https://github.com/channdev/goastra\" target=\"_blank\">Documentation</a>\n      </div>\n    </div>\n  `,\n  styles: [`\n    .container { min-height: 100vh; display: flex; flex-direction: column; align-items: center; justify-content: center; text-align: center; padding: 2rem; }\n    h1 { font-size: 3rem; margin-bottom: 1rem; background: linear-gradient(135deg, #3b82f6, #8b5cf6); -webkit-background-clip: text; -webkit-text-fill-color: transparent; }\n    p { color: #94a3b8; margin-bottom: 2rem; }\n    .links { display: flex; gap: 1rem; }\n    .links a { padding: 0.75rem 1.5rem; background: #3b82f6; color: white; border-radius: 8px; }\n  `]\n})\nexport class HomeComponent {}\n"

	envDev := `export const environment = { production: false, apiUrl: 'http://localhost:8080/api/v1' };`
	envProd := `export const environment = { production: true, apiUrl: '/api/v1' };`

	files := map[string]string{
		"web/package.json":                         packageJSON,
		"web/angular.json":                         angularJSON,
		"web/tsconfig.json":                        tsconfig,
		"web/tsconfig.app.json":                    tsconfigApp,
		"web/proxy.conf.json":                      proxyConf,
		"web/src/index.html":                       indexHTML,
		"web/src/main.ts":                          mainTS,
		"web/src/styles.css":                       stylesCSS,
		"web/src/app/app.component.ts":             appComponent,
		"web/src/app/app.config.ts":                appConfig,
		"web/src/app/app.routes.ts":                appRoutes,
		"web/src/app/features/home/home.component.ts": homeComponent,
		"web/src/environments/environment.ts":      envDev,
		"web/src/environments/environment.prod.ts": envProd,
	}

	for path, content := range files {
		fullPath := filepath.Join(projectPath, path)
		if err := os.WriteFile(fullPath, []byte(content), 0644); err != nil {
			return err
		}
	}

	return nil
}

func generateSchema(projectPath string) error {
	schemaGo := `package types

import "time"

type BaseModel struct {
	ID        uint      ` + "`json:\"id\"`" + `
	CreatedAt time.Time ` + "`json:\"created_at\"`" + `
	UpdatedAt time.Time ` + "`json:\"updated_at\"`" + `
}

type User struct {
	BaseModel
	Email  string ` + "`json:\"email\"`" + `
	Name   string ` + "`json:\"name\"`" + `
	Role   string ` + "`json:\"role\"`" + `
	Active bool   ` + "`json:\"active\"`" + `
}

type APIError struct {
	Code    string ` + "`json:\"code\"`" + `
	Message string ` + "`json:\"message\"`" + `
}
`

	goMod := `module schema

go 1.21
`

	if err := os.WriteFile(filepath.Join(projectPath, "schema/types/types.go"), []byte(schemaGo), 0644); err != nil {
		return err
	}
	return os.WriteFile(filepath.Join(projectPath, "schema/go.mod"), []byte(goMod), 0644)
}

func runCommand(dir, name string, args ...string) error {
	cmd := exec.Command(name, args...)
	cmd.Dir = dir
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

func installDependencies(projectPath string, skipBackend, skipAngular bool) error {
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

	if !skipAngular {
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
