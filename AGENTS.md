# Agent Configuration

This file provides context and instructions for AI agents working on this project.

## Project Overview
Mobile Order App Prototype. The core value proposition is investigating a "Single-User Order" experience.
The project is split into 3 binaries: `cmd/admin`, `cmd/restaurant`, `cmd/guest`.

## Tech Stack
- **Go**: Echo v4, GORM, SQLite.
- **Architecture**: Shared `internal/` (DB, Models, Email, Auth), Separate `cmd/` and `views/`.

## Critical Business Logic
> [!IMPORTANT]
> **Exclusive Guest Access**:
> The `Restaurant` model contains `CurrentGuestSession` and `GuestSessionExpires`.
> - **CheckIn**: MUST check if lock is free or expired. If free, generate session and lock.
> - **Menu**: MUST validate `session_id` against the lock.
> - **Checkout**: MUST clear the lock.
> - **Logic Location**: `internal/handlers/guest/handler.go`

## Development Features
- **Email Trap**: Use `email.LogSender` by default. Do not send real emails in dev.
- **Auto Browser**: `utils.OpenBrowser` is called on startup.

## Documentation Rules
- **Walkthroughs**: Must be written in **Japanese** (日本語).
- **Code Comments**: English is preferred, but Japanese is acceptable for complex logic.

## Commands
```bash
go run cmd/guest/main.go
go run cmd/admin/main.go
go run cmd/restaurant/main.go
```
