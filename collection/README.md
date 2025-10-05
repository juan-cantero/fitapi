# FitAPI Postman Collection

This directory contains the Postman collection for testing the FitAPI endpoints using Newman (Postman CLI).

## Prerequisites

1. **Node.js and npm** installed
2. **Newman** installed globally:
   ```bash
   npm install -g newman
   ```

3. **API Server running**:
   ```bash
   source .env && go run cmd/api/main.go
   ```

## Features

### 🔐 Smart Token Caching

The collection uses shell scripts for intelligent token management:

- **`generate-token.sh`**: Generates and caches authentication tokens
- **`run-tests.sh`**: Wrapper to run Newman with auto-token generation
- **Caches tokens** to `.token_cache.json` with expiration timestamps
- **Reuses valid tokens** - only generates new ones when expired or expiring soon (within 5 minutes)

Supabase tokens expire after **1 hour** by default. The scripts:
1. Check if cached token exists and is still valid
2. Only run `go run cmd/gettoken/main.go --json` when necessary
3. This makes your tests much faster! 🚀

No manual token management needed - just run `./collection/run-tests.sh`!

## Running Tests

### Quick Start (Recommended):
```bash
./collection/run-tests.sh
```

This automatically generates and caches authentication tokens!

### Run with detailed output:
```bash
./collection/run-tests.sh --verbose
```

### Manual Method:
If you prefer to run Newman directly:
```bash
newman run collection/fitapi.postman_collection.json \
  --env-var "auth_token=$(./collection/generate-token.sh)"
```

### Run with custom base URL:
```bash
./collection/run-tests.sh --env-var "base_url=http://localhost:3000"
```

## Collection Structure

### Public Endpoints (No Auth Required)
- **Health Check** - `GET /health`
  - Checks API server and database connectivity

### Protected Endpoints (Auth Required)
- **Get Current User** - `GET /api/me`
  - Returns authenticated user's ID and email
  - Uses auto-generated Bearer token

## How Smart Authentication Works

### Token Generation Flow:

1. **Run tests** with `./collection/run-tests.sh`
2. **Script calls** `generate-token.sh` which:
   - Checks `.token_cache.json` for a valid token
   - **If token is valid** (more than 5 minutes until expiration):
     - Uses cached token (instant! ⚡)
     - Logs: `✅ Using cached token (expires in X minutes)`
   - **If token expired or missing**:
     - Runs `go run cmd/gettoken/main.go --json` to generate fresh token
     - Saves to cache with expiration timestamp
     - Logs: `✅ New token generated and cached (expires in 3600s)`
3. **Token is passed** to Newman via `--env-var "auth_token=..."`
4. **Protected endpoints** automatically use this token in `Authorization: Bearer` header

### Cache File Structure

The cache file (`.token_cache.json`) is gitignored and contains:
```json
{
  "access_token": "eyJ...",
  "expires_in": 3600,
  "expires_at": 1759640378,
  "user_id": "6b37ab1f-b190-4072-9e50-5318d4bad35d",
  "email": "test@example.com"
}
```

## Troubleshooting

### "Failed to generate auth token"
- Make sure your `.env` file is configured correctly
- Verify Supabase credentials are valid
- Check that you can run `go run cmd/gettoken/main.go` manually

### Connection refused
- Ensure the API server is running on port 8080
- Check `PORT` in your `.env` file

### 401 Unauthorized
- The auto-generated token might have expired
- Newman will generate a fresh one on the next request
- Verify `SUPABASE_JWT_SECRET` matches your Supabase project

## Using with Postman GUI

You can also import this collection into Postman desktop app:
1. Open Postman
2. Click **Import**
3. Select `collection/fitapi.postman_collection.json`
4. The auto-authentication will work in the GUI too!

## Adding New Requests

When adding new endpoints to the collection:

1. **Public endpoints**: Set auth to "No Auth"
   ```json
   "auth": {
     "type": "noauth"
   }
   ```

2. **Protected endpoints**: Use the collection's default Bearer token auth (no need to specify anything - it inherits from collection level)

The auto-generated token will be used for all protected endpoints automatically.
