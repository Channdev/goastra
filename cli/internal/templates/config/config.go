package config

import "fmt"

func GoastraJSON(projectName string) string {
	return fmt.Sprintf(`{
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
}

func Gitignore() string {
	return `dist/
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
}

func EnvDevelopment(db string) string {
	var dbURL string
	if db == "mysql" {
		dbURL = "user:password@tcp(localhost:3306)/goastra_dev?parseTime=true"
	} else {
		dbURL = "postgres://user:password@localhost:5432/goastra_dev?sslmode=disable"
	}

	return fmt.Sprintf(`APP_ENV=development
API_URL=http://localhost:8080
PORT=8080
LOG_LEVEL=debug
DB_DRIVER=%s
DB_URL=%s
JWT_SECRET=dev-secret-change-in-production-32chars
JWT_EXPIRY=24h
CORS_ALLOWED_ORIGINS=http://localhost:4200
`, db, dbURL)
}

func EnvProduction(db string) string {
	return fmt.Sprintf(`APP_ENV=production
API_URL=https://api.example.com
PORT=8080
LOG_LEVEL=info
DB_DRIVER=%s
DB_URL=
JWT_SECRET=
JWT_EXPIRY=24h
CORS_ALLOWED_ORIGINS=https://example.com
`, db)
}

func EnvTest(db string) string {
	var dbURL string
	if db == "mysql" {
		dbURL = "user:password@tcp(localhost:3306)/goastra_test?parseTime=true"
	} else {
		dbURL = "postgres://user:password@localhost:5432/goastra_test?sslmode=disable"
	}

	return fmt.Sprintf(`APP_ENV=test
API_URL=http://localhost:8081
PORT=8081
LOG_LEVEL=error
DB_DRIVER=%s
DB_URL=%s
JWT_SECRET=test-secret-32-characters-long!!
JWT_EXPIRY=1h
CORS_ALLOWED_ORIGINS=*
`, db, dbURL)
}
