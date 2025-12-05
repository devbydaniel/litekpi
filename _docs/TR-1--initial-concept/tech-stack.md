# Tech Stack Decision

**Ticket:** TR-1
**Status:** Approved
**Last Updated:** 2025-12-05

---

## Summary

This document captures the technology choices for the Trackable KPI platform. Decisions are driven by the core constraints: **self-hostable via Docker Compose**, **low complexity**, and **low data volume**.

---

## Stack Overview

| Layer                | Technology              | Version Policy        |
| -------------------- | ----------------------- | --------------------- |
| **Backend**          | Go                      | Latest stable (1.21+) |
| **Frontend**         | React + TypeScript      | React 18+, Vite       |
| **Database**         | PostgreSQL              | 15+                   |
| **Containerization** | Docker + Docker Compose | Latest stable         |

---

## Backend: Go

### Choice Rationale

| Factor                   | Benefit                                                   |
| ------------------------ | --------------------------------------------------------- |
| **Single binary**        | No runtime dependencies, simple Docker images             |
| **Low memory footprint** | Ideal for self-hosted environments with limited resources |
| **Fast cold starts**     | Quick container restarts                                  |
| **Strong stdlib**        | HTTP server, JSON, crypto all built-in                    |
| **Concurrency model**    | Goroutines handle concurrent API requests efficiently     |

### Key Libraries

| Purpose     | Library                             | Notes                                                            |
| ----------- | ----------------------------------- | ---------------------------------------------------------------- |
| HTTP Router | `chi` or `echo`                     | Lightweight, idiomatic                                           |
| Database    | `pgx`                               | Native PostgreSQL driver, better performance than `database/sql` |
| Migrations  | `golang-migrate/migrate`            | SQL-based migrations                                             |
| JWT         | `golang-jwt/jwt`                    | Token generation/validation                                      |
| OAuth       | `golang.org/x/oauth2`               | Standard OAuth2 client                                           |
| Email       | `go-mail/mail` or stdlib `net/smtp` | SMTP sending                                                     |
| Scheduler   | `robfig/cron`                       | Cron-style job scheduling for email reports                      |
| Validation  | `go-playground/validator`           | Struct validation                                                |
| Config      | `env` + `envconfig` or `viper`      | Environment-based configuration                                  |

### Project Structure (proposed)

Organized by **bounded context** rather than by layer. Each context contains its own domain types, repository, and service.

```
backend/
├── cmd/
│   └── server/
│       └── main.go              # Entry point, wires everything together
├── internal/
│   ├── auth/                    # Authentication bounded context
│   │   ├── domain.go            # User, Session types
│   │   ├── repository.go        # User persistence
│   │   ├── service.go           # Login, register, OAuth, JWT logic
│   │   └── handler.go           # HTTP handlers for auth routes
│   │
│   ├── product/                 # Product bounded context
│   │   ├── domain.go            # Product, APIKey types
│   │   ├── repository.go        # Product persistence
│   │   ├── service.go           # CRUD, API key generation
│   │   └── handler.go           # HTTP handlers for product routes
│   │
│   ├── ingest/                  # Data ingestion bounded context
│   │   ├── domain.go            # DataPoint, Metric types
│   │   ├── repository.go        # Time-series data persistence
│   │   ├── service.go           # Validation, batch insert
│   │   └── handler.go           # HTTP handlers for ingestion API
│   │
│   ├── dashboard/               # Dashboard bounded context
│   │   ├── domain.go            # Dashboard, Widget types
│   │   ├── repository.go        # Dashboard persistence
│   │   ├── service.go           # CRUD, widget configuration
│   │   ├── query.go             # Time-series query logic
│   │   └── handler.go           # HTTP handlers for dashboard routes
│   │
│   ├── report/                  # Reporting bounded context
│   │   ├── domain.go            # ScheduledReport types
│   │   ├── repository.go        # Report schedule persistence
│   │   ├── service.go           # Scheduling, rendering, export
│   │   ├── scheduler.go         # Cron job definitions
│   │   └── handler.go           # HTTP handlers for report config
│   │
│   ├── team/                    # Team access bounded context
│   │   ├── domain.go            # Invite, Membership types
│   │   ├── repository.go        # Team membership persistence
│   │   ├── service.go           # Invite, accept, remove logic
│   │   └── handler.go           # HTTP handlers for team routes
│   │
│   └── platform/                # Shared infrastructure (not a bounded context)
│       ├── config/              # Configuration loading
│       ├── database/            # DB connection, migrations runner
│       ├── middleware/          # Auth middleware, logging, CORS
│       ├── router/              # Route registration, server setup
│       └── email/               # SMTP client
│
├── migrations/                  # SQL migration files
├── Dockerfile
└── go.mod
```

**Rationale:** This structure keeps related code together, making it easier to understand and modify a single feature without jumping across multiple directories. Cross-cutting infrastructure lives in `platform/`.

---

## Frontend: React + TypeScript

### Choice Rationale

| Factor             | Benefit                           |
| ------------------ | --------------------------------- |
| **Ecosystem**      | Best-in-class charting libraries  |
| **Developer pool** | Easy to find React developers     |
| **TypeScript**     | Type safety, better DX            |
| **Vite**           | Fast builds, great dev experience |

### Key Libraries

| Purpose          | Library                                           | Notes                                       |
| ---------------- | ------------------------------------------------- | ------------------------------------------- |
| Build tool       | Vite                                              | Fast HMR, optimized builds                  |
| Routing          | React Router                                      | Standard routing                            |
| State management | Zustand or React Context                          | Keep it simple, no Redux overhead           |
| HTTP client      | `fetch` + custom hooks or `@tanstack/react-query` | Caching, loading states                     |
| Charts           | Recharts or Apache ECharts                        | Recharts = simpler, ECharts = more powerful |
| UI components    | Tailwind CSS + Headless UI or shadcn/ui           | Utility-first, accessible                   |
| Forms            | React Hook Form                                   | Performant form handling                    |
| Date handling    | `date-fns`                                        | Lightweight date utilities                  |

### Project Structure (proposed)

Organized by **page**. Code is colocated with the page that uses it and only moved to `shared/` when used by multiple pages.

```
frontend/
├── src/
│   ├── pages/
│   │   ├── auth/
│   │   │   ├── LoginPage.tsx
│   │   │   ├── RegisterPage.tsx
│   │   │   ├── ResetPasswordPage.tsx
│   │   │   ├── components/          # Auth-specific components
│   │   │   ├── hooks/               # Auth-specific hooks
│   │   │   └── api.ts               # Auth API calls
│   │   │
│   │   ├── products/
│   │   │   ├── ProductListPage.tsx
│   │   │   ├── ProductSettingsPage.tsx
│   │   │   ├── components/
│   │   │   └── api.ts
│   │   │
│   │   ├── dashboard/
│   │   │   ├── DashboardPage.tsx
│   │   │   ├── DashboardEditPage.tsx
│   │   │   ├── components/          # Widgets, charts, etc.
│   │   │   ├── hooks/               # Query/aggregation hooks
│   │   │   └── api.ts
│   │   │
│   │   ├── reports/
│   │   │   ├── ReportsPage.tsx
│   │   │   ├── components/
│   │   │   └── api.ts
│   │   │
│   │   └── team/
│   │       ├── TeamPage.tsx
│   │       ├── components/
│   │       └── api.ts
│   │
│   ├── shared/                      # Only code used by multiple pages
│   │   ├── components/              # Shared UI (Button, Input, Modal, etc.)
│   │   ├── hooks/                   # Shared hooks (useAuth, useFetch, etc.)
│   │   ├── api/                     # Shared API client setup
│   │   ├── types/                   # Shared TypeScript types
│   │   └── utils/                   # Shared utilities
│   │
│   ├── App.tsx                      # Root component, routing
│   └── main.tsx                     # Entry point
│
├── public/
├── index.html
├── vite.config.ts
├── tailwind.config.js
└── package.json
```

**Rationale:** Colocation keeps related code together, making it easier to work on a single feature. Only extract to `shared/` when a component, hook, or utility is genuinely reused—this avoids premature abstraction.

---

## Database: PostgreSQL

### Choice Rationale

| Factor                   | Benefit                                          |
| ------------------------ | ------------------------------------------------ |
| **Battle-tested**        | Proven reliability                               |
| **Docker-friendly**      | Official image, easy setup                       |
| **JSON support**         | JSONB for flexible tag storage                   |
| **Aggregations**         | Built-in functions for sum, avg, count, min, max |
| **Sufficient for scale** | < 1,000 data points/day is trivial for Postgres  |

### Schema Approach

- **Time-series data**: Standard table with timestamp index; no need for TimescaleDB at this volume
- **Tags**: Stored as JSONB column for flexible querying
- **Indexes**: B-tree on timestamp, GIN on tags JSONB

### Why Not Specialized Time-Series DBs?

| Option      | Reason to Skip                                 |
| ----------- | ---------------------------------------------- |
| TimescaleDB | Adds complexity; not needed at < 1k points/day |
| ClickHouse  | Overkill; harder to self-host simply           |
| InfluxDB    | Different query paradigm; Postgres is simpler  |

---

## Authentication

### Approach

| Component          | Implementation                              |
| ------------------ | ------------------------------------------- |
| Session/Token      | JWT (stateless, stored in httpOnly cookie)  |
| Password hashing   | Argon2id (via `golang.org/x/crypto/argon2`) |
| OAuth providers    | Google, GitHub via `golang.org/x/oauth2`    |
| Email verification | Token-based, sent via SMTP                  |

### Configuration

OAuth providers are **optional**. If `OAUTH_GOOGLE_CLIENT_ID` / `OAUTH_GITHUB_CLIENT_ID` are not set, those login options are hidden in the UI.

---

## Email

### Approach

- SMTP-based sending via environment configuration
- No external email service dependencies (Sendgrid, Mailgun, etc.)
- Self-hosters configure their own SMTP (or use services like Mailgun SMTP)

### Environment Variables

```
SMTP_HOST=smtp.example.com
SMTP_PORT=587
SMTP_USER=user@example.com
SMTP_PASSWORD=secret
SMTP_FROM=noreply@example.com
```

---

## Scheduled Jobs

### Approach

- Built-in scheduler using `robfig/cron`
- Runs in the same process as the API server
- No external job queue (Redis, RabbitMQ) needed at this scale

### Jobs

| Job              | Schedule                               | Purpose                                  |
| ---------------- | -------------------------------------- | ---------------------------------------- |
| Email reports    | User-configured (daily/weekly/monthly) | Render dashboard snapshot, send via SMTP |
| Cleanup (future) | Optional                               | Data retention if added later            |

---

## Containerization

### Docker Compose Structure

```yaml
services:
  app:
    build: .
    ports:
      - "8080:8080"
    environment:
      - DATABASE_URL=postgres://...
      - SMTP_HOST=...
    depends_on:
      - db

  db:
    image: postgres:15-alpine
    volumes:
      - postgres_data:/var/lib/postgresql/data
    environment:
      - POSTGRES_DB=trackable
      - POSTGRES_USER=trackable
      - POSTGRES_PASSWORD=secret

volumes:
  postgres_data:
```

### Docker Image Strategy

| Stage         | Purpose                                                                    |
| ------------- | -------------------------------------------------------------------------- |
| Build stage   | Compile Go binary                                                          |
| Runtime stage | Minimal image (`alpine` or `scratch`) with binary + static frontend assets |

Target image size: **< 50MB** (Go binary + frontend build)

---

## Development Environment

### Prerequisites

- Go 1.21+
- Node.js 20+ (for frontend)
- Docker + Docker Compose
- PostgreSQL (or use Docker)

### Local Development

```bash
# Backend
cd backend && go run ./cmd/server

# Frontend
cd frontend && npm run dev

# Full stack via Docker
docker-compose up
```

---

## Future Considerations (Not in v1)

| Item               | When to Consider                                       |
| ------------------ | ------------------------------------------------------ |
| Redis              | If we add real-time features or caching layer          |
| TimescaleDB        | If data volume grows significantly (> 100k points/day) |
| Kubernetes/Helm    | If users request it; Docker Compose is priority        |
| CDN/Static hosting | If we offer a hosted SaaS version                      |

---

## Decision Log

| Date       | Decision                          | Rationale                                                  |
| ---------- | --------------------------------- | ---------------------------------------------------------- |
| 2025-12-05 | Go for backend                    | Self-hosting priority, low memory footprint, single binary |
| 2025-12-05 | React + TypeScript for frontend   | Best charting ecosystem, developer availability            |
| 2025-12-05 | PostgreSQL for database           | Simple, proven, sufficient for target volume               |
| 2025-12-05 | Built-in scheduler over job queue | Simplicity; no need for Redis at this scale                |
