# LK-1: Metrics Ingestion

## Problem Statement

LiteKPI needs the ability to ingest metrics from client applications. Products already exist with API keys, but there's no mechanism to send and store metric data points. Without this foundational capability, the platform cannot provide any KPI tracking or analysis functionality.

This feature establishes the core data ingestion pipeline that all future analysis and visualization features will depend on.

---

## User Scenarios

### Persona: Developer (API Consumer)

A developer integrating their application with LiteKPI to track business or technical metrics.

---

### Scenario 1: Send a Single Metric

**Who:** Developer  
**What:** Send one metric data point to LiteKPI  
**Why:** Track an event or measurement as it occurs  

**User Flow:**
1. Developer has a product API key from LiteKPI
2. Developer makes POST request to single metric endpoint with API key in header
3. Request body contains metric name, value, and optional timestamp/metadata
4. System validates the request
5. System stores the metric
6. System returns success response

#### Functional Requirements

| ID | Requirement |
|----|-------------|
| F1.1 | System SHALL accept POST requests to `/api/v1/ingest` |
| F1.2 | System SHALL authenticate requests via `X-API-Key` header |
| F1.3 | System SHALL require `name` field (string) |
| F1.4 | System SHALL require `value` field (number - integer or float) |
| F1.5 | System SHALL accept optional `timestamp` field (ISO 8601 format) |
| F1.6 | System SHALL default timestamp to server time if not provided |
| F1.7 | System SHALL accept optional `metadata` field (object with string values) |
| F1.8 | System SHALL associate the metric with the product identified by the API key |

#### Non-Functional Requirements

| ID | Requirement |
|----|-------------|
| N1.1 | Endpoint SHALL respond within 200ms under normal load |
| N1.2 | System SHALL use constant-time comparison for API key validation |

#### Acceptance Criteria

- [ ] Valid request with name + value returns 201 Created
- [ ] Valid request with name + value + timestamp + metadata returns 201 Created
- [ ] Missing API key returns 401 Unauthorized
- [ ] Invalid API key returns 401 Unauthorized
- [ ] Missing name returns 400 Bad Request with descriptive error
- [ ] Missing value returns 400 Bad Request with descriptive error
- [ ] Invalid timestamp format returns 400 Bad Request with descriptive error
- [ ] Non-string metadata values return 400 Bad Request with descriptive error

---

### Scenario 2: Send Multiple Metrics (Batch)

**Who:** Developer  
**What:** Send multiple metrics in a single request  
**Why:** Reduce HTTP overhead when tracking multiple metrics at once  

**User Flow:**
1. Developer has a product API key from LiteKPI
2. Developer makes POST request to batch endpoint with API key in header
3. Request body contains array of metrics (each with name, value, optional timestamp/metadata)
4. System validates ALL metrics in the batch
5. If any metric is invalid, entire batch is rejected
6. If all valid, system stores all metrics
7. System returns success response with count

#### Functional Requirements

| ID | Requirement |
|----|-------------|
| F2.1 | System SHALL accept POST requests to `/api/v1/ingest/batch` |
| F2.2 | System SHALL authenticate requests via `X-API-Key` header |
| F2.3 | System SHALL accept array of metric objects in request body |
| F2.4 | System SHALL enforce maximum batch size of 100 metrics |
| F2.5 | System SHALL validate ALL metrics before storing any |
| F2.6 | System SHALL reject entire batch if ANY metric is invalid |
| F2.7 | System SHALL store all metrics atomically (all or nothing) |
| F2.8 | System SHALL return count of successfully stored metrics |

#### Non-Functional Requirements

| ID | Requirement |
|----|-------------|
| N2.1 | Batch endpoint SHALL respond within 1000ms for 100 metrics |

#### Acceptance Criteria

- [ ] Valid batch of 5 metrics returns 201 Created with `{"count": 5}`
- [ ] Batch with 1 invalid metric out of 10 returns 400 and stores nothing
- [ ] Batch of 101 metrics returns 400 Bad Request (exceeds limit)
- [ ] Empty batch returns 400 Bad Request
- [ ] Batch with duplicate metrics (same name + timestamp) across items returns 400

---

### Scenario 3: Idempotent Ingestion

**Who:** Developer  
**What:** Re-send a metric without creating duplicates  
**Why:** Handle retries and network failures gracefully  

**User Flow:**
1. Developer sends metric (name="orders", timestamp="2024-01-15T10:00:00Z", value=42)
2. Network timeout occurs, developer doesn't receive response
3. Developer retries with identical request
4. System detects duplicate (same product + name + timestamp)
5. System rejects duplicate with appropriate error

#### Functional Requirements

| ID | Requirement |
|----|-------------|
| F3.1 | System SHALL define uniqueness as: product_id + name + timestamp |
| F3.2 | System SHALL reject duplicate metrics with 409 Conflict |
| F3.3 | System SHALL check for duplicates within batch submissions |

#### Acceptance Criteria

- [ ] Sending same metric twice returns 409 Conflict on second attempt
- [ ] Metrics with same name but different timestamps are both stored
- [ ] Metrics with same name and timestamp but different products are both stored
- [ ] Batch containing internal duplicates returns 400 Bad Request

---

## Data Validation Rules

### Metric Name

| Rule | Specification |
|------|---------------|
| Format | `snake_case` - lowercase alphanumeric and underscores only |
| Pattern | `^[a-z][a-z0-9_]*$` |
| Min length | 1 character |
| Max length | 128 characters |
| Examples (valid) | `page_views`, `api_latency_ms`, `order_count` |
| Examples (invalid) | `PageViews`, `api-latency`, `123_count`, `_private` |

### Metric Value

| Rule | Specification |
|------|---------------|
| Type | Number (integer or floating point) |
| Range | Standard IEEE 754 double precision |
| Examples (valid) | `42`, `3.14159`, `-10`, `0`, `1000000` |
| Examples (invalid) | `"42"`, `null`, `NaN`, `Infinity` |

### Timestamp

| Rule | Specification |
|------|---------------|
| Format | ISO 8601 with timezone |
| Examples (valid) | `2024-01-15T10:30:00Z`, `2024-01-15T10:30:00+01:00` |
| Bounds | No restrictions (any past or future timestamp accepted) |
| Default | Server time (UTC) when not provided |

### Metadata

| Rule | Specification |
|------|---------------|
| Type | Object (key-value pairs) |
| Key format | String, non-empty |
| Value format | String only |
| Max keys | 20 |
| Max key length | 64 characters |
| Max value length | 256 characters |
| Examples (valid) | `{"env": "prod", "region": "us-east"}` |
| Examples (invalid) | `{"count": 5}`, `{"nested": {"key": "val"}}` |

---

## API Request/Response Examples

### Single Metric - Minimal

**Request:**
```http
POST /api/v1/ingest
X-API-Key: pk_live_abc123...
Content-Type: application/json

{
  "name": "page_views",
  "value": 1
}
```

**Response (201 Created):**
```json
{
  "id": "550e8400-e29b-41d4-a716-446655440000",
  "name": "page_views",
  "value": 1,
  "timestamp": "2024-01-15T10:30:00Z",
  "metadata": null
}
```

### Single Metric - Full

**Request:**
```http
POST /api/v1/ingest
X-API-Key: pk_live_abc123...
Content-Type: application/json

{
  "name": "api_latency_ms",
  "value": 234.5,
  "timestamp": "2024-01-15T10:30:00Z",
  "metadata": {
    "endpoint": "/api/users",
    "method": "GET",
    "status": "200"
  }
}
```

**Response (201 Created):**
```json
{
  "id": "550e8400-e29b-41d4-a716-446655440001",
  "name": "api_latency_ms",
  "value": 234.5,
  "timestamp": "2024-01-15T10:30:00Z",
  "metadata": {
    "endpoint": "/api/users",
    "method": "GET",
    "status": "200"
  }
}
```

### Batch Metrics

**Request:**
```http
POST /api/v1/ingest/batch
X-API-Key: pk_live_abc123...
Content-Type: application/json

{
  "metrics": [
    {"name": "orders", "value": 5},
    {"name": "revenue", "value": 499.95, "metadata": {"currency": "USD"}},
    {"name": "api_calls", "value": 1000, "timestamp": "2024-01-15T10:00:00Z"}
  ]
}
```

**Response (201 Created):**
```json
{
  "count": 3
}
```

### Error Response Examples

**Invalid name (400):**
```json
{
  "error": "validation_failed",
  "message": "Invalid metric name 'Page-Views': must be snake_case (lowercase alphanumeric and underscores, starting with letter)"
}
```

**Duplicate (409):**
```json
{
  "error": "duplicate_metric",
  "message": "Metric with name 'orders' and timestamp '2024-01-15T10:00:00Z' already exists"
}
```

**Batch validation failure (400):**
```json
{
  "error": "validation_failed",
  "message": "Invalid metric at index 2: value is required"
}
```

---

## Out of Scope

The following are explicitly NOT part of this feature:

- **Querying/reading metrics** - Future feature for analysis
- **Aggregations** (sum, avg, count) - Future feature
- **Metric deletion** - Not planned
- **Metric updates** - Metrics are immutable once stored
- **Rate limiting** - May be added later if needed
- **Compression** - Not needed at current scale
- **Retention policies** - Metrics stored forever
- **Metric type definitions** - No schema registry, all metrics are ad-hoc

---

## Dependencies

| Dependency | Description |
|------------|-------------|
| Products | Must have existing product with API key |
| API Key Auth | Existing middleware to validate product API keys |

---

## Assumptions

1. Product API key authentication middleware already exists or will be created
2. Low volume (<1,000 metrics/day/product) - no need for specialized timeseries database
3. PostgreSQL JSONB is sufficient for metadata storage and future querying
4. Clients can handle 409 Conflict responses for duplicate detection
