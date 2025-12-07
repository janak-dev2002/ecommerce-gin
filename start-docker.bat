@echo off
REM Quick Start Script for Docker Deployment (Windows)

echo.
echo E-Commerce Backend - Docker Quick Start
echo ==========================================
echo.

REM Check if Docker is installed
docker --version >nul 2>&1
if errorlevel 1 (
    echo X Docker is not installed!
    echo Please install Docker Desktop from: https://www.docker.com/products/docker-desktop
    pause
    exit /b 1
)

REM Check if Docker Compose is installed
docker-compose --version >nul 2>&1
if errorlevel 1 (
    echo X Docker Compose is not installed!
    echo Please install Docker Compose
    pause
    exit /b 1
)

echo √ Docker is installed
echo √ Docker Compose is installed
echo.

REM Check if .env exists
if not exist .env (
    echo Creating .env file from template...
    copy .env.docker .env
    echo √ .env file created
    echo.
    echo IMPORTANT: Edit .env file and add your S3 credentials!
    echo You can skip S3 config for testing, but image upload won't work.
    echo.
    pause
)

echo.
echo Building and starting services...
echo This may take a few minutes on first run...
echo.

REM Build and start services
docker-compose up --build -d

echo.
echo Waiting for services to be ready...
timeout /t 15 /nobreak >nul

echo.
echo √ All services should be running!
echo.
echo Your application is available at:
echo    - API:              http://localhost:8080
echo    - Swagger UI:       http://localhost:8080/swagger/index.html
echo    - Health Check:     http://localhost:8080/health
echo    - phpMyAdmin:       http://localhost:8081
echo    - Redis Commander:  http://localhost:8082
echo.
echo View logs:
echo    docker-compose logs -f
echo.
echo Stop services:
echo    docker-compose down
echo.
echo Happy testing!
echo.
pause
