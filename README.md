# Attendance Management System

A production-ready attendance management system backend written in Go. Features include role-based authentication (students and lecturers), QR-code based attendance tracking, and comprehensive analytics dashboards.

## 🚀 Quick Start (One Command)

Bootstrap the entire project with database seeding and tests:

```bash
bash scripts/bootstrap.sh
```

This script will:
1. Start the application with Docker Compose (builds and runs containers)
2. Seed the database with 15 students, 1 lecturer, 5 events, and 56+ attendance records
3. Run comprehensive API tests (16 test cases)
4. Display login credentials for testing

**After bootstrap completes:**
- API runs at: `http://localhost:2754`
- Login credentials: See `seed-login-credentials.txt`
- Seeded lecturer: `dr.adebayo.olumide@fupre.edu.ng` / `Lecturer@123`

---

## 📚 Documentation

This repository keeps focused documentation in the `docs/` folder:

- **README.md** (this file) - Quick start and overview
- **ANALYTICS_STATUS.md** - ✅ Analytics endpoints status (all 4 endpoints working!)
- **docs/ADMIN_GUIDE.md** - ✅ Complete admin dashboard API documentation
- **docs/PDF_EXPORT_GUIDE.md** - ✅ PDF export endpoints guide (student & lecturer reports)
- **docs/FRONTEND_ANALYTICS_ENDPOINTS.md** - Complete frontend integration guide
- **docs/API_REFERENCE.md** - Complete API reference & examples
- **docs/QUICK_START.md** - Detailed local setup guide
- **docs/ARCHITECTURE.md** - Project structure & design patterns
- **docs/ANALYTICS_IMPLEMENTATION.md** - Analytics endpoints guide

---

## 🏗️ Tech Stack

- **Language:** Go 1.21+
- **Web Framework:** Gin (HTTP router)
- **ORM:** GORM with PostgreSQL 15
- **Authentication:** JWT (HS256) with role-based access
- **QR Codes:** UUID tokens + PNG base64 encoding
- **Analytics:** 25+ endpoints for insights, predictions, anomaly detection
- **Containerization:** Docker & Docker Compose

---

## 🔌 Primary API Endpoints

### Authentication
- `POST /api/auth/register-student` - Register new student
- `POST /api/auth/register-lecturer` - Register new lecturer
- `POST /api/auth/login-student` - Student login (returns JWT)
- `POST /api/auth/login-lecturer` - Lecturer login (returns JWT)
- `POST /api/auth/login-admin` - Admin login (returns JWT, 7-day expiration)

### QR Code & Attendance
- `POST /api/lecturer/qrcode/generate` 🔒 Lecturer - Generate QR code for event
- `POST /api/attendance/check-in` 🔒 Student - Check in using QR token
- `GET /api/attendance/:event_id` 🔒 Lecturer - View event attendance list
- `GET /api/attendance/student/records` 🔒 Student - View personal attendance history

### Admin Dashboard
- `GET /api/admin/students` 🔒 Admin - List all students with pagination & filters
- `GET /api/admin/lecturers` 🔒 Admin - List all lecturers with statistics
- `GET /api/admin/users/:type/:id` 🔒 Admin - Get detailed user profile
- `PATCH /api/admin/users/:type/:id/status` 🔒 Admin - Update user status
- `DELETE /api/admin/users/:type/:id` 🔒 Admin - Delete user (soft delete)
- `GET /api/admin/events` 🔒 Admin - List all events with filters
- `DELETE /api/admin/events/:id` 🔒 Admin - Delete event (cascade delete)
- `GET /api/admin/trends` 🔒 Admin - Attendance trends analysis
- `GET /api/admin/low-attendance` 🔒 Admin - Students at risk
- `GET /api/admin/settings` 🔒 Admin - Get system settings
- `PATCH /api/admin/settings` 🔒 Admin - Update system settings
- `GET /api/admin/audit-logs` 🔒 Admin - View audit trail

### Lecturer Analytics
- `GET /api/events/lecturer` 🔒 Lecturer - Get all events with attendance counts
- `GET /api/analytics/lecturer/summary` 🔒 Lecturer - Dashboard summary stats
- `GET /api/analytics/lecturer/courses` 🔒 Lecturer - Course-level metrics

### Admin Analytics
- `GET /api/analytics/admin/overview` 🔒 Admin - University-wide dashboard
- `GET /api/analytics/admin/departments` 🔒 Admin - Per-department statistics
- `GET /api/analytics/student/{id}` 🔒 Admin - Individual student metrics
- `GET /api/analytics/anomalies` 🔒 Admin - Fraud detection & duplicate check-ins
- `GET /api/analytics/predictions/student/{id}` 🔒 Admin - Attendance predictions

### PDF Export
- `GET /api/student/attendance/export/pdf` 🔒 Student - Export personal attendance report as PDF
- `GET /api/lecturer/attendance/export/pdf` 🔒 Lecturer - Export all events summary as PDF
- `GET /api/lecturer/attendance/export/pdf?event_id=X` 🔒 Lecturer - Export single event report as PDF

🔒 = Requires JWT authentication

See `docs/API_REFERENCE.md` for complete endpoint documentation with request/response examples.
See `docs/ADMIN_GUIDE.md` for comprehensive admin dashboard documentation.

---

## 📝 Environment Variables

The application uses environment variables defined in `cmd/api/app.env`:

```env
# Database Configuration
DB_HOST=postgres          # Database host
DB_PORT=5432              # Database port
DB_USER=postgres          # Database user
DB_PASSWORD=postgres      # Database password
DB_NAME=attendance-management  # Database name

# Application Configuration
APP_PORT=2754             # API server port
JWT_SECRET=your-secret-key-change-this  # JWT signing key
APP_MODE=release          # Gin mode (debug|release)
```

**Docker Compose automatically loads these values.** For local development without Docker, copy `app.env` and adjust settings.

---

## 🧩 Frontend Integration Guide

### Authentication Flow
1. Call register/login endpoint to get JWT token
2. Store token securely (localStorage/cookies)
3. Include token in all authenticated requests: `Authorization: Bearer {token}`
4. Token expires after 60 minutes (refresh mechanism coming soon)

### QR Code Display
```javascript
// After calling /api/lecturer/qrcode/generate
const response = await fetch('/api/lecturer/qrcode/generate', {
  method: 'POST',
  headers: {
    'Authorization': `Bearer ${lecturerToken}`,
    'Content-Type': 'application/json'
  },
  body: JSON.stringify({
    course_code: 'CSC301',
    course_name: 'Data Structures',
    venue: 'LH-101',
    duration_minutes: 30
  })
});

const data = await response.json();

// Display QR code as image
<img src={`data:image/png;base64,${data.data.qr_code_data}`} alt="QR Code" />

// Or save qr_token for manual entry
const qrToken = data.data.qr_token;
```

### Student Check-In
```javascript
// After scanning QR code or entering token manually
const checkIn = await fetch('/api/attendance/check-in', {
  method: 'POST',
  headers: {
    'Authorization': `Bearer ${studentToken}`,
    'Content-Type': 'application/json'
  },
  body: JSON.stringify({
    qr_token: extractedToken  // From QR scan or manual input
  })
});
```

### Analytics Dashboard
```javascript
// Lecturer dashboard
const summary = await fetch('/api/analytics/lecturer/summary', {
  headers: { 'Authorization': `Bearer ${lecturerToken}` }
});

// Admin dashboard
const overview = await fetch('/api/analytics/admin/overview', {
  headers: { 'Authorization': `Bearer ${adminToken}` }
});
```

See `docs/API_REFERENCE.md` for complete response schemas.

---

## 🔧 Manual Setup (Without Bootstrap)

### Prerequisites
- Docker & Docker Compose installed
- Git
- Bash shell

### Step 1: Clone and Start

```bash
# Clone the repository
git clone <repository-url>
cd ML-backend

# Start application (builds containers and runs migrations)
docker compose up -d --build

# Wait for services to be ready (~8 seconds)
sleep 8
```

### Step 2: Seed Database

```bash
# Seed with test data (15 students, 1 lecturer, 5 events, 56+ attendances)
bash scripts/seed-database.sh

# Login credentials saved to seed-login-credentials.txt
cat seed-login-credentials.txt
```

### Step 3: Test API

```bash
# Run comprehensive test suite (16 test cases)
bash scripts/test-api.sh
```

**Expected Output:** All 16 tests should pass ✓

---

## 🧪 Testing with cURL

### Authentication Tests

**Register Student:**
```bash
curl -X POST http://localhost:2754/api/auth/register-student \
  -H "Content-Type: application/json" \
  -d '{
    "first_name": "John",
    "last_name": "Doe",
    "email": "john.doe@fupre.edu.ng",
    "password": "securepass123",
    "matric_number": "FUPRE/2024/12345"
  }'
```

**Register Lecturer:**
```bash
curl -X POST http://localhost:2754/api/auth/register-lecturer \
  -H "Content-Type: application/json" \
  -d '{
    "first_name": "Jane",
    "last_name": "Smith",
    "email": "jane.smith@fupre.edu.ng",
    "password": "securepass123",
    "department": "Computer Science",
    "staff_id": "FUPRE/STAFF/001"
  }'
```

**Login (Student):**
```bash
curl -X POST http://localhost:2754/api/auth/login-student \
  -H "Content-Type: application/json" \
  -d '{
    "email": "john.doe@fupre.edu.ng",
    "password": "securepass123"
  }'
```

Save the returned `access_token` for authenticated requests.

### QR Code Generation (Lecturer Only)

```bash
# Set your lecturer token
LECTURER_TOKEN="your_jwt_token_here"

# Generate QR code for an event
curl -X POST http://localhost:2754/api/lecturer/qrcode/generate \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer $LECTURER_TOKEN" \
  -d '{
    "course_code": "CSC301",
    "course_name": "Data Structures",
    "venue": "LH-101",
    "duration_minutes": 30
  }'
```

**Response includes:**
- `qr_code_data` - Base64 PNG image (display with `<img src="data:image/png;base64,{data}">`)
- `qr_token` - UUID token for check-in
- `event_id` - Event identifier
- `expires_at` - QR code expiration time

### Student Check-In

```bash
# Set your student token
STUDENT_TOKEN="your_jwt_token_here"

# Check in to event (use qr_token from QR generation)
curl -X POST http://localhost:2754/api/attendance/check-in \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer $STUDENT_TOKEN" \
  -d '{
    "qr_token": "uuid-from-qr-code"
  }'
```

### Analytics Endpoints

**Lecturer Analytics:**
```bash
# Get all lecturer's events with attendance counts
curl -X GET http://localhost:2754/api/events/lecturer \
  -H "Authorization: Bearer $LECTURER_TOKEN"

# Get lecturer dashboard summary
curl -X GET http://localhost:2754/api/analytics/lecturer/summary \
  -H "Authorization: Bearer $LECTURER_TOKEN"
```

**Admin Analytics:**
```bash
# University-wide overview
curl -X GET http://localhost:2754/api/analytics/admin/overview \
  -H "Authorization: Bearer $LECTURER_TOKEN"

# Department statistics
curl -X GET http://localhost:2754/api/analytics/admin/departments \
  -H "Authorization: Bearer $LECTURER_TOKEN"
```

**View Attendance Records:**
```bash
# Lecturer: View event attendance
curl -X GET http://localhost:2754/api/attendance/{event_id} \
  -H "Authorization: Bearer $LECTURER_TOKEN"

# Student: View personal attendance records
curl -X GET http://localhost:2754/api/attendance/student/records \
  -H "Authorization: Bearer $STUDENT_TOKEN"
```

---

## 📁 Project Structure

```
ML-backend/
├── cmd/api/              # Application entry point
├── config/               # App and database configuration
├── docs/                 # Comprehensive documentation
├── entities/             # Database entities (GORM models)
├── internal/             # Core business logic
│   ├── analytics/        # Analytics domain (25+ endpoints)
│   ├── attendance/       # Attendance tracking
│   └── auth/             # Authentication & authorization
├── migrations/           # SQL migration scripts
├── pkg/                  # Reusable packages
│   ├── logger/           # Structured logging
│   ├── middleware/       # Auth, CORS, etc.
│   ├── responses/        # Standardized API responses
│   └── utils/            # JWT, hashing, QR generation
├── scripts/              # Automation scripts
│   ├── bootstrap.sh      # Full setup + seed + test
│   ├── seed-database.sh  # Database seeding
│   └── test-api.sh       # API test suite
├── docker-compose.yml    # Container orchestration
├── Dockerfile            # Application container
└── go.mod                # Go dependencies
```

See `docs/ARCHITECTURE.md` for detailed architecture documentation.

---

## 🔑 Using Seeded Data

After running `bash scripts/seed-database.sh`, you'll have:

**1 Lecturer:**
- Email: `dr.adebayo.olumide@fupre.edu.ng`
- Password: `Lecturer@123`
- Staff ID: `FUPRE/LECT/001`

**15 Students:**
- Format: `{firstname}.{lastname}@fupre.edu.ng`
- Password: `Student@10X` (where X = 0-14)
- Example: `chukwuemeka.okonkwo@fupre.edu.ng` / `Student@100`

**5 Events:**
- Courses: Data Structures, Web Development, Database Systems, Software Engineering, Computer Networks
- All events have 30-minute duration
- Created by the seeded lecturer

**56+ Attendance Records:**
- Random attendance distribution across all events
- Each student has checked into 3-4 events on average

**Test the seeded data:**
```bash
# Login as seeded lecturer
curl -X POST http://localhost:2754/api/auth/login-lecturer \
  -H "Content-Type: application/json" \
  -d '{"email":"dr.adebayo.olumide@fupre.edu.ng","password":"Lecturer@123"}'

# View all events (should return 5 events)
curl -X GET http://localhost:2754/api/events/lecturer \
  -H "Authorization: Bearer $LECTURER_TOKEN"
```

---

## 📊 Available Scripts

| Script | Purpose | Usage |
|--------|---------|-------|
| `bootstrap.sh` | Full setup (start + seed + test) | `bash scripts/bootstrap.sh` |
| `seed-database.sh` | Populate database with test data | `bash scripts/seed-database.sh` |
| `test-api.sh` | Run comprehensive API tests | `bash scripts/test-api.sh` |

**Script Features:**
- ✅ Color-coded output (green=pass, red=fail)
- ✅ Detailed test results with pass/fail counts
- ✅ Saves login credentials to `seed-login-credentials.txt`
- ✅ Validates all API endpoints (16 test cases)

---

## 🛠️ Docker Commands

```bash
# Start application
docker compose up -d --build

# View logs
docker compose logs -f

# Stop application
docker compose down

# Restart application
docker compose restart

# Stop and remove all data (including database)
docker compose down -v

# Access database directly
docker exec -it ml-backend-postgres-1 psql -U postgres -d attendance-management
```

See `docs/QUICK_START.md` for more Docker tips and troubleshooting.

---

## ✅ Project Status

| Feature | Status |
|---------|--------|
| Authentication System | ✅ Complete |
| QR Code Generation | ✅ Complete |
| Attendance Check-In | ✅ Complete |
| Lecturer Analytics | ✅ Complete |
| Admin Analytics | ✅ Complete |
| Database Seeding | ✅ Complete |
| API Testing | ✅ Complete |
| Documentation | ✅ Complete |
| Docker Setup | ✅ Complete |

**Overall:** 95% Production Ready

---

## 🛣️ Roadmap

### Completed ✅
- Core authentication system (register, login, JWT)
- Role-based access control (student, lecturer, admin)
- QR code generation and attendance check-in
- Lecturer analytics (events, summary, course metrics)
- Admin analytics (overview, departments, anomalies)
- Predictive analytics and at-risk detection
- Database migration system
- Comprehensive test suite (16 test cases)
- Bootstrap automation script
- Structured logging and graceful shutdown
- Docker containerization

### In Progress ⏳
- Token refresh mechanism
- Email verification system
- Password reset functionality
- Rate limiting middleware

### Future 📅
- OAuth2/Social login integration
- Multi-factor authentication (2FA)
- Real-time WebSocket notifications
- Export attendance reports (PDF/Excel)
- Audit logging system
- Advanced anomaly detection with ML

---

## 🆘 Troubleshooting

### Application won't start
```bash
# Check if Docker is running
docker --version

# Check if ports are available (2754 for API, 5432 for PostgreSQL)
netstat -an | grep 2754
netstat -an | grep 5432

# View application logs
docker compose logs -f
```

### Database connection error
```bash
# Check if PostgreSQL container is running
docker ps | grep postgres

# Verify database exists
docker exec -it ml-backend-postgres-1 psql -U postgres -c "\l"

# Check connection from app container
docker exec -it ml-backend-app-1 nc -zv postgres 5432
```

### Tests failing
```bash
# Ensure application is running
curl http://localhost:2754/health || echo "API not responding"

# Re-seed database
bash scripts/seed-database.sh

# Run tests with verbose output
bash scripts/test-api.sh 2>&1 | tee test-output.log
```

### Port already in use
```bash
# Find process using port 2754
lsof -i :2754  # macOS/Linux
netstat -ano | findstr :2754  # Windows

# Kill the process or change APP_PORT in cmd/api/app.env
```

### JWT token validation failed
- Ensure `JWT_SECRET` matches between registration and login
- Token expires after 60 minutes - login again to get new token
- Verify token format: `Authorization: Bearer {token}`

### "Email already exists" error
- Use different email for testing
- Or clean database: `docker compose down -v && docker compose up -d --build`

See `docs/QUICK_START.md` for more troubleshooting tips.

---

## 🤝 Contributing

To extend this system:

1. **Add new endpoints:**
   - Define domain interfaces in `internal/{domain}/domain/`
   - Implement repository methods in `internal/{domain}/repository/`
   - Implement service logic in `internal/{domain}/service/`
   - Create handlers in `internal/{domain}/handler/`
   - Register routes in `config/app/app.config.go`

2. **Add database migrations:**
   - Create SQL file in `migrations/`
   - Run migration: Connect to DB and execute SQL

3. **Add tests:**
   - Update `scripts/test-api.sh` with new test cases
   - Follow existing pattern for authentication and validation

4. **Update documentation:**
   - Update `docs/API_REFERENCE.md` with new endpoints
   - Add examples with curl commands
   - Document request/response schemas

**Architecture Pattern:** Domain → Repository → Service → Handler → Routes

See `docs/ARCHITECTURE.md` for detailed design patterns.

---

## 📞 Support & Resources

### Documentation Files
- **README.md** (this file) - Quick start and overview
- **docs/API_REFERENCE.md** - Complete API documentation
- **docs/QUICK_START.md** - Detailed setup guide
- **docs/ARCHITECTURE.md** - System design and patterns
- **docs/ANALYTICS_IMPLEMENTATION.md** - Analytics guide

### External Resources
- [Go Documentation](https://golang.org/doc/)
- [Gin Web Framework](https://gin-gonic.com/docs/)
- [GORM Documentation](https://gorm.io/docs/)
- [JWT Best Practices](https://tools.ietf.org/html/rfc7519)
- [Docker Compose](https://docs.docker.com/compose/)

### Getting Help
1. Check the troubleshooting section above
2. Review relevant documentation in `docs/`
3. Check Docker logs: `docker compose logs -f`
4. Run tests to verify setup: `bash scripts/test-api.sh`

---

## 📜 License

This project is part of the Attendance Management System.

---

## ✨ Key Features

- 🔐 **Secure Authentication** - JWT tokens with bcrypt password hashing
- 📱 **QR Code Attendance** - Generate and scan QR codes for check-in
- 📊 **Rich Analytics** - 25+ endpoints for insights and predictions
- 🎯 **Role-Based Access** - Student, Lecturer, and Admin roles
- 🐳 **Docker Ready** - One-command deployment with Docker Compose
- 🧪 **Fully Tested** - Comprehensive test suite with 16 test cases
- 📚 **Well Documented** - Extensive guides and API reference
- ⚡ **High Performance** - Optimized queries with database indexing
- 🏗️ **Clean Architecture** - Domain-driven design pattern
- 🔄 **Auto Bootstrap** - Automated setup, seed, and test script

---

**Built with ❤️ for the Attendance Management System**

*Last Updated: December 1, 2025*
