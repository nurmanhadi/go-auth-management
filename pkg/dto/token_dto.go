package dto

import (
	"auth-management/pkg/enum"

	"github.com/golang-jwt/jwt/v5"
)

type TokenResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}
type TokenModel struct {
	Role enum.ROLE `json:"role"`
	jwt.RegisteredClaims
}
