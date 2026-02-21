# Contributing to Nexus Super Node

Thank you for your interest in contributing to the Nexus Super Node! This document provides the technical context and guidelines needed to modify the core Go codebase.

## 🏗 Architecture Overview (The "Onion" Model)

We follow **Hexagonal Architecture** (Ports & Adapters) to keep business logic isolated from infrastructure.

### Directory Structure

```
/internal
  /core             # Pure Business Logic (No external deps)
    /domain         # Entities (e.g., Agent, User, Auction)
    /ports          # Interfaces (Repositories, Service definitions)
    /services       # Use Cases (e.g., AgentRunner, BidManager)
  /adapters         # Infrastructure Implementations
    /persistence    # Database (TiDB, Postgres)
    /messaging      # Redpanda / Kafka
    /workflow       # Temporal Client
    /gateway        # HTTP/WebSocket/GraphQL Handlers
```

### Key Rules
1.  **Dependencies Flow Inward**: `adapters` -> `core`. The Core never imports from Adapters.
2.  **Interfaces First**: Define behavior in `core/ports` before implementing in `adapters`.
3.  **Dependency Injection**: We use `uber-go/fx` for wiring components.

---

## 🛠 Development Workflow

### 1. Setting up the Dev Environment
Ensure you have the "Standard Stack" running:
```bash
docker-compose up -d
```

### 2. Making Changes
*   **Adding a Database Field**:
    1.  Modify `internal/core/domain/models.go`.
    2.  Create a migration in `db/migrations`.
    3.  Update `internal/adapters/persistence/tidb`.
*   **Adding a New API Endpoint**:
    1.  Define the Service method in `internal/core/ports/services.go`.
    2.  Implement logic in `internal/core/services`.
    3.  Expose via `internal/adapters/gateway` (Echo Handler or Hasura Action).

### 3. Testing
We distinguish between Unit Tests and Integration Tests.

*   **Unit Tests** (Fast, Mocked):
    *   Place in the same package (e.g., `service_test.go`).
    *   Use `stretchr/testify/mock` to mock Ports.
    ```bash
    go test ./internal/core/...
    ```

*   **Integration Tests** (Slow, Real DB/Redpanda):
    *   Requires `docker-compose` to be running.
    *   Located in `tests/integration` or `_test.go` files in adapters.
    ```bash
    go test -tags=integration ./internal/adapters/...
    ```

### 4. Linting
We use `golangci-lint`.
```bash
golangci-lint run
```

---

## 🧩 Adding a New Adapter

If you want to support a new database or message queue:
1.  Create a new folder in `internal/adapters/` (e.g., `scylla`).
2.  Implement the interfaces defined in `internal/core/ports`.
3.  Register the new module in `cmd/nexus-super-node/main.go` using `fx.Provide`.

## 📜 Code Style
*   **Errors**: Wrap errors with context. `fmt.Errorf("failed to process order: %w", err)`.
*   **Logging**: Use structured logging (`slog` or `zap`). Do not use `fmt.Println`.
*   **Configuration**: All config should come from `config/config.yaml` and be injected.

## 🤝 Pull Request Process
1.  Fork the repo.
2.  Create a feature branch (`feat/add-solana-support`).
3.  Ensure tests pass.
4.  Submit PR with a description of *Why* and *How*.
