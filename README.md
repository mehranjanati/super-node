# Super Node: Autonomous AI Agent Orchestrator

This project is a high-performance, event-driven architecture for autonomous AI agents, leveraging **Redpanda**, **Temporal**, and **Wasm** for low-latency MLOps and reasoning.

## 🏗 Architecture Overview

The system is designed as a **Super Node** that orchestrates the entire lifecycle of AI agents:
1.  **Ingestion**: Real-time event ingestion via Redpanda (Kafka-compatible).
2.  **Processing (The Brain)**: **Wasm-based** reward functions and data filtering running directly within the data pipeline (Redpanda Connect).
3.  **Orchestration**: **Temporal Workflows** manage long-running processes like training loops, model evaluation, and deployment.
4.  **Inference**: Low-latency inference using `llama.cpp` (GGUF models) or external LLM providers.
5.  **Real-time Comms**: LiveKit for audio/video/data streaming.

### Key Components

| Component | Technology | Role |
| :--- | :--- | :--- |
| **Event Bus** | Redpanda | High-throughput event streaming & storage. |
| **Data Pipeline** | Redpanda Connect | ETL, transformation, and **Wasm-based** reward scoring. |
| **Orchestrator** | Temporal | Reliable workflow execution (Training, Deployment). |
| **MLOps Logic** | Rust (Wasm) | Lightweight, portable reward functions & data curation. |
| **Inference** | llama-server | Local LLM inference (OpenAI compatible API). |
| **Backend** | Go (Echo + Uber Fx) | API Gateway, business logic, and system integration. |
| **Agent Orchestrator** | **VoltAgent** | Central nervous system managing MCP tools, Temporal workflows, and AI context. |

---

## ⚡ VoltAgent: The Reasoning Engine

**VoltAgent** is the core reasoning service that powers the Super Node. It is responsible for:
1.  **Tool Abstraction (MCP)**: Unifies local functions, remote MCP servers, and system workflows into a single tool interface.
2.  **Workflow Orchestration**: Triggers complex, long-running Temporal workflows (e.g., "Deploy Website", "Crypto Analysis") from simple chat commands.
3.  **Context Management**: (In Progress) Manages short-term and long-term memory for coherent conversations.

### Supported Tools
- **System Tools**:
  - `system__deploy_website`: Generates UI/Code, pushes to Git, builds Wasm, and deploys.
  - `system__crypto_analysis`: Deep market analysis pipeline with human-in-the-loop approval.
  - `system__human_handoff`: Escalate to a human operator via Matrix/LiveKit.
- **MCP Tools**: Connects to any Model Context Protocol (MCP) server.

---

## 🚀 Getting Started

### Prerequisites
- Docker & Docker Compose
- Go 1.21+
- Rust (for Wasm modules)

### Running the Stack

1.  **Start Infrastructure**:
    ```bash
    docker-compose up -d
    ```
    This spins up Redpanda, Temporal, Postgres, Redis, and other core services.

2.  **Compile Wasm Modules**:
    Navigate to `mlops-wasm/` and build the reward function:
    ```bash
    cd mlops-wasm
    cargo build --target wasm32-wasi --release
    # Ensure the output .wasm file is mounted to Redpanda Connect container
    ```

3.  **Run Super Node**:
    ```bash
    go run main.go
    ```

---

## 🧠 MLOps & Autonomous Training

The core innovation is the **Wasm-based Autonomous Loop**:

1.  **Data Collection**: User interactions and agent logs are published to Redpanda topics.
2.  **Wasm Processing**:
    - The `mlops_reward.wasm` module (running in Redpanda Connect) analyzes events in real-time.
    - It calculates a **Reward Score** based on feedback, syntax, and reasoning depth.
    - Low-quality data is filtered out *before* storage.
3.  **Continuous Training**:
    - A Temporal Workflow triggers periodically.
    - It uses the high-quality dataset to fine-tune models (simulated Wasm optimization).
    - If successful, the new model is automatically swapped into the inference engine.

---

## 📂 Project Structure

- `config/`: Configuration for Redpanda, Temporal, and Pipelines.
- `internal/`: Go backend source code (DDD structure).
- `mlops-wasm/`: Rust source for Wasm reward functions.
- `docker-compose.yml`: Infrastructure definition.

---

## 🛠 Development

### Adding a New Wasm Processor
1.  Edit `mlops-wasm/src/lib.rs`.
2.  Compile to Wasm.
3.  Update `config/training_pipeline.yaml` to reference the new function.
4.  Restart Redpanda Connect.

### Modifying Workflows
- Workflows are defined in `internal/workflows/`.
- Use the Temporal UI (localhost:8081) to debug execution.
