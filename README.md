# Finance-manager API service

This application is a web service for managing and counting the user's financial transactions

## Description

The service provides api-gateway for managing transaction categories, financial transactions and calculating statistics

The service consists of several microservices:
- [user-service](https://github.com/Anton9372/user-service) - Manages user-related operations and communicates using gRPC. View the gRPC contracts [here](https://github.com/Anton9372/user-service-contracts).

- [operation-service](https://github.com/Anton9372/operation-service) - Manages categories and transactions through HTTP

- [stats-service](https://github.com/Anton9372/stats-service) - Provides statistics and reports on financial operations via HTTP.

Detailed information about the api can be found at `http://host:port/swagger`

## List of technologies

- Golang net/http
- gRPC
- Protocol Buffers
- PostgreSQL
- Docker
- JWT