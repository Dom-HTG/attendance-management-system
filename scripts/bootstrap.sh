#!/bin/bash

# ============================================
# Bootstrap Script
# Attendance Management System
# ============================================
# This script:
# 1. Starts the application with Docker Compose
# 2. Seeds the database with test data
# 3. Runs comprehensive API tests
# ============================================

set -e  # Exit on error

# Colors
GREEN='\033[0;32m'
BLUE='\033[0;34m'
YELLOW='\033[1;33m'
RED='\033[0;31m'
NC='\033[0m' # No Color

# Get the script directory
SCRIPT_DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )"
PROJECT_ROOT="$(dirname "$SCRIPT_DIR")"

echo -e "${BLUE}=========================================${NC}"
echo -e "${BLUE}  Attendance Management System${NC}"
echo -e "${BLUE}  Bootstrap Script${NC}"
echo -e "${BLUE}=========================================${NC}"
echo ""

# Step 1: Start Docker Compose
echo -e "${YELLOW}[1/3] Starting application with Docker Compose...${NC}"
cd "$PROJECT_ROOT"
docker compose up -d --build

echo -e "${GREEN}✓${NC} Application started"
echo -e "${YELLOW}Waiting 8 seconds for services to be ready...${NC}"
sleep 8
echo ""

# Step 2: Seed Admin User
echo -e "${YELLOW}[2/4] Creating admin user...${NC}"
bash "$SCRIPT_DIR/seed-admin.sh"

if [ $? -eq 0 ]; then
    echo -e "${GREEN}✓${NC} Admin user created"
else
    echo -e "${YELLOW}⚠${NC} Admin user creation skipped (may already exist)"
fi
echo ""

# Step 3: Seed Database
echo -e "${YELLOW}[3/4] Seeding database with test data...${NC}"
bash "$SCRIPT_DIR/seed-database.sh"

if [ $? -eq 0 ]; then
    echo -e "${GREEN}✓${NC} Database seeded successfully"
else
    echo -e "${RED}✗${NC} Database seeding failed"
    exit 1
fi
echo ""

# Step 4: Run API Tests
echo -e "${YELLOW}[4/4] Running comprehensive API tests...${NC}"
bash "$SCRIPT_DIR/test-api.sh"

if [ $? -eq 0 ]; then
    echo ""
    echo -e "${GREEN}=========================================${NC}"
    echo -e "${GREEN}  Bootstrap Complete!${NC}"
    echo -e "${GREEN}=========================================${NC}"
    echo ""
    echo -e "${BLUE}Application is running at:${NC} http://localhost:2754"
    echo -e "${BLUE}Login credentials saved in:${NC} seed-login-credentials.txt"
    echo -e "${BLUE}Admin credentials:${NC} admin@fupre.edu.ng / Admin@2024"
    echo ""
    echo -e "${BLUE}Quick commands:${NC}"
    echo -e "  View logs:        ${YELLOW}docker compose logs -f${NC}"
    echo -e "  Stop application: ${YELLOW}docker compose down${NC}"
    echo -e "  Restart:          ${YELLOW}docker compose restart${NC}"
    echo -e "  Run tests only:   ${YELLOW}bash scripts/test-api.sh${NC}"
    echo ""
else
    echo -e "${RED}✗${NC} API tests failed"
    exit 1
fi
