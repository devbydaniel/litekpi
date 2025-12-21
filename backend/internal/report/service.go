package report

import (
	"context"
	"fmt"
	"strings"

	"github.com/google/uuid"

	"github.com/devbydaniel/litekpi/internal/kpi"
)

// Service handles report business logic.
type Service struct {
	repo       *Repository
	kpiService *kpi.Service
}

// NewService creates a new report service.
func NewService(repo *Repository, kpiService *kpi.Service) *Service {
	return &Service{
		repo:       repo,
		kpiService: kpiService,
	}
}

// CreateReport creates a new report.
func (s *Service) CreateReport(ctx context.Context, orgID uuid.UUID, req CreateReportRequest) (*Report, error) {
	name := strings.TrimSpace(req.Name)
	if name == "" {
		return nil, ErrReportNameEmpty
	}
	if len(name) > 255 {
		return nil, ErrReportNameTooLong
	}

	report, err := s.repo.CreateReport(ctx, orgID, name)
	if err != nil {
		return nil, fmt.Errorf("failed to create report: %w", err)
	}

	return report, nil
}

// ListReports returns all reports for an organization.
func (s *Service) ListReports(ctx context.Context, orgID uuid.UUID) ([]Report, error) {
	reports, err := s.repo.GetReportsByOrganizationID(ctx, orgID)
	if err != nil {
		return nil, fmt.Errorf("failed to list reports: %w", err)
	}

	return reports, nil
}

// GetReport returns a report with its KPIs after verifying organization ownership.
func (s *Service) GetReport(ctx context.Context, orgID, reportID uuid.UUID) (*ReportWithKPIs, error) {
	report, err := s.repo.GetReportByID(ctx, reportID)
	if err != nil {
		return nil, fmt.Errorf("failed to get report: %w", err)
	}
	if report == nil {
		return nil, ErrReportNotFound
	}
	if report.OrganizationID != orgID {
		return nil, ErrUnauthorized
	}

	kpis, err := s.kpiService.GetKPIsByReportID(ctx, reportID)
	if err != nil {
		return nil, fmt.Errorf("failed to get KPIs: %w", err)
	}

	return &ReportWithKPIs{
		Report: *report,
		KPIs:   kpis,
	}, nil
}

// UpdateReport updates a report's name.
func (s *Service) UpdateReport(ctx context.Context, orgID, reportID uuid.UUID, req UpdateReportRequest) (*Report, error) {
	report, err := s.repo.GetReportByID(ctx, reportID)
	if err != nil {
		return nil, fmt.Errorf("failed to get report: %w", err)
	}
	if report == nil {
		return nil, ErrReportNotFound
	}
	if report.OrganizationID != orgID {
		return nil, ErrUnauthorized
	}

	name := strings.TrimSpace(req.Name)
	if name == "" {
		return nil, ErrReportNameEmpty
	}
	if len(name) > 255 {
		return nil, ErrReportNameTooLong
	}

	if err := s.repo.UpdateReport(ctx, reportID, name); err != nil {
		return nil, fmt.Errorf("failed to update report: %w", err)
	}

	report.Name = name
	return report, nil
}

// DeleteReport deletes a report.
func (s *Service) DeleteReport(ctx context.Context, orgID, reportID uuid.UUID) error {
	report, err := s.repo.GetReportByID(ctx, reportID)
	if err != nil {
		return fmt.Errorf("failed to get report: %w", err)
	}
	if report == nil {
		return ErrReportNotFound
	}
	if report.OrganizationID != orgID {
		return ErrUnauthorized
	}

	if err := s.repo.DeleteReport(ctx, reportID); err != nil {
		return fmt.Errorf("failed to delete report: %w", err)
	}

	return nil
}

// ComputeReport returns a report with computed KPI values.
func (s *Service) ComputeReport(ctx context.Context, orgID, reportID uuid.UUID) (*ComputedReport, error) {
	report, err := s.repo.GetReportByID(ctx, reportID)
	if err != nil {
		return nil, fmt.Errorf("failed to get report: %w", err)
	}
	if report == nil {
		return nil, ErrReportNotFound
	}
	if report.OrganizationID != orgID {
		return nil, ErrUnauthorized
	}

	kpis, err := s.kpiService.GetKPIsByReportID(ctx, reportID)
	if err != nil {
		return nil, fmt.Errorf("failed to get KPIs: %w", err)
	}

	computedKPIs, err := s.kpiService.ComputeKPIs(ctx, kpis)
	if err != nil {
		return nil, fmt.Errorf("failed to compute KPIs: %w", err)
	}

	return &ComputedReport{
		Report: *report,
		KPIs:   computedKPIs,
	}, nil
}

// KPI operations delegated to kpi service

// CreateKPI creates a new KPI for a report.
func (s *Service) CreateKPI(ctx context.Context, orgID, reportID uuid.UUID, req kpi.CreateKPIRequest) (*kpi.KPI, error) {
	// Verify report ownership
	report, err := s.repo.GetReportByID(ctx, reportID)
	if err != nil {
		return nil, fmt.Errorf("failed to get report: %w", err)
	}
	if report == nil {
		return nil, ErrReportNotFound
	}
	if report.OrganizationID != orgID {
		return nil, ErrUnauthorized
	}

	return s.kpiService.CreateKPIForReport(ctx, orgID, reportID, req)
}

// GetKPIs returns all KPIs for a report.
func (s *Service) GetKPIs(ctx context.Context, orgID, reportID uuid.UUID) ([]kpi.KPI, error) {
	// Verify report ownership
	report, err := s.repo.GetReportByID(ctx, reportID)
	if err != nil {
		return nil, fmt.Errorf("failed to get report: %w", err)
	}
	if report == nil {
		return nil, ErrReportNotFound
	}
	if report.OrganizationID != orgID {
		return nil, ErrUnauthorized
	}

	return s.kpiService.GetKPIsByReportID(ctx, reportID)
}

// UpdateKPI updates a KPI.
func (s *Service) UpdateKPI(ctx context.Context, orgID, reportID, kpiID uuid.UUID, req kpi.UpdateKPIRequest) (*kpi.KPI, error) {
	// Verify report ownership
	report, err := s.repo.GetReportByID(ctx, reportID)
	if err != nil {
		return nil, fmt.Errorf("failed to get report: %w", err)
	}
	if report == nil {
		return nil, ErrReportNotFound
	}
	if report.OrganizationID != orgID {
		return nil, ErrUnauthorized
	}

	// Verify KPI belongs to this report
	existingKPI, err := s.kpiService.GetKPIByID(ctx, kpiID)
	if err != nil {
		return nil, err
	}
	if existingKPI.ReportID == nil || *existingKPI.ReportID != reportID {
		return nil, kpi.ErrKPINotFound
	}

	return s.kpiService.UpdateKPI(ctx, kpiID, req)
}

// DeleteKPI deletes a KPI.
func (s *Service) DeleteKPI(ctx context.Context, orgID, reportID, kpiID uuid.UUID) error {
	// Verify report ownership
	report, err := s.repo.GetReportByID(ctx, reportID)
	if err != nil {
		return fmt.Errorf("failed to get report: %w", err)
	}
	if report == nil {
		return ErrReportNotFound
	}
	if report.OrganizationID != orgID {
		return ErrUnauthorized
	}

	// Verify KPI belongs to this report
	existingKPI, err := s.kpiService.GetKPIByID(ctx, kpiID)
	if err != nil {
		return err
	}
	if existingKPI.ReportID == nil || *existingKPI.ReportID != reportID {
		return kpi.ErrKPINotFound
	}

	return s.kpiService.DeleteKPI(ctx, kpiID)
}

// ReorderKPIs reorders KPIs on a report.
func (s *Service) ReorderKPIs(ctx context.Context, orgID, reportID uuid.UUID, kpiIDs []uuid.UUID) error {
	// Verify report ownership
	report, err := s.repo.GetReportByID(ctx, reportID)
	if err != nil {
		return fmt.Errorf("failed to get report: %w", err)
	}
	if report == nil {
		return ErrReportNotFound
	}
	if report.OrganizationID != orgID {
		return ErrUnauthorized
	}

	return s.kpiService.ReorderKPIsForReport(ctx, reportID, kpiIDs)
}
