```markdown
# Node.js Project Setup Guide

This guide provides detailed instructions for setting up the Node.js development environment, configuring the project, and getting started with development.

## 📋 Table of Contents

- [Prerequisites](#prerequisites)
- [Quick Installation](#quick-installation)
- [Environment Configuration](#environment-configuration)
- [Development Setup](#development-setup)
- [IDE Configuration](#ide-configuration)
- [Database Setup](#database-setup)
- [Redis Setup](#redis-setup)
- [Testing Setup](#testing-setup)
- [Docker Setup](#docker-setup)
- [Common Issues](#common-issues)
- [Troubleshooting](#troubleshooting)

## Prerequisites

Before you begin, ensure you have the following installed:

### Required Software

| Software | Version | Purpose | Installation |
|----------|---------|---------|--------------|
| **Node.js** | 20.x or higher | JavaScript runtime | [Download](https://nodejs.org/) |
| **npm** | 10.x or higher | Package manager | Included with Node.js |
| **Git** | 2.x or higher | Version control | [Download](https://git-scm.com/) |
| **Docker** (optional) | 24.x or higher | Containerization | [Download](https://docker.com/) |
| **PostgreSQL** (optional) | 15.x or higher | Database | [Download](https://postgresql.org/) |
| **Redis** (optional) | 7.x or higher | Cache | [Download](https://redis.io/) |

### Verify Installations

Open your terminal and run:

```bash
# Check Node.js version
node --version
# Should output: v20.x.x

# Check npm version
npm --version
# Should output: 10.x.x

# Check Git version
git --version
# Should output: 2.x.x

# Check Docker version (if installed)
docker --version
docker-compose --version
```

## Quick Installation

### 1. Clone the Repository

```bash
# Clone with HTTPS
git clone https://github.com/username/node-project.git
cd node-project

# Or clone with SSH
git clone git@github.com:username/node-project.git
cd node-project
```

### 2. Install Dependencies

```bash
# Install all dependencies
npm install

# Install with specific npm registry
npm install --registry=https://registry.npmjs.org

# Install with yarn (if you prefer)
yarn install

# Install with pnpm
pnpm install
```

### 3. Set Up Environment Variables

```bash
# Copy example environment file
cp .env.example .env

# Edit the .env file with your configuration
nano .env
# or
code .env
```

### 4. Run Database Migrations

```bash
# Run migrations
npm run migrate

# Seed the database with initial data
npm run seed
```

### 5. Start the Development Server

```bash
# Start in development mode
npm run dev

# The server should start at http://localhost:3000
```

## Environment Configuration

### Environment Files

The project uses different environment files for different stages:

| File | Purpose |
|------|---------|
| `.env.example` | Template with example values |
| `.env` | Local development (git-ignored) |
| `.env.test` | Testing environment |
| `.env.staging` | Staging environment |
| `.env.production` | Production environment |

### Environment Variables

Create a `.env` file with the following variables:

```env
# ============================================================================
# Server Configuration
# ============================================================================

# Node environment (development, test, production)
NODE_ENV=development

# Server port
PORT=3000

# Host (0.0.0.0 for all interfaces)
HOST=localhost

# API base URL
API_URL=http://localhost:3000/api

# Client URL (for CORS)
CLIENT_URL=http://localhost:5173

# ============================================================================
# Database Configuration
# ============================================================================

# PostgreSQL connection
DB_HOST=localhost
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=postgres
DB_NAME=myapp_dev
DB_POOL_MIN=2
DB_POOL_MAX=10
DB_IDLE_TIMEOUT=10000

# MongoDB connection (if using MongoDB)
MONGODB_URI=mongodb://localhost:27017/myapp_dev

# ============================================================================
# Redis Configuration
# ============================================================================

REDIS_HOST=localhost
REDIS_PORT=6379
REDIS_PASSWORD=
REDIS_DB=0

# ============================================================================
# Authentication
# ============================================================================

# JWT secret (use a strong random string in production)
JWT_SECRET=your-super-secret-jwt-key-change-this
JWT_EXPIRES_IN=7d
JWT_REFRESH_SECRET=your-refresh-secret-key
JWT_REFRESH_EXPIRES_IN=30d

# OAuth providers
GOOGLE_CLIENT_ID=your-google-client-id
GOOGLE_CLIENT_SECRET=your-google-client-secret
GITHUB_CLIENT_ID=your-github-client-id
GITHUB_CLIENT_SECRET=your-github-client-secret

# ============================================================================
# Email Configuration
# ============================================================================

SMTP_HOST=smtp.gmail.com
SMTP_PORT=587
SMTP_USER=your-email@gmail.com
SMTP_PASS=your-app-password
EMAIL_FROM=noreply@example.com

# ============================================================================
# File Upload
# ============================================================================

UPLOAD_DIR=./uploads
MAX_FILE_SIZE=10485760 # 10MB
ALLOWED_FILE_TYPES=image/jpeg,image/png,image/gif,application/pdf

# ============================================================================
# Logging
# ============================================================================

LOG_LEVEL=debug
LOG_FORMAT=combined
LOG_DIR=./logs

# ============================================================================
# Rate Limiting
# ============================================================================

RATE_LIMIT_WINDOW=60000 # 1 minute
RATE_LIMIT_MAX=100

# ============================================================================
# Cache
# ============================================================================

CACHE_TTL=3600 # 1 hour
CACHE_CHECK_PERIOD=120

# ============================================================================
# External APIs
# ============================================================================

STRIPE_API_KEY=sk_test_...
STRIPE_WEBHOOK_SECRET=whsec_...
SENDGRID_API_KEY=SG...
AWS_ACCESS_KEY_ID=your-access-key
AWS_SECRET_ACCESS_KEY=your-secret-key
AWS_REGION=us-east-1
S3_BUCKET=myapp-uploads

# ============================================================================
# Security
# ============================================================================

CORS_ORIGIN=http://localhost:5173
SESSION_SECRET=your-session-secret
COOKIE_SECURE=false # Set to true in production with HTTPS
BCRYPT_ROUNDS=10

# ============================================================================
# Monitoring
# ============================================================================

SENTRY_DSN=https://key@sentry.io/project
NEW_RELIC_LICENSE_KEY=your-license-key
NEW_RELIC_APP_NAME=myapp-dev

# ============================================================================
# Feature Flags
# ============================================================================

ENABLE_ANALYTICS=true
ENABLE_NOTIFICATIONS=true
ENABLE_BACKGROUND_JOBS=true
MAINTENANCE_MODE=false
```

### Loading Environment Variables

```javascript
// Using dotenv
require('dotenv').config();

// Access variables
const port = process.env.PORT || 3000;
const dbHost = process.env.DB_HOST;
const jwtSecret = process.env.JWT_SECRET;

// Validate required variables
const requiredEnvVars = ['DB_HOST', 'DB_USER', 'JWT_SECRET'];
requiredEnvVars.forEach(varName => {
    if (!process.env[varName]) {
        throw new Error(`Missing required environment variable: ${varName}`);
    }
});
```

## Development Setup

### 1. Install Global Tools

```bash
# Install nodemon for auto-restart
npm install -g nodemon

# Install eslint for linting
npm install -g eslint

# Install prettier for code formatting
npm install -g prettier

# Install sequelize-cli for database migrations
npm install -g sequelize-cli

# Install knex for database queries (if using Knex)
npm install -g knex
```

### 2. Configure Git Hooks

```bash
# Install husky for git hooks
npm install --save-dev husky

# Initialize husky
npx husky install

# Add pre-commit hook for linting
npx husky add .husky/pre-commit "npm run lint"

# Add pre-push hook for tests
npx husky add .husky/pre-push "npm test"
```

### 3. Set Up VS Code Debugging

Create `.vscode/launch.json`:

```json
{
    "version": "0.2.0",
    "configurations": [
        {
            "type": "node",
            "request": "launch",
            "name": "Launch Program",
            "skipFiles": ["<node_internals>/**"],
            "program": "${workspaceFolder}/src/index.js",
            "envFile": "${workspaceFolder}/.env"
        },
        {
            "type": "node",
            "request": "attach",
            "name": "Attach to Process",
            "port": 9229
        },
        {
            "type": "node",
            "request": "launch",
            "name": "Jest Tests",
            "program": "${workspaceFolder}/node_modules/.bin/jest",
            "args": ["--runInBand", "--watch"],
            "console": "integratedTerminal",
            "internalConsoleOptions": "neverOpen"
        },
        {
            "type": "node",
            "request": "launch",
            "name": "Jest Current File",
            "program": "${workspaceFolder}/node_modules/.bin/jest",
            "args": ["${fileBasename}", "--config", "jest.config.js"],
            "console": "integratedTerminal"
        }
    ]
}
```

### 4. Set Up VS Code Settings

Create `.vscode/settings.json`:

```json
{
    "editor.formatOnSave": true,
    "editor.codeActionsOnSave": {
        "source.fixAll.eslint": true
    },
    "eslint.validate": [
        "javascript",
        "javascriptreact",
        "typescript",
        "typescriptreact"
    ],
    "files.exclude": {
        "**/node_modules": true,
        "**/dist": true,
        "**/coverage": true
    },
    "search.exclude": {
        "**/node_modules": true,
        "**/dist": true,
        "**/coverage": true
    },
    "files.watcherExclude": {
        "**/node_modules/**": true,
        "**/dist/**": true,
        "**/coverage/**": true
    }
}
```

### 5. Set Up EditorConfig

Create `.editorconfig`:

```ini
root = true

[*]
indent_style = space
indent_size = 2
end_of_line = lf
charset = utf-8
trim_trailing_whitespace = true
insert_final_newline = true
max_line_length = 100

[*.md]
trim_trailing_whitespace = false
max_line_length = 0

[*.{json,yml,yaml}]
indent_size = 2
```

## IDE Configuration

### VS Code Extensions

Recommended extensions for Node.js development:

```bash
# Install via command line
code --install-extension dbaeumer.vscode-eslint
code --install-extension esbenp.prettier-vscode
code --install-extension Orta.vscode-jest
code --install-extension ms-azuretools.vscode-docker
code --install-extension eamodio.gitlens
code --install-extension wix.vscode-import-cost
code --install-extension christian-kohler.npm-intellisense
code --install-extension christian-kohler.path-intellisense
code --install-extension formulahendry.auto-rename-tag
code --install-extension CoenraadS.bracket-pair-colorizer-2
code --install-extension streetsidesoftware.code-spell-checker
```

### WebStorm Configuration

1. **Node.js interpreter**: Settings → Languages & Frameworks → Node.js
2. **npm scripts**: Tools → Tasks and Contexts → npm
3. **ESLint**: Settings → Languages & Frameworks → JavaScript → Code Quality Tools → ESLint
4. **Prettier**: Settings → Tools → File Watchers → Add Prettier

## Database Setup

### PostgreSQL Setup

```bash
# Start PostgreSQL
sudo service postgresql start
# or
brew services start postgresql

# Create database user
sudo -u postgres createuser --interactive
# or
psql -c "CREATE USER myapp WITH PASSWORD 'myapp' SUPERUSER;"

# Create database
sudo -u postgres createdb -O myapp myapp_dev
sudo -u postgres createdb -O myapp myapp_test

# Run migrations
npm run migrate

# Run seeds
npm run seed
```

### MongoDB Setup

```bash
# Start MongoDB
sudo service mongod start
# or
brew services start mongodb-community

# Create database
mongo
use myapp_dev
db.createUser({
    user: "myapp",
    pwd: "password",
    roles: ["readWrite"]
})

# Import initial data
mongoimport --db myapp_dev --collection users --file ./data/users.json --jsonArray
```

### Database Migration Commands

```bash
# Create a new migration
npm run migrate:create --name=add-users-table

# Run migrations
npm run migrate

# Rollback last migration
npm run migrate:rollback

# Rollback all migrations
npm run migrate:rollback:all

# Create seed
npm run seed:create --name=demo-users

# Run seeds
npm run seed

# Reset database (drop, migrate, seed)
npm run db:reset
```

## Redis Setup

### Install Redis

```bash
# Ubuntu/Debian
sudo apt-get install redis-server

# macOS
brew install redis

# Windows (using WSL)
sudo apt-get install redis-server

# Start Redis
sudo service redis-server start
# or
redis-server

# Connect to Redis
redis-cli
```

### Redis Configuration

```bash
# Test Redis connection
redis-cli ping
# Should return: PONG

# Set configuration
redis-cli CONFIG SET maxmemory 256mb
redis-cli CONFIG SET maxmemory-policy allkeys-lru

# Monitor Redis
redis-cli MONITOR

# Flush all data
redis-cli FLUSHALL
```

## Testing Setup

### Install Testing Tools

```bash
# Install Jest and related packages
npm install --save-dev jest supertest @types/jest

# Install testing utilities
npm install --save-dev faker axios-mock-adapter

# Install coverage tools
npm install --save-dev nyc coveralls
```

### Test Configuration

Create `jest.config.js`:

```javascript
module.exports = {
    testEnvironment: 'node',
    testMatch: ['**/__tests__/**/*.js', '**/?(*.)+(spec|test).js'],
    testPathIgnorePatterns: ['/node_modules/', '/dist/'],
    collectCoverageFrom: [
        'src/**/*.js',
        '!src/**/*.test.js',
        '!src/index.js'
    ],
    coverageThreshold: {
        global: {
            branches: 80,
            functions: 80,
            lines: 80,
            statements: 80
        }
    },
    setupFilesAfterEnv: ['./jest.setup.js'],
    verbose: true,
    testTimeout: 10000
};
```

Create `jest.setup.js`:

```javascript
// Set test environment variables
process.env.NODE_ENV = 'test';
process.env.PORT = 3001;
process.env.DB_NAME = 'myapp_test';
process.env.JWT_SECRET = 'test-secret-key';

// Global setup
beforeAll(async () => {
    // Connect to test database
    // Setup test data
});

afterAll(async () => {
    // Cleanup
    // Disconnect from database
});
```

### Test Database Setup

```javascript
// test/setup.js
const { sequelize } = require('../src/models');

beforeAll(async () => {
    await sequelize.sync({ force: true });
});

afterAll(async () => {
    await sequelize.close();
});
```

## Docker Setup

### Dockerfile

Create `Dockerfile`:

```dockerfile
# Multi-stage build
FROM node:20-alpine AS builder

WORKDIR /app

# Copy package files
COPY package*.json ./
RUN npm ci --only=production

# Development stage
FROM node:20-alpine AS development

WORKDIR /app

# Install dependencies
COPY package*.json ./
RUN npm install

# Copy source code
COPY . .

# Expose port
EXPOSE 3000

# Start in development mode
CMD ["npm", "run", "dev"]

# Production stage
FROM node:20-alpine AS production

WORKDIR /app

# Copy production dependencies
COPY --from=builder /app/node_modules ./node_modules
COPY . .

# Create non-root user
RUN addgroup -g 1001 -S nodejs && \
    adduser -S nodejs -u 1001

USER nodejs

EXPOSE 3000

HEALTHCHECK --interval=30s --timeout=3s --start-period=5s --retries=3 \
    CMD node healthcheck.js

CMD ["node", "src/index.js"]
```

### Docker Compose

Create `docker-compose.yml`:

```yaml
version: '3.8'

services:
  app:
    build:
      context: .
      target: development
    container_name: node-app
    restart: unless-stopped
    ports:
      - "3000:3000"
      - "9229:9229" # Debug port
    environment:
      - NODE_ENV=development
      - DB_HOST=postgres
      - DB_PORT=5432
      - DB_USER=postgres
      - DB_PASSWORD=postgres
      - DB_NAME=myapp_dev
      - REDIS_HOST=redis
      - REDIS_PORT=6379
    volumes:
      - .:/app
      - /app/node_modules
    depends_on:
      - postgres
      - redis
    networks:
      - app-network

  postgres:
    image: postgres:15-alpine
    container_name: postgres
    restart: unless-stopped
    ports:
      - "5432:5432"
    environment:
      - POSTGRES_USER=postgres
      - POSTGRES_PASSWORD=postgres
      - POSTGRES_DB=myapp_dev
    volumes:
      - postgres-data:/var/lib/postgresql/data
    networks:
      - app-network
    healthcheck:
      test: ["CMD-SHELL", "pg_isready -U postgres"]
      interval: 10s
      timeout: 5s
      retries: 5

  redis:
    image: redis:7-alpine
    container_name: redis
    restart: unless-stopped
    ports:
      - "6379:6379"
    volumes:
      - redis-data:/data
    networks:
      - app-network
    healthcheck:
      test: ["CMD", "redis-cli", "ping"]
      interval: 10s
      timeout: 5s
      retries: 5

  pgadmin:
    image: dpage/pgadmin4
    container_name: pgadmin
    restart: unless-stopped
    ports:
      - "5050:80"
    environment:
      - PGADMIN_DEFAULT_EMAIL=admin@example.com
      - PGADMIN_DEFAULT_PASSWORD=admin
    volumes:
      - pgadmin-data:/var/lib/pgadmin
    depends_on:
      - postgres
    networks:
      - app-network
    profiles:
      - tools

volumes:
  postgres-data:
  redis-data:
  pgadmin-data:

networks:
  app-network:
    driver: bridge
```

### Docker Commands

```bash
# Build and start services
docker-compose up -d

# View logs
docker-compose logs -f app

# Execute command in container
docker-compose exec app npm run migrate

# Run tests in container
docker-compose exec app npm test

# Stop services
docker-compose down

# Stop and remove volumes
docker-compose down -v

# Rebuild specific service
docker-compose up -d --build app

# Scale services
docker-compose up -d --scale app=3
```

## Common Issues

### Installation Issues

| Issue | Solution |
|-------|----------|
| `npm install` fails with permission errors | Run `sudo chown -R $(whoami) ~/.npm` |
| Node Sass installation fails | Install with `npm install --force` |
| Python build errors | Install Python build tools: `npm install --global windows-build-tools` |
| Git not found | Install Git or use HTTPS URLs |

### Environment Issues

| Issue | Solution |
|-------|----------|
| `.env` file not loading | Ensure `dotenv` is required at the top of your entry file |
| Port already in use | Change PORT in `.env` or kill the process using the port |
| Database connection refused | Check if database is running and credentials are correct |
| Redis connection timeout | Verify Redis is running and host/port are correct |

### Database Issues

| Issue | Solution |
|-------|----------|
| Migration fails | Check migration syntax and rollback to fix |
| Connection pool exhausted | Increase pool size or check for connection leaks |
| Data truncation | Increase column size or check data types |
| Foreign key constraint fails | Ensure referenced data exists or disable constraints temporarily |

### Testing Issues

| Issue | Solution |
|-------|----------|
| Tests timeout | Increase timeout in jest.config.js |
| Database tests interfering | Use a separate test database and clean between tests |
| Mock not working | Check that mock is applied before the module is imported |
| Coverage too low | Add more tests or adjust coverage threshold |

## Troubleshooting

### Debug Mode

```bash
# Run with debug logs
DEBUG=app:* npm run dev

# Run with Node.js inspector
node --inspect src/index.js

# Debug with Chrome DevTools
chrome://inspect
```

### Logs

```bash
# Application logs
tail -f logs/app.log

# Nginx logs
tail -f /var/log/nginx/access.log
tail -f /var/log/nginx/error.log

# Database logs
tail -f /var/log/postgresql/postgresql.log
```

### Performance Profiling

```bash
# Use Node.js profiler
node --prof src/index.js
node --prof-process isolate-*.log > profile.txt

# Use clinic.js
npm install -g clinic
clinic doctor -- node src/index.js
clinic flame -- node src/index.js
```

### Memory Leak Detection

```bash
# Use heap dump
node --heapsnapshot-signal SIGUSR2 src/index.js
kill -USR2 <pid>

# Use memwatch-next
npm install memwatch-next
```

### Network Debugging

```bash
# Check open ports
netstat -tulpn | grep LISTEN

# Test database connection
telnet localhost 5432

# Test API endpoints
curl -X GET http://localhost:3000/api/health
curl -X POST http://localhost:3000/api/users -H "Content-Type: application/json" -d '{"name":"test"}'
```

## Quick Reference

### Useful Commands

```bash
# Development
npm run dev          # Start development server
npm run build        # Build for production
npm start           # Start production server

# Testing
npm test            # Run tests
npm run test:watch  # Run tests in watch mode
npm run test:coverage # Run tests with coverage

# Database
npm run migrate     # Run migrations
npm run migrate:rollback # Rollback migration
npm run seed        # Seed database

# Code Quality
npm run lint        # Run linter
npm run lint:fix    # Fix linting issues
npm run format      # Format code

# Docker
docker-compose up   # Start services
docker-compose down # Stop services
docker-compose logs # View logs
```

### Project Structure

```
project/
├── src/
│   ├── controllers/    # Request handlers
│   ├── models/         # Database models
│   ├── services/       # Business logic
│   ├── middleware/     # Express middleware
│   ├── utils/          # Utility functions
│   ├── config/         # Configuration files
│   ├── routes/         # API routes
│   └── index.js        # Entry point
├── tests/
│   ├── unit/           # Unit tests
│   ├── integration/    # Integration tests
│   └── fixtures/       # Test data
├── migrations/         # Database migrations
├── seeds/              # Seed data
├── logs/               # Application logs
├── uploads/            # File uploads
├── scripts/            # Utility scripts
├── .env.example        # Example environment
├── .eslintrc.js       # ESLint config
├── .prettierrc        # Prettier config
├── jest.config.js     # Jest config
├── docker-compose.yml # Docker compose
├── Dockerfile         # Docker build
└── package.json       # Dependencies
```

---

**Need more help?** Check our [FAQ](FAQ.md) or [open an issue](https://github.com/username/node-project/issues).

*Last Updated: 2024*
```

## ✅ **What this SETUP guide provides:**

| Section | Purpose |
|---------|---------|
| **Prerequisites** | Required software and version checks |
| **Quick Installation** | Fast-track setup in 5 steps |
| **Environment Configuration** | Detailed .env file with all variables |
| **Development Setup** | Global tools, git hooks, IDE config |
| **IDE Configuration** | VS Code extensions and settings |
| **Database Setup** | PostgreSQL, MongoDB, migrations |
| **Redis Setup** | Installation and configuration |
| **Testing Setup** | Jest configuration and test database |
| **Docker Setup** | Dockerfile, compose, commands |
| **Common Issues** | Solutions to frequent problems |
| **Troubleshooting** | Debugging and profiling techniques |
| **Quick Reference** | Command cheatsheet and project structure |

## 🎯 **Purpose as Test Data**

This file serves as an **expected output** for validating that your readme generator produces:
- ✅ Comprehensive setup documentation
- ✅ Environment configuration examples
- ✅ Docker configuration
- ✅ Testing setup
- ✅ Troubleshooting guides
- ✅ Command references