# LiteKPI

A metrics tracking and analytics platform for your products.

## Quick Start

### Prerequisites

- [Go 1.21+](https://golang.org/dl/)
- [Node.js 20+](https://nodejs.org/)
- [Docker](https://www.docker.com/) & Docker Compose
- [Air](https://github.com/cosmtrek/air) (for Go hot-reload)

### Development Setup

1. **Clone and setup environment**

   ```bash
   cp .env.example .env
   # Edit .env with your configuration
   ```

2. **Start the database**

   ```bash
   make db
   ```

3. **Run database migrations**

   ```bash
   make migrate
   ```

4. **Install dependencies**

   ```bash
   make install
   ```

5. **Start development servers**

   ```bash
   make dev
   ```

   Or start services individually:

   ```bash
   make dev-backend   # Go API with hot-reload (port 8080)
   make dev-frontend  # Vite dev server with HMR (port 5173)
   ```

6. **Access the application**

   - Frontend: http://localhost:5173
   - Backend API: http://localhost:8080
   - Health check: http://localhost:8080/health

## Project Structure

```
litekpi/
├── backend/                 # Go API server
│   ├── cmd/server/         # Application entry point
│   ├── internal/           # Private application code
│   │   ├── auth/           # Authentication module
│   │   ├── product/        # Product management
│   │   ├── ingest/         # Data ingestion
│   │   ├── dashboard/      # Dashboard management
│   │   ├── report/         # Report generation
│   │   ├── team/           # Team management
│   │   └── platform/       # Shared infrastructure
│   │       ├── config/     # Configuration
│   │       ├── database/   # Database connection
│   │       ├── middleware/ # HTTP middleware
│   │       ├── router/     # HTTP routing
│   │       └── email/      # Email service
│   └── migrations/         # Database migrations
│
├── frontend/               # React TypeScript app
│   ├── src/
│   │   ├── routes/         # TanStack Router pages
│   │   └── shared/         # Shared code
│   │       ├── components/ # UI components
│   │       ├── hooks/      # Custom hooks
│   │       ├── api/        # API client
│   │       ├── lib/        # Utilities
│   │       ├── types/      # TypeScript types
│   │       └── stores/     # Zustand stores
│   └── public/             # Static assets
│
├── _docs/                  # Documentation
│   └── TR-1--initial-concept/
│       ├── prd.md          # Product Requirements
│       └── tech-stack.md   # Technology Stack
│
├── docker-compose.yml      # Docker orchestration
├── Makefile               # Development commands
└── README.md              # This file
```

## Available Commands

| Command | Description |
|---------|-------------|
| `make dev` | Start all services with hot-reload |
| `make dev-backend` | Start backend only |
| `make dev-frontend` | Start frontend only |
| `make db` | Start PostgreSQL database |
| `make migrate` | Run database migrations |
| `make migrate-new name=X` | Create new migration |
| `make test` | Run all tests |
| `make build` | Build production images |
| `make clean` | Clean build artifacts |
| `make install` | Install all dependencies |
| `make fmt` | Format code |
| `make lint` | Lint code |

## Tech Stack

### Backend
- **Go 1.21** - Fast, compiled language
- **Chi** - Lightweight HTTP router
- **pgx** - PostgreSQL driver
- **golang-migrate** - Database migrations
- **JWT** - Authentication tokens

### Frontend
- **React 18** - UI library
- **TypeScript** - Type safety
- **Vite** - Build tool
- **TanStack Router** - Type-safe routing
- **TanStack Query** - Data fetching
- **Tailwind CSS** - Styling
- **shadcn/ui** - UI components
- **Zustand** - State management
- **Recharts** - Charts

### Infrastructure
- **PostgreSQL 15** - Database
- **Docker** - Containerization
- **Air** - Go hot-reload

## API Endpoints

| Method | Endpoint | Description |
|--------|----------|-------------|
| GET | `/health` | Health check |
| GET | `/api/v1/` | API status |

*More endpoints will be added as features are implemented.*

## Documentation

- [Product Requirements Document](./_docs/TR-1--initial-concept/prd.md)
- [Technology Stack](./_docs/TR-1--initial-concept/tech-stack.md)

## License

Private - All rights reserved.
