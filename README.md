# AUTH MANAGEMENT

Service for authentication management for microservices architecture.

## Tech Stack

- **Language:** Go
- **Database:** PostgreSQL
- **Cache:** Memcached
- **Message Broker:** LavinMQ / RabbitMQ
- **Containerization:** Docker
- **Orchestration:** Kubernetes

## Environment Variables

```bash
# Database Configuration
DB_HOST=localhost
DB_PORT=5432
DB_USERNAME=auth
DB_PASSWORD=auth
DB_NAME=auth_management

# JWT Configuration
JWT_SECRET=aksjmkjdkfiaosjkdjsdoquwqiw

# Cache Configuration
CACHE_HOST=localhost
CACHE_PORT=11211

# Message Broker Configuration
BROKER_HOST=localhost
BROKER_PORT=5672
BROKER_USERNAME=guest
BROKER_PASSWORD=guest
BROKER_VHOST=someone

# OpenTelemetry Configuration
OTLP_HOST=localhost
OTLP_PORT=4317
```

## API Documentation

### 1. User Registration

**Endpoint:** `POST /api/auth/register`

**Request Body:**
```json
{
  "username": "test",
  "password": "test"
}
```

**Response Codes:**
- `201` - Created (User successfully registered)
- `400` - Bad Request (Invalid input data)
- `409` - Conflict (Username already exists)

---

### 2. User Login

**Endpoint:** `POST /api/auth/login`

**Request Body:**
```json
{
  "username": "test",
  "password": "test"
}
```

**Response Body:**
```json
{
    "data": {
        {
        "access_token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
        "refresh_token": "08ca5a94f78c6744",
        }
    },
    "path": "api/auth/login"
}
```

**Response Codes:**
- `200` - OK (Login successful)
- `400` - Bad Request (Invalid input data)
- `401` - Unauthorized (Invalid credentials)

---

### 3. Refresh Access Token

**Endpoint:** `POST /api/auth/token`

**Request Body:**
```json
{
  "refresh_token": "08ca5a94f78c6744"
}
```

**Response Body:**
```json
{
    "data": {
        {
        "access_token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
        "refresh_token": "08ca5a94f78c6744",
        }
    },
    "path": "api/auth/token"
}
```

**Response Codes:**
- `200` - OK (Token refreshed successfully)
- `401` - Unauthorized (Invalid or expired refresh token)

---

## Getting Started

### Prerequisites
- Go 1.21+
- Docker & Docker Compose
- PostgreSQL 15+
- Memcached
- LavinMQ or RabbitMQ

### Installation

1. Clone the repository
```bash
git clone https://github.com/nurmanhadi/go-auth-management.git
cd auth-management
```

2. Set up environment variables
```bash
cp .env.example .env
# Edit .env with your configuration
```

3. Run with Docker Compose
```bash
docker-compose up -d
```

4. Run migrations
```bash
make migrate-up
```

### Development

```bash
# Run locally
go run cmd/main.go

# Run tests
go test ./...

# Build
go build -o bin/auth-service cmd/main.go
```

## Architecture Flow

### Registration Flow
1. User sends registration request to `POST /api/auth/register`
2. Service validates input data
3. Password is hashed using bcrypt
4. User data is stored in PostgreSQL database
5. **Message is published to message broker** (LavinMQ/RabbitMQ) for downstream services (e.g., user id)
6. Returns 201 Created response

### Login Flow
1. User sends login credentials to `POST /api/auth/login`
2. Service validates credentials against database
3. If valid, generates JWT access token and refresh token
4. **Refresh token is stored in Memcached** with TTL (Time To Live)
5. Returns access token and refresh token to client

### Token Refresh Flow
1. Client sends refresh token to `POST /api/auth/token`
2. Service validates refresh token from Memcached
3. If valid, generates new access token
4. Returns new access token to client

## Message Broker Events

### Published Events

#### `user.registered`
Published when a new user successfully registers.

**Payload:**
```json
{
  "event": "user.registered",
  "timestamp": "2025-11-12T10:30:00Z",
  "data": {
    "user_id": "uuid-here",
    "registered_at": "2025-11-12T10:30:00Z"
  }
}
```

**Exchange:** `auth.exchange`  
**Routing Key:** `user.registered`

## Cache Strategy

### Refresh Token Storage
- **Key Pattern:** `refresh:{token_value}`
- **Value:** User ID or session data (JSON)
- **TTL:** 7 days (configurable)
- **Purpose:** Fast validation and session management

### Cache Invalidation
- Refresh tokens are automatically expired after TTL
- Manual invalidation on logout (if implemented)
- Token rotation on refresh

## Security Considerations

- Store `JWT_SECRET` securely (use secrets management in production)
- Use strong passwords for database and broker
- Enable TLS/SSL for production deployments
- Implement rate limiting on authentication endpoints
- Use secure password hashing (bcrypt recommended)
- Set appropriate TTL for refresh tokens in cache
- Implement token rotation strategy
- Use HTTPS for all API endpoints

## License

This project is licensed under the MIT License.

## Author
**Nurman Hadi**  
Backend Developer (Golang, Microservices)  
GitHub: [nurmanhadi](https://github.com/nurmanhadi)