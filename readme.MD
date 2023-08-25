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
   go run . migrate

   ```