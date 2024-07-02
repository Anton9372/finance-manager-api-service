# Finance-manager API service

This application is a web service for managing and counting the user's financial transactions

## Description

The service provides api-gateway for managing transaction categories, financial transactions and calculating statistics

The service consists of several microservices communicating with each other using the http protocol:
- [user-service](https://github.com/Anton9372/user-service) - service for managing users

- [operation-service](https://github.com/Anton9372/operation-service) - service for managing categories and transactions

- [stats-service](https://github.com/Anton9372/stats-service) - service for obtaining statistics and report on financial operations

Detailed information about the api can be found at `http://localhost:10000/swagger`

## List of technologies
- Golang net/http
- PostgreSQL
- Docker
- JWT