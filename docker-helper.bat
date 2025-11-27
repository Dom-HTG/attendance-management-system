@echo off
REM Attendance Management System - Docker Helper Script (Windows)
REM Usage: docker-helper.bat [command]

setlocal enabledelayedexpansion

REM Configuration
set PROJECT_NAME=attendance-management
set COMPOSE_FILE=docker-compose.yml

REM Check if docker-compose.yml exists
if not exist "%COMPOSE_FILE%" (
    echo.
    echo [ERROR] docker-compose.yml not found in current directory
    echo.
    exit /b 1
)

REM Command handling
if "%1"=="" goto help
if "%1"=="up" goto up
if "%1"=="down" goto down
if "%1"=="restart" goto restart
if "%1"=="logs" goto logs
if "%1"=="logs-app" goto logs_app
if "%1"=="logs-db" goto logs_db
if "%1"=="status" goto status
if "%1"=="clean" goto clean
if "%1"=="build" goto build
if "%1"=="rebuild" goto rebuild
if "%1"=="shell-app" goto shell_app
if "%1"=="shell-db" goto shell_db
if "%1"=="test" goto test
if "%1"=="backup" goto backup
if "%1"=="help" goto help
if "%1"=="--help" goto help
if "%1"=="-h" goto help

echo [ERROR] Unknown command: %1
echo Run 'docker-helper.bat help' for available commands
exit /b 1

:up
echo.
echo ===== Starting services =====
echo.
docker-compose up -d
timeout /t 3
docker-compose ps
goto end

:down
echo.
echo ===== Stopping services =====
echo.
docker-compose down
goto end

:restart
echo.
echo ===== Restarting services =====
echo.
docker-compose restart
docker-compose ps
goto end

:logs
echo.
echo ===== Viewing logs (Ctrl+C to exit) =====
echo.
docker-compose logs -f
goto end

:logs_app
echo.
echo ===== Viewing app logs =====
echo.
docker-compose logs -f app
goto end

:logs_db
echo.
echo ===== Viewing database logs =====
echo.
docker-compose logs -f postgres
goto end

:status
echo.
echo ===== Service status =====
echo.
docker-compose ps
goto end

:clean
echo.
echo ===== Cleaning up (removing containers and volumes) =====
echo [WARNING] This will delete all data!
echo.
set /p confirm="Are you sure? (y/N): "
if /i "%confirm%"=="y" (
    docker-compose down -v
    echo [SUCCESS] Cleaned up
) else (
    echo [WARNING] Cleanup cancelled
)
goto end

:build
echo.
echo ===== Building services =====
echo.
docker-compose build --no-cache
echo [SUCCESS] Build complete
goto end

:rebuild
echo.
echo ===== Rebuilding and starting =====
echo.
docker-compose up -d --build
timeout /t 3
echo [SUCCESS] Rebuild and start complete
docker-compose ps
goto end

:shell_app
echo.
echo ===== Opening shell in app container =====
echo.
docker-compose exec app sh
goto end

:shell_db
echo.
echo ===== Opening PostgreSQL shell =====
echo.
docker-compose exec postgres psql -U postgres -d attendance-management
goto end

:test
echo.
echo ===== Testing API endpoints =====
echo.

REM Test student registration
echo [TEST] Testing student registration...
for /f "delims=" %%i in ('curl -s -X POST http://localhost:2754/api/auth/register-student -H "Content-Type: application/json" -d "{\"first_name\":\"Test\",\"last_name\":\"User\",\"email\":\"test@example.com\",\"password\":\"TestPassword123\",\"matric_number\":\"TEST-2024-001\"}"') do set REGISTER_RESPONSE=%%i

if "!REGISTER_RESPONSE!"=="" (
    echo [ERROR] No response from server. Is app running?
    goto end
)

echo !REGISTER_RESPONSE! | find "successfully registered" >nul
if %errorlevel%==0 (
    echo [SUCCESS] Student registration successful
) else (
    echo [ERROR] Student registration failed
    echo Response: !REGISTER_RESPONSE!
    goto end
)

REM Test student login
echo [TEST] Testing student login...
for /f "delims=" %%i in ('curl -s -X POST http://localhost:2754/api/auth/login-student -H "Content-Type: application/json" -d "{\"email\":\"test@example.com\",\"password\":\"TestPassword123\"}"') do set LOGIN_RESPONSE=%%i

echo !LOGIN_RESPONSE! | find "access_token" >nul
if %errorlevel%==0 (
    echo [SUCCESS] Student login successful
    echo Response contains access token
) else (
    echo [ERROR] Student login failed
    echo Response: !LOGIN_RESPONSE!
)
goto end

:backup
echo.
echo ===== Backing up database =====
echo.
for /f "tokens=2-4 delims=/ " %%a in ('date /t') do (set mydate=%%c%%a%%b)
for /f "tokens=1-2 delims=/:" %%a in ('time /t') do (set mytime=%%a%%b)
set BACKUP_FILE=backup_%mydate%_%mytime%.sql
docker-compose exec -T postgres pg_dump -U postgres -d attendance-management > "%BACKUP_FILE%"
if %errorlevel%==0 (
    echo [SUCCESS] Database backed up to %BACKUP_FILE%
) else (
    echo [ERROR] Backup failed
)
goto end

:help
echo.
echo ===== Docker Helper - Available Commands =====
echo.
echo Service Management:
echo   up              - Start all services in background
echo   down            - Stop all services
echo   restart         - Restart all services
echo   status          - Show service status
echo   logs            - View all logs ^(follow mode^)
echo   logs-app        - View app logs only
echo   logs-db         - View database logs only
echo.
echo Building:
echo   build           - Build images
echo   rebuild         - Rebuild and start services
echo.
echo Access:
echo   shell-app       - Open shell in app container
echo   shell-db        - Open PostgreSQL shell
echo.
echo Testing ^& Backup:
echo   test            - Run API tests
echo   backup          - Backup database
echo.
echo Cleanup:
echo   clean           - Remove all containers and volumes
echo.
echo Help:
echo   help            - Show this message
echo.

:end
exit /b 0
