# Authentication System - Attendance Management

## Overview

This document describes the complete end-to-end authentication system implemented for the Attendance Management System. The system supports role-based authentication for both **Students** and **Lecturers** with JWT token-based authorization.

## Features

✅ **Student Registration** - Register new students with email validation
✅ **Lecturer Registration** - Register new lecturers with department assignment
✅ **Student Login** - Secure login with JWT token generation
✅ **Lecturer Login** - Secure login with JWT token generation
✅ **Password Hashing** - Bcrypt-based password hashing (cost factor: 10)
✅ **JWT Token Generation** - 60-minute token expiration
✅ **Role-Based Access** - Tokens include user role information
✅ **Input Validation** - Binding validation with required fields and format checking

## API Endpoints

### Authentication Routes

Base URL: `/api/auth`

#### 1. Register Student
```http
POST /api/auth/register-student
Content-Type: application/json

{
  "first_name": "John",
  "last_name": "Doe",
  "email": "john.doe@student.edu",
  "password": "securePassword123",
  "matric_number": "STU-2024-001"
}
```

**Success Response (201):**
```json
{
  "success": true,
  "message": "Student successfully registered",
  "data": {
    "message": "Student successfully registered. Please login with your credentials."
  }
}
```

**Error Response (400/500):**
```json
{
  "success": false,
  "error_message": "Unable to bind request body",
  "error": "..."
}
```

---

#### 2. Register Lecturer
```http
POST /api/auth/register-lecturer
Content-Type: application/json

{
  "first_name": "Jane",
  "last_name": "Smith",
  "email": "jane.smith@lecturer.edu",
  "password": "securePassword123",
  "department": "Computer Science",
  "staff_id": "STAFF-2024-001"
}
```

**Success Response (201):**
```json
{
  "success": true,
  "message": "Lecturer successfully registered",
  "data": {
    "message": "Lecturer successfully registered. Please login with your credentials."
  }
}
```

---

#### 3. Login Student
```http
POST /api/auth/login-student
Content-Type: application/json

{
  "email": "john.doe@student.edu",
  "password": "securePassword123"
}
```

**Success Response (200):**
```json
{
  "success": true,
  "message": "Login successful",
  "data": {
    "message": "Student login successful",
    "access_token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
    "user": {
      "id": 1,
      "first_name": "John",
      "last_name": "Doe",
      "email": "john.doe@student.edu",
      "matric_number": "STU-2024-001",
      "role": "student",
      "created_at": "2025-11-27T10:30:00Z"
    }
  }
}
```

**Error Response (401):**
```json
{
  "success": false,
  "error_message": "Invalid email or password",
  "error": null
}
```

---

#### 4. Login Lecturer
```http
POST /api/auth/login-lecturer
Content-Type: application/json

{
  "email": "jane.smith@lecturer.edu",
  "password": "securePassword123"
}
```

**Success Response (200):**
```json
{
  "success": true,
  "message": "Login successful",
  "data": {
    "message": "Lecturer login successful",
    "access_token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
    "user": {
      "id": 1,
      "first_name": "Jane",
      "last_name": "Smith",
      "email": "jane.smith@lecturer.edu",
      "department": "Computer Science",
      "staff_id": "STAFF-2024-001",
      "role": "lecturer",
      "created_at": "2025-11-27T10:30:00Z"
    }
  }
}
```

---

## Database Schema

### Student Table
```sql
CREATE TABLE students (
  id BIGINT PRIMARY KEY AUTO_INCREMENT,
  first_name VARCHAR(255) NOT NULL,
  last_name VARCHAR(255) NOT NULL,
  email VARCHAR(255) NOT NULL UNIQUE,
  role VARCHAR(50) DEFAULT 'student',
  password VARCHAR(255) NOT NULL,
  matric_number VARCHAR(50) NOT NULL UNIQUE,
  created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
);
```

### Lecturer Table
```sql
CREATE TABLE lecturers (
  id BIGINT PRIMARY KEY AUTO_INCREMENT,
  first_name VARCHAR(255) NOT NULL,
  last_name VARCHAR(255) NOT NULL,
  email VARCHAR(255) NOT NULL UNIQUE,
  role VARCHAR(50) DEFAULT 'lecturer',
  password VARCHAR(255) NOT NULL,
  department VARCHAR(255) NOT NULL,
  staff_id VARCHAR(50) NOT NULL UNIQUE,
  created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
);
```

---

## JWT Token Structure

The JWT token contains the following claims:

```json
{
  "id": 1,
  "email": "john.doe@student.edu",
  "role": "student",
  "exp": 1700000000,
  "iat": 1699996400
}
```

**Token Expiration:** 60 minutes
**Algorithm:** HS256 (HMAC SHA-256)

---

## Input Validation Rules

### Student Registration
- `first_name`: Required
- `last_name`: Required
- `email`: Required, must be valid email format
- `password`: Required, minimum 6 characters
- `matric_number`: Required

### Lecturer Registration
- `first_name`: Required
- `last_name`: Required
- `email`: Required, must be valid email format
- `password`: Required, minimum 6 characters
- `department`: Required
- `staff_id`: Required

### Login (Both)
- `email`: Required, must be valid email format
- `password`: Required

---

## Error Handling

All error responses follow this format:

```json
{
  "success": false,
  "error_message": "Human readable error message",
  "error": "Technical error details (optional)"
}
```

### Common Status Codes
- `200 OK` - Successful login
- `201 Created` - Successful registration
- `400 Bad Request` - Invalid request body or validation error
- `401 Unauthorized` - Invalid credentials
- `404 Not Found` - User not found
- `500 Internal Server Error` - Database or server error

---

## Security Features

1. **Password Hashing**: Bcrypt with cost factor of 10
2. **JWT Tokens**: Signed with HS256 algorithm
3. **Input Validation**: Server-side validation for all inputs
4. **Email Uniqueness**: Database constraints prevent duplicate emails
5. **Role-Based Design**: Separate login endpoints for students and lecturers

---

## Architecture

### Project Structure

```
internal/auth/
├── domain/
│   └── auth.go              # Interfaces and DTOs
├── repository/
│   └── auth.repository.go   # Database operations
└── service/
    └── auth.service.go      # Business logic

pkg/
├── utils/
│   ├── hashPassword.go      # Password hashing utility
│   └── jwt.go               # JWT token generation/validation
└── responses/
    ├── success.response.go  # Success response formatter
    └── failure.response.go  # Error response formatter

entities/
└── entities.go              # Student and Lecturer models

config/
└── app/
    └── app.config.go        # Route configuration and DI
```

### Layers

1. **Domain Layer** (`auth.go`)
   - Defines interfaces for repository and service
   - Contains DTOs for requests and responses

2. **Repository Layer** (`auth.repository.go`)
   - Handles all database operations
   - Implements CRUD operations for Student and Lecturer
   - Query methods for authentication

3. **Service Layer** (`auth.service.go`)
   - Business logic for registration and login
   - Password hashing and validation
   - JWT token generation
   - HTTP request/response handling

4. **Utility Layer**
   - Password hashing with Bcrypt
   - JWT token generation and validation

---

## Configuration

### Environment Variables (`.env`)

```bash
# SERVER
APP_PORT=":2754"

# DATABASE
DB_HOST=localhost
DB_USER=postgres
DB_PASSWORD=<your_password>
DB_NAME="attendance-management"
DB_PORT=5432

# DATABASE POOLING
POOL_MAX_OPEN_CONN=5
POOL_MAX_IDLE_CONN=3
POOL_MAX_CONN_TIMEOUT=1m

# JWT
JWT_SECRET=your-super-secret-key-change-in-production
```

---

## Dependencies

```
github.com/gin-gonic/gin v1.10.0           - HTTP web framework
github.com/golang-jwt/jwt/v5 v5.2.1        - JWT token generation
golang.org/x/crypto v0.36.0                - Bcrypt password hashing
gorm.io/gorm v1.25.12                      - ORM
gorm.io/driver/postgres v1.5.11            - PostgreSQL driver
```

---

## Testing Guide

### Using curl or Postman

**1. Register a Student:**
```bash
curl -X POST http://localhost:2754/api/auth/register-student \
  -H "Content-Type: application/json" \
  -d '{
    "first_name": "John",
    "last_name": "Doe",
    "email": "john@example.com",
    "password": "password123",
    "matric_number": "STU-001"
  }'
```

**2. Login as Student:**
```bash
curl -X POST http://localhost:2754/api/auth/login-student \
  -H "Content-Type: application/json" \
  -d '{
    "email": "john@example.com",
    "password": "password123"
  }'
```

**3. Register a Lecturer:**
```bash
curl -X POST http://localhost:2754/api/auth/register-lecturer \
  -H "Content-Type: application/json" \
  -d '{
    "first_name": "Jane",
    "last_name": "Smith",
    "email": "jane@example.com",
    "password": "password123",
    "department": "Computer Science",
    "staff_id": "STAFF-001"
  }'
```

**4. Login as Lecturer:**
```bash
curl -X POST http://localhost:2754/api/auth/login-lecturer \
  -H "Content-Type: application/json" \
  -d '{
    "email": "jane@example.com",
    "password": "password123"
  }'
```

---

## Future Enhancements

- [ ] Email verification on registration
- [ ] Password reset functionality
- [ ] Refresh token implementation
- [ ] Multi-factor authentication (2FA)
- [ ] Rate limiting on login attempts
- [ ] Token blacklist on logout
- [ ] OAuth2/Social login integration
- [ ] Session management

---

## Troubleshooting

### Issue: "Database connection error"
**Solution**: Ensure PostgreSQL is running and credentials in `.env` are correct.

### Issue: "JWT token generation failed"
**Solution**: Check that `JWT_SECRET` is set in environment variables.

### Issue: "Invalid email or password"
**Solution**: Verify credentials are correct and user exists in database.

### Issue: "Cannot bind request body"
**Solution**: Ensure request JSON matches the expected format and all required fields are provided.

---

## Support

For issues or questions, please create an issue in the repository or contact the development team.
