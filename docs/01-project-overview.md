# Project Overview

## Fitness API

A comprehensive REST API for fitness tracking and workout management built with Go, Gin framework, and Supabase.

## Purpose

Enable users to:
- Create and manage exercise libraries (public and private)
- Design workout templates with detailed exercise parameters
- Track actual workout sessions with performance metrics
- Manage equipment and filter exercises by equipment
- Upload images for exercises and workouts
- Monitor progress with personal records and historical data

## Key Features

### Exercise Management
- Public exercise library (shared with all users)
- Private exercises (user-specific)
- Multi-equipment support per exercise
- Image uploads for exercises

### Workout Templates
- Create reusable workout plans
- Configure sets, reps, weight, tempo, rest times
- Support for supersets, dropsets, warmups, and cooldowns
- Order exercises in specific sequences
- Track intensity and RPE (Rate of Perceived Exertion) targets

### Workout Sessions
- Start sessions from templates or create ad-hoc
- Real-time workout tracking
- Log actual performance vs planned
- Track mood, energy levels, heart rate
- Session ratings and notes

### Performance Tracking
- Individual exercise logs with actual performance
- Personal record (PR) detection
- Form rating and RPE tracking
- Historical comparison with previous performances

### Equipment Management
- User-defined equipment library
- Link exercises to equipment
- Filter exercises by available equipment

## Tech Stack

- **Language**: Go 1.25+
- **Web Framework**: Gin
- **Database**: PostgreSQL (via Supabase)
- **Authentication**: Supabase Auth (JWT-based)
- **Storage**: Supabase Storage (for images)
- **Database Driver**: pgx/v5
- **Migrations**: golang-migrate

## Project Status

Currently in development. See [PROJECT_PLAN.md](../PROJECT_PLAN.md) for progress tracking.

### Completed
- ✅ Epic 1: Project Setup & Configuration

### In Progress
- 🔄 Epic 2: Database Schema & Migrations

### Upcoming
- Epic 3: Authentication & User Management
- Epic 4: Equipment Management
- Epic 5: Exercise Management
- Epic 6: Workout Management
- Epic 7: Image Upload & Storage
- Epic 8: API Documentation & Testing
- Epic 9: Optimization & Polish
