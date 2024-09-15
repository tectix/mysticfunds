# Auth Service

The Auth Service is a gRPC-based microservice that handles user authentication and authorization for the MysticFunds project.

## Features

- User registration
- User login
- Token refresh
- Token validation
- User logout

## Prerequisites

- Go 1.16 or later
- PostgreSQL
- Protocol Buffers compiler (protoc)
- [golang-migrate](https://github.com/golang-migrate/migrate) for database migrations

## Configuration

The service uses a `config.yaml` file for configuration. Here's an example of the configuration:

```yaml
SERVICE_NAME: auth-service
GRPC_PORT: 50051
LOG_LEVEL: info
JWT_SECRET: your_jwt_secret_here

DB_HOST: localhost
DB_PORT: 5432
DB_USER: postgres
DB_PASSWORD: password
DB_NAME: auth
```

Ensure all required fields are set correctly in this file before running the service.

## Setup

1. Ensure you're in the `cmd/auth-service/` directory.
2. Review and update the `config.yaml` file as needed.
3. Run database migrations:
   ```
   make migrate
   ```

## Building

To build the service, run:

```
make build
```

## Running

To start the service, run:

```
make run
```

## Testing

To run the tests for this service:

```
make test
```

## Database Migrations

To run migrations:
```
make migrate
```

To rollback the last migration:
```
make migrate-down
```

## API

The Auth Service provides the following gRPC endpoints:

1. `Register`: Register a new user
2. `Login`: Authenticate a user
3. `RefreshToken`: Get a new token
4. `ValidateToken`: Check if a token is valid
5. `Logout`: Invalidate a token

For detailed API documentation, refer to the `proto/auth/auth.proto` file.

## Development

### Regenerating gRPC Code

If you make changes to the `auth.proto` file, regenerate the gRPC code by running:

```
make proto
```

## Integrating with Other Services

Other services can use the Auth Service for user authentication by:

1. Including the auth service protobuf definitions
2. Creating a gRPC client to communicate with the auth service
3. Using the `ValidateToken` method to check if a user's token is valid before processing requests

## Troubleshooting

- If you encounter database connection issues, make sure your PostgreSQL server is running and the connection details in `config.yaml` are correct.
- For "connection refused" errors, check if the auth service is running and listening on the expected port (50051 by default).
- If token validation fails, ensure that the `JWT_SECRET` is consistent across all services.
