


#  Go Transaction Processing Engine

A lightweight backend service built in **Golang** that simulates a card transaction processing system (similar to a payment authorization engine). It supports card validation, secure PIN verification, transaction handling, and history tracking using in-memory storage.

---

##  Features

* Card management (in-memory)
* Secure PIN validation using SHA-256 hashing
* Transaction processing:

  * Withdraw
  * Top-up
* Balance inquiry
* Transaction history tracking
* RESTful APIs
* Thread-safe handling (important for concurrent requests)

---

##  Tech Stack

* **Language:** Go (Golang)
* **Storage:** In-memory (HashMap / map)
* **Hashing:** SHA-256
* **API:** net/http (standard library)

---

##  Project Structure

```
go-transaction/
│── main.go
│── handlers/
│── models/
│── service/
│── storage/
│── utils/
│── README.md
```

---

##  Setup Instructions

### 1. Clone the repository

```bash
git clone https://github.com/Gitesg/go-transcation.git
cd go-transcation
```

### 2. Run the application

```bash
go run main.go
```

Server will start at:

```
http://localhost:8080
```

---

##  API Endpoints

### 1. Process Transaction

**POST** `/api/transaction`

#### Request

```json
{
  "cardNumber": "4123456789012345",
  "pin": "1234",
  "type": "withdraw",
  "amount": 200
}
```

#### Success Response

```json
{
  "status": "SUCCESS",
  "respCode": "00",
  "balance": 800
}
```

#### Error Responses

* Invalid Card

```json
{
  "status": "FAILED",
  "respCode": "05",
  "message": "Invalid card"
}
```

* Invalid PIN

```json
{
  "status": "FAILED",
  "respCode": "06",
  "message": "Invalid PIN"
}
```

* Insufficient Balance

```json
{
  "status": "FAILED",
  "respCode": "99",
  "message": "Insufficient balance"
}
```

---

### 2. Get Balance

**GET** `/api/card/balance/{cardNumber}`

#### Example

```bash
curl http://localhost:8080/api/card/balance/4123456789012345
```

---

### 3. Get Transaction History

**GET** `/api/card/transactions/{cardNumber}`

#### Example

```bash
curl http://localhost:8080/api/card/transactions/4123456789012345
```

---

##  Security

* PIN is stored using **SHA-256 hashing**
* Plaintext PIN is never stored
* PIN is never logged
* Comparison is always done using hashed values

---

##  Sample Test Using Curl

### Withdraw

```bash
curl -X POST http://localhost:8080/api/transaction \
-H "Content-Type: application/json" \
-d '{
  "cardNumber":"4123456789012345",
  "pin":"1234",
  "type":"withdraw",
  "amount":200
}'
```

---

##

### Card

| Field      | Description         |
| ---------- | ------------------- |
| cardNumber | Unique identifier   |
| cardHolder | Name of card holder |
| pinHash    | SHA256 hashed PIN   |
| balance    | Available balance   |
| status     | ACTIVE / BLOCKED    |

---

### Transaction

| Field         | Description         |
| ------------- | ------------------- |
| transactionId | Unique ID           |
| cardNumber    | Associated card     |
| type          | withdraw / topup    |
| amount        | Transaction amount  |
| status        | SUCCESS / FAILED    |
| timestamp     | Time of transaction |

---



* **In-memory storage** for simplicity and speed
* **Stateless APIs**
* **Separation of concerns** (handler → service → storage)
* Designed keeping **concurrency & race conditions** in mind

---


---



