# Warehouse Inventory Management System

## Technical Design Document

---

# PHASE 1 — SYSTEM UNDERSTANDING

## 1.1 Product Lifecycle

Real-world warehouse systems manage products through distinct stages:

1. **Procurement** – Product is identified for purchase, supplier is selected
2. **Purchase Order Creation** – PO generated with expected quantities, prices, delivery dates
3. **Receiving** – Goods arrive at dock and are counted and inspected for quality
4. **Put-away** – Items moved from receiving area to storage locations
5. **Storage** – Items held in bins/racks until needed
6. **Allocation** – Items reserved for specific orders
7. **Picking** – Items retrieved from storage locations
8. **Packing** – Items packed into shipping containers
9. **Shipping** – Goods dispatched to customers
10. **Returns** – Defective or unsold items returned to inventory

---

## 1.2 Warehouse Structure

Typical warehouse hierarchy:

* **Warehouse** → Top-level facility (example: `DC-NewYork`)
* **Zone** → Logical area (Receiving, Storage, Shipping)
* **Aisle** → Row within zone
* **Rack** → Storage unit inside aisle
* **Bin / Location** → Smallest storage unit (`A-01-03-02`)

### Location Types

* **Fixed** – Permanent storage positions
* **Dynamic** – Pooled storage
* **Pick-face** – Front picking locations
* **Bulk** – Secondary storage for replenishment

---

## 1.3 Stock Movements

Every inventory change creates a **movement record**.

Movement Types:

* RECEIPT
* PUTAWAY
* PICK
* PACK
* SHIP
* TRANSFER_IN
* TRANSFER_OUT
* ADJUSTMENT
* RETURN
* DAMAGE
* EXPIRY

---

## 1.4 Barcode Scanning Workflow

### Barcode Types

* **UPC / EAN** → Product identification
* **SSCC** → Serial Shipping Container Code (18 digits)
* **Location barcode** → Warehouse location code
* **Batch / Lot barcode** → Production batch identifier
* **Pallet barcode** → Pallet identifier

### Scanning Operations

Receiving:

1. Scan PO
2. Scan product
3. Scan quantity
4. Confirm

Put-away:

1. Scan item
2. Scan location
3. Confirm placement

Picking:

1. Scan order
2. Scan location
3. Scan item
4. Confirm quantity

Packing:

1. Scan items
2. Scan carton
3. Scan shipping label

Cycle count:

1. Scan location
2. Count items
3. Submit variance

---

## 1.5 Purchase Order Flow

1. Create PO
2. Send to Supplier
3. Supplier Confirms
4. Delivery Date Set
5. Goods Arrive
6. Receiving

Receiving steps:

* Scan PO number
* Scan product barcode
* Enter received quantity
* Perform quality inspection
* Record damages or shortages

7. Partial or Full Receipt
8. Put-away

Put-away steps:

* Generate tasks
* Scan location
* Scan items
* Confirm placement

9. Inventory Updated
10. Invoice Matching

---

## 1.6 Sales Order Flow

1. Order received
2. Validate order
3. Check inventory
4. Allocate stock
5. Generate pick tasks

Picking:

* Scanner directs picker to location
* Scan location
* Scan product
* Confirm quantity

Packing:

* Consolidate items
* Pack cartons
* Record dimensions

Shipping:

* Assign carrier
* Generate tracking number
* Print labels
* Create manifest

Final Steps:

* Shipment confirmed
* Inventory deducted
* Invoice generated

---

## 1.7 Stock Transfer Between Warehouses

1. Transfer request created
2. Source warehouse approval
3. Stock reserved
4. Picking at source
5. Packing for transit
6. Shipping
7. Transit tracking
8. Receiving at destination
9. Quality check
10. Put-away
11. Transfer complete

---

## 1.8 Batch Tracking

Purpose: **Traceability**

Fields:

* Batch Number
* Manufacturing Date
* Expiry Date
* Supplier Lot
* Origin Country

Use Cases:

* Product recall
* Quality control
* Regulatory compliance

---

## 1.9 Expiry Tracking

Two strategies:

**FIFO** – First In First Out
**FEFO** – First Expired First Out

Expiry alerts:

* 30 days → warning
* 7 days → critical
* Expired → blocked from picking

---

## 1.10 Stock Audits

### Cycle Counting

Continuous counting system:

* ABC classification
* Random verification
* Variance investigation

### Physical Inventory

* Full warehouse shutdown counting
* Compare physical vs system inventory
* Post adjustments

---

# PHASE 2 — SYSTEM ARCHITECTURE

## 2.1 Recommended Go Project Structure

```
warehouse-inventory/

cmd/
 ├── api/
 ├── worker/
 └── migrate/

internal/
 ├── config/
 ├── domain/
 ├── service/
 ├── repository/
 ├── handler/
 ├── middleware/
 ├── event/
 ├── queue/
 └── worker/

pkg/
 ├── logger
 ├── validator
 ├── errors
 └── tracing

api/
 ├── openapi
 └── proto

migrations/
scripts/
docker/
tests/

go.mod
```

---

## 2.2 Service Layer Architecture

```
API Layer (Gin/Fiber)
        ↓
Service Layer (Business Logic)
        ↓
Repository Layer
        ↓
PostgreSQL + Redis
```

---

## 2.3 Repository Pattern

Benefits:

* Dependency injection
* Testable architecture
* Decoupled data layer
* Replaceable storage implementation

---

## 2.4 Domain Models

Core entities:

* Product
* Warehouse
* Inventory
* StockMovement
* Batch
* Order

Important concept:

```
Inventory.Version → used for optimistic locking
```

---

## 2.5 API Structure

Example REST endpoints:

```
/api/v1/products
/api/v1/products/:id
/api/v1/warehouses
/api/v1/inventory
/api/v1/orders/purchase
/api/v1/orders/sales
/api/v1/transfers
/api/v1/batches
/api/v1/reports
/api/v1/audit/logs
```

---

## 2.6 Background Workers

Worker types:

* Stock recalculation worker
* Report generation worker
* Alert worker
* Sync worker

Architecture:

```
API → Queue → Worker → Processor
```

---

## 2.7 Event System

Example events:

* product.created
* product.updated
* inventory.adjusted
* order.created
* order.shipped
* transfer.completed
* batch.created
* expiry.alert

---

## 2.8 Technology Stack

| Component      | Technology              |
| -------------- | ----------------------- |
| HTTP Framework | Gin / Fiber             |
| Database       | PostgreSQL              |
| Cache          | Redis                   |
| Event System   | Kafka / RabbitMQ        |
| ORM            | GORM / sqlx             |
| Validation     | go-playground/validator |
| Config         | Viper                   |
| Logging        | Zerolog / Zap           |
| Tracing        | OpenTelemetry           |

---

# PHASE 3 — DATABASE DESIGN

## Core Tables

* products
* warehouses
* locations
* batches
* inventory
* stock_movements
* purchase_orders
* purchase_order_items
* sales_orders
* sales_order_items
* transfers
* transfer_items
* barcodes
* audit_logs

### Example Constraints

* Unique SKU
* Batch per product uniqueness
* Non-negative quantities
* Reserved ≤ quantity

### Important Indexes

```
inventory(product_id, warehouse_id)
stock_movements(created_at)
batches(expiry_date)
sales_orders(status, order_date)
```

---

# PHASE 4 — CORE MODULES

## Product Management

Handles:

* Product master data
* SKU validation
* Barcode management
* Product activation

---

## Warehouse Management

Handles:

* Warehouse configuration
* Location structure
* Capacity limits

---

## Inventory Tracking

Handles:

* Stock quantities
* Reserved inventory
* Batch integration
* Expiry monitoring
* Movement history

---

## Order Management

### Purchase Orders

* Supplier orders
* Receiving workflow
* Put-away tasks

### Sales Orders

* Order validation
* Allocation
* Picking
* Shipping

---

## Barcode Module

Handles:

* Barcode formats
* Scan processing
* Product lookup
* Location scanning

---

## Reporting Module

Reports include:

* Inventory levels
* Movement history
* Stock aging
* Turnover metrics
* Expiry alerts

---

## Audit Logs

Captures:

* entity changes
* user actions
* before/after data
* compliance history

---

# PHASE 5 — CONCURRENCY & CONSISTENCY

## Optimistic Locking

Inventory record contains:

```
version INT
```

Update only succeeds when version matches.

---

## Row Level Locking

```
SELECT ... FOR UPDATE
```

Used for:

* stock allocation
* transfers
* inventory adjustments

---

## Transaction Strategy

Critical operations wrapped inside:

```
BEGIN TRANSACTION
...
COMMIT
```

Guarantees **ACID consistency**.

---

# PHASE 6 — BACKGROUND JOBS

Async tasks:

* stock recalculation
* expiry alerts
* large report generation
* notification sending
* data synchronization
* cleanup tasks

Worker architecture:

```
API → Kafka → Worker → Processor
```

---

# PHASE 7 — PERFORMANCE OPTIMIZATION

## Redis Caching

Cache examples:

| Data         | TTL   |
| ------------ | ----- |
| product data | long  |
| inventory    | short |
| locations    | long  |

Cache invalidation triggered on inventory updates.

---

## Event Queue

Kafka topics:

* inventory.events
* order.events
* stock.movements
* alerts

---

## Database Optimization

Techniques:

* composite indexes
* read replicas
* connection pooling
* pagination
* batch inserts

---

# PHASE 8 — IMPLEMENTATION PLAN

## Step 1 — Foundation Setup (Week 1-2)

* initialize Go project
* setup Gin/Fiber
* configure PostgreSQL
* configure Redis
* logging setup
* configuration system
* middleware
* Docker dev environment

---

## Step 2 — Database Schema (Week 2-3)

* write migrations
* create domain models
* repository interfaces
* CRUD operations
* seed data
* integration tests

---

## Step 3 — Core Inventory Logic (Week 3-5)

* product management
* warehouse/location system
* inventory service
* stock movements
* batch tracking
* expiry logic
* inventory adjustments

---

## Step 4 — Order System (Week 5-7)

* purchase orders
* receiving
* sales orders
* allocation
* picking
* shipping
* transfer orders

---

## Step 5 — Background Workers (Week 7-8)

* Kafka/RabbitMQ setup
* worker framework
* stock recalculation worker
* expiry alert worker
* report worker
* notification worker

---

## Step 6 — Caching Layer (Week 8-9)

* Redis caching
* cache invalidation
* distributed locking

---

## Step 7 — Event System (Week 9-10)

* event publishing
* Kafka topics
* event consumers
* audit events

---

## Step 8 — Reporting (Week 10-11)

* inventory reports
* movement reports
* expiry reports
* CSV export
* scheduled reports

---

## Step 9 — Testing & Finalization (Week 11-12)

* unit tests
* integration tests
* performance testing
* load testing
* API documentation
* security review
* deployment configuration

---

# END OF DOCUMENT
