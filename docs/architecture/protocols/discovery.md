# Nexus Discovery & Routing Protocol (NDRP)

This document details the mechanisms for **Content Routing** (finding services) and **Peer Routing** (finding nodes) within the Nexus DePIN.

## 1. The Challenge
In a network of 100,000+ nodes and millions of agents, how do you find "The best Agent for analyzing BioTech stocks" without a central registry?

**Solution**: A modified **Kademlia DHT** with **Semantic Tagging**.

---

## 2. Distributed Hash Table (DHT) Structure

Nexus uses `libp2p-kad-dht` with custom validator logic.

### Key-Value Schema

| Record Type | Key Structure | Value Content | TTL |
| :--- | :--- | :--- | :--- |
| **Peer Record** | `/nexus/peer/<PeerID>` | `[Multiaddr_1, Multiaddr_2]` | 1 Hour |
| **Agent Record** | `/nexus/agent/<AgentDID>` | `Provider_PeerID` | 24 Hours |
| **Service Record** | `/nexus/service/<Hash(Tag)>` | `[Provider_A, Provider_B, ...]` | 10 Min |

### Service Discovery (The "Yellow Pages")
When a node wants to advertise a service (e.g., "Llama3 Inference"):
1.  It hashes the tag: `Key = SHA256("service:inference:llama3")`.
2.  It finds the `k` closest peers to that Key (XOR distance).
3.  It stores its PeerID on those peers.

**Query Flow:**
```mermaid
graph LR
    User[User Client] -->|FindProviders("service:inference:llama3")| NodeA[Super Node A]
    NodeA -->|Recursive Lookup| NodeB[Super Node B]
    NodeB -->|Found Record| NodeC[Super Node C (Holder)]
    NodeC -- Return Provider List --> User
```

---

## 3. Semantic Routing (Smart Discovery)

Standard DHTs only support exact matches. Nexus adds a **GossipSub-based Semantic Layer**.

### The Problem
A user asks: "I need an agent that knows about Quantum Physics."
There is no exact tag `service:agent:quantum-physics`.

### The Solution: Vector Gossip
1.  **Embedding**: The query is converted to a vector embedding.
2.  **Broadcast**: The vector is gossiped to a subset of "Index Nodes".
3.  **Similarity Search**: Index Nodes check their local TiDB Vector store for agents with matching descriptions.
4.  **Response**: Nodes return a list of candidate Agent DIDs.

---

## 4. Relay V2 & NAT Traversal

Most Edge Nodes (users running agents on laptops) are behind NATs.

### Circuit Relay
1.  **Reservation**: Edge Node connects to a Super Node and reserves a relay slot.
2.  **Advertisement**: Edge Node advertises a relay address: `/ip4/1.2.3.4/tcp/4001/p2p/SuperNodeID/p2p-circuit/p2p/EdgeNodeID`.
3.  **Connection**: When another peer connects to this address, the Super Node tunnels the traffic.

### DCUtR (Direct Connection Upgrade through Relay)
After the initial relayed handshake, the two peers attempt **Hole Punching** to establish a direct P2P connection, bypassing the Super Node to save bandwidth.

---

## 5. Geo-Routing & Latency Optimization

For latency-sensitive tasks (e.g., Real-time Voice AI), routing matters.
-   **Proximity**: The DHT prefers peers with lower RTT (Round Trip Time).
-   **Region Tags**: Records include region metadata (`us-east`, `eu-central`).
-   **Logic**: `GetClosestPeers(Key)` filters for peers in the same region first.
