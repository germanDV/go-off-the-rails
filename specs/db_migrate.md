# Database Migration Tool Specification

## Overview
A minimalistic Go package for running SQL database migrations within the `gotr` tool.

## Goals
- Simple package with no external dependencies beyond the database driver
- Transaction safety (each migration runs in a transaction)
- Sequential number-based ordering
- Migration tracking via database table

There's no _up_ vs _down_ distinction because down migrations are sensitive and not even possible many times. If you want to undo something, simply write a new migration with a DROP statement or whatever the case may be.

## Directory Structure
```
db/migrations/
  001_create_users_table.sql
  002_add_user_email_index.sql
```

## Migration File Naming Convention
- Format: `{VERSION}_{DESCRIPTION}.sql`
- `VERSION`: Sequential number (e.g., 001, 002, 003)
- `DESCRIPTION`: Short, snake_case description

## Database Schema

### Migration Tracking Table
```sql
CREATE TABLE IF NOT EXISTS schema_migrations (
    version VARCHAR(255) PRIMARY KEY,
    applied_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);
```

## Package API

```go
package migrate

// Migrator handles database migrations
type Migrator struct {
    db          *sql.DB
    migrationsDir string
    tableName   string
}

// New creates a new Migrator instance
func New(db *sql.DB, migrationsDir string) *Migrator

// Up runs all pending migrations and returns the applied ones
func (m *Migrator) Up() ([]MigrationStatus, error)

// Down rolls back the single last applied migration
func (m *Migrator) Down() error

// DownAll rolls back all applied migrations
func (m *Migrator) DownAll() error

// Status returns the current migration status
func (m *Migrator) Status() ([]MigrationStatus, error)
```

## Implementation Details

### Core Components

1. **Migrator** - Main migration engine
   - Scans migration directory
   - Parses migration files
   - Executes migrations
   - Manages transaction safety

2. **Migration Record** - Tracks applied migrations
   - Version string
   - Applied timestamp

### Key Behaviors

- **Migrations**: Execute in version order, skip already applied
- **Transactions**: Each migration runs in its own transaction
- **Locking**: Use advisory locks to prevent concurrent migrations
- **Error handling**: Stop on first error, leave transaction rolled back

### Example Usage

```go
package main

import (
    "database/sql"
    "log"

    "github.com/yourproject/gotr/db/migrate"
    _ "github.com/mattn/go-sqlite3"
)

func main() {
    db, err := sql.Open("sqlite3", "./app.db")
    if err != nil {
        log.Fatal(err)
    }
    defer db.Close()

    migrator := migrate.New(db, "./db/migrations")

    // Run all pending migrations
    if err := migrator.Run(); err != nil {
        log.Fatal(err)
    }

    // Check status
    status, err := migrator.Status()
    if err != nil {
        log.Fatal(err)
    }
    for _, s := range status {
        log.Printf("Migration %s: applied=%v at %s\n", s.Version, s.Applied, s.AppliedAt)
    }
}
```

## Database Support
- SQLite (primary target, as used by gotr)

## Error Handling
- Invalid migration files: Return error before execution
- Database connection errors: Fail fast with clear message
- Migration failure: Rollback transaction, report which migration failed


