# GoAstra

A production-ready full-stack framework combining Go backend with Angular frontend.

**Author:** [channdev](https://github.com/channdev)

## Features

- **Multiple API Types** - REST (Gin), GraphQL (gqlgen), tRPC (Connect-Go)
- **ORM Options** - SQLx for raw SQL or Ent ORM (Prisma-like experience)
- **Angular Frontend** - Modern TypeScript SPA with standalone components
- **Database Support** - MySQL & PostgreSQL with auto-detection from env vars
- **Database Migrations** - Native migration system with batch support
- **TypeSync** - Auto-generate TypeScript interfaces from Go structs
- **Code Generation** - Scaffold APIs, modules, CRUD, GraphQL, tRPC, and Ent schemas
- **JWT Authentication** - Built-in auth with refresh tokens
- **Stylish Logging** - Color-coded request logs with handler file info
- **Environment Management** - Development, production, and test configs

## Installation

### Prerequisites

- **Go 1.21+** - [Download Go](https://go.dev/dl/)
- **Node.js 18+** - [Download Node.js](https://nodejs.org/)
- **MySQL 8+** or **PostgreSQL 14+** (optional)

### Install via Go

```bash
go install github.com/channdev/goastra/cli/goastra@latest
```

### Install from Source

```bash
git clone https://github.com/channdev/goastra.git
cd goastra/cli
go build -o goastra.exe ./goastra

# Move to PATH (Linux/macOS)
sudo mv goastra /usr/local/bin/

# Or add to PATH (Windows)
# Move goastra.exe to a directory in your PATH
```

### Verify Installation

```bash
goastra version
goastra --help
```

## Quick Start

```bash
# Create a REST API project (default)
goastra new my-app --db mysql

# Create a GraphQL project with Ent ORM
goastra new my-app --api graphql --orm ent --db mysql

# Create a tRPC project
goastra new my-app --api trpc --db postgres

# Navigate and start development
cd my-app
goastra dev
```

Your app will be available at:
- Frontend: http://localhost:4200
- Backend: http://localhost:8080

## CLI Commands

| Command | Description |
|---------|-------------|
| `goastra new <name>` | Create a new project |
| `goastra dev` | Start development servers |
| `goastra build` | Build for production |
| `goastra start` | Start production server |
| `goastra migrate` | Run database migrations |
| `goastra generate` | Generate code (api, module, crud) |
| `goastra typesync` | Sync Go types to TypeScript |
| `goastra test` | Run test suites |

---

## Creating a New Project

```bash
goastra new <project-name> [flags]
```

| Flag | Default | Description |
|------|---------|-------------|
| `--api` | `rest` | API type (`rest`, `graphql`, `trpc`) |
| `--orm` | `sqlx` | ORM type (`sqlx`, `ent`) |
| `--db` | `postgres` | Database driver (`postgres`, `mysql`) |
| `-t, --template` | `default` | Project template (`default`, `minimal`) |
| `--skip-angular` | `false` | Skip Angular frontend generation |
| `--skip-backend` | `false` | Skip Go backend generation |

### API Types

| Type | Framework | Description |
|------|-----------|-------------|
| `rest` | Gin | Traditional REST API with JSON endpoints |
| `graphql` | gqlgen | Type-safe GraphQL with playground |
| `trpc` | Connect-Go | Type-safe RPC with Protocol Buffers |

### ORM Options

| Type | Description |
|------|-------------|
| `sqlx` | Raw SQL with type-safe query helpers |
| `ent` | Prisma-like ORM by Facebook with code generation |

**Examples:**

```bash
# REST API with SQLx (default)
goastra new my-app --db mysql

# GraphQL API with Ent ORM
goastra new my-app --api graphql --orm ent --db mysql

# tRPC API with SQLx
goastra new my-app --api trpc --db postgres

# REST with Ent ORM and minimal template
goastra new my-app --api rest --orm ent -t minimal

# Backend only (no Angular)
goastra new my-api --skip-angular
```

---

## Database Migrations

GoAstra includes a powerful migration system for managing database schema changes.

### Migration Commands

| Command | Description |
|---------|-------------|
| `goastra migrate` | Run all pending migrations |
| `goastra migrate:status` | Show migration status |
| `goastra migrate:rollback` | Rollback the last batch |
| `goastra migrate:reset` | Rollback all migrations |
| `goastra migrate:refresh` | Reset and re-run all migrations |
| `goastra migrate:fresh` | Drop all tables and re-run migrations |
| `goastra migrate:make <name>` | Create a new migration file |

### Database Configuration

GoAstra automatically loads database config from `.env.development`. Configure using simple environment variables:

**MySQL Configuration:**
```bash
MYSQL_HOST=localhost
MYSQL_USERNAME=root
MYSQL_PASSWORD=
MYSQL_DATABASE=myapp_dev
MYSQL_PORT=3306
```

**PostgreSQL Configuration:**
```bash
DB_URL=postgres://user:password@localhost:5432/myapp_dev?sslmode=disable
```

**Or use individual PostgreSQL vars:**
```bash
POSTGRES_HOST=localhost
POSTGRES_USER=postgres
POSTGRES_PASSWORD=secret
POSTGRES_DB=myapp_dev
POSTGRES_PORT=5432
```

### Creating Migrations

```bash
# Create a new migration
goastra migrate:make create_users_table

# Create with table template
goastra migrate:make create_products_table --create
```

This creates a file in `app/database/migrations/`:

```sql
-- GoAstra Migration
-- @up
CREATE TABLE users (
    id INT AUTO_INCREMENT PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    email VARCHAR(255) UNIQUE NOT NULL,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

-- @down
DROP TABLE IF EXISTS users;
```

### Running Migrations

```bash
# Run all pending migrations
goastra migrate

# Run specific number of migrations
goastra migrate --step=3

# Check migration status
goastra migrate:status

# Rollback last batch
goastra migrate:rollback

# Rollback specific number
goastra migrate:rollback --step=2

# Fresh start (drop all & re-run)
goastra migrate:fresh
```

---

## Development

### Start Development Servers

```bash
# Start both backend and frontend
goastra dev

# Backend only
goastra dev --backend

# Frontend only
goastra dev --frontend

# Custom ports
goastra dev -p 3000 --frontend-port 4300
```

### Request Logging

GoAstra provides stylish, color-coded request logs:

```
02:48:42 | 200 |   55.3ms | auth.go      POST    "/api/v1/auth/register"
02:48:56 | 201 |   12.1ms | auth.go      POST    "/api/v1/auth/login"
02:49:01 | 200 |    0.8ms | users.go     GET     "/api/v1/users"
02:49:05 | 404 |    0.2ms | users.go     GET     "/api/v1/users/123"
```

- **Green** status for 2xx responses
- **Yellow** status for 4xx responses
- **Red** status for 5xx responses
- Handler file name shown for easy debugging

---

## Code Generation

GoAstra provides powerful code generators for all API types and ORMs.

### Generate REST API Endpoint

```bash
goastra generate api product
```

Creates:
- `app/internal/handlers/product_handler.go`
- `app/internal/services/product_service.go`
- `app/internal/repository/product_repository.go`

### Generate GraphQL Schema & Resolvers

```bash
goastra generate graphql product
```

Creates:
- `app/graph/product.graphqls` - GraphQL schema with types and operations
- `app/graph/product.resolvers.go` - Resolver implementations

After generation:
```bash
cd app && go generate ./...
```

### Generate tRPC Proto & Service

```bash
goastra generate trpc product
```

Creates:
- `app/proto/v1/product.proto` - Protocol Buffer definitions
- `app/internal/rpc/product_service.go` - Connect-Go service

After generation:
```bash
cd app && buf generate
```

### Generate Ent Schema

```bash
goastra generate ent product
```

Creates:
- `app/ent/schema/product.go` - Ent entity schema

After generation:
```bash
cd app && go generate ./ent
```

### Generate Angular Module

```bash
goastra generate module product
```

Creates:
- `web/src/app/features/product/product.component.ts`
- `web/src/app/features/product/product.service.ts`

### Generate Full CRUD Stack

```bash
goastra generate crud product
```

Creates both backend API and frontend module with list, detail, create, edit, and delete operations.

---

## Type Synchronization

Define your types in `schema/types/`:

```go
type Product struct {
    ID    uint    `json:"id"`
    Name  string  `json:"name"`
    Price float64 `json:"price"`
}
```

Run:

```bash
goastra typesync
```

Generates TypeScript:

```typescript
export interface Product {
    id: number;
    name: string;
    price: number;
}
```

---

## Project Structure

### REST API Project (default)

```
my-app/
├── app/                      # Go backend
│   ├── cmd/server/           # Entry point
│   ├── internal/             # Internal packages
│   │   ├── auth/             # JWT authentication
│   │   ├── config/           # Configuration
│   │   ├── database/         # Database connection (SQLx)
│   │   ├── handlers/         # HTTP handlers
│   │   ├── logger/           # Structured logging (Zap)
│   │   ├── middleware/       # CORS, Auth, Logger
│   │   ├── models/           # Data models
│   │   ├── repository/       # Data access layer
│   │   ├── router/           # Route registration
│   │   ├── services/         # Business logic
│   │   └── validator/        # Request validation
│   └── migrations/           # SQL migration files
├── web/                      # Angular frontend
│   └── src/app/
│       ├── core/services/    # API & Auth services
│       └── features/         # Feature modules
└── .env.development          # Environment config
```

### GraphQL API Project

```
my-app/
├── app/
│   ├── cmd/server/           # gqlgen server entry
│   ├── graph/                # GraphQL layer
│   │   ├── schema.graphqls   # GraphQL schema
│   │   ├── resolver.go       # Root resolver
│   │   ├── generated/        # gqlgen generated code
│   │   └── model/            # GraphQL models
│   ├── gqlgen.yml            # gqlgen config
│   └── internal/...          # Shared packages
├── web/
│   └── src/app/core/services/
│       └── graphql.service.ts  # Apollo client
└── codegen.yml               # GraphQL codegen config
```

### tRPC API Project

```
my-app/
├── app/
│   ├── cmd/server/           # Connect-Go server
│   ├── proto/v1/             # Protocol Buffer definitions
│   │   └── service.proto
│   ├── internal/rpc/         # RPC implementations
│   │   ├── gen/              # Generated code
│   │   ├── service.go        # Service implementations
│   │   └── interceptor.go    # Logging/auth interceptors
│   ├── buf.yaml              # Buf config
│   └── buf.gen.yaml          # Code generation config
├── web/
│   └── src/app/core/services/
│       └── trpc.service.ts   # Connect-Web client
└── buf.gen.yaml              # Frontend code generation
```

### Ent ORM Project (any API type)

```
my-app/app/
├── ent/                      # Ent ORM
│   ├── schema/               # Entity schemas
│   │   └── user.go
│   └── generate.go           # go:generate directive
└── internal/database/
    └── database.go           # Ent client setup
```

---

## Configuration

### goastra.json

```json
{
  "name": "my-app",
  "version": "1.0.0",
  "api": {
    "type": "rest",
    "prefix": "/api/v1"
  },
  "backend": {
    "port": 8080
  },
  "frontend": {
    "port": 4200
  },
  "database": {
    "driver": "mysql",
    "migrationsPath": "app/database/migrations"
  }
}
```

### Environment Variables

```bash
# Application
APP_ENV=development
PORT=8080
LOG_LEVEL=debug

# MySQL
MYSQL_HOST=localhost
MYSQL_USERNAME=root
MYSQL_PASSWORD=
MYSQL_DATABASE=myapp_dev
MYSQL_PORT=3306

# JWT
JWT_SECRET=your-secret-key-min-32-chars
JWT_EXPIRY=24h

# CORS
CORS_ALLOWED_ORIGINS=http://localhost:4200
```

---

## Building for Production

```bash
# Build both backend and frontend
goastra build

# Start production server
goastra start
```

The build process:
1. Compiles Go backend to `./bin/server`
2. Builds Angular frontend to `./public/browser`
3. Backend serves static files in production mode

---

## License

MIT License - see [LICENSE](LICENSE) for details.

## Contributing

1. Fork the repository
2. Create your feature branch
3. Commit your changes
4. Push to the branch
5. Create a Pull Request

## Links

- [GitHub](https://github.com/channdev/goastra)
- [Issues](https://github.com/channdev/goastra/issues)
