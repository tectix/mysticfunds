package auth

import (
	"context"
	"strings"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

type AuthInterceptor struct {
	jwtSecretKey  string
	publicMethods map[string]bool
}

func NewAuthInterceptor(jwtSecretKey string, publicMethods []string) *AuthInterceptor {
	public := make(map[string]bool)
	for _, method := range publicMethods {
		public[method] = true
	}

	return &AuthInterceptor{
		jwtSecretKey:  jwtSecretKey,
		publicMethods: public,
	}
}

func (interceptor *AuthInterceptor) Unary() grpc.UnaryServerInterceptor {
	return func(
		ctx context.Context,
		req interface{},
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler,
	) (interface{}, error) {
		if !interceptor.publicMethods[info.FullMethod] {
			var err error
			ctx, err = interceptor.authorize(ctx)
			if err != nil {
				return nil, err
			}
		}

		return handler(ctx, req)
	}
}

func (interceptor *AuthInterceptor) authorize(ctx context.Context) (context.Context, error) {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return ctx, status.Errorf(codes.Unauthenticated, "metadata is not provided")
	}

	values := md[AuthorizationHeader]
	if len(values) == 0 {
		return ctx, status.Errorf(codes.Unauthenticated, "authorization token is not provided")
	}

	accessToken := values[0]
	if !strings.HasPrefix(accessToken, BearerSchema) {
		return ctx, status.Errorf(codes.Unauthenticated, "invalid authorization token format")
	}

	tokenString := strings.TrimPrefix(accessToken, BearerSchema)

	claims, err := ValidateToken(tokenString, interceptor.jwtSecretKey)
	if err != nil {
		return ctx, status.Errorf(codes.Unauthenticated, "invalid token: %v", err)
	}

	return context.WithValue(ctx, UserIDKey, claims.UserID), nil
}
