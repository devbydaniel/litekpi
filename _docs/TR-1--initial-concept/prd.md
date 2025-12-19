# PRD: LiteKPI - KPI Tracking Platform

**Ticket:** TR-1
**Status:** Draft
**Last Updated:** 2025-12-05

---

## 1. Problem Statement

Product builders (solo developers, founders, small teams) need a simple way to track custom KPIs from their applications. Existing solutions are either:

- Too complex and enterprise-focused (Mixpanel, Amplitude)
- Too infrastructure-focused (Prometheus, Grafana)
- Not self-hostable / privacy-respecting

**Goal:** Build a lightweight, self-hostable web application where users can send arbitrary time-series metrics via API and visualize them in customizable dashboards.

---

## 2. User Scenarios and User Flows

### 2.1 Personas

| Persona            | Description                                                                  |
| ------------------ | ---------------------------------------------------------------------------- |
| **Solo Developer** | Individual building side projects, wants quick visibility into product usage |
| **Team Member**    | Part of a small team, needs shared access to product metrics                 |

---

### 2.2 Scenario: Account Registration & Setup

**Persona:** Solo Developer / Team Member

**Flow:**

1. User navigates to the application
2. User registers via email/password OR OAuth (Google/GitHub)
3. User lands on an empty dashboard state prompting them to create their first product

#### Functional Requirements

| ID     | Requirement                                                             |
| ------ | ----------------------------------------------------------------------- |
| FR-1.1 | System shall allow user registration via email and password             |
| FR-1.2 | System shall allow user registration/login via Google OAuth             |
| FR-1.3 | System shall allow user registration/login via GitHub OAuth             |
| FR-1.4 | System shall require email verification for email/password registration |
| FR-1.5 | System shall support password reset via email                           |

#### Non-Functional Requirements

| ID      | Requirement                                                         |
| ------- | ------------------------------------------------------------------- |
| NFR-1.1 | OAuth login shall complete within 3 seconds under normal conditions |
| NFR-1.2 | Passwords shall be hashed using bcrypt or argon2                    |

#### Acceptance Criteria

- [ ] User can register with email/password and receives verification email
- [ ] User can register/login with Google OAuth
- [ ] User can register/login with GitHub OAuth
- [ ] User can reset forgotten password
- [ ] After registration, user sees empty state with CTA to create first product

#### Dependencies

- SMTP configuration for sending emails
- OAuth app credentials for Google and GitHub

#### Assumptions

- Users have valid email addresses
- Self-hosters will configure their own OAuth apps or disable OAuth

---

### 2.3 Scenario: Product Creation & API Key Management

**Persona:** Solo Developer / Team Member

**Flow:**

1. User clicks "Create Product"
2. User enters product name
3. System creates product and generates an API key
4. User copies the API key for use in their application
5. User can view, regenerate, or revoke API keys

#### Functional Requirements

| ID     | Requirement                                                                    |
| ------ | ------------------------------------------------------------------------------ |
| FR-2.1 | User shall be able to create multiple products                                 |
| FR-2.2 | Each product shall have a unique, auto-generated API key                       |
| FR-2.3 | User shall be able to view their API key (shown once at creation, then masked) |
| FR-2.4 | User shall be able to regenerate an API key (invalidates the old one)          |
| FR-2.5 | User shall be able to delete a product (deletes all associated data)           |
| FR-2.6 | Products shall be isolated—data from one product is not visible in another     |

#### Non-Functional Requirements

| ID      | Requirement                                                       |
| ------- | ----------------------------------------------------------------- |
| NFR-2.1 | API keys shall be cryptographically secure (min 32 bytes entropy) |
| NFR-2.2 | API key validation shall complete within 50ms                     |

#### Acceptance Criteria

- [ ] User can create a product with a name
- [ ] Upon creation, API key is displayed and user can copy it
- [ ] API key is masked after initial display (user can regenerate if lost)
- [ ] Regenerating API key invalidates the previous key immediately
- [ ] Deleting a product prompts for confirmation and removes all data

#### Dependencies

- None

#### Assumptions

- Users understand that regenerating an API key will break existing integrations

---

### 2.4 Scenario: Sending Data via API

**Persona:** Solo Developer (integrating from their application)

**Flow:**

1. Developer includes API key in their application
2. Application sends HTTP POST requests with metric data
3. System validates API key and ingests data
4. Data becomes available for visualization

#### Functional Requirements

| ID     | Requirement                                                                      |
| ------ | -------------------------------------------------------------------------------- |
| FR-3.1 | System shall accept data points via HTTP POST                                    |
| FR-3.2 | Each data point shall include: metric name, value (numeric), timestamp           |
| FR-3.3 | Each data point may include optional tags (key-value pairs) for filtering        |
| FR-3.4 | System shall validate API key and reject requests with invalid keys              |
| FR-3.5 | System shall validate data point schema and return clear errors for invalid data |
| FR-3.6 | System shall support batch ingestion (multiple data points per request)          |
| FR-3.7 | If timestamp is omitted, system shall use server receive time                    |

**Data Point Schema:**

```json
{
  "metric": "string (required)",
  "value": "number (required)",
  "timestamp": "ISO8601 string (optional)",
  "tags": {
    "key": "value (optional, all values are strings)"
  }
}
```

#### Non-Functional Requirements

| ID      | Requirement                                                           |
| ------- | --------------------------------------------------------------------- |
| NFR-3.1 | Ingestion endpoint shall handle up to 100 requests/second per product |
| NFR-3.2 | Ingestion endpoint shall respond within 200ms under normal load       |
| NFR-3.3 | Batch requests shall support up to 1000 data points per request       |

#### Acceptance Criteria

- [ ] Valid API key + valid data → 200 OK, data stored
- [ ] Invalid API key → 401 Unauthorized
- [ ] Valid API key + invalid data → 400 Bad Request with descriptive error
- [ ] Batch request with mixed valid/invalid points → clear per-point errors
- [ ] Data with tags can be filtered by those tags in dashboards

#### Dependencies

- None

#### Assumptions

- Clients will send data in JSON format
- Timestamps are in UTC

---

### 2.5 Scenario: Creating and Viewing Dashboards

**Persona:** Solo Developer / Team Member

**Flow:**

1. User navigates to a product
2. User creates a new dashboard (or views existing ones)
3. User adds widgets to the dashboard
4. User configures each widget: selects metric, time granularity, filters
5. Dashboard displays visualizations based on configuration

#### Functional Requirements

| ID     | Requirement                                                    |
| ------ | -------------------------------------------------------------- |
| FR-4.1 | User shall be able to create multiple dashboards per product   |
| FR-4.2 | User shall be able to name and rename dashboards               |
| FR-4.3 | User shall be able to delete dashboards                        |
| FR-4.4 | User shall be able to add widgets to a dashboard               |
| FR-4.5 | User shall be able to remove widgets from a dashboard          |
| FR-4.6 | User shall be able to reorder/rearrange widgets on a dashboard |

**Widget Configuration:**

| ID      | Requirement                                                                                         |
| ------- | --------------------------------------------------------------------------------------------------- |
| FR-4.7  | Each widget shall be configured with: metric name, visualization type, time range, time granularity |
| FR-4.8  | Widget shall support optional tag filters (e.g., show only data where `country=US`)                 |
| FR-4.9  | Supported visualization types: time series chart (line), aggregate value (single number), table     |
| FR-4.10 | Supported time granularities: minute, hour, day, week, month                                        |
| FR-4.11 | Supported aggregations for aggregate values: sum, count, average, min, max                          |
| FR-4.12 | Time range shall be configurable: last N hours/days/weeks/months, or custom date range              |

#### Non-Functional Requirements

| ID      | Requirement                                                    |
| ------- | -------------------------------------------------------------- |
| NFR-4.1 | Dashboard shall load within 2 seconds for typical data volumes |
| NFR-4.2 | Widgets shall update without full page reload                  |

#### Acceptance Criteria

- [ ] User can create a dashboard and it appears in the dashboard list
- [ ] User can add a time series widget showing a metric over time
- [ ] User can add an aggregate widget showing a single computed value
- [ ] User can add a table widget showing metric data in tabular form
- [ ] User can filter widget data by tags
- [ ] Changing time granularity updates the visualization
- [ ] Empty state is shown when no data matches the configuration

#### Dependencies

- Data ingestion must be working (FR-3.x)

#### Assumptions

- Users understand basic time-series concepts (granularity, aggregation)

---

### 2.6 Scenario: Team Access

**Persona:** Team Member

**Flow:**

1. Product owner invites team member via email
2. Team member receives invite and creates account (or logs in if existing)
3. Team member now has access to the product and all its dashboards

#### Functional Requirements

| ID     | Requirement                                                              |
| ------ | ------------------------------------------------------------------------ |
| FR-5.1 | Product owner shall be able to invite users by email                     |
| FR-5.2 | Invited user shall receive email with invite link                        |
| FR-5.3 | Upon accepting invite, user gains access to the product                  |
| FR-5.4 | All team members have equal access (view and edit dashboards, view data) |
| FR-5.5 | Product owner shall be able to remove team members                       |
| FR-5.6 | User shall see all products they have access to (owned + invited)        |

#### Non-Functional Requirements

| ID      | Requirement                                   |
| ------- | --------------------------------------------- |
| NFR-5.1 | Invite emails shall be sent within 30 seconds |

#### Acceptance Criteria

- [ ] Owner can invite a user by email
- [ ] Invited user receives email with working invite link
- [ ] After accepting, user sees the product in their product list
- [ ] Team member can view and edit dashboards
- [ ] Owner can remove a team member, who then loses access immediately

#### Dependencies

- SMTP configuration (FR-1.x)

#### Assumptions

- Simple access model (no roles/permissions beyond owner vs member)
- "Owner" is the user who created the product

---

### 2.7 Scenario: Scheduled Email Reports

**Persona:** Solo Developer / Team Member

**Flow:**

1. User navigates to a dashboard
2. User configures a scheduled report: selects frequency, recipients
3. System sends dashboard snapshot via email at scheduled times

#### Functional Requirements

| ID     | Requirement                                                   |
| ------ | ------------------------------------------------------------- |
| FR-6.1 | User shall be able to schedule email reports for a dashboard  |
| FR-6.2 | Supported frequencies: daily, weekly, monthly                 |
| FR-6.3 | User shall be able to specify recipient email addresses       |
| FR-6.4 | Email shall contain a rendered snapshot of the dashboard      |
| FR-6.5 | User shall be able to edit or delete scheduled reports        |
| FR-6.6 | User shall be able to set the time of day for report delivery |

#### Non-Functional Requirements

| ID      | Requirement                                               |
| ------- | --------------------------------------------------------- |
| NFR-6.1 | Reports shall be sent within 15 minutes of scheduled time |

#### Acceptance Criteria

- [ ] User can create a daily report for a dashboard
- [ ] Report email arrives at the scheduled time
- [ ] Email contains visual snapshot of all widgets on the dashboard
- [ ] User can add multiple recipients
- [ ] User can delete a scheduled report and emails stop

#### Dependencies

- SMTP configuration
- Dashboard rendering (FR-4.x)

#### Assumptions

- Dashboard snapshot is a static image or HTML rendering, not interactive

---

### 2.8 Scenario: Exporting Data

**Persona:** Solo Developer / Team Member

**Flow:**

1. User navigates to a dashboard or specific widget
2. User clicks "Export"
3. User downloads data as CSV

#### Functional Requirements

| ID     | Requirement                                                             |
| ------ | ----------------------------------------------------------------------- |
| FR-7.1 | User shall be able to export widget data as CSV                         |
| FR-7.2 | Export shall respect current widget configuration (time range, filters) |
| FR-7.3 | CSV shall include: timestamp, metric name, value, and all tags          |

#### Non-Functional Requirements

| ID      | Requirement                                                      |
| ------- | ---------------------------------------------------------------- |
| NFR-7.1 | Export shall complete within 30 seconds for up to 1 million rows |

#### Acceptance Criteria

- [ ] User can export a widget's data as CSV
- [ ] Downloaded CSV contains expected columns and data
- [ ] Filters applied to widget are reflected in exported data

#### Dependencies

- Dashboard and widget configuration (FR-4.x)

#### Assumptions

- Large exports may need to be handled asynchronously (future enhancement)

---

## 3. Out of Scope

The following are explicitly **not** included in this version:

| Item                                                            | Rationale                                    |
| --------------------------------------------------------------- | -------------------------------------------- |
| Alerts/notifications                                            | Deferred—user indicated not needed for v1    |
| Advanced visualizations (pie charts, gauges, heatmaps, funnels) | Keep it simple; time series focus            |
| Role-based permissions                                          | Simple team model is sufficient for now      |
| PDF export                                                      | CSV covers primary export needs              |
| Public/shared dashboard links                                   | Can be added later                           |
| Data retention policies                                         | Data kept forever for now                    |
| Kubernetes/Helm deployment                                      | Docker Compose is priority                   |
| Mobile app                                                      | Web-only for v1                              |
| Real-time streaming updates                                     | Near-real-time (page refresh) is acceptable  |
| Pre-built integrations (Stripe, Segment, etc.)                  | API-first approach; users send data directly |

---

## 4. Technical Constraints

| Constraint         | Detail                                                                  |
| ------------------ | ----------------------------------------------------------------------- |
| Self-hostable      | Must run via single `docker-compose up` command                         |
| Low complexity     | Target < 1,000 data points/day per product; no need for complex scaling |
| Configurable email | SMTP settings must be configurable via environment variables            |
| Configurable OAuth | OAuth can be disabled if self-hoster doesn't configure it               |

---

## 5. Open Questions

| #   | Question                                                                               | Impact                                          |
| --- | -------------------------------------------------------------------------------------- | ----------------------------------------------- |
| 1   | Should there be a hosted/SaaS version, or purely self-hosted?                          | Affects pricing model, infrastructure decisions |
| 2   | Should dashboard time zone be user-configurable or product-configurable?               | Affects data display                            |
| 3   | What happens to data when a team member is removed—do their created dashboards remain? | Affects data ownership model                    |
| 4   | Should there be a limit on number of products/dashboards per user for self-hosted?     | Affects resource planning                       |

---

## 6. Success Metrics

Once deployed, success can be measured by:

- User can go from registration to seeing first data point in < 10 minutes
- Self-hosting setup completes in < 5 minutes with Docker Compose
- Dashboard loads in < 2 seconds with typical data volumes

---

## 7. Glossary

| Term           | Definition                                                                                            |
| -------------- | ----------------------------------------------------------------------------------------------------- |
| **Product**    | An organizational container representing a user's application. Has its own API key and isolated data. |
| **Metric**     | A named measurement (e.g., "signups", "page_views", "revenue")                                        |
| **Data point** | A single measurement: metric name + numeric value + timestamp + optional tags                         |
| **Tag**        | A key-value pair attached to a data point for filtering/grouping (e.g., `country: "US"`)              |
| **Widget**     | A visualization component on a dashboard (chart, number, table)                                       |
| **Dashboard**  | A collection of widgets displaying metrics for a product                                              |
