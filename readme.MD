# BookNest

![Library](https://i.imgur.com/hGaRvQ5.jpg)

This project only for personal uses and education purpose..

## Prerequisites

Before you begin, ensure you have met the following requirements:

1. Golang 1.20 or above: [Installation Guide](https://golang.org/doc/install)
2. PostgreSQL 13 or above: [Download and Installation](https://www.postgresql.org/download/)

## Installation

1. Clone the repository:

   ```sh
   git clone https://github.com/dpurbosakti/BookNest.git
   cd BookNest
   ```

2. Ensure the uuid-ossp extension is enabled in your PostgreSQL database:

   Run the following SQL command:
   ```sh
   CREATE EXTENSION IF NOT EXISTS "uuid-ossp";
   ```

3. Run the migration to set up the database:
   ```sh
   go run . migrate_up

   ```

## Run The App

To run the app you need to run this following command
```sh
go run . serve
```
**Note:** When running the application using the `go run . serve` command, make sure you are in the root directory of the cloned repository. This is the directory where the `main.go` file is located. Running the command from this directory ensures that the application is correctly started and served.


## More Commands
```sh
go run . migrate_down #drop all tables
go run . seeder # seeding book table
```

## API Supports

- AUTH
   - [x] Login
   - [x] Login By Google
   - [x] Login By Twitter
   - [x] Login By Github
- USERS
   - [x] Create User
   - [x] Verify Email
   - [x] Refresh Verification Code 
   - [x] Get List
   - [x] Get Detail
   - [x] Update
   - [x] Delete
- BOOKS
   - [x] Create Book
   - [x] Get List
   - [x] Get Detail
   - [x] Update
   - [x] Delete
   - [x] Return
- ORDERS
   - [x] Create Order
   - [x] Accept
   - [x] Reject
   - [x] Get List
   - [x] Get Detail
- ADDRESS
   - [x] Create Address
   - [x] Get Detail
   - [x] Update
- COURIERS
   - [x] Get Couriers From Biteship
   - [x] Get List
   - [x] Check Rates

