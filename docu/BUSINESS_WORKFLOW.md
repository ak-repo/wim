# Warehouse Inventory Management System - Business Workflow Guide

This document explains the project in business language first and technical language second. It is written for beginners, warehouse staff, analysts, founders, students, and developers who want to understand what the system does in the real world.

# 1. SYSTEM OVERVIEW

## What problem does this system solve?

A warehouse business has one big daily challenge:

- It must know what products it has
- where those products are stored
- how much stock is free to sell
- how much stock is already promised
- what is coming in from suppliers
- what is going out to customers
- what moved, when, and why

If this is managed with spreadsheets, phone calls, memory, or paper, the business usually faces:

- missing items
- wrong shipments
- duplicate selling of the same stock
- expiry losses
- slow receiving and shipping
- no clear history when something goes wrong

This Warehouse Inventory Management System solves that by giving the business one central operating system for inventory movement.

In simple words:

- products are defined once
- warehouses and locations are mapped
- every stock change is recorded
- inbound stock is received and stored
- outbound stock is reserved and shipped
- transfers between warehouses are controlled
- barcode-based actions reduce mistakes
- audit logs keep a business history

## Who uses it?

Different people use the same system for different reasons.

### Warehouse staff

They use the system to:

- receive goods from suppliers
- place goods into storage locations
- pick goods for customer orders
- pack and ship orders
- count stock and report mismatches

### Warehouse manager

They use the system to:

- monitor stock levels
- check movement history
- review warehouse activity
- investigate errors
- track expiring items

### Purchasing team

They use the system to:

- create purchase orders
- track expected deliveries
- compare what was ordered versus what arrived

### Sales or customer operations team

They use the system to:

- create customer sales orders
- check whether stock is available
- track order status from creation to shipment

### Business admin or compliance team

They use the system to:

- review audit logs
- check reports
- monitor inventory changes
- support investigations and compliance work

### Supplier and customer

They may not directly log into the system, but the system tracks the business work done for them:

- suppliers send goods in through purchase orders
- customers receive goods through sales orders and shipments

## Real-world analogy

Think of the system like an airport control tower for products.

- Products are the planes
- warehouses are the airports
- locations are the gates and parking spots
- purchase orders are arriving flights
- sales orders are departing flights
- transfers are domestic repositioning flights
- barcodes are the plane identity checks
- stock movements are the flight logs
- audit logs are the security and control records

Without the control tower, planes still move, but nobody truly knows where they are, what has landed, what has departed, and what went wrong.

# 2. COMPLETE END-TO-END STORY

This is the full warehouse story from start to finish.

Imagine a company sells packaged foods and household products. It keeps inventory in one or more warehouses and wants to avoid stock confusion.

## Part 1: Procurement - deciding what to buy

The business notices that stock for cooking oil is running low. The purchasing team checks demand and sees that more units are needed.

So they create a purchase order, also called a PO.

A purchase order is an official buying document that says:

- what product is being purchased
- how many units are needed
- from which supplier
- at what agreed price
- for which warehouse
- and when the goods are expected

In the codebase, purchase orders are created through the order service and stored in PostgreSQL as business records.

At this point, nothing has arrived yet. The system is only preparing for inbound stock.

## Part 2: Supplier sends the goods

The supplier receives the PO and prepares the shipment.

Later, a truck arrives at the receiving dock of the warehouse.

Warehouse staff starts the receiving process. They:

- look up the purchase order
- unload the goods
- count what arrived
- inspect for damage or shortage
- confirm the receiving location

If the PO expected 500 units but only 450 arrived, the system can mark the purchase order as partially received instead of fully received.

This is important because the business must separate:

- what was ordered
- what actually arrived

In the current codebase, receiving a purchase order:

- updates received quantities on PO items
- creates inventory records
- creates a `RECEIPT` stock movement
- updates the PO status to `PARTIAL` or `RECEIVED`

So the system does not just say "goods came in." It records exactly what happened.

## Part 3: Receiving area to storage - put-away

After goods are received, they may still be at the dock or in a temporary receiving area. They are in the building, but they are not yet properly stored.

The next step is put-away.

Put-away means moving received items from the receiving area to their real storage location.

For example:

- Warehouse: Main Warehouse
- Zone: Storage
- Aisle: A-01
- Rack: R-03
- Bin: B-02

The worker scans the item and then scans the location barcode. This confirms that the product is now stored there.

Business value:

- staff can find products later
- location discipline reduces lost stock
- picking becomes faster and more accurate

The technical document mentions put-away as a workflow. In the current codebase, the location-based inventory structure supports this operational model even though put-away is not exposed as a dedicated API step yet.

## Part 4: Storage - goods wait until needed

Now the products sit in storage.

The business needs the system to answer simple questions at any moment:

- how many units do we have?
- in which warehouse?
- in which location?
- in which batch?
- how many are still free?
- how many are already promised?

That is the role of the inventory system.

The inventory record in this project tracks:

- product
- warehouse
- location
- optional batch
- total quantity
- reserved quantity
- version number for safety

So inventory is not just one number. It is stock with context.

## Part 5: Customer order arrives

Now a customer places an order for the same cooking oil.

The sales team enters a sales order into the system.

A sales order tells the business:

- who the customer is
- what items they want
- how many they want
- from which warehouse the goods will ship
- where the goods should be delivered

When the order is first created, it is only a demand signal. It does not yet mean stock has left the warehouse.

In the codebase, a newly created sales order starts as:

- `PENDING`
- allocation status `UNALLOCATED`

That means the business knows the customer wants the goods, but the system has not yet reserved stock.

## Part 6: Allocation and reservation - promising stock safely

Before warehouse staff starts picking, the system must decide whether the requested stock is truly available.

This is where allocation happens.

Allocation means the system chooses inventory to fulfill the order.

Reservation means that chosen stock is now promised to that customer order and should not be given to anyone else.

Example:

- Total stock = 200 units
- Reserved for older orders = 50 units
- Available for new orders = 150 units

If a new customer requests 120 units, the system can reserve them.

If a new customer requests 180 units, the system should reject or delay the allocation because only 150 are actually free.

In the current codebase, allocation:

- runs inside a database transaction
- locks the sales order row
- locks the sales order items
- locks matching inventory rows
- checks total available stock
- increases `reserved_quantity`
- marks sales order items as allocated
- updates the sales order to `PROCESSING` and `ALLOCATED`

This prevents overselling and reduces race conditions.

## Part 7: Picking - physically collecting items

Once stock is allocated, warehouse workers can pick it.

Picking means walking to the correct storage location and taking the required quantity off the shelf.

The picker follows instructions such as:

- go to warehouse WH-001
- zone storage
- aisle A-01
- rack R-03
- bin B-02
- pick 120 units

In a barcode-driven warehouse, the picker scans:

- the order
- the location
- the product
- the quantity confirmation

This helps avoid common mistakes like:

- going to the wrong shelf
- picking the wrong product
- picking too many or too few

In this project, the full business concept of picking exists in the technical design and barcode workflow. The current shipping logic consumes reserved stock and records shipment movements, while a dedicated pick task module can be added later as the warehouse process becomes more advanced.

## Part 8: Packing - preparing goods for dispatch

After picking, the goods are brought to a packing station.

Packing means:

- confirming the right items are together
- putting them into cartons or pallets
- labeling the shipment
- preparing documents or tracking labels

Business value:

- fewer wrong shipments
- faster carrier handoff
- clearer package identity

The barcode workflow in the design document includes carton and shipping label scanning to support this step.

## Part 9: Shipping - stock officially leaves the business

Shipping is the moment the warehouse confirms the order has been dispatched.

This is one of the most important inventory moments.

Why?

Because reserved stock should now become shipped stock.

That means:

- total quantity decreases
- reserved quantity decreases
- shipment history is recorded
- order status changes to shipped

In the current codebase, shipping a sales order:

- runs in a database transaction
- locks the sales order and related rows
- reduces quantity from reserved inventory
- reduces reserved quantity
- creates `SHIP` stock movement records
- updates sales order item shipped quantities
- updates the sales order status to `SHIPPED`
- publishes an `order.shipped` event

So shipping is not just a status change. It is an inventory event.

## Part 10: Returns - stock comes back

Sometimes customers return goods.

Common reasons:

- wrong item shipped
- damaged goods
- product no longer wanted
- quality issue

When returned stock arrives, warehouse staff inspects it.

Then one of several things may happen:

- good stock goes back into saleable inventory
- damaged stock goes into a damaged area
- expired stock is blocked
- suspect goods are quarantined for review

The technical design includes return and damage movement types so every return outcome can be recorded.

## Part 11: Internal transfer - moving stock between warehouses

Now imagine the company has two warehouses.

- Warehouse A has too much stock
- Warehouse B is running low

The business creates an internal transfer.

This is not a supplier purchase and not a customer sale. It is an internal company movement.

The transfer process is:

1. create transfer request
2. approve transfer
3. pick stock from source warehouse
4. ship it out internally
5. mark it in transit
6. receive it at destination warehouse
7. store it in a destination location
8. mark transfer completed

In the current codebase, transfer logic is already implemented.

When a transfer is shipped:

- source inventory is reduced
- `TRANSFER_OUT` stock movements are created
- transfer status becomes `IN_TRANSIT`

When a transfer is received:

- destination inventory is created or increased
- `TRANSFER_IN` stock movements are created
- transfer item quantities are marked received
- transfer status becomes `COMPLETED`

This gives the business full visibility from warehouse to warehouse.

## Part 12: Traceability, expiry, and reporting

For some products, the business must know more than just quantity.

It must also know:

- which batch a unit belongs to
- when it was made
- when it expires
- where it came from

That is why batch tracking exists.

This matters for:

- food
- medicine
- cosmetics
- regulated goods
- recall situations

The system can also generate reports for:

- inventory levels
- movement history
- stock aging
- expiring inventory

These reports help managers make business decisions, not just technical decisions.

## Part 13: Audit and business trust

Every serious warehouse system needs trust.

If numbers change, the business should know:

- what changed
- who changed it
- when it changed
- what the old value was
- what the new value became

That is why audit logs exist.

This codebase writes audit entries for major actions like:

- product creation and updates
- inventory adjustments
- purchase order updates
- sales order updates
- transfer changes

This helps with investigations, compliance, and accountability.

In one sentence, the entire story is:

- the business buys goods, receives them, stores them, promises them safely, ships them correctly, transfers them when needed, tracks every movement, and keeps a permanent history.

# 3. CORE CONCEPTS (SIMPLIFIED)

## Warehouse

A warehouse is the physical building where goods are stored.

Example:

- Main warehouse in New York
- overflow warehouse in New Jersey

Business meaning:

- it is the top-level place that owns stock

## Location

A location is the exact address of stock inside the warehouse.

The hierarchy is usually:

- Zone -> big area
- Aisle -> row
- Rack -> shelf structure
- Bin -> smallest slot

Example:

- `ST-A01-R03-B02`

This means the product is not just "in the warehouse." It is in a known spot.

Analogy:

- Warehouse = city
- Zone = district
- Aisle = street
- Rack = building
- Bin = apartment number

## Inventory

Inventory means the stock the business owns and tracks.

But inventory is more than quantity. In this project, inventory includes:

- product
- warehouse
- location
- batch
- total quantity
- reserved quantity

So inventory answers both:

- how much do we have?
- where exactly is it?

## Stock movement

A stock movement is a record of any inventory change.

Common movement types in the design:

- `RECEIPT`
- `PUTAWAY`
- `PICK`
- `PACK`
- `SHIP`
- `TRANSFER_IN`
- `TRANSFER_OUT`
- `ADJUSTMENT`
- `RETURN`
- `DAMAGE`
- `EXPIRY`

Analogy:

- a bank account uses transactions to explain money changes
- a warehouse uses stock movements to explain stock changes

## Batch

A batch is a group of units produced or received together.

Example:

- 1,000 yogurt cups made on the same day from the same production run

Why it matters:

- recalls
- quality control
- expiry management
- supplier traceability

## Expiry

Expiry is the date after which a product should no longer be sold or used.

The system supports expiry-related thinking such as:

- warnings 30 days before expiry
- critical alerts 7 days before expiry
- blocking already expired items from being picked

Common stock rotation rules:

- FIFO = first in, first out
- FEFO = first expired, first out

For food and medicine, FEFO is usually more useful than simple FIFO.

## Allocation vs reservation

These terms are related but slightly different.

### Reservation

Reservation means stock is held for a specific order.

It is still physically on the shelf, but business-wise it is no longer free.

### Allocation

Allocation means the system decides which stock should satisfy the order.

In many warehouse systems, allocation also implies reservation.

Simple example:

- Total stock = 100
- Reserved stock = 30
- Available stock = 70

Those 30 units are not gone yet, but they are already promised.

# 4. HOW EACH MODULE WORKS (DETAILED)

## Product Management

### What it does

It stores master data about products.

In this project, product data includes fields like:

- SKU
- name
- description
- category
- unit of measure
- barcode
- dimensions and weight
- active or inactive status

### Why it exists

Because every warehouse action begins with a product identity.

If products are not clearly defined, the business cannot:

- receive correctly
- pick correctly
- report correctly
- scan correctly

### When it is used

- when a new product is introduced
- when a barcode is assigned
- when a product description changes
- when a product should be deactivated

## Warehouse Management

### What it does

It stores warehouse data and location structure.

In this project, warehouse records include code, name, address, and active status. Locations capture zone, aisle, rack, bin, location code, and location type.

### Why it exists

Because the business must know where stock can physically exist.

### When it is used

- when setting up a warehouse
- when adding new locations
- when reorganizing storage
- when controlling picking and put-away discipline

## Inventory System

### What it does

It tracks stock quantities and availability.

It also connects inventory to:

- locations
- batches
- movements
- audit history

### Why it exists

Because inventory accuracy is the heart of warehouse control.

### When it is used

- during receiving
- during adjustments
- during allocation
- during shipping
- during transfers
- during counting and investigations

Business note:

- inventory tells the business what is physically there
- available inventory tells the business what can still be promised

## Order Management

### Purchase Orders

#### What it does

Handles buying stock from suppliers.

#### Why it exists

Because inbound stock should not arrive without business expectation and record.

#### When it is used

- before supplier delivery
- during receiving
- during inbound discrepancy handling

### Sales Orders

#### What it does

Handles customer demand from order creation to allocation and shipping.

#### Why it exists

Because outbound stock must be controlled to avoid overselling and wrong dispatch.

#### When it is used

- when customers place orders
- when stock is allocated
- when orders are shipped

## Transfer System

### What it does

Moves stock from one warehouse to another.

### Why it exists

Because one warehouse may have excess stock while another has shortage.

### When it is used

- regional rebalancing
- emergency replenishment
- internal distribution planning

## Barcode System

### What it does

Supports fast scanning of:

- products
- locations
- batches
- cartons
- pallets

### Why it exists

Because typing codes manually is slow and error-prone.

### When it is used

- receiving
- put-away
- picking
- packing
- cycle counting
- shipping

## Batch System

### What it does

Tracks product lots with manufacturing and expiry information.

### Why it exists

Because some businesses need traceability and expiry control.

### When it is used

- receiving lot-controlled products
- checking expiring stock
- running recalls

## Reporting Module

### What it does

Provides operational views such as:

- inventory report
- movement report
- expiry report

### Why it exists

Because managers need decisions, not just raw transactions.

### When it is used

- daily operations review
- stock investigation
- expiry planning
- business reporting

## Audit Logs

### What it does

Captures business history of major record changes.

### Why it exists

Because the business needs traceability and accountability.

### When it is used

- after major create or update actions
- during compliance or incident review
- when checking who changed what

# 5. STEP-BY-STEP WORKFLOWS

## Purchase Flow

1. Business sees stock is low or demand is increasing.
2. Purchasing team chooses a supplier.
3. A purchase order is created.
4. PO contains product, quantity, target warehouse, price, and notes.
5. Supplier receives the PO.
6. Supplier prepares shipment.
7. Warehouse waits for arrival.

## Receiving Flow

1. Truck reaches receiving dock.
2. Staff match the goods to a purchase order.
3. Staff count the delivered units.
4. Staff inspect quality and damage.
5. Staff record full receipt or partial receipt.
6. System increases stock in the receiving or chosen location.
7. System creates a `RECEIPT` movement.
8. PO status changes to `PARTIAL` or `RECEIVED`.

## Put-away Flow

1. Received stock waits in the inbound area.
2. Worker gets put-away instruction.
3. Worker takes the stock to a proper storage location.
4. Worker scans item barcode.
5. Worker scans location barcode.
6. System confirms placement.
7. Inventory now points to the final storage location.

## Picking Flow

1. Customer sales order is created.
2. System checks whether stock is available.
3. System allocates and reserves stock.
4. Pick list or pick task is generated.
5. Picker goes to the instructed location.
6. Picker scans location.
7. Picker scans product.
8. Picker confirms quantity.
9. Picked items move to packing area.

## Packing Flow

1. Picked items arrive at packing station.
2. Staff verify order contents.
3. Staff put items into cartons or onto pallets.
4. Staff scan carton or shipping container.
5. Label and tracking information are created.
6. Order becomes ready for dispatch.

## Shipping Flow

1. Packed goods are assigned to carrier.
2. Shipment is confirmed.
3. Reserved stock is deducted from actual stock.
4. `SHIP` movement is recorded.
5. Sales order status changes to `SHIPPED`.
6. Customer delivery can now proceed.

## Transfer Flow

1. Business identifies stock imbalance between warehouses.
2. Transfer request is created.
3. Source warehouse approves transfer.
4. Stock is selected for transfer.
5. Source warehouse ships the transfer.
6. System records `TRANSFER_OUT`.
7. Transfer status becomes `IN_TRANSIT`.
8. Destination warehouse receives the transfer.
9. System records `TRANSFER_IN`.
10. Destination inventory is increased.
11. Transfer status becomes `COMPLETED`.

# 5A. FLOW DIAGRAMS

These diagrams are simple visual shortcuts. They are written in plain text so they are easy to read in any editor.

## Full End-to-End Warehouse Story

```text
Procurement
    |
    v
Purchase Order Created
    |
    v
Supplier Sends Goods
    |
    v
Receiving at Warehouse
    |
    v
Put-away Into Storage Locations
    |
    v
Inventory Available for Orders
    |
    v
Sales Order Created
    |
    v
Allocation / Reservation
    |
    v
Picking
    |
    v
Packing
    |
    v
Shipping
    |
    v
Customer Receives Goods
    |
    v
Return or Close Order
```

## Purchase and Receiving Flow

```text
Stock Need Identified
    |
    v
Create Purchase Order
    |
    v
Send PO to Supplier
    |
    v
Supplier Confirms
    |
    v
Goods Arrive at Dock
    |
    v
Count + Inspect Goods
    |
    +------------------+
    |                  |
    v                  v
Full Receipt       Partial Receipt
    |                  |
    +--------+---------+
             |
             v
Create RECEIPT Movement
             |
             v
Stock Ready for Put-away
```

## Put-away Flow

```text
Received Stock in Inbound Area
    |
    v
Assign Storage Location
    |
    v
Move Product to Shelf/Bin
    |
    v
Scan Item
    |
    v
Scan Location
    |
    v
Confirm Placement
    |
    v
Inventory Updated by Location
```

## Sales Order to Shipment Flow

```text
Customer Places Order
    |
    v
Create Sales Order
    |
    v
Check Available Inventory
    |
    +----------------------+
    |                      |
    v                      v
Enough Stock           Not Enough Stock
    |                      |
    v                      v
Reserve / Allocate      Backorder / Reject / Review
    |
    v
Generate Pick Work
    |
    v
Pick Items
    |
    v
Pack Items
    |
    v
Ship Order
    |
    v
Reduce Inventory + Record SHIP Movement
```

## Inventory Logic Flow

```text
                 +-------------------+
                 |   Total Quantity   |
                 +-------------------+
                           |
                           v
                 +-------------------+
                 | Reserved Quantity  |
                 +-------------------+
                           |
                           v
                 +-------------------+
                 | Available Quantity |
                 | = Total - Reserved |
                 +-------------------+
```

## Transfer Flow Between Warehouses

```text
Warehouse B Needs Stock
    |
    v
Create Transfer Request
    |
    v
Source Warehouse Approves
    |
    v
Pick Stock from Source
    |
    v
Ship Transfer Out
    |
    v
Record TRANSFER_OUT
    |
    v
Status = IN_TRANSIT
    |
    v
Destination Receives Stock
    |
    v
Record TRANSFER_IN
    |
    v
Put-away at Destination
    |
    v
Status = COMPLETED
```

## Returns Flow

```text
Customer Sends Goods Back
    |
    v
Warehouse Receives Return
    |
    v
Inspect Condition
    |
    +-------------------+-------------------+
    |                   |                   |
    v                   v                   v
Good for Resale     Damaged Goods      Expired / Blocked
    |                   |                   |
    v                   v                   v
Return to Stock      Move to Damage      Move to Blocked Area
    |                   |                   |
    +---------+---------+-------------------+
              |
              v
       Record Return Outcome
```

## Event and Worker Flow

```text
User Action in API
    |
    v
Business Record Saved in PostgreSQL
    |
    v
Event Published to Kafka
    |
    v
Worker Reads Message
    |
    v
Worker Pool Processes Job
    |
    +------------------+------------------+
    |                  |                  |
    v                  v                  v
Stock Recalc      Expiry Check      Report Generation
```

# 6. INVENTORY LOGIC (CRITICAL)

## How stock increases

Stock increases when:

- goods are received from a supplier
- stock is transferred in from another warehouse
- a valid customer return goes back into saleable inventory
- a manual positive adjustment is made after investigation

## How stock decreases

Stock decreases when:

- a customer order is shipped
- stock is transferred out to another warehouse
- damaged stock is written off
- expired stock is removed
- a manual negative adjustment is made after investigation

## Reserved vs available stock

This is one of the most important concepts in the whole project.

- Total quantity = all stock recorded physically
- Reserved quantity = stock already promised to orders
- Available quantity = total quantity minus reserved quantity

Example:

- quantity = 500
- reserved = 120
- available = 380

So even though 500 units are in the building, only 380 can be offered to a new customer.

In the codebase, `AvailableQty()` is explicitly modeled in the inventory domain as:

- `Quantity - ReservedQty`

## Why stock inconsistency happens

Real warehouses face stock inconsistency for many reasons:

- wrong count during receiving
- product placed into wrong bin
- item picked without scan
- item damaged but not reported
- delayed transaction entry
- duplicate or concurrent updates
- incomplete transfer process
- poor handling of returns

This project addresses inconsistency through:

- location-based inventory
- movement history
- reservation logic
- audit logging
- transactions for critical workflows
- optimistic locking on inventory updates

## Role of stock movements

Stock movement records answer the question:

- why did the stock number change?

Each movement explains:

- movement type
- product
- warehouse
- from-location or to-location
- quantity
- business reference such as order or transfer
- notes
- time of action

That makes stock movements the warehouse equivalent of an accounting ledger.

Without movement records, you only have numbers.

With movement records, you have explanation.

# 7. REAL-WORLD PROBLEMS AND HOW THE SYSTEM SOLVES THEM

## Problem: Lost items

What happens in real life:

- staff know stock entered the warehouse but cannot find it later

How the system helps:

- inventory is tied to warehouse locations
- barcode scanning confirms placement
- movements show where stock went last

## Problem: Wrong shipments

What happens in real life:

- the customer receives the wrong SKU or wrong quantity

How the system helps:

- sales order structure tells staff exactly what is needed
- allocation reserves correct stock
- scanning verifies location and item
- packing stage adds another check

## Problem: Expired goods shipped to customers

What happens in real life:

- expired stock remains on shelf and gets picked accidentally

How the system helps:

- batch and expiry tracking exist
- expiry alerts can be generated
- FEFO-style thinking can prioritize soonest-expiring stock first
- expired stock can be blocked from outbound use

## Problem: Stock mismatch between system and reality

What happens in real life:

- system says 1,000 units, physical count says 940

How the system helps:

- movement records support investigation
- adjustments can be posted with notes
- audit logs preserve who changed records
- concurrency controls reduce conflicting updates

## Problem: Overselling the same stock

What happens in real life:

- two orders try to use the same inventory

How the system helps:

- reserved quantity separates free stock from promised stock
- allocation is wrapped in database transactions
- inventory rows are locked during critical operations

## Problem: No proof of who changed data

What happens in real life:

- inventory changed, but nobody knows how or why

How the system helps:

- audit logs record business actions
- before/after snapshots can be stored
- entity change history becomes reviewable

# 8. TECH SIMPLIFIED (NO JARGON FIRST)

This section explains why the codebase uses PostgreSQL, Redis, Kafka, and worker pools.

## Why PostgreSQL?

Simple explanation:

- PostgreSQL is the system's official long-term memory.

The business needs one trusted place for:

- products
- warehouses
- locations
- inventory
- purchase orders
- sales orders
- transfers
- stock movements
- batches
- audit logs

Why PostgreSQL fits this project:

- it is reliable for business-critical records
- it supports transactions well
- it handles structured relationships clearly
- it works well for reporting and filtering
- it supports row-level locking for safe allocation and transfer logic

How it appears in the codebase:

- repositories store and read business data from PostgreSQL
- order allocation and shipping use SQL transactions and `FOR UPDATE`
- transfer shipping and receiving also use transactions

Business meaning:

- PostgreSQL protects the truth of the warehouse.

## Why Redis?

Simple explanation:

- Redis is the system's fast short-term memory.

Some data is requested very frequently, such as:

- inventory by ID
- available stock totals
- repeated lookups that should return quickly

Why Redis fits this project:

- it is extremely fast
- it reduces repeated load on PostgreSQL
- it improves response speed for high-read operations

How it appears in the codebase:

- the inventory repository caches inventory records by ID
- the inventory repository caches total available quantity by product and warehouse

Business meaning:

- Redis helps the system answer common questions faster.

Important note:

- Redis is not the main source of truth here
- PostgreSQL remains the source of truth
- Redis is a speed helper

## Why Kafka?

Simple explanation:

- Kafka is the system's message highway.

When important business events happen, the system should be able to notify other parts without slowing down the main user action.

Example:

- an order ships
- inventory changes
- a transfer completes
- expiry alerts need to be checked

Why Kafka fits this project:

- it allows events to be published once and consumed later
- it decouples immediate user actions from background processing
- it supports growth as the system expands

How it appears in the codebase:

- API-side services publish events using a Kafka publisher
- event types include product, inventory, order, transfer, and expiry events
- the worker process consumes Kafka messages and turns them into jobs

Business meaning:

- Kafka helps the warehouse system react to business events without forcing everything to happen in the same request.

## Why background workers?

Simple explanation:

- workers are back-office operators inside the system.

Not every task should happen while the user is waiting on the API response.

Some work is better done in the background, such as:

- stock recalculation
- expiry checks
- report generation

Why workers fit this project:

- they keep the API responsive
- they isolate slower tasks
- they let the platform handle operational work continuously

How it appears in the codebase:

- there is a separate `cmd/worker/main.go`
- Kafka messages are read by a processor
- jobs are submitted into a worker pool
- workers handle topics like stock recalculation, expiry alert checking, and report generation

Business meaning:

- workers let the warehouse keep operating quickly while support tasks happen behind the scenes.

## Why a worker pool?

Simple explanation:

- a worker pool is a team of background workers instead of just one worker.

Why that matters:

- if many messages arrive, one worker may be too slow
- a pool lets several jobs run at the same time
- retries can be managed when a job temporarily fails

How it appears in the codebase:

- the pool has a configurable size
- jobs are placed into a queue
- multiple workers pull jobs from that queue
- failed jobs are retried with delay
- repeated failures are logged as dead-letter style outcomes

Business meaning:

- worker pools help the system scale operational tasks without blocking day-to-day warehouse users.

## Putting it together

- PostgreSQL = official record book
- Redis = quick memory for fast lookup
- Kafka = message highway for events
- Worker pool = back-office team processing delayed work

Together they create a system that is both:

- reliable for business truth
- responsive for day-to-day operations

# 9. ADVANCED CONCEPTS (SIMPLIFIED)

## Optimistic locking

Simple explanation:

- optimistic locking stops two users from silently overwriting each other's inventory changes.

Analogy:

- two clerks read the same stock card showing 100 units
- clerk A updates it first
- clerk B tries to update using the old card view
- system says, "this record already changed, reload first"

How it appears in the codebase:

- inventory has a `version` field
- `UpdateWithVersion` only succeeds if the expected version still matches
- otherwise the repository returns a concurrent update error

Business value:

- reduces silent stock corruption

## Transactions

Simple explanation:

- a transaction means several related actions succeed together or fail together.

Analogy:

- if you transfer money from one bank account to another, you cannot allow only half the action to happen

In this warehouse system:

- allocating stock should not reserve some lines and forget others
- shipping should not reduce inventory without updating order status
- transfer shipment should not subtract stock without recording movement

How it appears in the codebase:

- allocation, shipping, transfer shipping, and transfer receiving run inside PostgreSQL transactions

Business value:

- prevents broken workflows and half-finished stock updates

## Row-level locking

Simple explanation:

- row-level locking temporarily protects the exact database rows being changed.

Analogy:

- if one cashier is counting a cash drawer, another cashier should not edit the same drawer at the same time

How it appears in the codebase:

- critical SQL queries use `FOR UPDATE`
- this is applied during allocation, shipping, and transfer processing

Business value:

- prevents two operations from fighting over the same stock at the same moment

## Event-driven architecture

Simple explanation:

- event-driven architecture means the system reacts to meaningful business events.

Analogy:

- in a restaurant, when an order is placed, the kitchen, billing, and display board all react to the same event in different ways

In this project, example events include:

- `product.created`
- `inventory.adjusted`
- `order.created`
- `order.allocated`
- `order.shipped`
- `transfer.created`
- `transfer.completed`
- `inventory.expiry_alert`

Business value:

- one action can trigger many follow-up processes cleanly
- the system becomes easier to grow over time

# 10. MENTAL MODEL SUMMARY

The simplest mental model is:

- Warehouse system = Brain + Memory + Map + Movers + History

## Brain

The service layer decides business rules.

Examples:

- can this sales order be allocated?
- can this transfer be shipped?
- should this PO be partial or received?

## Memory

- PostgreSQL stores the official business truth
- Redis stores fast temporary lookup data

## Map

The warehouse and location structure tells the system where inventory lives.

## Movers

- API handles immediate business actions
- Kafka carries event messages
- worker pool handles background operational jobs

## History

- stock movements explain quantity changes
- audit logs explain who changed business records

## One-line memory trick

Remember the project like this:

- Product management tells the system what an item is
- warehouse management tells it where an item can live
- inventory tells it how much exists and what is still free
- purchase orders bring stock in
- sales orders send stock out
- transfers move stock between warehouses
- barcodes reduce human error
- stock movements explain quantity change
- audit logs explain record change

## Final business summary

This project is a warehouse control system that helps a business buy stock, receive it, store it, reserve it safely, ship it accurately, transfer it between sites, monitor expiry, investigate issues, and keep a permanent business history.

# APPENDIX: HOW THE CURRENT CODEBASE IS ORGANIZED

This is a business-friendly explanation of the project structure.

## API application

Path:

- `cmd/api/main.go`

Purpose:

- starts the REST API
- connects to PostgreSQL and Redis
- runs migrations
- creates repositories and services
- creates Kafka publisher
- exposes business endpoints

## Worker application

Path:

- `cmd/worker/main.go`

Purpose:

- starts the background worker process
- reads Kafka messages
- pushes messages into a worker pool
- handles background jobs such as stock recalculation and expiry checks

## Domain layer

Path:

- `internal/domain/domain.go`

Purpose:

- defines core business objects like product, warehouse, inventory, order, batch, transfer, barcode, and audit log

## Repository layer

Path:

- `internal/repository/postgres/*`

Purpose:

- stores and retrieves data from PostgreSQL
- uses Redis for selected inventory caching

## Service layer

Path:

- `internal/service/*`

Purpose:

- contains business rules
- examples: allocation, shipping, receiving, adjustments, transfer state changes, report generation

## Event layer

Path:

- `internal/event/*`

Purpose:

- defines event types
- publishes events to Kafka

## Queue and worker layer

Paths:

- `internal/queue/processor.go`
- `internal/worker/pool.go`

Purpose:

- reads Kafka messages
- converts them into jobs
- processes jobs through a configurable worker pool with retries

## HTTP modules available in the API

The codebase currently exposes routes for:

- products
- warehouses
- warehouse locations
- inventory
- batches
- purchase orders
- sales orders
- transfers
- reports
- audit logs
- health checks

That means the project is already structured as a full warehouse operations platform, even though some advanced warehouse steps like dedicated pick-task or put-away-task orchestration can still be expanded later.
