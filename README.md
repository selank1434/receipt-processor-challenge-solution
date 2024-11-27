# Receipt Processor Go API with Docker Setup

This repository contains a Go-based API that processes receipts and calculates points. It is containerized using Docker.

## Prerequisites

Before proceeding, ensure you have the following installed:
- [Go](https://golang.org/dl/)
- [Docker](https://www.docker.com/get-started)

### `main.go` Code Overview

The code defines an API with the following functionality:
- Accepts POST requests at `/receipts/process` to process receipt data and calculate points.
- Responds to GET requests at `/receipts/{id}/points` to retrieve the points associated with a receipt.

## Step-by-Step Setup

### 1. Clone the Repository

### 2. Build using go 

- go mod init receipt-process-api
- go mod tidy

### 3. Open docker desktop(make sure the docker daemon is running)

### 4. Build the Docker Image

docker build -t receipt-process-api .

### 5. Start up the API using docker 

docker run -p 8080:8080 receipt-process-api

