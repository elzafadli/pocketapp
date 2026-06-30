# 📚 Pakuningratan - Clean Architecture Service Template

## 📖 About

**Pakuningratan** is a production-ready Go service skeleton that follows Clean Architecture principles. It serves as a template and guidance for creating new microservices with best practices, standardized structure, and common features pre-implemented.

This project provides a solid foundation for building scalable, maintainable, and testable Go services with:

- Clean Architecture layered structure
- Dependency Injection using gontainer
- Graceful shutdown handling
- Comprehensive logging
- Database and cache support
- RESTful API with Fiber framework
- Configuration management
- Docker support

---

## ✨ Features

### Core Features

- **🏗️ Clean Architecture**
  - Well-organized layered structure (Domain, Application, Adapter)
  - Separation of concerns
  - Dependency inversion principle
  - Easy to test and maintain

- **🔄 Dependency Injection**
  - Container-based DI using `gontainer`
  - Automatic service registration
  - Lifecycle management (Startup/Shutdown)

- **🛡️ Graceful Shutdown**
  - Signal handling (SIGTERM, SIGINT)
  - Proper resource cleanup
  - Server shutdown with timeout
  - Error handling during shutdown

- **📝 Logging**
  - Structured logging with `golog`
  - File-based logging with rotation
  - Configurable log levels
  - Request tracing support

- **⚙️ Configuration Management**
  - YAML-based configuration
  - Environment variable support
  - Automatic env variable override
  - Type-safe configuration structs

- **🗄️ Database Support**
  - GORM integration
  - PostgreSQL driver
  - Connection pooling
  - Auto-migration support

- **💾 Cache Support**
  - Redis client integration
  - Support for single/sentinel/cluster modes
  - Connection pooling

- **🌐 REST API**
  - Fiber v2 framework
  - Custom middleware (logging, recovery, CORS)
  - Request ID tracking
  - Error handling middleware
  - Authentication middleware (BasicAuth, API Key)

- **✅ Validation**
  - Struct validation using `validator/v10`
  - Custom error messages
  - Translation support


- **🔐 Authentication**
  - BasicAuth protection middleware
  - API Key authentication middleware
  - Configurable authentication credentials
  - Domain-level authentication errors

- **🐳 Docker Support**
  - Multi-stage Dockerfile
  - Docker Compose setup
  - Production-ready image

- **🛠️ Development Tools**
  - Air for hot-reload
  - Makefile with common tasks
  - Test coverage enforcement
  - Code generation support

---

## 🏗️ Architecture

### Clean Architecture Layers

```
┌─────────────────────────────────────────┐
│          Presentation Layer             │
│         (REST API Handlers)             │
└─────────────────┬───────────────────────┘
                  │
┌─────────────────▼───────────────────────┐
│         Application Layer               │
│      (Use Cases / Services)             │
└─────────────────┬───────────────────────┘
                  │
┌─────────────────▼───────────────────────┐
│            Domain Layer                 │
│      (Entities / Business Logic)        │
└─────────────────┬───────────────────────┘
                  │
┌─────────────────▼───────────────────────┐
│         Infrastructure Layer            │
│   (Database, Cache, External APIs)      │
└─────────────────────────────────────────┘
```

### Project Structure

```
pakuningratan/
├── cmd/                   # Command-line entry points
│   ├── cmd.go             # Cobra command setup
│   └── service.go         # Service command implementation
├── config/                # Configuration management
│   ├── model.go           # Configuration structs
│   └── config.yaml        # Default configuration file
├── internal/
│   ├── adapter/           # Infrastructure adapters
│   │   ├── repository/    # Data access layer
│   │   │   ├── database/  # Database repositories
│   │   │   └── cache/     # Cache repositories
│   │   ├── rest/          # REST API adapter (Fiber)
│   │   └── toronto/       # External service adapters
│   ├── application/       # Application layer (use cases)
│   │   ├── api/           # API handlers
│   │   ├── service/       # Business logic services
│   │   └── api.go         # API route registration
│   ├── bootstrap/         # Application bootstrap
│   │   ├── bootstrap.go   # Main bootstrap logic
│   │   ├── adapter.go     # Adapter registration
│   │   └── application.go # Application registration
│   ├── domain/            # Domain layer (entities)
│   │   ├── auth/          # Authentication domain entities
│   │   └── shared/        # Shared domain entities
│   └── pkg/               # Shared packages
│       ├── validator/     # Validation utilities
│       ├── formatter/     # Response formatters
│       └── custommiddleware/ # Custom middleware (auth, log, error, etc.)
├── bin/                   # Build output directory
├── logs/                  # Log files directory
├── main.go                # Application entry point
├── Dockerfile             # Docker build configuration
├── docker-compose.yml     # Docker Compose configuration
├── Makefile               # Build automation
├── .air.toml              # Air hot-reload configuration
├── go.mod                 # Go module dependencies
└── README.md              # This file
```

### Dependency Flow

```
main.go
  └── cmd.Execute()
      └── cmd.service
          └── bootstrap.Run()
              ├── RegisterDatabase()
              ├── RegisterCache()
              ├── RegisterRest()
              ├── RegisterRepository()
              ├── RegisterService()
              └── RegisterApi()
```

---

## 🚀 Getting Started

### Prerequisites

- **Go 1.25.5** or later ([Download](https://golang.org/dl/))
- **PostgreSQL** 12+ (for database)
- **Redis** 6+ (for caching)
- **Docker** and **Docker Compose** (optional, for containerized setup)
- **Air** for hot-reload during development ([Installation](https://github.com/cosmtrek/air))

### Installation

#### 1. Clone the Repository

```bash
git clone <repository-url>
cd pakuningratan
```

#### 2. Install Dependencies

```bash
go mod download
go mod tidy
```

#### 3. Setup Environment Variables

Create a `.env` file in the root directory:

```bash
# Copy from sample if available, or create manually
cp .env.example .env
```

Configure your environment variables. The application supports both YAML configuration and environment variable overrides. Environment variables use underscore notation and are automatically mapped to configuration keys.

Example `.env` file:

```env
# Database
DATABASE_USER=postgres
DATABASE_PASSWORD=postgres

# Redis
REDIS_PASSWORD=secret
```

#### 4. Configure Database

Ensure PostgreSQL is running and create the database:

```sql
CREATE DATABASE pakuningratan;
```

The application uses GORM for database operations. Migrations are handled automatically by GORM's AutoMigrate feature.

#### 5. Configure Redis

Ensure Redis is running on the configured host and port.

---

## 💻 Development

### Running the Application

#### Development Mode (Hot Reload)

Use Air for automatic reloading during development:

```bash
make watch
# or
air -c .air.toml
```

This will:

- Watch for file changes
- Automatically rebuild and restart the application
- Exclude test files and build artifacts

#### Production Mode

Build and run the binary:

```bash
# Build
make build

# Run
./bin/pakuningratan service
# or
make run
```

#### Using Docker Compose

Run the entire stack (application + dependencies):

```bash
# Start all services
docker-compose up -d

# View logs
docker-compose logs -f pakuningratan

# Stop services
docker-compose stop

# Stop and remove containers
docker-compose down
```

### Available Make Commands

```bash
make help          # Show all available commands
make watch         # Run in development mode with hot-reload
make build         # Build binary for local architecture
make run           # Run the application
make test          # Run tests with coverage check
make coverage      # Show detailed test coverage
make coverage-web  # Generate HTML coverage report
make bench         # Run benchmark tests
make check         # Run SonarQube code quality check
```

### Configuration

The application uses a two-tier configuration system:

1. **YAML Configuration** (`config.yaml`) - Base configuration
2. **Environment Variables** (`.env`) - Override YAML values

#### Configuration File Structure

```yaml
app: pakuningratan
app_version: v0.0.1
env: development

http:
  port: 8000
  write_timeout: 30 # seconds
  read_timeout: 30 # seconds

log:
  file_location: "logs"
  file_tdr_location: "logs"
  file_max_size: 20 # megabytes
  file_max_backup: 10 # number of backups
  file_max_age: 30 # days
  stdout: true

database:
  host: 127.0.0.1
  port: 5432
  name: pakuningratan
  user: postgres
  password: postgres
  ssl_mode: disable
  max_idle_conn: 10
  conn_max_lifetime: 1h
  max_open_conn: 100

redis:
  mode: single # single/sentinel/cluster
  address: 127.0.0.1
  port: 6379
  password: ""

basic_auths: user:password,user2:password2 # Comma-separated username:password pairs
internal_api_key: changemeinprod # API key for internal service authentication
```

#### Environment Variable Override

Environment variables automatically override YAML values. Use underscore notation:

- `HTTP_PORT` → `http.port`
- `DATABASE_HOST` → `database.host`
- `LOG_STDOUT` → `log.stdout`
- `BASIC_AUTHS` → `basic_auths`
- `INTERNAL_API_KEY` → `internal_api_key`

#### Authentication Configuration

- **BasicAuth**: Configure multiple username:password pairs separated by commas
  - Format: `"username1:password1,username2:password2"`
  - Example: `BASIC_AUTHS="admin:secret,user:pass123"`
- **API Key**: Set a secret key for API key authentication
  - Used with `X-Api-Key` header
  - Example: `INTERNAL_API_KEY="your-secret-key-here"`

---

## 🧪 Testing

### Running Tests

```bash
# Run all tests with coverage check
make test

# View detailed coverage
make coverage

# Generate HTML coverage report
make coverage-web

# Run benchmarks
make bench
```

### Test Coverage

The project enforces a minimum test coverage threshold (default: 15%). Tests must meet this threshold to pass CI/CD pipelines.

### Writing Tests

- Place test files alongside source files with `_test.go` suffix
- Use `testify` for assertions
- Mock dependencies using `go.uber.org/mock`
- Generate mocks: `go generate ./...`

---

## 📡 API Endpoints

### Health Check Endpoints

#### Ping

```http
GET /ping
```

Simple health check endpoint.

**Response:**

```json
{
  "message": "pong"
}
```

#### Ready

```http
GET /ready
```

Readiness check endpoint. Returns service status.

**Response:**

```json
{
  "database": "ready",
  "cache": "ready"
}
```

### Template Endpoints (Example)

All template endpoints require BasicAuth authentication. Include credentials in the Authorization header.

#### Create Template

```http
POST /v1/templates
Authorization: Basic <base64(username:password)>
Content-Type: application/json

{
  "name": "Example Template",
  "description": "Template description"
}
```

#### List Templates

```http
GET /v1/templates
Authorization: Basic <base64(username:password)>
```

#### Get Template by ID

```http
GET /v1/templates/:id
Authorization: Basic <base64(username:password)>
```

#### Update Template

```http
PUT /v1/templates/:id
Authorization: Basic <base64(username:password)>
Content-Type: application/json

{
  "name": "Updated Template",
  "description": "Updated description"
}
```

**Note**: Replace `<base64(username:password)>` with the base64-encoded credentials. For example, if username is `user` and password is `password`, the value would be `dXNlcjpwYXNzd29yZA==` (base64 of "user:password").

### Using API Key Authentication

Alternatively, you can use API Key authentication by including the `X-Api-Key` header:

```http
GET /v1/templates
X-Api-Key: your-api-key-here
```

The API key must match the `internal_api_key` configured in your configuration file or environment variable.

### API Response Format

#### Success Response

```json
{
  "status": "00",
  "message": "success",
  "data": { ... }
}
```

#### Error Response

```json
{
  "status": "PAKU05",
  "message": "sample error message",
  "trace_id": "uuid-here"
}
```

#### Authentication Error Response

```json
{
  "status": "PAKU04",
  "message": "auth: invalid basic auth",
  "trace_id": "uuid-here"
}
```

Or for API Key authentication:

```json
{
  "status": "PAKU04",
  "message": "auth: invalid api key",
  "trace_id": "uuid-here"
}
```

---

## 🐳 Docker

### Building Docker Image

```bash
docker build -t pakuningratan:latest .
```

### Docker Compose

The `docker-compose.yml` file includes:

- Application service
- Network configuration
- Volume mounts for logs
- Environment variable support

### Production Deployment

The Dockerfile uses a multi-stage build:

1. **Builder stage**: Compiles Go binary
2. **Runtime stage**: Minimal Alpine image with binary

Key features:

- CGO disabled for static binary
- Timezone configuration (Asia/Jakarta)
- Minimal image size
- Security best practices

---

## 🔧 Configuration Details

### HTTP Server Configuration

- **Port**: Server listening port (default: 8000)
- **Read Timeout**: Maximum duration for reading request (seconds)
- **Write Timeout**: Maximum duration for writing response (seconds)

### Database Configuration

- **Connection Pooling**: Configured via `max_idle_conn`, `max_open_conn`
- **Connection Lifetime**: Maximum time a connection can be reused
- **SSL Mode**: PostgreSQL SSL connection mode

### Redis Configuration

- **Mode**: `single`, `sentinel`, or `cluster`
- **Address**: Redis server address
- **Port**: Redis server port
- **Password**: Redis authentication password (if required)

### Logging Configuration

- **File Location**: Directory for log files
- **Max Size**: Maximum log file size before rotation (MB)
- **Max Backup**: Number of backup files to keep
- **Max Age**: Maximum age of log files (days)
- **Stdout**: Enable/disable console output

### Authentication Configuration

- **BasicAuths**: Comma-separated list of username:password pairs
  - Format: `"user1:pass1,user2:pass2"`
  - Used for BasicAuth middleware protection
  - Example: `BASIC_AUTHS="admin:secret123,user:password"`
- **InternalApiKey**: Secret key for API key authentication
  - Used with `X-Api-Key` header
  - Should be changed in production
  - Example: `INTERNAL_API_KEY="your-production-key"`

---

## 🏛️ Architecture Principles

### Clean Architecture Benefits

1. **Independence**: Business logic independent of frameworks
2. **Testability**: Easy to test without external dependencies
3. **Independence of UI**: Business logic works with any interface
4. **Independence of Database**: Business logic independent of data storage
5. **Independence of External Services**: Business logic doesn't know about external services

### Dependency Rule

Dependencies point inward:

- Outer layers depend on inner layers
- Inner layers don't know about outer layers
- Domain layer has no dependencies

### Project Philosophies

1. **Keep things flat**: Avoid unnecessary subdirectories unless truly needed
2. **Clarity over cleverness**: If code is hard to understand, it needs improvement
3. **Simplicity first**: Don't over-engineer; complexity will come naturally
4. **Readability matters**: Code should be self-explanatory
5. **Fast understandability over fast performance**: Go is fast enough; prioritize clarity
6. **Copyability**: Good code is easy to copy and adapt

---

## 📊 Monitoring & Observability

### Logging

The application uses structured logging with:

- Request ID tracking (trace ID)
- Log levels (Info, Error, Panic)
- File rotation
- Configurable output (file/stdout)

### Health Checks

- `/ping`: Basic health check
- `/ready`: Readiness probe for Kubernetes

---

## 🔒 Security Considerations

- **Authentication & Authorization**
  - BasicAuth middleware for endpoint protection
  - API Key authentication support
  - Configurable authentication credentials
  - Domain-level authentication error handling
- **Input Security**
  - Input validation using struct validators
  - SQL injection prevention via GORM
  - Error message sanitization
- **Infrastructure Security**
  - CORS configuration
  - Request ID for tracing
  - Graceful shutdown prevents data loss
  - Secure configuration management

---

## 🚢 Deployment

### Building for Production

```bash
# Build binary
make build

# Or build with specific GOOS/GOARCH
CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o ./bin/pakuningratan
```

### Environment Setup

1. Set production environment variables
2. Configure database connection
3. Configure Redis connection
4. Set appropriate log levels

### Kubernetes Deployment

The service is Kubernetes-ready with:

- Health check endpoints (`/ping`, `/ready`)
- Graceful shutdown support
- ConfigMap/Secret support via environment variables
- Resource limits configuration

---

## 🤝 Contributing

1. Follow the project's architecture principles
2. Write tests for new features
3. Ensure test coverage meets minimum threshold
4. Update documentation for API changes
5. Follow Go code style guidelines
6. Keep functions small and focused
7. Add comments for complex logic

---

## 📞 Support

For issues, questions, or contributions:

- **Email**: tommy@runsystem.id
- **Issues**: [GitLab/GitHub Issues URL]

---

## 📚 Additional Resources

- [Clean Architecture by Robert C. Martin](https://blog.cleancoder.com/uncle-bob/2012/08/13/the-clean-architecture.html)
- [Go Best Practices](https://golang.org/doc/effective_go)
- [Fiber Documentation](https://docs.gofiber.io/)
- [GORM Documentation](https://gorm.io/docs/)
