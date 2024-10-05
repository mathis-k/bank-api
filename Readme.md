# Bank API

A simple JSON API server for a banking system built with Go (Golang), Gorilla/mux (no Gin), MongoDB, and JWT for authorization and authentication.

## Features

- **User Authentication**: Secure JWT-based authorization and authentication.
- **Account Management**: Create, view, update, and delete bank accounts.
- **Transaction Management**: Perform and track transactions between accounts.
- **Database**: MongoDB for data storage.

## Tech Stack

- **Golang**: Backend server
- **Gorilla/mux**: HTTP router and dispatcher
- **MongoDB**: NoSQL database
- **JWT**: Secure token-based authentication

## Installation

1. Clone the repository:

   ```bash
   git clone https://github.com/mathis-k/bank-api.git
   cd bank-api
2. Install dependencies:

   ```bash
   go mod download
3. Set up your environment variables:

   ```bash
    cp .env.example .env
4. Install dependencies:

   ```bash
   go run main.go

## API Endpoints

### Authentication

- **POST /api/auth/register**: Register a new user \
Request Body:
  ```json
    {
      "first_name": "John",
      "last_name": "Doe",
      "email": "john.doe@example.com",
      "password": "password123"
    }
- **POST /api/auth/login**: Login an existing user \
  Request Body:
  ```json
    {
      "email": "john.doe@example.com",
      "password": "password123"
    }
### Users

- **GET /api/user**: Get the current user
- **PUT /api/user**: Update the current user \
  Request Body:
  ```json
    {
      "first_name": "John",
      "last_name": "Doe",
      "email": "john.doe@example.com"
    }
  or
    {
      "last_name": "Doey"
    }

### Accounts

- **GET /api/accounts**: Get all accounts for the current user
- **GET /api/accounts/{number}**: Get an account by ID for the current user
- **POST /api/accounts**: Create a new account for the current user
- **DELETE /api/accounts/{number}**: Delete an account by ID for the current user

### Transactions

- **GET /api/transactions**: Get all transactions for the current user
- **GET /api/transactions/{id}**: Get a transaction by ID for the current user
- **GET /api/transactions/account/{number}**: Get all transactions for an account from the current user
- **POST /api/transactions/account/{number}/deposit**: Deposit funds into an account from the current user \
  Request Body:
  ```json
    {
      "amount": 150.00
    }
- **POST /api/transactions/account/{number}/withdraw**: Withdraw funds from an account from the current user \
  Request Body:
  ```json
    {
      "amount": 150.00
    }
- **POST /api/transactions/account/{number}/transfer**: Transfer funds from an account of the current user to another account \
  Request Body:
  ```json
    {
      "amount": 150.00,
      "to_account": "7252934484834"
    }


## Project Structure

```bash 
.
├── main.go
├── routes/
├── controllers/
├── models/
├── middleware/
├── utils/
├── .env.example
├── go.mod
└── go.sum