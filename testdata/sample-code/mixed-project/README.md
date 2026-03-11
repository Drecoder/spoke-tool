```markdown
# Mixed-Language Project

A modern microservices-based application demonstrating integration between Go, Node.js, and Python services with Docker orchestration, monitoring, and CI/CD pipelines.

## 📋 Table of Contents

- [Overview](#overview)
- [Architecture](#architecture)
- [Prerequisites](#prerequisites)
- [Quick Start](#quick-start)
- [Development Setup](#development-setup)
- [Services](#services)
- [API Documentation](#api-documentation)
- [Testing](#testing)
- [Deployment](#deployment)
- [Monitoring](#monitoring)
- [Contributing](#contributing)
- [License](#license)

## 🎯 Overview

This project showcases a complete microservices architecture with:

- **Go Backend API** - High-performance REST API with PostgreSQL, Redis, and JWT authentication
- **Node.js Frontend** - Modern web application with React and real-time features
- **Python Processor** - Data processing service with ML capabilities and batch jobs
- **Message Queue** - RabbitMQ for async communication between services
- **Monitoring Stack** - Prometheus, Grafana, ELK stack, and Jaeger tracing
- **Container Orchestration** - Docker Compose for local development and production

## 🏗️ Architecture

```
┌─────────────────────────────────────────────────────────────────────┐
│                         Load Balancer (Traefik)                      │
│                              :80 :443                                │
└─────────────────────────────────────────────────────────────────────┘
                                    │
                    ┌───────────────┼───────────────┐
                    │               │               │
                    ▼               ▼               ▼
          ┌─────────────────┐ ┌─────────────────┐ ┌─────────────────┐
          │   Go API        │ │ Node.js Frontend│ │ Python Processor│
          │   Port: 8080    │ │   Port: 3000    │ │   Port: 8000    │
          │   - REST API    │ │ - React SPA     │ │ - Data Pipeline │
          │   - Auth/JWT    │ │ - WebSocket     │ │ - ML Models     │
          │   - GORM        │ │ - State Mgmt    │ │ - Celery Tasks  │
          └────────┬────────┘ └────────┬────────┘ └────────┬────────┘
                   │                   │                    │
                   └───────────────────┼────────────────────┘
                                      │
                    ┌─────────────────┼─────────────────┐
                    │                 │                 │
                    ▼                 ▼                 ▼
          ┌─────────────────┐ ┌─────────────────┐ ┌─────────────────┐
          │   PostgreSQL    │ │     Redis       │ │    RabbitMQ     │
          │   Port: 5432    │ │   Port: 6379    │ │   Port: 5672    │
          │   - Main DB     │ │ - Cache         │ │ - Message Queue │
          │   - User Data   │ │ - Sessions      │ │ - Task Queue    │
          │   - Orders      │ │ - Rate Limiting │ │ - Events        │
          └─────────────────┘ └─────────────────┘ └─────────────────┘
                                      │
                    ┌─────────────────┼─────────────────┐
                    │                 │                 │
                    ▼                 ▼                 ▼
          ┌─────────────────┐ ┌─────────────────┐ ┌─────────────────┐
          │   Prometheus    │ │    Grafana      │ │    Jaeger       │
          │   Port: 9090    │ │   Port: 3001    │ │   Port: 16686   │
          │   - Metrics     │ │ - Dashboards    │ │ - Tracing       │
          │   - Alerts      │ │ - Visualizations│ │ - Performance   │
          └─────────────────┘ └─────────────────┘ └─────────────────┘
```

## ✅ Prerequisites

- **Docker** 24.0+
- **Docker Compose** 2.20+
- **Git** 2.40+
- **Make** 4.0+ (optional)
- **Node.js** 20+ (for local development)
- **Go** 1.21+ (for local development)
- **Python** 3.11+ (for local development)

## 🚀 Quick Start

### 1. Clone the Repository

```bash
git clone https://github.com/yourusername/mixed-project.git
cd mixed-project
```

### 2. Set Up Environment Variables

```bash
# Copy example environment files
cp .env.example .env
cp go-server/.env.example go-server/.env
cp web-client/.env.example web-client/.env
cp scripts/.env.example scripts/.env

# Edit configuration as needed
nano .env
```

### 3. Start All Services

```bash
# Using Docker Compose
docker-compose up -d

# Or using Make
make up

# View logs
docker-compose logs -f
```

### 4. Access the Applications

| Service | URL | Credentials |
|---------|-----|-------------|
| Web Frontend | http://localhost:3000 | - |
| API Gateway | http://localhost:8080 | - |
| API Docs (Swagger) | http://localhost:8080/swagger | - |
| PostgreSQL Admin | http://localhost:8081 | postgres/postgres |
| Redis Commander | http://localhost:8082 | admin/admin |
| RabbitMQ Management | http://localhost:15672 | guest/guest |
| Prometheus | http://localhost:9090 | - |
| Grafana | http://localhost:3001 | admin/admin |
| Jaeger UI | http://localhost:16686 | - |
| Kibana | http://localhost:5601 | - |
| Mailhog | http://localhost:8025 | - |
| MinIO Console | http://localhost:9001 | minioadmin/minioadmin |

## 🔧 Development Setup

### Local Development (Without Docker)

#### Go API

```bash
cd go-server

# Install dependencies
go mod download

# Run database migrations
go run cmd/migrate/main.go

# Start development server with hot reload
air

# Or run normally
go run cmd/server/main.go
```

#### Node.js Frontend

```bash
cd web-client

# Install dependencies
npm install

# Start development server
npm run dev

# Build for production
npm run build

# Run tests
npm test
```

#### Python Processor

```bash
cd scripts

# Create virtual environment
python -m venv venv
source venv/bin/activate  # On Windows: venv\Scripts\activate

# Install dependencies
pip install -r requirements.txt
pip install -r requirements-dev.txt

# Start development server
uvicorn main:app --reload --port 8000

# Run Celery worker
celery -A tasks worker --loglevel=info
```

### Docker Development

```bash
# Build all services
docker-compose build

# Start specific service
docker-compose up -d go-api

# View logs for specific service
docker-compose logs -f go-api

# Execute command in container
docker-compose exec go-api /bin/sh

# Run database migrations
docker-compose exec go-api go run cmd/migrate/main.go

# Run tests in container
docker-compose run test
```

### Makefile Commands

```bash
make help           # Show available commands
make build          # Build all services
make up             # Start all services
make down           # Stop all services
make logs           # View logs
make ps             # List services
make test           # Run tests
make lint           # Run linters
make fmt            # Format code
make migrate        # Run database migrations
make seed           # Seed database
make backup         # Create database backup
make restore        # Restore database
make clean          # Clean up resources
```

## 📦 Services

### Go API Service

**Location:** `./go-server`

**Technologies:**
- Go 1.21 with Gin framework
- GORM for database ORM
- PostgreSQL 15
- Redis 7 for caching
- JWT authentication
- Swagger documentation
- Prometheus metrics

**Key Features:**
- RESTful API design
- JWT authentication
- Role-based access control
- Rate limiting
- Request validation
- Comprehensive logging
- Health checks
- Graceful shutdown

**Configuration:**

```env
# go-server/.env
APP_ENV=development
PORT=8080
DB_HOST=localhost
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=postgres
DB_NAME=go_api_db
REDIS_HOST=localhost
REDIS_PORT=6379
JWT_SECRET=your-secret-key
LOG_LEVEL=debug
```

### Node.js Frontend

**Location:** `./web-client`

**Technologies:**
- React 18 with TypeScript
- Redux Toolkit for state management
- React Query for data fetching
- React Router for navigation
- Socket.io for real-time updates
- TailwindCSS for styling
- Jest for testing
- Vite for build tooling

**Key Features:**
- Responsive design
- Real-time updates via WebSocket
- Form validation
- Error boundaries
- Loading skeletons
- Dark mode support
- PWA capabilities
- Accessibility (WCAG 2.1)

**Configuration:**

```env
# web-client/.env
VITE_API_URL=http://localhost:8080
VITE_WS_URL=ws://localhost:8080
VITE_GA_TRACKING_ID=UA-XXXXX-Y
VITE_SENTRY_DSN=https://key@sentry.io/project
```

### Python Processor

**Location:** `./scripts`

**Technologies:**
- Python 3.11 with FastAPI
- SQLAlchemy for database
- Celery for task queue
- Pandas for data processing
- Scikit-learn for ML models
- Pytest for testing
- Black for formatting

**Key Features:**
- Async request handling
- Background task processing
- Data validation with Pydantic
- Machine learning inference
- Batch job scheduling
- File upload/download
- Data export (CSV, JSON, Excel)

**Configuration:**

```env
# scripts/.env
PROCESSOR_ENV=development
DATABASE_URL=postgresql://postgres:postgres@localhost/processor_db
REDIS_URL=redis://localhost:6379/2
CELERY_BROKER_URL=amqp://guest:guest@localhost:5672//
CELERY_RESULT_BACKEND=redis://localhost:6379/2
MODEL_PATH=/app/models
```

## 📚 API Documentation

### Authentication

All API endpoints (except public ones) require JWT authentication.

**Login**

```bash
curl -X POST http://localhost:8080/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{"email":"user@example.com","password":"password"}'
```

**Response**

```json
{
  "token": "eyJhbGciOiJIUzI1NiIs...",
  "user": {
    "id": 1,
    "email": "user@example.com",
    "name": "John Doe",
    "role": "user"
  }
}
```

### Users API

| Method | Endpoint | Description | Auth |
|--------|----------|-------------|------|
| GET | `/api/v1/users` | List users | Admin |
| GET | `/api/v1/users/{id}` | Get user | User |
| POST | `/api/v1/users` | Create user | Public |
| PUT | `/api/v1/users/{id}` | Update user | User |
| DELETE | `/api/v1/users/{id}` | Delete user | Admin |

### Products API

| Method | Endpoint | Description | Auth |
|--------|----------|-------------|------|
| GET | `/api/v1/products` | List products | Public |
| GET | `/api/v1/products/{id}` | Get product | Public |
| POST | `/api/v1/products` | Create product | Admin |
| PUT | `/api/v1/products/{id}` | Update product | Admin |
| DELETE | `/api/v1/products/{id}` | Delete product | Admin |

### Orders API

| Method | Endpoint | Description | Auth |
|--------|----------|-------------|------|
| GET | `/api/v1/orders` | List orders | User |
| GET | `/api/v1/orders/{id}` | Get order | User |
| POST | `/api/v1/orders` | Create order | User |
| PUT | `/api/v1/orders/{id}/status` | Update status | Admin |

## 🧪 Testing

### Running Tests

```bash
# Run all tests
make test

# Run Go tests
cd go-server && go test ./... -v

# Run Node.js tests
cd web-client && npm test

# Run Python tests
cd scripts && pytest -v

# Run integration tests
make test-integration

# Run with coverage
make test-coverage
```

### Test Structure

```
tests/
├── unit/           # Unit tests
├── integration/    # Integration tests
├── e2e/           # End-to-end tests
├── fixtures/      # Test data
└── mocks/         # Mock objects
```

## 🚢 Deployment

### Production Deployment

```bash
# Set production environment
export APP_ENV=production
export TAG=v1.0.0

# Build and deploy
docker-compose -f docker-compose.prod.yml build
docker-compose -f docker-compose.prod.yml up -d
```

### Kubernetes Deployment

```bash
# Apply Kubernetes manifests
kubectl apply -f k8s/namespace.yaml
kubectl apply -f k8s/configmap.yaml
kubectl apply -f k8s/secrets.yaml
kubectl apply -f k8s/deployments/
kubectl apply -f k8s/services/
kubectl apply -f k8s/ingress.yaml
```

### CI/CD Pipeline

The project includes GitHub Actions workflows for:

- **CI**: Lint, test, build on every push
- **CD**: Deploy to staging on main branch
- **Release**: Create release and deploy to production on tags

## 📊 Monitoring

### Prometheus Metrics

Each service exposes metrics at `/metrics`:

- **Go API**: http://localhost:8080/metrics
- **Node.js**: http://localhost:3000/metrics
- **Python**: http://localhost:8000/metrics

### Grafana Dashboards

Access Grafana at http://localhost:3001

Pre-configured dashboards:
- **Service Overview** - Request rates, errors, latencies
- **Database** - Connection pools, query performance
- **System** - CPU, memory, disk usage
- **Business** - Orders, users, products metrics

### Logging with ELK Stack

- **Elasticsearch**: http://localhost:9200
- **Kibana**: http://localhost:5601
- **Logstash**: Collects logs from all services

### Tracing with Jaeger

Access Jaeger UI at http://localhost:16686

View distributed traces across services:
- API calls between services
- Database queries
- External service calls
- Async operations

## 🤝 Contributing

### Development Workflow

1. Fork the repository
2. Create a feature branch (`git checkout -b feature/amazing-feature`)
3. Commit changes (`git commit -m 'Add amazing feature'`)
4. Push to branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

### Code Style

- **Go**: Follow [Go Code Review Comments](https://github.com/golang/go/wiki/CodeReviewComments)
- **JavaScript/TypeScript**: Use ESLint and Prettier
- **Python**: Follow PEP 8, use Black for formatting

### Commit Convention

We follow [Conventional Commits](https://www.conventionalcommits.org/):

```
feat: add new feature
fix: bug fix
docs: documentation update
style: code formatting
refactor: code restructuring
test: add tests
chore: maintenance tasks
```

### Branch Strategy

- `main` - Production-ready code
- `develop` - Development branch
- `feature/*` - New features
- `bugfix/*` - Bug fixes
- `release/*` - Release preparation

## 📄 License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## 👥 Authors

- **John Doe** - *Lead Developer* - [@johndoe](https://github.com/johndoe)
- **Jane Smith** - *Backend Developer* - [@janesmith](https://github.com/janesmith)
- **Bob Johnson** - *Frontend Developer* - [@bobjohnson](https://github.com/bobjohnson)

## 🙏 Acknowledgments

- Go community for excellent tooling
- Node.js ecosystem
- Python data science stack
- Docker for containerization

## 📞 Support

- 📧 Email: support@example.com
- 💬 Discord: [Join our Discord](https://discord.gg/example)
- 🐛 GitHub Issues: [Create an issue](https://github.com/yourusername/mixed-project/issues)

## 📈 Project Status

| Service | Status | Coverage | Performance |
|---------|--------|----------|-------------|
| Go API | ✅ Stable | 92% | 🟢 Good |
| Node.js | ✅ Stable | 88% | 🟢 Good |
| Python | 🟡 Beta | 76% | 🟡 Improving |

## 🗺️ Roadmap

### Q1 2024
- [x] Initial architecture
- [x] Basic API endpoints
- [x] Docker setup

### Q2 2024
- [ ] Real-time features with WebSockets
- [ ] Machine learning integration
- [ ] Performance optimization

### Q3 2024
- [ ] GraphQL API
- [ ] Mobile app support
- [ ] Multi-region deployment

## 📊 Performance Benchmarks

| Operation | Go API | Python | Node.js |
|-----------|--------|--------|---------|
| GET /users | 5ms | 15ms | 8ms |
| POST /users | 10ms | 25ms | 12ms |
| Database query | 2ms | 5ms | 3ms |
| Cache hit | <1ms | 2ms | 1ms |

## 🔒 Security

- JWT authentication with refresh tokens
- Rate limiting per IP/user
- SQL injection prevention
- XSS protection
- CORS configuration
- HTTPS in production
- Secrets management with Vault

## 📚 Additional Resources

- [API Documentation](http://localhost:8080/swagger)
- [Postman Collection](./postman/collection.json)
- [Architecture Decision Records](./docs/adr/)
- [Contributing Guidelines](./CONTRIBUTING.md)
- [Code of Conduct](./CODE_OF_CONDUCT.md)

---

**Built with ❤️ using Go, Node.js, Python, and Docker**

[![Go Report Card](https://goreportcard.com/badge/github.com/yourusername/mixed-project)](https://goreportcard.com/report/github.com/yourusername/mixed-project)
[![Node.js CI](https://github.com/yourusername/mixed-project/actions/workflows/node.yml/badge.svg)](https://github.com/yourusername/mixed-project/actions/workflows/node.yml)
[![Python CI](https://github.com/yourusername/mixed-project/actions/workflows/python.yml/badge.svg)](https://github.com/yourusername/mixed-project/actions/workflows/python.yml)
[![Docker pulls](https://img.shields.io/docker/pulls/yourusername/mixed-project)](https://hub.docker.com/r/yourusername/mixed-project)
[![License](https://img.shields.io/github/license/yourusername/mixed-project)](LICENSE)
```

## ✅ **What this README demonstrates:**

| Section | Content |
|---------|---------|
| **Table of Contents** | Complete navigation with anchor links |
| **Overview** | Project description and technology stack |
| **Architecture** | ASCII diagram showing service interaction |
| **Prerequisites** | Required software and versions |
| **Quick Start** | Step-by-step setup instructions |
| **Development Setup** | Local development without Docker |
| **Services** | Detailed service descriptions with tech stack |
| **API Documentation** | Endpoint tables and examples |
| **Testing** | Test commands and structure |
| **Deployment** | Production deployment options |
| **Monitoring** | Prometheus, Grafana, ELK, Jaeger |
| **Contributing** | Workflow, code style, commit convention |
| **Project Status** | Service status and roadmap |
| **Performance** | Benchmark comparisons |
| **Security** | Security features |
| **Resources** | Links to additional documentation |

This is a production-grade README that demonstrates comprehensive project documentation for testing your code analyzer! 📚