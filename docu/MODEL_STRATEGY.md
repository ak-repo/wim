# Model Strategy

This project should keep models separated by boundary, not mixed inside one file.

## Folder-By-Folder Strategy

### `internal/domain`
- Keep core business entities only.
- No JSON tags.
- No DB tags.
- No `sql.Null*` fields unless the domain truly needs them.
- Safe for service/business logic.

Examples:
- `User`
- `Product`
- `Warehouse`
- `Inventory`
- `PurchaseOrder`
- `SalesOrder`
- `Transfer`

### `internal/repository/postgresmodel`
- Keep Postgres row structs only.
- Use `db` tags.
- Use `sql.Null*` and `uuid.NullUUID` here when DB nullability requires it.
- Safe for query scan/write operations.

Examples:
- `UserRow`
- `ProductRow`
- `WarehouseRow`
- `InventoryRow`
- `PurchaseOrderRow`

### `internal/http/dto`
- Keep request/response structs only.
- Use JSON tags only.
- Never expose secrets like `PasswordHash`.
- This is the correct place for POST request structs.

Examples:
- `RegisterRequest`
- `CreateProductRequest`
- `CreateWarehouseRequest`
- `CreatePurchaseOrderRequest`
- `CreateSalesOrderRequest`
- `CreateTransferRequest`
- `UserResponse`
- `ProductResponse`

### `internal/service`
- Use domain models and service input/output structs.
- Business-use-case input structs are allowed here when they are not HTTP-specific.
- Services should not depend on HTTP DTOs directly.

### `internal/http/handler`
- Decode request DTOs.
- Convert request DTOs to service inputs.
- Convert domain/service results to response DTOs.

### `internal/repository`
- Define repository interfaces.
- Return domain models to the service layer.
- Concrete Postgres code can map `postgresmodel.*Row` to `domain.*`.

## Recommended Mapping Direction

Use this direction consistently:

```text
HTTP request DTO -> service input -> domain -> repository row
repository row -> domain -> response DTO
```

## Rule For Future Modules

- Do not create three structs automatically for every table.
- Create separate structs only when the boundary actually differs.
- If DB shape and domain shape are identical enough, repository may scan straight into domain.
- If API shape differs from domain, create response DTOs in `internal/http/dto`.

## Module Coverage

This repository now has struct placeholders for:

- identity/auth
- users
- products
- warehouses
- locations
- inventory
- stock movements
- purchase orders
- sales orders
- transfers
- batches
- reports
- audit logs

This keeps the project production-safe as modules grow because each layer owns its own representation.
