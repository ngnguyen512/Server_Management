package jwtha

import (
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v4"
)

// JwtConfig holds the necessary fields for JwtTokenizer
type JwtConfig struct {
	SecretKey     []byte
	SigningMethod jwt.SigningMethod
	Expiration    time.Duration
}

// JwtTokenizer implements the ITokenizer interface using JWT
type JwtTokenizer struct {
	config JwtConfig
}

// NewJwtTokenizer initializes a new JwtTokenizer with the given config
func NewJwtTokenizer(config JwtConfig) *JwtTokenizer {
	return &JwtTokenizer{config: config}
}

// Gencode generates a token for the given input data
func (t *JwtTokenizer) Gencode(data map[string]interface{}) (string, error) {
	token := jwt.NewWithClaims(t.config.SigningMethod, jwt.MapClaims{
		"data": data,
		"exp":  time.Now().Add(t.config.Expiration).Unix(),
	})
	tokenString, err := token.SignedString(t.config.SecretKey)
	if err != nil {
		return "", err
	}
	return tokenString, nil
}

// Parsecode parses the given token string and returns the data
func (t *JwtTokenizer) Parsecode(tokenString string) (map[string]interface{}, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return t.config.SecretKey, nil
	})
	if err != nil {
		return nil, err
	}
	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		return claims["data"].(map[string]interface{}), nil
	} else {
		return nil, fmt.Errorf("invalid token")
	}
}
