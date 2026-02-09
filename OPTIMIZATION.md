# Performance Optimization Techniques

This document details the performance optimization techniques implemented in the User Management Service to ensure high throughput, low latency, and efficient resource management.

## 1. Database Layer (PostgreSQL)

The service uses the `pgxpool` package from the `pgx/v5` driver for advanced connection pooling.

### Connection Pooling
Efficient database connection management is critical for performance. Instead of opening and closing a connection for every request, we maintain a pool of reusable connections.

- **`MaxConns` (25)**: Limits the maximum number of concurrent database connections. This prevents the service from overwhelming the database server.
- **`MinConns` (5)**: Ensures a minimum number of connections are always ready, reducing latency for the first few requests after a period of inactivity.
- **`MaxConnLifetime` (1 hour)**: Periodically recycles connections to prevent memory leaks and handle stale connections that might have been dropped by the network/database.
- **`MaxConnIdleTime` (10 minutes)**: Closes connections that have been idle for too long, freeing up resources on the database server.

**Implementation Detail**:
```go
config.MaxConns = 25
config.MinConns = 5
config.MaxConnLifetime = time.Hour
config.MaxConnIdleTime = 10 * time.Minute
```
*Location: [internal/database/postgres.go](file:///c:/Users/ksneh/.gemini/antigravity/scratch/user-management-service/internal/database/postgres.go)*

---

## 2. Service & Repository Layer

### Request Timeouts & Deadlines
Every database operation is wrapped in a `context.WithTimeout`. This ensures that if the database is slow or unresponsive, the service doesn't hang indefinitely, preventing resource exhaustion (like goroutine leaks).

- **Standard Timeout (5s)**: A strict 5-second limit is placed on all SQL queries.

**Implementation Detail**:
```go
ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
defer cancel()
// ... database operation using ctx ...
```
*Location: [internal/repository/user_repository.go](file:///c:/Users/ksneh/.gemini/antigravity/scratch/user-management-service/internal/repository/user_repository.go)*

---

## 3. API Layer (Optimization & Monitoring)

### Performance Profiling (pprof)
The service includes built-in support for Go's `pprof` tool. This allows real-time analysis of CPU usage, memory allocation, and goroutine counts.

- **Endpoint**: `http://localhost:6060/debug/pprof/`
- **Tracing**: Supports execution tracing to identify bottlenecks in request handling.

### Efficient Data Handling
- **GraphQL**: Reduces over-fetching and under-fetching by allowing clients to request exactly what they need.
- **CORS Management**: Properly configured CORS to handle preflight requests efficiently.

---

---

## Performance Metrics (Recorded on 2026-02-09)

The following metrics were captured using built-in Go benchmarks (`go test -bench . -benchmem`) against a local PostgreSQL instance.

| Operation | Throughput (Iterations) | Latency (ns/op) | Memory (B/op) | Allocations (allocs/op) |
| :--- | :--- | :--- | :--- | :--- |
| **CreateUser** | 6,302 | 198,851 (~0.2ms) | 1,089 | 23 |
| **GetAllUsers** | 279 | 5,225,482 (~5.2ms) | 2,366,350 | 51,168 |

### Observation:
- **Write Performance**: The 0.2ms latency for `CreateUser` demonstrates efficient connection pooling and minimal overhead in the repository layer.
- **Read Scalability**: `GetAllUsers` performance reflects the data volume in the table. The memory allocations are consistent with the data structures being hydrated.

---

*Note: These benchmarks were run using the dedicated performance test suite in [internal/repository/performance_test.go](file:///c:/Users/ksneh/.gemini/antigravity/scratch/user-management-service/internal/repository/performance_test.go).*
