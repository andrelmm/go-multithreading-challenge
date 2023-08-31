# Go Multithreading Challenge

This is a small Go project that demonstrates how to make simultaneous queries to two different ZIP code (CEP) APIs and return the response from the fastest API.

### Prerequisites

- Go (Golang)

## How to Run

1. Clone this repository to your local directory:

   ```sh
   git clone https://github.com/andrelmm/go-multithreading-challenge.git
   cd go-multithreading-challenge

2. Run the project:

    ```sh
    go run main.go

3. Follow the instructions to input a valid ZIP code (CEP) when prompted.

## Features

- Makes simultaneous requests to two ZIP code (CEP) query APIs.
- Returns the response from the fastest API in JSON format.
- Properly handles connection errors, response reading, and timeouts.
