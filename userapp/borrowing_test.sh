#!/bin/bash

BASE_URL="http://localhost:8000"
AUTH="user:password"

echo "--- Memulai Test Peminjaman Buku (curl) ---"

# =============================================
# Skenario 1: ID customer sama
# (Menguji aturan "1 orang cuma bisa 1 stock buku")
# =============================================
echo ""
echo "Skenario 1: Menguji batas peminjaman aktif untuk orang yang sama..."

for i in 1 2 3; do
  echo "[CustomerID: 1 | Request #$i]"
  curl -s -o /dev/null -w "Status: %{http_code}\n" \
    -X POST "$BASE_URL/v1/borrowing" \
    -u "$AUTH" \
    -H "Content-Type: application/json" \
    -d '{
      "customer_id": 1,
      "due_at": "2026-06-02T17:00:00Z",
      "book_ids": [5]
    }'
done

# =============================================
# Skenario 2: ID customer berbeda
# (Menguji batas stock buku)
# =============================================
echo ""
echo "Skenario 2: Menguji batas stock buku dengan ID customer berbeda..."

for customer_id in 1 2 3 4; do
  echo "[CustomerID: $customer_id]"
  curl -s -o /dev/null -w "Status: %{http_code}\n" \
    -X POST "$BASE_URL/v1/borrowing" \
    -u "$AUTH" \
    -H "Content-Type: application/json" \
    -d "{
      \"customer_id\": $customer_id,
      \"due_at\": \"2026-06-02T17:00:00Z\",
      \"book_ids\": [5]
    }"
done
