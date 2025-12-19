package demo

import (
	"context"
	"fmt"
	"math"
	"math/rand"
	"time"

	"github.com/google/uuid"

	"github.com/devbydaniel/litekpi/internal/ingest"
	"github.com/devbydaniel/litekpi/internal/product"
)

// Service handles demo product creation business logic.
type Service struct {
	productService *product.Service
	ingestService  *ingest.Service
}

// NewService creates a new demo service.
func NewService(productService *product.Service, ingestService *ingest.Service) *Service {
	return &Service{
		productService: productService,
		ingestService:  ingestService,
	}
}

// CreateDemoProduct creates a demo product with sample measurements for the last 30 days.
func (s *Service) CreateDemoProduct(ctx context.Context, orgID uuid.UUID) (*product.CreateProductResponse, error) {
	// Create the demo product
	response, err := s.productService.CreateProduct(ctx, orgID, product.CreateProductRequest{
		Name: "Demo Product",
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create demo product: %w", err)
	}

	// Generate demo measurements for last 30 days
	if err := s.createDemoMeasurements(ctx, response.Product.ID); err != nil {
		// Rollback: delete the product if measurements fail
		s.productService.DeleteProduct(ctx, orgID, response.Product.ID)
		return nil, fmt.Errorf("failed to create demo measurements: %w", err)
	}

	return response, nil
}

// createDemoMeasurements generates realistic demo data for the last 30 days.
func (s *Service) createDemoMeasurements(ctx context.Context, productID uuid.UUID) error {
	now := time.Now().UTC()
	var metrics []ingest.IngestRequest

	for i := 30; i >= 0; i-- {
		date := now.AddDate(0, 0, -i)
		isWeekend := date.Weekday() == time.Saturday || date.Weekday() == time.Sunday
		timestamp := date.Format(time.RFC3339)

		// Daily Active Users - base 1200, 5% growth trend, weekend dip
		dauBase := 1200.0 * (1 + 0.05*float64(30-i)/30)
		if isWeekend {
			dauBase *= 0.8
		}
		dauValue := dauBase * (0.85 + rand.Float64()*0.30) // ±15% noise
		metrics = append(metrics, ingest.IngestRequest{
			Name:      "daily_active_users",
			Value:     math.Round(dauValue),
			Timestamp: timestamp,
		})

		// Revenue - base 2500, weekend dip, with currency metadata
		revenueBase := 2500.0
		if isWeekend {
			revenueBase *= 0.7
		}
		revenueValue := revenueBase * (0.75 + rand.Float64()*0.50) // ±25% noise
		metrics = append(metrics, ingest.IngestRequest{
			Name:      "revenue",
			Value:     math.Round(revenueValue*100) / 100,
			Timestamp: timestamp,
			Metadata:  map[string]string{"currency": "usd"},
		})

		// Page Views - web and mobile with source metadata
		webBase := 8000.0
		mobileBase := 4000.0
		if isWeekend {
			webBase *= 0.75
			mobileBase *= 0.75
		}

		metrics = append(metrics, ingest.IngestRequest{
			Name:      "page_views",
			Value:     math.Round(webBase * (0.80 + rand.Float64()*0.40)),
			Timestamp: timestamp,
			Metadata:  map[string]string{"source": "web"},
		})
		metrics = append(metrics, ingest.IngestRequest{
			Name:      "page_views",
			Value:     math.Round(mobileBase * (0.80 + rand.Float64()*0.40)),
			Timestamp: date.Add(time.Second).Format(time.RFC3339), // +1s to avoid duplicate
			Metadata:  map[string]string{"source": "mobile"},
		})
	}

	// Batch insert (max 100 per batch)
	for i := 0; i < len(metrics); i += 100 {
		end := i + 100
		if end > len(metrics) {
			end = len(metrics)
		}
		batch := ingest.BatchIngestRequest{Metrics: metrics[i:end]}
		if _, err := s.ingestService.IngestBatch(ctx, productID, batch); err != nil {
			return err
		}
	}

	return nil
}
