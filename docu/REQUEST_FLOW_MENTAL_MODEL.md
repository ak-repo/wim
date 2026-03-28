# Request Flow Mental Model

This document explains the request flow for this codebase in a simple mental-model style.

## Main Idea

One API request should move through the system like this:

```text
Client
  -> HTTP Router
  -> HTTP Handler
  -> Request DTO
  -> Service Input
  -> Service Logic
  -> Repository Interface
  -> Postgres Repository
  -> Postgres Row Model
  -> Database
  -> Postgres Row Model
  -> Domain Model
  -> Service Output
  -> Response DTO
  -> HTTP Response
  -> Client
```

## Layer Responsibility

### 1. Router
- File area: `internal/http/router`
- Job: connect URL + HTTP method to a handler
- Example: `POST /api/v1/auth/register`

Router should only decide:
- which endpoint is called
- which middleware runs
- which handler should execute

Router should not:
- contain business logic
- talk to database
- build SQL queries

### 2. Handler
- File area: `internal/http/handler`
- Job: HTTP boundary logic

Handler should:
- read JSON body
- decode into request struct from `internal/http/dto`
- validate basic request shape
- call service with service input
- map service/domain result into response DTO
- return HTTP status code + JSON

Handler should not:
- contain business rules
- talk directly to DB
- build repository queries

### 3. Request DTO
- File area: `internal/http/dto`
- Example: `RegisterRequest`, `CreateProductRequest`
- Job: represent incoming API payload

This layer is HTTP-only.

Example:

```text
POST /api/v1/auth/register
body -> dto.RegisterRequest
```

### 4. Service Input
- File area: `internal/service`
- Example: `RegisterInput`, `LoginInput`
- Job: represent business-use-case input

This layer is not tied to HTTP.

Reason:
- service should work even if tomorrow input comes from CLI, worker, or gRPC instead of REST

### 5. Service
- File area: `internal/service`
- Job: business logic and use cases

Service should:
- validate business rules
- call repository interfaces
- call identity/auth helpers
- coordinate transactions/workflows
- return domain result or service result

Service should not:
- know JSON or HTTP details
- know SQL query shapes
- know `sql.NullString`

### 6. Repository Interface
- File area: `internal/repository`
- Job: define persistence behavior needed by service

Example:

```text
UserRepository.Create(...)
UserRepository.GetByEmail(...)
```

This gives service a clean contract without depending on Postgres details.

### 7. Repository Implementation
- File area: `internal/repository` and later concrete Postgres files
- Job: execute DB operations

Repository implementation should:
- run SQL
- scan rows
- map DB row structs to domain models
- handle persistence-specific details

Repository should not:
- return HTTP responses
- contain route logic

### 8. Postgres Row Model
- File area: `internal/repository/postgresmodel`
- Job: represent DB scan/write shapes

Use this layer for:
- `db` tags
- `sql.NullString`
- `sql.NullTime`
- `uuid.NullUUID`

This layer exists because DB shape is often different from domain shape.

### 9. Domain Model
- File area: `internal/domain`
- Job: business entity representation

Examples:
- `User`
- `Product`
- `Warehouse`
- `Inventory`

Domain should stay clean:
- no JSON tags
- no DB tags
- no HTTP concerns

### 10. Response DTO
- File area: `internal/http/dto`
- Job: represent API output

Example:
- `UserResponse`
- `ProductResponse`
- `AuthResponse`

Important:
- never expose internal fields like `PasswordHash`

## Real Example: Register User Flow

```text
1. Client sends POST /api/v1/auth/register
2. Router sends request to AuthHandler.Register
3. Handler decodes JSON into dto.RegisterRequest
4. Handler converts dto.RegisterRequest -> service.RegisterInput
5. Service validates business rules
6. Service calls UserRepository.ExistsByEmail
7. Service calls identity.PasswordHasher.Hash
8. Service creates domain.User
9. Service calls UserRepository.Create
10. Repository converts domain.User -> postgresmodel.UserRow if needed
11. Repository inserts into PostgreSQL
12. Repository returns domain.User
13. Service returns domain.User or auth result
14. Handler converts domain.User -> dto.UserResponse
15. Handler writes HTTP 201 JSON response
```

## Real Example: Login Flow

```text
1. Client sends POST /api/v1/auth/login
2. Router -> AuthHandler.Login
3. Handler decodes dto.LoginRequest
4. Handler converts dto.LoginRequest -> service.LoginInput
5. Service loads user through UserRepository.GetByEmail
6. Service verifies password through identity.PasswordHasher.Compare
7. Service issues token through identity.TokenManager
8. Service returns AuthResult
9. Handler maps AuthResult -> dto.AuthResponse
10. Handler writes HTTP 200 JSON response
```

## Mapping Rules

Keep mapping direction like this:

```text
HTTP request DTO -> service input
service input -> domain model
domain model -> postgres row model
postgres row model -> domain model
domain model -> response DTO
```

Do not do this:
- repository row -> API response directly
- handler -> repository directly
- service -> JSON DTO directly

## Simple Folder Mental Model

```text
internal/domain
  business truth

internal/http/dto
  API request and response shapes

internal/http/handler
  HTTP in/out boundary

internal/http/router
  endpoint wiring

internal/service
  business workflow and use cases

internal/repository
  persistence contracts

internal/repository/postgresmodel
  database row shapes

internal/identity
  password, token, claims, auth helpers
```

## Golden Rules

- Handler knows HTTP.
- Service knows business.
- Repository knows persistence.
- Domain knows core entities.
- DTO knows API payload shape.
- Postgres model knows DB row shape.

## Short Memory Trick

Think of the flow like this:

```text
request comes in as DTO
DTO becomes service input
service works on domain
repository stores through row model
result comes back as domain
domain becomes response DTO
response goes out
```

That is the safe and scalable pattern for this project.
