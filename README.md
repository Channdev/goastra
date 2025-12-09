# GoAstra

A production-ready full-stack framework combining Go backend with Angular frontend.

**Author:** [channdev](https://github.com/channdev)

## Features

- **Go Backend** - High-performance REST API with Gin framework
- **Angular Frontend** - Modern TypeScript SPA with standalone components
- **Database Migrations** - Native migration system with MySQL & PostgreSQL support
- **TypeSync** - Auto-generate TypeScript interfaces from Go structs
- **Code Generation** - Scaffold CRUD operations with a single command
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
# Create a new project with MySQL
goastra new my-app --db mysql

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
| `-t, --template` | `default` | Project template (`default`, `minimal`) |
| `--db` | `postgres` | Database driver (`postgres`, `mysql`) |
| `--skip-angular` | `false` | Skip Angular frontend generation |
| `--skip-backend` | `false` | Skip Go backend generation |

**Examples:**

```bash
# Create with MySQL (recommended)
goastra new my-app --db mysql

# Create with minimal template
goastra new my-app -t minimal

# Create backend only (no Angular)
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

### Generate API Endpoint

```bash
goastra generate api product
```

Creates:
- `app/internal/handlers/product_handler.go`
- `app/internal/services/product_service.go`
- `app/internal/repository/product_repository.go`

### Generate Angular Module

```bash
goastra generate module product
```

Creates:
- `web/src/app/features/product/product.component.ts`
- `web/src/app/features/product/product.service.ts`

### Generate Full CRUD

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

```
my-app/
├── app/                      # Go backend
│   ├── cmd/server/           # Entry point
│   ├── internal/             # Internal packages
│   │   ├── auth/             # JWT authentication
│   │   ├── config/           # Configuration
│   │   ├── database/         # Database connection
│   │   ├── handlers/         # HTTP handlers
│   │   ├── logger/           # Structured logging (Zap)
│   │   ├── middleware/       # CORS, Auth, Logger
│   │   ├── models/           # Data models
│   │   ├── repository/       # Data access layer
│   │   ├── router/           # Route registration
│   │   ├── services/         # Business logic
│   │   └── validator/        # Request validation
│   └── database/
│       └── migrations/       # SQL migration files
├── web/                      # Angular frontend
│   └── src/
│       ├── app/
│       │   ├── core/         # Core services
│       │   ├── features/     # Feature modules
│       │   └── shared/       # Shared components
│       └── environments/     # Environment configs
├── schema/                   # Shared type definitions
├── .env.development          # Dev environment
├── .env.production           # Prod environment
├── .env.test                 # Test environment
└── goastra.json              # Framework config
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
