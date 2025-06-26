# Testing and CI/CD Documentation

This document explains the testing strategy and CI/CD pipeline for the Muse backend GraphQL API.

## Testing Strategy

### 1. Test Types

#### Unit Tests

- **Location**: `*_test.go` files next to source code
- **Purpose**: Test individual functions and methods in isolation
- **Coverage**: Repository methods, utility functions, converters
- **Command**: `go test -short ./...`

#### Integration Tests

- **Location**: `integration_test.go`
- **Purpose**: Test GraphQL resolvers end-to-end with real database
- **Coverage**: Full GraphQL operations, authentication, database interactions
- **Command**: `go test -v ./...`

#### Performance Tests

- **Location**: `*_bench_test.go` files
- **Purpose**: Benchmark critical operations
- **Coverage**: Database queries, GraphQL operations
- **Command**: `go test -bench=. ./...`

### 2. Running Tests Locally

#### Prerequisites

```bash
# Install dependencies
go mod download

# Set up test database (optional for integration tests)
createdb muse_test
```

#### Running Different Test Types

```bash
# Unit tests only (fast, no database required)
go test -short ./...

# All tests including integration tests (requires database)
go test ./...

# Tests with coverage
go test -coverprofile=coverage.out ./...
go tool cover -html=coverage.out

# Benchmarks
go test -bench=. -benchmem ./...

# Specific integration test
go test -run TestGraphQLIntrospection -v
```

#### Environment Variables for Testing

```bash
# Database (for integration tests)
export DB_HOST=localhost
export DB_PORT=5432
export DB_USER=postgres
export DB_PASSWORD=password
export DB_NAME=muse_test
export DB_SSL_MODE=disable

# Redis (optional)
export REDIS_HOST=localhost
export REDIS_PORT=6379

# Spotify (for API integration tests)
export SPOTIFY_CLIENT_ID=your_client_id
export SPOTIFY_CLIENT_SECRET=your_client_secret

# JWT Secret
export JWT_SECRET=test-jwt-secret
```

### 3. Test Organization

```
backend/
├── integration_test.go           # GraphQL integration tests
├── internal/
│   ├── repository/
│   │   └── postgres/
│   │       ├── user_repository_test.go      # Unit tests for user repo
│   │       ├── album_repository_test.go     # Unit tests for album repo
│   │       └── ...
│   └── models/
│       └── model_test.go         # Unit tests for models
└── graph/
    └── resolver_test.go          # Unit tests for resolvers (if any)
```

## CI/CD Pipeline

### 1. Pipeline Overview

The GitHub Actions pipeline consists of 8 jobs that run in sequence and parallel:

```
lint ──┬── build ──┬── integration-tests ──┬── docker ── deploy-staging
       ├── unit-tests                      │
       ├── security                        └── performance
       └── ...
```

### 2. Pipeline Jobs

#### Job 1: Lint and Format Check

- **Triggers**: All pushes and PRs
- **Actions**:
  - Go module verification
  - `golangci-lint` linting
  - `gofmt` formatting check
  - `go vet` static analysis

#### Job 2: Build Application

- **Triggers**: After lint passes
- **Actions**:
  - Compile Go application
  - Upload build artifact
  - Verify binary works

#### Job 3: Unit Tests

- **Triggers**: After lint passes
- **Actions**:
  - Run unit tests with race detection
  - Generate coverage report
  - Upload coverage to Codecov

#### Job 4: Integration Tests

- **Triggers**: After lint and build pass
- **Services**: PostgreSQL 15, Redis 7
- **Actions**:
  - Run database migrations
  - Execute integration tests
  - Generate integration coverage

#### Job 5: Security Scan

- **Triggers**: After lint passes
- **Actions**:
  - Run `gosec` security scanner
  - Upload SARIF results to GitHub Security

#### Job 6: Docker Build & Push

- **Triggers**: Only on `main` branch pushes
- **Actions**:
  - Build Docker image with caching
  - Push to Docker Hub
  - Tag with branch, SHA, and `latest`

#### Job 7: Deploy to Staging

- **Triggers**: After Docker build on `main` branch
- **Environment**: `staging` (requires approval)
- **Actions**:
  - Deploy to staging environment
  - Run smoke tests

#### Job 8: Performance Tests

- **Triggers**: PRs and `main` branch
- **Services**: PostgreSQL 15
- **Actions**:
  - Run performance benchmarks
  - Upload benchmark results

### 3. Required Secrets

Set these secrets in your GitHub repository settings:

```bash
# Docker Hub (for image publishing)
DOCKER_USERNAME=your_dockerhub_username
DOCKER_PASSWORD=your_dockerhub_password

# Spotify API (for integration tests)
SPOTIFY_CLIENT_ID=your_spotify_client_id
SPOTIFY_CLIENT_SECRET=your_spotify_client_secret

# Deployment (optional)
STAGING_DEPLOY_URL=your_staging_deployment_webhook
PROD_DEPLOY_URL=your_production_deployment_webhook
```

### 4. Branch Protection

Recommended branch protection rules for `main`:

```yaml
# .github/branch-protection.yml
protection_rules:
  main:
    required_status_checks:
      strict: true
      contexts:
        - "lint"
        - "build"
        - "unit-tests"
        - "integration-tests"
        - "security"
    enforce_admins: true
    required_pull_request_reviews:
      required_approving_review_count: 1
      dismiss_stale_reviews: true
    restrictions: null
```

## Development Workflow

### 1. Pre-commit Hooks (Recommended)

Install pre-commit hooks to catch issues early:

```bash
# Install pre-commit
pip install pre-commit

# Install hooks (create .pre-commit-config.yaml first)
pre-commit install
```

Example `.pre-commit-config.yaml`:

```yaml
repos:
  - repo: local
    hooks:
      - id: go-fmt
        name: go fmt
        entry: gofmt
        language: system
        args: [-w]
        files: \.go$

      - id: go-vet
        name: go vet
        entry: go vet
        language: system
        files: \.go$
        pass_filenames: false

      - id: go-test
        name: go test
        entry: go test -short
        language: system
        files: \.go$
        pass_filenames: false
```

### 2. Feature Development Flow

1. **Create Feature Branch**

   ```bash
   git checkout -b feature/new-feature
   ```

2. **Develop with Tests**

   ```bash
   # Write code and tests
   # Run tests locally
   go test -short ./...
   ```

3. **Pre-commit Checks**

   ```bash
   # Format code
   gofmt -w .

   # Run linter
   golangci-lint run

   # Run tests
   go test ./...
   ```

4. **Create Pull Request**

   - CI pipeline runs automatically
   - All checks must pass before merge
   - Code review required

5. **Merge to Main**
   - Triggers full pipeline including deployment
   - Docker image built and pushed
   - Staging deployment (if configured)

### 3. Hotfix Flow

For urgent production fixes:

1. **Create Hotfix Branch from Main**

   ```bash
   git checkout main
   git checkout -b hotfix/urgent-fix
   ```

2. **Fast-track Review**

   - Minimal but focused changes
   - Priority code review
   - All CI checks still required

3. **Deploy**
   - Merge triggers automatic deployment
   - Monitor staging before production

## Monitoring and Observability

### 1. Test Results

- **GitHub Actions**: View test results in Actions tab
- **Codecov**: Coverage reports and trends
- **Security**: SARIF results in Security tab

### 2. Performance Monitoring

- **Benchmark Results**: Stored as artifacts
- **Performance Trends**: Compare across commits
- **Database Query Performance**: Monitor slow queries

### 3. Deployment Health

- **Health Endpoints**: `/health` endpoint for monitoring
- **Smoke Tests**: Basic functionality verification
- **Error Tracking**: Monitor application errors

## Troubleshooting

### Common Issues

#### Tests Fail Locally But Pass in CI

- Check environment variables
- Ensure database is clean between tests
- Verify Go version compatibility

#### Docker Build Fails

- Check Dockerfile syntax
- Verify build context includes necessary files
- Check for platform-specific dependencies

#### Integration Tests Timeout

- Database connection issues
- Increase timeout values
- Check service health in CI

#### Deployment Fails

- Verify secrets are set correctly
- Check deployment target accessibility
- Review deployment logs

### Getting Help

1. **Check CI Logs**: Detailed error messages in GitHub Actions
2. **Local Reproduction**: Run the exact commands from CI locally
3. **Database Issues**: Check connection strings and migrations
4. **Performance Issues**: Run benchmarks locally to identify bottlenecks

## Best Practices

### 1. Test Writing

- Write tests before or alongside code (TDD)
- Use table-driven tests for multiple scenarios
- Mock external dependencies in unit tests
- Use real database for integration tests

### 2. CI/CD

- Keep pipeline fast (< 10 minutes total)
- Fail fast on critical issues
- Cache dependencies for speed
- Use parallel jobs when possible

### 3. Security

- Scan dependencies regularly
- Use least privilege for secrets
- Review security reports promptly
- Keep dependencies updated

### 4. Performance

- Set performance budgets
- Monitor key metrics
- Optimize slow tests
- Profile critical paths

This testing and CI/CD strategy ensures high code quality, fast feedback, and reliable deployments for the Muse backend API.
