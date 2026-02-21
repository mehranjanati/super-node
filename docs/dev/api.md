# API Reference

This document details the public interfaces of the Nexus Super Node.

## 1. Authentication & Security
*All API endpoints (except Health Check) require authentication.*

### Authentication Flow
1.  **Login/Connect**:
    *   Client signs a message with their Wallet (Ethereum/Solana) or sends a Telegram InitData.
    *   Endpoint: `POST /v1/auth/login`
2.  **Receive JWT**:
    *   Server verifies signature.
    *   Returns a **JWT Access Token** (Short-lived) and **Refresh Token** (Long-lived).
3.  **Authenticated Requests**:
    *   Pass the token in the header: `Authorization: Bearer <token>`

### Roles & Permissions
*   `anonymous`: Read-only access to public market data.
*   `user`: Can manage their own Agents and view their data.
*   `node_operator`: Full system access (Admin RPC).

---

## 2. GraphQL API (Hasura)
*Endpoint: `https://<node-url>/v1/graphql`*

Used for reading state, historical data, and user profiles.

### Queries

#### Get Market Signals
```graphql
query GetMarketSignals {
  market_signals(limit: 10, order_by: {created_at: desc}) {
    symbol
    price
    source
    created_at
  }
}
```

#### Get Agent Status
```graphql
query GetAgentStatus($user_id: String!) {
  agents(where: {owner_id: {_eq: $user_id}}) {
    id
    name
    status # IDLE, RUNNING, TRAINING
    last_execution_at
    logs(limit: 5) {
      message
      level
    }
  }
}
```

### Subscriptions (Real-time)
```graphql
subscription MonitorTrades {
  trades(where: {status: {_eq: "PENDING"}}) {
    id
    symbol
    amount
    side
  }
}
```

---

## 2. WebSocket Gateway
*Endpoint: `wss://<node-url>/ws`*

Used for high-frequency signaling, bi-directional communication, and ephemeral state.

### Connection
Connect with an Auth Token (JWT) in the header or query param.

### Message Protocol (JSON)

#### 1. Subscribe to Topic
```json
{
  "op": "subscribe",
  "topic": "market.btc.usd"
}
```

#### 2. Publish User Intent
```json
{
  "op": "intent",
  "payload": {
    "agent_id": "QmHash...",
    "command": "pause_trading"
  }
}
```

#### 3. Server Event (Push)
```json
{
  "op": "event",
  "topic": "agent.execution",
  "data": {
    "agent_id": "QmHash...",
    "result": "Order Placed",
    "tx_hash": "0x123..."
  }
}
```

---

## 3. Admin RPC (gRPC)
*Port: 50051 (Protected via mTLS)*

Used by the Node Operator and internal services (Temporal/Cocoon).

### Service: `NodeControl`

#### `RegisterAgent(Manifest)`
Registers a new agent CID in the local registry.

#### `PruneCache(Request)`
Clears Redis and temporary IPFS blocks to free up space.

## 4. Real-Time Media & Chat (LiveKit & Matrix)

The Super Node acts as a gateway to these services, handling authentication and room management.

### LiveKit (Audio/Video)
*Endpoint: `wss://<node-url>/livekit`*

#### `GenerateToken(RoomName, ParticipantName)`
Returns a JWT to join a WebRTC room.
- **Input**: `{ room: "daily-standup", identity: "agent-007" }`
- **Output**: `{ token: "eyJhbG..." }`

#### `AgentConnect(RoomID)`
Instructs an AI Agent to join a room as a participant.
- **Capabilities**: The agent can subscribe to audio tracks (transcription) and publish audio (TTS).

### Matrix (Chat/Signaling)
*Endpoint: `https://<node-url>/_matrix/client/v3`*

We run a federated **Synapse** (or **Dendrite**) instance. Agents interact via standard Matrix Client-Server API.
- **Login**: `POST /login`
- **Sync**: `GET /sync` (Long polling for new messages)
- **Send**: `PUT /rooms/{roomId}/send/m.room.message/{txnId}`

> **Note**: All Agent-to-Agent communication should happen over Matrix for auditability and encryption.

