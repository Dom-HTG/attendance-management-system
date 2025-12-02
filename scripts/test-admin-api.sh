#!/bin/bash

# Admin API Test Script
# Comprehensive testing of all admin endpoints

set -e

API_URL="${API_URL:-http://localhost:2754}"

# Colors
GREEN='\033[0;32m'
BLUE='\033[0;34m'
YELLOW='\033[1;33m'
RED='\033[0;31m'
NC='\033[0m' # No Color

# Test counter
TESTS_RUN=0
TESTS_PASSED=0
TESTS_FAILED=0

echo -e "${BLUE}======================================${NC}"
echo -e "${BLUE}  Admin API Test Suite${NC}"
echo -e "${BLUE}======================================${NC}"
echo ""

# Helper function to test endpoints
test_endpoint() {
    local test_name="$1"
    local method="$2"
    local endpoint="$3"
    local expected_status="$4"
    local data="$5"
    
    TESTS_RUN=$((TESTS_RUN + 1))
    
    if [ -n "$data" ]; then
        response=$(curl -s -w "\n%{http_code}" -X "$method" "${API_URL}${endpoint}" \
            -H "Authorization: Bearer $ADMIN_TOKEN" \
            -H "Content-Type: application/json" \
            -d "$data")
    else
        response=$(curl -s -w "\n%{http_code}" -X "$method" "${API_URL}${endpoint}" \
            -H "Authorization: Bearer $ADMIN_TOKEN")
    fi
    
    status_code=$(echo "$response" | tail -n1)
    body=$(echo "$response" | sed '$d')
    
    if [ "$status_code" -eq "$expected_status" ]; then
        echo -e "${GREEN}✓${NC} $test_name (HTTP $status_code)"
        TESTS_PASSED=$((TESTS_PASSED + 1))
        return 0
    else
        echo -e "${RED}✗${NC} $test_name (Expected HTTP $expected_status, got $status_code)"
        TESTS_FAILED=$((TESTS_FAILED + 1))
        return 1
    fi
}

# Test 1: Admin Login
echo -e "${YELLOW}[1/15] Testing Admin Login${NC}"
LOGIN_RESPONSE=$(curl -s -X POST "${API_URL}/api/auth/login-admin" \
    -H "Content-Type: application/json" \
    -d '{"email":"admin@fupre.edu.ng","password":"Admin@2024"}')

ADMIN_TOKEN=$(echo "$LOGIN_RESPONSE" | python3 -c "import sys, json; print(json.load(sys.stdin)['data']['access_token'])" 2>/dev/null || echo "")

if [ -n "$ADMIN_TOKEN" ]; then
    echo -e "${GREEN}✓${NC} Admin login successful"
    TESTS_RUN=$((TESTS_RUN + 1))
    TESTS_PASSED=$((TESTS_PASSED + 1))
else
    echo -e "${RED}✗${NC} Admin login failed"
    TESTS_RUN=$((TESTS_RUN + 1))
    TESTS_FAILED=$((TESTS_FAILED + 1))
    echo ""
    echo -e "${RED}Cannot proceed without admin token. Exiting.${NC}"
    exit 1
fi
echo ""

# Test 2: Get All Students
echo -e "${YELLOW}[2/15] Testing GET /api/admin/students${NC}"
test_endpoint "Get all students with pagination" "GET" "/api/admin/students?page=1&limit=5" 200
echo ""

# Test 3: Get All Lecturers
echo -e "${YELLOW}[3/15] Testing GET /api/admin/lecturers${NC}"
test_endpoint "Get all lecturers with pagination" "GET" "/api/admin/lecturers?page=1&limit=5" 200
echo ""

# Test 4: Get Student Detail
echo -e "${YELLOW}[4/15] Testing GET /api/admin/users/student/:id${NC}"
test_endpoint "Get student details" "GET" "/api/admin/users/student/1" 200
echo ""

# Test 5: Get Lecturer Detail
echo -e "${YELLOW}[5/15] Testing GET /api/admin/users/lecturer/:id${NC}"
test_endpoint "Get lecturer details" "GET" "/api/admin/users/lecturer/1" 200
echo ""

# Test 6: Get All Events
echo -e "${YELLOW}[6/15] Testing GET /api/admin/events${NC}"
test_endpoint "Get all events with pagination" "GET" "/api/admin/events?page=1&limit=5" 200
echo ""

# Test 7: Get All Events (Filtered by Status)
echo -e "${YELLOW}[7/15] Testing GET /api/admin/events?status=expired${NC}"
test_endpoint "Get expired events" "GET" "/api/admin/events?status=expired&limit=3" 200
echo ""

# Test 8: Get Attendance Trends
echo -e "${YELLOW}[8/15] Testing GET /api/admin/trends${NC}"
test_endpoint "Get attendance trends" "GET" "/api/admin/trends?period=weekly" 200
echo ""

# Test 9: Get Low Attendance Students
echo -e "${YELLOW}[9/15] Testing GET /api/admin/low-attendance${NC}"
test_endpoint "Get low attendance students" "GET" "/api/admin/low-attendance?threshold=75&limit=10" 200
echo ""

# Test 10: Get System Settings
echo -e "${YELLOW}[10/15] Testing GET /api/admin/settings${NC}"
test_endpoint "Get system settings" "GET" "/api/admin/settings" 200
echo ""

# Test 11: Update System Settings
echo -e "${YELLOW}[11/15] Testing PATCH /api/admin/settings${NC}"
test_endpoint "Update system settings" "PATCH" "/api/admin/settings" 200 '{"low_attendance_threshold":70}'
echo ""

# Test 12: Get Audit Logs
echo -e "${YELLOW}[12/15] Testing GET /api/admin/audit-logs${NC}"
test_endpoint "Get audit logs" "GET" "/api/admin/audit-logs?page=1&limit=10" 200
echo ""

# Test 13: Search Students
echo -e "${YELLOW}[13/15] Testing GET /api/admin/students?search=test${NC}"
test_endpoint "Search students by keyword" "GET" "/api/admin/students?search=test&limit=5" 200
echo ""

# Test 14: Filter Lecturers by Department
echo -e "${YELLOW}[14/15] Testing GET /api/admin/lecturers?department=Computer${NC}"
test_endpoint "Filter lecturers by department" "GET" "/api/admin/lecturers?department=Computer%20Science&limit=5" 200
echo ""

# Test 15: Unauthorized Access (No Token)
echo -e "${YELLOW}[15/15] Testing Unauthorized Access${NC}"
UNAUTH_RESPONSE=$(curl -s -w "%{http_code}" -X GET "${API_URL}/api/admin/students" -o /dev/null)
if [ "$UNAUTH_RESPONSE" -eq 401 ]; then
    echo -e "${GREEN}✓${NC} Unauthorized access blocked (HTTP 401)"
    TESTS_RUN=$((TESTS_RUN + 1))
    TESTS_PASSED=$((TESTS_PASSED + 1))
else
    echo -e "${RED}✗${NC} Unauthorized access not properly blocked (HTTP $UNAUTH_RESPONSE)"
    TESTS_RUN=$((TESTS_RUN + 1))
    TESTS_FAILED=$((TESTS_FAILED + 1))
fi
echo ""

# Summary
echo -e "${BLUE}======================================${NC}"
echo -e "${BLUE}  Test Results${NC}"
echo -e "${BLUE}======================================${NC}"
echo -e "  Total Tests: $TESTS_RUN"
echo -e "  ${GREEN}Passed: $TESTS_PASSED${NC}"

if [ $TESTS_FAILED -gt 0 ]; then
    echo -e "  ${RED}Failed: $TESTS_FAILED${NC}"
    echo ""
    echo -e "${YELLOW}Some tests failed. Check the output above for details.${NC}"
    exit 1
else
    echo -e "  Failed: 0"
    echo ""
    echo -e "${GREEN}All tests passed! ✓${NC}"
    exit 0
fi
