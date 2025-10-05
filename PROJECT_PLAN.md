# Fitness API - Project Plan

## Project Overview
A REST API built with Gin and Supabase for fitness tracking, allowing users to manage exercises, workouts, and equipment with image uploads.

## Tech Stack
- **Backend**: Go with Gin framework
- **Database**: Supabase (PostgreSQL)
- **Auth**: Supabase built-in authentication
- **Storage**: Supabase Storage for images
- **Migration**: golang-migrate with SQL files
- **Supabase Client**: supabase-go

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

## Epic 3: Authentication & User Management
- [ ] Implement Supabase auth middleware
- [ ] Create user registration endpoint
- [ ] Create user login endpoint
- [ ] Implement JWT token validation
- [ ] Create user context extraction from token
- [ ] Add protected route middleware

## Epic 4: Equipment Management
- [ ] Create equipment model/struct
- [ ] Implement POST /api/equipment (create equipment)
- [ ] Implement GET /api/equipment (list all equipment)
- [ ] Implement GET /api/equipment/:id (get single equipment)
- [ ] Implement PUT /api/equipment/:id (update equipment)
- [ ] Implement DELETE /api/equipment/:id (delete equipment)

## Epic 5: Exercise Management
- [ ] Create exercise model/struct with public/private flag
- [ ] Implement POST /api/exercises (create exercise)
- [ ] Implement GET /api/exercises (list exercises - public + user's private)
- [ ] Implement GET /api/exercises/:id (get single exercise)
- [ ] Implement PUT /api/exercises/:id (update exercise - only if owner)
- [ ] Implement DELETE /api/exercises/:id (delete exercise - only if owner)
- [ ] Implement GET /api/exercises?equipment_id=X (filter by equipment)
- [ ] Add exercise-equipment relationship endpoints

## Epic 6: Workout Management
- [ ] Create workout model/struct
- [ ] Implement POST /api/workouts (create workout with exercises)
- [ ] Implement GET /api/workouts (list user's workouts)
- [ ] Implement GET /api/workouts/:id (get workout with exercises)
- [ ] Implement PUT /api/workouts/:id (update workout)
- [ ] Implement DELETE /api/workouts/:id (delete workout)
- [ ] Add/remove exercises to/from workout

## Epic 7: Image Upload & Storage
- [ ] Configure Supabase Storage buckets (exercises, workouts)
- [ ] Implement POST /api/exercises/:id/image (upload exercise image)
- [ ] Implement POST /api/workouts/:id/image (upload workout image)
- [ ] Implement DELETE endpoints for images
- [ ] Add image URL to exercise/workout responses

## Epic 8: API Documentation & Testing
- [ ] Add input validation for all endpoints
- [ ] Implement proper error handling
- [ ] Add API documentation (comments/Swagger)
- [ ] Create sample requests/responses
- [ ] Write unit tests for core functionality
- [ ] Write integration tests

## Epic 9: Optimization & Polish
- [ ] Add pagination for list endpoints
- [ ] Implement search functionality for exercises
- [ ] Add rate limiting
- [ ] Optimize database queries
- [ ] Add logging middleware
- [ ] Performance testing

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
