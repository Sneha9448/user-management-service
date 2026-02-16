# User Management Microservice

A production-ready REST API in Go for managing users, using PostgreSQL (pgx) and Gorilla Mux.

## Features
- **RESTful API**: Full CRUD operations for user management.
- **GraphQL API**: Query and mutate users via `/graphql`.
- **PostgreSQL**: Robust connection pooling with `pgx`.
- **Configuration**: Environment variable based configuration.
- **Reliability**: Context timeouts on database operations (5s).
- **Architecture**: Clean, modular structure (Handlers, Repository, Models).
- **Optimizations**: Detailed performance tuning at DB and API layers. See [OPTIMIZATION.md](OPTIMIZATION.md) for details.

## Prerequisites

- Go 1.21+
- PostgreSQL 12+
- Postman (for API testing)

## Setup & Running

### 1. Database Setup
Ensure you have a PostgreSQL database running.

#### Run Migrations
Execute the SQL in `migrations/001_create_users_table.sql` on your database. You can use a tool like `psql` or a GUI (pgAdmin, DBeaver).

```sql
-- Example using psql
psql "postgres://postgres:postgres@localhost:5432/postgres" -f migrations/001_create_users_table.sql
```

### 2. Configuration
Set the following environment variables. You can set them in your terminal or create a `.env` file loader (not included by default) or just export them.

**Linux/Mac:**
```bash
export DATABASE_URL="postgres://user:password@localhost:5432/dbname"
export PORT="8080"
```

**Windows (PowerShell):**
```powershell
$env:DATABASE_URL="postgres://user:password@localhost:5432/dbname"
$env:PORT="8080"
```

*Note: The default `DATABASE_URL` is `postgres://postgres:postgres@localhost:5432/postgres?sslmode=disable` if not specified.*

### 3. Start the Service
Run the application using the Go CLI:

```bash
go run cmd/server/main.go
```

You should see:
```
Connected to PostgreSQL successfully
Server listening on port 8080
```

## Testing with Postman

1. Open Postman.
2. Click **Import** -> **Upload Files**.
3. Select `postman_collection.json` from this directory.
4. You will see nine requests:
   - **Health Check**: `GET localhost:8080/health`
   - **Create User (REST)**: `POST localhost:8080/users`
   - **Get User (REST)**: `GET localhost:8080/users/1`
   - **Update User (REST)**: `PUT localhost:8080/users/1`
   - **Delete User (REST)**: `DELETE localhost:8080/users/1`
   - **Get All Users (GraphQL)**: `POST localhost:8080/graphql`
   - **Create User (GraphQL)**: `POST localhost:8080/graphql`
   - **Update User (GraphQL)**: `POST localhost:8080/graphql`
   - **Delete User (GraphQL)**: `POST localhost:8080/graphql`

## REST API Endpoints

### 1. Health Check
- **URL**: `/health`
- **Method**: `GET`
- **Response**: `200 OK`
  ```json
  { "status": "healthy" }
  ```

### 2. Create User
- **URL**: `/users`
- **Method**: `POST`
- **Body**:
  ```json
  {
      "name": "Jane Doe",
      "email": "jane@example.com"
  }
  ```
- **Response**: `201 Created`

### 3. Get User
- **URL**: `/users/{id}`
- **Method**: `GET`
- **Response**: `200 OK`

### 4. Update User
- **URL**: `/users/{id}`
- **Method**: `PUT`
- **Body**:
  ```json
  {
      "name": "Updated Name",
      "email": "updated@example.com"
  }
  ```
- **Response**: `200 OK`

### 5. Delete User
- **URL**: `/users/{id}`
- **Method**: `DELETE`
- **Response**: `204 No Content`

## GraphQL API Endpoints

All GraphQL requests are sent to `/graphql` via `POST`.

### 1. Get All Users (Query)
**Query:**
```graphql
query {
  users {
    id
    name
    email
  }
}
```

### 2. Create User (Mutation)
**Mutation:**
```graphql
mutation {
  createUser(name: "GraphQL User", email: "graphql@example.com") {
    id
    name
    email
  }
}
```

### 3. Update User (Mutation)
**Mutation:**
```graphql
mutation {
  updateUser(id: 1, name: "Updated GraphQL", email: "updatedgql@example.com") {
    id
    name
    email
  }
}
```

### 4. Delete User (Mutation)
**Mutation:**
```graphql
mutation {
  deleteUser(id: 1)
}
```
