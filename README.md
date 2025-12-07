# GoAstra

A production-ready full-stack framework combining Go backend with Angular frontend.

**Author:** [channdev](https://github.com/channdev)

## Features

- **Go Backend** - High-performance REST API with Gin framework
- **Angular Frontend** - Modern TypeScript SPA with standalone components
- **TypeSync** - Auto-generate TypeScript interfaces from Go structs
- **Code Generation** - Scaffold CRUD operations with a single command
- **JWT Authentication** - Built-in auth with refresh tokens
- **Environment Management** - Development, production, and test configs
- **Database Support** - PostgreSQL and MySQL with migrations

## Installation

### Prerequisites

Make sure you have the following installed:

- **Go 1.21+** - [Download Go](https://go.dev/dl/)
- **Node.js 18+** - [Download Node.js](https://nodejs.org/)
- **PostgreSQL 14+** or **MySQL 8+** (optional)

### Install via Go

```bash
go install github.com/channdev/goastra/cli/goastra@latest
```

### Install from Source

```bash
git clone https://github.com/channdev/goastra.git
cd goastra/cli/goastra
go build -o goastra .

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
# Create a new project
goastra new my-app

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
| `goastra generate api <name>` | Generate REST API endpoint |
| `goastra generate module <name>` | Generate Angular module |
| `goastra generate crud <name>` | Generate full-stack CRUD |
| `goastra typesync` | Sync Go types to TypeScript |
| `goastra test` | Run test suites |

### New Project Options

```bash
goastra new <project-name> [flags]
```

| Flag | Default | Description |
|------|---------|-------------|
| `-t, --template` | `default` | Project template (`default`, `minimal`) |
| `--db` | `postgres` | Database driver (`postgres`, `mysql`) |
| `--skip-angular` | `false` | Skip Angular frontend generation |
| `--skip-backend` | `false` | Skip Go backend generation |
| `--graphql` | `false` | Use GraphQL instead of REST |

**Examples:**

```bash
# Create with default template and PostgreSQL
goastra new my-app

# Create with minimal template
goastra new my-app -t minimal

# Create with MySQL database
goastra new my-app --db mysql

# Create backend only (no Angular)
goastra new my-api --skip-angular
```

## Project Structure

```
my-app/
├── app/                    # Go backend
│   ├── cmd/server/         # Entry point
│   ├── internal/           # Internal packages
│   │   ├── auth/           # JWT authentication
│   │   ├── config/         # Configuration loading
│   │   ├── database/       # Database connection
│   │   ├── handlers/       # HTTP handlers
│   │   ├── logger/         # Structured logging
│   │   ├── middleware/     # CORS, Auth, RequestID
│   │   ├── models/         # Data models
│   │   ├── repository/     # Data access layer
│   │   ├── router/         # Route registration
│   │   ├── services/       # Business logic
│   │   └── validator/      # Request validation
│   └── migrations/         # Database migrations
├── web/                    # Angular frontend
│   └── src/
│       ├── app/
│       │   ├── core/       # Core services
│       │   ├── features/   # Feature modules
│       │   └── shared/     # Shared components
│       └── environments/   # Environment configs
├── schema/                 # Shared type definitions
├── .env.development        # Dev environment
├── .env.production         # Prod environment
├── .env.test               # Test environment
└── goastra.json            # Framework config
```

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
  }
}
```

### Environment Variables

```bash
APP_ENV=development
PORT=8080
DB_URL=postgres://user:pass@localhost:5432/mydb
JWT_SECRET=your-secret-key
```

## Development

### Backend Only

```bash
goastra dev --backend
```

### Frontend Only

```bash
goastra dev --frontend
```

### Custom Ports

```bash
goastra dev -p 3000 --frontend-port 4300
```

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
- `web/src/app/features/product/product.module.ts`
- `web/src/app/features/product/product.component.ts`
- `web/src/app/features/product/product.service.ts`

### Generate Full CRUD

```bash
goastra generate crud product
```

Creates both backend API and frontend module with list, detail, create, edit, and delete components.

## Type Synchronization

Define your types in `schema/types/`:

```go
type Product struct {
    ID    uint   `json:"id"`
    Name  string `json:"name"`
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
- [Documentation](https://github.com/channdev/goastra#readme)
