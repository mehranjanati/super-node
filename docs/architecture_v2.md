# Nexus Super Node - Architecture V2

This document contains updated architectural diagrams incorporating **Hasura** (Data Access), **Unsloth** (Model Fine-tuning), and **Rivet** (Agent Logic).

## 1. Unified Request Flow (v2)
*Integrates Hasura for reads, Rivet for logic definition, and Unsloth for background learning.*

```mermaid
sequenceDiagram
    autonumber
    box "Omni-Client Layer" #e6f7ff
        participant User
        participant Client as Mini-App / Portal
        participant GQL as Hasura (GraphQL)
    end
    
    box "Super Node Edge" #fff0f6
        participant API as Go API / Gateway
        participant Auth as Identity Manager
    end
    
    box "Execution & Logic" #f9f0ff
        participant Temporal as Temporal (Orchestrator)
        participant Wasm as Wazero (Runtime)
        participant Rivet as Rivet Virtual Machine
    end

    box "Data & AI Ops" #fcffe6
        participant TiDB as TiDB (State/Vector)
        participant RP as Redpanda (Signals)
        participant Unsloth as Unsloth Service (GPU)
        participant IPFS as IPFS (Storage)
    end

    Note over User, GQL: Read Path (Fast)
    User->>Client: View Dashboard / Market Data
    Client->>GQL: Query: subscription { market_signals, agent_status }
    GQL->>TiDB: SQL Query (Live State)
    TiDB-->>Client: Real-time Data Stream

    Note over User, API: Write/Action Path
    User->>Client: Action: "Optimize Strategy"
    Client->>API: POST /execute (Intent, Params)
    API->>Auth: Verify Identity (Git/Crypto)
    API->>Temporal: Start Workflow "Agent_Optimization"

    Note over Temporal, Rivet: Logic Execution
    Temporal->>Wasm: Load Agent Container
    Wasm->>IPFS: Fetch Logic Graph (Rivet File .rivet-project)
    Wasm->>Rivet: Execute Logic Graph
    Rivet->>TiDB: Query Vector Memory (Context)
    
    Note over Rivet, Unsloth: Self-Improvement Loop
    alt Performance Below Threshold
        Rivet->>RP: Emit "Training_Required" Event
        RP->>Temporal: Trigger "FineTune_Model" Workflow
        Temporal->>Unsloth: Start Training (BaseModel + New Data)
        Unsloth->>TiDB: Fetch Historical Success/Fail Patterns
        Unsloth->>Unsloth: Fine-tune (LoRA/QLoRA)
        Unsloth->>IPFS: Save New Adapter
        Unsloth->>TiDB: Update Agent Config (Model Version)
    end

    Rivet-->>API: Action Result
    API-->>Client: Notification
```

## 2. The "Super Node" Stack - Component Map
*Shows where every service lives and how they interconnect.*

```mermaid
graph TD
    subgraph "Frontend & Access Layer"
        Portal[Web Portal (SvelteKit)]
        TG[Telegram Mini App]
        Frames[Farcaster Frames]
        Hasura[Hasura GraphQL Engine]
        
        Portal --> Hasura
        TG --> Hasura
        Frames --> GoAPI
    end

    subgraph "Core Services (The Nexus)"
        GoAPI[Go Gateway (Echo/gRPC)]
        Temporal[Temporal Server]
        Wasm[Wazero Runtime]
        RivetVM[Rivet Logic Runner]
        
        GoAPI --> Temporal
        Temporal --> Wasm
        Wasm --> RivetVM
    end

    subgraph "Event & Signal Fabric"
        RP[Redpanda (Kafka API)]
        Connect[Redpanda Connect (Benthos)]
        
        GoAPI --> RP
        Wasm --> RP
        RP --> Connect
    end

    subgraph "Data & Memory Layer"
        Redis[Redis (Hot Cache)]
        TiDB[TiDB (HTAP - SQL + Vector)]
        MinIO[MinIO (S3 - Raw Data)]
        IPFS[IPFS (Decentralized Content)]
        
        Hasura --> TiDB
        GoAPI --> Redis
        Connect --> MinIO
        Wasm --> IPFS
    end

    subgraph "AI & Compute Plane"
        Unsloth[Unsloth AI Service]
        GPU[(NVIDIA GPUs)]
        
        Unsloth --> GPU
        Temporal --> Unsloth
        Unsloth --> MinIO
        Unsloth --> IPFS
    end

    %% Key Relationships
    RivetVM -.-> |"Executes"| IPFS
    RivetVM -.-> |"Remembers"| TiDB
    Unsloth -.-> |"Reads Training Data"| TiDB
    Unsloth -.-> |"Saves Models"| IPFS
    Hasura -.-> |"Subscribes"| RP
```

## 3. Data Hierarchy & Lifecycle
*Detailed view of how data moves from Raw to Smart to Archived.*

```mermaid
graph TD
    Input((External Signals)) --> |"Ingest"| RP[Redpanda]

    subgraph "Hot (Milliseconds)"
        RP --> |"Stream"| Redis
        Redis --> |"Sub"| Client
    end

    subgraph "Warm (Seconds)"
        RP --> |"Process"| Temporal
        Temporal --> |"Write State"| TiDB
        Hasura --> |"Read"| TiDB
    end

    subgraph "Cold (Minutes/Hours)"
        RP --> |"Batch"| MinIO[MinIO (Raw Logs)]
        MinIO --> |"ETL"| Unsloth[Unsloth Training]
        Unsloth --> |"New Model"| IPFS
    end

    subgraph "Archive (Forever)"
        IPFS --> |"Pin"| Filecoin
    end
```
