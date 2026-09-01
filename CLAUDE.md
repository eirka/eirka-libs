# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Repository Overview

eirka-libs is a collection of Go packages shared across different eirka services. These libraries provide common functionality for authentication, database access, Redis caching, configuration, validation, and more.

## Development Commands

### Running Tests

To run all tests in the repository:
```bash
go test ./...
```

To run tests for a specific package:
```bash
go test github.com/eirka/eirka-libs/user
go test github.com/eirka/eirka-libs/validate
go test github.com/eirka/eirka-libs/csrf
```

To run a specific test:
```bash
go test github.com/eirka/eirka-libs/user -run TestIsValidName
```

To enable verbose test output:
```bash
go test -v ./...
```

## Architecture

The repository is organized into domain-specific packages:

1. **user**: User authentication, authorization, and management
   - Handles user creation, password hashing/validation, and permissions
   - Provides JWT token generation and validation
   - Implements secret management with support for key rotation
   - Auth middleware for validating JWT tokens with backward compatibility

2. **db**: Database connection and transaction management
   - Connection pooling and initialization
   - Transaction support
   - Provides test mocks for database testing

3. **redis**: Redis cache management
   - Connection pooling
   - Key management and mutex locks
   - Data storage and retrieval methods

4. **config**: Application configuration
   - Various limits and settings for the application
   - Configuration for external services like Amazon S3
   - Session configuration for JWT secret management
   - Centralized configuration loading from `/etc/pram/pram.conf`

5. **validate**: Request validation
   - Parameter validation for API requests
   - Utility functions for common validation tasks

6. **csrf**: CSRF protection
   - Token generation and validation
   - Middleware for protecting routes

7. **audit**: Auditing system
   - Records user actions
   - Provides audit trail functionality

## Testing

The codebase uses the standard Go testing package with additional libraries:

1. **stretchr/testify**: For assertions and better test output
2. **DATA-DOG/go-sqlmock**: For mocking database connections
3. **redigomock**: For mocking Redis connections
4. **tempredis**: For creating temporary Redis servers during tests

Tests follow a consistent pattern:
- Each test function is prefixed with `Test`
- Most tests use mock databases/Redis to avoid external dependencies
- Tests are organized by package, matching the production code structure

## Common Patterns

1. **Error Handling**: The codebase centralizes errors in the `errors` package and generally follows Go's idiomatic error handling.

2. **Dependency Injection**: Database and Redis connections are initialized once and then accessed through getter functions.

3. **Middleware**: Packages like `csrf` and `validate` provide middleware functions designed to work with the Gin web framework.

4. **Mocking**: The codebase extensively uses mocking for tests, particularly for database and Redis operations.

5. **Secret Management**: JWT secrets are stored in the central configuration file and support rotation through a primary/secondary secret mechanism. The SecretManager ensures thread-safe access to secrets.

## Test prerequisites

- Only `redis/` needs a real server: `tempredis.Start` in redis/redis_test.go, mutex_test.go, keys_test.go execs `redis-server -` and panics (no skip) if it is not on PATH; `make preflight` checks this.
- `go test -count=1 ./...` takes ~18s wall (redis ~16s, user ~18s); not hung.
- Runs rewrite `redis/dump.rdb` (gitignored `*.rdb`).
- Gate before PR: `make check` (build + gofmt/go vet + `go test -count=1 ./...`).

## Fleet conventions

Tags/releases are cut by the maintainer only, by convention (nothing in-repo enforces it) - never tag here; consumers eirka-admin, eirka-get, eirka-index, eirka-post pin an explicit version with `go get github.com/eirka/eirka-libs@vX.Y.Z` because proxy.golang.org serves a cached version list and lags fresh tags.
