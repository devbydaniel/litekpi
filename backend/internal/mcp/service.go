package mcp

import (
	"context"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"strings"

	"github.com/google/uuid"

	"github.com/devbydaniel/litekpi/internal/datasource"
)

const (
	apiKeyPrefix = "lkmcp_"
	apiKeyBytes  = 32
)

// Service handles MCP API key business logic.
type Service struct {
	repo      *Repository
	dsService *datasource.Service
}

// NewService creates a new MCP service.
func NewService(repo *Repository, dsService *datasource.Service) *Service {
	return &Service{repo: repo, dsService: dsService}
}

// CreateKey creates a new MCP API key and returns the plain key.
func (s *Service) CreateKey(ctx context.Context, orgID, userID uuid.UUID, req CreateKeyRequest) (*CreateKeyResponse, error) {
	name := strings.TrimSpace(req.Name)
	if name == "" {
		return nil, ErrKeyNameEmpty
	}

	// Validate at least one data source is selected
	if len(req.DataSourceIDs) == 0 {
		return nil, ErrNoDataSourcesSelected
	}

	// Validate all data sources belong to the organization
	for _, dsID := range req.DataSourceIDs {
		_, err := s.dsService.GetDataSource(ctx, orgID, dsID)
		if err != nil {
			return nil, ErrInvalidDataSource
		}
	}

	plainKey, keyHash, err := generateAPIKey()
	if err != nil {
		return nil, fmt.Errorf("failed to generate API key: %w", err)
	}

	key, err := s.repo.Create(ctx, orgID, name, keyHash, userID, req.DataSourceIDs)
	if err != nil {
		return nil, fmt.Errorf("failed to create MCP API key: %w", err)
	}

	return &CreateKeyResponse{
		Key:    *key,
		APIKey: plainKey,
	}, nil
}

// ListKeys returns all MCP API keys for an organization.
func (s *Service) ListKeys(ctx context.Context, orgID uuid.UUID) ([]MCPAPIKey, error) {
	keys, err := s.repo.GetByOrganizationID(ctx, orgID)
	if err != nil {
		return nil, fmt.Errorf("failed to list MCP API keys: %w", err)
	}

	if keys == nil {
		keys = []MCPAPIKey{}
	}

	return keys, nil
}

// DeleteKey deletes an MCP API key after verifying organization ownership.
func (s *Service) DeleteKey(ctx context.Context, orgID, keyID uuid.UUID) error {
	key, err := s.repo.GetByID(ctx, keyID)
	if err != nil {
		return fmt.Errorf("failed to get MCP API key: %w", err)
	}
	if key == nil {
		return ErrKeyNotFound
	}
	if key.OrganizationID != orgID {
		return ErrUnauthorized
	}

	if err := s.repo.Delete(ctx, keyID); err != nil {
		return fmt.Errorf("failed to delete MCP API key: %w", err)
	}

	return nil
}

// UpdateKey updates the data sources for an MCP API key.
func (s *Service) UpdateKey(ctx context.Context, orgID, keyID uuid.UUID, req UpdateKeyRequest) (*MCPAPIKey, error) {
	// Validate at least one data source is selected
	if len(req.DataSourceIDs) == 0 {
		return nil, ErrNoDataSourcesSelected
	}

	// Get the key and verify ownership
	key, err := s.repo.GetByID(ctx, keyID)
	if err != nil {
		return nil, fmt.Errorf("failed to get MCP API key: %w", err)
	}
	if key == nil {
		return nil, ErrKeyNotFound
	}
	if key.OrganizationID != orgID {
		return nil, ErrUnauthorized
	}

	// Validate all data sources belong to the organization
	for _, dsID := range req.DataSourceIDs {
		_, err := s.dsService.GetDataSource(ctx, orgID, dsID)
		if err != nil {
			return nil, ErrInvalidDataSource
		}
	}

	// Update the data sources
	if err := s.repo.UpdateDataSources(ctx, keyID, req.DataSourceIDs); err != nil {
		return nil, fmt.Errorf("failed to update MCP API key: %w", err)
	}

	// Return the updated key
	key.AllowedDataSourceIDs = req.DataSourceIDs
	return key, nil
}

// ValidateKey validates an API key and returns the associated key record.
func (s *Service) ValidateKey(ctx context.Context, apiKey string) (*MCPAPIKey, error) {
	if apiKey == "" {
		return nil, ErrKeyNotFound
	}

	keyHash := hashAPIKey(apiKey)
	key, err := s.repo.GetByAPIKeyHash(ctx, keyHash)
	if err != nil {
		return nil, fmt.Errorf("failed to validate API key: %w", err)
	}
	if key == nil {
		return nil, ErrKeyNotFound
	}

	// Update last used timestamp asynchronously (fire and forget)
	go func() {
		_ = s.repo.UpdateLastUsed(context.Background(), key.ID)
	}()

	return key, nil
}

// generateAPIKey generates a new API key and its hash.
func generateAPIKey() (plainKey, hash string, err error) {
	bytes := make([]byte, apiKeyBytes)
	if _, err := rand.Read(bytes); err != nil {
		return "", "", err
	}

	plainKey = apiKeyPrefix + base64.URLEncoding.EncodeToString(bytes)
	hash = hashAPIKey(plainKey)

	return plainKey, hash, nil
}

// hashAPIKey hashes an API key using SHA-256.
func hashAPIKey(apiKey string) string {
	hashBytes := sha256.Sum256([]byte(apiKey))
	return hex.EncodeToString(hashBytes[:])
}
