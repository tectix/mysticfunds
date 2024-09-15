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

## Setup

1. Clone the repository:
   ```
   git clone https://github.com/Alinoureddine1/mysticfunds.git
   cd mysticfunds/cmd/auth-service
   ```

2. Install dependencies:
   ```
   go mod tidy
   ```

3. Set up the database:
   - Create a PostgreSQL database for the auth service
   - Run the migrations (refer to the main project README for migration instructions)

4. Configure the service:
   - Copy `config.yaml.example` to `config.yaml`
   - Edit `config.yaml` and set the appropriate values for your environment

## Building

To build the service, run:

```
go build -o auth-service
```

## Running

To start the service, run:

```
./auth-service
```

By default, the service will listen on port 50051. You can change this in the `config.yaml` file.

## Usage

The Auth Service provides the following gRPC endpoints:

1. `Register`: Register a new user
   - Input: username, email, password
   - Output: JWT token, user ID

2. `Login`: Authenticate a user
   - Input: username, password
   - Output: JWT token, user ID

3. `RefreshToken`: Get a new token
   - Input: existing JWT token
   - Output: new JWT token, user ID

4. `ValidateToken`: Check if a token is valid
   - Input: JWT token
   - Output: validation status, user ID

5. `Logout`: Invalidate a token
   - Input: JWT token
   - Output: success status

To interact with the service, you need to use a gRPC client. You can find an example client in `cmd/auth-client/main.go`.

## Configuration

The service uses a `config.yaml` file for configuration. Here's an example:

```yaml
SERVICE_NAME: auth-service
GRPC_PORT: 50051
LOG_LEVEL: info
JWT_SECRET: your_jwt_secret_here

DB_HOST: localhost
DB_PORT: 5432
DB_USER: postgres
DB_PASSWORD: your_password
DB_NAME: auth_db
```

## Development

### Regenerating gRPC code

If you make changes to the `auth.proto` file, you need to regenerate the gRPC code. From the project root, run:

```
protoc --go_out=. --go_opt=paths=source_relative --go-grpc_out=. --go-grpc_opt=paths=source_relative proto/auth/auth.proto
```

## Integrating with Other Services

Other services can use the Auth Service for user authentication by:

1. Including the auth service protobuf definitions
2. Creating a gRPC client to communicate with the auth service
3. Using the `ValidateToken` method to check if a user's token is valid before processing requests

## Troubleshooting

- If you encounter database connection issues, make sure your PostgreSQL server is running and the connection details in `config.yaml` are correct.
- For "connection refused" errors, check if the auth service is running and listening on the expected port.
- If token validation fails, ensure that the `JWT_SECRET` is consistent across all services.
