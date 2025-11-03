# High-Level Design (HLD) — Customer Plan Management Microservice

## 1. Overview

**Project:** Customer Plan Management Service  
**Language:** Golang (Go 1.21+)  
**Target Platform:** Rakuten Mobile BSS Platform  
**Purpose:**  
To manage customer mobile plans — including creation, updates, subscriptions, and event propagation to downstream systems such as Billing, Inventory, and Notifications.

This service is a core component of the Business Support System (BSS) and demonstrates:
- Microservice design and domain modeling  
- Production-grade Go coding standards  
- Resilience, observability, and distributed-systems readiness  
- Integration with Kafka and PostgreSQL  

---

## 2. Architecture Overview

### 2.1 Logical Components

| Component | Responsibility |
|------------|----------------|
| **API Layer** | Exposes RESTful endpoints for Plan and Subscription management |
| **Service Layer** | Core business logic, validation, and orchestration |
| **Repository Layer** | Data access and persistence (PostgreSQL) |
| **Event Producer** | Publishes Kafka messages for plan/subscription lifecycle events |
| **Event Consumer (optional)** | Processes inbound Kafka events for sync/update operations |
| **Observability** | Handles metrics, structured logging, tracing, and health checks |

### 2.2 Diagram (Text Representation)

```
      ┌──────────────────────────┐
      │   Client (BSS UI / API)  │
      └────────────┬─────────────┘
                   │
          REST over HTTP (JSON)
                   │
     ┌─────────────┴─────────────────┐
     │  Plan Management Microservice │
     │  ───────────────────────────  │
     │  • API Layer (chi)            │
     │  • Service Layer              │
     │  • Repository Layer (pgx)     │
     │  • Kafka Producer/Consumer    │
     │  • Prometheus + OpenTelemetry │
     └──────────────┬────────────────┘
                    │
  ┌─────────────────┼───────────────────┐
  │                 │                   │
┌─────────────┐┌─────────────┐┌────────────────┐
│ PostgreSQL  ││ Kafka Broker││ Downstream Svc │
│ (State)     ││ (Events)    ││ Billing, Inv.  │
└─────────────┘└─────────────┘└────────────────┘
```


## 3. API Design

### 3.1 REST Endpoints

| Method | Endpoint | Description |
|---------|-----------|-------------|
| **POST** | `/plans` | Create a new mobile plan |
| **GET** | `/plans` | List all active plans |
| **GET** | `/plans/{id}` | Fetch details of a plan |
| **PUT** | `/plans/{id}` | Update an existing plan |
| **POST** | `/customers/{customer_id}/subscribe` | Subscribe a customer to a plan |
| **POST** | `/customers/{customer_id}/unsubscribe` | Unsubscribe a customer |
| **GET** | `/customers/{customer_id}/subscriptions` | Get all subscriptions for a customer |
| **GET** | `/healthz` | Liveness probe |
| **GET** | `/readyz` | Readiness probe |
| **GET** | `/metrics` | Prometheus metrics endpoint |

### 3.2 Sample Payloads

**Create Plan Request**
```json
{
  "code": "RM-UL-30D",
  "name": "Unlimited 30 Days",
  "price_cents": 1999,
  "duration_days": 30,
  "data_mb": 30720
}
```

**Subscribe Request**
```json
{
  "plan_id": "uuid",
  "start_date": "2025-11-03",
  "auto_renew": true
}
```

**Error Response**
```json
{
  "error": {
    "code": "SUBSCRIPTION_ALREADY_EXISTS",
    "message": "Customer already subscribed to this plan"
  }
}
```
## 4. Data Model

### 4.1 Entity Relationship Overview

```nginx
Customer ───< CustomerSubscription >─── Plan
```

### 4.2 Tables
**plans**
| Column                  | Type        | Notes                  |
| ----------------------- | ----------- | ---------------------- |
| id                      | UUID (PK)   | generated              |
| code                    | TEXT        | unique plan identifier |
| name                    | TEXT        | plan name              |
| price_cents             | BIGINT      | price in minor units   |
| currency                | CHAR(3)     | e.g., INR              |
| duration_days           | INT         | plan validity          |
| data_mb                 | BIGINT      | included data          |
| active                  | BOOLEAN     | status flag            |
| created_at / updated_at | timestamptz | audit fields           |

**customer_subscriptions**
| Column                  | Type        | Notes                      |
| ----------------------- | ----------- | -------------------------- |
| id                      | UUID (PK)   | generated                  |
| customer_id             | UUID        | external system ID         |
| plan_id                 | UUID        | FK to plans                |
| start_date / end_date   | DATE        | computed duration          |
| status                  | TEXT        | ACTIVE, CANCELLED, EXPIRED |
| auto_renew              | BOOLEAN     | auto renewal enabled       |
| created_at / updated_at | timestamptz | audit fields               |

**plan_events**
| Column      | Type        | Notes                       |
| ----------- | ----------- | --------------------------- |
| id          | BIGSERIAL   | PK                          |
| event_type  | TEXT        | e.g. `subscription.created` |
| resource_id | UUID        | plan/subscription id        |
| payload     | JSONB       | serialized event            |
| created_at  | timestamptz | timestamp                   |

## 5. Event Architecture
### 5.1 Topics
| Topic                 | Publisher    | Description                       |
| --------------------- | ------------ | --------------------------------- |
| `plan.events`         | Plan service | Plan lifecycle events             |
| `subscription.events` | Plan service | Subscription lifecycle events     |
| `billing.requests`    | (optional)   | Events consumed by billing system |

### 5.2 Example Message (JSON)
```json
{
  "event_id": "uuid",
  "event_type": "subscription.created",
  "timestamp": "2025-11-03T09:00:00Z",
  "subscription": {
    "id": "uuid",
    "customer_id": "uuid",
    "plan_id": "uuid",
    "start_date": "2025-11-03",
    "status": "ACTIVE",
    "auto_renew": true,
    "price_cents": 1999
  }
}
```

### 5.3 Delivery Semantics

1. At-least-once event publishing using retry & backoff.
2. Events contain event_id for idempotency handling downstream.
3. Consumers use offset commits and deduplication based on event_id.

## 6. Observability
| Aspect               | Implementation                                                          |
| -------------------- | ----------------------------------------------------------------------- |
| **Metrics**          | Prometheus (`/metrics`) via `promhttp` middleware                       |
| **Tracing**          | OpenTelemetry SDK for Go (HTTP + DB spans)                              |
| **Logging**          | `zap` or `zerolog`, JSON structured logs with `trace_id` & `request_id` |
| **Health Checks**    | `/healthz` (DB + Kafka ping) and `/readyz` (migration applied)          |
| **Business Metrics** | `subscriptions_created_total`, `subscriptions_cancelled_total`          |


**Example Prometheus metric:**
```lua
http_requests_total{path="/plans",method="GET",status="200"} 42
```

## 7. Error Handling & Resilience
| Concern                       | Strategy                                   |
| ----------------------------- | ------------------------------------------ |
| **Validation errors**         | Return 400 with structured JSON            |
| **Duplicate requests**        | Idempotency key for subscription endpoints |
| **Transient DB errors**       | Retry with exponential backoff             |
| **Kafka failure**             | Queue and retry with async worker          |
| **Downstream unavailability** | Circuit breaker pattern (future scope)     |

**Error format**
```json
{ "error": { "code": "PLAN_NOT_FOUND", "message": "Plan does not exist" } }
```

## 8. Testing Strategy
| Type                    | Scope                                                | Tools                      |
| ----------------------- | ---------------------------------------------------- | -------------------------- |
| **Unit Tests**          | Service layer: business rules, date calc, validation | `testing` pkg              |
| **Integration Tests**   | HTTP + DB end-to-end                                 | `httptest`, Docker Compose |
| **Kafka Tests**         | Mock producer verification                           | Mock writer via `kafka-go` |
| **Observability Tests** | Verify `/metrics` and `/healthz`                     | curl assertions            |

**Example test case matrix:**
| Scenario                       | Expected Behavior       |
| ------------------------------ | ----------------------- |
| Create plan with valid payload | 201 Created             |
| Subscribe to plan twice        | 409 Conflict            |
| Invalid plan_id in subscribe   | 404 Not Found           |
| Kafka produce fails            | Warning log, no crash   |
| DB down                        | 503 Service Unavailable |

## 9. Deployment & Configuration

### 9.1 Environment Variables
| Key                           | Description                 |
| ----------------------------- | --------------------------- |
| `DB_URL`                      | Postgres connection string  |
| `KAFKA_BROKERS`               | Comma-separated broker list |
| `SERVICE_PORT`                | HTTP port (default 8080)    |
| `OTEL_EXPORTER_OTLP_ENDPOINT` | OpenTelemetry collector     |
| `PROMETHEUS_PORT`             | Metrics endpoint port       |

### 9.2 Docker Compose (Development)
```yaml
version: "3.8"
services:
  postgres:
    image: postgres:15
    environment:
      POSTGRES_USER: admin
      POSTGRES_PASSWORD: admin
      POSTGRES_DB: plansvc
    ports: [ "5432:5432" ]

  kafka:
    image: bitnami/kafka:3
    ports: [ "9092:9092" ]

  planservice:
    build: .
    ports: [ "8080:8080" ]
    depends_on: [ postgres, kafka ]
    environment:
      DB_URL: postgres://admin:admin@postgres:5432/plansvc?sslmode=disable
      KAFKA_BROKERS: kafka:9092
```

## 10. Security Considerations
| Concern          | Mitigation                                                       |
| ---------------- | ---------------------------------------------------------------- |
| Input validation | Strong JSON schema validation                                    |
| Authentication   | JWT or OAuth2 token validation middleware (stubbed in prototype) |
| Authorization    | Scoped API keys or customer ownership check                      |
| Secrets          | Environment variables or Vault integration                       |
| Network          | TLS termination at ingress layer                                 |

## 11. Non-Functional Requirements
| Aspect               | Target                                            |
| -------------------- | ------------------------------------------------- |
| **Availability**     | 99.9%                                             |
| **Latency**          | <100ms per API call (95th percentile)             |
| **Throughput**       | 1k requests/sec (scalable horizontally)           |
| **Data consistency** | Strong consistency in DB, eventual across systems |
| **Durability**       | Kafka ACK=all, Postgres replication-ready         |
| **Scalability**      | Stateless API, multi-instance friendly            |

## 12. Future Enhancements
1. Replace JSON with Protobuf for Kafka messages.
2. Introduce Saga orchestration for multi-system plan activation.
3. Add RBAC and per-tenant rate limiting.
4. Implement outbox pattern for guaranteed event delivery.
5. Extend to GraphQL or gRPC interface.
6. Automated CI/CD pipeline with test coverage thresholds.
