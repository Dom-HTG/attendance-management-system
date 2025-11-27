# 🐳 Docker Setup Guide - Attendance Management System

## Overview

This guide explains how to use Docker and Docker Compose to run the entire Attendance Management System locally, including the PostgreSQL database and the Go application.

---

## 📋 Prerequisites

- **Docker** - [Install Docker Desktop](https://www.docker.com/products/docker-desktop)
- **Docker Compose** - Usually included with Docker Desktop
- **Your project files** - All files in the `attendance-management` directory

Verify installation:
```bash
docker --version
docker-compose --version
```

---

## 🚀 Quick Start (3 Steps)

### Step 1: Prepare Environment
```bash
# Navigate to project directory
cd attendance-management

# Verify your .env file exists (should be at cmd/api/app.env)
# Update it with your preferred settings if needed
cat cmd/api/app.env
```

**Default .env values:**
```
APP_PORT=:2754
DB_HOST=localhost
DB_USER=postgres
DB_PASSWORD=Ac101bb101
DB_NAME=attendance-management
DB_PORT=5432
POOL_MAX_OPEN_CONN=5
POOL_MAX_IDLE_CONN=3
POOL_MAX_CONN_TIMEOUT=1m
JWT_SECRET=your-super-secret-key-change-in-production-use-env-var
```

### Step 2: Build and Start Services
```bash
# Build and start all services (database + app)
docker-compose up --build

# Or run in background
docker-compose up -d --build
```

**Expected Output:**
```
Creating attendance-management-db  ... done
Creating attendance-management-app ... done
Attaching to attendance-management-db, attendance-management-app
...
app      | Database connection established successfully..
app      | [GIN-debug] Listening and serving HTTP on :2754
```

### Step 3: Test the API
```bash
# Register a student
curl -X POST http://localhost:2754/api/auth/register-student \
  -H "Content-Type: application/json" \
  -d '{
    "first_name": "John",
    "last_name": "Doe",
    "email": "john@example.com",
    "password": "TestPassword123",
    "matric_number": "STU-2024-001"
  }'

# Expected: 201 Created response
```

**✅ Done!** Your system is running in Docker!

---

## 📚 Docker Compose Components

### PostgreSQL Service (`postgres`)
- **Image:** postgres:15-alpine
- **Container Name:** attendance-management-db
- **Port:** 5432 (internal) → 5432 (host)
- **Volume:** postgres_data (persistent data)
- **Health Check:** Verifies database is ready

### Go Application Service (`app`)
- **Build:** Uses Dockerfile
- **Container Name:** attendance-management-app
- **Port:** 2754 (internal) → 2754 (host)
- **Depends On:** postgres service
- **Network:** attendance-network (internal communication)

### Network
- **Name:** attendance-network
- **Type:** Bridge
- **Purpose:** Allows services to communicate by service name

### Volume
- **Name:** postgres_data
- **Type:** Local
- **Purpose:** Persists database data between container restarts

---

## 🔧 Common Commands

### Start Services
```bash
# Start in foreground (see logs)
docker-compose up

# Start in background (detached mode)
docker-compose up -d

# Build and start
docker-compose up --build

# Build only (don't start)
docker-compose build
```

### Stop Services
```bash
# Stop all services (keep data)
docker-compose stop

# Stop and remove containers
docker-compose down

# Stop and remove everything including volumes
docker-compose down -v
```

### View Logs
```bash
# View all logs
docker-compose logs

# View app logs only
docker-compose logs app

# View database logs only
docker-compose logs postgres

# Follow logs in real-time
docker-compose logs -f

# View last 50 lines
docker-compose logs --tail=50
```

### Access Services
```bash
# Open shell in app container
docker-compose exec app sh

# Open shell in database container
docker-compose exec postgres sh

# Run database commands
docker-compose exec postgres psql -U postgres -d attendance-management -c "SELECT * FROM students;"
```

### Health Check
```bash
# See service status
docker-compose ps

# Check logs for errors
docker-compose logs app
```

---

## 🗄️ Database Management

### Access PostgreSQL
```bash
# Connect to database with psql
docker-compose exec postgres psql -U postgres -d attendance-management

# Once connected, useful commands:
\dt                    # List tables
\d students            # Describe students table
SELECT * FROM students;   # View all students
\q                     # Quit
```

### Backup Database
```bash
# Create backup file
docker-compose exec postgres pg_dump -U postgres -d attendance-management > backup.sql

# Restore backup
docker-compose exec -T postgres psql -U postgres -d attendance-management < backup.sql
```

### Reset Database
```bash
# Delete all data (but keep structure)
docker-compose exec postgres psql -U postgres -d attendance-management -c "
  TRUNCATE TABLE users CASCADE;
  TRUNCATE TABLE students CASCADE;
  TRUNCATE TABLE lecturers CASCADE;
"

# Delete database and volume completely
docker-compose down -v
docker-compose up -d
```

---

## 🐛 Troubleshooting

### Issue: "Port 2754 already in use"
**Solution:** Either stop the conflicting container or change the port in `docker-compose.yml`:
```yaml
ports:
  - "2755:2754"  # Change 2755 to your desired port
```

### Issue: "Port 5432 already in use"
**Solution:** Either stop the conflicting database or change the port:
```yaml
# In postgres service
ports:
  - "5433:5432"  # Change 5433 to your desired port
```

### Issue: "Database connection failed"
**Solution:** Check if postgres service is healthy:
```bash
docker-compose ps
# Status should show "healthy" for postgres service

docker-compose logs postgres
# Check for any error messages
```

### Issue: "Cannot connect to app on localhost:2754"
**Solution:**
```bash
# Check if app is running
docker-compose ps

# Check app logs
docker-compose logs app

# Verify port mapping
docker-compose ps app
```

### Issue: "Cannot find .env file"
**Solution:** 
```bash
# Make sure .env is in cmd/api/ directory
ls -la cmd/api/app.env

# If missing, create it
touch cmd/api/app.env
# Then add required variables
```

### Issue: "Build fails"
**Solution:**
```bash
# Clean build
docker-compose build --no-cache

# Check logs
docker-compose logs

# Try rebuilding with verbose output
docker-compose build --verbose
```

---

## 📝 Docker Compose Environment Variables

The `docker-compose.yml` file reads variables from `.env` file. Make sure your `.env` file has:

```env
# Server Configuration
APP_PORT=:2754

# Database Configuration
DB_HOST=localhost          # Changed to "postgres" by docker-compose
DB_USER=postgres
DB_PASSWORD=Ac101bb101
DB_NAME=attendance-management
DB_PORT=5432

# Database Pooling
POOL_MAX_OPEN_CONN=5
POOL_MAX_IDLE_CONN=3
POOL_MAX_CONN_TIMEOUT=1m

# JWT Configuration
JWT_SECRET=your-super-secret-key-change-in-production
```

**Note:** The `docker-compose.yml` automatically overrides `DB_HOST` to `"postgres"` (the service name) so containers can communicate internally.

---

## 🔐 Security Notes

### For Development (What We Have)
✅ Database credentials in .env
✅ Simple JWT secret
✅ Volumes for local testing

### For Production
⚠️ Use Docker secrets instead of .env
⚠️ Use strong JWT_SECRET
⚠️ Implement environment-specific configs
⚠️ Use private container registry
⚠️ Enable SSL/TLS
⚠️ Use read-only volumes where possible

---

## 🧪 Testing with Docker

### Option 1: curl (from host machine)
```bash
# Everything works the same as before
curl -X POST http://localhost:2754/api/auth/register-student \
  -H "Content-Type: application/json" \
  -d '{...}'
```

### Option 2: curl (inside app container)
```bash
# Open shell in app container
docker-compose exec app sh

# Inside container, you can access the app at http://localhost:2754
curl -X POST http://localhost:2754/api/auth/register-student ...
```

### Option 3: Postman
- Import `postman_collection.json`
- Change base URL to `http://localhost:2754` (already the default)
- Run your tests

---

## 📊 Useful Docker Compose Commands Reference

| Command | Purpose |
|---------|---------|
| `docker-compose up` | Start all services |
| `docker-compose down` | Stop all services |
| `docker-compose ps` | List services and status |
| `docker-compose logs` | View service logs |
| `docker-compose exec <service> <cmd>` | Run command in service |
| `docker-compose build` | Build images |
| `docker-compose pull` | Pull base images |
| `docker-compose restart` | Restart services |
| `docker-compose config` | Validate compose file |

---

## 🎯 Development Workflow

### Daily Workflow

**Morning - Start Services:**
```bash
cd attendance-management
docker-compose up -d
# Services start in background
```

**During Day - Access Services:**
```bash
# View logs
docker-compose logs -f app

# Test API (in another terminal)
curl -X POST http://localhost:2754/api/auth/register-student ...

# Connect to database if needed
docker-compose exec postgres psql -U postgres -d attendance-management
```

**Evening - Stop Services:**
```bash
# Stop (keeps data)
docker-compose stop

# Or just close the terminal if running in foreground
# Ctrl+C
```

**Next Day - Resume:**
```bash
# Services still exist with all data
docker-compose start

# Or restart from scratch
docker-compose down -v
docker-compose up -d --build
```

---

## 📈 Scaling and Performance

### Increase Resources
Edit `docker-compose.yml` to add resource limits:
```yaml
services:
  app:
    # ... other config ...
    deploy:
      resources:
        limits:
          cpus: '1'
          memory: 512M
        reservations:
          cpus: '0.5'
          memory: 256M
```

### Multiple App Instances (Load Balancing)
```yaml
services:
  app1:
    # ... app config ...
    ports:
      - "2754:2754"
  
  app2:
    # ... app config ...
    ports:
      - "2755:2754"
```

### Database Optimization
```yaml
services:
  postgres:
    # ... postgres config ...
    environment:
      # Increase shared buffers
      POSTGRES_INITDB_ARGS: "-c shared_buffers=256MB -c max_connections=200"
```

---

## 🔍 Inspect Docker Resources

### List All Containers
```bash
docker ps -a
```

### List All Images
```bash
docker images
```

### View Container Details
```bash
docker inspect attendance-management-app
docker inspect attendance-management-db
```

### View Network Details
```bash
docker network inspect attendance-network
```

### View Volume Details
```bash
docker volume inspect attendance-management-db_postgres_data
```

---

## 🗑️ Cleanup

### Remove Containers Only
```bash
docker-compose down
```

### Remove Containers and Volumes
```bash
docker-compose down -v
```

### Remove Images Too
```bash
docker-compose down --rmi all
```

### Remove Unused Docker Resources
```bash
docker system prune
docker system prune -a  # Remove unused images too
```

---

## 📚 Files Reference

- **`Dockerfile`** - Defines how to build the app image
- **`docker-compose.yml`** - Defines all services and how they connect
- **`cmd/api/app.env`** - Environment variables for configuration

---

## ✅ Verification Checklist

- [ ] Docker installed: `docker --version`
- [ ] Docker Compose installed: `docker-compose --version`
- [ ] .env file exists at `cmd/api/app.env`
- [ ] Dockerfile exists in project root
- [ ] docker-compose.yml exists in project root
- [ ] Run `docker-compose up -d`
- [ ] Check status: `docker-compose ps`
- [ ] Both services show "healthy" or "Up"
- [ ] Test API: `curl http://localhost:2754/api/auth/register-student`
- [ ] See success response (201 or validation error)

---

## 🎓 Next Steps

1. **Understand the Architecture** - Read how services communicate
2. **Test Different Scenarios** - Try all API endpoints
3. **Explore Logs** - Use `docker-compose logs` to understand flow
4. **Database Management** - Connect to postgres and explore schema
5. **Customize Configuration** - Modify .env and rebuild

---

## 📞 Quick Reference

**Start everything:**
```bash
docker-compose up -d
```

**See what's running:**
```bash
docker-compose ps
```

**Watch logs:**
```bash
docker-compose logs -f
```

**Stop everything:**
```bash
docker-compose down
```

**Access app shell:**
```bash
docker-compose exec app sh
```

**Access database:**
```bash
docker-compose exec postgres psql -U postgres -d attendance-management
```

---

**Happy Dockering! 🐳**
