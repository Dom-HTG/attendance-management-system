#!/bin/bash

# Database Seeding Script for Attendance Management System
# Populates database with Nigerian students, lecturer, and attendance records

set -e

API_URL="${API_URL:-http://localhost:2754}"

# Colors for output
GREEN='\033[0;32m'
BLUE='\033[0;34m'
YELLOW='\033[1;33m'
RED='\033[0;31m'
NC='\033[0m' # No Color

echo -e "${BLUE}======================================${NC}"
echo -e "${BLUE}  Database Seeding Script${NC}"
echo -e "${BLUE}======================================${NC}"
echo ""

# Wait for API to be ready
echo -e "${YELLOW}Waiting for API to be ready...${NC}"
for i in {1..30}; do
    HTTP_CODE=$(curl -s -o /dev/null -w "%{http_code}" "${API_URL}/api/auth/login/student" -X POST -H "Content-Type: application/json" -d '{}' 2>/dev/null)
    if [ "$HTTP_CODE" != "000" ]; then
        echo -e "${GREEN}✓ API is ready${NC}"
        break
    fi
    if [ $i -eq 30 ]; then
        echo -e "${RED}✗ API failed to start${NC}"
        exit 1
    fi
    sleep 1
done
echo ""

# Arrays of Nigerian names
FIRST_NAMES=("Chukwuemeka" "Oluwaseun" "Adewale" "Chidinma" "Chiamaka" "Olumide" "Funmilayo" "Obiageli" "Tunde" "Yewande" "Chinedu" "Folake" "Ikenna" "Amara" "Babajide")
LAST_NAMES=("Okonkwo" "Adeyemi" "Nwosu" "Olagunju" "Eze" "Balogun" "Okafor" "Adeleke" "Chukwu" "Oyedepo" "Nnaji" "Akinyemi" "Obiora" "Taiwo" "Ugwu")
DEPARTMENTS=("Computer Science" "Electrical Engineering" "Mechanical Engineering" "Civil Engineering" "Computer Science")

# Store login credentials
LOGIN_CREDENTIALS=()

echo -e "${BLUE}=== Creating Admin User ===${NC}"
ADMIN_EMAIL="admin@fupre.edu.ng"
ADMIN_PASSWORD="Admin@2024"

# Admin user is created via database migration, not API registration
# We'll just verify we can login
ADMIN_LOGIN=$(curl -s -X POST "${API_URL}/api/auth/login-admin" \
  -H "Content-Type: application/json" \
  -d "{
    \"email\": \"${ADMIN_EMAIL}\",
    \"password\": \"${ADMIN_PASSWORD}\"
  }")

ADMIN_TOKEN=$(echo "$ADMIN_LOGIN" | python3 -c "import sys, json; print(json.load(sys.stdin)['data']['access_token'])" 2>/dev/null || echo "")

if [ -z "$ADMIN_TOKEN" ]; then
    echo -e "${YELLOW}⚠ Admin user needs to be created via migration${NC}"
    echo -e "${YELLOW}  Run the migration: migrations/add_admin_functionality.sql${NC}"
else
    echo -e "${GREEN}✓ Admin user verified${NC}"
    echo -e "  Email: ${ADMIN_EMAIL}"
    echo -e "  Password: ${ADMIN_PASSWORD}"
    LOGIN_CREDENTIALS+=("Admin|${ADMIN_EMAIL}|${ADMIN_PASSWORD}")
fi
echo ""

echo -e "${BLUE}=== Registering Lecturer ===${NC}"
LECTURER_EMAIL="dr.adebayo.olumide@fupre.edu.ng"
LECTURER_PASSWORD="Lecturer@123"

LECTURER_RESPONSE=$(curl -s -X POST "${API_URL}/api/auth/register-lecturer" \
    -H "Content-Type: application/json" \
    -d "{
        \"email\": \"${LECTURER_EMAIL}\",
        \"password\": \"${LECTURER_PASSWORD}\",
        \"first_name\": \"Adebayo\",
        \"last_name\": \"Olumide\",
        \"staff_id\": \"FUPRE/LECT/001\",
        \"department\": \"Computer Science\"
    }")

if echo "$LECTURER_RESPONSE" | grep -q '"success":true'; then
    echo -e "${GREEN}✓ Lecturer registered: Dr. Adebayo Olumide${NC}"
    echo -e "  Email: ${LECTURER_EMAIL}"
    echo -e "  Password: ${LECTURER_PASSWORD}"
    echo -e "  Staff ID: FUPRE/LECT/001"
    LOGIN_CREDENTIALS+=("LECTURER|Dr. Adebayo Olumide|${LECTURER_EMAIL}|${LECTURER_PASSWORD}|FUPRE/LECT/001")
else
    echo -e "${RED}✗ Failed to register lecturer${NC}"
    echo "$LECTURER_RESPONSE"
fi
echo ""

# Login lecturer to get token
LECTURER_LOGIN=$(curl -s -X POST "${API_URL}/api/auth/login-lecturer" \
    -H "Content-Type: application/json" \
    -d "{
        \"email\": \"${LECTURER_EMAIL}\",
        \"password\": \"${LECTURER_PASSWORD}\"
    }")

LECTURER_TOKEN=$(echo "$LECTURER_LOGIN" | grep -o '"access_token":"[^"]*"' | cut -d'"' -f4)

if [ -z "$LECTURER_TOKEN" ]; then
    echo -e "${RED}✗ Failed to login lecturer${NC}"
    exit 1
fi

echo -e "${GREEN}✓ Lecturer logged in successfully${NC}"
echo ""

echo -e "${BLUE}=== Registering 15 Students ===${NC}"
STUDENT_IDS=()

for i in {0..14}; do
    FIRST_NAME="${FIRST_NAMES[$i]}"
    LAST_NAME="${LAST_NAMES[$i]}"
    MATRIC="FUPRE/$(printf "%04d" $((2021 + i % 3)))/$(printf "%05d" $((10000 + i)))"
    EMAIL="${FIRST_NAME,,}.${LAST_NAME,,}@fupre.edu.ng"
    PASSWORD="Student@$((100 + i))"
    DEPARTMENT="${DEPARTMENTS[$((i % 5))]}"
    LEVEL=$((100 * ((i % 4) + 1)))
    
    STUDENT_RESPONSE=$(curl -s -X POST "${API_URL}/api/auth/register-student" \
        -H "Content-Type: application/json" \
        -d "{
            \"email\": \"${EMAIL}\",
            \"password\": \"${PASSWORD}\",
            \"first_name\": \"${FIRST_NAME}\",
            \"last_name\": \"${LAST_NAME}\",
            \"matric_number\": \"${MATRIC}\"
        }")
    
    if echo "$STUDENT_RESPONSE" | grep -q '"success":true'; then
        STUDENT_ID=$(echo "$STUDENT_RESPONSE" | grep -o '"id":[0-9]*' | head -1 | cut -d':' -f2)
        STUDENT_IDS+=("$STUDENT_ID")
        echo -e "${GREEN}✓ Student $((i+1)): ${FIRST_NAME} ${LAST_NAME}${NC}"
        echo -e "  Email: ${EMAIL}"
        echo -e "  Password: ${PASSWORD}"
        echo -e "  Matric: ${MATRIC}"
        LOGIN_CREDENTIALS+=("STUDENT|${FIRST_NAME} ${LAST_NAME}|${EMAIL}|${PASSWORD}|${MATRIC}")
    else
        echo -e "${RED}✗ Failed to register student: ${FIRST_NAME} ${LAST_NAME}${NC}"
    fi
done
echo ""

echo -e "${BLUE}=== Creating Course Events ===${NC}"
COURSES_NAMES=("Introduction to Computer Science" "Data Structures and Algorithms" "Database Systems" "Software Engineering" "Machine Learning")
COURSES_CODES=("CSC101" "CSC201" "CSC301" "CSC401" "CSC501")
EVENT_IDS=()

for i in "${!COURSES_NAMES[@]}"; do
    COURSE_NAME="${COURSES_NAMES[$i]}"
    COURSE_CODE="${COURSES_CODES[$i]}"
    
    # Create events that are currently active (started 30 mins ago, ending in 1.5 hours)
    START_DATE=$(date -u -d "30 minutes ago" +"%Y-%m-%dT%H:%M:%SZ" 2>/dev/null || date -u -v-30M +"%Y-%m-%dT%H:%M:%SZ")
    END_DATE=$(date -u -d "90 minutes" +"%Y-%m-%dT%H:%M:%SZ" 2>/dev/null || date -u -v+90M +"%Y-%m-%dT%H:%M:%SZ")
    
    QR_RESPONSE=$(curl -s -X POST "${API_URL}/api/lecturer/qrcode/generate" \
        -H "Content-Type: application/json" \
        -H "Authorization: Bearer ${LECTURER_TOKEN}" \
        -d "{
            \"course_name\": \"${COURSE_NAME}\",
            \"course_code\": \"${COURSE_CODE}\",
            \"start_time\": \"${START_DATE}\",
            \"end_time\": \"${END_DATE}\",
            \"venue\": \"Lecture Hall $(((RANDOM % 5) + 1))\",
            \"department\": \"Computer Science\"
        }")
    
    if echo "$QR_RESPONSE" | grep -q '"event_id"'; then
        EVENT_ID=$(echo "$QR_RESPONSE" | grep -o '"event_id":[0-9]*' | cut -d':' -f2)
        QR_TOKEN=$(echo "$QR_RESPONSE" | grep -o '"qr_token":"[^"]*"' | cut -d'"' -f4)
        EVENT_IDS+=("${EVENT_ID}:${QR_TOKEN}:${START_DATE}:${COURSE_NAME}")
        echo -e "${GREEN}✓ Event created: ${COURSE_NAME} (${COURSE_CODE})${NC}"
    else
        echo -e "${RED}✗ Failed to create event: ${COURSE_NAME} (${COURSE_CODE})${NC}"
    fi
done
echo ""

echo -e "${BLUE}=== Simulating Attendance Records ===${NC}"
ATTENDANCE_COUNT=0
EVENT_COUNTER=0

# Debug: Show how many credentials we have
STUDENT_CRED_COUNT=0
for CRED in "${LOGIN_CREDENTIALS[@]}"; do
    if [[ $CRED == STUDENT* ]]; then
        STUDENT_CRED_COUNT=$((STUDENT_CRED_COUNT + 1))
    fi
done
echo -e "${YELLOW}Found ${STUDENT_CRED_COUNT} student credentials${NC}"
echo ""

for EVENT_DATA in "${EVENT_IDS[@]}"; do
    IFS=':' read -r EVENT_ID QR_TOKEN START_DATE COURSE_NAME <<< "$EVENT_DATA"
    EVENT_COUNTER=$((EVENT_COUNTER + 1))
    
    # Random number of students attending (60-100% attendance)
    STUDENT_CREDS=()
    for CRED in "${LOGIN_CREDENTIALS[@]}"; do
        if [[ $CRED == STUDENT* ]]; then
            STUDENT_CREDS+=("$CRED")
        fi
    done
    
    NUM_TOTAL=${#STUDENT_CREDS[@]}
    
    # Ensure at least some students if credentials exist
    if [ "$NUM_TOTAL" -eq 0 ]; then
        echo -e "${YELLOW}⚠ No student credentials found, skipping event ${EVENT_ID}${NC}"
        continue
    fi
    
    NUM_ATTENDING=$(( (NUM_TOTAL * (60 + RANDOM % 41)) / 100 ))
    [ "$NUM_ATTENDING" -lt 1 ] && NUM_ATTENDING=1
    
    # Shuffle students (preserve full lines, not split on spaces)
    SHUFFLED_CREDS=()
    while IFS= read -r line; do
        SHUFFLED_CREDS+=("$line")
    done < <(printf '%s\n' "${STUDENT_CREDS[@]}" | shuf)
    
    SUCCESSFUL_CHECKINS=0
    for ((j=0; j<NUM_ATTENDING; j++)); do
        STUDENT_CRED="${SHUFFLED_CREDS[$j]}"
        IFS='|' read -r ROLE NAME EMAIL PASSWORD MATRIC <<< "$STUDENT_CRED"
        
        # Skip if credentials are empty
        if [ -z "$EMAIL" ] || [ -z "$PASSWORD" ]; then
            continue
        fi
        
        # Login student
        STUDENT_LOGIN=$(curl -s -X POST "${API_URL}/api/auth/login-student" \
            -H "Content-Type: application/json" \
            -d "{
                \"email\": \"${EMAIL}\",
                \"password\": \"${PASSWORD}\"
            }")
        
        STUDENT_TOKEN=$(echo "$STUDENT_LOGIN" | grep -o '"access_token":"[^"]*"' | cut -d'"' -f4)
        
        if [ -n "$STUDENT_TOKEN" ]; then
            # Check-in
            CHECK_IN_RESPONSE=$(curl -s -X POST "${API_URL}/api/attendance/check-in" \
                -H "Content-Type: application/json" \
                -H "Authorization: Bearer ${STUDENT_TOKEN}" \
                -d "{
                    \"qr_token\": \"${QR_TOKEN}\"
                }")
            
            # Check for success - look for message or student_id field
            if echo "$CHECK_IN_RESPONSE" | grep -qE '"(message|student_id)":'; then
                if ! echo "$CHECK_IN_RESPONSE" | grep -q '"error"'; then
                    ATTENDANCE_COUNT=$((ATTENDANCE_COUNT + 1))
                    SUCCESSFUL_CHECKINS=$((SUCCESSFUL_CHECKINS + 1))
                fi
            fi
        fi
    done
    
    echo -e "${GREEN}✓ Event ${EVENT_COUNTER}/5: Recorded ${SUCCESSFUL_CHECKINS}/${NUM_ATTENDING} attendance(s) for ${COURSE_NAME:-Event} (ID: ${EVENT_ID})${NC}"
done

echo ""
echo -e "${GREEN}Total attendance records created: ${ATTENDANCE_COUNT}${NC}"
echo ""

# Save login credentials to file
CREDENTIALS_FILE="seed-login-credentials.txt"
echo "# Login Credentials - Generated on $(date)" > "$CREDENTIALS_FILE"
echo "# API URL: ${API_URL}" >> "$CREDENTIALS_FILE"
echo "" >> "$CREDENTIALS_FILE"

echo "=== ADMIN ===" >> "$CREDENTIALS_FILE"
for CRED in "${LOGIN_CREDENTIALS[@]}"; do
    if [[ $CRED == Admin* ]]; then
        IFS='|' read -r ROLE EMAIL PASSWORD <<< "$CRED"
        echo "Name: System Administrator" >> "$CREDENTIALS_FILE"
        echo "Email: $EMAIL" >> "$CREDENTIALS_FILE"
        echo "Password: $PASSWORD" >> "$CREDENTIALS_FILE"
        echo "Role: admin" >> "$CREDENTIALS_FILE"
        echo "" >> "$CREDENTIALS_FILE"
    fi
done

echo "=== LECTURER ===" >> "$CREDENTIALS_FILE"
for CRED in "${LOGIN_CREDENTIALS[@]}"; do
    if [[ $CRED == LECTURER* ]]; then
        IFS='|' read -r ROLE NAME EMAIL PASSWORD STAFF_ID <<< "$CRED"
        echo "Name: $NAME" >> "$CREDENTIALS_FILE"
        echo "Email: $EMAIL" >> "$CREDENTIALS_FILE"
        echo "Password: $PASSWORD" >> "$CREDENTIALS_FILE"
        echo "Staff ID: $STAFF_ID" >> "$CREDENTIALS_FILE"
        echo "" >> "$CREDENTIALS_FILE"
    fi
done

echo "=== STUDENTS ===" >> "$CREDENTIALS_FILE"
for CRED in "${LOGIN_CREDENTIALS[@]}"; do
    if [[ $CRED == STUDENT* ]]; then
        IFS='|' read -r ROLE NAME EMAIL PASSWORD MATRIC <<< "$CRED"
        echo "Name: $NAME" >> "$CREDENTIALS_FILE"
        echo "Email: $EMAIL" >> "$CREDENTIALS_FILE"
        echo "Password: $PASSWORD" >> "$CREDENTIALS_FILE"
        echo "Matric: $MATRIC" >> "$CREDENTIALS_FILE"
        echo "" >> "$CREDENTIALS_FILE"
    fi
done

echo -e "${BLUE}======================================${NC}"
echo -e "${GREEN}✓ Database seeded successfully!${NC}"
echo -e "${BLUE}======================================${NC}"
echo ""
echo -e "${YELLOW}Login credentials saved to: ${CREDENTIALS_FILE}${NC}"
echo ""
echo -e "${BLUE}Summary:${NC}"
echo -e "  • 1 Admin user (via migration)"
echo -e "  • 1 Lecturer registered"
echo -e "  • 15 Students registered"
echo -e "  • ${#EVENT_IDS[@]} Course events created"
echo -e "  • ${ATTENDANCE_COUNT} Attendance records created"
echo ""
echo -e "${GREEN}You can now login with any of the credentials in ${CREDENTIALS_FILE}${NC}"
