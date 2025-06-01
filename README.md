# Loyalty Points Management API

## Project Overview

This project provides a RESTful API built with **Golang 1.20+** for managing user loyalty points. It enables user registration and authentication, tracks transactions that earn points, supports points redemption with mandatory reasons, and offers paginated, filtered views of points history. Security is ensured through JWT-based authentication using access and refresh tokens. Data persistence is handled with **PostgreSQL** or **MySQL** via the **GORM** ORM.

---

## Database Schema

### 1. Users

| Field           | Type         | Description               |
| --------------- | ------------ | ------------------------- |
| id              | uint         | Primary key               |
| email           | varchar(255) | Unique, required          |
| password        | varchar(255) | Hashed password, required |
| points\_balance | int64        | Current loyalty points    |
| created\_at     | timestamp    | Managed automatically     |
| updated\_at     | timestamp    | Managed automatically     |

### 2. Transactions

| Field               | Type         | Description                    |
| ------------------- | ------------ | ------------------------------ |
| id                  | uint         | Primary key                    |
| transaction\_id     | varchar(100) | Unique external transaction ID |
| user\_id            | uint         | Foreign key to Users           |
| transaction\_amount | float64      | Amount involved in transaction |
| category            | varchar(100) | Transaction category           |
| transaction\_date   | timestamp    | Date of transaction            |
| product\_code       | varchar(100) | Associated product code        |
| created\_at         | timestamp    | Managed automatically          |

### 3. Points Records

| Field                    | Type         | Description                     |
| ------------------------ | ------------ | ------------------------------- |
| id                       | uint         | Primary key                     |
| user\_id                 | uint         | Foreign key to Users            |
| points                   | int64        | Points earned or redeemed       |
| type                     | varchar(50)  | Record type (earn/redeem)       |
| reason                   | varchar(255) | Mandatory reason for the change |
| related\_transaction\_id | \*uint       | Optional FK to a transaction    |
| created\_at              | timestamp    | Managed automatically           |

---

## Key Features

* **User Registration & Login:** Secure signup with email verification and hashed passwords.
* **JWT Authentication:** Access tokens passed in Authorization headers, refresh tokens stored as HttpOnly cookies.
* **Transaction Management:** Add transactions that contribute to points accumulation.
* **Points Redemption:** Redeem points with a required reason for auditing.
* **Points History:** View paginated and filtered records by date range and type.
* **Persistent Storage:** MySQL database integration via GORM.

---

## Technologies Used

* **Golang 1.20+**
* **Gorilla Mux** for HTTP routing
* **GORM** as the ORM
* **PostgreSQL / MySQL** for the database
* **JWT** for authentication
* **TOML / Environment configs** for configuration
* Custom logging utilities

---

## Setup & Installation

### 1. Clone the repository

```bash
git https://github.com/jeeva1019/loyalty-management.git
cd loyalty-management
```

### 2. Install dependencies

```bash
go mod tidy
```

### 3. Prepare your database

Create a MySQL or PostgreSQL database named `loyalty_point_system` (or as configured in `dbconfig.toml`).

### 4. Run the API server

```bash
go run cmd/main.go
```

The server will start on the configured port (default **8080**).

---

## API Usage

### 1. `POST /signup`

Register a new user with email and password.

### 2. `POST /login`

Authenticate a user and receive access and refresh JWT tokens.

### 3. `POST /api/transaction`

Add a transaction to earn points (requires authentication).

### 4. `POST /api/points/redeem`

Redeem loyalty points with a required reason (requires authentication).

### 5. `GET /api/points/balance?page=1&page_size=20`

Get the current points balance along with paginated transaction history.

### 6. `GET /api/points/history?start_date=2024-06-01&end_date=2025-06-02&start=1&end=20&type=earn`

Get paginated points history filtered by:

| Query Parameter | Description                                |
| --------------- | ------------------------------------------ |
| start\_date     | Start date filter (YYYY-MM-DD)             |
| end\_date       | End date filter (YYYY-MM-DD)               |
| start           | Page number (default: 1)                   |
| end             | Page size (default: 20)                    |
| type            | Filter by points record type (earn/redeem) |

Got it! Here's how you can update the README to mention that the Postman collection is in the **`document`** folder:

---

## Postman Collection

For your convenience, a Postman collection is provided to help you quickly test all API endpoints with example requests and responses.

**Import the Postman collection file located at:**
`document/LoyalityPoints.postman_collection.json`

This collection includes:

* **SignUp**: Create a new user by sending a `POST` request to `/signup` with email and password.
* **LogIn**: Authenticate by sending a `POST` request to `/login` with email and password; returns a JWT access token in the response header and a refresh token as a cookie.
* **PointBalance**: Retrieve the user’s points balance with a `GET` request to `/api/points/balance`, supporting pagination through query parameters.
* **PointsHistory**: Fetch the user’s points transaction history via a `GET` request to `/api/points/history`, filtered by date range, pagination, and transaction type.
* **Transaction**: Record a new transaction and earn points by sending a `POST` request to `/api/transaction`.
* **RedeemPoint**: Redeem points by recording a transaction through a `POST` request to `/api/points/redeem`.
