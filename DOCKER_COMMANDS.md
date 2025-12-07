# Docker Commands Reference Card

## 🚀 Quick Start Commands

### Start Everything
```bash
# Start all services (production mode)
docker-compose up -d

# Start with rebuild
docker-compose up --build -d

# Start and view logs
docker-compose up
```

### Stop Everything
```bash
# Stop all services
docker-compose down

# Stop and remove volumes (⚠️ DELETES DATA)
docker-compose down -v
```

## 📊 Monitoring & Logs

### View Logs
```bash
# All services
docker-compose logs

# Follow logs in real-time
docker-compose logs -f

# Specific service
docker-compose logs app
docker-compose logs mysql
docker-compose logs redis

# Last 100 lines
docker-compose logs --tail=100 app
```

### Check Status
```bash
# List running containers
docker-compose ps

# See resource usage
docker stats

# Check health
docker-compose ps
```

## 🔧 Managing Services

### Restart Services
```bash
# Restart all
docker-compose restart

# Restart specific service
docker-compose restart app
docker-compose restart mysql
```

### Execute Commands in Containers
```bash
# Enter container shell
docker-compose exec app sh
docker-compose exec mysql bash

# Run a command
docker-compose exec app ls -la
docker-compose exec mysql mysql -uroot -prootpassword
```

## 🗄️ Database Commands

### MySQL
```bash
# Access MySQL CLI
docker-compose exec mysql mysql -uroot -prootpassword ecommerce_go

# Show tables
docker-compose exec mysql mysql -uroot -prootpassword ecommerce_go -e "SHOW TABLES;"

# Backup database
docker-compose exec mysql mysqldump -uroot -prootpassword ecommerce_go > backup.sql

# Restore database
docker-compose exec -T mysql mysql -uroot -prootpassword ecommerce_go < backup.sql
```

### Redis
```bash
# Access Redis CLI
docker-compose exec redis redis-cli

# Test connection
docker-compose exec redis redis-cli ping

# Get all keys
docker-compose exec redis redis-cli KEYS "*"

# Flush all data (⚠️ DELETES ALL CACHE)
docker-compose exec redis redis-cli FLUSHALL
```

## 🧹 Cleanup Commands

### Remove Stopped Containers
```bash
docker-compose rm
```

### Remove Images
```bash
# Remove project images
docker-compose down --rmi all

# Remove unused images
docker image prune

# Remove all unused images
docker image prune -a
```

### Remove Volumes
```bash
# Remove unused volumes
docker volume prune

# List volumes
docker volume ls

# Remove specific volume
docker volume rm ecommerce-gin_mysql-data
```

### Complete Cleanup
```bash
# Remove everything (⚠️ NUCLEAR OPTION)
docker system prune -a --volumes
```

## 🔄 Development Workflow

### Rebuild After Code Changes
```bash
# Rebuild and restart app
docker-compose up --build -d app

# Or stop, rebuild, and start
docker-compose stop app
docker-compose build app
docker-compose up -d app
```

### View Real-time Logs During Development
```bash
# Follow app logs
docker-compose logs -f app

# Follow all logs
docker-compose logs -f
```

## 🐛 Debugging

### Check Container Logs
```bash
# See errors in app
docker-compose logs app | grep -i error

# See MySQL errors
docker-compose logs mysql | grep -i error
```

### Inspect Container
```bash
# Get detailed container info
docker inspect ecommerce-app

# Check environment variables
docker inspect ecommerce-app | grep -A 20 Env
```

### Test Network Connectivity
```bash
# Ping from app to mysql
docker-compose exec app ping mysql

# Check DNS resolution
docker-compose exec app nslookup mysql
```

## 📦 Build Commands

### Build Images
```bash
# Build all images
docker-compose build

# Build without cache
docker-compose build --no-cache

# Build specific service
docker-compose build app
```

### Tag and Push (for deployment)
```bash
# Build and tag
docker build -t your-registry/ecommerce-app:v1.0.0 .

# Push to registry
docker push your-registry/ecommerce-app:v1.0.0
```

## 🎯 Useful One-Liners

### Quick Health Check
```bash
# Check if app is responding
curl http://localhost:8080/health
```

### See Disk Usage
```bash
docker system df
```

### Get Container IP
```bash
docker inspect -f '{{range .NetworkSettings.Networks}}{{.IPAddress}}{{end}}' ecommerce-app
```

### Copy Files from Container
```bash
# Copy from container to host
docker cp ecommerce-app:/app/uploads/image.jpg ./local-image.jpg

# Copy from host to container
docker cp ./local-file.txt ecommerce-app:/app/
```

## 🔐 Security Commands

### Check for Vulnerabilities
```bash
# Scan image for vulnerabilities
docker scan ecommerce-app:latest
```

### View Container Processes
```bash
docker-compose top
```

## 📊 Performance Monitoring

### Resource Usage
```bash
# Real-time stats
docker stats

# Stats for specific container
docker stats ecommerce-app
```

### Container Events
```bash
# Watch events in real-time
docker events

# Filter events
docker events --filter container=ecommerce-app
```

## 🎓 Tips & Tricks

### Auto-remove Containers on Exit
```bash
docker-compose up --force-recreate --remove-orphans
```

### Scale Services (if needed)
```bash
# Run multiple instances
docker-compose up -d --scale app=3
```

### Export/Import Images
```bash
# Export image
docker save ecommerce-app:latest | gzip > ecommerce-app.tar.gz

# Import image
docker load < ecommerce-app.tar.gz
```

---

## 🆘 Quick Troubleshooting

**Problem: Port already in use**
```bash
# Find process using port
netstat -ano | findstr :8080  # Windows
lsof -i :8080                  # Linux/Mac

# Change port in docker-compose.yml
ports:
  - "9000:8080"  # Use port 9000 instead
```

**Problem: Container keeps restarting**
```bash
# Check logs
docker-compose logs app

# Check if there's an error in config
docker-compose config
```

**Problem: Can't connect to database**
```bash
# Check if MySQL is ready
docker-compose logs mysql | grep "ready for connections"

# Test connection
docker-compose exec app ping mysql
```

**Problem: Out of disk space**
```bash
# Clean up
docker system prune -a --volumes
```

---

**Save this file for quick reference! 📌**
