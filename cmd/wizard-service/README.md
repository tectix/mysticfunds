# Wizard Service

The Wizard Service is a gRPC-based microservice that handles wizard-related operations for the MysticFunds project.

## Features

- Create, retrieve, update, and delete wizards
- List wizards with pagination
- Join and leave guilds
- Manage wizard attributes (name, realm, element)
- Track wizard mana balance

## Prerequisites

- Go 1.16 or later
- PostgreSQL
- Protocol Buffers compiler (protoc)
- [golang-migrate](https://github.com/golang-migrate/migrate) for database migrations

## Configuration

The service uses a `config.yaml` file for configuration. Here's an example of the configuration:

```yaml
SERVICE_NAME: wizard-service
GRPC_PORT: 50052
LOG_LEVEL: info
JWT_SECRET: your_jwt_secret_here

DB_HOST: localhost
DB_PORT: 5432
DB_USER: postgres
DB_PASSWORD: password
DB_NAME: wizard
```

Ensure all required fields are set correctly in this file before running the service.

## Setup

1. Ensure you're in the `cmd/wizard-service/` directory.
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

The Wizard Service provides the following gRPC endpoints:

1. `CreateWizard`: Create a new wizard
2. `GetWizard`: Retrieve details of a specific wizard
3. `UpdateWizard`: Update information for an existing wizard
4. `ListWizards`: Retrieve a list of wizards with pagination
5. `DeleteWizard`: Delete a wizard
6. `JoinGuild`: Make a wizard join a guild
7. `LeaveGuild`: Make a wizard leave their current guild

For detailed API documentation, refer to the `proto/wizard/wizard.proto` file.

## Development

### Regenerating gRPC Code

If you make changes to the `wizard.proto` file, regenerate the gRPC code by running:

```
make proto
```

## Troubleshooting

- If you encounter database connection issues, make sure your PostgreSQL server is running and the connection details in `config.yaml` are correct.
- For "connection refused" errors, check if the wizard service is running and listening on the expected port (50052 by default).

