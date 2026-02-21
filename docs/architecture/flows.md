# System Data Flows & Lifecycles

This document details the critical data paths within the Nexus Super Node. It covers how users authenticate, how agents react to the world, and how the system learns from its mistakes.

## 1. Authentication & Identity Federation
*How we bridge Web2 (Telegram/GitHub) and Web3 (Wallets) into a single identity.*

The Nexus Super Node uses a **Federated Identity Model**. A single user ID (DID) can be linked to multiple authentication providers.

### The Auth Sequence
1.  **Challenge Generation**: Client requests a login challenge (nonce).
2.  **Signing**:
    -   **Web3**: User signs the nonce with their Private Key (MetaMask/Phantom).
    -   **Telegram**: The Mini App provides `initData` signed by Telegram's bot token.
    -   **GitHub**: OAuth exchange provides a valid Access Token.
3.  **Verification**: The Super Node verifies the signature/token against the provider.
4.  **Session Creation**:
    -   If valid, the system looks up the `UserDID` in TiDB.
    -   A JWT is issued with claims: `sub: DID`, `roles: [pro, dev]`, `scope: [read:market, write:agent]`.

```mermaid
sequenceDiagram
    participant Client
    participant Gateway as Go Gateway
    participant Auth as Auth Service
    participant TiDB
    participant Provider as External (TG/GitHub)

    Client->>Gateway: POST /auth/login (provider="telegram", payload=initData)
    Gateway->>Auth: Verify Payload
    Auth->>Provider: Validate Signature (Public Key / API)
    Provider-->>Auth: Valid + User Metadata (username, id)
    
    Auth->>TiDB: Find DID by ProviderID
    alt User Exists
        TiDB-->>Auth: Returns DID
    else New User
        Auth->>TiDB: Create New DID + Link ProviderID
    end

    Auth-->>Gateway: Issue JWT (Access + Refresh)
    Gateway-->>Client: 200 OK { token: "eyJ..." }
```

---

## 2. Agent Awakening (The "Hot" Path)
*From Market Signal to Wasm Execution in milliseconds.*

This is the core loop. We do **not** poll the database. We react to events.

### The Trigger Chain
1.  **Ingestion**: Market data flows into Redpanda topics (`market.btc.usd`).
2.  **Filtering**: The **Signal Router** (in Go) reads the stream. It checks: *"Which agents are subscribed to BTC > 90k?"*
3.  **Orchestration**: If a condition matches, a Temporal Workflow is triggered.
4.  **Execution**: Temporal spins up a Wasm worker, injects the context, and runs the agent's logic.

```mermaid
sequenceDiagram
    participant Market as Market Feed
    participant RP as Redpanda
    participant Router as Signal Router
    participant Temporal
    participant Wasm as Agent Runtime
    participant TiDB

    Market->>RP: Publish { symbol: "BTC", price: 91000 }
    RP->>Router: Consume Event
    
    Router->>Router: Match Subscriptions (Bloom Filter)
    Note right of Router: "Agent_007 wants BTC > 90k"

    Router->>Temporal: ExecuteWorkflow("Agent_007", Payload)
    Temporal->>Wasm: Start Worker
    Wasm->>TiDB: Fetch Agent Memory (Vector Store)
    Wasm->>Wasm: Run Strategy Logic
    
    alt Decision = BUY
        Wasm->>RP: Publish { type: "ORDER", action: "BUY" }
    end
```

---

## 3. The AI Training Loop (The "Cold" Path)
*How the system gets smarter over time.*

Agents produce logs and outcomes. If an agent loses money or fails a task, we don't just log it; we use it as a negative training example.

### The Optimization Pipeline
1.  **Data Collection**: All agent decisions are logged to MinIO (via Redpanda Connect).
2.  **Evaluation**: A daily cron job (or real-time trigger) scores agent performance.
3.  **Fine-Tuning**:
    -   High-performing traces become "Golden Datasets".
    -   **Unsloth** is triggered to fine-tune the base SLM (Small Language Model) with these new examples.
4.  **Deployment**: The new model weights (LoRA adapter) are saved to IPFS and hot-swapped into the running agents.

```mermaid
sequenceDiagram
    participant Agent
    participant MinIO as MinIO (Logs)
    participant Eval as Evaluator Service
    participant Unsloth
    participant IPFS
    participant Registry as Model Registry

    Agent->>MinIO: Log Execution Trace (Input + Output + Outcome)
    
    loop Daily Optimization
        Eval->>MinIO: Scan for High/Low Performance
        Eval->>Eval: Curate Dataset (JSONL)
        Eval->>Unsloth: Start Fine-Tuning Job (BaseModel + Dataset)
        
        Unsloth->>Unsloth: Train (GPU Accelerated)
        Unsloth->>IPFS: Upload New LoRA Adapter
        IPFS-->>Unsloth: Return CID (QmHash...)
        
        Unsloth->>Registry: Update Agent Config (Model = NewCID)
    end
```
