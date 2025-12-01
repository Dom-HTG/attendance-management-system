# Quick start — Run locally

This file describes the minimal steps to run the Attendance Management System locally for development and testing.

Prerequisites
- Go (1.20+ recommended; the project uses Go modules)
- PostgreSQL database
- Environment: you can run on WSL or native Windows (see .env / app.env in repo)

1) Configure environment
- Copy `config/app/app.env` or set environment variables used by the app (DB connection, APP_PORT)
- Typical vars: APP_PORT, DB_HOST, DB_PORT, DB_USER, DB_PASSWORD, DB_NAME, JWT_SECRET

2) Start database (example using docker-compose)
```bash
docker compose up -d
```

3) Build & run the server
```bash
# from project root
go build ./cmd/api
./api    # or `go run ./cmd/api` for development
```

4) Auto-migrate
- The application runs GORM auto-migrations on start for created entities (events, user_attendance, users). Confirm the DB tables exist.

5) Test endpoints
- Use the Postman collection (if available) or use `docs/API.md` for all endpoints and concrete request/response examples.
- Typical flow: register a lecturer, login lecturer -> generate QR code -> register a student, login student -> student checks in using qr_token -> verify attendance via lecturer endpoint.

Notes
- The server logs and errors will be printed to stdout. Check logs for DB connection issues.
- If tokens expire, re-login. Tokens are signed with the configured JWT_SECRET.
- For local development, set APP_PORT=2754 (default used by docs) or adapt BASE_URL accordingly.
