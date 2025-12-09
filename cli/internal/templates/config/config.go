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
	if db == "mysql" {
		return `APP_ENV=development
API_URL=http://localhost:8080
PORT=8080
LOG_LEVEL=debug

# MySQL Configuration
MYSQL_HOST=localhost
MYSQL_USERNAME=root
MYSQL_PASSWORD=
MYSQL_DATABASE=goastra_dev
MYSQL_PORT=3306

JWT_SECRET=dev-secret-change-in-production-32chars
JWT_EXPIRY=24h
CORS_ALLOWED_ORIGINS=http://localhost:4200
`
	}

	return `APP_ENV=development
API_URL=http://localhost:8080
PORT=8080
LOG_LEVEL=debug

# PostgreSQL Configuration
DB_URL=postgres://user:password@localhost:5432/goastra_dev?sslmode=disable

JWT_SECRET=dev-secret-change-in-production-32chars
JWT_EXPIRY=24h
CORS_ALLOWED_ORIGINS=http://localhost:4200
`
}

func EnvProduction(db string) string {
	if db == "mysql" {
		return `APP_ENV=production
API_URL=https://api.example.com
PORT=8080
LOG_LEVEL=info

# MySQL Configuration
MYSQL_HOST=
MYSQL_USERNAME=
MYSQL_PASSWORD=
MYSQL_DATABASE=
MYSQL_PORT=3306

JWT_SECRET=
JWT_EXPIRY=24h
CORS_ALLOWED_ORIGINS=https://example.com
`
	}

	return `APP_ENV=production
API_URL=https://api.example.com
PORT=8080
LOG_LEVEL=info

# PostgreSQL Configuration
DB_URL=

JWT_SECRET=
JWT_EXPIRY=24h
CORS_ALLOWED_ORIGINS=https://example.com
`
}

func EnvTest(db string) string {
	if db == "mysql" {
		return `APP_ENV=test
API_URL=http://localhost:8081
PORT=8081
LOG_LEVEL=error

# MySQL Configuration
MYSQL_HOST=localhost
MYSQL_USERNAME=root
MYSQL_PASSWORD=
MYSQL_DATABASE=goastra_test
MYSQL_PORT=3306

JWT_SECRET=test-secret-32-characters-long!!
JWT_EXPIRY=1h
CORS_ALLOWED_ORIGINS=*
`
	}

	return `APP_ENV=test
API_URL=http://localhost:8081
PORT=8081
LOG_LEVEL=error

# PostgreSQL Configuration
DB_URL=postgres://user:password@localhost:5432/goastra_test?sslmode=disable

JWT_SECRET=test-secret-32-characters-long!!
JWT_EXPIRY=1h
CORS_ALLOWED_ORIGINS=*
`
}
