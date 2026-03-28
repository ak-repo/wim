# Auth Structure Guide

This document explains exactly where each auth-related struct belongs in this codebase.

## Goal

Keep auth production-safe by separating concerns clearly:

- HTTP payload structs stay in the HTTP layer
- business structs stay in the service/domain layer
- database row structs stay in the repository layer
- token and password logic stay in the identity layer

## Where Each Auth Struct Belongs

### `internal/http/dto`

Use this folder for API request and response payloads.

Current auth DTOs:
- `RegisterRequest`
- `LoginRequest`
- `RefreshTokenRequest`
- `UserResponse`
- `AuthResponse`

These structs are for:
- JSON decoding
- JSON encoding
- API contract only

These structs should not:
- contain DB tags
- contain SQL null types
- contain password hashes in responses

### `internal/service`

Use this folder for business-use-case inputs and outputs.

Current auth service structs:
- `RegisterInput`
- `LoginInput`
- `RefreshInput`
- `AuthResult`

These structs are for:
- service method inputs
- service workflow outputs
- app-level logic independent from HTTP

These structs should not:
- use JSON tags just because they may later be returned by HTTP
- depend on request bodies directly

### `internal/domain`

Use this folder for core business entities.

Current auth domain structs:
- `User`
- `RefreshToken`
- `Role`

These structs are for:
- core business data
- internal business logic
- system truth inside the application

These structs should not:
- use DB tags
- use JSON tags
- expose transport concerns

### `internal/repository/postgresmodel`

Use this folder for Postgres row models.

Current auth persistence structs:
- `UserRow`
- `RefreshTokenRow`

These structs are for:
- scanning SQL rows
- handling nullable database fields
- matching database column layout

These structs may use:
- `db` tags
- `sql.NullString`
- `sql.NullTime`
- `uuid.NullUUID`

### `internal/identity`

Use this folder for shared identity and auth mechanics.

Current identity structs and contracts:
- `Claims`
- `PasswordHasher`
- `TokenManager`
- `BcryptPasswordHasher`
- `JWTTokenManager`

These types are for:
- password hashing
- token issuing and parsing
- authenticated identity context

### `internal/http/middleware`

Use this folder for HTTP authentication enforcement.

Current middleware behavior:
- parse bearer token
- validate token via `identity.TokenManager`
- place claims into request context

### `internal/http/handler`

Use this folder for auth entrypoints.

Current auth handler responsibilities:
- decode `dto.RegisterRequest`
- convert into `service.RegisterInput`
- call service
- map result into `dto.UserResponse` or `dto.AuthResponse`

## Flow For Register

```text
dto.RegisterRequest
  -> service.RegisterInput
  -> domain.User
  -> postgresmodel.UserRow (inside repository mapping)
  -> database
  -> domain.User
  -> dto.UserResponse
```

## Flow For Login

```text
dto.LoginRequest
  -> service.LoginInput
  -> repository loads postgresmodel.UserRow
  -> repository maps to domain.User
  -> service verifies password and creates tokens
  -> service.AuthResult
  -> dto.AuthResponse
```

## Flow For Refresh

```text
dto.RefreshTokenRequest
  -> service.RefreshInput
  -> repository loads postgresmodel.RefreshTokenRow
  -> repository maps to domain.RefreshToken
  -> service rotates token and issues new access token
  -> service.AuthResult
  -> dto.AuthResponse
```

## Production Safety Rules

- Never return `PasswordHash` in HTTP response structs.
- Never put `sql.Null*` types in domain or HTTP DTOs.
- Never let handlers talk directly to the repository.
- Never let repository return HTTP DTOs.
- Never let public registration assign privileged roles directly.
- Keep access token parsing in middleware, not in every handler.

## Short Rule

Use this memory trick:

```text
HTTP DTO in
service input in
domain inside
postgres row near DB
domain back out
HTTP DTO out
```
