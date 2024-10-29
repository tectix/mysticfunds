package middleware

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"google.golang.org/grpc"
)

type Config struct {
	JWTSecret string
	AuthConn  *grpc.ClientConn
}

func Setup(r chi.Router, config *Config) {
	publicPaths := []string{
		"/api/v1/auth/login",
		"/api/v1/auth/register",
		"/health",
		"/metrics",
	}
	authMiddleware := NewAuthMiddleware(config.JWTSecret, config.AuthConn, publicPaths)

	r.Use(authMiddleware.Handler)

	r.Group(func(r chi.Router) {
		r.Use(func(next http.Handler) http.Handler {
			return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				if IsPublicPath(r.URL.Path) {
					next.ServeHTTP(w, r)
					return
				}
				authMiddleware.Handler(next).ServeHTTP(w, r)
			})
		})
	})
}
