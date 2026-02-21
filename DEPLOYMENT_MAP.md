# Modern Nexus Super Node Deployment Map (DPIN Architecture)

This document outlines the architecture for the **Nexus Super Node**, designed as a robust, scalable, and decentralized infrastructure hub. It leverages **Temporal** for orchestration, **LiveKit** for real-time media, and **Rust Wasm** for high-performance MLOps and Edge Intelligence.

## 🏗️ Architecture Layers

### 1. **Infrastructure Layer (State & Events)**
*   **PostgreSQL 14:** Durable storage for Temporal, Hasura, and Nexus business data.
*   **TiDB:** Distributed SQL database for high-scale relational data.
*   **Redpanda:** High-performance, Kafka-compatible event streaming for decoupled communication.
*   **Redis:** Fast key-value store for LiveKit room state and caching.

### 2. **Media & Communication Layer**
*   **LiveKit Server:** WebRTC SFU for high-quality audio/video.
*   **LiveKit SIP:** Native bridge for telephony. Features: Call Queuing, Eavesdropping, and AI Agent Transfers.

### 3. **Orchestration Layer (The "Nervous System")**
*   **Temporal Server:** Manages stateful, long-running workflows (Sagas). Ensures reliability and retries (e.g., Training Loops, Deployment Pipelines).
*   **Temporal UI:** Visibility into workflow execution.

### 4. **Application Layer (The "Brain")**
*   **Nexus Super Node (Go):** The core service.
    *   **Embedded Wasm Logic:** Runs internal plugins and business rules compiled to Wasm.
    *   **Temporal Worker:** Executes workflow activities (e.g., `TrainingWorkflow`, `DeploymentWorkflow`).
    *   **Event Consumer:** Reacts to Redpanda events for real-time processing.

### 5. **Edge/MLOps Layer (The "Capabilities") - *UPDATED***
*   **Wasm MLOps (Rust):**
    *   **Reward Function:** A Rust-based Wasm module running directly within the Data Pipeline (Redpanda Connect).
    *   **Logic:** Filters high-quality training data based on feedback, syntax, and reasoning depth.
    *   **Efficiency:** Eliminates heavy Python containers for data preprocessing; runs at near-native speed.
*   **Inference Engine:**
    *   **llama-server:** Runs GGUF models (e.g., Llama 3, Mistral) with low latency.
    *   **Unsloth Integration:** Future-proofed for fine-tuning loops orchestrated by Temporal.

---

## 📋 Deployment Checklist

- [ ] **Docker Environment:**
    - [x] Docker Desktop installed.
    - [x] "Use containerd for pulling and storing images" enabled.
    - [x] "Enable Wasm" feature enabled.

- [ ] **Configuration:**
    - [ ] `docker-compose.yml` validated (Frontend removed, Wasm Agent enabled).
    - [ ] LiveKit keys and SIP config validated.
    - [ ] `mlops-wasm` compiled to `.wasm` and mounted.

- [ ] **Execution:**
    - [ ] Build `nexus-super-node`.
    - [ ] Compile Rust Wasm modules (`cargo build --target wasm32-wasi`).
    - [ ] Start Core Infra (DB, Redis, Redpanda).
    - [ ] Start Temporal & LiveKit.
    - [ ] Start Super Node.

- [ ] **Verification:**
    - [ ] Verify Temporal UI access (localhost:8081).
    - [ ] Verify Redpanda Connect is processing events with Wasm.
    - [ ] Verify Super Node connection to Temporal.

---

## 🚀 Strategic Vision: DPIN Ecosystem
In a Decentralized Physical Infrastructure Network (DPIN), this Super Node acts as a **Gateway**. 
*   It orchestrates resources (Compute, Bandwidth, Storage).
*   **Wasm Agents** allow third-party developers to deploy "Capabilities" to your node securely (sandboxed) and efficiently.
*   **Temporal** ensures that complex interactions (e.g., "Rent this GPU for 5 mins", "Process this audio stream") are tracked and paid for reliably.
