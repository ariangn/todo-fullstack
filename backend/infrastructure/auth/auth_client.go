package auth

import (
    "errors"
    "os"
    "time"

    "github.com/golang-jwt/jwt/v4"
)

// defines methods for generating/verifying JWTs
type AuthClientInterface interface {
    GenerateToken(userID string, ttl time.Duration) (string, error)
    ValidateToken(tokenString string) (jwt.MapClaims, error)
}

type AuthClient struct {
    secretKey []byte
}

func NewAuthClient() AuthClientInterface {
    secret := os.Getenv("JWT_SECRET")
    if secret == "" {
        panic("JWT_SECRET must be set in environment")
    }
    return &AuthClient{secretKey: []byte(secret)}
}

func (a *AuthClient) GenerateToken(userID string, ttl time.Duration) (string, error) {
    claims := jwt.MapClaims{
        "sub": userID,
        "exp": time.Now().Add(ttl).Unix(),
    }
    token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
    return token.SignedString(a.secretKey)
}

func (a *AuthClient) ValidateToken(tokenString string) (jwt.MapClaims, error) {
    token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
        if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
            return nil, errors.New("unexpected signing method")
        }
        return a.secretKey, nil
    })
    if err != nil {
        return nil, err
    }
    if !token.Valid {
        return nil, errors.New("invalid token")
    }
    claims, ok := token.Claims.(jwt.MapClaims)
    if !ok {
        return nil, errors.New("could not parse claims")
    }
    return claims, nil
}
