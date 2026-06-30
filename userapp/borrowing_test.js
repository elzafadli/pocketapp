const BASE_URL = "http://localhost:8000";
const credentials = "user:password";
const authHeader = "Basic " + Buffer.from(credentials).toString("base64");

async function sendBorrowRequest(customerID, bookIDs) {
  const payload = {
    customer_id: customerID,
    due_at: "2026-06-02T17:00:00Z",
    book_ids: bookIDs,
  };

  try {
    const response = await fetch(`${BASE_URL}/v1/borrowing`, {
      method: "POST",
      headers: {
        "Content-Type": "application/json",
        Authorization: authHeader,
      },
      body: JSON.stringify(payload),
    });

    const status = response.status;
    const body = await response.json();
    console.log(
      `[CustomerID: ${customerID}] Status: ${status}`,
      JSON.stringify(body),
    );
  } catch (error) {
    console.error(`[CustomerID: ${customerID}] Error:`, error.message);
  }
}

async function runTest() {
  console.log("--- Memulai Test Peminjaman Buku (Node.js) ---");

  // Skenario 1: ID customer sama (Menguji aturan "1 orang cuma bisa 1 stock buku")
  console.log(
    "\nSkenario 1: Menguji batas peminjaman aktif untuk orang yang sama...",
  );
  await Promise.all([
    sendBorrowRequest(1, [5]),
    sendBorrowRequest(1, [5]),
    sendBorrowRequest(1, [5]),
  ]);

  // Skenario 2: ID customer berbeda (Menguji batas stock buku)
  console.log(
    "\nSkenario 2: Menguji batas stock buku dengan ID customer berbeda...",
  );
  await Promise.all([
    sendBorrowRequest(1, [5]),
    sendBorrowRequest(2, [5]),
    sendBorrowRequest(3, [5]),
    sendBorrowRequest(4, [5]),
  ]);
}

runTest();
