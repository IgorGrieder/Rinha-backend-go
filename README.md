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

### POST /payments

```json
{
  "correlationId": "4a7901b8-7d26-4d9d-aa19-4dc1c7cf60b3",
  "amount": 19.9,
  "requestedAt": "2025-07-15T12:34:56.000Z"
}
```

- Response HTTP 200 - OK

```json
{
  "message": "payment processed successfully"
}
```

correlationId type UUID.
amount type decimal.
requestedAt timestamp in ISO format (UTC).

Response:
message is an always-present field of type string.

### GET /payments/service-health

Response HTTP 200 - OK

```json
{
  "failing": false,
  "minResponseTime": 100
}
```

Request:

There are no request parameters. However, this endpoint imposes a rate limit of 1 call every 5 seconds. If this limit is exceeded, you will receive an HTTP 429 - Too Many Requests error response.

Response:

- failing is an always-present boolean field that indicates if the Payments endpoint is available. If it returns true, requests to the Payments endpoint will receive HTTP 5XX errors.
- minResponseTime is an always-present integer indicating the best possible response time (in milliseconds) for the Payments endpoint. For example, if the value is 100, no response will be faster than 100ms.
