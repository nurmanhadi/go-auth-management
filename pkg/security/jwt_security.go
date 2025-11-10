package security

import (
	"auth-management/pkg/dto"
	"auth-management/pkg/enum"
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
