package middleware

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	"github.com/golang-jwt/jwt/v5"
)

type userIDKey struct{}

var publicPaths = map[string]struct{}{
	"/v1/auth/register": {},
	"/v1/auth/login":    {},
}

func AuthMiddleware(next http.Handler, secret []byte) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if _, ok := publicPaths[r.URL.Path]; ok {
			next.ServeHTTP(w, r)
			return
		}

		auth := r.Header.Get("Authorization")
		if !strings.HasPrefix(auth, "Bearer ") {
			http.Error(w, "missing token", http.StatusUnauthorized)
			return
		}
		tokenStr := strings.TrimPrefix(auth, "Bearer ")

		claims := jwt.MapClaims{}
		token, err := jwt.ParseWithClaims(tokenStr, claims, func(t *jwt.Token) (interface{}, error) {
			if t.Method.Alg() != jwt.SigningMethodHS256.Alg() {
				return nil, fmt.Errorf("unexpected alg")
			}
			return secret, nil
		})
		if err != nil || !token.Valid {
			http.Error(w, "invalid token", http.StatusUnauthorized)
			return
		}

		if claims["iss"] != "auth-service" || claims["aud"] != "gateway" {
			http.Error(w, "invalid claims", http.StatusUnauthorized)
			return
		}

		sub := fmt.Sprintf("%v", claims["sub"])
		ctx := context.WithValue(r.Context(), userIDKey{}, sub)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
