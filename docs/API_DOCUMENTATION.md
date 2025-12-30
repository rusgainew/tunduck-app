# Tunduc API System - API Documentation

## Overview

Tunduc API System is an enterprise-grade API for managing ESF (Electronic Waybill) documents, organizations, and users. It features:

- **Authentication**: JWT-based user registration and login with secure logout
- **Caching**: Redis-backed caching for improved performance
- **Rate Limiting**: DDoS protection with configurable rate limits
- **Monitoring**: Prometheus metrics and health checks
- **Testing**: Comprehensive unit and integration tests

## Quick Start

### Base URL

- Development: `http://localhost:8080`
- Production: `https://api.example.com`

### Authentication

Most endpoints require authentication using JWT tokens. Include the token in the `Authorization` header:

```bash
Authorization: Bearer <your_jwt_token>
```

## API Endpoints

### Authentication Endpoints

#### 1. User Registration

**Endpoint**: `POST /api/auth/register`

**Description**: Create a new user account

**Rate Limit**: 30 requests/minute per IP

**Request Body**:

```json
{
  "username": "johndoe",
  "email": "john@example.com",
  "password": "SecurePassword123!"
}
```

**Validation Rules**:

- Username: 3-50 characters
- Email: Valid email format
- Password: Minimum 8 characters

**Success Response** (201 Created):

```json
{
  "id": 1,
  "username": "johndoe",
  "email": "john@example.com"
}
```

**Error Responses**:

- 400 Bad Request: Invalid input or validation failed
- 409 Conflict: User already exists
- 429 Too Many Requests: Rate limit exceeded

**Example**:

```bash
curl -X POST http://localhost:8080/api/auth/register \
  -H "Content-Type: application/json" \
  -d '{
    "username": "johndoe",
    "email": "john@example.com",
    "password": "SecurePassword123!"
  }'
```

---

#### 2. User Login

**Endpoint**: `POST /api/auth/login`

**Description**: Authenticate user and receive JWT token

**Rate Limit**: 30 requests/minute per IP

**Request Body**:

```json
{
  "email": "john@example.com",
  "password": "SecurePassword123!"
}
```

**Success Response** (200 OK):

```json
{
  "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
  "user": {
    "id": 1,
    "username": "johndoe",
    "email": "john@example.com"
  }
}
```

**Error Responses**:

- 401 Unauthorized: Invalid email or password
- 429 Too Many Requests: Rate limit exceeded

**Example**:

```bash
curl -X POST http://localhost:8080/api/auth/login \
  -H "Content-Type: application/json" \
  -d '{
    "email": "john@example.com",
    "password": "SecurePassword123!"
  }'
```

---

#### 3. User Logout

**Endpoint**: `POST /api/auth/logout`

**Description**: Invalidate JWT token and add to blacklist

**Authentication**: Required (Bearer token)

**Rate Limit**: 5 requests/minute per user

**Success Response** (200 OK):

```json
{
  "message": "Logged out successfully"
}
```

**Error Responses**:

- 401 Unauthorized: Missing or invalid token
- 429 Too Many Requests: Rate limit exceeded

**Example**:

```bash
curl -X POST http://localhost:8080/api/auth/logout \
  -H "Authorization: Bearer YOUR_JWT_TOKEN"
```

---

### System Endpoints

#### 4. Health Check

**Endpoint**: `GET /health`

**Description**: Check system health including PostgreSQL and Redis

**Rate Limit**: 120 requests/minute per IP

**Response** (200 OK - All systems UP):

```json
{
  "status": "UP",
  "timestamp": "2025-12-28T10:30:00Z",
  "uptime": "2h15m30s",
  "components": [
    {
      "name": "PostgreSQL",
      "status": "UP",
      "response_time": "3ms",
      "message": "Database connected successfully",
      "last_checked": "2025-12-28T10:30:00Z"
    },
    {
      "name": "Redis",
      "status": "UP",
      "response_time": "1ms",
      "message": "Redis connected successfully",
      "last_checked": "2025-12-28T10:30:00Z"
    }
  ]
}
```

**Error Response** (503 Service Unavailable - System DOWN):

```json
{
  "status": "DOWN",
  "timestamp": "2025-12-28T10:30:00Z",
  "uptime": "1h45m20s",
  "components": [
    {
      "name": "PostgreSQL",
      "status": "DOWN",
      "response_time": "5000ms",
      "message": "Failed to connect: connection timeout",
      "last_checked": "2025-12-28T10:30:00Z"
    }
  ]
}
```

**Example**:

```bash
curl http://localhost:8080/health
```

**Use Cases**:

- Kubernetes liveness/readiness probes
- Load balancer health monitoring
- Monitoring dashboards
- Auto-scaling triggers

---

#### 5. Prometheus Metrics

**Endpoint**: `GET /metrics`

**Description**: Prometheus-compatible metrics endpoint

**Rate Limit**: 180 requests/minute per IP

**Content Type**: `text/plain`

**Metrics Included**:

- HTTP metrics (requests, duration, size)
- Cache metrics (hits, misses, operations)
- Database metrics (query time, errors)
- Auth metrics (login, logout attempts)
- Business metrics (users registered, documents created)

**Example**:

```bash
curl http://localhost:8080/metrics | head -20
```

---

## Rate Limiting

The API implements rate limiting to protect against DDoS attacks and abuse:

### Rate Limit Categories

| Endpoint             | Limit       | Window   |
| -------------------- | ----------- | -------- |
| `/api/auth/register` | 30 req/min  | 1 minute |
| `/api/auth/login`    | 30 req/min  | 1 minute |
| `/api/auth/logout`   | 5 req/min   | 1 minute |
| `/health`            | 120 req/min | 1 minute |
| `/metrics`           | 180 req/min | 1 minute |

### Rate Limit Response

When rate limit is exceeded, the API returns HTTP 429:

```json
{
  "error": "Too Many Requests",
  "message": "Rate limit exceeded. Please try again later.",
  "reset": 1735382400
}
```

The `reset` field contains a Unix timestamp indicating when the limit resets.

**Handling Rate Limits**:

```bash
# Check rate limit in response headers
curl -i http://localhost:8080/api/auth/login \
  -H "Content-Type: application/json" \
  -d '{"email":"test@example.com","password":"password"}'
```

---

## Error Handling

### Standard Error Response

All errors follow a consistent format:

```json
{
  "error": "Error Type",
  "message": "Detailed error message",
  "timestamp": "2025-12-28T10:30:00Z"
}
```

### HTTP Status Codes

| Code | Meaning               | Example                    |
| ---- | --------------------- | -------------------------- |
| 200  | OK                    | Request successful         |
| 201  | Created               | Resource created           |
| 400  | Bad Request           | Invalid input format       |
| 401  | Unauthorized          | Missing/invalid token      |
| 409  | Conflict              | Resource already exists    |
| 429  | Too Many Requests     | Rate limit exceeded        |
| 500  | Internal Server Error | Server error               |
| 503  | Service Unavailable   | System down (health check) |

---

## Caching

The API uses Redis for caching to improve performance:

### Cache Behavior

- **User Cache**: 1-hour TTL
  - Cached on first login/registration lookup
  - Invalidated on profile changes
- **Organization Cache**: 2-hour TTL
  - Cached on first access
  - Automatic warming on startup
- **Document Cache**: 30-minute TTL
  - Cached on first retrieval
  - Fast subsequent accesses

### Cache Warming

On application startup, frequently accessed data is preloaded into cache for instant availability.

---

## Authentication Details

### JWT Token Structure

Tokens are JWT format with three parts:

```
eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiJ1c2VyMTIzIiwiZW1haWwiOiJqb2huQGV4YW1wbGUuY29tIiwiaWF0IjoxNzM1MjY2MDAwfQ.signature
```

**Token Expiration**: 24 hours (configurable)

### Token Blacklist

When you logout, the token is immediately added to a Redis blacklist, preventing its use for future requests.

---

## Examples

### Complete Authentication Flow

```bash
# 1. Register a new user
curl -X POST http://localhost:8080/api/auth/register \
  -H "Content-Type: application/json" \
  -d '{
    "username": "alice",
    "email": "alice@example.com",
    "password": "SecurePassword123!"
  }'

# Response:
# {
#   "id": 42,
#   "username": "alice",
#   "email": "alice@example.com"
# }

# 2. Login
curl -X POST http://localhost:8080/api/auth/login \
  -H "Content-Type: application/json" \
  -d '{
    "email": "alice@example.com",
    "password": "SecurePassword123!"
  }'

# Response:
# {
#   "token": "eyJhbGciOiJIUzI1NiIs...",
#   "user": {
#     "id": 42,
#     "username": "alice",
#     "email": "alice@example.com"
#   }
# }

# 3. Use token for protected endpoint
TOKEN="eyJhbGciOiJIUzI1NiIs..."
curl -X GET http://localhost:8080/api/user/profile \
  -H "Authorization: Bearer $TOKEN"

# 4. Logout
curl -X POST http://localhost:8080/api/auth/logout \
  -H "Authorization: Bearer $TOKEN"

# Response:
# {
#   "message": "Logged out successfully"
# }
```

### Monitor System Health

```bash
# Check health
curl http://localhost:8080/health | jq .

# For Kubernetes probe configuration:
# livenessProbe:
#   httpGet:
#     path: /health
#     port: 8080
#   initialDelaySeconds: 10
#   periodSeconds: 10
#
# readinessProbe:
#   httpGet:
#     path: /health
#     port: 8080
#   initialDelaySeconds: 5
#   periodSeconds: 5
```

---

## API Documentation

Interactive Swagger UI documentation is available at:

```
http://localhost:8080/swagger/index.html
```

The Swagger UI provides:

- Interactive endpoint testing
- Request/response examples
- Parameter descriptions
- Error code documentation

---

## Development

### Running Tests

```bash
# Run all tests
go test ./...

# Run with coverage
go test -cover ./...

# Run specific test
go test -run TestJWTToken ./pkg/auth/...

# Integration tests only
go test -v ./...  # includes integration tests that connect to Redis/DB
```

### Building

```bash
# Build binary
go build -o api ./cmd/api

# Build with version info
go build -ldflags="-X main.Version=1.0.0" -o api ./cmd/api
```

---

## Performance

### Benchmarks

- **Login endpoint**: ~50-100ms (including database query)
- **Cached user lookup**: ~1-5ms (Redis cache hit)
- **Health check**: ~5ms (two parallel probes with 5s timeout)
- **Rate limiter check**: <1ms (Redis atomic operation)

### Optimization Tips

1. **Use caching headers**: Response includes cache information
2. **Batch requests**: Reduce number of API calls
3. **Connection pooling**: Database connections are pooled
4. **Monitor metrics**: Use `/metrics` endpoint with Prometheus

---

## Support

For issues, questions, or contributions:

- GitHub: https://github.com/rusgainew/tunduck-app
- Email: support@example.com
- Documentation: See `/swagger/index.html` for interactive API docs
