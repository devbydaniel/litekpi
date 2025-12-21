package dashboard

import (
	"errors"
	"time"

	"github.com/google/uuid"
)

// Error definitions
var (
	ErrDashboardNotFound   = errors.New("dashboard not found")
	ErrWidgetNotFound      = errors.New("widget not found")
	ErrUnauthorized        = errors.New("unauthorized")
	ErrDashboardNameEmpty  = errors.New("dashboard name is required")
	ErrCannotDeleteDefault = errors.New("cannot delete default dashboard")
	ErrInvalidChartType    = errors.New("invalid chart type")
	ErrInvalidDateRange    = errors.New("invalid date range")
	ErrTitleTooLong        = errors.New("title exceeds maximum length of 128 characters")
)

// Dashboard represents a dashboard in the system.
type Dashboard struct {
	ID             uuid.UUID `json:"id"`
	Name           string    `json:"name"`
	OrganizationID uuid.UUID `json:"organizationId"`
	IsDefault      bool      `json:"isDefault"`
	CreatedAt      time.Time `json:"createdAt"`
	UpdatedAt      time.Time `json:"updatedAt"`
}

// Widget represents a chart widget on a dashboard.
type Widget struct {
	ID              uuid.UUID  `json:"id"`
	DashboardID     uuid.UUID  `json:"dashboardId"`
	DataSourceID    uuid.UUID  `json:"dataSourceId"`
	Title           *string    `json:"title,omitempty"`
	MeasurementName string     `json:"measurementName"`
	ChartType       string     `json:"chartType"` // area, bar, line
	DateRange       string     `json:"dateRange"` // last_7_days, last_30_days, custom
	DateFrom        *time.Time `json:"dateFrom,omitempty"`
	DateTo          *time.Time `json:"dateTo,omitempty"`
	SplitBy         *string    `json:"splitBy,omitempty"`
	Filters         []Filter   `json:"filters"`
	Position        int        `json:"position"`
	CreatedAt       time.Time  `json:"createdAt"`
	UpdatedAt       time.Time  `json:"updatedAt"`
}

// Filter represents a metadata filter for a widget.
type Filter struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

// Valid chart types
var validChartTypes = map[string]bool{
	"area": true,
	"bar":  true,
	"line": true,
}

// Valid date ranges
var validDateRanges = map[string]bool{
	"last_7_days":  true,
	"last_30_days": true,
	"custom":       true,
}

// Request/Response types

// CreateDashboardRequest is the request body for creating a dashboard.
type CreateDashboardRequest struct {
	Name string `json:"name"`
}

// UpdateDashboardRequest is the request body for updating a dashboard.
type UpdateDashboardRequest struct {
	Name string `json:"name"`
}

// CreateWidgetRequest is the request body for creating a widget.
type CreateWidgetRequest struct {
	DataSourceID    uuid.UUID  `json:"dataSourceId"`
	MeasurementName string     `json:"measurementName"`
	Title           *string    `json:"title,omitempty"`
	ChartType       string     `json:"chartType"`
	DateRange       string     `json:"dateRange"`
	DateFrom        *time.Time `json:"dateFrom,omitempty"`
	DateTo          *time.Time `json:"dateTo,omitempty"`
	SplitBy         *string    `json:"splitBy,omitempty"`
	Filters         []Filter   `json:"filters,omitempty"`
}

// UpdateWidgetRequest is the request body for updating a widget.
type UpdateWidgetRequest struct {
	Title     *string    `json:"title,omitempty"`
	ChartType string     `json:"chartType"`
	DateRange string     `json:"dateRange"`
	DateFrom  *time.Time `json:"dateFrom,omitempty"`
	DateTo    *time.Time `json:"dateTo,omitempty"`
	SplitBy   *string    `json:"splitBy,omitempty"`
	Filters   []Filter   `json:"filters,omitempty"`
}

// ReorderWidgetsRequest is the request body for reordering widgets.
type ReorderWidgetsRequest struct {
	WidgetIDs []uuid.UUID `json:"widgetIds"`
}

// DashboardWithWidgets is a dashboard with its widgets.
type DashboardWithWidgets struct {
	Dashboard Dashboard `json:"dashboard"`
	Widgets   []Widget  `json:"widgets"`
}

// ListDashboardsResponse is the response for listing dashboards.
type ListDashboardsResponse struct {
	Dashboards []Dashboard `json:"dashboards"`
}

// MessageResponse is a generic response with a message.
type MessageResponse struct {
	Message string `json:"message"`
}

// ErrorResponse represents an API error.
type ErrorResponse struct {
	Error string `json:"error"`
}
