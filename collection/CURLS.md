# cURL Commands for FitAPI Testing

Quick reference for testing the API endpoints with cURL.

## Setup

### 1. Get Authentication Token

```bash
# Generate token (creates test user if doesn't exist)
go run cmd/gettoken/main.go

# Or get JSON output for scripting
TOKEN=$(go run cmd/gettoken/main.go --json | jq -r '.access_token')
```

### 2. Alternative: Use SKIP_AUTH (Development Only)

Set in `.env`:
```bash
SKIP_AUTH=true
```

Then restart server:
```bash
source .env && go run cmd/api/main.go
```

---

## Health Check (No Auth Required)

```bash
curl -X GET http://localhost:8080/health | jq
```

**Expected Response:**
```json
{
  "database": "connected",
  "status": "ok",
  "supabase": true
}
```

---

## Authentication Test

```bash
TOKEN=$(go run cmd/gettoken/main.go --json | jq -r '.access_token')

curl -X GET http://localhost:8080/api/me \
  -H "Authorization: Bearer $TOKEN" | jq
```

**Expected Response:**
```json
{
  "email": "test@example.com",
  "message": "Authentication successful!",
  "user_id": "6b37ab1f-b190-4072-9e50-5318d4bad35d"
}
```

---

## Equipment Endpoints

### Create Equipment

```bash
TOKEN=$(go run cmd/gettoken/main.go --json | jq -r '.access_token')

curl -X POST http://localhost:8080/api/equipment \
  -H "Authorization: Bearer $TOKEN" \
  -H 'Content-Type: application/json' \
  -d '{
    "name": "Barbell",
    "description": "Olympic barbell 20kg"
  }' | jq
```

**Expected Response (201 Created):**
```json
{
  "id": "550e8400-e29b-41d4-a716-446655440000",
  "name": "Barbell",
  "description": "Olympic barbell 20kg",
  "user_id": "6b37ab1f-b190-4072-9e50-5318d4bad35d",
  "created_at": "2025-10-05T13:00:00Z",
  "updated_at": "2025-10-05T13:00:00Z"
}
```

### List All Equipment

```bash
TOKEN=$(go run cmd/gettoken/main.go --json | jq -r '.access_token')

curl -X GET http://localhost:8080/api/equipment \
  -H "Authorization: Bearer $TOKEN" | jq
```

**Expected Response (200 OK):**
```json
[
  {
    "id": "550e8400-e29b-41d4-a716-446655440000",
    "name": "Barbell",
    "description": "Olympic barbell 20kg",
    "user_id": "6b37ab1f-b190-4072-9e50-5318d4bad35d",
    "created_at": "2025-10-05T13:00:00Z",
    "updated_at": "2025-10-05T13:00:00Z"
  }
]
```

### Get Single Equipment

```bash
TOKEN=$(go run cmd/gettoken/main.go --json | jq -r '.access_token')
EQUIPMENT_ID="your-equipment-id-here"

curl -X GET "http://localhost:8080/api/equipment/$EQUIPMENT_ID" \
  -H "Authorization: Bearer $TOKEN" | jq
```

**Expected Response (200 OK):**
```json
{
  "id": "550e8400-e29b-41d4-a716-446655440000",
  "name": "Barbell",
  "description": "Olympic barbell 20kg",
  "user_id": "6b37ab1f-b190-4072-9e50-5318d4bad35d",
  "created_at": "2025-10-05T13:00:00Z",
  "updated_at": "2025-10-05T13:00:00Z"
}
```

### Update Equipment

```bash
TOKEN=$(go run cmd/gettoken/main.go --json | jq -r '.access_token')
EQUIPMENT_ID="your-equipment-id-here"

curl -X PUT "http://localhost:8080/api/equipment/$EQUIPMENT_ID" \
  -H "Authorization: Bearer $TOKEN" \
  -H 'Content-Type: application/json' \
  -d '{
    "name": "Barbell",
    "description": "Olympic barbell 20kg (updated)"
  }' | jq
```

**Expected Response (200 OK):**
```json
{
  "id": "550e8400-e29b-41d4-a716-446655440000",
  "name": "Barbell",
  "description": "Olympic barbell 20kg (updated)",
  "user_id": "6b37ab1f-b190-4072-9e50-5318d4bad35d",
  "created_at": "2025-10-05T13:00:00Z",
  "updated_at": "2025-10-05T13:01:00Z"
}
```

### Delete Equipment

```bash
TOKEN=$(go run cmd/gettoken/main.go --json | jq -r '.access_token')
EQUIPMENT_ID="your-equipment-id-here"

curl -X DELETE "http://localhost:8080/api/equipment/$EQUIPMENT_ID" \
  -H "Authorization: Bearer $TOKEN" \
  -w "\nStatus: %{http_code}\n"
```

**Expected Response (204 No Content):**
```
Status: 204
```

---

## Complete Test Flow

```bash
#!/bin/bash

# 1. Get token
echo "=== Getting Authentication Token ==="
TOKEN=$(go run cmd/gettoken/main.go --json | jq -r '.access_token')
echo "Token obtained"
echo

# 2. Create equipment
echo "=== Creating Equipment ==="
CREATED=$(curl -s -X POST http://localhost:8080/api/equipment \
  -H "Authorization: Bearer $TOKEN" \
  -H 'Content-Type: application/json' \
  -d '{"name":"Dumbbells","description":"20kg pair"}')
echo "$CREATED" | jq
EQUIPMENT_ID=$(echo "$CREATED" | jq -r '.id')
echo

# 3. List equipment
echo "=== Listing Equipment ==="
curl -s -X GET http://localhost:8080/api/equipment \
  -H "Authorization: Bearer $TOKEN" | jq
echo

# 4. Get by ID
echo "=== Getting Equipment by ID ==="
curl -s -X GET "http://localhost:8080/api/equipment/$EQUIPMENT_ID" \
  -H "Authorization: Bearer $TOKEN" | jq
echo

# 5. Update
echo "=== Updating Equipment ==="
curl -s -X PUT "http://localhost:8080/api/equipment/$EQUIPMENT_ID" \
  -H "Authorization: Bearer $TOKEN" \
  -H 'Content-Type: application/json' \
  -d '{"name":"Dumbbells","description":"20kg pair (updated)"}' | jq
echo

# 6. Delete
echo "=== Deleting Equipment ==="
curl -s -X DELETE "http://localhost:8080/api/equipment/$EQUIPMENT_ID" \
  -H "Authorization: Bearer $TOKEN" \
  -w "Status: %{http_code}\n"
echo

echo "=== Test Complete ==="
```

---

## Error Examples

### 401 Unauthorized (Missing Token)

```bash
curl -X GET http://localhost:8080/api/equipment | jq
```

**Response:**
```json
{
  "error": "missing authorization header"
}
```

### 401 Unauthorized (Invalid Token)

```bash
curl -X GET http://localhost:8080/api/equipment \
  -H "Authorization: Bearer invalid-token" | jq
```

**Response:**
```json
{
  "error": "invalid or expired token"
}
```

### 400 Bad Request (Validation Error)

```bash
TOKEN=$(go run cmd/gettoken/main.go --json | jq -r '.access_token')

curl -X POST http://localhost:8080/api/equipment \
  -H "Authorization: Bearer $TOKEN" \
  -H 'Content-Type: application/json' \
  -d '{"name":""}' | jq
```

**Response:**
```json
{
  "error": "Key: 'CreateEquipmentRequest.Name' Error:Field validation for 'Name' failed on the 'required' tag"
}
```

### 404 Not Found

```bash
TOKEN=$(go run cmd/gettoken/main.go --json | jq -r '.access_token')

curl -X GET http://localhost:8080/api/equipment/nonexistent-id \
  -H "Authorization: Bearer $TOKEN" | jq
```

**Response:**
```json
{
  "error": "equipment not found"
}
```

### 403 Forbidden (Not Owner)

Trying to access equipment owned by another user:

**Response:**
```json
{
  "error": "you don't have permission to access this equipment"
}
```

---

## Tips

### Save Token for Multiple Requests

```bash
# Save token to file
go run cmd/gettoken/main.go --json > /tmp/token.json
TOKEN=$(cat /tmp/token.json | jq -r '.access_token')

# Use in multiple requests
curl -H "Authorization: Bearer $TOKEN" http://localhost:8080/api/equipment
curl -H "Authorization: Bearer $TOKEN" http://localhost:8080/api/me
```

### Pretty Print JSON

```bash
# Add | jq at the end
curl http://localhost:8080/health | jq

# Or use jq colors
curl http://localhost:8080/health | jq -C | less -R
```

### Show HTTP Headers

```bash
curl -v http://localhost:8080/health
```

### Show Only Status Code

```bash
curl -s -o /dev/null -w "%{http_code}" http://localhost:8080/health
```

---

## Automated Testing with Newman

Instead of manual cURL commands, use Newman:

```bash
./collection/run-tests.sh
```

See `collection/README.md` for details.
