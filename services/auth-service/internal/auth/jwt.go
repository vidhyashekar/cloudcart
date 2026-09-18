package auth

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

// JWTManager is responsible for generating and validating JWT tokens.
type JWTManager struct {
	SecretKey string
	ExpiresIn time.Duration
}

// Claims represents the JWT claims.
type Claims struct {
	UserID uint   `json:"user_id"`
	Email  string `json:"email"`
	Role   string `json:"role"`
	jwt.RegisteredClaims
}

// NewJWTManager creates a new instance of JWTManager with the provided secret key and expiration duration.
func NewJWTManager(secretKey string, expiresIn time.Duration) *JWTManager {
	return &JWTManager{
		SecretKey: secretKey,
		ExpiresIn: expiresIn,
	}
}

// GenerateToken generates a JWT token for the given user ID, email, and role.
func (m *JWTManager) GenerateToken(
	userID uint,
	email string,
	role string,
) (string, error) {

	claims := Claims{
		UserID: userID,
		Email:  email,
		Role:   role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(m.ExpiresIn)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	return token.SignedString([]byte(m.SecretKey))
}

func (m *JWTManager) ValidateToken(tokenString string) (*Claims, error) {
	token, err := jwt.ParseWithClaims(
		tokenString,
		&Claims{},
		func(token *jwt.Token) (interface{}, error) {
			if token.Method != jwt.SigningMethodHS256 {
				return nil, errors.New("unexpected signing method")
			}

			return []byte(m.SecretKey), nil
		},
	)

	if err != nil {
		return nil, err
	}

	claims, ok := token.Claims.(*Claims)
	if !ok || !token.Valid {
		return nil, errors.New("invalid token")
	}

	return claims, nil
}
