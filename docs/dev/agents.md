# Building Agents on Nexus

This guide explains how to build, compile, and deploy AI Agents on the Nexus Super Node.

## 1. The Agent Architecture
A Nexus Agent is not just a script; it's a **portable, sovereign logic unit**.
-   **Logic**: Defined visually in **Rivet** (or written in Go/Rust).
-   **Runtime**: Compiled to **WebAssembly (Wasm)** for sandboxed execution.
-   **Storage**: Stored on **IPFS** (Content-Addressed).
-   **Identity**: Registered on the Super Node with a unique DID.

---

## 2. Developing with Rivet (Visual Logic)

[Rivet](https://rivet.ironcladapp.com/) is our primary tool for defining agent behavior. It allows you to create graphs of nodes (Prompt, If/Else, Vector Search).

### Step-by-Step
1.  **Install Rivet**: Download the IDE from the official site.
2.  **Create a Project**: Start a new `.rivet-project`.
3.  **Define Inputs**:
    -   `user_intent`: What the user wants (e.g., "Analyze BTC").
    -   `market_data`: Real-time signals injected by the Super Node.
4.  **Add Logic Nodes**:
    -   Use the **Chat Node** to call LLMs (handled by the Super Node via Unsloth/Cocoon).
    -   Use the **Vector Database Node** to recall past memories.
5.  **Export**: Save the graph.

---

## 3. Compiling to Wasm

To run on the Super Node, your logic must be bundled into a Wasm module.

### Using the Nexus CLI
We provide a CLI tool to wrap your Rivet graph or Go code.

```bash
# Install Nexus CLI
npm install -g @nexus/cli

# Scaffolding a new agent
nexus init my-trading-agent --template rivet

# Build the Wasm binary
cd my-trading-agent
nexus build
# Output: dist/agent.wasm
```

### Writing in Go (Advanced)
If you prefer code over graphs:

```go
package main

import "github.com/nexus/sdk-go/agent"

func main() {
    agent.HandleRequest(func(ctx agent.Context, input []byte) ([]byte, error) {
        // Your logic here
        price := ctx.GetMarketPrice("BTC")
        if price > 90000 {
            return ctx.ExecuteTrade("SELL", 0.1)
        }
        return nil, nil
    })
}
```

---

## 4. Deploying to IPFS

Once built, the agent must be published to the decentralized storage layer.

```bash
nexus publish dist/agent.wasm
```
**Output:**
```
Successfully uploaded to IPFS!
CID: QmXyZ123456789ABCDEF...
```

---

## 5. Registering on the Super Node

Now, tell the Super Node to run this CID when specific conditions are met.

### The Manifest (`nexus.yaml`)
```yaml
name: "btc-sniper-v1"
version: "1.0.0"
cid: "QmXyZ123..."
trigger:
  type: "market_signal"
  condition: "BTC > 90000"
permissions:
  - "trade:execute"
  - "memory:read"
privacy: "confidential" # Uses Cocoon TEE
```

### Registration Command
```bash
nexus register nexus.yaml --node https://supernode.example.com
```

---

## 6. Advanced: Fine-Tuning Your Agent

If your agent requires specialized knowledge or a specific speaking style, you can fine-tune a dedicated LLM model using **Unsloth**.
See the [AI Model Training Guide](../ai/unsloth.md) for details on how to create and deploy custom LoRA adapters for your agents.

---

## 7. Adding Eyes & Ears (LiveKit & Matrix)

Nexus Agents are multimodal. They can join voice calls and chat rooms.

### Voice & Vision (LiveKit)
In your Rivet graph or Go code, you can trigger a **LiveKit Connection**.
```go
// Go Example
ctx.JoinRoom("room-123", agent.WithCamera(), agent.WithMicrophone())
```
The Super Node handles the WebRTC negotiation. The agent receives:
- **Audio Stream**: Automatically transcribed to text.
- **Video Frames**: Passed to Vision LLMs (like GPT-4o or Llava).

### Secure Chat (Matrix)
Agents use Matrix for persistent, encrypted communication.
- **Trigger**: `on_matrix_message`
- **Action**: `send_matrix_message`
This allows agents to exist in Element, Beeper, or any Matrix client.

---

## 9. Testing & Debugging
Before deploying to the mainnet, you should verify your agent locally.

### Local Simulation
Use the `nexus-cli` to mock inputs and run the Wasm.

```bash
# Run agent with mock input
nexus run dist/agent.wasm --input '{"intent": "buy_btc"}' --mock-market '{"BTC": 95000}'
```

### Logs & Tracing
Agents print to `stdout`, which is captured by the Super Node.
*   **In Rivet**: Use the "Log" node to inspect values during graph execution.
*   **In Go**: Use `agent.Log("Checking price...")`.
*   **Viewing Logs**:
    ```bash
    nexus logs -f agent-id-123
    ```

---

## 10. The Nexus SDK (Client-Side)

If you are building a **Telegram Mini App** or **React Portal** that interacts with this agent:

```typescript
import { NexusClient } from '@nexus/sdk-js';

const nexus = new NexusClient({
  nodeUrl: 'wss://supernode.example.com',
  auth: { telegram: window.Telegram.WebApp.initData }
});

// Subscribe to Agent updates
nexus.agent('btc-sniper-v1').on('action', (event) => {
  console.log('Agent executed trade:', event);
});

// Manually trigger the agent
await nexus.agent('btc-sniper-v1').execute({
  intent: "Force re-evaluation"
});
```
