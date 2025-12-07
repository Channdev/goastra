want you to act as a senior software engineer and design a full-stack framework named GoAstra. Produce a highly technical, internal-RFC level document. The framework uses Go for the backend, Angular (TypeScript) for the frontend, and includes a Go-based CLI tool named goastra.

Define everything with production-grade detail, including architecture, CLI behavior, file structure, and environment handling. No emojis and no casual tone.

Framework Requirements

Tech Stack

Backend: Go

Frontend: Angular (TypeScript)

Communication: REST or GraphQL (choose one and justify)

Shared Types: Go structs auto-generate TypeScript interfaces and Angular services

CLI tool written in Go using Cobra or an equivalent library

Environment Management

The framework must support environment selection via .env files.

Define .env.development, .env.production, and .env.test.

The CLI must load the correct env file using either flags or auto-detection:

goastra dev → loads .env.development

goastra build → loads .env.production

goastra test → loads .env.test

Provide a recommended environment variable structure, such as:

APP_ENV=development
API_URL=http://localhost:8080
DB_URL=postgres://user:pass@localhost:5432/db
JWT_SECRET=...


Describe how environment variables propagate into:

Go backend configuration loader

Angular environment files during build

CLI command execution and dev server startup

Project Structure
Design a full monorepo layout:

goastra-project/
 ├── app/                 # Go backend
 ├── web/                 # Angular frontend
 ├── schema/              # Shared Go type definitions for codegen
 ├── cli/                 # GoAstra CLI source
 ├── .env.development
 ├── .env.production
 ├── .env.test
 └── goastra.json         # Framework config file


Backend Architecture

Routing system (controllers, middleware, validation)

Database abstraction (ORM or query builder, migrations)

Authentication (JWT, refresh tokens, RBAC)

Error handling standards

Logging conventions and log levels

Environment-aware config loader reading from .env.* files

Optional: WebSocket or SSE realtime support

Frontend Architecture

Modular Angular structure

Routing design (lazy-loaded modules, guards, role-based protection)

Codegen structure for auto-generated models, services, and CRUD UIs

Environment injection from GoAstra build pipeline

State management strategy (RxJS services or NgRx)

Code Generation System

Go → TypeScript model generator

Go → Angular service generator

CRUD module generator (components, routes, service, forms)

Schema loader for Go structs

TypeScript writer module

CLI Requirements
Design the following commands in detail:

goastra new <project>

Creates a full monorepo

Sets up Go backend, Angular frontend

Generates environment files

goastra dev

Runs Go backend with live reload

Runs Angular dev server

Loads .env.development

goastra build

Builds backend binary

Builds Angular production app

Loads .env.production

goastra generate api <name>

goastra generate module <name>

goastra generate crud <name>

CRUD for backend and frontend, including routing

goastra typesync

Generates TypeScript from Go structs

goastra test

Runs backend + frontend tests

Loads .env.test

Explain all CLI flags, configuration resolution, and how embedded templates work using Go’s embed package.