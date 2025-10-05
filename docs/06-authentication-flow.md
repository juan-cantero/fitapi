# Authentication Flow & Security

## Overview

This document explains how authentication works in the Fitness API, including JWT tokens, Supabase Auth integration, and security architecture.

## Table of Contents

1. [Authentication Architecture](#authentication-architecture)
2. [Supabase Keys Explained](#supabase-keys-explained)
3. [Complete Auth Flow](#complete-auth-flow)
4. [JWT Tokens Deep Dive](#jwt-tokens-deep-dive)
5. [Middleware Implementation](#middleware-implementation)
6. [Security Model](#security-model)
7. [RLS vs API Security](#rls-vs-api-security)
8. [Testing Authentication](#testing-authentication)

---

## Authentication Architecture

### Our Approach: API Gateway Pattern

```
┌─────────────┐
│  Frontend   │  Step 1: Login via Supabase Auth
│  (Browser)  │  ────────────────────────────────►  Supabase Auth
└──────┬──────┘                                      (validates password)
       │                                                     │
       │  Step 2: Receives JWT token  ◄──────────────────────┘
       │
       │  Step 3: Sends requests with token
       │  Authorization: Bearer <token>
       ↓
┌─────────────┐
│  Your API   │  Step 4: Validates token
│  (Go/Gin)   │  Step 5: Extracts user_id from token
└──────┬──────┘  Step 6: Enforces authorization
       │
       │  Step 7: Queries with validated user_id
       ↓
┌─────────────┐
│  Database   │  Step 8: Returns data
│ (Postgres)  │
└─────────────┘
```

### Why This Architecture?

**✅ Pros:**
- Full control over business logic
- Database agnostic (can switch from Supabase later)
- Security enforced in server code (easier to understand)
- Industry standard approach
- Better for complex authorization rules

**❌ Alternative (Direct Database Access):**
- Frontend → Supabase Client → Database
- Requires Row Level Security (RLS)
- Less control over business logic
- Locked into Supabase

---

## Supabase Keys Explained

Supabase provides three types of keys. Understanding them is crucial for security.

### 1. Anon Key (Public - Safe to Expose)

**What it is:**
```
eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpc3MiOiJzdXBhYmFzZSIsInJlZiI6ImNyZ3pvdW1id3Fva2hqYWtvbGd1Iiwicm9sZSI6ImFub24iLCJpYXQiOjE3NTk1NDE1NDMsImV4cCI6MjA3NTExNzU0M30.7kONuMC8oL1DNKPVs0_8gD0fcrsDVGmGKC9DjqGr2MQ
```

**Purpose:**
- Frontend uses it to call Supabase services
- Identifies which Supabase project you're using

**Can be exposed?**
- ✅ YES - It's meant to be public
- Lives in frontend code
- Can be seen in browser DevTools

**What it can do:**
- Call Auth API (login, signup, logout)
- Query database (if RLS policies allow)
- Use Storage API (if policies allow)

**What it CANNOT do:**
- Bypass RLS policies
- Access data without proper authentication
- Perform admin operations

**Think of it as:** A visitor badge that gets you into the building, but RLS/policies control which rooms you can enter.

**Where to find it:**
- Supabase Dashboard → Settings → API → Project API keys → `anon` `public`

**Usage:**
```javascript
// Frontend
const supabase = createClient(
  'https://project.supabase.co',
  'ANON_KEY_HERE'  // ← Public, safe to expose
)
```

### 2. Service Role Key (Secret - NEVER EXPOSE)

**What it is:**
```
eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpc3MiOiJzdXBhYmFzZSIsInJlZiI6InByb2plY3QiLCJyb2xlIjoic2VydmljZV9yb2xlIiwiaWF0IjoxNjUwMDAwMDAwLCJleHAiOjE5NjU1MDAwMDB9.SECRET_SIGNATURE
```

**Purpose:**
- Backend admin operations
- Bypass all RLS policies
- System-level access

**Can be exposed?**
- ❌ NO - Keep it secret!
- Server-side only
- Never commit to git
- Never send to frontend

**What it can do:**
- Everything (bypass all security)
- Admin operations
- Direct database access
- Manage users

**Think of it as:** Master key that opens everything.

**Where to find it:**
- Supabase Dashboard → Settings → API → Project API keys → `service_role` `secret`

**Usage:**
```go
// Backend only - for admin operations
supabaseClient := supa.NewClient(url, SERVICE_ROLE_KEY, nil)
// Can bypass RLS, manage users, etc.
```

**⚠️ WARNING:** Never use this key for regular API operations. Only for admin tasks.

### 3. JWT Secret (Secret - NEVER EXPOSE)

**What it is:**
```
your-super-secret-jwt-secret-never-share-this
```

**Purpose:**
- Sign and verify JWT tokens
- Proves tokens are authentic

**Can be exposed?**
- ❌ NO - Keep it secret!
- Server-side only

**What it's used for:**
- Supabase uses it to sign tokens when users login
- Your API uses it to verify tokens are authentic

**Where to find it:**
- Supabase Dashboard → Settings → API → JWT Settings → JWT Secret

**Usage:**
```go
// Validate JWT manually (alternative to using Supabase client)
token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
    return []byte(jwtSecret), nil
})
```

### Summary Table

| Key Type | Safe to Expose? | Used By | Purpose |
|----------|----------------|---------|---------|
| **Anon Key** | ✅ Yes | Frontend | Identify project, call public APIs |
| **Service Role Key** | ❌ No | Backend (admin only) | Bypass security for admin tasks |
| **JWT Secret** | ❌ No | Backend | Verify token authenticity |

---

## Complete Auth Flow

### Step 1: User Registration (Frontend → Supabase)

**Frontend code:**
```javascript
const { data, error } = await supabase.auth.signUp({
  email: 'user@example.com',
  password: 'password123'
})
```

**What happens:**
1. Request sent to Supabase Auth API
2. Supabase hashes password with bcrypt
3. Creates user in `auth.users` table
4. Sends confirmation email (if enabled)
5. Returns JWT token

**Response:**
```json
{
  "access_token": "eyJhbGci...",
  "token_type": "bearer",
  "expires_in": 3600,
  "refresh_token": "...",
  "user": {
    "id": "uuid-here",
    "email": "user@example.com"
  }
}
```

### Step 2: User Login (Frontend → Supabase)

**Frontend code:**
```javascript
const { data, error } = await supabase.auth.signInWithPassword({
  email: 'user@example.com',
  password: 'password123'
})
```

**What happens:**
1. Request sent to Supabase Auth API
2. Supabase retrieves password hash from database
3. Compares submitted password with hash
4. If match: creates JWT token (signed with JWT secret)
5. If no match: returns error

**Response:** Same as signup (includes `access_token`)

### Step 3: Store Token (Frontend)

**Frontend code:**
```javascript
const token = data.session.access_token
localStorage.setItem('supabase_token', token)
// or use cookies, sessionStorage, etc.
```

**Important:** Token expires after 1 hour by default. Use refresh token to get new one.

### Step 4: Make Authenticated Request (Frontend → Your API)

**Frontend code:**
```javascript
const token = localStorage.getItem('supabase_token')

fetch('http://localhost:8080/api/exercises', {
  method: 'POST',
  headers: {
    'Authorization': `Bearer ${token}`,
    'Content-Type': 'application/json'
  },
  body: JSON.stringify({
    name: 'Bench Press',
    description: 'Chest exercise'
  })
})
```

**Request format:**
```http
POST /api/exercises HTTP/1.1
Host: localhost:8080
Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...
Content-Type: application/json

{
  "name": "Bench Press",
  "description": "Chest exercise"
}
```

### Step 5: Validate Token (Your API Middleware)

**Go middleware code:**
```go
func AuthRequired(supabaseClient *supa.Client) gin.HandlerFunc {
    return func(c *gin.Context) {
        // 1. Extract token from Authorization header
        authHeader := c.GetHeader("Authorization")
        if authHeader == "" {
            c.JSON(401, gin.H{"error": "missing authorization header"})
            c.Abort()
            return
        }

        // 2. Remove "Bearer " prefix
        token := strings.TrimPrefix(authHeader, "Bearer ")

        // 3. Validate token with Supabase
        user, err := supabaseClient.Auth.User(c.Request.Context(), token)
        if err != nil {
            c.JSON(401, gin.H{"error": "invalid token"})
            c.Abort()
            return
        }

        // 4. Store user_id in context
        c.Set("user_id", user.ID)
        c.Set("user_email", user.Email)

        // 5. Continue to handler
        c.Next()
    }
}
```

**What validation checks:**
- ✅ Token signature is valid (proves it's from Supabase)
- ✅ Token not expired
- ✅ Token not revoked
- ✅ User exists

### Step 6: Use User ID in Handler (Your API)

**Go handler code:**
```go
func (h *ExerciseHandler) Create(c *gin.Context) {
    // 1. Get validated user_id from context (set by middleware)
    userID, exists := c.Get("user_id")
    if !exists {
        c.JSON(401, gin.H{"error": "unauthorized"})
        return
    }

    // 2. Parse request body
    var req CreateExerciseRequest
    if err := c.ShouldBindJSON(&req); err != nil {
        c.JSON(400, gin.H{"error": err.Error()})
        return
    }

    // 3. Create exercise with user_id from token
    exercise := &Exercise{
        Name:        req.Name,
        Description: req.Description,
        UserID:      userID.(string), // ← From validated token!
    }

    // 4. Save to database
    if err := h.service.CreateExercise(c.Request.Context(), exercise); err != nil {
        c.JSON(500, gin.H{"error": err.Error()})
        return
    }

    c.JSON(201, exercise)
}
```

**Key point:** `user_id` comes from the **validated token**, not from the request body. User cannot fake this!

### Step 7: Query Database (Your API)

**Repository code:**
```go
func (r *ExerciseRepository) Create(ctx context.Context, exercise *Exercise) error {
    query := `
        INSERT INTO exercises (id, name, description, user_id, created_at, updated_at)
        VALUES (gen_random_uuid(), $1, $2, $3, NOW(), NOW())
        RETURNING id, created_at, updated_at
    `

    return r.db.Pool.QueryRow(
        ctx,
        query,
        exercise.Name,
        exercise.Description,
        exercise.UserID, // ← From validated token
    ).Scan(&exercise.ID, &exercise.CreatedAt, &exercise.UpdatedAt)
}
```

### Step 8: Return Response (Your API → Frontend)

**Response:**
```json
{
  "id": "uuid-here",
  "name": "Bench Press",
  "description": "Chest exercise",
  "user_id": "user-uuid",
  "is_public": false,
  "created_at": "2025-10-04T10:00:00Z",
  "updated_at": "2025-10-04T10:00:00Z"
}
```

---

## JWT Tokens Deep Dive

### What is a JWT?

**JWT (JSON Web Token)** is a secure way to transmit information between parties as a JSON object.

### JWT Structure

A JWT has three parts separated by dots (`.`):

```
eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIxMjM0NTY3ODkwIiwibmFtZSI6IkpvaG4gRG9lIiwiaWF0IjoxNTE2MjM5MDIyfQ.SflKxwRJSMeKKF2QT4fwpMeJf36POk6yJV_adQssw5c

    Header (red)              Payload (purple)           Signature (cyan)
```

#### 1. Header

**Base64 encoded JSON:**
```json
{
  "alg": "HS256",  // Algorithm used (HMAC SHA-256)
  "typ": "JWT"     // Type: JWT
}
```

**Purpose:** Describes how the token is signed.

#### 2. Payload (Claims)

**Base64 encoded JSON:**
```json
{
  "sub": "923ac887-1234-5678-9abc-def012345678",  // Subject (user_id)
  "email": "user@example.com",
  "role": "authenticated",
  "aud": "authenticated",                         // Audience
  "iat": 1759541543,                             // Issued at (timestamp)
  "exp": 1759545143                              // Expires at (timestamp)
}
```

**Standard claims:**
- `sub` - Subject (user ID)
- `iat` - Issued at (when token was created)
- `exp` - Expiration time
- `aud` - Audience (who the token is for)

**Custom claims:** Can add any data you need.

**Important:** Payload is **NOT encrypted**, only encoded. Anyone can decode it (but can't modify it without breaking signature).

#### 3. Signature

**How it's created:**
```javascript
HMACSHA256(
  base64UrlEncode(header) + "." + base64UrlEncode(payload),
  JWT_SECRET  // ← Only Supabase knows this!
)
```

**Purpose:**
- Proves token is authentic
- Detects tampering
- Can only be created by someone with the JWT secret

### Why JWT is Secure

#### Attack 1: Modify Payload

**Hacker tries:**
```javascript
// Decode the token
const payload = base64Decode(token.split('.')[1])
// { sub: "hacker-id", ... }

// Change user_id
payload.sub = "victim-id"

// Re-encode
const newPayload = base64Encode(payload)

// Create fake token
const fakeToken = header + "." + newPayload + "." + signature
```

**What happens:**
```go
// Your API validates token
token, err := jwt.Parse(fakeToken, func(token *jwt.Token) (interface{}, error) {
    return []byte(jwtSecret), nil
})

// JWT library recalculates signature
realSignature = HMACSHA256(header + "." + newPayload, jwtSecret)

// Compares signatures
if realSignature != providedSignature {
    return error  // ❌ TAMPERING DETECTED!
}
```

**Result:** Token rejected because signature doesn't match modified payload.

#### Attack 2: Create Fake Token

**Hacker tries:**
```javascript
const fakeToken = {
  header: { alg: "HS256", typ: "JWT" },
  payload: { sub: "victim-id", email: "victim@example.com" },
  signature: "hacker-random-string"
}
```

**What happens:**
```go
// Your API tries to verify
realSignature = HMACSHA256(header + payload, jwtSecret)
// realSignature = "abc123xyz"

// Compare with provided signature
providedSignature = "hacker-random-string"

if realSignature != providedSignature {
    return error  // ❌ INVALID SIGNATURE!
}
```

**Result:** Without the JWT secret, hacker can't create valid signatures.

#### Attack 3: Steal Token

**Hacker steals valid token** (e.g., via XSS attack)

**What happens:**
- ✅ Token is valid (real signature)
- ✅ Hacker can use it until it expires
- ⚠️ Token expires after 1 hour
- ✅ User can revoke token (logout)

**Mitigation:**
- Use HTTPS only (prevent interception)
- Short expiration times
- HttpOnly cookies (prevent XSS access)
- Implement logout/token revocation

### Token Expiration

**Default:** Tokens expire after 3600 seconds (1 hour)

**After expiration:**
- API rejects the token
- Frontend must use refresh token to get new access token

**Refresh token flow:**
```javascript
// Access token expired, use refresh token
const { data, error } = await supabase.auth.refreshSession({
  refresh_token: storedRefreshToken
})

// Get new access token
const newToken = data.session.access_token
```

---

## Middleware Implementation

### Purpose of Middleware

Middleware runs **before** your handler, allowing you to:
- Validate authentication
- Extract user information
- Reject unauthorized requests
- Share data with handlers via context

### Auth Middleware Structure

```go
// internal/middleware/auth.go
package middleware

import (
    "strings"
    "github.com/gin-gonic/gin"
    supa "github.com/supabase-community/supabase-go"
)

func AuthRequired(supabaseClient *supa.Client) gin.HandlerFunc {
    return func(c *gin.Context) {
        // Authentication logic here
        // If valid: c.Set("user_id", userID) and c.Next()
        // If invalid: c.JSON(401, ...) and c.Abort()
    }
}
```

### Using Middleware in Routes

```go
// cmd/api/main.go
router := gin.Default()

// Public routes (no auth)
router.GET("/health", healthHandler)
router.GET("/api/exercises/public", listPublicExercises)

// Protected routes (auth required)
api := router.Group("/api")
api.Use(middleware.AuthRequired(supabaseClient)) // ← Apply middleware
{
    // All these routes require authentication
    api.POST("/exercises", createExercise)
    api.GET("/exercises", listExercises)
    api.PUT("/exercises/:id", updateExercise)
    api.DELETE("/exercises/:id", deleteExercise)
}
```

### Middleware Execution Flow

```
Request arrives
    ↓
Gin router matches route
    ↓
Middleware.AuthRequired runs
    ├─ Validates token
    ├─ Extracts user_id
    ├─ Stores in context: c.Set("user_id", userID)
    └─ Calls c.Next()
    ↓
Handler runs (createExercise)
    ├─ Gets user_id: c.Get("user_id")
    ├─ Uses it in business logic
    └─ Returns response
    ↓
Response sent to client
```

### Context Sharing

**Middleware sets data:**
```go
c.Set("user_id", user.ID)
c.Set("user_email", user.Email)
c.Set("user_role", user.Role)
```

**Handler retrieves data:**
```go
userID, exists := c.Get("user_id")
if !exists {
    c.JSON(401, gin.H{"error": "unauthorized"})
    return
}

// Use it
fmt.Printf("Request from user: %s\n", userID)
```

---

## Security Model

### Our Security Layers

```
┌──────────────────────────────────────┐
│  Layer 1: HTTPS                      │  ← Encryption in transit
│  (Prevents man-in-the-middle)        │
└──────────────────────────────────────┘
             ↓
┌──────────────────────────────────────┐
│  Layer 2: JWT Token Validation       │  ← Authentication
│  (Proves user identity)              │
└──────────────────────────────────────┘
             ↓
┌──────────────────────────────────────┐
│  Layer 3: Authorization in Code      │  ← Access control
│  (Checks permissions)                │
└──────────────────────────────────────┘
             ↓
┌──────────────────────────────────────┐
│  Layer 4: Database Queries           │  ← Data filtering
│  (Filter by user_id)                 │
└──────────────────────────────────────┘
```

### Authentication vs Authorization

**Authentication:** WHO are you?
- "Prove you're user-123"
- Handled by JWT validation
- Middleware checks token

**Authorization:** WHAT can you do?
- "Can user-123 delete this exercise?"
- Handled in handler code
- Business logic checks ownership

### Example: Deleting an Exercise

```go
func (h *ExerciseHandler) Delete(c *gin.Context) {
    // Authentication (who are you?)
    userID := c.Get("user_id").(string)  // From middleware

    // Get exercise ID from URL
    exerciseID := c.Param("id")

    // Get exercise from database
    exercise, err := h.repo.GetByID(ctx, exerciseID)
    if err != nil {
        c.JSON(404, gin.H{"error": "exercise not found"})
        return
    }

    // Authorization (can you do this?)
    if exercise.UserID != userID {
        c.JSON(403, gin.H{"error": "forbidden: you don't own this exercise"})
        return
    }

    // Perform action
    if err := h.repo.Delete(ctx, exerciseID); err != nil {
        c.JSON(500, gin.H{"error": err.Error()})
        return
    }

    c.JSON(204, nil)
}
```

### Common Authorization Patterns

#### 1. User Owns Resource

```go
if resource.UserID != userID {
    return ErrForbidden
}
```

#### 2. Public or Owned by User

```go
query := `
    SELECT * FROM exercises
    WHERE is_public = true OR user_id = $1
`
```

#### 3. Admin Only

```go
userRole := c.Get("user_role").(string)
if userRole != "admin" {
    return ErrForbidden
}
```

---

## RLS vs API Security

### Row Level Security (RLS)

**What it is:** PostgreSQL feature that filters rows based on policies.

**When you need it:** When frontend talks **directly** to database.

```
Frontend → Supabase Client → Database (with RLS)
```

**Example RLS policy:**
```sql
-- Users can only see their own exercises
CREATE POLICY "users_own_exercises"
ON exercises
FOR SELECT
USING (auth.uid() = user_id);
```

**How it works:**
1. Frontend sends query: `SELECT * FROM exercises WHERE user_id = 'victim-id'`
2. RLS adds condition: `AND auth.uid() = user_id`
3. Even if query specifies different user_id, RLS enforces real user's ID

### API Security (Our Approach)

**What it is:** Security enforced in your server code.

**When you use it:** When frontend talks to **your API**.

```
Frontend → Your API → Database (no RLS needed)
```

**Example API security:**
```go
// Handler enforces security
userID := getValidatedUserID(token)  // From token, can't be faked
db.Query("SELECT * FROM exercises WHERE user_id = $1", userID)
```

**How it works:**
1. Middleware validates token (can't be bypassed)
2. Extracts real user_id from token (can't be faked)
3. Handler uses validated user_id in queries
4. User can't tamper with server code

### Comparison

| Aspect | RLS (Direct DB Access) | API Security (Our Approach) |
|--------|------------------------|------------------------------|
| **Where enforced** | Database | Server code |
| **Can be bypassed?** | No (enforced by Postgres) | No (enforced by server) |
| **Easy to understand?** | Complex SQL policies | Clear code logic |
| **Flexibility** | Limited to SQL | Full programming language |
| **Best for** | Simple apps, rapid prototyping | Complex business logic |
| **Required for us** | ❌ No | ✅ Yes |

### Why We Don't Use RLS

**Our architecture:**
- Frontend → API → Database
- API validates all requests
- Security in code (easier to maintain)
- More control over authorization logic
- Can switch databases later (not locked into Supabase)

**If we used RLS:**
- Would be redundant (API already secures)
- Adds complexity
- Harder to debug
- Less flexible

---

## Testing Authentication

### Tool 1: Get Token Script

We created `cmd/gettoken/main.go` to easily get JWT tokens for testing.

**Usage:**
```bash
# Create user and get token (default: test@example.com / test123456)
go run cmd/gettoken/main.go

# Create user with custom credentials
go run cmd/gettoken/main.go myemail@test.com mypassword123
```

**Output:**
```
✅ User created successfully!

🎉 Authentication successful!

📋 Copy this token for testing:
─────────────────────────────────────────────────────────
eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...
─────────────────────────────────────────────────────────

👤 User ID: 923ac887-1234-5678-9abc-def012345678
📧 Email: test@example.com
⏰ Expires in: 3600 seconds
```

### Tool 2: curl Commands

**Get token and save to variable:**
```bash
TOKEN=$(go run cmd/gettoken/main.go 2>/dev/null | grep -A1 "Copy this token" | tail -1 | tr -d ' ')
```

**Test authenticated endpoint:**
```bash
curl http://localhost:8080/api/exercises \
  -H "Authorization: Bearer $TOKEN"
```

**Create exercise:**
```bash
curl -X POST http://localhost:8080/api/exercises \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Bench Press",
    "description": "Chest exercise",
    "is_public": false
  }'
```

### Tool 3: Postman Collection

**Setup:**
1. Create new collection: "Fitness API"
2. Add environment variable: `token`
3. Create request: "Get Exercises"
   - Method: GET
   - URL: `http://localhost:8080/api/exercises`
   - Headers: `Authorization: Bearer {{token}}`

**Get token:**
```bash
go run cmd/gettoken/main.go
# Copy token to Postman environment variable
```

### Tool 4: Decode JWT (Debug)

**Online:** https://jwt.io

**Command line:**
```bash
# Get token
TOKEN=$(go run cmd/gettoken/main.go 2>/dev/null | grep -A1 "Copy this token" | tail -1 | tr -d ' ')

# Decode payload (second part of JWT)
echo $TOKEN | cut -d'.' -f2 | base64 -d | jq
```

**Output:**
```json
{
  "sub": "923ac887-1234-5678-9abc-def012345678",
  "email": "test@example.com",
  "role": "authenticated",
  "iat": 1759541543,
  "exp": 1759545143
}
```

---

## Summary

### Key Concepts

1. **Supabase Auth** handles user registration and login
2. **JWT tokens** prove user identity (cryptographically signed)
3. **Anon key** is public, used by frontend to call Supabase
4. **JWT secret** is private, used to validate tokens
5. **Middleware** validates tokens before handlers run
6. **user_id** extracted from validated token (can't be faked)
7. **Authorization** enforced in handler code
8. **RLS not needed** because API enforces security

### Security Flow

```
User logs in (password validated by Supabase)
    ↓
Receives JWT token (signed by Supabase)
    ↓
Sends token with every request
    ↓
Middleware validates token (checks signature)
    ↓
Extracts user_id from token
    ↓
Handler uses user_id for authorization
    ↓
Database query filtered by user_id
    ↓
Returns only authorized data
```

### Best Practices

✅ **Do:**
- Always validate tokens in middleware
- Use user_id from validated token (never from request)
- Use HTTPS in production
- Set short token expiration times
- Implement logout/token revocation
- Store tokens securely (HttpOnly cookies preferred)

❌ **Don't:**
- Trust user_id from request body/headers
- Expose service role key or JWT secret
- Skip token validation
- Store tokens in localStorage (XSS risk)
- Use the same token for different environments

### Next Steps

Now that you understand authentication:
1. Implement auth middleware (Epic 3)
2. Protect your API endpoints
3. Build frontend that integrates with flow
4. Add refresh token handling
5. Implement logout functionality

Authentication is the foundation of your API security - get it right, and everything else follows!
