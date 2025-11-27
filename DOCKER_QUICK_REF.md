# 🐳 Docker Quick Reference

## 📚 Table of Contents
- [Installation](#-installation)
- [Quick Start](#-quick-start)
- [Essential Commands](#-essential-commands)
- [Helper Script](#-helper-script)
- [Troubleshooting](#-troubleshooting)
- [Common Issues](#-common-issues)

---

## 💻 Installation

### Windows
1. Download [Docker Desktop for Windows](https://www.docker.com/products/docker-desktop)
2. Run installer
3. Restart computer
4. Open PowerShell and verify:
   ```powershell
   docker --version
   docker-compose --version
   ```

### macOS
1. Download [Docker Desktop for Mac](https://www.docker.com/products/docker-desktop)
2. Open .dmg and drag Docker to Applications
3. Open Docker from Applications
4. Open terminal and verify:
   ```bash
   docker --version
   docker-compose --version
   ```

### Linux (Ubuntu/Debian)
```bash
# Install Docker
sudo apt-get update
sudo apt-get install docker.io docker-compose

# Verify installation
docker --version
docker-compose --version

# Add current user to docker group (optional, avoids sudo)
sudo usermod -aG docker $USER
```

---

## 🚀 Quick Start

### 1️⃣ Start Everything (One Command)
```bash
# Navigate to project directory
cd attendance-management

# Start all services
docker-compose up -d

# That's it! Services are running
```

### 2️⃣ Verify Services Are Running
```bash
docker-compose ps

# Expected output:
# NAME                    STATUS
# attendance-management-db   Up (healthy)
# attendance-management-app  Up
```

### 3️⃣ Test the API
```bash
curl -X POST http://localhost:2754/api/auth/register-student \
  -H "Content-Type: application/json" \
  -d '{"first_name":"John","last_name":"Doe","email":"john@example.com","password":"test123","matric_number":"STU001"}'

# Expected: 201 Created response
```

### 4️⃣ Stop Everything (One Command)
```bash
docker-compose down

# Or with data removal:
docker-compose down -v
```

---

## 🎯 Essential Commands

### Start/Stop Services

| Command | What it does |
|---------|------------|
| `docker-compose up` | Start services in foreground |
| `docker-compose up -d` | Start services in background |
| `docker-compose down` | Stop all services |
| `docker-compose restart` | Restart all services |
| `docker-compose stop` | Pause all services (keep data) |
| `docker-compose start` | Resume paused services |

### View Information

| Command | What it does |
|---------|------------|
| `docker-compose ps` | List all services and status |
| `docker-compose logs` | View all logs |
| `docker-compose logs -f` | Follow logs in real-time |
| `docker-compose logs app` | View app logs only |
| `docker-compose config` | Validate docker-compose.yml |
| `docker ps -a` | List all containers |
| `docker images` | List all images |

### Access Services

| Command | What it does |
|---------|------------|
| `docker-compose exec app sh` | Open shell in app |
| `docker-compose exec postgres psql -U postgres` | Open database shell |
| `docker-compose exec postgres pg_dump -U postgres -d attendance-management > backup.sql` | Backup database |

### Build & Maintenance

| Command | What it does |
|---------|------------|
| `docker-compose build` | Build images |
| `docker-compose up --build` | Rebuild and start |
| `docker-compose down -v` | Remove everything including data |
| `docker system prune` | Clean up unused resources |

---

## 🛠️ Helper Script

### Linux/macOS (Bash)
```bash
# Make script executable
chmod +x docker-helper.sh

# View available commands
./docker-helper.sh help

# Common commands
./docker-helper.sh up          # Start services
./docker-helper.sh down        # Stop services
./docker-helper.sh logs        # View logs
./docker-helper.sh test        # Test API
./docker-helper.sh shell-app   # Access app shell
./docker-helper.sh shell-db    # Access database
```

### Windows (Batch)
```cmd
REM View available commands
docker-helper.bat help

REM Common commands
docker-helper.bat up          # Start services
docker-helper.bat down        # Stop services
docker-helper.bat logs        # View logs
docker-helper.bat test        # Test API
docker-helper.bat shell-app   # Access app shell
docker-helper.bat shell-db    # Access database
```

---

## ✅ Verification Checklist

After running `docker-compose up -d`:

- [ ] Check status: `docker-compose ps` (should show 2 services UP)
- [ ] Check logs: `docker-compose logs` (should see successful startup)
- [ ] Test API:
  ```bash
  curl http://localhost:2754/api/auth/register-student \
    -H "Content-Type: application/json" \
    -d '{"first_name":"Test","last_name":"User","email":"test@example.com","password":"test123","matric_number":"TST001"}'
  ```
  Expected: `201 Created` response
- [ ] Check database:
  ```bash
  docker-compose exec postgres psql -U postgres -d attendance-management -c "SELECT * FROM students;"
  ```
  Should show your test user

---

## 🔧 Troubleshooting

### "docker: command not found"
**Cause:** Docker not installed or not in PATH
**Solution:**
1. Install Docker Desktop
2. Restart terminal
3. Verify: `docker --version`

### "Cannot connect to Docker daemon"
**Cause:** Docker daemon not running
**Solution:**
- **Windows:** Start Docker Desktop application
- **macOS:** Start Docker Desktop application
- **Linux:** `sudo systemctl start docker`

### "Permission denied while trying to connect to Docker daemon"
**Cause:** User not in docker group
**Solution (Linux only):**
```bash
sudo usermod -aG docker $USER
# Then logout and login again
```

### "Port 2754 already in use"
**Cause:** Another service using the port
**Solution:**
1. Find what's using the port:
   ```bash
   # Windows
   netstat -ano | findstr :2754
   
   # macOS/Linux
   lsof -i :2754
   ```
2. Either stop the other service or change the port in docker-compose.yml:
   ```yaml
   ports:
     - "2755:2754"  # Use port 2755 instead
   ```

### "Cannot find .env file"
**Cause:** .env file not at correct location
**Solution:**
```bash
# Verify .env location
ls cmd/api/app.env

# If missing, create it
cat > cmd/api/app.env << EOF
APP_PORT=:2754
DB_USER=postgres
DB_PASSWORD=Ac101bb101
DB_NAME=attendance-management
DB_PORT=5432
POOL_MAX_OPEN_CONN=5
POOL_MAX_IDLE_CONN=3
POOL_MAX_CONN_TIMEOUT=1m
JWT_SECRET=your-secret-key
EOF
```

### "Database connection refused"
**Cause:** PostgreSQL not ready or not running
**Solution:**
```bash
# Check if postgres service is healthy
docker-compose ps postgres

# View postgres logs
docker-compose logs postgres

# Wait a bit longer and try again
sleep 5
docker-compose ps
```

### "Build fails"
**Cause:** Various build issues
**Solution:**
```bash
# Clean rebuild
docker-compose down -v
docker-compose build --no-cache
docker-compose up

# Check detailed logs
docker-compose logs app
```

---

## 🆘 Common Issues

### Issue: Slow Build
**Solution:** Use .dockerignore file (already provided)

### Issue: Database Won't Connect
**Solution:**
1. Check postgres is healthy: `docker-compose ps`
2. Check logs: `docker-compose logs postgres`
3. Verify credentials match in .env and docker-compose.yml
4. Try rebuilding: `docker-compose down -v && docker-compose up`

### Issue: Can't Access API from Outside Container
**Cause:** Port not properly mapped
**Solution:**
```bash
# Verify port mapping
docker-compose ps

# Should show: 0.0.0.0:2754->2754/tcp

# Try connecting via localhost
curl http://localhost:2754/api/auth/...
```

### Issue: Data Disappears After Container Restart
**Cause:** Using wrong volume or no volume
**Solution:** Ensure docker-compose.yml has volume mapping:
```yaml
volumes:
  postgres_data:  # This persists the database
```

### Issue: Running Out of Disk Space
**Cause:** Docker layers, unused images, unused containers
**Solution:**
```bash
# Clean up unused resources
docker system prune

# Remove unused images
docker image prune

# Remove unused volumes
docker volume prune
```

---

## 📊 Useful Commands Cheat Sheet

```bash
# ===== Quick Start =====
docker-compose up -d              # Start all
docker-compose ps                 # Check status
docker-compose down               # Stop all

# ===== Logs & Debug =====
docker-compose logs               # View logs
docker-compose logs -f app         # Follow app logs
docker-compose logs postgres       # Database logs

# ===== Access Services =====
docker-compose exec app sh        # App shell
docker-compose exec postgres psql -U postgres -d attendance-management
# Run SQL commands in database

# ===== Build & Images =====
docker-compose build              # Build images
docker-compose build --no-cache   # Force rebuild
docker images                      # List images

# ===== Cleanup =====
docker-compose down -v            # Stop and remove data
docker system prune               # Clean unused

# ===== Testing =====
curl http://localhost:2754/api/auth/register-student \
  -H "Content-Type: application/json" \
  -d '{...}'

# ===== Backup & Restore =====
docker-compose exec -T postgres pg_dump -U postgres -d attendance-management > backup.sql
docker-compose exec -T postgres psql -U postgres -d attendance-management < backup.sql
```

---

## 🎓 Best Practices

1. **Always use `docker-compose.yml`** - Don't run containers manually
2. **Use `.env` for configuration** - Keep secrets out of compose file
3. **Tag images properly** - Include version numbers
4. **Use health checks** - Verify services are ready
5. **Set resource limits** - Prevent runaway containers
6. **Use named volumes** - For persistent data
7. **Keep images small** - Use alpine base images
8. **Update regularly** - Keep images and Docker updated

---

## 📞 When Things Go Wrong

1. **Check logs first:** `docker-compose logs`
2. **Verify services are running:** `docker-compose ps`
3. **Rebuild if needed:** `docker-compose down -v && docker-compose up`
4. **Check Docker is running:** `docker ps`
5. **Restart Docker daemon** (last resort)

---

## 🔗 Useful Links

- [Docker Documentation](https://docs.docker.com/)
- [Docker Compose Documentation](https://docs.docker.com/compose/)
- [Docker Hub](https://hub.docker.com/)
- [Docker Best Practices](https://docs.docker.com/develop/dev-best-practices/)

---

**Happy Dockering! 🐳**
