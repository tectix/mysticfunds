# Mana Service

The Mana Service is a gRPC-based microservice that handles mana (currency) transactions and investments for the MysticFunds project.

## Features

- Transfer mana between wizards
- Check mana balances
- List transaction history with pagination
- Investment management
  - Create investments
  - Track investment returns
  - Automated investment completion
  - Risk-based return calculations
- Investment types with different risk levels and returns
- Scheduled investment processing

## Prerequisites

- Go 1.16 or later
- PostgreSQL
- Protocol Buffers compiler (protoc)
- [golang-migrate](https://github.com/golang-migrate/migrate) for database migrations

## Configuration

The service uses a `config.yaml` file for configuration. Copy `config.example.yaml` to `config.yaml` and update the values accordingly.

## Setup

1. Ensure you're in the `cmd/mana-service/` directory.
2. Copy `config.example.yaml` to `config.yaml` and update it.
3. Run database migrations:
   ```
   make migrate-up
   ```

## Building

To build the service:
```
make build
```

## Running

To start the service:
```
make run
```

## Testing

To run the tests:
```
make test
```

## Database Migrations

- Run migrations: `make migrate-up`
- Rollback last migration: `make migrate-down`
- Check migration status: `make migration-status`

## API

The Mana Service provides the following gRPC endpoints:

1. `TransferMana`: Transfer mana between wizards
2. `GetManaBalance`: Get a wizard's current mana balance
3. `ListTransactions`: List mana transactions with pagination
4. `CreateInvestment`: Create a new investment
5. `GetInvestments`: Get a wizard's investments
6. `GetInvestmentTypes`: List available investment types

For detailed API documentation, refer to `proto/mana/mana.proto`.

## Investment System

### Investment Types
- Novice Spell Bond (Risk Level 1): Low-risk, short-term
- Mystic Market Fund (Risk Level 2): Balanced investment
- Elemental Ventures (Risk Level 3): Higher risk/return
- Dragon's Hoard (Risk Level 4): High-risk, long-term
- Phoenix Rising (Risk Level 5): Maximum risk/return

### Risk and Returns
- Each risk level affects potential return variance
- Higher risk levels have greater potential returns but also higher volatility
- Automated return calculation based on risk level and market conditions

## Development

### Regenerating gRPC Code
After modifying the proto file:
```
make proto
```

## Troubleshooting

- For database connection issues, verify PostgreSQL is running and config.yaml settings
- Check logs for investment scheduling issues
- Ensure sufficient mana balance for transactions and investments