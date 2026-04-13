# --- Test commands for PowerShell ---
# Run these AFTER docker-compose up --build

# 1. Create a successful order (amount 500.00 = within limit)
$body = '{"customer_id": "cust-1", "item_name": "Laptop", "amount": 50000}'
Invoke-RestMethod -Uri "http://localhost:9080/orders" -Method POST -ContentType "application/json" -Body $body

# 2. Create an order that gets DECLINED (amount 2000.00 = over limit)
$body2 = '{"customer_id": "cust-1", "item_name": "Sports Car", "amount": 200000}'
Invoke-RestMethod -Uri "http://localhost:9080/orders" -Method POST -ContentType "application/json" -Body $body2

# 3. Get order by ID (replace ORDER_ID with real ID from step 1)
# Invoke-RestMethod -Uri "http://localhost:9080/orders/ORDER_ID" -Method GET

# 4. Cancel a Pending order
# Invoke-RestMethod -Uri "http://localhost:9080/orders/ORDER_ID/cancel" -Method PATCH

# 5. Get payment by order ID
# Invoke-RestMethod -Uri "http://localhost:9081/payments/ORDER_ID" -Method GET
