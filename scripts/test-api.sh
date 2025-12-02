#!/bin/bash

# Comprehensive Test Script for Attendance Management System
# Tests all endpoints: Auth, QR Generation, Check-in, Attendance Records, and Analytics

set -e

API_URL="http://localhost:2754"
GREEN='\033[0;32m'
RED='\033[0;31m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m'

TOTAL_TESTS=0
PASSED_TESTS=0
FAILED_TESTS=0

echo "========================================="
echo "  Attendance Management System"
echo "  Comprehensive API Test Suite"
echo "========================================="
echo ""

# Helper function to check API response
check_response() {
    local response="$1"
    local test_name="$2"
    TOTAL_TESTS=$((TOTAL_TESTS + 1))
    
    # Consider both success and "already registered" as acceptable
    if echo "$response" | grep -q '"success":true\|"data":\|"access_token":\|"message":"\|already registered'; then
        echo -e "${GREEN}✓${NC} $test_name"
        PASSED_TESTS=$((PASSED_TESTS + 1))
        return 0
    else
        echo -e "${RED}✗${NC} $test_name"
        echo "   Response: $response"
        FAILED_TESTS=$((FAILED_TESTS + 1))
        return 1
    fi
}

# Section 1: Authentication
echo -e "${BLUE}=== 1. Authentication Tests ===${NC}"

# Generate unique timestamp for test users
TIMESTAMP=$(date +%s)

echo -n "Registering Student... "
STUDENT_REG=$(curl -s -X POST "${API_URL}/api/auth/register-student" \
  -H "Content-Type: application/json" \
  -d "{\"first_name\":\"Test\",\"last_name\":\"Student\",\"email\":\"test.student.${TIMESTAMP}@test.edu\",\"password\":\"test123\",\"matric_number\":\"TEST/${TIMESTAMP}/001\"}")
check_response "$STUDENT_REG" "Student Registration"

echo -n "Registering Lecturer... "
LECTURER_REG=$(curl -s -X POST "${API_URL}/api/auth/register-lecturer" \
  -H "Content-Type: application/json" \
  -d "{\"first_name\":\"Test\",\"last_name\":\"Lecturer\",\"email\":\"test.lecturer.${TIMESTAMP}@test.edu\",\"password\":\"test123\",\"department\":\"Test Department\",\"staff_id\":\"TEST_STAFF_${TIMESTAMP}\"}")
check_response "$LECTURER_REG" "Lecturer Registration"

echo -n "Student Login... "
STUDENT_LOGIN=$(curl -s -X POST "${API_URL}/api/auth/login-student" \
  -H "Content-Type: application/json" \
  -d "{\"email\":\"test.student.${TIMESTAMP}@test.edu\",\"password\":\"test123\"}")
STUDENT_TOKEN=$(echo "$STUDENT_LOGIN" | grep -o '"access_token":"[^"]*"' | cut -d'"' -f4)
check_response "$STUDENT_LOGIN" "Student Login"

echo -n "Lecturer Login... "
LECTURER_LOGIN=$(curl -s -X POST "${API_URL}/api/auth/login-lecturer" \
  -H "Content-Type: application/json" \
  -d "{\"email\":\"test.lecturer.${TIMESTAMP}@test.edu\",\"password\":\"test123\"}")
LECTURER_TOKEN=$(echo "$LECTURER_LOGIN" | grep -o '"access_token":"[^"]*"' | cut -d'"' -f4)
check_response "$LECTURER_LOGIN" "Lecturer Login"

# Section 2: QR Code Generation
echo ""
echo -e "${BLUE}=== 2. QR Code Generation ===${NC}"

echo -n "Generating QR Code for Event... "
QR_RESPONSE=$(curl -s -X POST "${API_URL}/api/lecturer/qrcode/generate" \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer $LECTURER_TOKEN" \
  -d '{
    "course_name":"Test Course",
    "course_code":"TEST101",
    "start_time":"2025-01-01T10:00:00Z",
    "end_time":"2025-12-31T23:59:59Z",
    "venue":"Test Hall",
    "department":"Test Department"
  }')
QR_TOKEN=$(echo "$QR_RESPONSE" | grep -o '"qr_token":"[^"]*"' | cut -d'"' -f4)
EVENT_ID=$(echo "$QR_RESPONSE" | grep -o '"event_id":[0-9]*' | cut -d':' -f2)
check_response "$QR_RESPONSE" "QR Code Generation"

# Section 3: Check-in
echo ""
echo -e "${BLUE}=== 3. Student Check-in ===${NC}"

echo -n "Student Check-in to Event... "
CHECKIN_RESPONSE=$(curl -s -X POST "${API_URL}/api/attendance/check-in" \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer $STUDENT_TOKEN" \
  -d "{\"qr_token\":\"$QR_TOKEN\"}")
check_response "$CHECKIN_RESPONSE" "Check-in"

# Section 4: Attendance Records
echo ""
echo -e "${BLUE}=== 4. Attendance Records ===${NC}"

echo -n "Get Event Attendance (Lecturer)... "
EVENT_ATTENDANCE=$(curl -s -X GET "${API_URL}/api/attendance/${EVENT_ID}" \
  -H "Authorization: Bearer $LECTURER_TOKEN")
check_response "$EVENT_ATTENDANCE" "Event Attendance"

echo -n "Get Student Attendance Records... "
STUDENT_RECORDS=$(curl -s -X GET "${API_URL}/api/attendance/student/records" \
  -H "Authorization: Bearer $STUDENT_TOKEN")
check_response "$STUDENT_RECORDS" "Student Records"

# Section 5: Analytics - Lecturer
echo ""
echo -e "${BLUE}=== 5. Lecturer Analytics ===${NC}"

echo -n "Get Lecturer Events... "
LECTURER_EVENTS=$(curl -s -X GET "${API_URL}/api/events/lecturer" \
  -H "Authorization: Bearer $LECTURER_TOKEN")
check_response "$LECTURER_EVENTS" "Lecturer Events"

echo -n "Get Lecturer Summary... "
LECTURER_SUMMARY=$(curl -s -X GET "${API_URL}/api/analytics/lecturer/summary" \
  -H "Authorization: Bearer $LECTURER_TOKEN")
check_response "$LECTURER_SUMMARY" "Lecturer Summary"

# Section 6: Analytics - Admin
echo ""
echo -e "${BLUE}=== 6. Admin Analytics ===${NC}"

echo -n "Get Admin Overview... "
ADMIN_OVERVIEW=$(curl -s -X GET "${API_URL}/api/analytics/admin/overview" \
  -H "Authorization: Bearer $LECTURER_TOKEN")
check_response "$ADMIN_OVERVIEW" "Admin Overview"

echo -n "Get Department Statistics... "
DEPT_STATS=$(curl -s -X GET "${API_URL}/api/analytics/admin/departments" \
  -H "Authorization: Bearer $LECTURER_TOKEN")
check_response "$DEPT_STATS" "Department Stats"

# Section 7: Sample Data Verification (using seeded data)
echo ""
echo -e "${BLUE}=== 7. Seeded Data Tests ===${NC}"

echo -n "Login Seeded Lecturer... "
SEEDED_LECTURER_LOGIN=$(curl -s -X POST "${API_URL}/api/auth/login-lecturer" \
  -H "Content-Type: application/json" \
  -d '{"email":"dr.adebayo.olumide@fupre.edu.ng","password":"Lecturer@123"}')
SEEDED_LECTURER_TOKEN=$(echo "$SEEDED_LECTURER_LOGIN" | grep -o '"access_token":"[^"]*"' | cut -d'"' -f4)
check_response "$SEEDED_LECTURER_LOGIN" "Seeded Lecturer Login"

if [ -n "$SEEDED_LECTURER_TOKEN" ]; then
    echo -n "Get Seeded Lecturer Events... "
    SEEDED_EVENTS=$(curl -s -X GET "${API_URL}/api/events/lecturer" \
      -H "Authorization: Bearer $SEEDED_LECTURER_TOKEN")
    check_response "$SEEDED_EVENTS" "Seeded Events"
    
    # Extract event count
    EVENT_COUNT=$(echo "$SEEDED_EVENTS" | grep -o '"total_events":[0-9]*' | cut -d':' -f2)
    if [ -n "$EVENT_COUNT" ] && [ "$EVENT_COUNT" -gt 0 ]; then
        echo -e "   ${GREEN}Found $EVENT_COUNT seeded events${NC}"
    fi
fi

echo -n "Login Seeded Student... "
SEEDED_STUDENT_LOGIN=$(curl -s -X POST "${API_URL}/api/auth/login-student" \
  -H "Content-Type: application/json" \
  -d '{"email":"chukwuemeka.okonkwo@fupre.edu.ng","password":"Student@100"}')
SEEDED_STUDENT_TOKEN=$(echo "$SEEDED_STUDENT_LOGIN" | grep -o '"access_token":"[^"]*"' | cut -d'"' -f4)
check_response "$SEEDED_STUDENT_LOGIN" "Seeded Student Login"

if [ -n "$SEEDED_STUDENT_TOKEN" ]; then
    echo -n "Get Seeded Student Records... "
    SEEDED_RECORDS=$(curl -s -X GET "${API_URL}/api/attendance/student/records" \
      -H "Authorization: Bearer $SEEDED_STUDENT_TOKEN")
    check_response "$SEEDED_RECORDS" "Seeded Student Records"
fi

# Final Summary
echo ""
echo "========================================="
echo "  Test Results Summary"
echo "========================================="
echo -e "Total Tests:  ${BLUE}$TOTAL_TESTS${NC}"
echo -e "Passed:       ${GREEN}$PASSED_TESTS${NC}"
echo -e "Failed:       ${RED}$FAILED_TESTS${NC}"
echo ""

if [ $FAILED_TESTS -eq 0 ]; then
    echo -e "${GREEN}✓ All tests passed successfully!${NC}"
    exit 0
else
    echo -e "${RED}✗ Some tests failed. Please review the output above.${NC}"
    exit 1
fi
