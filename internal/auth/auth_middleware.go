package auth

import (
	"context"
	"net/http"
	"strings"

	"github.com/rs/zerolog/log"
)

type AuthMiddleware struct {
	JwtService JwtService
}

func NewAuthMiddleware(jwtService JwtService) AuthMiddleware {
	return AuthMiddleware{
		JwtService: jwtService,
	}
}

type usernameContextKey int

const UsernameContextKey usernameContextKey = 0

func (h AuthMiddleware) Authenticate(next http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		const bearerScheme string = "Bearer "

		auth := r.Header.Get("X-Forwarded-Authorization")

		if !strings.HasPrefix(auth, bearerScheme) {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		token := auth[len(bearerScheme):]
		claims, err := h.JwtService.GetClaims(token)
		if err != nil {
			log.Error().Err(err).Msg("Error getting JWT Token claims")
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		username := claims.Subject
		ctxWithUsername := context.WithValue(r.Context(), UsernameContextKey, username)
		rWithUsername := r.WithContext(ctxWithUsername)
		next.ServeHTTP(w, rWithUsername)
	})
}
