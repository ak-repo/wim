# Warehouse Inventory Management - Recommended Implementation Order

This document captures the recommended module implementation order for this project, based on the business dependencies described in `docu/BUSINESS_WORKFLOW.md`.

## Recommended Order

1. Authentication + User Management
2. Authorization / Roles / Permissions
3. Product Management
4. Warehouse Management
5. Location Management
6. Inventory Core
7. Stock Movement Ledger
8. Purchase Orders + Receiving
9. Put-away
10. Sales Orders
11. Allocation / Reservation
12. Picking
13. Packing
14. Shipping
15. Transfers
16. Batch + Expiry
17. Barcode / Scanning
18. Returns + Adjustments + Damage
19. Audit Logs
20. Reporting
21. Events, Kafka, Workers
22. Redis Caching / Performance Layer

## Why This Order Fits The Project

- Authentication comes first because business records need identity fields such as `created_by`, `performed_by`, and audit ownership.
- Authorization comes immediately after authentication because warehouse staff, managers, purchasing, sales, and admin users need different access levels.
- Product management must exist before purchase orders, sales orders, batch tracking, barcode workflows, or inventory operations can work correctly.
- Warehouse and location management must exist before receiving, put-away, picking, shipping, or transfers can be modeled properly.
- Inventory core should be implemented before order workflows because allocation and reservation depend on quantity, reserved quantity, and available quantity rules.
- Inbound flows should be completed before outbound flows because purchase receiving creates stock and sales shipping consumes stock.
- Picking, packing, and shipping should come only after sales orders and allocation are stable, because they depend on reserved inventory.
- Transfers should come after warehouse and inventory workflows are reliable, since transfers are controlled stock-out and stock-in operations.
- Batch, expiry, and barcode capabilities are important cross-cutting warehouse enhancements, but they are safer to add after the core stock flows are working.
- Reporting, Kafka, workers, and Redis should come after the transactional source-of-truth flows are stable.

## Practical Phased Roadmap

### Phase 1

- Authentication
- User management
- Roles and permissions
- Session or JWT support
- Password reset and account lifecycle
- User identity propagation for audit fields

### Phase 2

- Product management
- Warehouse management
- Location management

### Phase 3

- Inventory core
- Stock movement ledger
- Manual inventory adjustments

### Phase 4

- Purchase orders
- Receiving workflow
- Put-away workflow

### Phase 5

- Sales orders
- Allocation and reservation
- Shipping workflow

### Phase 6

- Picking
- Packing
- Transfers

### Phase 7

- Batch tracking
- Expiry control
- Barcode workflows
- Returns, damage, and exception flows

### Phase 8

- Audit reporting
- Kafka and worker processing
- Redis caching
- Performance optimization

## Important Note For This Repo

Before implementing authentication deeply, standardize the user identity type across the project. The current `internal/domain/user.go` uses an `int` ID while the rest of the domain and schema references user-related values as UUIDs. That should be aligned first so auth and audit can evolve on a consistent foundation.
