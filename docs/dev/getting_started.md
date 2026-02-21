# Getting Started with Nexus Super Node

Welcome, Developer! This guide will help you set up a local instance of the Nexus Super Node for development and testing.

## 📋 Prerequisites

Before you begin, ensure you have the following installed:

*   **Go**: v1.22 or higher
*   **Docker & Docker Compose**: For running dependent services (Redpanda, TiDB, MinIO).
*   **Node.js & pnpm**: For the web portal (optional, if working on frontend).
*   **Temporal CLI**: For monitoring workflows.
    ```bash
    brew install temporal
    ```
*   **Redpanda (rpk)**: For managing topics.
    ```bash
    brew install redpanda-data/tap/redpanda
    ```

---

## 🚀 Quick Start (Local)

### 1. Clone the Repository
```bash
git clone https://github.com/your-org/nexus-super-node.git
cd nexus-super-node
```

### 2. Start Infrastructure
We use Docker Compose to spin up the "Nervous System" and "Data Layer".

```bash
# Starts Redpanda, TiDB, Redis, MinIO, Temporal, and LiveKit
docker-compose up -d
```
*Wait about 30-60 seconds for all services to initialize.*

**Note on LiveKit:**
For local development (especially on macOS), we map LiveKit ports directly to the host:
- **7880**: HTTP/WebSocket (API & Signal)
- **7881**: TCP (RTC)
- **7882**: UDP (RTC)
- **5060**: SIP (VoIP)

Ensure these ports are free on your machine.

### 3. Initialize the Database
Run the migration scripts to set up the TiDB schema.

```bash
go run cmd/migrate/main.go up
```

### 4. Configure Environment
Copy the example config and adjust if necessary.

```bash
cp config/config.example.yaml config/config.yaml
```
*Tip: The default config is pre-tuned for the local docker-compose setup.*

### 5. Generate Node Identity (DePIN)
A Super Node needs a cryptographic identity to participate in the network.

```bash
# Generates a new Ed25519 private key in ./data/keystore/
go run cmd/nexus-cli/main.go keygen
```
*Note: Keep this key safe! It controls your staked tokens and reputation.*

### 6. Run the Super Node
Start the core Go service.

```bash
go run cmd/nexus-super-node/main.go
```
You should see logs indicating connection to Redpanda and Temporal.

---

## 🧪 Verifying the Installation

### Check Health
Visit `http://localhost:8080/health` in your browser. You should see `{"status": "ok"}`.

### Test Event Flow
1.  **Produce a Market Signal**:
    ```bash
    rpk topic produce market-data -n 1
    # Type: {"symbol": "BTC", "price": 95000}
    ```
2.  **Check Logs**:
    The Super Node terminal should show: `Signal Received: BTC $95000`.

---

## 🛠 Project Structure

-   `/cmd`: Entry points (Main server, Migrations).
-   `/internal/core`: Business logic (Services, Domain Models).
-   `/internal/adapters`: Integrations (Redpanda, TiDB, Matrix).
-   `/internal/workflow`: Temporal Workflow definitions.
-   `/agents`: Pre-compiled Wasm agents (for testing).

## 📚 Next Steps
-   **[Building Your First Agent](agents.md)**: Learn how to write logic that runs on the node.
-   **[API Reference](api.md)**: Explore the GraphQL and REST endpoints.
