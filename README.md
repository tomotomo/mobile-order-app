# Mobile Order App Prototype

A prototype mobile ordering system built with Go and Echo.
This application demonstrates a unique constraint: **Only one guest is allowed to access the ordering interface of a specific restaurant at a time.**

## Architecture
The system is split into three separate applications sharing a common database:

1.  **Guest App** (Port: **8080**): Public interface for customers.
    *   **Features**: Check-in, Menu View, Checkout.
    *   **Constraint**: Exclusive 1-user lock per restaurant.
2.  **Admin App** (Port: **8081**): Internal tool for System Admins.
    *   **Features**: Create Restaurants, Send Welcome Emails (Trap/SMTP).
3.  **Restaurant App** (Port: **8082**): Management tool for Restaurant Managers.
    *   **Features**: Add Menu Items, Update Stock/Sold-Out status.

## Tech Stack
- **Language**: Go 1.25+
- **Framework**: Echo v4
- **Database**: SQLite (via GORM)
- **Frontend**: Vanilla HTML/JS (SPA-like per app)

## Getting Started

### Prerequisites
- Go installed

### Running the Apps
You need to run these commands in separate terminals. The browser will open automatically.

**1. Guest App**
```bash
go run cmd/guest/main.go
# Opens http://localhost:8080
```

**2. Admin App**
```bash
go run cmd/admin/main.go
# Opens http://localhost:8081
# Login: admin@example.com / admin123
```

**3. Restaurant App**
```bash
go run cmd/restaurant/main.go
# Opens http://localhost:8082
# Login: Use credentials created via Admin App
```

## Features
- **Auto-Open Browser**: Apps automatically open the default browser on startup.
- **Email Trap**: Emails (e.g., welcome messages) are logged to stdout by default. Configure `internal/email` to use real SMTP.
