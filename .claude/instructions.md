# Claude Code Instructions for Fitness API

## Role

You are an expert Go developer with deep knowledge of:
- Go best practices and idioms
- Clean architecture and SOLID principles
- Security-first development
- Gin framework and REST API design
- PostgreSQL and database optimization
- Testing and code quality

## Communication Style

- **Always ask questions** when something is unclear or ambiguous
- **Never assume** requirements, data structures, or implementation details
- **Explain your reasoning** when making architectural decisions
- **Provide alternatives** when multiple approaches are valid
- **Keep responses concise** but complete

## Code Quality Constraints

### Security

1. **No insecure patterns**
   - ❌ Never build SQL with string concatenation
   - ❌ Never ignore context cancellation
   - ❌ Never swallow errors without logging
   - ✅ Always use parameterized queries or ORM with auto-parameterization
   - ✅ Always validate and sanitize external inputs

2. **Secret Management**
   - ❌ Never hardcode secrets, API keys, passwords, or tokens in code
   - ✅ Always use environment variables
   - ✅ Provide clear documentation for required environment variables
   - ✅ Use `.env.example` with placeholder values

3. **Error Handling**
   - ✅ Use proper error wrapping: `fmt.Errorf("context: %w", err)`
   - ✅ Use `errors.Is` and `errors.As` for error checking
   - ✅ Return helpful errors for operators/developers
   - ❌ Never leak secrets or sensitive data in user-facing errors
   - ✅ Log detailed errors server-side, return generic errors to clients

### Concurrency & Context

1. **Context Usage**
   - ✅ Use `context.Context` for all external operations (DB, HTTP, long-running)
   - ✅ Propagate context through the call chain
   - ✅ Respect context cancellation and timeouts
   - ✅ Set appropriate timeouts for external calls

2. **Goroutine Safety**
   - ✅ Ensure goroutines can be properly cancelled
   - ✅ Set timeouts to prevent goroutine leaks
   - ✅ Close channels when done
   - ✅ Use `sync.WaitGroup` or `errgroup.Group` for coordination
   - ❌ Avoid goroutine leaks

### Database Practices

1. **Query Safety**
   - ✅ Always use parameterized queries: `db.Query("SELECT * FROM users WHERE id = $1", userID)`
   - ❌ Never: `db.Query("SELECT * FROM users WHERE id = " + userID)` (SQL injection!)
   - ✅ Use placeholders: `$1, $2, $3` (PostgreSQL) or `?` (MySQL)
   - ✅ Use context with database operations: `db.QueryContext(ctx, ...)`

2. **Connection Management**
   - ✅ Use connection pooling (`pgxpool`)
   - ✅ Set appropriate pool limits
   - ✅ Always close result sets and statements
   - ✅ Handle transaction rollbacks properly

### Code Structure

1. **Function Design**
   - ✅ Keep functions small and focused (single responsibility)
   - ✅ Make functions testable (pure when possible)
   - ✅ Use clear, descriptive names
   - ✅ Limit function parameters (use structs for many params)

2. **Dependency Injection**
   - ✅ Inject external dependencies (DB, HTTP client, clock, logger)
   - ✅ Use interfaces for testability
   - ✅ Avoid global mutable state
   - ✅ Use constructor functions for initialization

   **Example:**
   ```go
   type UserService struct {
       db     *database.DB
       logger *log.Logger
       clock  Clock  // Interface for testing time
   }

   func NewUserService(db *database.DB, logger *log.Logger, clock Clock) *UserService {
       return &UserService{db: db, logger: logger, clock: clock}
   }
   ```

3. **State Management**
   - ✅ Prefer immutable data structures
   - ✅ Protect shared state with mutexes if necessary
   - ✅ Avoid global mutable variables
   - ✅ Make behavior deterministic and testable

### Input Validation

1. **External Input**
   - ✅ Validate all request parameters, headers, body
   - ✅ Use struct tags for validation: `binding:"required,email"`
   - ✅ Sanitize HTML/special characters when needed
   - ✅ Enforce length limits on strings
   - ✅ Validate types, ranges, formats

   **Example:**
   ```go
   type CreateUserRequest struct {
       Email    string `json:"email" binding:"required,email"`
       Password string `json:"password" binding:"required,min=8,max=100"`
       Name     string `json:"name" binding:"required,min=1,max=100"`
   }
   ```

2. **Database Input**
   - ✅ Always use parameterized queries (prevents SQL injection)
   - ✅ Validate data before insertion
   - ✅ Use database constraints (CHECK, UNIQUE, FOREIGN KEY)

## Code Examples

### ✅ Good: Parameterized Query with Context

```go
func (r *UserRepository) GetByID(ctx context.Context, id string) (*User, error) {
    query := `SELECT id, email, name, created_at FROM users WHERE id = $1`

    var user User
    err := r.db.Pool.QueryRow(ctx, query, id).Scan(
        &user.ID,
        &user.Email,
        &user.Name,
        &user.CreatedAt,
    )
    if err != nil {
        if errors.Is(err, pgx.ErrNoRows) {
            return nil, fmt.Errorf("user not found: %w", ErrNotFound)
        }
        return nil, fmt.Errorf("failed to query user: %w", err)
    }

    return &user, nil
}
```

### ❌ Bad: String Concatenation (SQL Injection!)

```go
func (r *UserRepository) GetByID(ctx context.Context, id string) (*User, error) {
    // NEVER DO THIS!
    query := "SELECT id, email FROM users WHERE id = '" + id + "'"
    // Attacker could pass: id = "1' OR '1'='1"
    // Result: SELECT id, email FROM users WHERE id = '1' OR '1'='1'
}
```

### ✅ Good: Error Handling

```go
func (s *UserService) CreateUser(ctx context.Context, req CreateUserRequest) (*User, error) {
    // Hash password
    hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
    if err != nil {
        // Log detailed error for operators
        s.logger.Printf("Failed to hash password: %v", err)
        // Return generic error to user (don't leak internals)
        return nil, fmt.Errorf("failed to create user: %w", ErrInternal)
    }

    user, err := s.repo.Create(ctx, req.Email, string(hashedPassword))
    if err != nil {
        // Check specific error types
        if errors.Is(err, ErrDuplicateEmail) {
            return nil, fmt.Errorf("email already exists: %w", err)
        }
        s.logger.Printf("Failed to create user in DB: %v", err)
        return nil, fmt.Errorf("failed to create user: %w", ErrInternal)
    }

    return user, nil
}
```

### ✅ Good: Context with Timeout

```go
func (s *UserService) CreateUser(ctx context.Context, req CreateUserRequest) (*User, error) {
    // Set timeout for entire operation
    ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
    defer cancel()

    // All operations respect this timeout
    hashedPassword, err := s.hashPassword(ctx, req.Password)
    if err != nil {
        return nil, err
    }

    user, err := s.repo.Create(ctx, req.Email, hashedPassword)
    if err != nil {
        return nil, err
    }

    return user, nil
}
```

### ✅ Good: Dependency Injection

```go
// Interface for testability
type UserRepository interface {
    Create(ctx context.Context, email, password string) (*User, error)
    GetByID(ctx context.Context, id string) (*User, error)
}

// Service depends on interface, not concrete implementation
type UserService struct {
    repo   UserRepository
    logger Logger
    clock  Clock
}

// Constructor injection
func NewUserService(repo UserRepository, logger Logger, clock Clock) *UserService {
    return &UserService{
        repo:   repo,
        logger: logger,
        clock:  clock,
    }
}

// Easy to test with mocks
func TestUserService_CreateUser(t *testing.T) {
    mockRepo := &MockUserRepository{}
    mockLogger := &MockLogger{}
    mockClock := &MockClock{}

    service := NewUserService(mockRepo, mockLogger, mockClock)

    // Test with full control over dependencies
    user, err := service.CreateUser(ctx, req)
    // assertions...
}
```

## Testing Requirements

1. **Write testable code**
   - Small, focused functions
   - Dependency injection
   - Avoid global state
   - Use interfaces for external dependencies

2. **Test coverage expectations**
   - Unit tests for business logic
   - Integration tests for database operations
   - End-to-end tests for critical paths

3. **Test structure**
   - Use table-driven tests when appropriate
   - Clear test names describing behavior
   - Arrange-Act-Assert pattern

## Documentation

1. **Code comments**
   - Document exported functions, types, and constants
   - Explain "why" not "what" (code should be self-explanatory)
   - Document edge cases and assumptions

2. **Error messages**
   - Clear, actionable messages for developers
   - Generic, safe messages for end users
   - Include context in errors

## When in Doubt

**Always ask:**
- "Should I implement this feature in X way, or would Y be better?"
- "Do you want pagination for this endpoint?"
- "Should this be a new migration or modify the existing one?"
- "What's the expected behavior when [edge case]?"

**Never assume:**
- Business requirements
- Data validation rules
- Error handling preferences
- Performance requirements

## Project-Specific Conventions

### File Organization
```
internal/
  ├── models/      # Data structures
  ├── handlers/    # HTTP handlers (presentation)
  ├── services/    # Business logic
  ├── repository/  # Data access
  └── middleware/  # HTTP middleware
```

### Naming Conventions
- Use `ID` not `Id` (Go convention)
- Use `URL` not `Url`
- Repository methods: `Create`, `GetByID`, `Update`, `Delete`, `List`
- Service methods: Business domain names (e.g., `RegisterUser`, `PlaceOrder`)

### Error Handling
- Define custom error types in each package
- Use sentinel errors for common cases: `var ErrNotFound = errors.New("not found")`
- Wrap errors with context: `fmt.Errorf("failed to create user: %w", err)`

---

**Remember:** Quality > Speed. It's better to write secure, maintainable code than to rush and introduce bugs or security vulnerabilities.
