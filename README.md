# Warehouse Inventory Management (WIM)

A RESTful API for managing warehouse inventory, purchase orders, and sales orders.

## Architecture Overview

```
┌─────────────────────────────────────────────────────────────────────────┐
│                           WIM System                                    │
├─────────────────────────────────────────────────────────────────────────┤
│                                                                         │
│  ┌──────────────┐      ┌──────────────┐      ┌──────────────┐         │
│  │   Clients    │──────▶│   REST API   │──────▶│  Services   │         │
│  │  (Frontend)  │      │    (Gin)     │      │  (Business) │         │
│  └──────────────┘      └──────────────┘      └──────┬───────┘         │
│                                                      │                  │
│                       ┌───────────────────────────────┼────────────┐   │
│                       │                               │            │   │
│                       ▼                               ▼            ▼   │
│              ┌────────────────┐           ┌────────────────┐          │
│              │  PostgreSQL     │           │     Redis      │          │
│              │  (Persistence)  │           │    (Cache)     │          │
│              └────────────────┘           └────────────────┘          │
│                                                      │                  │
│                                                      ▼                  │
│              ┌────────────────┐      ┌─────────────────────┐          │
│              │    Kafka       │─────▶│     Worker          │          │
│              │    Queue       │      │  (Async Processing)│          │
│              └────────────────┘      └─────────────────────┘          │
│                                                                         │
└─────────────────────────────────────────────────────────────────────────┘
```

## Technology Stack

- **API Framework**: Gin (Go)
- **Database**: PostgreSQL
- **Cache**: Redis
- **Message Queue**: Kafka
- **Configuration**: Viper

## Data Models

### Core Entities

```
┌─────────────┐     ┌─────────────┐     ┌─────────────┐
│  Product    │     │  Warehouse  │     │  Inventory  │
├─────────────┤     ├─────────────┤     ├─────────────┤
│ SKU         │     │ Code        │     │ ProductID   │
│ Name        │     │ Name        │────▶│ WarehouseID │
│ Description │     │ Address     │     │ LocationID  │
│ Category    │     │ Locations   │     │ BatchID     │
│ Barcode     │     └─────────────┘     │ Quantity    │
└─────────────┘                         │ ReservedQty │
                                        └─────────────┘
                                              │
                         ┌────────────────────┼────────────────────┐
                         │                    │                    │
                         ▼                    ▼                    ▼
                  ┌────────────┐       ┌────────────┐       ┌────────────┐
                  │PurchaseOrder│      │SalesOrder │       │StockMovement│
                  ├────────────┤       ├────────────┤       ├────────────┤
                  │PONumber    │       │OrderNumber│       │Type        │
                  │Supplier   │       │Customer  │        │From/To     │
                  │Items      │       │Items     │        │Quantity    │
                  │Status     │       │Status    │        │Reference   │
                  └────────────┘       └────────────┘       └────────────┘
```

## API Endpoints

### Products
| Method | Endpoint | Description |
|--------|----------|-------------|
| GET | `/api/v1/products` | List all products |
| GET | `/api/v1/products/:id` | Get product by ID |
| POST | `/api/v1/products` | Create new product |
| PUT | `/api/v1/products/:id` | Update product |
| DELETE | `/api/v1/products/:id` | Delete product |

### Warehouses
| Method | Endpoint | Description |
|--------|----------|-------------|
| GET | `/api/v1/warehouses` | List all warehouses |
| GET | `/api/v1/warehouses/:id` | Get warehouse by ID |
| POST | `/api/v1/warehouses` | Create new warehouse |
| PUT | `/api/v1/warehouses/:id` | Update warehouse |
| DELETE | `/api/v1/warehouses/:id` | Delete warehouse |
| GET | `/api/v1/warehouses/:id/locations` | List warehouse locations |
| POST | `/api/v1/warehouses/:id/locations` | Create warehouse location |

### Inventory
| Method | Endpoint | Description |
|--------|----------|-------------|
| GET | `/api/v1/inventory` | List all inventory |
| GET | `/api/v1/inventory/warehouse/:warehouse_id` | Get inventory by warehouse |
| GET | `/api/v1/inventory/product/:product_id` | Get inventory by product |
| POST | `/api/v1/inventory/adjust` | Adjust inventory quantity |

### Purchase Orders
| Method | Endpoint | Description |
|--------|----------|-------------|
| GET | `/api/v1/orders/purchase` | List purchase orders |
| GET | `/api/v1/orders/purchase/:id` | Get purchase order by ID |
| POST | `/api/v1/orders/purchase` | Create purchase order |
| POST | `/api/v1/orders/purchase/:id/receive` | Receive purchase order |

### Sales Orders
| Method | Endpoint | Description |
|--------|----------|-------------|
| GET | `/api/v1/orders/sales` | List sales orders |
| GET | `/api/v1/orders/sales/:id` | Get sales order by ID |
| POST | `/api/v1/orders/sales` | Create sales order |
| POST | `/api/v1/orders/sales/:id/allocate` | Allocate inventory to order |
| POST | `/api/v1/orders/sales/:id/ship` | Ship sales order |

### Health Check
| Method | Endpoint | Description |
|--------|----------|-------------|
| GET | `/api/v1/health` | Check API health |

## Workflows

### 1. Purchase Order Flow (Receiving Goods)

```
┌─────────────┐     ┌─────────────┐     ┌─────────────┐     ┌─────────────┐
│  Create PO  │────▶│  PO Status  │────▶│   Receive   │────▶│  Inventory  │
│             │     │   =PENDING  │     │   Goods     │     │  Updated    │
└─────────────┘     └─────────────┘     └─────────────┘     └─────────────┘
                                                │
                                                ▼
                                    ┌─────────────────────┐
                                    │ StockMovement       │
                                    │ (Type: RECEIPT)     │
                                    └─────────────────────┘
```

**Steps:**
1. Create a Purchase Order with items
2. When goods arrive, use `/receive` endpoint
3. System creates inventory records for each item
4. Stock movement is recorded

---

### 2. Sales Order Flow (Shipping Goods)

```
┌─────────────┐     ┌─────────────┐     ┌─────────────┐     ┌─────────────┐
│  Create SO  │────▶│  Allocate   │────▶│   Ship      │────▶│  Inventory  │
│             │     │  Inventory  │     │   Order     │     │  Decreased  │
└─────────────┘     └─────────────┘     └─────────────┘     └─────────────┘
       │                   │                   │
       ▼                   ▼                   ▼
┌─────────────┐     ┌─────────────┐     ┌─────────────────────┐
│ SO Status:  │     │ SO Status:  │     │ StockMovement       │
│ PENDING     │     │ PROCESSING  │     │ (Type: SHIP)       │
└─────────────┘     └─────────────┘     └─────────────────────┘
```

**Steps:**
1. Create a Sales Order
2. Allocate inventory (reserves stock)
3. Ship order (reduces actual inventory)
4. Stock movement is recorded

---

### 3. Inventory Adjustment Flow

```
┌─────────────┐     ┌─────────────┐     ┌─────────────┐
│   Client    │────▶│   Adjust    │────▶│  Inventory  │
│   Request   │     │   API       │     │  Updated    │
└─────────────┘     └─────────────┘     └─────────────┘
                                               │
                                               ▼
                                    ┌─────────────────────┐
                                    │ StockMovement       │
                                    │ (Type: ADJUSTMENT)  │
                                    └─────────────────────┘
```

---

### 4. Async Processing (Kafka Worker)

```
┌──────────┐     ┌──────────┐     ┌──────────┐     ┌──────────┐
│  Kafka   │────▶│  Worker │────▶│  Topic   │────▶│  Action  │
│  Message │     │  Pool   │     │  Handler │     │          │
└──────────┘     └──────────┘     └──────────┘     └──────────┘
                                                         │
                                    ┌────────────────────┤
                                    │                    │
                                    ▼                    ▼
                          ┌────────────────┐    ┌────────────────┐
                          │ stock_recalc   │    │  expiry_alert  │
                          └────────────────┘    └────────────────┘
```

**Supported Topics:**
- `inventory.stock_recalculation` - Recalculate stock levels
- `inventory.expiry_alert` - Alert for expiring inventory
- `reports.generation` - Generate reports

---

## Project Structure

```
warehouse-inventory/
├── cmd/
│   ├── api/
│   │   └── main.go          # API server entry point
│   └── worker/
│       └── main.go           # Kafka worker entry point
├── internal/
│   ├── config/
│   │   └── config.go         # Configuration loading
│   ├── domain/
│   │   └── domain.go         # Domain models
│   ├── event/
│   │   └── event.go          # Event definitions
│   ├── handler/
│   │   ├── inventory/        # Inventory handlers
│   │   ├── order/            # Order handlers
│   │   ├── product/          # Product handlers
│   │   └── warehouse/        # Warehouse handlers
│   ├── middleware/
│   │   └── middleware.go     # HTTP middleware
│   ├── queue/
│   │   └── processor.go      # Kafka message processor
│   ├── repository/
│   │   └── postgres/         # Database repositories
│   ├── service/
│   │   ├── inventory/        # Inventory business logic
│   │   ├── order/            # Order business logic
│   │   ├── product/          # Product business logic
│   │   └── warehouse/        # Warehouse business logic
│   └── worker/
│       ├── errors.go         # Worker errors
│       └── pool.go           # Worker pool
├── pkg/
│   ├── errors/
│   │   └── errors.go         # Error definitions
│   ├── logger/
│   │   └── logger.go         # Logging utility
│   └── validator/
│       └── validator.go      # Input validation
├── go.mod
├── go.sum
└── config.yaml               # Configuration file
```

## Configuration

Create a `config.yaml` file:

```yaml
server:
  port: 8080

database:
  host: "localhost"
  port: 5432
  user: "postgres"
  password: "password"
  database: "wim"
  sslmode: "disable"
  max_conns: 20

redis:
  host: "localhost"
  port: 6379
  password: ""
  db: 0

kafka:
  brokers:
    - "localhost:9092"
  topic: "wim-events"
  group_id: "wim-worker"

worker:
  pool_size: 5
  queue_size: 100

loglevel: "info"
```

## Running the Application

### Prerequisites
- Go 1.21+
- PostgreSQL
- Redis
- Kafka

### Start API Server
```bash
go run cmd/api/main.go
```

### Start Worker
```bash
go run cmd/worker/main.go
```

## Example Requests

### Create Product
```bash
curl -X POST http://localhost:8080/api/v1/products \
  -H "Content-Type: application/json" \
  -d '{
    "sku": "PROD-001",
    "name": "Sample Product",
    "description": "A sample product",
    "category": "Electronics",
    "unitOfMeasure": "pcs",
    "barcode": "123456789"
  }'
```

### Create Warehouse
```bash
curl -X POST http://localhost:8080/api/v1/warehouses \
  -H "Content-Type: application/json" \
  -d '{
    "code": "WH-001",
    "name": "Main Warehouse",
    "addressLine1": "123 Main St",
    "city": "New York",
    "state": "NY",
    "postalCode": "10001",
    "country": "USA"
  }'
```

### Create Purchase Order
```bash
curl -X POST http://localhost:8080/api/v1/orders/purchase \
  -H "Content-Type: application/json" \
  -d '{
    "supplierId": "uuid-here",
    "warehouseId": "uuid-here",
    "totalAmount": 1000.00,
    "notes": "Sample PO"
  }'
```

### Receive Purchase Order
```bash
curl -X POST http://localhost:8080/api/v1/orders/purchase/{id}/receive \
  -H "Content-Type: application/json" \
  -d '{
    "quantity": 100,
    "location": "uuid-here"
  }'
```

### Create Sales Order
```bash
curl -X POST http://localhost:8080/api/v1/orders/sales \
  -H "Content-Type: application/json" \
  -d '{
    "customerId": "uuid-here",
    "warehouseId": "uuid-here",
    "shippingMethod": "Express",
    "shippingAddress": "123 Customer St",
    "totalAmount": 500.00
  }'
```

## License

MIT
