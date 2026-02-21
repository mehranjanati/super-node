# SIP Integration Guide

The **LiveKit SIP** service acts as a bridge between the traditional telephony world (SIP) and the modern real-time web (WebRTC/LiveKit). This allows phones to call into Nexus rooms and AI Agents to make outbound calls to real phone numbers.

## 1. What is LiveKit SIP?

LiveKit SIP is a specialized participant in a LiveKit room. It:
- **Inbound**: Listens for incoming SIP calls and connects them to a specific room.
- **Outbound**: Allows the LiveKit server to initiate a SIP call to a phone number.

## 2. Configuration (No Built-in UI)

LiveKit does **not** provide a built-in graphical UI for managing SIP Trunks or Dispatch Rules. Configuration is handled in three ways:

### A. Static Configuration (`sip.yaml`)
Used for fixed setups. You define your trunks and rules directly in the YAML file.

```yaml
sip_trunk:
  - name: my-twilio-trunk
    transport: udp
    address: your-provider.pstn.twilio.com
    username: your_username
    password: your_password

dispatch_rule:
  - rule: ".*" # Match all incoming numbers
    room_prefix: "phone_call_"
```

### B. LiveKit CLI
For dynamic management without restarting the service.

```bash
# Create a SIP Trunk
lk sip trunk create --name "my-trunk" --address "..." --username "..." --password "..."

# Create a Dispatch Rule
lk sip dispatch-rule create --trunk-id "ST_..." --rule ".*" --room-prefix "call_"
```

### C. Server SDK (Go/Node.js)
The most powerful way. You can build your own custom UI (e.g., in your Super Node dashboard) that calls the LiveKit API to manage trunks.

```go
// Example in Go
client := livekit.NewSipClient(url, apiKey, apiSecret)
_, err := client.CreateSipTrunk(ctx, &livekit.CreateSipTrunkRequest{
    Trunk: &livekit.SipTrunkInfo{
        Name: "User_Custom_Trunk",
        Address: "sip.provider.com",
        // ...
    },
})
```

## 3. Key Concepts

### SIP Trunk
Think of this as your "phone line" provider (Twilio, Zadarma, Telnyx). It connects the Super Node to the Global Telephony Network (PSTN).

### Inbound Trunk
A trunk that accepts calls from the outside world.

### Outbound Trunk
A trunk used by the Super Node to call a real person's phone.

### Dispatch Rule
A logic layer that decides what happens when a call arrives.
- **Direct to Room**: Call is placed in a specific room.
- **Metadata-based**: Room name is derived from the caller's ID or the number dialed.

## 4. UI Strategy for Nexus

Since there is no default UI, the Nexus Super Node strategy is:
1.  **Management via Hasura**: Store SIP configurations in TiDB via Hasura.
2.  **Orchestrator Sync**: A background worker in the Super Node watches these TiDB records and uses the LiveKit SDK to sync them with the running `livekit-sip` service.
3.  **Custom Portal**: Users manage their "Phone Numbers" and "Trunks" through the Nexus Portal SPA, which talks to Hasura.
