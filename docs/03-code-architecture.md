# Code Architecture

## Project Structure

```
fitapi/
├── cmd/api/              # Application entry points
│   └── main.go          # Main server initialization
├── config/              # Configuration management
│   └── config.go        # Environment variable handling
├── internal/            # Private application code
│   ├── database/        # Database connection
│   ├── models/          # Data models (TODO)
│   ├── handlers/        # HTTP request handlers (TODO)
│   ├── middleware/      # Custom middleware (TODO)
│   ├── services/        # Business logic (TODO)
│   └── repository/      # Data access layer (TODO)
├── migrations/          # Database migrations
├── docs/               # Documentation
└── pkg/                # Public reusable packages (TODO)
```

## Architecture Pattern

We follow a **Clean Architecture** approach with clear separation of concerns:

```
HTTP Request
    ↓
Handler (HTTP Layer)
    ↓
Service (Business Logic)
    ↓
Repository (Data Access)
    ↓
Database
```

### Layers

#### 1. **Handler Layer** (`internal/handlers/`)
- Handles HTTP requests/responses
- Request validation
- Response formatting
- No business logic

```go
// Example handler structure
func (h *ExerciseHandler) CreateExercise(c *gin.Context) {
    var req CreateExerciseRequest
    if err := c.ShouldBindJSON(&req); err != nil {
        c.JSON(400, gin.H{"error": err.Error()})
        return
    }

    exercise, err := h.service.CreateExercise(c.Request.Context(), req)
    if err != nil {
        c.JSON(500, gin.H{"error": err.Error()})
        return
    }

    c.JSON(201, exercise)
}
```

#### 2. **Service Layer** (`internal/services/`)
- Contains business logic
- Orchestrates operations
- Independent of HTTP concerns

```go
// Example service structure
type ExerciseService struct {
    repo ExerciseRepository
}

func (s *ExerciseService) CreateExercise(ctx context.Context, req CreateExerciseRequest) (*Exercise, error) {
    // Business logic here
    // Validation, transformations, etc.
    return s.repo.Create(ctx, exercise)
}
```

#### 3. **Repository Layer** (`internal/repository/`)
- Database operations only
- SQL queries
- No business logic

```go
// Example repository structure
type ExerciseRepository struct {
    db *database.DB
}

func (r *ExerciseRepository) Create(ctx context.Context, exercise *Exercise) error {
    query := `INSERT INTO exercises (name, description, ...) VALUES ($1, $2, ...)`
    _, err := r.db.Pool.Exec(ctx, query, exercise.Name, exercise.Description, ...)
    return err
}
```

#### 4. **Model Layer** (`internal/models/`)
- Data structures
- Type definitions
- No logic

```go
type Exercise struct {
    ID          string    `json:"id"`
    Name        string    `json:"name"`
    Description string    `json:"description"`
    IsPublic    bool      `json:"is_public"`
    UserID      string    `json:"user_id"`
    ImageURL    *string   `json:"image_url,omitempty"`
    CreatedAt   time.Time `json:"created_at"`
    UpdatedAt   time.Time `json:"updated_at"`
}
```

## Dependency Injection

We use constructor-based dependency injection:

```go
// main.go
func main() {
    cfg := config.Load()
    db := database.New(cfg.DatabaseURL)

    // Repositories
    exerciseRepo := repository.NewExerciseRepository(db)

    // Services
    exerciseService := services.NewExerciseService(exerciseRepo)

    // Handlers
    exerciseHandler := handlers.NewExerciseHandler(exerciseService)

    // Routes
    router := gin.Default()
    router.POST("/api/exercises", exerciseHandler.CreateExercise)
}
```

## Configuration Management

### Environment Variables (`config/config.go`)

```go
type Config struct {
    SupabaseURL string
    SupabaseKey string
    DatabaseURL string
    Port        string
    GinMode     string
}

func Load() *Config {
    godotenv.Load()
    return &Config{
        SupabaseURL: getEnv("SUPABASE_URL", ""),
        // ... other fields
    }
}
```

**Why**: Centralizes configuration, type-safe, easy to test

## Database Connection

### Connection Pool (`internal/database/database.go`)

```go
type DB struct {
    Pool *pgxpool.Pool
}

func New(databaseURL string) (*DB, error) {
    config, err := pgxpool.ParseConfig(databaseURL)
    if err != nil {
        return nil, err
    }

    pool, err := pgxpool.NewWithConfig(context.Background(), config)
    if err != nil {
        return nil, err
    }

    return &DB{Pool: pool}, nil
}
```

**Why**: Connection pooling for performance, graceful shutdown, health checks

## Middleware

### Authentication Middleware (TODO)

```go
func AuthMiddleware(supabase *supa.Client) gin.HandlerFunc {
    return func(c *gin.Context) {
        token := c.GetHeader("Authorization")

        // Validate JWT with Supabase
        user, err := validateToken(token)
        if err != nil {
            c.AbortWithStatusJSON(401, gin.H{"error": "unauthorized"})
            return
        }

        // Store user in context
        c.Set("user_id", user.ID)
        c.Next()
    }
}
```

### Logging Middleware (TODO)

```go
func LoggerMiddleware() gin.HandlerFunc {
    return func(c *gin.Context) {
        start := time.Now()
        c.Next()
        duration := time.Since(start)

        log.Printf("%s %s %d %v", c.Request.Method, c.Request.URL.Path, c.Writer.Status(), duration)
    }
}
```

## Error Handling

### Standard Error Response

```go
type ErrorResponse struct {
    Error   string `json:"error"`
    Message string `json:"message,omitempty"`
    Code    string `json:"code,omitempty"`
}
```

### Error Types

```go
var (
    ErrNotFound      = errors.New("resource not found")
    ErrUnauthorized  = errors.New("unauthorized")
    ErrValidation    = errors.New("validation error")
    ErrInternal      = errors.New("internal server error")
)
```

## Testing Strategy

### Unit Tests
- Test services in isolation
- Mock repositories
- Focus on business logic

### Integration Tests
- Test handlers with real DB
- Use test database
- Test full request/response cycle

### Example Test

```go
func TestCreateExercise(t *testing.T) {
    // Setup
    db := setupTestDB(t)
    repo := repository.NewExerciseRepository(db)
    service := services.NewExerciseService(repo)

    // Execute
    exercise, err := service.CreateExercise(ctx, CreateExerciseRequest{
        Name: "Bench Press",
    })

    // Assert
    assert.NoError(t, err)
    assert.Equal(t, "Bench Press", exercise.Name)
}
```

## Security Considerations

### Authentication Flow

1. User logs in via Supabase Auth
2. Receives JWT token
3. Include token in Authorization header: `Bearer <token>`
4. Middleware validates token with Supabase
5. Extract user_id from token
6. Use user_id for authorization

### Authorization Rules

- Users can only modify their own resources
- Public exercises are read-only for non-owners
- Private exercises only visible to owner
- Workouts always private to owner

### Data Validation

- All input validated at handler level
- Use struct tags for automatic validation
- Custom validators for business rules

```go
type CreateExerciseRequest struct {
    Name        string `json:"name" binding:"required,min=3,max=100"`
    Description string `json:"description" binding:"max=500"`
    IsPublic    bool   `json:"is_public"`
}
```

## Performance Optimization

### Database
- Use connection pooling (pgxpool)
- Add indexes on frequently queried columns
- Use prepared statements
- Batch operations when possible

### Caching (Future)
- Cache public exercises
- Cache user sessions
- Use Redis for distributed caching

### Pagination
- Implement cursor-based pagination
- Limit default page size
- Index pagination columns

## Logging

### Structured Logging (Future)

```go
log.WithFields(log.Fields{
    "user_id": userID,
    "exercise_id": exerciseID,
    "action": "create",
}).Info("Exercise created")
```

## API Versioning

Future consideration:
- `/api/v1/exercises`
- `/api/v2/exercises`

Currently using unversioned `/api/` routes for MVP.
