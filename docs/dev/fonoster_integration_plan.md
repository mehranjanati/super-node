# Future Roadmap: Scaling Voice with Fonoster

This document outlines the migration path from the MVP architecture (LiveKit SIP) to a Production-Grade architecture using **Fonoster**.

## 1. Why Migration is Needed?

While **LiveKit SIP** is excellent for low-latency media bridging, it is not a full-featured Programmable Voice Platform. As the Nexus Super Node scales to support thousands of concurrent calls and complex telephony workflows, we will need a dedicated "Telephony Orchestrator".

| Feature | LiveKit SIP (MVP) | Fonoster (Production) |
| :--- | :--- | :--- |
| **Role** | Media Bridge (SIP <-> WebRTC) | Programmable Voice Platform |
| **Scale** | Single Node / Vertical Scaling | Cloud-Native / Horizontal Scaling |
| **Logic** | Simple Dispatch Rules (Regex) | Full Programmable Scripts (JavaScript) |
| **Features** | Inbound/Outbound Calls | IVR, Recording, TTS, STT, Billing |
| **Tenancy** | Single Tenant | Multi-Tenant (Project-based) |

## 2. The Target Architecture

In the Production phase, **Fonoster** sits between the PSTN (Public Telephone Network) and the Nexus Super Node.

```mermaid
graph TD
    subgraph "Public Telephony Network (PSTN)"
        Twilio[Twilio / Zadarma]
        Mobile[User Mobile Phone]
    end

    subgraph "Fonoster Cluster (The Gatekeeper)"
        Proxy[SIP Proxy (Kamailio)]
        AppServer[Applications Server]
        MediaServer[Media Server (Asterisk/RTPProxy)]
    end

    subgraph "Nexus Super Node"
        LKSIP[LiveKit SIP Bridge]
        LK[LiveKit Server]
        Agent[AI Agent (Rivet/Wasm)]
    end

    Mobile --> |"SIP Call"| Twilio
    Twilio --> |"SIP Invite"| Proxy
    Proxy --> |"Route"| AppServer
    AppServer --> |"Execute Logic (IVR/Auth)"| MediaServer
    MediaServer --> |"Forward Approved Call"| LKSIP
    LKSIP --> |"WebRTC"| LK
    LK <--> Agent
```

## 3. Integration Strategy

Since we are designing the UI and Database now, we can ensure a seamless migration later.

### A. Database Compatibility
Our `sip_trunks` table in TiDB is designed to be compatible with Fonoster's `Projects` and `Trunks` schemas.
- **MVP**: The Super Node reads this table and configures LiveKit SIP.
- **Production**: The Super Node reads this table and calls the **Fonoster API** to provision resources.

### B. UI Consistency
The user interface in the Nexus Portal will remain **identical**.
- Users still see "My Phone Numbers" and "Call Rules".
- The backend implementation changes from "Direct LiveKit Config" to "Fonoster API Calls".

## 4. Migration Steps

1.  **Deploy Fonoster Cluster**: Set up Fonoster on Kubernetes alongside the Super Node.
2.  **Point Trunks to Fonoster**: Change the SIP URI in the provider (e.g., Twilio) to point to Fonoster's Ingress instead of LiveKit directly.
3.  **Create "Bridge Application"**: Write a simple Fonoster Script that forwards calls to LiveKit.
    ```javascript
    // Fonoster Script (Example)
    const { VoiceServer } = require("@fonoster/voice");

    const server = new VoiceServer({ base: '/voice' });

    server.use(async (req, res) => {
      // 1. Play Welcome Message
      await res.say("Welcome to Nexus AI. Connecting you to your agent...");

      // 2. Forward to LiveKit SIP URI
      await res.dial("sip:my-room@livekit-sip.supernode.local");
    });
    ```
4.  **Switch Backend Logic**: Update the Super Node to provision Fonoster Projects instead of updating `sip.yaml`.

## 5. Success Metrics
- **Zero Downtime**: The migration should not drop active calls.
- **Enhanced Features**: Users immediately gain access to Call Recording and IVR menus without changing their agent logic.
