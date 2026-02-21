# Nexus Super Node - Documentation Roadmap & Task List

This document outlines the comprehensive plan to create full documentation for the Nexus Super Node project. It covers architecture, developer guides, user manuals, and API specifications.

## Phase 1: Core Stack (Current Focus)
- [x] **Infrastructure Setup**
    - [x] Docker Compose / K3s configuration.
    - [x] Redpanda Connect (Benthos) integration.
    - [x] TiDB + PostgreSQL setup.
- [x] **Service Integration**
    - [x] Go Backend (Nexus).
    - [x] Hasura GraphQL Engine.
    - [x] Temporal Workflow Engine.
    - [x] LiveKit & Matrix Integration.
- [x] **Developer Guides**
    - [x] Contributing Guide (`docs/dev/contributing.md`).
    - [x] API Reference & Auth (`docs/dev/api.md`).
    - [x] Agent Testing & Debugging (`docs/dev/agents.md`).

## Phase 2: DePIN Evolution (Future Roadmap)
- [x] **DePIN Roadmap Document (`docs/architecture/depin_roadmap.md`)**
- [ ] **Network Operations (`docs/ops/networking.md`)**
    - [x] P2P Architecture Concepts (libp2p, QUIC).
    - [x] Discovery (DHT) & Rendezvous Protocols (`docs/architecture/protocols/discovery.md`).
    - [x] Agent-to-Agent Protocol (`docs/architecture/protocols/agent_protocol.md`).
- [ ] **Secure Compute Layer**
    - [x] Cocoon Integration Analysis (`docs/architecture/cocoon_integration.md`).
- [ ] **Node Operator Guide (Basic)**
    - [x] Hardware Requirements (`docs/ops/running_node.md`).

## 📅 Execution Plan
1. Create folder structure in `docs/`.
2. Generate Phase 1 diagrams (Mermaid).
3. Draft the "Getting Started" guide.
4. Document the Hasura <-> Go <-> Temporal interaction.
