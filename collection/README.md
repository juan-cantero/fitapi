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

### 🔐 Automatic Token Generation

The collection includes a **pre-request script** that automatically generates a fresh authentication token before each request by running:
```bash
go run cmd/gettoken/main.go
```

This means you don't need to manually copy/paste tokens - Newman will handle authentication automatically!

## Running Tests

### Run the entire collection:
```bash
newman run collection/fitapi.postman_collection.json
```

### Run with detailed output:
```bash
newman run collection/fitapi.postman_collection.json --verbose
```

### Run a specific request:
```bash
newman run collection/fitapi.postman_collection.json --folder "Get Current User"
```

### Run with custom environment:
You can override variables like base_url:
```bash
newman run collection/fitapi.postman_collection.json \
  --env-var "base_url=http://localhost:3000"
```

## Collection Structure

### Public Endpoints (No Auth Required)
- **Health Check** - `GET /health`
  - Checks API server and database connectivity

### Protected Endpoints (Auth Required)
- **Get Current User** - `GET /api/me`
  - Returns authenticated user's ID and email
  - Uses auto-generated Bearer token

## How Auto-Authentication Works

1. Before each request, the collection runs a pre-request script
2. The script executes `go run cmd/gettoken/main.go`
3. This creates a test user in Supabase and returns an access token
4. The token is stored in the `auth_token` collection variable
5. Protected endpoints automatically use this token in the `Authorization: Bearer` header

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
