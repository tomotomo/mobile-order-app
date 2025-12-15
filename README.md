# Mobile Order App Prototype

A prototype mobile ordering system built with Go and Echo.
This application demonstrates a unique constraint: **Only one guest is allowed to access the ordering interface of a specific restaurant at a time.**

## Features
- **System Admin**: Create and manage Restaurant accounts.
- **Restaurant Manager**: Manage menu items, prices, and stock/sold-out status.
- **Guest Access**:
  - No login required (Session-based).
  - **Single-User Lock**: If a guest enters a restaurant, others are blocked until that guest checks out or the session times out (5 mins).

## Tech Stack
- **Language**: Go 1.25+
- **Framework**: Echo v4
- **Database**: SQLite (via GORM)
- **Frontend**: Vanilla HTML/JS

## getting Started

### Prerequisites
- Go installed

### Running the App
1. Clone repository
2. Run the server:
   ```bash
   go run main.go
   ```
3. Open `http://localhost:8080`

### Default Accounts
| Role | Email | Password |
|------|-------|----------|
| **Admin** | `admin@example.com` | `admin123` |
