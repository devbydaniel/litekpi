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
)

const (
	apiKeyPrefix = "lkmcp_"
	apiKeyBytes  = 32
)

// Service handles MCP API key business logic.
type Service struct {
	repo *Repository
}

// NewService creates a new MCP service.
func NewService(repo *Repository) *Service {
	return &Service{repo: repo}
}

// CreateKey creates a new MCP API key and returns the plain key.
func (s *Service) CreateKey(ctx context.Context, orgID, userID uuid.UUID, req CreateKeyRequest) (*CreateKeyResponse, error) {
	name := strings.TrimSpace(req.Name)
	if name == "" {
		return nil, ErrKeyNameEmpty
	}

	plainKey, keyHash, err := generateAPIKey()
	if err != nil {
		return nil, fmt.Errorf("failed to generate API key: %w", err)
	}

	key, err := s.repo.Create(ctx, orgID, name, keyHash, userID)
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
