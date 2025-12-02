# Database Seeding Script

## Overview
This script populates the database with dummy data for testing and development purposes.

## What it creates:
- **1 Lecturer**: Dr. Adebayo Olumide (Computer Science Department)
- **15 Students**: Nigerian names from various engineering departments
- **5 Course Events**: CS courses (CSC101-CSC501)
- **Random Attendance Records**: 60-100% attendance rate per event

## Usage

### First time setup (clean database):
```bash
docker compose down -v
docker compose up -d
sleep 15
./scripts/seed-database.sh
```

### Quick reseed (existing database will have duplicates):
```bash
./scripts/seed-database.sh
```

## Login Credentials

After running the script, all login credentials are saved to `seed-login-credentials.txt` in the project root.

### Default Credentials:

**Lecturer:**
- Email: `dr.adebayo.olumide@fupre.edu.ng`
- Password: `Lecturer@123`
- Staff ID: `FUPRE/LECT/001`

**Students:**
- Emails follow pattern: `firstname.lastname@fupre.edu.ng`
- Passwords: `Student@100` through `Student@114`
- Matric Numbers: `FUPRE/2021-2023/10000-10014`

## Features

- ✅ Waits for API to be ready before seeding
- ✅ Creates diverse Nigerian student names
- ✅ Generates future-dated events (to allow check-ins)
- ✅ Simulates realistic attendance patterns (60-100% per event)
- ✅ Saves all credentials to file
- ✅ Color-coded console output
- ✅ Handles duplicate registrations gracefully

## API Endpoints Used

- `POST /api/auth/register-lecturer`
- `POST /api/auth/register-student`
- `POST /api/auth/login-lecturer`
- `POST /api/auth/login-student`
- `POST /api/lecturer/qrcode/generate`
- `POST /api/attendance/check-in`

## Environment Variables

- `API_URL`: API base URL (default: `http://localhost:2754`)

## Notes

- The script uses `@fupre.edu.ng` email domain for all users
- Events are created with dates 0-7 days in the future
- Attendance is randomized per event (different students attend different classes)
- All passwords meet the minimum security requirements
