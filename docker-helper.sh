#!/bin/bash

# Attendance Management System - Docker Helper Script
# Usage: ./docker-helper.sh [command]

set -e

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Configuration
PROJECT_NAME="attendance-management"
COMPOSE_FILE="docker-compose.yml"

# Helper functions
print_header() {
    echo -e "${BLUE}═══════════════════════════════════════${NC}"
    echo -e "${BLUE}$1${NC}"
    echo -e "${BLUE}═══════════════════════════════════════${NC}"
}

print_success() {
    echo -e "${GREEN}✓ $1${NC}"
}

print_error() {
    echo -e "${RED}✗ $1${NC}"
}

print_warning() {
    echo -e "${YELLOW}⚠ $1${NC}"
}

# Check if docker-compose.yml exists
if [ ! -f "$COMPOSE_FILE" ]; then
    print_error "docker-compose.yml not found in current directory"
    exit 1
fi

# Command handling
case "${1:-help}" in
    up)
        print_header "Starting services..."
        docker-compose up -d
        sleep 3
        print_success "Services started"
        docker-compose ps
        ;;
    
    down)
        print_header "Stopping services..."
        docker-compose down
        print_success "Services stopped"
        ;;
    
    restart)
        print_header "Restarting services..."
        docker-compose restart
        print_success "Services restarted"
        docker-compose ps
        ;;
    
    logs)
        print_header "Viewing logs (Ctrl+C to exit)..."
        docker-compose logs -f
        ;;
    
    logs-app)
        print_header "Viewing app logs..."
        docker-compose logs -f app
        ;;
    
    logs-db)
        print_header "Viewing database logs..."
        docker-compose logs -f postgres
        ;;
    
    status)
        print_header "Service status"
        docker-compose ps
        ;;
    
    clean)
        print_header "Cleaning up (removing containers and volumes)..."
        print_warning "This will delete all data!"
        read -p "Are you sure? (y/N) " -n 1 -r
        echo
        if [[ $REPLY =~ ^[Yy]$ ]]; then
            docker-compose down -v
            print_success "Cleaned up"
        else
            print_warning "Cleanup cancelled"
        fi
        ;;
    
    build)
        print_header "Building services..."
        docker-compose build --no-cache
        print_success "Build complete"
        ;;
    
    rebuild)
        print_header "Rebuilding and starting..."
        docker-compose up -d --build
        sleep 3
        print_success "Rebuild and start complete"
        docker-compose ps
        ;;
    
    shell-app)
        print_header "Opening shell in app container..."
        docker-compose exec app sh
        ;;
    
    shell-db)
        print_header "Opening PostgreSQL shell..."
        docker-compose exec postgres psql -U postgres -d attendance-management
        ;;
    
    test)
        print_header "Testing API endpoints..."
        
        # Register student
        print_warning "Testing student registration..."
        REGISTER_RESPONSE=$(curl -s -X POST http://localhost:2754/api/auth/register-student \
            -H "Content-Type: application/json" \
            -d '{
                "first_name": "Test",
                "last_name": "User",
                "email": "test@example.com",
                "password": "TestPassword123",
                "matric_number": "TEST-2024-001"
            }')
        
        if echo "$REGISTER_RESPONSE" | grep -q "successfully registered"; then
            print_success "Student registration successful"
        else
            print_error "Student registration failed"
            echo "$REGISTER_RESPONSE"
        fi
        
        # Login student
        print_warning "Testing student login..."
        LOGIN_RESPONSE=$(curl -s -X POST http://localhost:2754/api/auth/login-student \
            -H "Content-Type: application/json" \
            -d '{
                "email": "test@example.com",
                "password": "TestPassword123"
            }')
        
        if echo "$LOGIN_RESPONSE" | grep -q "access_token"; then
            print_success "Student login successful"
            TOKEN=$(echo "$LOGIN_RESPONSE" | grep -o '"access_token":"[^"]*' | cut -d'"' -f4)
            echo "Token: ${TOKEN:0:50}..."
        else
            print_error "Student login failed"
            echo "$LOGIN_RESPONSE"
        fi
        ;;
    
    backup)
        print_header "Backing up database..."
        BACKUP_FILE="backup_$(date +%Y%m%d_%H%M%S).sql"
        docker-compose exec -T postgres pg_dump -U postgres -d attendance-management > "$BACKUP_FILE"
        print_success "Database backed up to $BACKUP_FILE"
        ;;
    
    restore)
        if [ -z "$2" ]; then
            print_error "Please specify backup file: ./docker-helper.sh restore <backup_file>"
            exit 1
        fi
        
        if [ ! -f "$2" ]; then
            print_error "Backup file not found: $2"
            exit 1
        fi
        
        print_header "Restoring database from $2..."
        print_warning "This will overwrite existing data!"
        read -p "Are you sure? (y/N) " -n 1 -r
        echo
        if [[ $REPLY =~ ^[Yy]$ ]]; then
            docker-compose exec -T postgres psql -U postgres -d attendance-management < "$2"
            print_success "Database restored"
        else
            print_warning "Restore cancelled"
        fi
        ;;
    
    help|--help|-h)
        print_header "Docker Helper - Available Commands"
        echo
        echo -e "${GREEN}Service Management:${NC}"
        echo "  up              - Start all services in background"
        echo "  down            - Stop all services"
        echo "  restart         - Restart all services"
        echo "  status          - Show service status"
        echo "  logs            - View all logs (follow mode)"
        echo "  logs-app        - View app logs only"
        echo "  logs-db         - View database logs only"
        echo
        echo -e "${GREEN}Building:${NC}"
        echo "  build           - Build images"
        echo "  rebuild         - Rebuild and start services"
        echo
        echo -e "${GREEN}Access:${NC}"
        echo "  shell-app       - Open shell in app container"
        echo "  shell-db        - Open PostgreSQL shell"
        echo
        echo -e "${GREEN}Testing & Backup:${NC}"
        echo "  test            - Run API tests"
        echo "  backup          - Backup database"
        echo "  restore <file>  - Restore database from backup"
        echo
        echo -e "${GREEN}Cleanup:${NC}"
        echo "  clean           - Remove all containers and volumes"
        echo
        echo -e "${GREEN}Help:${NC}"
        echo "  help            - Show this message"
        echo
        ;;
    
    *)
        print_error "Unknown command: $1"
        echo "Run './docker-helper.sh help' for available commands"
        exit 1
        ;;
esac
