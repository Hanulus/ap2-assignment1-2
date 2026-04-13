# AP2 Assignment 2 — gRPC Migration

**Student:** Chingiz Uraimov  
**Group:** SE-2405  

## Repository Links
- **Proto Repository:** https://github.com/Hanulus/ap2-protos
- **Generated Code Repository:** https://github.com/Hanulus/ap2-generated

---

## Architecture

```
[Client]
   |
   | REST (POST /orders, GET /orders/:id)
   v
[Order Service :9080]
   |
   | gRPC (ProcessPayment)
   v
[Payment Service :9082 gRPC]

[gRPC Client]
   |
   | gRPC Server-side Streaming (SubscribeToOrderUpdates)
   v
[Order Service :9083]  <-- polls real DB every second
```

### What Changed from Assignment 1
- `Order Service → Payment Service` call migrated from **REST to gRPC**
- Added **Server-side Streaming** RPC for real-time order status updates
- Clean Architecture preserved: only transport/repository layers changed, use cases untouched

---

## How to Run

### Prerequisites
- Docker & Docker Compose

### Start all services
```bash
docker-compose up --build
```

This starts:
- `payment-service` REST `:9081`, gRPC `:9082`
- `order-service` REST `:9080`, streaming gRPC `:9083`

---

## API Endpoints

### Create Order (triggers gRPC call to Payment Service internally)
```bash
curl -X POST http://localhost:9080/orders \
  -H "Content-Type: application/json" \
  -d '{"customer_id":"cust-1","item_name":"Book","amount":500}'
```

### Get Order
```bash
curl http://localhost:9080/orders/<order_id>
```

### Cancel Order
```bash
curl -X PATCH http://localhost:9080/orders/<order_id>/cancel
```

### Trigger a streaming update (update status in DB)
```sql
UPDATE orders SET status = 'Cancelled' WHERE id = '<order_id>';
```
The streaming subscriber receives the new status within 1 second.

---

## Contract-First Flow

1. `.proto` files live in `ap2-protos` repo
2. On every push to `main`, GitHub Actions runs `protoc` and pushes `.pb.go` files to `ap2-generated`
3. Services import: `go get github.com/Hanulus/ap2-generated@v1.0.0`

### Local generation
```bash
chmod +x generate.sh
./generate.sh
```

---

## Bonus: gRPC Logging Interceptor
Payment Service logs every incoming gRPC call:
```
[gRPC] method=/payment.PaymentService/ProcessPayment duration=1.2ms err=<nil>
```
