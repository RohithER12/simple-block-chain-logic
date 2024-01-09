# Simple Blockchain Logic

## Overview

This project implements a simple blockchain library in Go, allowing users to manage a decentralized ledger for book checkouts. It leverages the Gin framework for handling HTTP requests and responses.

## Features

- **Blockchain Management:** Create, read, and validate blocks in a decentralized blockchain.
- **Book Checkouts:** Record book checkouts with associated information.

## Installation

1. **Clone the repository:**

    ```bash
    git clone https://github.com/RohithER12/simple-block-chain-logic.git
    ```

2. **Change into the project directory:**

    ```bash
    cd simple-block-chain-logic
    ```

3. **Install dependencies:**

    ```bash
    go mod download
    ```

## Usage

1. **Run the application:**

    ```bash
    go run main.go
    ```

2. **Access the API using the provided endpoints:**

    - **GET /:** Retrieves the entire blockchain.
    - **POST /:** Adds a new block to the blockchain with book checkout information.
    - **POST /new:** Adds a new book to the system.

## Endpoints

- **GET /:**
  - Retrieves the entire blockchain.

- **POST /:**
  - Adds a new block to the blockchain with book checkout information.

- **POST /new:**
  - Adds a new book to the system.

## Blockchain Structure

The blockchain consists of interconnected blocks, each containing information about a book checkout. The blocks are linked using cryptographic hashes. The structure of a block is as follows:

```go
type Block struct {
    Pos       int
    Data      BookCheckOut
    TimeStamp string
    Hash      string
    PrevHash  string
}
