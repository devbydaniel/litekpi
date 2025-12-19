package auth

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

var (
	ErrInvalidToken = errors.New("invalid or expired token")
)

const (
	oauthSetupTokenExpiry = 15 * time.Minute
)

// JWTService handles JWT token operations.
type JWTService struct {
	secret     []byte
	expiration time.Duration
}

// NewJWTService creates a new JWT service.
func NewJWTService(secret string) *JWTService {
	return &JWTService{
		secret:     []byte(secret),
		expiration: 7 * 24 * time.Hour, // 7 days
	}
}

// Claims represents the JWT claims.
type Claims struct {
	UserID         string `json:"userId"`
	Email          string `json:"email"`
	OrganizationID string `json:"organizationId"`
	Role           string `json:"role"`
	jwt.RegisteredClaims
}

// OAuthSetupClaims represents claims for OAuth setup tokens.
type OAuthSetupClaims struct {
	OAuthAccountID string `json:"oauthAccountId"`
	jwt.RegisteredClaims
}

// GenerateToken generates a new JWT token for a user.
func (j *JWTService) GenerateToken(userID uuid.UUID, email string, organizationID uuid.UUID, role Role) (string, error) {
	claims := &Claims{
		UserID:         userID.String(),
		Email:          email,
		OrganizationID: organizationID.String(),
		Role:           string(role),
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(j.expiration)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
			Issuer:    "litekpi",
			Subject:   userID.String(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(j.secret)
}

// ValidateToken validates a JWT token and returns the claims.
func (j *JWTService) ValidateToken(tokenString string) (*Claims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, ErrInvalidToken
		}
		return j.secret, nil
	})

	if err != nil {
		return nil, ErrInvalidToken
	}

	claims, ok := token.Claims.(*Claims)
	if !ok || !token.Valid {
		return nil, ErrInvalidToken
	}

	return claims, nil
}

// GenerateOAuthSetupToken generates a short-lived token for OAuth setup completion.
func (j *JWTService) GenerateOAuthSetupToken(oauthAccountID uuid.UUID) (string, error) {
	claims := &OAuthSetupClaims{
		OAuthAccountID: oauthAccountID.String(),
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(oauthSetupTokenExpiry)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
			Issuer:    "litekpi",
			Subject:   "oauth-setup",
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(j.secret)
}

// ValidateOAuthSetupToken validates an OAuth setup token and returns the OAuth account ID.
func (j *JWTService) ValidateOAuthSetupToken(tokenString string) (uuid.UUID, error) {
	token, err := jwt.ParseWithClaims(tokenString, &OAuthSetupClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return uuid.Nil, ErrInvalidToken
		}
		return j.secret, nil
	})

	if err != nil {
		return uuid.Nil, ErrInvalidToken
	}

	claims, ok := token.Claims.(*OAuthSetupClaims)
	if !ok || !token.Valid {
		return uuid.Nil, ErrInvalidToken
	}

	oauthAccountID, err := uuid.Parse(claims.OAuthAccountID)
	if err != nil {
		return uuid.Nil, ErrInvalidToken
	}

	return oauthAccountID, nil
}
