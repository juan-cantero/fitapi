# Fitness API

A REST API built with Go, Gin, and Supabase for fitness tracking. Manage exercises, workouts, and equipment with image uploads.

## Tech Stack

- **Backend**: Go with Gin framework
- **Database**: Supabase (PostgreSQL)
- **Auth**: Supabase built-in authentication
- **Storage**: Supabase Storage for images
- **Migration**: golang-migrate

## Prerequisites

- Go 1.25+
- Supabase account and project
- golang-migrate CLI

## Project Structure

```
fitapi/
├── cmd/api/          # Application entry point
├── config/           # Configuration management
├── internal/
│   └── database/     # Database connection
├── .env              # Environment variables (not in git)
├── .env.example      # Example environment file
└── PROJECT_PLAN.md   # Development roadmap
```

## Setup

1. **Clone the repository**
   ```bash
   git clone https://github.com/juan-cantero/fitapi.git
   cd fitapi
   ```

2. **Install dependencies**
   ```bash
   go mod download
   ```

3. **Configure environment variables**

   Copy `.env.example` to `.env` and fill in your Supabase credentials:
   ```bash
   cp .env.example .env
   ```

   Edit `.env` with your values:
   ```env
   SUPABASE_URL=your-supabase-project-url
   SUPABASE_KEY=your-supabase-anon-key
   DATABASE_URL=postgresql://postgres:password@db.your-project.supabase.co:5432/postgres
   PORT=8080
   GIN_MODE=debug
   ```

   **Where to find credentials:**
   - `SUPABASE_URL`: Supabase Dashboard → Settings → API → Project URL
   - `SUPABASE_KEY`: Supabase Dashboard → Settings → API → Project API keys → `anon` `public`
   - `DATABASE_URL`: Supabase Dashboard → Settings → Database → Connection String → URI

4. **Run the server**
   ```bash
   go run cmd/api/main.go
   ```

5. **Test the API**
   ```bash
   curl http://localhost:8080/health
   ```

   Expected response:
   ```json
   {"database":"connected","status":"ok","supabase":true}
   ```

## Configuration Files

- **`.env`** - Contains secrets and environment-specific configuration (gitignored)
- **`.env.example`** - Template for environment variables
- **`config/config.go`** - Configuration loader

## Development

See [PROJECT_PLAN.md](PROJECT_PLAN.md) for the complete development roadmap and progress tracking.

## API Documentation

Coming soon...
