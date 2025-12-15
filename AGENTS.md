# Agent Configuration

This file provides context and instructions for AI agents working on this project.

## Project Overview
Mobile Order App Prototype. The core value proposition of this prototype is investigating a "Single-User Order" experience to prevent kitchen congestion.

## Tech Stack
- **Go**: Latest version. Use standard library where possible, but Echo for HTTP.
- **ORM**: GORM with SQLite (`toeic_system.db` -> renamed to `mobile_order.db` conceptually, but currently `toeic_system.db` in code).
- **Frontend**: Single `views/index.html` file serving as a SPA-like interface.

## Critical Business Logic
> [!IMPORTANT]
> **Exclusive Guest Access**:
> The `Restaurant` model contains `CurrentGuestSession` and `GuestSessionExpires`.
> - **CheckIn**: MUST check if lock is free or expired. If free, generate session and lock.
> - **Menu**: MUST validate `session_id` against the lock.
> - **Checkout**: MUST clear the lock to allow next customer.
> **DO NOT REMOVE OR BYPASS THIS LOCKING MECHANISM.**

## Coding Conventions
- **Handlers**: specific logic belongs in `internal/handlers`.
- **Models**: Database interactions managed via GORM models in `internal/models`.
- **Error Handling**: Return JSON errors with appropriate HTTP status codes (423 Locked, 401 Unauthorized).

## Commands
- Run server: `go run main.go`
- Add dependency: `go get <package>`
