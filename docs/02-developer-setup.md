# Developer Setup

## Prerequisites

- **Go**: Version 1.25 or higher
- **Supabase Account**: For database and authentication
- **golang-migrate**: For database migrations
- **Git**: Version control
- **Postman** (optional): For API testing

## Initial Setup

### 1. Clone the Repository

```bash
git clone https://github.com/juan-cantero/fitapi.git
cd fitapi
```

### 2. Install Go Dependencies

```bash
go mod download
```

### 3. Install golang-migrate

```bash
go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest
```

Verify installation:
```bash
migrate -version
```

### 4. Set Up Supabase Project

1. Go to [supabase.com](https://supabase.com)
2. Create a new project
3. Wait for the project to be provisioned

### 5. Configure Environment Variables

Create `.env` file from the example:

```bash
cp .env.example .env
```

Edit `.env` with your Supabase credentials:

```env
# Supabase Configuration
SUPABASE_URL=https://your-project.supabase.co
SUPABASE_KEY=your-anon-key
DATABASE_URL=postgresql://postgres:your-password@db.your-project.supabase.co:5432/postgres

# Server Configuration
PORT=8080
GIN_MODE=debug
```

#### Where to Find Credentials

**SUPABASE_URL**:
- Dashboard → Settings → API → Project URL

**SUPABASE_KEY**:
- Dashboard → Settings → API → Project API keys → `anon` `public`

**DATABASE_URL**:
- Dashboard → Settings → Database → Connection String → URI
- Select "Use connection pooling" with mode "Transaction"

### 6. Run Database Migrations

```bash
# Run all migrations
make migrate-up

# Or manually
migrate -path migrations -database "${DATABASE_URL}" up
```

### 7. Run the Server

```bash
go run cmd/api/main.go
```

Server will start on `http://localhost:8080`

### 8. Test the API

```bash
curl http://localhost:8080/health
```

Expected response:
```json
{"database":"connected","status":"ok","supabase":true}
```

## Development Workflow

### Branch Strategy

- `main` - Production-ready code
- `epic-N-feature-name` - Feature branches for each epic

### Creating a New Feature

```bash
# Pull latest main
git checkout main
git pull origin main

# Create feature branch
git checkout -b epic-N-feature-name

# Make changes, commit
git add .
git commit -m "Your commit message"

# Push to remote
git push -u origin epic-N-feature-name

# Merge to main (after review)
git checkout main
git merge epic-N-feature-name
git push origin main
```

### Running Migrations

**Create a new migration**:
```bash
migrate create -ext sql -dir migrations -seq migration_name
```

**Apply migrations**:
```bash
migrate -path migrations -database "${DATABASE_URL}" up
```

**Rollback last migration**:
```bash
migrate -path migrations -database "${DATABASE_URL}" down 1
```

**Check migration version**:
```bash
migrate -path migrations -database "${DATABASE_URL}" version
```

## Project Structure

```
fitapi/
├── cmd/api/              # Application entry point
│   └── main.go          # Main server file
├── config/              # Configuration management
│   └── config.go        # Env variable loader
├── internal/            # Private application code
│   └── database/        # Database connection
│       └── database.go  # DB pool setup
├── migrations/          # Database migrations (SQL files)
├── docs/               # Documentation
├── .env                # Environment variables (gitignored)
├── .env.example        # Environment template
├── go.mod              # Go dependencies
└── PROJECT_PLAN.md     # Development roadmap
```

## Troubleshooting

### Port Already in Use

```bash
# Kill process on port 8080
fuser -k 8080/tcp

# Or use different port in .env
PORT=8081
```

### Database Connection Issues

1. Verify DATABASE_URL is correct
2. Check Supabase project is running
3. Ensure IP is whitelisted (if applicable)
4. Test connection:
   ```bash
   psql "${DATABASE_URL}"
   ```

### Migration Errors

```bash
# Force version (use carefully)
migrate -path migrations -database "${DATABASE_URL}" force VERSION

# Drop all tables and re-migrate (development only!)
migrate -path migrations -database "${DATABASE_URL}" drop
migrate -path migrations -database "${DATABASE_URL}" up
```

## Useful Commands

```bash
# Run tests
go test ./...

# Run with hot reload (install air first)
air

# Format code
go fmt ./...

# Lint code
golangci-lint run

# Update dependencies
go get -u ./...
go mod tidy
```

## Environment Modes

**Development** (`GIN_MODE=debug`):
- Verbose logging
- Detailed error messages
- Auto-reload recommended

**Production** (`GIN_MODE=release`):
- Minimal logging
- Generic error messages
- Performance optimized
