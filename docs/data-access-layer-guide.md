# Data Access Layer Guide

This guide explains how to add a new storage backend to the AI API Gateway's data access layer.

## Overview

The data access layer follows a repository pattern with clear separation between:
- **Domain entities**: Pure data structures in `internal/domain/entity/`
- **Repository interfaces**: Contracts in `internal/domain/port/repository.go`
- **Repository implementations**: Concrete storage implementations in `internal/infrastructure/repository/`

## Architecture

```
Service Layer
    ↓ (uses)
Repository Interface (port)
    ↓ (implemented by)
Repository Implementation (infrastructure)
    ↓ (uses)
Database/Storage
```

## Adding a New Storage Backend

### Step 1: Define the Repository Interface

Add the interface in `internal/domain/port/repository.go`:

```go
// Example: UserRepository interface
type UserRepository interface {
    Create(user *entity.User) error
    GetByID(id string) (*entity.User, error)
    GetByEmail(email string) (*entity.User, error)
    Update(user *entity.User) error
    Delete(id string) error
    List(page, pageSize int) ([]*entity.User, int, error)
}
```

**Key principles:**
- Interface should be storage-agnostic
- Use domain entity types
- Return errors for all failure cases
- Support pagination for list operations

### Step 2: Implement the Repository

Create implementation in `internal/infrastructure/repository/`:

```go
// user_repository.go
package repository

import (
    "github.com/ai-api-gateway/auth-service/internal/domain/entity"
    "github.com/ai-api-gateway/auth-service/internal/domain/port"
    "gorm.io/gorm"
)

type GormUserRepository struct {
    db *gorm.DB
}

func NewGormUserRepository(db *gorm.DB) port.UserRepository {
    return &GormUserRepository{db: db}
}

func (r *GormUserRepository) Create(user *entity.User) error {
    return r.db.Create(user).Error
}

func (r *GormUserRepository) GetByID(id string) (*entity.User, error) {
    var user entity.User
    err := r.db.Where("id = ?", id).First(&user).Error
    if err != nil {
        return nil, err
    }
    return &user, nil
}

// ... implement other methods
```

**Implementation guidelines:**
- Prefix with storage type (e.g., `GormUserRepository`, `PostgresUsageRecordRepository`)
- Handle database-specific errors appropriately
- Use transactions where needed (e.g., for multi-record operations)

### Step 3: Wire the Implementation

In your service initialization code:

```go
// Example in auth-service/cmd/server/main.go
func main() {
    // Initialize database connection
    db, err := gorm.Open(sqlite.Open("auth.db"), &gorm.Config{})
    if err != nil {
        log.Fatal(err)
    }

    // Auto-migrate schema
    db.AutoMigrate(&entity.User{}, &entity.APIKey{})

    // Create repository implementation
    userRepo := repository.NewGormUserRepository(db)
    apiKeyRepo := repository.NewGormAPIKeyRepository(db)

    // Create service with repositories
    authService := application.NewAuthService(userRepo, apiKeyRepo, nil)
}
```

### Step 4: Add Database Migration

For GORM-based implementations, use AutoMigrate in initialization:

```go
db.AutoMigrate(&entity.YourEntity{})
```

For manual migrations, create migration files in `internal/infrastructure/migration/`:

```go
// migration.go
func RunMigrations(db *gorm.DB) error {
    return db.AutoMigrate(
        &entity.User{},
        &entity.APIKey{},
        &entity.UsageRecord{},
    )
}
```

## Testing Strategies

### Unit Tests

Test repository logic with an in-memory database:

```go
func TestGormUserRepository_Create(t *testing.T) {
    // Setup in-memory SQLite
    db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
    if err != nil {
        t.Fatal(err)
    }
    db.AutoMigrate(&entity.User{})

    repo := NewGormUserRepository(db)

    // Test
    user := &entity.User{ID: "test-1", Name: "Test User"}
    err = repo.Create(user)
    assert.NoError(t, err)

    // Verify
    fetched, err := repo.GetByID("test-1")
    assert.NoError(t, err)
    assert.Equal(t, "Test User", fetched.Name)
}
```

### Integration Tests

Test with a real database using Docker or testcontainers:

```go
func TestUserRepository_Integration(t *testing.T) {
    if testing.Short() {
        t.Skip("skipping integration test")
    }

    // Setup test database container
    // Run repository tests
    // Cleanup
}
```

## Adding a New Storage Backend (e.g., PostgreSQL)

### 1. Create PostgreSQL Repository Implementation

```go
// postgres_user_repository.go
package repository

import (
    "database/sql"
    "github.com/ai-api-gateway/auth-service/internal/domain/entity"
)

type PostgresUserRepository struct {
    db *sql.DB
}

func NewPostgresUserRepository(db *sql.DB) port.UserRepository {
    return &PostgresUserRepository{db: db}
}

func (r *PostgresUserRepository) Create(user *entity.User) error {
    query := `INSERT INTO users (id, name, email, role, status) VALUES ($1, $2, $3, $4, $5)`
    _, err := r.db.Exec(query, user.ID, user.Name, user.Email, user.Role, user.Status)
    return err
}

// ... implement other methods with SQL queries
}
```

### 2. Add Configuration

```go
// config.go
type Config struct {
    StorageBackend string `yaml:"storage_backend"` // "sqlite" | "postgres"
    SQLitePath     string `yaml:"sqlite_path"`
    PostgresURL    string `yaml:"postgres_url"`
}
```

### 3. Factory Pattern for Repository Selection

```go
func NewUserRepository(cfg *Config) (port.UserRepository, error) {
    switch cfg.StorageBackend {
    case "postgres":
        db, err := sql.Open("postgres", cfg.PostgresURL)
        if err != nil {
            return nil, err
        }
        return repository.NewPostgresUserRepository(db), nil
    case "sqlite":
        db, err := gorm.Open(sqlite.Open(cfg.SQLitePath), &gorm.Config{})
        if err != nil {
            return nil, err
        }
        return repository.NewGormUserRepository(db), nil
    default:
        return repository.NewGormUserRepository(db), nil
    }
}
```

## Best Practices

1. **Interface First**: Always define the interface before implementation
2. **Storage Agnostic**: Keep domain entities free of storage-specific details
3. **Error Handling**: Return descriptive errors, don't panic in repositories
4. **Transactions**: Use transactions for multi-record operations
5. **Pagination**: Support pagination for list operations to avoid loading all data
6. **Testing**: Test both unit (in-memory) and integration (real database) scenarios
7. **Migrations**: Keep migrations in version control and run them on startup
8. **Connection Pooling**: Configure appropriate connection pool settings for production

## Example: Adding PostgreSQL Support to Billing Service

See the token-tracker-implementation change for a complete example of adding PostgreSQL support:
- Task 4.1: Implement PostgreSQL repository for UsageRecord
- Task 4.2: Add database migration scripts
- Task 4.3: Add configuration flag to select backend

## Current Storage Backends

| Service | Current Backend | Alternative |
|---------|----------------|-------------|
| auth-service | SQLite (GORM) | PostgreSQL |
| billing-service | SQLite (GORM) | PostgreSQL (planned) |
| router-service | In-memory | - |
| provider-service | In-memory | - |
| gateway-service | N/A (no storage) | - |
| monitor-service | In-memory | - |
