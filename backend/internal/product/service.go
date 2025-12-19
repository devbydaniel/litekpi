package product

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

// Service handles product business logic.
type Service struct {
	repo *Repository
}

// NewService creates a new product service.
func NewService(repo *Repository) *Service {
	return &Service{repo: repo}
}

// CreateProduct creates a new product and returns the plain API key.
func (s *Service) CreateProduct(ctx context.Context, orgID uuid.UUID, req CreateProductRequest) (*CreateProductResponse, error) {
	name := strings.TrimSpace(req.Name)
	if name == "" {
		return nil, ErrProductNameEmpty
	}

	plainKey, keyHash, err := generateAPIKey()
	if err != nil {
		return nil, fmt.Errorf("failed to generate API key: %w", err)
	}

	product, err := s.repo.CreateProduct(ctx, orgID, name, keyHash)
	if err != nil {
		return nil, fmt.Errorf("failed to create product: %w", err)
	}

	return &CreateProductResponse{
		Product: *product,
		APIKey:  plainKey,
	}, nil
}

// ListProducts returns all products for an organization.
func (s *Service) ListProducts(ctx context.Context, orgID uuid.UUID) ([]Product, error) {
	products, err := s.repo.GetProductsByOrganizationID(ctx, orgID)
	if err != nil {
		return nil, fmt.Errorf("failed to list products: %w", err)
	}

	if products == nil {
		products = []Product{}
	}

	return products, nil
}

// DeleteProduct deletes a product after verifying organization ownership.
func (s *Service) DeleteProduct(ctx context.Context, orgID, productID uuid.UUID) error {
	product, err := s.repo.GetProductByID(ctx, productID)
	if err != nil {
		return fmt.Errorf("failed to get product: %w", err)
	}
	if product == nil {
		return ErrProductNotFound
	}
	if product.OrganizationID != orgID {
		return ErrUnauthorized
	}

	if err := s.repo.DeleteProduct(ctx, productID); err != nil {
		return fmt.Errorf("failed to delete product: %w", err)
	}

	return nil
}

// RegenerateAPIKey regenerates the API key for a product.
func (s *Service) RegenerateAPIKey(ctx context.Context, orgID, productID uuid.UUID) (*RegenerateKeyResponse, error) {
	product, err := s.repo.GetProductByID(ctx, productID)
	if err != nil {
		return nil, fmt.Errorf("failed to get product: %w", err)
	}
	if product == nil {
		return nil, ErrProductNotFound
	}
	if product.OrganizationID != orgID {
		return nil, ErrUnauthorized
	}

	plainKey, keyHash, err := generateAPIKey()
	if err != nil {
		return nil, fmt.Errorf("failed to generate API key: %w", err)
	}

	if err := s.repo.UpdateAPIKeyHash(ctx, productID, keyHash); err != nil {
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
