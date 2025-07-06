# MysticFunds - Magical Investment Platform

A magical investment platform built with Go microservices and a modern web frontend. MysticFunds allows wizards to manage their mana, create investments, and track their magical portfolio growth.

![MysticFunds Dashboard](https://img.shields.io/badge/Status-Ready_for_Magic-purple?style=for-the-badge&logo=magic)

## Quick Start

### Docker (Recommended)

The easiest way to get MysticFunds running:

```bash
# Clone the repository
git clone https://github.com/Alinoureddine1/mysticfunds.git
cd mysticfunds

# Start everything with Docker
docker compose up -d

# Open http://localhost:8080 in your browser and start investing!
```

### Local Development

```bash
# Clone the repository
git clone https://github.com/Alinoureddine1/mysticfunds.git
cd mysticfunds

# Start everything (databases, services, frontend)
make start

# Open http://localhost:8080 in your browser and start investing!
```

The startup command will show you the service URLs:
```
MysticFunds is running!
Web Interface: http://localhost:8080
Auth Service:  localhost:50051
Wizard Service: localhost:50052
Mana Service:  localhost:50053
```

The system will:
- Stop any existing services
- Create PostgreSQL user `mysticfunds` (if needed)
- Create PostgreSQL databases (auth, wizard, mana)
- Run database migrations
- Build all microservices
- Start all services in the background
- Launch the API Gateway and web frontend

**Note**: On first run, you may need to enter your system password to create the PostgreSQL user.

## Management Commands

| Command | Purpose | Usage |
|---------|---------|-------|
| `make start` | Start entire system | One-command startup |
| `make stop` | Stop all services | Graceful shutdown |
| `make status` | Check system status | Monitor health |
| `make dev` | Development mode | Auto-restart on changes |
| `make test` | Run all tests | Comprehensive test suite |

## Architecture Overview

```
┌─────────────────┐    ┌─────────────────┐
│   Web Browser   │────│   API Gateway   │
│  (Frontend UI)  │    │   (Port 8080)   │
└─────────────────┘    └─────────────────┘
                                │
                ┌───────────────┼───────────────┐
                │               │               │
        ┌───────▼────┐  ┌───────▼────┐  ┌───────▼────┐
        │    Auth    │  │   Wizard   │  │    Mana    │
        │  Service   │  │  Service   │  │  Service   │
        │(Port 50051)│  │(Port 50052)│  │(Port 50053)│
        └────────────┘  └────────────┘  └────────────┘
                │               │               │
        ┌───────▼────┐  ┌───────▼────┐  ┌───────▼────┐
        │   Auth DB  │  │ Wizard DB  │  │  Mana DB   │
        │(PostgreSQL)│  │(PostgreSQL)│  │(PostgreSQL)│
        └────────────┘  └────────────┘  └────────────┘
```

### Core Services

| Service | Port | Purpose | Database |
|---------|------|---------|----------|
| **Auth Service** | 50051 | User authentication & JWT tokens | `auth` |
| **Wizard Service** | 50052 | Wizard profiles & guild management | `wizard` |  
| **Mana Service** | 50053 | Mana transfers & investment platform | `mana` |
| **API Gateway** | 8080 | HTTP REST API & frontend hosting | - |

## Features

### Current Features

#### Core Systems
- **User Authentication** - Secure registration, login, logout with JWT
- **Wizard Management** - Create wizards with 10 different elements and 10 mystical realms
- **Mana System** - Transfer mana between wizards, track balances
- **Investment Platform** - Create investments with different risk levels and returns
- **Portfolio Dashboard** - View stats, recent activity, investment performance

#### Magical World Features
- **10 Mystical Realms** - Each with unique lore, artifacts, and characteristics
  - Pyrrhian Flame, Zepharion Heights, Terravine Hollow, Thalorion Depths, Virelya
  - Umbros, Nyxthar, Aetherion, Chronarxis, Technarok
- **10 Elemental Schools** - Fire, Water, Earth, Air, Light, Shadow, Spirit, Metal, Time, Void
- **Jobs System** - 50+ diverse magical jobs across all realms and elements
- **Marketplace** - Trade artifacts, spells, and enhancement scrolls
- **Spell Collection** - Learn and cast spells from different magical schools
- **Guild System** - Join guilds and participate in collaborative activities

#### User Experience
- **Responsive UI** - Beautiful web interface that works on all devices
- **Real-time Job Progress** - Live updates on wizard job assignments and earnings
- **Advanced Filtering** - Element-based filtering for jobs and marketplace items
- **Wizard Exploration** - Discover other wizards across the realms
- **Transaction History** - Complete audit trail of all mana movements
- **Interactive Dashboard** - Real-time notifications and balance updates

### Planned Features
- **Spell Casting System** - Active spell usage with mana costs
- **Realm Bonuses** - Element-based job and investment bonuses
- **Advanced Guild Features** - Guild wars, collaborative investments, shared rewards
- **PvP Arena** - Wizard battles and tournaments
- **Advanced Analytics** - Portfolio optimization and risk analysis
- **Achievements System** - Unlock rewards and titles
- **Seasonal Events** - Special realm events and limited-time content

## Prerequisites

### Docker Setup (Recommended)

The Docker setup handles everything automatically! You only need:

- **Docker** - [Install Docker Desktop](https://www.docker.com/products/docker-desktop/)
- **Git** - For cloning the repository

That's it! Docker will handle PostgreSQL, Go, and all dependencies.

### Local Development Setup

If you prefer to run without Docker, ensure you have:

#### Required
- **Go 1.21.5+** - [Download Go](https://golang.org/dl/)
- **PostgreSQL 12+** - [Install PostgreSQL](https://www.postgresql.org/download/)
- **Make** - Usually pre-installed on Unix systems
- **Git** - For cloning the repository

### Database Configuration
The system will automatically create a dedicated PostgreSQL user:
- **Host**: localhost
- **Port**: 5432
- **User**: mysticfunds
- **Password**: mysticfunds
- **Databases**: auth, wizard, mana (auto-created)

### Platform-Specific Setup

#### macOS
```bash
# Install PostgreSQL
brew install postgresql
brew services start postgresql

# The system will create the mysticfunds user automatically
```

#### Ubuntu/Debian
```bash
# Install PostgreSQL
sudo apt update
sudo apt install postgresql postgresql-contrib
sudo systemctl start postgresql

# The system will create the mysticfunds user automatically
```

### Optional Tools
- **golang-migrate** - The system includes migrate in the build process
- **psql** - PostgreSQL command line tool (included with PostgreSQL)

## Manual Setup (Advanced Users)

If you prefer to run commands manually or need more control:

### 1. Database Setup
```bash
# Create databases
make create-dbs

# Run migrations (if golang-migrate is installed)
make migrate-up

# Check migration status
make migration-status
```

### 2. Build Services
```bash
# Build all services
make build

# Or build individually (each service has its own Makefile)
cd cmd/auth-service && make build
cd cmd/wizard-service && make build  
cd cmd/mana-service && make build
cd cmd/api-gateway && make build
```

### 3. Start Services 
```bash
# Start all services in background
make run

# Or start individually (in separate terminals)
cd cmd/auth-service && ./bin/auth
cd cmd/wizard-service && ./bin/wizard
cd cmd/mana-service && ./bin/mana
cd cmd/api-gateway && ./bin/api-gateway
```

## Docker Setup

### Running with Docker

MysticFunds includes a complete Docker setup that handles everything automatically:

```bash
# Start everything (PostgreSQL, migrations, all services)
docker compose up -d

# View logs
docker compose logs -f

# Stop everything
docker compose down

# Stop and remove data (fresh start)
docker compose down -v
```

### What Docker Provides

The Docker setup automatically:
- **PostgreSQL Database** - Runs PostgreSQL 15 with `mysticfunds` user
- **Database Creation** - Creates `auth`, `wizard`, and `mana` databases
- **Automatic Migrations** - Runs all database migrations on startup
- **Service Orchestration** - Starts all microservices in correct order
- **Health Checks** - Monitors service health and dependencies
- **Network Configuration** - Sets up inter-service communication
- **Volume Persistence** - Persists database data between restarts

### Docker Commands

```bash
# Development workflow
docker compose up -d          # Start in background
docker compose logs -f        # Follow logs
docker compose ps             # View running services
docker compose restart        # Restart all services
docker compose down           # Stop all services

# Individual service management
docker compose restart auth-service    # Restart specific service
docker compose logs wizard-service     # View specific service logs

# Database management
docker compose exec postgres psql -U mysticfunds -d wizard  # Connect to database
docker compose down -v                                      # Reset all data
```

### Docker Architecture

```
┌─────────────────┐    ┌─────────────────┐
│   Web Browser   │────│   API Gateway   │
│  (Frontend UI)  │    │   (Port 8080)   │
└─────────────────┘    └─────────────────┘
                                │
                ┌───────────────┼───────────────┐
                │               │               │
        ┌───────▼────┐  ┌───────▼────┐  ┌───────▼────┐
        │    Auth    │  │   Wizard   │  │    Mana    │
        │  Service   │  │  Service   │  │  Service   │
        │(Port 50051)│  │(Port 50052)│  │(Port 50053)│
        └─────┬──────┘  └─────┬──────┘  └─────┬──────┘
              │               │               │
              └───────────────┼───────────────┘
                              │
                      ┌───────▼────┐
                      │ PostgreSQL │
                      │ Container  │
                      │(Port 5432) │
                      └────────────┘
```

## Configuration

Each service can be configured via YAML files in `cmd/{service}/config.yaml`:

### Database Configuration
```yaml
DB_HOST: localhost
DB_PORT: 5432
DB_USER: mysticfunds
DB_PASSWORD: mysticfunds
DB_NAME: service_name  # auth, wizard, or mana
```

### Service Configuration
```yaml
SERVICE_NAME: auth-service
GRPC_PORT: 50051
LOG_LEVEL: info
JWT_SECRET: your_jwt_secret_here
```

### API Gateway Configuration
```yaml
SERVICE_NAME: api-gateway
HTTP_PORT: 8080
LOG_LEVEL: info
```

## Development

### Project Structure
```
mysticfunds/
├── cmd/                        # Service entry points
│   ├── auth-service/           # Authentication service
│   ├── wizard-service/         # Wizard management service
│   ├── mana-service/           # Mana and investment service
│   └── api-gateway/            # HTTP REST API gateway
├── internal/                   # Business logic implementations
│   ├── auth/                   # Auth service logic
│   ├── wizard/                 # Wizard service logic
│   └── mana/                   # Mana service logic
├── pkg/                        # Shared packages
│   ├── auth/                   # JWT and auth utilities
│   ├── config/                 # Configuration management
│   ├── database/               # Database utilities
│   └── logger/                 # Logging utilities
├── proto/                      # Protocol buffer definitions
├── migrations/                 # Database migration files
├── web/                        # Frontend assets
│   ├── index.html              # Main web interface
│   └── assets/                 # CSS, JS, images
├── logs/                       # Service logs (created at runtime)
├── .github/workflows/          # CI/CD pipeline
├── .golangci.yml              # Linting configuration
├── LORE.md                     # Complete world lore and development reference
├── LICENSE                     # MIT License
└── Makefile                    # Build automation
```

### Available Make Commands

#### Building
```bash
make build              # Build all services and gateway
make build-auth         # Build auth service only
make build-wizard       # Build wizard service only
make build-mana         # Build mana service only
make build-gateway      # Build API gateway only
```

#### Running
```bash
make run               # Run all services (requires multiple terminals)
make run-auth          # Run auth service only
make run-wizard        # Run wizard service only
make run-mana          # Run mana service only
make run-gateway       # Run API gateway only
```

#### Database Management
```bash
make create-dbs        # Create all PostgreSQL databases
make migrate-up        # Apply all database migrations
make migrate-down      # Rollback database migrations
make migration-status  # Check current migration status
make nuke              # ⚠️ Drop all databases (destructive!)
```

#### Development
```bash
make proto             # Generate protobuf code from .proto files
make test              # Run all unit tests
make clean             # Remove all compiled binaries
make help              # Show all available commands
```

### Running Tests
```bash
# Run all tests
make test

# Run specific service tests
go test ./internal/auth -v
go test ./internal/wizard -v
go test ./internal/mana -v
```

### Adding New Features

1. **Update Protocol Buffers** (if needed)
   ```bash
   # Edit proto files in proto/
   make proto  # Regenerate Go code
   ```

2. **Implement Service Logic**
   ```bash
   # Add logic in internal/{service}/
   # Update service implementations
   ```

3. **Update API Gateway**
   ```bash
   # Add new routes in cmd/api-gateway/main.go
   ```

4. **Add Frontend Features**
   ```bash
   # Update web/assets/js/ files
   # Add new UI components
   ```

5. **Database Changes**
   ```bash
   # Create migration files in migrations/{service}/
   make migrate-up
   ```

## API Documentation

### Base URL
```
http://localhost:8080/api
```

### Authentication Endpoints
| Method | Endpoint | Description | Auth Required |
|--------|----------|-------------|---------------|
| POST | `/auth/register` | Register new user | No |
| POST | `/auth/login` | User login | No |
| POST | `/auth/refresh` | Refresh JWT token | No |
| POST | `/auth/logout` | User logout | Yes |

### Wizard Endpoints
| Method | Endpoint | Description | Auth Required |
|--------|----------|-------------|---------------|
| GET | `/wizards` | List all wizards | Yes |
| POST | `/wizards` | Create new wizard | Yes |
| GET | `/wizards/{id}` | Get wizard details | Yes |
| PUT | `/wizards/{id}` | Update wizard | Yes |
| DELETE | `/wizards/{id}` | Delete wizard | Yes |

### Mana Endpoints
| Method | Endpoint | Description | Auth Required |
|--------|----------|-------------|---------------|
| GET | `/mana/balance/{wizardId}` | Get mana balance | Yes |
| POST | `/mana/transfer` | Transfer mana between wizards | Yes |
| GET | `/mana/transactions/{wizardId}` | Get transaction history | Yes |

### Investment Endpoints
| Method | Endpoint | Description | Auth Required |
|--------|----------|-------------|---------------|
| GET | `/mana/investment-types` | List available investment types | Yes |
| POST | `/mana/investments` | Create new investment | Yes |
| GET | `/mana/investments` | Get wizard's investments | Yes |

### Example API Usage

#### Register a new user
```bash
curl -X POST http://localhost:8080/api/auth/register \
  -H "Content-Type: application/json" \
  -d '{
    "username": "gandalf",
    "email": "gandalf@middleearth.com", 
    "password": "youshallnotpass"
  }'
```

#### Create a wizard
```bash
curl -X POST http://localhost:8080/api/wizards \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer YOUR_JWT_TOKEN" \
  -d '{
    "name": "Gandalf the Grey",
    "realm": "Middle-earth",
    "element": "Light"
  }'
```

## Monitoring & Troubleshooting

### System Status
```bash
# Quick status check
make status

# View logs
make logs
```

### Log Files
```bash
# View real-time logs
tail -f logs/api-gateway.log
tail -f logs/auth-service.log
tail -f logs/wizard-service.log
tail -f logs/mana-service.log

# View all logs
ls -la logs/
```

### Common Issues & Solutions

#### PostgreSQL Connection Failed
```bash
# Check if PostgreSQL is running
make status

# Start PostgreSQL
# macOS:
brew services start postgresql
# Ubuntu:
sudo systemctl start postgresql

# Test connection manually
psql -h localhost -p 5432 -U mysticfunds -d postgres
```

#### Port Already in Use
```bash
# Find process using port
lsof -i :8080
lsof -i :50051

# Kill process if needed
kill -9 <PID>
```

#### Services Won't Start
```bash
# Check build status
make build

# View service logs
cat logs/auth-service.log
cat logs/wizard-service.log
cat logs/mana-service.log
cat logs/api-gateway.log
```

#### Database Migration Errors
```bash
# Check migration status
make migration-status

# Reset databases (destructive!)
make nuke
make create-dbs
make migrate-up
```

#### Frontend Not Loading
```bash
# Check API Gateway status
make status

# Verify API Gateway is serving files
curl http://localhost:8080

# Check browser console for errors
# Open Developer Tools > Console
```

### Performance Monitoring

#### CPU and Memory Usage
```bash
# View detailed process information
make status

# Monitor system resources
top
htop  # if available
```

#### Database Performance
```bash
# Connect to PostgreSQL
psql -h localhost -p 5432 -U postgres

# Check database sizes
\l+

# View active connections
SELECT * FROM pg_stat_activity;
```

## Security Features

- **JWT Authentication** - Secure token-based authentication
- **Password Hashing** - bcrypt password encryption
- **Input Validation** - Comprehensive request validation
- **CORS Support** - Configurable cross-origin requests
- **SQL Injection Prevention** - Parameterized database queries
- **Audit Logging** - Complete transaction history
- **Token Expiration** - Automatic session management

## Contributing

We welcome contributions to MysticFunds! Here's how to get started:

### 1. Fork & Clone
```bash
git clone https://github.com/Alinoureddine1/mysticfunds.git
cd mysticfunds
```

### 2. Create Feature Branch
```bash
git checkout -b feature/amazing-new-feature
```

### 3. Make Changes
- Follow existing code style
- Add tests for new functionality
- Update documentation as needed

### 4. Test Your Changes
```bash
# Run tests
make test

# Test full system
make start
# Test functionality in browser at http://localhost:8080
make stop
```

### 5. Submit Pull Request
- Write clear commit messages
- Include description of changes
- Reference any related issues

## License

This project is licensed under the MIT License with Commercial Use Restrictions - see the [LICENSE](LICENSE) file for details.

## Support & Community

- **Bug Reports** - [GitHub Issues](https://github.com/Alinoureddine1/mysticfunds/issues)
- **Feature Requests** - [GitHub Discussions](https://github.com/Alinoureddine1/mysticfunds/discussions)
- **Documentation** - This README and inline code comments
- **Community** - [Discussions](https://github.com/Alinoureddine1/mysticfunds/discussions)

## Getting Started Guide

### First-Time Users

1. **Clone and Start**
   ```bash
   git clone https://github.com/Alinoureddine1/mysticfunds.git
   cd mysticfunds
   make start
   ```

2. **Open Web Interface**
   - Navigate to `http://localhost:8080`
   - Click "Register" to create an account

3. **Create Your First Wizard**
   - After registration, go to "Wizards" tab
   - Click "Create Wizard"
   - Choose a name, realm, and element

4. **Start Investing**
   - Go to "Investments" tab
   - Select your wizard and investment type
   - Choose amount and create investment

5. **Transfer Mana**
   - Go to "Mana" tab
   - Select source wizard
   - Enter recipient wizard ID and amount

### Demo Data
The system starts with rich pre-configured content:

#### Investment Types
- **Crystal Mining** - Low risk, steady returns
- **Dragon Trading** - Medium risk, good returns  
- **Arcane Ventures** - High risk, high rewards

#### Sample Wizards (27 lore-friendly characters)
- Wizards from all 10 realms with authentic magical names
- Various experience levels and specializations
- Complete with guild affiliations and backstories

#### Jobs (50+ available)
- **4+ jobs per element** across all 10 elemental schools
- Difficulty ranges from Easy to Legendary
- Diverse job types: Research, Combat, Maintenance, Diplomacy, etc.
- Authentic magical locations and requirements

#### Marketplace Items
- **Artifacts** - Magical weapons, armor, and accessories with passive effects
- **Spells** - Learnable magic from different schools and elements  
- **Enhancement Scrolls** - Boost wizard abilities and skills
- **Lore-friendly descriptions** and realm-specific origins

---

## World Lore & Development

MysticFunds includes a rich fantasy world with deep lore and consistent magical systems. For detailed information about the realms, elements, guilds, and development guidelines, see [LORE.md](LORE.md).

### Quick Lore Reference
- **10 Mystical Realms** - From volcanic Pyrrhian Flame to void-touched Nyxthar
- **10 Elemental Schools** - Complete magical system with unique characteristics
- **Rich Guild System** - Organizations spanning all realms with unique purposes
- **Consistent Naming** - Lore-friendly names and descriptions throughout
- **Development Guidelines** - Color schemes, naming conventions, and game mechanics

---

*May your investments be ever magical!*

## About

**MysticFunds** is a passion project created by [Ali Noureddine](https://github.com/Alinoureddine1) - a magical investment platform that combines modern web development with rich fantasy world-building.

**License**: MIT License with Commercial Use Restrictions - see [LICENSE](LICENSE) file for details.

