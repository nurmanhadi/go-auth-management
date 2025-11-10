package middleware

import (
	"auth-management/pkg/response"
	"auth-management/pkg/security"
	"context"
	"net/http"
	"strings"
)

func JwtVerify(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		header := r.Header.Get("Authorization")
		if header == "" {
			response.Except(http.StatusUnauthorized, "token null")
			return
		}
		token := strings.Split(header, " ")
		if len(token) != 2 || token[0] != "Bearer" {
			response.Except(http.StatusUnauthorized, "value authorization most be Bearer example 'Authorization: Bearer T'")
			return
		}
		claim, err := security.JwtVerify(token[1])
		if err != nil {
			response.Except(http.StatusUnauthorized, err.Error())
			return
		}
		ctx := context.WithValue(r.Context(), "role", claim.Role)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
