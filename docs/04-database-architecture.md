# Database Architecture

## Overview

The database is designed to support a comprehensive fitness tracking system with workout templates, session tracking, and performance logging.

## Database Schema

### Entity Relationship Diagram

```
users (Supabase Auth)
    ↓ creates
equipment
    ↓ used in
exercises ←→ exercise_equipment (many-to-many)
    ↓ included in
workouts ←→ workout_exercises (many-to-many with details)
    ↓ executed as
workout_sessions
    ↓ contains
exercise_logs (performance data)
```

## Core Tables

### 1. Users (Managed by Supabase Auth)

Supabase automatically manages the `auth.users` table with:
- `id` (UUID) - Primary key
- `email` - User email
- `encrypted_password` - Hashed password
- `created_at`, `updated_at` - Timestamps

**Note**: We reference `auth.users(id)` in our custom tables using `user_id`.

### 2. Equipment

Stores gym equipment definitions (dumbbells, barbells, resistance bands, etc.).

```sql
CREATE TABLE equipment (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL REFERENCES auth.users(id) ON DELETE CASCADE,
    name TEXT NOT NULL,
    description TEXT,
    created_at TIMESTAMPTZ DEFAULT NOW(),
    updated_at TIMESTAMPTZ DEFAULT NOW()
);
```

**Fields**:
- `id` - Unique identifier
- `user_id` - Owner of the equipment (custom equipment per user)
- `name` - Equipment name (e.g., "Barbell", "Dumbbells")
- `description` - Optional details
- `created_at`, `updated_at` - Timestamps

**Indexes**:
- `user_id` - Fast lookup of user's equipment

### 3. Exercises

Exercise library with public/private visibility.

```sql
CREATE TABLE exercises (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL REFERENCES auth.users(id) ON DELETE CASCADE,
    name TEXT NOT NULL,
    description TEXT,
    is_public BOOLEAN DEFAULT FALSE,
    image_url TEXT,
    created_at TIMESTAMPTZ DEFAULT NOW(),
    updated_at TIMESTAMPTZ DEFAULT NOW()
);
```

**Fields**:
- `id` - Unique identifier
- `user_id` - Creator of the exercise
- `name` - Exercise name (e.g., "Bench Press")
- `description` - Instructions, form tips
- `is_public` - TRUE = visible to all users, FALSE = private
- `image_url` - Supabase Storage URL for exercise image
- `created_at`, `updated_at` - Timestamps

**Indexes**:
- `user_id` - User's exercises
- `is_public` - Public exercises lookup
- Composite: `(is_public, user_id)` - Filter public + user's private

### 4. Exercise Equipment (Junction Table)

Links exercises to equipment (many-to-many relationship).

```sql
CREATE TABLE exercise_equipment (
    exercise_id UUID NOT NULL REFERENCES exercises(id) ON DELETE CASCADE,
    equipment_id UUID NOT NULL REFERENCES equipment(id) ON DELETE CASCADE,
    PRIMARY KEY (exercise_id, equipment_id)
);
```

**Purpose**: One exercise can use multiple equipment, and one equipment can be used in multiple exercises.

**Example**:
- Bench Press → uses Barbell + Bench
- Dumbbell Curl → uses Dumbbells

### 5. Workouts (Templates)

Workout plans/templates created by users.

```sql
CREATE TABLE workouts (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL REFERENCES auth.users(id) ON DELETE CASCADE,
    name TEXT NOT NULL,
    description TEXT,
    image_url TEXT,
    created_at TIMESTAMPTZ DEFAULT NOW(),
    updated_at TIMESTAMPTZ DEFAULT NOW()
);
```

**Fields**:
- `id` - Unique identifier
- `user_id` - Owner of the workout
- `name` - Workout name (e.g., "Push Day A")
- `description` - Workout notes, goals
- `image_url` - Supabase Storage URL
- `created_at`, `updated_at` - Timestamps

**Note**: Workouts are always private to the user.

### 6. Workout Exercises (Junction Table with Details)

Links workouts to exercises with detailed parameters.

```sql
CREATE TABLE workout_exercises (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    workout_id UUID NOT NULL REFERENCES workouts(id) ON DELETE CASCADE,
    exercise_id UUID NOT NULL REFERENCES exercises(id) ON DELETE CASCADE,
    order_index INTEGER NOT NULL,
    sets INTEGER DEFAULT 1,
    reps INTEGER,
    weight_kg REAL,
    duration_seconds INTEGER,
    distance_meters REAL,
    rest_time_seconds INTEGER DEFAULT 60,
    intensity_percentage REAL,
    tempo TEXT,
    notes TEXT,
    is_superset BOOLEAN DEFAULT FALSE,
    superset_group_id UUID,
    is_dropset BOOLEAN DEFAULT FALSE,
    is_warmup BOOLEAN DEFAULT FALSE,
    is_cooldown BOOLEAN DEFAULT FALSE,
    target_rpe INTEGER CHECK (target_rpe BETWEEN 1 AND 10),
    created_at TIMESTAMPTZ DEFAULT NOW(),
    updated_at TIMESTAMPTZ DEFAULT NOW()
);
```

**Fields**:
- `id` - Unique identifier (allows multiple instances of same exercise)
- `workout_id` - Parent workout
- `exercise_id` - Exercise being performed
- `order_index` - Order in workout (1, 2, 3...)
- `sets` - Target number of sets
- `reps` - Target repetitions
- `weight_kg` - Target weight
- `duration_seconds` - For timed exercises (planks, running)
- `distance_meters` - For distance exercises (running, rowing)
- `rest_time_seconds` - Rest between sets
- `intensity_percentage` - % of 1RM (one-rep max)
- `tempo` - Lifting tempo (e.g., "3-1-2-0")
- `notes` - Exercise-specific notes
- `is_superset` - Part of a superset
- `superset_group_id` - Groups exercises performed back-to-back
- `is_dropset` - Dropset indicator
- `is_warmup` / `is_cooldown` - Warmup/cooldown flags
- `target_rpe` - Rate of Perceived Exertion (1-10)

**Indexes**:
- `workout_id` - Get all exercises in a workout
- `(workout_id, order_index)` - Ordered exercise list

### 7. Workout Sessions (Actual Workouts)

Records of actual workout performances.

```sql
CREATE TABLE workout_sessions (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL REFERENCES auth.users(id) ON DELETE CASCADE,
    workout_id UUID REFERENCES workouts(id) ON DELETE SET NULL,
    name TEXT,
    started_at TIMESTAMPTZ NOT NULL,
    completed_at TIMESTAMPTZ,
    duration_minutes INTEGER,
    status TEXT CHECK (status IN ('planned', 'in_progress', 'completed', 'cancelled', 'paused')),
    location TEXT,
    weather_conditions TEXT,
    energy_level_start INTEGER CHECK (energy_level_start BETWEEN 1 AND 10),
    energy_level_end INTEGER CHECK (energy_level_end BETWEEN 1 AND 10),
    perceived_exertion INTEGER CHECK (perceived_exertion BETWEEN 1 AND 10),
    mood_before TEXT,
    mood_after TEXT,
    calories_burned INTEGER,
    heart_rate_avg INTEGER,
    heart_rate_max INTEGER,
    notes TEXT,
    workout_rating INTEGER CHECK (workout_rating BETWEEN 1 AND 5),
    created_at TIMESTAMPTZ DEFAULT NOW(),
    updated_at TIMESTAMPTZ DEFAULT NOW()
);
```

**Fields**:
- `id` - Session identifier
- `user_id` - User performing the workout
- `workout_id` - Template used (NULL if ad-hoc workout)
- `name` - Custom session name
- `started_at` - When workout started
- `completed_at` - When finished (NULL if in progress)
- `duration_minutes` - Total workout duration
- `status` - Current state (planned/in_progress/completed/cancelled/paused)
- `location` - Where workout happened
- `weather_conditions` - Environmental factors
- `energy_level_start/end` - Energy before/after (1-10)
- `perceived_exertion` - Overall RPE (1-10)
- `mood_before/after` - Subjective mood
- `calories_burned` - Estimated calories
- `heart_rate_avg/max` - Heart rate metrics
- `notes` - Session notes
- `workout_rating` - How good was the workout (1-5 stars)

**Indexes**:
- `user_id` - User's workout history
- `(user_id, started_at)` - Chronological history
- `status` - Active workouts

### 8. Exercise Logs (Performance Data)

Individual exercise performances within a session.

```sql
CREATE TABLE exercise_logs (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    workout_session_id UUID NOT NULL REFERENCES workout_sessions(id) ON DELETE CASCADE,
    exercise_id UUID NOT NULL REFERENCES exercises(id) ON DELETE CASCADE,
    workout_exercise_id UUID REFERENCES workout_exercises(id) ON DELETE SET NULL,
    order_index INTEGER NOT NULL,
    sets_completed INTEGER DEFAULT 0,
    sets_planned INTEGER DEFAULT 1,
    reps_completed INTEGER,
    reps_planned INTEGER,
    weight_kg REAL,
    duration_seconds INTEGER,
    distance_meters REAL,
    rest_time_seconds INTEGER,
    intensity_percentage REAL,
    rpe INTEGER CHECK (rpe BETWEEN 1 AND 10),
    form_rating INTEGER CHECK (form_rating BETWEEN 1 AND 5),
    equipment_used TEXT,
    notes TEXT,
    is_personal_record BOOLEAN DEFAULT FALSE,
    previous_best_weight REAL,
    previous_best_reps INTEGER,
    previous_best_duration INTEGER,
    created_at TIMESTAMPTZ DEFAULT NOW(),
    updated_at TIMESTAMPTZ DEFAULT NOW()
);
```

**Fields**:
- `id` - Log identifier
- `workout_session_id` - Parent session
- `exercise_id` - Exercise performed
- `workout_exercise_id` - Reference to template (if from template)
- `order_index` - Order in session
- `sets_completed/planned` - Actual vs planned sets
- `reps_completed/planned` - Actual vs planned reps
- `weight_kg` - Actual weight used
- `duration_seconds` - Actual duration
- `distance_meters` - Actual distance
- `rest_time_seconds` - Actual rest
- `intensity_percentage` - Actual intensity
- `rpe` - Actual Rate of Perceived Exertion (1-10)
- `form_rating` - How good was form (1-5)
- `equipment_used` - JSON array of equipment IDs actually used
- `notes` - Exercise notes
- `is_personal_record` - PR flag
- `previous_best_*` - Previous records for comparison

**Indexes**:
- `workout_session_id` - Session's exercises
- `exercise_id` - Exercise history
- `(exercise_id, is_personal_record)` - Find PRs
- `(user_id, exercise_id, created_at)` - User's exercise progression (via join)

## Relationships Summary

### One-to-Many
- `users` → `equipment` (one user has many equipment)
- `users` → `exercises` (one user creates many exercises)
- `users` → `workouts` (one user has many workouts)
- `users` → `workout_sessions` (one user has many sessions)
- `workouts` → `workout_exercises` (one workout has many exercises)
- `workout_sessions` → `exercise_logs` (one session has many logs)

### Many-to-Many
- `exercises` ↔ `equipment` (via `exercise_equipment`)
- `workouts` ↔ `exercises` (via `workout_exercises`)

## Data Types

### UUIDs vs Auto-increment IDs
We use UUIDs (`UUID`) for:
- Distributed system support
- No collision risk
- Security (non-sequential)
- Easy merging/syncing

### Timestamps
- `TIMESTAMPTZ` - Timezone-aware timestamps
- `created_at` - Record creation time
- `updated_at` - Last modification time (updated via trigger)

### Constraints
- `CHECK` constraints for valid ranges (RPE 1-10, ratings 1-5)
- `NOT NULL` for required fields
- `DEFAULT` values for common cases
- `ON DELETE CASCADE` - Delete child records when parent deleted
- `ON DELETE SET NULL` - Keep child but remove reference

## Indexes

### Performance Indexes

```sql
-- Equipment
CREATE INDEX idx_equipment_user_id ON equipment(user_id);

-- Exercises
CREATE INDEX idx_exercises_user_id ON exercises(user_id);
CREATE INDEX idx_exercises_is_public ON exercises(is_public);
CREATE INDEX idx_exercises_public_user ON exercises(is_public, user_id);

-- Exercise Equipment
CREATE INDEX idx_exercise_equipment_exercise ON exercise_equipment(exercise_id);
CREATE INDEX idx_exercise_equipment_equipment ON exercise_equipment(equipment_id);

-- Workouts
CREATE INDEX idx_workouts_user_id ON workouts(user_id);

-- Workout Exercises
CREATE INDEX idx_workout_exercises_workout ON workout_exercises(workout_id);
CREATE INDEX idx_workout_exercises_order ON workout_exercises(workout_id, order_index);

-- Workout Sessions
CREATE INDEX idx_workout_sessions_user ON workout_sessions(user_id);
CREATE INDEX idx_workout_sessions_user_date ON workout_sessions(user_id, started_at);
CREATE INDEX idx_workout_sessions_status ON workout_sessions(status);

-- Exercise Logs
CREATE INDEX idx_exercise_logs_session ON exercise_logs(workout_session_id);
CREATE INDEX idx_exercise_logs_exercise ON exercise_logs(exercise_id);
CREATE INDEX idx_exercise_logs_pr ON exercise_logs(exercise_id, is_personal_record);
```

## Triggers

### Auto-update `updated_at` timestamps

```sql
CREATE OR REPLACE FUNCTION update_updated_at_column()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = NOW();
    RETURN NEW;
END;
$$ language 'plpgsql';

-- Apply to all tables with updated_at
CREATE TRIGGER update_equipment_updated_at BEFORE UPDATE ON equipment
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER update_exercises_updated_at BEFORE UPDATE ON exercises
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

-- ... repeat for other tables
```

## Migrations Strategy

### Migration Files
- Located in `migrations/` directory
- Numbered sequentially: `001_initial.up.sql`, `001_initial.down.sql`
- `up.sql` - Apply migration
- `down.sql` - Rollback migration

### Migration Order
1. Core tables (equipment, exercises)
2. Junction tables (exercise_equipment)
3. Workout tables (workouts, workout_exercises)
4. Session tables (workout_sessions, exercise_logs)
5. Indexes and triggers
6. Sample data (optional, development only)

## Security Considerations

### Row Level Security (RLS)

Supabase supports PostgreSQL RLS policies:

```sql
-- Enable RLS on tables
ALTER TABLE exercises ENABLE ROW LEVEL SECURITY;

-- Users can only see public exercises or their own
CREATE POLICY exercises_select_policy ON exercises
    FOR SELECT
    USING (is_public = true OR auth.uid() = user_id);

-- Users can only insert their own exercises
CREATE POLICY exercises_insert_policy ON exercises
    FOR INSERT
    WITH CHECK (auth.uid() = user_id);

-- Users can only update their own exercises
CREATE POLICY exercises_update_policy ON exercises
    FOR UPDATE
    USING (auth.uid() = user_id);
```

### Data Isolation
- All user data tied to `user_id`
- `ON DELETE CASCADE` ensures cleanup
- RLS policies enforce access control
- No cross-user data leakage

## Query Patterns

### Get User's Exercises (Public + Private)

```sql
SELECT e.*
FROM exercises e
WHERE e.user_id = $1 OR e.is_public = true
ORDER BY e.created_at DESC;
```

### Get Workout with Exercises

```sql
SELECT
    w.*,
    json_agg(
        json_build_object(
            'exercise', e,
            'details', we
        ) ORDER BY we.order_index
    ) as exercises
FROM workouts w
LEFT JOIN workout_exercises we ON we.workout_id = w.id
LEFT JOIN exercises e ON e.id = we.exercise_id
WHERE w.id = $1
GROUP BY w.id;
```

### Get Exercise Performance History

```sql
SELECT
    el.*,
    ws.started_at,
    ws.completed_at
FROM exercise_logs el
JOIN workout_sessions ws ON ws.id = el.workout_session_id
WHERE el.exercise_id = $1
    AND ws.user_id = $2
    AND ws.status = 'completed'
ORDER BY ws.started_at DESC
LIMIT 10;
```

## Future Enhancements

### Possible Additions
- **Exercise categories** (strength, cardio, flexibility)
- **Muscle groups** (chest, back, legs)
- **Exercise variations** (parent-child relationships)
- **Workout programs** (multi-week plans)
- **Social features** (sharing, following)
- **Achievements/badges**
- **1RM calculator** (estimated one-rep max)
- **Progressive overload tracking**
- **Deload weeks**
- **Body measurements** (weight, body fat %)

### Scaling Considerations
- Partitioning large tables by date (workout_sessions, exercise_logs)
- Archiving old data
- Caching frequently accessed data (public exercises)
- Read replicas for analytics
