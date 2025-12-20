# LiteKPI

[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)

A lightweight, self-hostable KPI tracking platform. Track custom metrics from your applications and visualize them in beautiful dashboards.

## Features

- **Simple Metrics API** - Send any numeric metric with timestamps and metadata tags
- **Flexible Dashboards** - Visualize your KPIs with line charts, bar charts, and more
- **Metadata Filtering** - Slice and dice your data by custom tags
- **Multi-Product Support** - Track metrics for multiple applications
- **Team Access** - Invite team members to view and manage dashboards
- **OAuth Support** - Optional Google and GitHub authentication
- **Self-Hosted** - Run on your own infrastructure with Docker Compose

## Quick Start (Self-Hosting)

### Prerequisites

- [Docker](https://www.docker.com/) & Docker Compose v2+
- A domain name (optional but recommended for production)

### 1. Clone the Repository

```bash
git clone https://github.com/devbydaniel/litekpi.git
cd litekpi
```

### 2. Configure Environment

```bash
cp .env.example .env
```

Edit `.env` with your settings:

```bash
# Required - change these!
POSTGRES_PASSWORD=your-secure-database-password
JWT_SECRET=your-secure-jwt-secret-min-32-chars

# URLs where your app will be accessible
APP_URL=https://kpi.example.com       # Frontend URL
API_URL=https://api.kpi.example.com   # Backend API URL

# Optional - Email (for verification & password reset)
SMTP_HOST=smtp.example.com
SMTP_PORT=587
SMTP_USER=your-email@example.com
SMTP_PASSWORD=your-email-password
SMTP_FROM=noreply@example.com

# Optional - OAuth Providers
OAUTH_GOOGLE_CLIENT_ID=your-google-client-id
OAUTH_GOOGLE_CLIENT_SECRET=your-google-client-secret
OAUTH_GITHUB_CLIENT_ID=your-github-client-id
OAUTH_GITHUB_CLIENT_SECRET=your-github-client-secret
```

### 3. Start the Application

```bash
docker compose up -d
```

This starts:

- **PostgreSQL** database on port 5432 (internal)
- **Backend API** on port 8080
- **Frontend** on port 80

### 4. Verify Installation

```bash
# Check all services are running
docker compose ps

# Check API health
curl http://localhost:8080/health
```

### 5. Access the Application

Open `http://localhost` (or your configured `APP_URL`) and create your account.

## Configuration Reference

### Required Environment Variables

| Variable            | Description                                           |
| ------------------- | ----------------------------------------------------- |
| `POSTGRES_PASSWORD` | PostgreSQL database password                          |
| `JWT_SECRET`        | Secret key for JWT tokens (min 32 characters)         |
| `APP_URL`           | Frontend URL (e.g., `https://kpi.example.com`)        |
| `API_URL`           | Backend API URL (e.g., `https://api.kpi.example.com`) |

### Optional Environment Variables

| Variable                     | Default   | Description                |
| ---------------------------- | --------- | -------------------------- |
| `POSTGRES_USER`              | `litekpi` | PostgreSQL username        |
| `POSTGRES_DB`                | `litekpi` | PostgreSQL database name   |
| `SERVER_PORT`                | `8080`    | Backend server port        |
| `SMTP_HOST`                  | -         | SMTP server hostname       |
| `SMTP_PORT`                  | `587`     | SMTP server port           |
| `SMTP_USER`                  | -         | SMTP username              |
| `SMTP_PASSWORD`              | -         | SMTP password              |
| `SMTP_FROM`                  | -         | From address for emails    |
| `OAUTH_GOOGLE_CLIENT_ID`     | -         | Google OAuth client ID     |
| `OAUTH_GOOGLE_CLIENT_SECRET` | -         | Google OAuth client secret |
| `OAUTH_GITHUB_CLIENT_ID`     | -         | GitHub OAuth client ID     |
| `OAUTH_GITHUB_CLIENT_SECRET` | -         | GitHub OAuth client secret |

## Usage Guide

### Creating a Product

1. Log in to LiteKPI
2. Click **"New Product"** from the Products page
3. Give your product a name (e.g., "My SaaS App")
4. Copy the generated **API Key** - you'll need this to send metrics

### Sending Metrics

Use the HTTP API to send metrics from your application. Authenticate with the `X-API-Key` header.

#### Single Metric

```bash
curl -X POST https://api.kpi.example.com/api/v1/ingest \
  -H "Content-Type: application/json" \
  -H "X-API-Key: your-api-key" \
  -d '{
    "name": "daily_active_users",
    "value": 1250,
    "timestamp": "2024-01-15T00:00:00Z",
    "metadata": {
      "plan": "pro",
      "region": "us-east"
    }
  }'
```

#### Batch Metrics (up to 100 per request)

```bash
curl -X POST https://api.kpi.example.com/api/v1/ingest/batch \
  -H "Content-Type: application/json" \
  -H "X-API-Key: your-api-key" \
  -d '{
    "measurements": [
      {
        "name": "page_views",
        "value": 5420,
        "timestamp": "2024-01-15T00:00:00Z"
      },
      {
        "name": "signups",
        "value": 23,
        "timestamp": "2024-01-15T00:00:00Z",
        "metadata": {"source": "google"}
      }
    ]
  }'
```

### Metric Schema

| Field       | Type   | Required | Description                             |
| ----------- | ------ | -------- | --------------------------------------- |
| `name`      | string | Yes      | Metric name (snake_case, max 128 chars) |
| `value`     | number | Yes      | Numeric value                           |
| `timestamp` | string | No       | ISO 8601 timestamp (defaults to now)    |
| `metadata`  | object | No       | Key-value tags for filtering            |

### Metadata Constraints

- Maximum 20 keys per measurement
- Key names: max 64 characters
- Values: max 256 characters

### Example: Tracking from Different Languages

<details>
<summary><strong>JavaScript / Node.js</strong></summary>

```javascript
async function trackMetric(name, value, metadata = {}) {
  await fetch("https://api.kpi.example.com/api/v1/ingest", {
    method: "POST",
    headers: {
      "Content-Type": "application/json",
      "X-API-Key": process.env.LITEKPI_API_KEY,
    },
    body: JSON.stringify({
      name,
      value,
      timestamp: new Date().toISOString(),
      metadata,
    }),
  });
}

// Usage
await trackMetric("orders_total", 15000, { currency: "usd" });
```

</details>

<details>
<summary><strong>Python</strong></summary>

```python
import requests
from datetime import datetime
import os

def track_metric(name: str, value: float, metadata: dict = None):
    requests.post(
        'https://api.kpi.example.com/api/v1/ingest',
        headers={
            'Content-Type': 'application/json',
            'X-API-Key': os.environ['LITEKPI_API_KEY'],
        },
        json={
            'name': name,
            'value': value,
            'timestamp': datetime.utcnow().isoformat() + 'Z',
            'metadata': metadata or {},
        },
    )

# Usage
track_metric('monthly_revenue', 45000, {'plan': 'enterprise'})
```

</details>

<details>
<summary><strong>Go</strong></summary>

```go
package main

import (
    "bytes"
    "encoding/json"
    "net/http"
    "os"
    "time"
)

type Metric struct {
    Name      string            `json:"name"`
    Value     float64           `json:"value"`
    Timestamp string            `json:"timestamp"`
    Metadata  map[string]string `json:"metadata,omitempty"`
}

func TrackMetric(name string, value float64, metadata map[string]string) error {
    metric := Metric{
        Name:      name,
        Value:     value,
        Timestamp: time.Now().UTC().Format(time.RFC3339),
        Metadata:  metadata,
    }

    body, _ := json.Marshal(metric)
    req, _ := http.NewRequest("POST", "https://api.kpi.example.com/api/v1/ingest", bytes.NewBuffer(body))
    req.Header.Set("Content-Type", "application/json")
    req.Header.Set("X-API-Key", os.Getenv("LITEKPI_API_KEY"))

    _, err := http.DefaultClient.Do(req)
    return err
}
```

</details>

<details>
<summary><strong>cURL (one-liner)</strong></summary>

```bash
curl -X POST https://api.kpi.example.com/api/v1/ingest \
  -H "Content-Type: application/json" \
  -H "X-API-Key: $LITEKPI_API_KEY" \
  -d '{"name":"active_users","value":42}'
```

</details>

### Viewing Metrics

1. Go to **Products** and click on your product
2. Select a metric from the list
3. Use the toolbar to:
   - Change chart type (line, bar)
   - Adjust date range
   - Split by metadata key (e.g., see users by plan)
   - Filter by metadata values

## Production Deployment

### Using a Reverse Proxy (Recommended)

For production, place a reverse proxy (Nginx, Caddy, Traefik) in front of LiteKPI to handle:

- SSL/TLS termination
- Domain routing
- Rate limiting

Example with **Caddy** (automatic HTTPS):

```Caddyfile
kpi.example.com {
    reverse_proxy frontend:80
}

api.kpi.example.com {
    reverse_proxy backend:8080
}
```

Example with **Nginx**:

```nginx
server {
    listen 443 ssl;
    server_name kpi.example.com;

    ssl_certificate /etc/ssl/certs/kpi.example.com.pem;
    ssl_certificate_key /etc/ssl/private/kpi.example.com.key;

    location / {
        proxy_pass http://localhost:80;
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
    }
}

server {
    listen 443 ssl;
    server_name api.kpi.example.com;

    ssl_certificate /etc/ssl/certs/kpi.example.com.pem;
    ssl_certificate_key /etc/ssl/private/kpi.example.com.key;

    location / {
        proxy_pass http://localhost:8080;
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
    }
}
```

### Database Backups

Back up your PostgreSQL data regularly:

```bash
# Create backup
docker compose exec db pg_dump -U litekpi litekpi > backup.sql

# Restore backup
docker compose exec -T db psql -U litekpi litekpi < backup.sql
```

### Updating LiteKPI

```bash
# Pull latest changes
git pull

# Rebuild and restart
docker compose down
docker compose build
docker compose up -d
```

## Troubleshooting

### Check Logs

```bash
# All services
docker compose logs

# Specific service
docker compose logs backend
docker compose logs frontend
docker compose logs db
```

### Database Connection Issues

```bash
# Check database is running
docker compose ps db

# Test connection
docker compose exec db psql -U litekpi -c "SELECT 1"
```

### Reset Everything

```bash
# Stop and remove containers, networks, and volumes
docker compose down -v

# Start fresh
docker compose up -d
```

## API Reference

| Method   | Endpoint                            | Description          |
| -------- | ----------------------------------- | -------------------- |
| `GET`    | `/health`                           | Health check         |
| `POST`   | `/api/v1/auth/register`             | Register new account |
| `POST`   | `/api/v1/auth/login`                | Login                |
| `POST`   | `/api/v1/auth/logout`               | Logout               |
| `GET`    | `/api/v1/products`                  | List products        |
| `POST`   | `/api/v1/products`                  | Create product       |
| `GET`    | `/api/v1/products/:id`              | Get product          |
| `DELETE` | `/api/v1/products/:id`              | Delete product       |
| `POST`   | `/api/v1/ingest`                    | Ingest single metric |
| `POST`   | `/api/v1/ingest/batch`              | Ingest batch metrics |
| `GET`    | `/api/v1/products/:id/measurements` | List measurements    |

Full API documentation available at `/swagger/` when running the backend.

---

## Development

See [AGENTS.md](./AGENTS.md) for development setup and architecture details.

### Quick Dev Setup

```bash
# Install dependencies
make install

# Start dev services (PostgreSQL + Mailcatcher)
make dev-services

# Run migrations
make migrate

# Start backend & frontend with hot-reload
make dev
```

Access:

- Frontend: http://localhost:5173
- Backend API: http://localhost:8080
- Swagger UI: http://localhost:8080/swagger/
- Mailcatcher: http://localhost:1080

## License

MIT License - see [LICENSE](LICENSE) for details.
