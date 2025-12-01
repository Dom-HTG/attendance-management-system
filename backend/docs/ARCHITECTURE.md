# Architecture Summary

This is a short architecture overview to give context to developers integrating with the project.

Structure
- `cmd/api` - main application entry
- `config` - app configuration and dependency injection (wiring services, repos, middleware)
- `internal/auth` - authentication domain, repository and service
- `internal/attendance` - attendance domain, repository and service
- `entities` - GORM entity definitions for users, events, attendance records
- `pkg/middleware` - auth middleware and role-based middleware
- `pkg/utils` - helpers such as QR code generation

Patterns
- Layered design (handlers -> services -> repositories -> entities)
- Dependency injection in `config/app/app.config.go`
- GORM for ORM and AutoMigrate on start
- JWT for authentication, role enforced by middleware

Database
- PostgreSQL (configure via environment)
- Core tables: users (students/lecturers), events, user_attendance
- Unique constraints: QR token uniqueness, user-attendance uniqueness to prevent duplicates

QR Code
- QR token: UUID v4 stored on `events` table
- QR image: base64 PNG returned on QR generation - frontend can display via data URI
- Check-in uses `qr_token` field to find corresponding event and insert attendance

Security notes
- Always use HTTPS in production
- JWT secret must be set in environment and kept secret
- Rate limit endpoints (recommended for production)

Running and local development
- See `docs/QUICK_START.md` for minimal run instructions

Where to look in code
- Handlers: `internal/*/service/*.go` (contains Gin handlers)
- Repositories: `internal/*/repository/*.go` (DB access)
- DTOs: `internal/*/domain/*.go` for request/response shapes
- Middleware: `pkg/middleware/auth.middleware.go`
- QR helper: `pkg/utils/qrcode.go`
