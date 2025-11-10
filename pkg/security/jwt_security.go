package security

import (
	"auth-management/pkg/dto"
	"auth-management/pkg/enum"
	"fmt"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

func JwtGenerateAccessToken(id string, role enum.ROLE, secret []byte) (string, error) {
	claims := dto.TokenModel{
		Role: role,
		RegisteredClaims: jwt.RegisteredClaims{
			ID:        uuid.NewString(),
			Subject:   id,
			Issuer:    "auth management",
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(15 * time.Minute)),
		},
	}
	token, err := jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString(secret)
	if err != nil {
		return "", err
	}
	return token, nil
}
func JwtGenerateRefreshToken(id string, role enum.ROLE, secret []byte) (string, int32, error) {
	exp := jwt.NewNumericDate(time.Now().Add((24 * time.Hour) * 7))
	claims := dto.TokenModel{
		Role: role,
		RegisteredClaims: jwt.RegisteredClaims{
			ID:        uuid.NewString(),
			Subject:   id,
			Issuer:    "auth management",
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			ExpiresAt: exp,
		},
	}
	token, err := jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString(secret)
	if err != nil {
		return "", 0, err
	}
	return token, int32(exp.Time.Unix()), nil
}
func JwtVerify(tokenString string) (*dto.TokenModel, error) {
	secret := []byte(os.Getenv("JWT_SECRET"))
	token, err := jwt.ParseWithClaims(tokenString, &dto.TokenModel{}, func(t *jwt.Token) (any, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", t.Header["alg"])
		}
		return secret, nil
	})
	if err != nil {
		return nil, err
	} else if claims, ok := token.Claims.(*dto.TokenModel); ok && token.Valid {
		return claims, nil
	} else {
		return nil, fmt.Errorf("unknown claim type, cannot processed: %v", err)
	}
}
