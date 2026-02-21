# High-Level System Overview

## 1. The Nexus Super Node Concept
**"One Core, Infinite Interfaces"**

The Nexus Super Node is not merely a backend API; it is a **Decentralized Operating System for AI Agents**. In the traditional web architecture, a backend serves a specific frontend (e.g., a React app). In the Nexus architecture, the Super Node acts as a universal **Logic Hub** that serves users wherever they exist—whether on Telegram, Farcaster, a CLI, or a VR environment.

### Core Philosophy
1.  **Sovereignty**: Users own their data (Identity, Wallet, History) via cryptographic keys, not database entries.
2.  **Continuity**: An agent started on Telegram can continue its task on the Web Portal without losing context. The state lives in the Super Node, not the client.
3.  **Fractal Scalability**: The network grows by replicating "Super Nodes". Each node is self-contained but can federate with others to form a global mesh.

---

## 2. Architectural Pillars

### A. The Compute Plane (Wazero & Temporal)
*Why not just Docker containers?*
We use **WebAssembly (Wasm)** via **Wazero** for agent logic because:
-   **Security**: Wasm provides a strict sandbox. An agent cannot access the host filesystem unless explicitly allowed.
-   **Portability**: A compiled agent (`agent.wasm`) runs identically on a MacBook, a Linux Server, or an Edge Device.
-   **Startup Time**: Wasm modules instantiate in microseconds, whereas containers take seconds.

**Temporal** acts as the "OS Scheduler". It guarantees that workflows (like "Monitor BTC price for 24 hours") run to completion, even if the server crashes. It provides **Durable Execution**.

### B. The Data Plane (TiDB & IPFS)
*Why Hybrid?*
Data comes in different "temperatures" and requires different storage engines:
-   **Hot Data (Redis)**: User session state, ephemeral WebSocket signals.
-   **Smart Data (TiDB)**: This is our **HTAP (Hybrid Transactional/Analytical Processing)** engine. It stores:
    -   **Transactional**: User balances, active orders (Row-store).
    -   **Analytical**: Historical market trends for AI training (Column-store).
    -   **Vector**: Semantic embeddings for Agent Memory (Vector Search).
-   **Cold Content (IPFS)**: Large blobs that must be tamper-proof. Agent code (`.wasm`), AI Models (`.gguf`), and trade logs are stored here. Content-Addressing (CID) ensures trust.

### C. The Nervous System (Redpanda)
*More than a Message Queue*
**Redpanda** is the single source of truth for events. In Nexus, **nothing happens without an event**.
-   **Signal Distribution**: Market data enters Redpanda and is "fan-out" broadcasted to thousands of agents instantly.
-   **Event Sourcing**: We don't just update a database row; we append an event (`OrderPlaced`). This allows us to "replay" history to train better AI models.

### D. The AI Learning Center (Unsloth & Rivet)
*The Self-Improving Loop*
-   **Rivet**: Provides a visual programming environment for defining agent logic (DAGs). It allows non-coders to build complex behaviors.
-   **Unsloth**: A specialized fine-tuning pipeline. When an agent fails a task, the data is tagged. Unsloth uses this data to fine-tune a Small Language Model (SLM) specifically for that task, making the agent smarter and cheaper over time.

### E. The Secure Compute Plane (Cocoon Integration)
*The Privacy Shield*
For tasks involving sensitive user data (Private Keys, Personal Messages), the Super Node offloads execution to **Cocoon Workers**. 
- **Confidentiality**: Uses Intel TDX / AMD SEV-SNP to ensure data is encrypted even in RAM.
- **Verifiability**: Remote Attestation proves that the AI model hasn't been tampered with.
- **Settlement**: Integrated with the TON blockchain for seamless micropayments.

### F. The Real-Time Communication Plane (LiveKit & Matrix)
*The Senses & Voice*
For agents to truly interact with humans, they need to see, hear, and speak.
-   **LiveKit (Audio/Video)**: Provides WebRTC infrastructure for agents to join voice calls, process video feeds, and stream screen data.
    -   *Use Case*: An AI Tutor agent seeing your screen and guiding you through code.
-   **Matrix (Messaging/Signaling)**: A decentralized communication layer that replaces proprietary APIs (like Telegram/Discord bots).
    -   *Use Case*: An agent joining a group chat, encrypting messages end-to-end, and coordinating with other agents securely.

### G. Internal DePIN Services (The Hidden Machinery)
To function as a sovereign node in the DePIN mesh, the Super Node runs these background daemons:
1.  **Identity Manager**: Manages the node's **Private Key** (Ed25519/Secp256k1). It signs every auction bid and computation result to prove authenticity.
2.  **Resource Sentinel**: A background worker that continuously polls the host's GPU (via NVML) and CPU stats. It feeds this data to the Auction Engine so the node never over-commits.
3.  **Auction Engine**: Listens to the global `auction.compute` GossipSub topic. If the Sentinel reports available capacity, it automatically calculates a price and broadcasts a **Signed Bid** within milliseconds.

---

## 3. The Fractal DevOps Model

The architecture is designed to be **Fractal**. This means the structure of the whole network resembles the structure of a single node.

### Hierarchy of Nodes
1.  **Edge Node (User-Run)**:
    -   Runs on a laptop or Raspberry Pi.
    -   Hosts personal agents.
    -   Connects to a Super Node for heavy compute (LLM inference).
2.  **Super Node (Provider-Run)**:
    -   Runs on high-end servers (GPU clusters).
    -   Provides services (Inference, Storage, Routing) to thousands of Edge Nodes.
    -   Earns network tokens for providing resources.
3.  **The Nexus Mesh (Global)**:
    -   The collection of all Super Nodes connected via **libp2p**.
    -   Uses a Distributed Hash Table (DHT) for peer discovery.
    -   "Where is User X?" -> DHT Query -> "User X is on Super Node Y".

### GitOps & Deployment
We treat infrastructure as code.
-   **Repo-Driven**: Pushing to the `main` branch of an Agent's repo triggers a build pipeline.
-   **Wasm Artifact**: The code is compiled to Wasm and pinned to IPFS.
-   **Registry Update**: The new IPFS CID is registered on the Super Node via a smart contract or system transaction.
-   **Instant Rollout**: All running instances of that agent update automatically in the next execution cycle.

---

## 4. Global Network Topology

```mermaid
graph TD
    subgraph "Region: North America"
        SN1[Super Node 1]
        SN2[Super Node 2]
        SN1 <--> |"GossipSub (Sync)"| SN2
    end

    subgraph "Region: Europe"
        SN3[Super Node 3]
    end

    SN1 <--> |"libp2p / WAN"| SN3

    subgraph "Edge Layer (Clients)"
        ClientA[User A (Telegram)]
        ClientB[User B (Web)]
        ClientC[User C (CLI)]
    end

    ClientA -.-> |"WebSocket"| SN1
    ClientB -.-> |"gRPC"| SN1
    ClientC -.-> |"QUIC"| SN3
```
