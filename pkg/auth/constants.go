package auth

type ContextKey string

const (
	// UserIDKey is the key used to store and retrieve the user ID from the context
	UserIDKey ContextKey = "user_id"

	// AuthorizationHeader is the key for the authorization header in the metadata
	AuthorizationHeader string = "authorization"

	// BearerSchema is the prefix for the bearer token
	BearerSchema string = "Bearer "
)

// gRPC method names
const (
	MethodRegister string = "/auth.AuthService/Register"
	MethodLogin    string = "/auth.AuthService/Login"
)
