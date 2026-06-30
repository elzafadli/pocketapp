# seedapp

![seedapp](assets/seedapp.png)

seedapp is a service for managing tenant database migrations for application Accounting Plus.

## 🚀 Features

- Clean Architecture implementation
- RESTful API using Fiber framework
- PostgreSQL database integration with migrations
- Dependency injection using Gontainer
- Structured logging
- Environment configuration management
- Docker support
- Hot reload development support

## 📋 Prerequisites

- Go 1.21.3 or higher
- Docker and Docker Compose
- PostgreSQL
- Make (for using Makefile commands)

## 🛠 Tech Stack

- **Framework:** [Fiber v2](https://github.com/gofiber/fiber)
- **Database:** PostgreSQL with [pgx](https://github.com/jackc/pgx) and [sqlx](https://github.com/jmoiron/sqlx)
- **SQL Query Builder:** [Squirrel](https://github.com/Masterminds/squirrel)
- **Migration:** [sql-migrate](https://github.com/rubenv/sql-migrate)
- **Configuration:** [Viper](https://github.com/spf13/viper)
- **CLI:** [Cobra](https://github.com/spf13/cobra)
- **Validation:** [validator/v10](https://github.com/go-playground/validator)
- **Logging:** [golog](https://github.com/runsystemid/golog)
- **Dependency Injection:** [gontainer](https://github.com/runsystemid/gontainer)

## 📁 Project Structure

```
seedapp/
├── bin/             # Compiled binaries
├── cmd/             # Command line interface entry points
├── config/          # Configuration files
├── internal/        # Private application code
│   ├── adapter/     # External adapters (HTTP, DB, etc.)
│   ├── application/ # Application business logic
│   ├── bootstrap/   # Application bootstrapping
│   ├── domain/      # Domain models and interfaces
│   └── pkg/         # Internal shared packages
├── logs/            # Log files
└── scripts/         # Utility scripts
```

## 🚀 Getting Started

1. Clone the repository

```bash
git clone git@gitlab.runsystemdev.com:runsystem_dev/acc/seedapp.git
cd seedapp
```

2. Copy the environment file and configure it

```bash
cp .env.sample .env
# Edit .env with your configurations
```

3. Prepare the pre-requisites

```bash
make install
make build
```

4. Run the application

```bash
# Using make commands
make run

# Or using Go directly
go run main.go service
```

## 🛠 Development

### Available Make Commands

```bash
make build      # Build the application
make run        # Run the application
make test       # Run tests
make coverage   # Run tests and check coverage
```

### Hot Reload Development

The project includes `.air.toml` configuration for hot reload during development. To use it:

```bash
go install github.com/cosmtrek/air@latest
air
```

## 🐳 Docker

Build and run using Docker:

```bash
# Build the image
docker build -t seedapp .

# Run the container
docker run -p 8080:8080 seedapp
```

Or using Docker Compose:

```bash
docker-compose up -d
```

## FAQ

**Q:** Where is the api documentation?

**A:** Refer to ht [Accounting Plus Blueprint](http://gitlab.runsystemdev.com/runsystem_dev/documentation/runs-api). Clone the repo and run it with [Bruno App](https://www.usebruno.com/).

## Commit Guidelines

### Commit Message Format

Please use the following format for commit message:
`[PLANE-ID]: Commit message`

**Example:**
`[HNP-1]: Add additonal unit test for logger`

### Merge Request Format

**Title:** `[PLANE-ID]: Title`

**Description:**

```
Changelogs:
- Restructure logger package
- Add additional unit test for logger
```

## Philosophies and Rules

- **Keep Things Flat:** Avoid creating subdirectories/subfiles unless necessary. Most issues can be resolved with good naming.
- **Clarity Over Complexity:** If something is hard to understand, it's a valid concern. Assume the code can be improved.
- **Simplicity is Key:** Avoid unnecessary complexity. Challenges will arise naturally.
- **Ease of Understanding:** Code should be easy to read and understand.
- **Prioritize Understandability Over Performance:** Golang is fast, and cloud resources are affordable. Focus on simplicity without sacrificing common sense.
- **Copasability:** If it's good, it's copasable. If it's really copasable, then it's perfect.
