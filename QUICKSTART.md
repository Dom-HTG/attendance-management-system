# Quick Start Guide - Authentication System

## 🚀 Getting Started in 5 Minutes

### Prerequisites
- PostgreSQL running locally
- Go 1.23.0 or higher
- Postman or curl

---

## Step 1: Database Setup

```bash
# Create the database
psql -U postgres -c "CREATE DATABASE \"attendance-management\";"

# Verify connection
psql -U postgres -d attendance-management -c "\dt"
```

**Expected Output:** Empty (tables will be created automatically by GORM)

---

## Step 2: Run the Application

```bash
# Navigate to project directory
cd attendance-management

# Download dependencies
go mod download

# Run the application
go run cmd/api/main.go
```

**Expected Output:**
```
Database connection established successfully..
[GIN-debug] Loaded HTML Templates (0): 
[GIN-debug] Listening and serving HTTP on :2754
```

---

## Step 3: Test Registration & Login

### Option A: Using Postman

1. Import `postman_collection.json` into Postman
2. Run requests in this order:
   - Student Registration
   - Student Login
   - Lecturer Registration
   - Lecturer Login

### Option B: Using curl

**Register a Student:**
```bash
curl -X POST http://localhost:2754/api/auth/register-student \
  -H "Content-Type: application/json" \
  -d '{
    "first_name": "Alice",
    "last_name": "Johnson",
    "email": "alice@student.edu",
    "password": "SecurePass123",
    "matric_number": "STU-2024-001"
  }'
```

**Expected Response (201):**
```json
{
  "success": true,
  "message": "Student successfully registered",
  "data": {
    "message": "Student successfully registered. Please login with your credentials."
  }
}
```

---

**Login as Student:**
```bash
curl -X POST http://localhost:2754/api/auth/login-student \
  -H "Content-Type: application/json" \
  -d '{
    "email": "alice@student.edu",
    "password": "SecurePass123"
  }'
```

**Expected Response (200):**
```json
{
  "success": true,
  "message": "Login successful",
  "data": {
    "message": "Student login successful",
    "access_token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpZCI6MSwi...",
    "user": {
      "id": 1,
      "first_name": "Alice",
      "last_name": "Johnson",
      "email": "alice@student.edu",
      "matric_number": "STU-2024-001",
      "role": "student",
      "created_at": "2025-11-27T15:30:45.123456Z"
    }
  }
}
```

---

**Register a Lecturer:**
```bash
curl -X POST http://localhost:2754/api/auth/register-lecturer \
  -H "Content-Type: application/json" \
  -d '{
    "first_name": "Dr.",
    "last_name": "Smith",
    "email": "smith@lecturer.edu",
    "password": "SecurePass123",
    "department": "Computer Science",
    "staff_id": "STAFF-2024-001"
  }'
```

---

**Login as Lecturer:**
```bash
curl -X POST http://localhost:2754/api/auth/login-lecturer \
  -H "Content-Type: application/json" \
  -d '{
    "email": "smith@lecturer.edu",
    "password": "SecurePass123"
  }'
```

---

## Step 4: Verify in Database

```bash
# Connect to the database
psql -U postgres -d attendance-management

# Check tables created
\dt

# View students
SELECT id, first_name, last_name, email, matric_number, role FROM students;

# View lecturers
SELECT id, first_name, last_name, email, department, staff_id, role FROM lecturers;
```

---

## 🔑 Key Endpoints

| Endpoint | Method | Purpose |
|----------|--------|---------|
| `/api/auth/register-student` | POST | Register new student |
| `/api/auth/login-student` | POST | Login as student (get token) |
| `/api/auth/register-lecturer` | POST | Register new lecturer |
| `/api/auth/login-lecturer` | POST | Login as lecturer (get token) |

---

## ⚙️ Configuration

### Environment Variables (`.env`)
```bash
# Server
APP_PORT=":2754"

# Database
DB_HOST=localhost
DB_USER=postgres
DB_PASSWORD=Ac101bb101
DB_NAME="attendance-management"
DB_PORT=5432

# JWT
JWT_SECRET=your-super-secret-key-change-in-production-use-env-var
```

**Change in Production:** Always update `JWT_SECRET` with a strong random key!

---

## 📝 Request/Response Examples

### Success Response Format
```json
{
  "success": true,
  "message": "Operation successful",
  "data": {
    // Response data here
  }
}
```

### Error Response Format
```json
{
  "success": false,
  "error_message": "Human readable error",
  "error": "Technical details (optional)"
}
```

### HTTP Status Codes
- `200` - OK (Login successful)
- `201` - Created (Registration successful)
- `400` - Bad Request (Validation error)
- `401` - Unauthorized (Invalid credentials)
- `500` - Internal Server Error (Database/server error)

---

## 🧪 Common Test Scenarios

### ✅ Valid Flow
1. Register student with valid data → 201 Created
2. Login with correct credentials → 200 OK + JWT token

### ❌ Error Cases

**Missing Required Field:**
```bash
curl -X POST http://localhost:2754/api/auth/register-student \
  -H "Content-Type: application/json" \
  -d '{"first_name": "John"}'
  # Expected: 400 Bad Request
```

**Invalid Email Format:**
```bash
curl -X POST http://localhost:2754/api/auth/register-student \
  -H "Content-Type: application/json" \
  -d '{
    "first_name": "John",
    "last_name": "Doe",
    "email": "invalid-email",
    "password": "pass123",
    "matric_number": "STU-001"
  }'
  # Expected: 400 Bad Request
```

**Password Too Short:**
```bash
curl -X POST http://localhost:2754/api/auth/register-student \
  -H "Content-Type: application/json" \
  -d '{
    "first_name": "John",
    "last_name": "Doe",
    "email": "john@example.com",
    "password": "123",
    "matric_number": "STU-001"
  }'
  # Expected: 400 Bad Request (minimum 6 characters)
```

**Invalid Credentials Login:**
```bash
curl -X POST http://localhost:2754/api/auth/login-student \
  -H "Content-Type: application/json" \
  -d '{
    "email": "john@example.com",
    "password": "wrongpassword"
  }'
  # Expected: 401 Unauthorized
```

---

## 🛠️ Troubleshooting

### Issue: "Connection refused"
**Solution:** Ensure your server is running with `go run cmd/api/main.go`

### Issue: "Database connection error"
**Solution:** Verify PostgreSQL is running and credentials in `.env` are correct

### Issue: "Port already in use"
**Solution:** Change `APP_PORT` in `.env` to a different port (e.g., `:2755`)

### Issue: "Email already exists"
**Solution:** Use a different email for testing (emails must be unique)

### Issue: Token not working in other endpoints
**Solution:** Other endpoints may require middleware to validate the token (to be implemented)

---

## 📚 Documentation Files

- **AUTH_SYSTEM.md** - Complete API documentation with all endpoints
- **BUILD_SUMMARY.md** - Detailed build process and changes made
- **postman_collection.json** - Ready-to-use Postman collection
- **this file** - Quick start guide

---

## 💡 Next Features to Implement

1. **Logout Endpoint** - Invalidate tokens
2. **Refresh Token** - Extend session
3. **Forgot Password** - Email reset link
4. **Token Validation Middleware** - Protect routes
5. **Email Verification** - On registration
6. **Rate Limiting** - Prevent brute force attacks

---

## ✨ You're All Set!

The authentication system is ready to use. Start with:

```bash
# Terminal 1: Run the server
go run cmd/api/main.go

# Terminal 2: Test endpoints
curl -X POST http://localhost:2754/api/auth/register-student ...
```

Happy coding! 🚀
