# Rinha de Backend

<!--toc:start-->

- [Rinha de Backend](#rinha-de-backend)
  - [Objective](#objective)
- [API endpoints](#api-endpoints)
  - [/POST /payments](#post-payments)
  - [/GET /payments-summary](#get-payments-summary)
- [Payment Processors endpoints](#payment-processors-endpoints)
<!--toc:end-->

## Objective

The main idea off the project is to have our backend as an proxy between two payment processor services. The catch of the project is that both of the services will suffer instabilities, meaning that they will become unavailable for some time and our target is to process the maximum amount of payments. We will receive an percentage of each successful transaction.

## API endpoints

### /POST /payments

```json
{
  "correlationId": "4a7901b8-7d26-4d9d-aa19-4dc1c7cf60b3",
  "amount": 19.9
}
```

- HTTP 2XX 0 -> Valid response
- `correlationId` type UUID.
- `amount` type decimal.

### /GET /payments-summary

```json
{
  "default": {
    "totalRequests": 43236,
    "totalAmount": 415542345.98
  },
  "fallback": {
    "totalRequests": 423545,
    "totalAmount": 329347.34
  }
}
```

- Query Parameters:
- `from` optional field in ISO UTC.
- `to` optional filed in ISO UTC.

- `default.totalRequests` type int.
- `default.totalAmount` type decimal.
- `fallback.totalRequests` type int.
- `fallback.totalAmount` type decimal.

## Payment Processors endpoints
