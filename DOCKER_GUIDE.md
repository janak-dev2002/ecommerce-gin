# Docker Guide for E-Commerce Backend

## 📦 What is Docker?

**Docker** is a platform that packages your application and all its dependencies into a **container**. Think of it like a shipping container for software - it contains everything your app needs to run, and it works the same way on any computer.

### Key Concepts for Beginners

#### 1. **Docker Image**
- A blueprint or template for your application
- Contains your code, runtime, libraries, and dependencies
- Like a recipe for making a cake

#### 2. **Docker Container**
- A running instance of an image
- Like the actual cake made from the recipe
- Isolated environment that runs your application

#### 3. **Dockerfile**
- A text file with instructions to build a Docker image
- Lists all the steps needed to create your application image

#### 4. **Docker Compose**
- A tool to define and run multiple containers together
- Uses a YAML file to configure all services
- Allows you to start everything with one command

#### 5. **Docker Volume**
- Persistent storage for containers
- Data survives even when containers are stopped or deleted
- Like a hard drive for your containers

#### 6. **Docker Network**
- Allows containers to communicate with each other
- Like a private network for your services

## 🏗️ Project Architecture

Our Docker setup includes:

```
┌─────────────────────────────────────────────────┐
│                  Host Machine                    │
│                                                  │
│  ┌────────────────────────────────────────────┐ │
│  │         Docker Network                     │ │
│  │                                            │ │
│  │  ┌──────────┐  ┌────────┐  ┌──────────┐  │ │
│  │  │   Go     │  │ MySQL  │  │  Redis   │  │ │
│  │  │   App    │←→│   DB   │  │  Cache   │  │ │
│  │  │  :8080   │  │ :3306  │  │  :6379   │  │ │
│  │  └──────────┘  └────────┘  └──────────┘  │ │
│  │       ↓             ↓           ↓         │ │
│  │  ┌──────────┐  ┌────────┐  ┌──────────┐  │ │
│  │  │ Uploads  │  │ MySQL  │  │  Redis   │  │ │
│  │  │ Volume   │  │ Volume │  │  Volume  │  │ │
│  │  └──────────┘  └────────┘  └──────────┘  │ │
│  │                                            │ │
│  │  Optional Management Tools:                │ │
│  │  ┌──────────────┐  ┌──────────────────┐   │ │
│  │  │ phpMyAdmin   │  │ Redis Commander  │   │ │
│  │  │    :8081     │  │      :8082       │   │ │
│  │  └──────────────┘  └──────────────────┘   │ │
│  └────────────────────────────────────────────┘ │
│                                                  │
│  Ports Exposed:                                  │
│  - 8080: Application API                        │
│  - 3307: MySQL Database                         │
│  - 6380: Redis Cache                            │
│  - 8081: phpMyAdmin (DB Manager)                │
│  - 8082: Redis Commander (Cache Manager)        │
└─────────────────────────────────────────────────┘
```

## 📋 Prerequisites

### Install Docker

**Windows:**
1. Download [Docker Desktop for Windows](https://www.docker.com/products/docker-desktop)
2. Run the installer
3. Restart your computer
4. Verify installation: `docker --version`

**Mac:**
1. Download [Docker Desktop for Mac](https://www.docker.com/products/docker-desktop)
2. Drag to Applications folder
3. Open Docker Desktop
4. Verify installation: `docker --version`

**Linux:**
```bash
# Ubuntu/Debian
sudo apt-get update
sudo apt-get install docker.io docker-compose
sudo systemctl start docker
sudo systemctl enable docker

# Verify
docker --version
docker-compose --version
```

## 🚀 Quick Start

### Option 1: Production Setup (Recommended for Testing)

This will run everything in Docker containers.

#### Step 1: Prepare Environment File

```bash
# Copy the Docker environment template
cp .env.docker .env

# Edit .env and add your S3 credentials
# You can skip S3 config for now and add it later
```

#### Step 2: Build and Start All Services

```bash
# Build and start all containers
docker-compose up --build

# Or run in background (detached mode)
docker-compose up -d --build
```

**What happens:**
- Downloads MySQL, Redis, and other images (first time only)
- Builds your Go application image
- Creates network and volumes
- Starts all containers
- MySQL creates the database
- App connects to MySQL and Redis

#### Step 3: Access the Application

Wait about 30 seconds for all services to start, then:

- **API**: http://localhost:8080
- **Swagger UI**: http://localhost:8080/swagger/index.html
- **Health Check**: http://localhost:8080/health
- **phpMyAdmin**: http://localhost:8081 (username: `root`, password: `rootpassword`)
- **Redis Commander**: http://localhost:8082

#### Step 4: Test the API

```bash
# Check if it's running
curl http://localhost:8080/health

# Create a test user
curl -X POST http://localhost:8080/auth/signup \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Test User",
    "email": "test@example.com",
    "password": "password123",
    "role": "customer"
  }'
```

### Option 2: Development Setup

For development, run only MySQL and Redis in Docker, and run the Go app locally.

#### Step 1: Start Only Database Services

```bash
# Start MySQL and Redis
docker-compose -f docker-compose.dev.yml up -d
```

#### Step 2: Update Your Local .env

```env
# Use localhost instead of service names
DB_HOST=localhost
DB_PORT=3307
REDIS_HOST=localhost
REDIS_PORT=6380
# ... rest of your config
```

#### Step 3: Run Application Locally

```bash
# Run with hot reload
go run ./cmd/api/main.go
```

This is useful when you're actively coding and want to see changes immediately.

## 🎮 Docker Commands Cheat Sheet

### Starting and Stopping

```bash
# Start all services
docker-compose up

# Start in background (detached)
docker-compose up -d

# Stop all services
docker-compose down

# Stop and remove volumes (WARNING: deletes data!)
docker-compose down -v

# Restart a specific service
docker-compose restart app

# Rebuild and start (after code changes)
docker-compose up --build
```

### Viewing Logs

```bash
# View logs from all services
docker-compose logs

# Follow logs in real-time
docker-compose logs -f

# View logs from specific service
docker-compose logs app
docker-compose logs mysql
docker-compose logs redis

# Last 100 lines
docker-compose logs --tail=100 app
```

### Managing Containers

```bash
# List running containers
docker-compose ps

# List all containers (including stopped)
docker ps -a

# Enter a container shell
docker-compose exec app sh
docker-compose exec mysql bash

# Run a command in container
docker-compose exec app ls -la
docker-compose exec mysql mysql -uroot -prootpassword
```

### Database Operations

```bash
# Access MySQL CLI
docker-compose exec mysql mysql -uroot -prootpassword ecommerce_go

# Backup database
docker-compose exec mysql mysqldump -uroot -prootpassword ecommerce_go > backup.sql

# Restore database
docker-compose exec -T mysql mysql -uroot -prootpassword ecommerce_go < backup.sql

# Access Redis CLI
docker-compose exec redis redis-cli
```

### Cleaning Up

```bash
# Remove stopped containers
docker-compose rm

# Remove unused images
docker image prune

# Remove unused volumes
docker volume prune

# Remove everything (WARNING: deletes all data!)
docker system prune -a --volumes
```

## 📁 File Explanations

### 1. `Dockerfile`

This file builds your Go application image using a multi-stage build:

**Stage 1: Builder**
- Uses Go 1.25.3 Alpine (lightweight Linux)
- Installs build tools
- Downloads dependencies
- Compiles your application

**Stage 2: Runtime**
- Uses minimal Alpine image
- Copies only the compiled binary
- Creates non-root user for security
- Sets up healthcheck

### 2. `.dockerignore`

Lists files Docker should ignore when building:
- `.env` files (contain secrets)
- Git files
- IDE files
- Test files
- Documentation

### 3. `docker-compose.yml`

Defines all services and how they work together:

**mysql service:**
- Database for storing data
- Port 3307 (to avoid conflicts)
- Persistent volume for data
- Healthcheck to ensure it's ready

**redis service:**
- Cache for performance
- Port 6380
- Persistent storage enabled

**app service:**
- Your Go application
- Built from Dockerfile
- Depends on MySQL and Redis
- Environment variables configured

**phpmyadmin service:**
- Web interface for MySQL
- Port 8081
- Easy database management

**redis-commander service:**
- Web interface for Redis
- Port 8082
- View cached data

### 4. `.env.docker`

Template for environment variables in Docker:
- Uses Docker service names (`mysql`, `redis`) instead of `localhost`
- Different ports (3307, 6380)

## 🔧 Configuration

### Changing Ports

Edit `docker-compose.yml`:

```yaml
services:
  app:
    ports:
      - "9000:8080"  # Change host port to 9000
```

### Changing Database Password

Edit `docker-compose.yml`:

```yaml
mysql:
  environment:
    MYSQL_ROOT_PASSWORD: NewSecurePassword
    MYSQL_PASSWORD: NewUserPassword

app:
  environment:
    DB_PASSWORD: NewUserPassword
```

### Adding Environment Variables

Edit `docker-compose.yml` under `app.environment`:

```yaml
app:
  environment:
    NEW_CONFIG: value
```

## 🐛 Troubleshooting

### Problem: Containers won't start

**Solution:**
```bash
# Check logs
docker-compose logs

# Check if ports are already in use
netstat -ano | findstr :8080
netstat -ano | findstr :3307

# Stop conflicting services
```

### Problem: Database connection refused

**Solution:**
```bash
# Wait for MySQL to be ready (takes ~30 seconds)
docker-compose logs mysql

# Check MySQL health
docker-compose ps

# Manually test connection
docker-compose exec mysql mysql -uroot -prootpassword
```

### Problem: App can't connect to Redis

**Solution:**
```bash
# Check Redis is running
docker-compose logs redis

# Test Redis connection
docker-compose exec redis redis-cli ping
# Should return: PONG
```

### Problem: "Cannot build" error

**Solution:**
```bash
# Make sure you're in the project directory
cd /path/to/ecommerce-gin

# Clean and rebuild
docker-compose down
docker-compose build --no-cache
docker-compose up
```

### Problem: Out of disk space

**Solution:**
```bash
# Remove unused images
docker image prune -a

# Remove unused volumes
docker volume prune

# See disk usage
docker system df
```

## 📊 Monitoring

### Check Container Status

```bash
# See all running containers
docker-compose ps

# Check resource usage
docker stats

# Check container health
docker inspect ecommerce-app | grep -A 10 Health
```

### View Database

**Option 1: phpMyAdmin**
- Open: http://localhost:8081
- Login: `root` / `rootpassword`
- Browse tables visually

**Option 2: Command Line**
```bash
docker-compose exec mysql mysql -uroot -prootpassword ecommerce_go -e "SHOW TABLES;"
```

### View Redis Cache

**Option 1: Redis Commander**
- Open: http://localhost:8082
- Browse keys visually

**Option 2: Command Line**
```bash
docker-compose exec redis redis-cli
> KEYS *
> GET product:some-slug
```

## 🚀 Deployment

### Building for Production

```bash
# Build optimized image
docker-compose build --no-cache

# Tag for registry
docker tag ecommerce-app:latest your-registry/ecommerce-app:v1.0.0

# Push to registry
docker push your-registry/ecommerce-app:v1.0.0
```

### Environment Variables for Production

Create `.env.production`:

```env
APP_ENV=production
# Use strong passwords!
DB_PASSWORD=very-secure-password-here
JWT_SECRET=super-secret-jwt-key-32-chars-min
# Add real S3 credentials
S3_KEY=real-access-key
S3_SECRET=real-secret-key
```

## 📚 Best Practices

### 1. Never Commit Secrets

```bash
# Add to .gitignore
echo ".env" >> .gitignore
echo ".env.*" >> .gitignore
echo "!.env.docker" >> .gitignore
```

### 2. Use Volumes for Data

- Database data persists in volumes
- Uploads persist in volumes
- Data survives container restarts

### 3. Check Logs Regularly

```bash
# Monitor in real-time
docker-compose logs -f app
```

### 4. Backup Your Data

```bash
# Backup database
docker-compose exec mysql mysqldump -uroot -prootpassword ecommerce_go > backup-$(date +%Y%m%d).sql

# Backup volumes
docker run --rm -v ecommerce-gin_mysql-data:/data -v $(pwd):/backup alpine tar czf /backup/mysql-backup.tar.gz /data
```

## 🎓 Learning Resources

- [Docker Official Documentation](https://docs.docker.com/)
- [Docker Compose Documentation](https://docs.docker.com/compose/)
- [Docker Hub](https://hub.docker.com/) - Find pre-built images
- [Play with Docker](https://labs.play-with-docker.com/) - Online Docker playground

## ❓ FAQ

**Q: Do I need to install MySQL and Redis locally?**
A: No! Docker containers include everything.

**Q: Can I use my local MySQL?**
A: Yes, use `docker-compose.dev.yml` which only runs MySQL/Redis.

**Q: How do I stop everything?**
A: `docker-compose down`

**Q: Will I lose my data when stopping containers?**
A: No, data is stored in Docker volumes and persists.

**Q: Can I access the database from outside Docker?**
A: Yes, on port 3307: `mysql -h 127.0.0.1 -P 3307 -u ecommerce_user -p`

**Q: How do I update my code?**
A: Make changes, then run `docker-compose up --build`

## 🆘 Getting Help

If you encounter issues:

1. Check logs: `docker-compose logs`
2. Verify services are healthy: `docker-compose ps`
3. Restart services: `docker-compose restart`
4. Check this guide's troubleshooting section

---

**Happy Dockerizing! 🐳**
