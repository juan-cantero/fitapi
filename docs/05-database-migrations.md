# Database Migrations

## What is a Database Migration?

A **database migration** is a version-controlled way to make changes to your database schema over time.

### Why Use Migrations?

**Problem without migrations:**
- You manually run SQL commands in production
- No history of what changed and when
- Hard to sync database changes across team members
- Difficult to rollback if something goes wrong
- Development and production databases drift apart

**Solution with migrations:**
- All database changes are tracked in code (version control)
- Changes are reproducible and testable
- Easy to apply/rollback changes
- Team members stay in sync
- Clear history of schema evolution

### Real-World Example

Imagine you're building the fitness API:

**Week 1:** You create the `exercises` table
**Week 2:** You realize you need to add an `image_url` column
**Week 3:** You need to add an index on `is_public` for performance

Without migrations, you'd manually run SQL commands and hope you remember what you did. With migrations, each change is a numbered file that can be applied in order.

## Migration Files Explained

### File Naming Convention

```
001_create_equipment_table.up.sql      # Apply migration (going "up")
001_create_equipment_table.down.sql    # Rollback migration (going "down")
002_create_exercises_table.up.sql
002_create_exercises_table.down.sql
...
```

**Format:** `{version}_{description}.{direction}.sql`

- `version`: Sequential number (001, 002, 003...)
- `description`: What this migration does (use snake_case)
- `direction`: `up` (apply) or `down` (rollback)

### Up Migration (*.up.sql)

This file contains SQL to **apply** the change.

**Example: 001_create_equipment_table.up.sql**
```sql
CREATE TABLE equipment (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name TEXT NOT NULL,
    ...
);

CREATE INDEX idx_equipment_user_id ON equipment(user_id);
```

**When you run:** Adds the table and index to your database.

### Down Migration (*.down.sql)

This file contains SQL to **undo** the change.

**Example: 001_create_equipment_table.down.sql**
```sql
DROP TRIGGER IF EXISTS update_equipment_updated_at ON equipment;
DROP TABLE IF EXISTS equipment CASCADE;
```

**When you run:** Removes the table (rollback to previous state).

### Migration Sequence

Migrations are applied in **numerical order**:

```
001 → 002 → 003 → 004 → 005 → 006 → 007
```

Each migration builds on the previous one. You can't apply migration 005 without first applying 001-004.

## How golang-migrate Works

### The Migration Tool

We use **[golang-migrate](https://github.com/golang-migrate/migrate)**, a popular Go migration library.

### Migration State Tracking

golang-migrate creates a special table in your database:

```sql
CREATE TABLE schema_migrations (
    version BIGINT PRIMARY KEY,
    dirty BOOLEAN NOT NULL
);
```

**Fields:**
- `version`: Current migration version (e.g., 7 after running all our migrations)
- `dirty`: `false` = migration completed successfully, `true` = migration failed mid-way (needs manual fix)

### How It Knows What to Run

1. **Reads migration files** in `migrations/` directory
2. **Checks current version** in `schema_migrations` table
3. **Applies pending migrations** that are newer than current version
4. **Updates version** after each successful migration

**Example:**

```
Current version: 3
Migration files: 001, 002, 003, 004, 005
Pending: 004, 005

Run migrations → Applies 004, then 005
New version: 5
```

## The `cmd/migrate/main.go` Program

### Purpose

A custom Go program to run migrations programmatically (instead of using the CLI).

### Code Walkthrough

```go
package main

import (
    "database/sql"
    "fmt"
    "log"
    "os"

    "github.com/golang-migrate/migrate/v4"
    "github.com/golang-migrate/migrate/v4/database/postgres"
    _ "github.com/golang-migrate/migrate/v4/source/file"
    "github.com/joho/godotenv"
    _ "github.com/lib/pq"
)
```

**Imports explained:**
- `database/sql`: Standard Go database interface
- `migrate/v4`: Migration library
- `database/postgres`: Postgres-specific driver for migrations
- `source/file`: Reads migration files from filesystem
- `godotenv`: Loads `.env` file
- `lib/pq`: Postgres driver (imported for side effects with `_`)

### Step-by-Step Execution

#### 1. Load Environment Variables

```go
if err := godotenv.Load(); err != nil {
    log.Println("No .env file found")
}

databaseURL := os.Getenv("DATABASE_URL")
if databaseURL == "" {
    log.Fatal("DATABASE_URL not set")
}
```

**What it does:**
- Loads `.env` file to get `DATABASE_URL`
- Exits if `DATABASE_URL` not found

#### 2. Open Database Connection

```go
db, err := sql.Open("postgres", databaseURL)
if err != nil {
    log.Fatalf("Failed to connect to database: %v", err)
}
defer db.Close()
```

**What it does:**
- Creates connection to Postgres using `lib/pq` driver
- `defer db.Close()` ensures connection is closed when program exits

#### 3. Create Migration Driver

```go
driver, err := postgres.WithInstance(db, &postgres.Config{})
if err != nil {
    log.Fatalf("Failed to create driver: %v", err)
}
```

**What it does:**
- Wraps the database connection with Postgres-specific migration features
- Handles Postgres-specific SQL syntax and behaviors

#### 4. Create Migration Instance

```go
m, err := migrate.NewWithDatabaseInstance(
    "file://migrations",  // Source: read from migrations/ directory
    "postgres",           // Database name
    driver,               // Database driver from step 3
)
if err != nil {
    log.Fatalf("Failed to create migrate instance: %v", err)
}
```

**What it does:**
- `file://migrations`: Tells migrate to read `.sql` files from `migrations/` folder
- Combines the file source with database driver
- Creates migration controller

#### 5. Run Migrations

```go
if err := m.Up(); err != nil && err != migrate.ErrNoChange {
    log.Fatalf("Failed to run migrations: %v", err)
}
```

**What it does:**
- `m.Up()`: Applies all pending migrations
- `migrate.ErrNoChange`: Ignored (means already up-to-date, not an error)
- Any other error: Fatal (stops program)

#### 6. Report Status

```go
version, dirty, err := m.Version()
if err != nil {
    log.Printf("Migration completed successfully")
} else {
    fmt.Printf("Migration completed successfully. Current version: %d, Dirty: %v\n", version, dirty)
}
```

**What it does:**
- Gets current migration version
- Reports success and version number
- Shows if database is in "dirty" state (incomplete migration)

## Running Migrations

### Method 1: Using Our Go Program

```bash
go run cmd/migrate/main.go
```

**Pros:**
- Uses `.env` file automatically
- Cross-platform (works on any OS with Go)
- Can be customized with Go code

**Output:**
```
Migration completed successfully. Current version: 7, Dirty: false
```

### Method 2: Using CLI (if installed)

```bash
migrate -path migrations -database "$DATABASE_URL" up
```

**Requires:**
- `migrate` CLI tool installed
- `DATABASE_URL` environment variable set

## Migration Version Control

### Tracking Migration State

The `schema_migrations` table tracks:

```sql
SELECT * FROM schema_migrations;
```

**Example output:**
```
 version | dirty
---------+-------
       7 | false
```

**Meaning:**
- Database is at version 7 (all 7 migrations applied)
- Not dirty (last migration succeeded)

### What Happens When You Run Migrations

**First time (empty database):**
```
Version: None → 1 → 2 → 3 → 4 → 5 → 6 → 7
```

**Already up-to-date:**
```
Current version: 7
No pending migrations
Output: "no change"
```

**New migration added (008):**
```
Current version: 7
Pending: 008
After running: Version = 8
```

### Dirty State Recovery

**What is a "dirty" migration?**

If a migration fails halfway through:

```sql
-- Migration 008_add_new_column.up.sql
ALTER TABLE exercises ADD COLUMN difficulty INTEGER; -- ✅ Succeeds
ALTER TABLE exercises ADD CONSTRAINT invalid_syntax; -- ❌ Fails
```

**Result:**
- `difficulty` column was added
- Constraint failed
- Migration marked as "dirty"
- Version stuck at 8 (dirty=true)

**How to fix:**

1. **Check what partially applied:**
   ```sql
   \d exercises  -- See if difficulty column exists
   ```

2. **Manually fix the issue:**
   ```sql
   ALTER TABLE exercises DROP COLUMN difficulty;  -- Undo partial changes
   ```

3. **Force version back:**
   ```bash
   migrate -path migrations -database "$DATABASE_URL" force 7
   ```

4. **Fix the migration file** (correct the SQL syntax)

5. **Re-run:**
   ```bash
   go run cmd/migrate/main.go
   ```

## Common Migration Operations

### Check Current Version

**Using our program:**
```go
// Modify cmd/migrate/main.go to just check version:
version, dirty, _ := m.Version()
fmt.Printf("Version: %d, Dirty: %v\n", version, dirty)
```

**Using CLI:**
```bash
migrate -path migrations -database "$DATABASE_URL" version
```

### Rollback Last Migration

**Using CLI:**
```bash
migrate -path migrations -database "$DATABASE_URL" down 1
```

**What happens:**
- Runs the `.down.sql` file for the current version
- Decrements version by 1

**Example:**
```
Before: version = 7
Run: down 1
Executes: 007_create_exercise_logs_table.down.sql
After: version = 6
```

### Rollback All Migrations

```bash
migrate -path migrations -database "$DATABASE_URL" down -all
```

**⚠️ WARNING:** This will drop ALL tables! Only use in development.

### Apply Specific Number of Migrations

**Go up 2 versions:**
```bash
migrate -path migrations -database "$DATABASE_URL" up 2
```

**Example:**
```
Current version: 3
Run: up 2
Applies: 004, 005
New version: 5
```

### Force Version (Recovery)

**When dirty or broken:**
```bash
migrate -path migrations -database "$DATABASE_URL" force 7
```

**What it does:**
- Sets version to 7 without running any migrations
- Marks as not dirty
- Use when you've manually fixed a broken migration

## Creating New Migrations

### Step 1: Create Migration Files

**Naming pattern:**
```bash
touch migrations/008_add_exercise_categories.up.sql
touch migrations/008_add_exercise_categories.down.sql
```

**Version number:** Next sequential number (008 after 007)

### Step 2: Write Up Migration

**migrations/008_add_exercise_categories.up.sql:**
```sql
-- Add categories table
CREATE TABLE exercise_categories (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name TEXT NOT NULL UNIQUE,
    description TEXT,
    created_at TIMESTAMPTZ DEFAULT NOW()
);

-- Add category to exercises
ALTER TABLE exercises ADD COLUMN category_id UUID REFERENCES exercise_categories(id);

-- Create index
CREATE INDEX idx_exercises_category ON exercises(category_id);
```

### Step 3: Write Down Migration

**migrations/008_add_exercise_categories.down.sql:**
```sql
-- Remove category from exercises
ALTER TABLE exercises DROP COLUMN IF EXISTS category_id;

-- Drop categories table
DROP TABLE IF EXISTS exercise_categories CASCADE;
```

### Step 4: Test Migration

**Apply:**
```bash
go run cmd/migrate/main.go
```

**Verify:**
```sql
-- Check tables exist
SELECT * FROM exercise_categories;
SELECT category_id FROM exercises LIMIT 1;
```

**Rollback test:**
```bash
migrate -path migrations -database "$DATABASE_URL" down 1
```

**Verify rollback:**
```sql
-- Tables should be gone
SELECT * FROM exercise_categories;  -- Should error
```

**Re-apply:**
```bash
go run cmd/migrate/main.go
```

## Best Practices

### 1. Never Modify Existing Migrations

❌ **Wrong:**
```
Edit migrations/003_create_exercise_equipment.up.sql
```

✅ **Right:**
```
Create migrations/008_modify_exercise_equipment.up.sql
```

**Why?** Once a migration is applied in production, modifying it creates inconsistency.

### 2. Always Write Down Migrations

Every `.up.sql` needs a corresponding `.down.sql`.

**Why?** Allows rollback if issues found in production.

### 3. Test Migrations Locally First

```bash
# Apply
go run cmd/migrate/main.go

# Test your app

# Rollback
migrate -path migrations -database "$DATABASE_URL" down 1

# Re-apply
go run cmd/migrate/main.go
```

**Why?** Catch errors before production deployment.

### 4. Make Migrations Idempotent (Where Possible)

```sql
-- Good: Won't fail if table exists
CREATE TABLE IF NOT EXISTS equipment (...);

-- Good: Won't fail if column doesn't exist
ALTER TABLE exercises DROP COLUMN IF EXISTS old_column;

-- Good: Safe index creation
CREATE INDEX IF NOT EXISTS idx_name ON table(column);
```

**Why?** Safer re-runs if something goes wrong.

### 5. One Logical Change Per Migration

❌ **Wrong:** One migration that adds 5 tables and modifies 3 existing ones

✅ **Right:** Separate migrations for each major change

**Why?** Easier to debug, rollback, and understand history.

### 6. Add Comments in Migration Files

```sql
-- Add image_url column to store Supabase Storage URLs
ALTER TABLE exercises ADD COLUMN image_url TEXT;

-- Index for faster image lookups
CREATE INDEX idx_exercises_image ON exercises(image_url) WHERE image_url IS NOT NULL;
```

**Why?** Future you (and teammates) will thank you.

### 7. Backup Before Major Migrations

**Before running in production:**
```bash
# Backup database
pg_dump $DATABASE_URL > backup_before_migration_008.sql
```

**Why?** Safety net for critical data.

## Troubleshooting

### Error: "no change"

**Message:** `no change` or `migrate.ErrNoChange`

**Meaning:** Database already up-to-date, no pending migrations.

**Action:** Not an error, all good!

### Error: "dirty database"

**Message:** `Dirty database version X. Fix and force version.`

**Meaning:** Previous migration failed halfway.

**Fix:**
1. Check what version is dirty
2. Manually inspect/fix database
3. Force version: `migrate force X`
4. Fix migration file if needed
5. Re-run

### Error: "file does not exist"

**Message:** `error: file does not exist`

**Causes:**
- Running from wrong directory
- Migration files not in `migrations/` folder

**Fix:**
```bash
# Ensure you're in project root
cd /path/to/fitapi

# Verify files exist
ls migrations/

# Run migration
go run cmd/migrate/main.go
```

### Error: "no such host"

**Message:** `dial tcp: lookup db.xxx.supabase.co: no such host`

**Cause:** Wrong DATABASE_URL (probably direct connection instead of pooler)

**Fix:** Use connection pooler URL:
```
postgresql://postgres.PROJECT:PASSWORD@aws-1-us-east-2.pooler.supabase.com:6543/postgres
```

## Summary

### Key Concepts

1. **Migrations** = Version-controlled database changes
2. **Up migrations** = Apply changes (*.up.sql)
3. **Down migrations** = Rollback changes (*.down.sql)
4. **Version tracking** = `schema_migrations` table stores current version
5. **golang-migrate** = Tool that manages migration execution

### Migration Workflow

```
1. Create migration files (008_description.up.sql + .down.sql)
   ↓
2. Write SQL for up migration (CREATE, ALTER, etc.)
   ↓
3. Write SQL for down migration (DROP, reverse changes)
   ↓
4. Test locally (apply → test → rollback → re-apply)
   ↓
5. Commit to git
   ↓
6. Deploy (run: go run cmd/migrate/main.go)
   ↓
7. Verify in production
```

### Quick Reference

```bash
# Run all pending migrations
go run cmd/migrate/main.go

# Check current version (CLI)
migrate -path migrations -database "$DATABASE_URL" version

# Rollback one migration (CLI)
migrate -path migrations -database "$DATABASE_URL" down 1

# Force version (recovery)
migrate -path migrations -database "$DATABASE_URL" force 7

# Create new migration files
touch migrations/008_name.up.sql migrations/008_name.down.sql
```

Migrations keep your database schema in sync, version-controlled, and safely deployable across all environments!
