package middleware

import (
	"context"
	"net/http"
	"strings"

	"github.com/Alinoureddine1/mysticfunds/pkg/auth"
	authpb "github.com/Alinoureddine1/mysticfunds/proto/auth"
	"google.golang.org/grpc"
)

type AuthMiddleware struct {
	jwtSecret   string
	authClient  authpb.AuthServiceClient
	publicPaths map[string]bool
}

func NewAuthMiddleware(jwtSecret string, authConn *grpc.ClientConn, publicPaths []string) *AuthMiddleware {
	paths := make(map[string]bool)
	for _, path := range publicPaths {
		paths[path] = true
	}

	return &AuthMiddleware{
		jwtSecret:   jwtSecret,
		authClient:  authpb.NewAuthServiceClient(authConn),
		publicPaths: paths,
	}
}

// Middleware function to authenticate requests
func (m *AuthMiddleware) Handler(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Check if the path is public
		if m.publicPaths[r.URL.Path] {
			next.ServeHTTP(w, r)
			return
		}

		// Get token from Authorization header
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			http.Error(w, "Authorization header is required", http.StatusUnauthorized)
			return
		}

		// Check Bearer scheme
		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			http.Error(w, "Invalid authorization header format", http.StatusUnauthorized)
			return
		}

		tokenString := parts[1]

		// Validate token locally
		claims, err := auth.ValidateToken(tokenString, m.jwtSecret)
		if err != nil {
			switch err {
			case auth.ErrExpiredToken:
				http.Error(w, "Token has expired", http.StatusUnauthorized)
			case auth.ErrInvalidToken:
				http.Error(w, "Invalid token", http.StatusUnauthorized)
			default:
				http.Error(w, "Authentication failed", http.StatusUnauthorized)
			}
			return
		}

		// Validate token with auth service
		resp, err := m.authClient.ValidateToken(r.Context(), &authpb.ValidateTokenRequest{
			Token: tokenString,
		})
		if err != nil {
			http.Error(w, "Failed to validate token with auth service", http.StatusUnauthorized)
			return
		}

		if !resp.IsValid {
			http.Error(w, "Invalid token", http.StatusUnauthorized)
			return
		}

		// Add user information to request context
		ctx := context.WithValue(r.Context(), UserIDKey, claims.UserID)
		r = r.WithContext(ctx)

		next.ServeHTTP(w, r)
	})
}

type contextKey string

const (
	UserIDKey contextKey = "userID"
)

func GetUserIDFromContext(ctx context.Context) (int64, bool) {
	userID, ok := ctx.Value(UserIDKey).(int64)
	return userID, ok
}

func IsPublicPath(path string) bool {
	publicPaths := map[string]bool{
		"/api/v1/auth/login":    true,
		"/api/v1/auth/register": true,
		"/health":               true,
		"/metrics":              true,
	}
	return publicPaths[path]
}

func NewAuthClient(conn *grpc.ClientConn) authpb.AuthServiceClient {
	return authpb.NewAuthServiceClient(conn)
}

type AuthError struct {
	Message string
	Code    int
}

func (e *AuthError) Error() string {
	return e.Message
}

var (
	ErrMissingToken          = &AuthError{Message: "Missing authentication token", Code: http.StatusUnauthorized}
	ErrInvalidTokenFormat    = &AuthError{Message: "Invalid token format", Code: http.StatusUnauthorized}
	ErrTokenValidationFailed = &AuthError{Message: "Token validation failed", Code: http.StatusUnauthorized}
	ErrServiceUnavailable    = &AuthError{Message: "Authentication service unavailable", Code: http.StatusServiceUnavailable}
)
