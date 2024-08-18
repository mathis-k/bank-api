# Bank API

A simple JSON API server for a banking system built with Go (Golang), Gorilla/mux, MongoDB, and JWT for authorization and authentication.

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

- **POST /api/auth/register**: Register a new user
- **POST /api/auth/login**: Login an existing user
- **GET /api/auth/logout**: Logout a user
- **GET /api/auth/user**: Get the current user
- **GET /api/auth/refresh**: Refresh the JWT token

### Users

- **GET /api/user**: Get the current user
- **PUT /api/user**: Update the current user
- **GET /api/user/accounts**: Get all accounts for the current user
- **GET /api/user/transactions**: Get all transactions for the current user

### Accounts

- **GET /api/accounts**: Get all accounts for the current user
- **GET /api/accounts/{id}**: Get an account by ID for the current user
- **POST /api/accounts**: Create a new account for the current user
- **PUT /api/accounts/{id}**: Update an account by ID for the current user
- **DELETE /api/accounts/{id}**: Delete an account by ID for the current user

### Transactions

- **GET /api/transactions**: Get all transactions for the current user
- **GET /api/transactions/{id}**: Get a transaction by ID for the current user
- **POST /api/transactions**: Create a new transaction for the current user
- **GET /api/transactions/account/{id}**: Get all transactions for an account for the current user
- **GET /api/transactions/account/{id}/balance**: Get the balance for an account for the current user
- **Post /api/transactions/account/{id}/deposit**: Deposit funds into an account for the current user
- **POST /api/transactions/account/{id}/withdraw**: Withdraw funds from an account for the current user
- **GET /api/transactions/account/{id}/transactions**: Get all transactions for an account for the current user

### For Admins

- **GET /api/admin/users**: Get all users
- **GET /api/admin/users/{id}**: Get a user by ID
- **DELETE /api/admin/users/{id}**: Delete a user by ID
- **GET /api/admin/accounts**: Get all accounts
- **GET /api/admin/accounts/{id}**: Get an account by ID
- **DELETE /api/admin/accounts/{id}**: Delete an account by ID
- **GET /api/admin/transactions**: Get all transactions
- **GET /api/admin/transactions/{id}**: Get a transaction by ID
