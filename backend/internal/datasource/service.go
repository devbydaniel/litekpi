package datasource

import (
	"context"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"strings"

	"github.com/google/uuid"
)

const (
	apiKeyPrefix = "lk_"
	apiKeyBytes  = 32
)

// Service handles data source business logic.
type Service struct {
	repo *Repository
}

// NewService creates a new data source service.
func NewService(repo *Repository) *Service {
	return &Service{repo: repo}
}

// CreateDataSource creates a new data source and returns the plain API key.
func (s *Service) CreateDataSource(ctx context.Context, orgID uuid.UUID, req CreateDataSourceRequest) (*CreateDataSourceResponse, error) {
	name := strings.TrimSpace(req.Name)
	if name == "" {
		return nil, ErrDataSourceNameEmpty
	}

	plainKey, keyHash, err := generateAPIKey()
	if err != nil {
		return nil, fmt.Errorf("failed to generate API key: %w", err)
	}

	ds, err := s.repo.CreateDataSource(ctx, orgID, name, keyHash)
	if err != nil {
		return nil, fmt.Errorf("failed to create data source: %w", err)
	}

	return &CreateDataSourceResponse{
		DataSource: *ds,
		APIKey:     plainKey,
	}, nil
}

// ListDataSources returns all data sources for an organization.
func (s *Service) ListDataSources(ctx context.Context, orgID uuid.UUID) ([]DataSource, error) {
	dataSources, err := s.repo.GetDataSourcesByOrganizationID(ctx, orgID)
	if err != nil {
		return nil, fmt.Errorf("failed to list data sources: %w", err)
	}

	if dataSources == nil {
		dataSources = []DataSource{}
	}

	return dataSources, nil
}

// GetDataSource returns a single data source after verifying organization ownership.
func (s *Service) GetDataSource(ctx context.Context, orgID, dataSourceID uuid.UUID) (*DataSource, error) {
	ds, err := s.repo.GetDataSourceByID(ctx, dataSourceID)
	if err != nil {
		return nil, fmt.Errorf("failed to get data source: %w", err)
	}
	if ds == nil {
		return nil, ErrDataSourceNotFound
	}
	if ds.OrganizationID != orgID {
		return nil, ErrUnauthorized
	}

	return ds, nil
}

// DeleteDataSource deletes a data source after verifying organization ownership.
func (s *Service) DeleteDataSource(ctx context.Context, orgID, dataSourceID uuid.UUID) error {
	ds, err := s.repo.GetDataSourceByID(ctx, dataSourceID)
	if err != nil {
		return fmt.Errorf("failed to get data source: %w", err)
	}
	if ds == nil {
		return ErrDataSourceNotFound
	}
	if ds.OrganizationID != orgID {
		return ErrUnauthorized
	}

	if err := s.repo.DeleteDataSource(ctx, dataSourceID); err != nil {
		return fmt.Errorf("failed to delete data source: %w", err)
	}

	return nil
}

// RegenerateAPIKey regenerates the API key for a data source.
func (s *Service) RegenerateAPIKey(ctx context.Context, orgID, dataSourceID uuid.UUID) (*RegenerateKeyResponse, error) {
	ds, err := s.repo.GetDataSourceByID(ctx, dataSourceID)
	if err != nil {
		return nil, fmt.Errorf("failed to get data source: %w", err)
	}
	if ds == nil {
		return nil, ErrDataSourceNotFound
	}
	if ds.OrganizationID != orgID {
		return nil, ErrUnauthorized
	}

	plainKey, keyHash, err := generateAPIKey()
	if err != nil {
		return nil, fmt.Errorf("failed to generate API key: %w", err)
	}

	if err := s.repo.UpdateAPIKeyHash(ctx, dataSourceID, keyHash); err != nil {
		return nil, fmt.Errorf("failed to update API key: %w", err)
	}

	return &RegenerateKeyResponse{
		APIKey: plainKey,
	}, nil
}

// generateAPIKey generates a new API key and its hash.
func generateAPIKey() (plainKey, hash string, err error) {
	bytes := make([]byte, apiKeyBytes)
	if _, err := rand.Read(bytes); err != nil {
		return "", "", err
	}

	plainKey = apiKeyPrefix + base64.URLEncoding.EncodeToString(bytes)

	hashBytes := sha256.Sum256([]byte(plainKey))
	hash = hex.EncodeToString(hashBytes[:])

	return plainKey, hash, nil
}
