# Fitness API - Project Plan

## Project Overview
A REST API built with Gin and Supabase for fitness tracking, allowing users to manage exercises, workouts, and equipment with image uploads.

## Tech Stack
- **Backend**: Go with Gin framework
- **Database**: Supabase (PostgreSQL) with pgx driver
- **Auth**: Supabase Auth with JWT validation
- **Storage**: Supabase Storage for images
- **Migration**: golang-migrate with SQL files
- **Architecture**: Clean Architecture (Handlers → Services → Repositories)
- **Testing**: Go native testing with mock repositories, Newman for API tests
- **API Testing**: Postman collections with Newman CLI

---

## Epic 1: Project Setup & Configuration
- [x] Initialize Go module and project structure
- [x] Install dependencies (Gin, Supabase client, golang-migrate)
- [x] Set up Supabase project and get credentials
- [x] Create .env configuration file
- [x] Set up database connection
- [x] Configure Supabase Storage buckets for images

## Epic 2: Database Schema & Migrations
- [x] Design database schema (users, exercises, workouts, equipment, relationships)
- [x] Create migration files for all tables
- [x] Set up many-to-many relationship tables (workout_exercises, exercise_equipment)
- [x] Add indexes for performance (user_id, public flag, etc.)
- [x] Run initial migrations

## Epic 3: Authentication & User Management ✅
- [x] Implement JWT token validation middleware
- [x] Create user context extraction from token
- [x] Add protected route middleware to /api
- [x] Create test endpoint /api/me
- [x] Add SKIP_AUTH environment variable for development
- [x] Create Postman collection with Newman support
- [x] Implement smart token caching for tests

## Epic 4: Equipment Management (Clean Architecture)
**Architecture Setup:**
- [ ] Create models/equipment.go (domain model)
- [ ] Create repositories/equipment.go (interface + PostgreSQL implementation)
- [ ] Create services/equipment.go (business logic layer)
- [ ] Create handlers/equipment.go (HTTP/Gin handlers)
- [ ] Add dependency injection in main.go

**CRUD Endpoints:**
- [ ] POST /api/equipment - Create equipment
- [ ] GET /api/equipment - List all equipment
- [ ] GET /api/equipment/:id - Get single equipment
- [ ] PUT /api/equipment/:id - Update equipment
- [ ] DELETE /api/equipment/:id - Delete equipment

**Testing:**
- [ ] Create mock repository for unit tests
- [ ] Write unit tests for service layer (business logic)
- [ ] Write unit tests for handlers (HTTP layer)
- [ ] Add equipment endpoints to Postman collection
- [ ] Test with Newman

**Validation & Error Handling:**
- [ ] Add input validation (required fields, max lengths)
- [ ] Implement proper error responses (400, 404, 500)
- [ ] Add ownership validation (user can only modify their equipment)

## Epic 5: Exercise Management (Clean Architecture)
**Architecture Setup:**
- [ ] Create models/exercise.go with public/private flag
- [ ] Create repositories/exercise.go (interface + implementation)
- [ ] Create services/exercise.go (with visibility logic)
- [ ] Create handlers/exercise.go
- [ ] Add dependency injection in main.go

**CRUD Endpoints:**
- [ ] POST /api/exercises - Create exercise
- [ ] GET /api/exercises - List exercises (public + user's private)
- [ ] GET /api/exercises/:id - Get single exercise
- [ ] PUT /api/exercises/:id - Update exercise (ownership check)
- [ ] DELETE /api/exercises/:id - Delete exercise (ownership check)

**Advanced Features:**
- [ ] GET /api/exercises?equipment_id=X - Filter by equipment
- [ ] GET /api/exercises?is_public=true - Filter public exercises
- [ ] POST /api/exercises/:id/equipment - Link equipment to exercise
- [ ] DELETE /api/exercises/:id/equipment/:equipment_id - Unlink equipment

**Testing:**
- [ ] Create mock repository
- [ ] Unit tests for service layer (visibility rules, ownership)
- [ ] Unit tests for handlers
- [ ] Add to Postman collection
- [ ] Test with Newman

**Validation:**
- [ ] Input validation
- [ ] Ownership validation
- [ ] Public/private visibility rules

## Epic 6: Workout Management (Clean Architecture)
**Architecture Setup:**
- [ ] Create models/workout.go and models/workout_exercise.go
- [ ] Create repositories/workout.go (interface + implementation)
- [ ] Create services/workout.go (with nested exercise logic)
- [ ] Create handlers/workout.go
- [ ] Add dependency injection in main.go

**CRUD Endpoints:**
- [ ] POST /api/workouts - Create workout with exercises
- [ ] GET /api/workouts - List user's workouts
- [ ] GET /api/workouts/:id - Get workout with exercises
- [ ] PUT /api/workouts/:id - Update workout
- [ ] DELETE /api/workouts/:id - Delete workout (cascade exercises)

**Workout-Exercise Relationship:**
- [ ] POST /api/workouts/:id/exercises - Add exercise to workout
- [ ] PUT /api/workouts/:id/exercises/:exercise_id - Update sets/reps/order
- [ ] DELETE /api/workouts/:id/exercises/:exercise_id - Remove exercise

**Testing:**
- [ ] Create mock repositories
- [ ] Unit tests for service layer (nested operations)
- [ ] Unit tests for handlers
- [ ] Add to Postman collection
- [ ] Test with Newman

**Validation:**
- [ ] Input validation (sets, reps, order)
- [ ] Ownership validation
- [ ] Exercise existence validation

## Epic 7: Image Upload & Storage (Clean Architecture)
**Architecture Setup:**
- [ ] Create storage/interface.go (StorageService interface)
- [ ] Create storage/supabase.go (real Supabase Storage implementation)
- [ ] Create storage/mock.go (mock for unit tests)
- [ ] Add USE_MOCK_STORAGE environment variable
- [ ] Update services to use StorageService interface

**Storage Configuration:**
- [ ] Configure Supabase Storage buckets (exercises, workouts)
- [ ] Set up bucket policies (public read, authenticated write)

**Upload Endpoints:**
- [ ] POST /api/exercises/:id/image - Upload exercise image
- [ ] POST /api/workouts/:id/image - Upload workout image
- [ ] DELETE /api/exercises/:id/image - Delete exercise image
- [ ] DELETE /api/workouts/:id/image - Delete workout image

**Testing:**
- [ ] Unit tests with mock storage
- [ ] Test file validation (size, type)
- [ ] Add to Postman collection (multipart/form-data)

**Validation:**
- [ ] File type validation (jpg, png, webp)
- [ ] File size limits
- [ ] Ownership validation

## Epic 8: Testing & Documentation
**Unit Testing:**
- [ ] Ensure all services have unit tests (using mocks)
- [ ] Ensure all handlers have unit tests
- [ ] Add table-driven tests for edge cases
- [ ] Achieve >80% code coverage

**Integration Testing:**
- [ ] Set up test database (testcontainers or separate DB)
- [ ] Write integration tests for critical flows
- [ ] Test authentication flow end-to-end

**API Documentation:**
- [ ] Add godoc comments to all exported functions
- [ ] Create OpenAPI/Swagger specification
- [ ] Document error responses
- [ ] Add usage examples

**Postman/Newman:**
- [ ] Complete Postman collection for all endpoints
- [ ] Add test assertions in Postman
- [ ] Create Newman CI/CD script

## Epic 9: Optimization & Polish
**Performance:**
- [ ] Add pagination to list endpoints (limit, offset)
- [ ] Implement search for exercises (by name, description)
- [ ] Add database indexes (already done in migrations)
- [ ] Add query optimization (eager loading, N+1 prevention)

**Observability:**
- [ ] Add structured logging (zerolog or zap)
- [ ] Add request ID middleware
- [ ] Add response time logging

**Security & Reliability:**
- [ ] Add rate limiting middleware
- [ ] Add request timeout middleware
- [ ] Add CORS configuration
- [ ] Add graceful shutdown
- [ ] Add health check with DB connectivity

**Developer Experience:**
- [ ] Add Makefile with common commands
- [ ] Add Docker Compose for local development
- [ ] Update documentation with deployment guide

---

## Database Schema (Preliminary)

### Tables
1. **users** (managed by Supabase Auth)
2. **equipment** (id, name, description, user_id, created_at, updated_at)
3. **exercises** (id, name, description, is_public, user_id, image_url, created_at, updated_at)
4. **workouts** (id, name, description, user_id, image_url, created_at, updated_at)
5. **exercise_equipment** (exercise_id, equipment_id)
6. **workout_exercises** (workout_id, exercise_id, sets, reps, order)

---

## API Endpoints Summary

### Auth
- POST /api/auth/register
- POST /api/auth/login

### Equipment
- GET /api/equipment
- POST /api/equipment
- GET /api/equipment/:id
- PUT /api/equipment/:id
- DELETE /api/equipment/:id

### Exercises
- GET /api/exercises (with filters: ?equipment_id=X, ?is_public=true)
- POST /api/exercises
- GET /api/exercises/:id
- PUT /api/exercises/:id
- DELETE /api/exercises/:id
- POST /api/exercises/:id/image
- DELETE /api/exercises/:id/image

### Workouts
- GET /api/workouts
- POST /api/workouts
- GET /api/workouts/:id
- PUT /api/workouts/:id
- DELETE /api/workouts/:id
- POST /api/workouts/:id/image
- DELETE /api/workouts/:id/image
- POST /api/workouts/:id/exercises
- DELETE /api/workouts/:id/exercises/:exercise_id

---

## Next Steps
1. Start with Epic 1: Project Setup & Configuration
2. Move to Epic 2: Database Schema & Migrations
3. Progress through epics sequentially
