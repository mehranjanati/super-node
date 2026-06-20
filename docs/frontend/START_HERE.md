> **⚠️ توجه (آپدیت استراتژی):** این سند ممکن است شامل جزئیات قدیمی مربوط به نسخه‌های اولیه فرانت‌اند باشد. سیستم اکنون از رویکرد **Template-Driven** و **GitOps کاربر‌محور** استفاده می‌کند. برای مشاهده استراتژی نهایی حتماً به [MASTER_PLAN.md](../../MASTER_PLAN.md) و [MVP_DEVELOPMENT_PLAN.md](../MVP_DEVELOPMENT_PLAN.md) مراجعه کنید.

# 🎯 Nexus Portal - Complete Development Package

## 📦 What You Have Now

### 🔷 **Backend Integration Ready**
Your Nexus Portal now has **complete integration specifications** for your real backend at `localhost:3001`:

1. **✅ SIWE Authentication** - Sign-In with Ethereum wallet
2. **✅ WebSocket Trading** - Real-time market alerts (BUY/SELL/HOLD)
3. **✅ App Builder** - AI-powered app generation (2s simulated delay)
4. **✅ Streaming Chat** - Chunked transfer with typing effect

---

## 📚 Documentation (8 files)

### Core Documentation
1. **README.md** - Project overview
2. **DEVELOPMENT_PLAN.md** - Original 8-week TDD plan
3. **QUICK_START.md** - TDD workflow guide
4. **CHECKLIST.md** - Task tracker

### Backend Integration (NEW!)
5. **BACKEND_INTEGRATION.md** ⭐ - Complete integration guide with code
6. **IMPLEMENTATION_ROADMAP.md** ⭐ - Revised 6-week backend-first plan
7. **NEXT_STEPS.md** - What to do right now
8. **design_specs.md** - Design specifications

---

## 🏗️ Code & Configuration (13 files)

### Testing Infrastructure
- `vitest.config.ts` - Unit test configuration
- `playwright.config.ts` - E2E test configuration
- `tests/setup.ts` - Global test setup
- `tests/unit/components/Button.test.ts` - Example unit test
- `tests/unit/stores/agents.test.ts` - Store tests
- `tests/e2e/dashboard.spec.ts` - E2E tests

### Source Code
- `src/lib/types/index.ts` - TypeScript definitions
- `src/lib/stores/agents.ts` - Agents store
- `src/lib/components/ui/Button.svelte` - Example component

### CI/CD
- `.github/workflows/test.yml` - Test automation
- `.github/workflows/deploy.yml` - Deployment automation

### Workflow
- `.agent/workflows/setup-testing.md` - Automated setup

### Configuration
- `package.json` - Updated with test scripts

---

## 🎨 Visual Assets (3 diagrams)

1. **Architecture Diagram** - Tech stack overview
2. **Development Roadmap** - 8-week timeline
3. **Backend Integration Flow** ⭐ - NEW! Shows SIWE, WebSocket, APIs

---

## 🔥 Two Development Paths

You now have **TWO complete implementation plans**:

### Path A: Original TDD Plan (8 weeks)
**File:** `DEVELOPMENT_PLAN.md`
- Week 1: Testing + Design System
- Week 2-3: Social Dashboard
- Week 4: Agent Management
- Week 5: Communication (Matrix/LiveKit)
- Week 6: CMS Integration
- Week 7: Polish
- Week 8: Deploy

**Best for:** Learning TDD, building a social network

---

### Path B: Backend Integration Plan (6 weeks) ⭐ RECOMMENDED
**File:** `IMPLEMENTATION_ROADMAP.md`
- Week 1: Auth (SIWE) + WebSocket
- Week 2: Trading Module (Real alerts)
- Week 3: App Builder (Simulated)
- Week 4: Chat (Streaming)
- Week 5: Dashboard + Integration
- Week 6: Polish + Deploy

**Best for:** Integrating with your existing backend

---

## 🚀 Your Next Steps (Choose One Path)

### Option 1: Backend Integration (Recommended)

#### Today (2-3 hours):
```bash
# 1. Install all dependencies
npm install -D vitest @vitest/ui @testing-library/svelte @testing-library/jest-dom @testing-library/user-event jsdom @playwright/test @vitest/coverage-v8

# 2. Install Web3 dependencies
npm install ethers wagmi viem @rainbow-me/rainbowkit

# 3. Start backend (in separate terminal)
cd ../backend  # Your backend directory
npm run dev    # Should start on http://localhost:3001

# 4. Start frontend
npm run dev    # Should start on http://localhost:5173

# 5. Verify tests work
npm test
```

#### This Week:
1. Read `BACKEND_INTEGRATION.md` (complete code examples)
2. Implement SIWE authentication
3. Connect to WebSocket
4. Build trading dashboard
5. Test real-time alerts

#### Follow:
- `IMPLEMENTATION_ROADMAP.md` for weekly plan
- `BACKEND_INTEGRATION.md` for code examples

---

### Option 2: Original TDD Plan

#### Today (2-3 hours):
```bash
# 1. Run automated setup
/setup-testing

# 2. Verify tests work
npm test
npm run test:e2e

# 3. Read QUICK_START.md
# 4. Build your first component (Card)
```

#### This Week:
1. Complete testing infrastructure
2. Build design system
3. Create 5 core UI components
4. Setup mock data

#### Follow:
- `DEVELOPMENT_PLAN.md` for weekly plan
- `QUICK_START.md` for TDD workflow
- `CHECKLIST.md` for task tracking

---

## 📖 Quick Reference Guide

| I want to... | Read this file... |
|--------------|-------------------|
| **Integrate with backend** | `BACKEND_INTEGRATION.md` ⭐ |
| **See backend roadmap** | `IMPLEMENTATION_ROADMAP.md` ⭐ |
| **Learn TDD workflow** | `QUICK_START.md` |
| **See full TDD plan** | `DEVELOPMENT_PLAN.md` |
| **Track my progress** | `CHECKLIST.md` |
| **Get started now** | `NEXT_STEPS.md` |
| **Understand architecture** | Backend Integration Flow diagram |

---

## 🔌 Backend Integration Highlights

### SIWE Authentication Flow
```typescript
// 1. Get nonce
const nonce = await authService.getNonce(address);

// 2. Sign with wallet
const signature = await authService.signMessage(address, nonce);

// 3. Verify and get token
const user = await authService.verify(address, signature, nonce);
```

### WebSocket Trading Alerts
```typescript
// Connect to market feed
marketWS.connect();

// Listen for alerts
marketWS.alerts.subscribe((alerts) => {
  // Display in UI
});

// Approve trade
await marketWS.approveTrade(alert);
```

### App Builder
```typescript
// Generate app
const result = await builderService.generate({
  prompt: 'Create a todo app'
});

// Display file structure
fileStructure = result.fileStructure;
```

### Streaming Chat
```typescript
// Send message with streaming
await chatService.sendMessage(message, (chunk) => {
  // Display chunk with typing effect
});
```

---

## 🎯 Module Status

| Module | Status | Integration | Priority |
|--------|--------|-------------|----------|
| **Crypto Trading** | ✅ REAL | WebSocket `/ws/market` | HIGH |
| **App Builder** | ⚙️ SIMULATED | POST `/builder/generate` | MEDIUM |
| **Chat** | 💬 MOCKED_STREAM | POST `/api/chat` | MEDIUM |
| **Auth** | 🔐 SIWE | GET/POST `/auth/*` | CRITICAL |

---

## 📊 Success Metrics

### Backend Integration
- ✅ SIWE auth works with real backend
- ✅ WebSocket receives live market alerts
- ✅ Trade approval sends signal to Temporal
- ✅ Builder generates file structure (2s delay)
- ✅ Chat streams responses with typing effect

### Testing
- ✅ 80%+ test coverage
- ✅ All E2E tests pass
- ✅ Mock backend for tests

### Performance
- ✅ Lighthouse score > 90
- ✅ WebSocket auto-reconnects
- ✅ Smooth streaming animations

---

## 🛠️ Tech Stack

### Frontend
- **Framework:** SvelteKit 2.x + TypeScript
- **Styling:** TailwindCSS 4.x
- **State:** Svelte Stores
- **Testing:** Vitest + Playwright

### Web3
- **Auth:** SIWE (ethers.js)
- **Wallet:** RainbowKit + wagmi

### Backend Integration
- **HTTP:** Fetch API
- **WebSocket:** Native WebSocket API
- **Streaming:** Chunked Transfer Encoding

---

## 🎨 UI Requirements (from backend spec)

### Dashboard
- Show active agents
- Display recent market alerts
- Quick actions for each module

### Trading Page
- Real-time alert cards
- Connection status indicator
- Approve trade modal
- Filter by strategy (BUY/SELL/HOLD)

### Builder Page
- Split screen layout
- Left: Prompt input
- Right: File tree + Code preview
- Loading state (2s spinner)

### Chat Page
- Message history
- Streaming responses
- Typing effect animation
- Markdown rendering

### Settings Page
- User config form
- Stored in localStorage
- Theme toggle
- Notification preferences

---

## 🚨 Critical Implementation Notes

### WebSocket Flow
```
1. Connect to ws://localhost:3001/ws/market
2. Listen for JSON messages
3. If message.strategy is 'BUY' or 'SELL':
   - Trigger UI alert/modal
   - Show "Approve Trade" button
4. On approve:
   - Send signal to Temporal (simulation for now)
   - Log "Signal sent (Simulation)"
```

### Auth Flow
```
1. GET /auth/nonce → receive nonce
2. Sign with Wallet → generate signature
3. POST /auth/verify → receive JWT token
4. Store token in localStorage
5. Include in all API requests
```

### Missing Backend Features
- ⚠️ Endpoint to send `approve_trade` signal to Temporal
- 📝 UI should log "Signal sent (Simulation)" for now
- 🔜 Will be implemented in backend later

---

## 🎉 You're Ready!

You now have:
- ✅ Complete TDD infrastructure
- ✅ Backend integration specifications
- ✅ Two implementation paths
- ✅ Code examples for all modules
- ✅ Testing strategy
- ✅ CI/CD pipelines
- ✅ Comprehensive documentation

### Choose Your Path:

**Path A (Backend Integration):**
```bash
# Start here
cat BACKEND_INTEGRATION.md
cat IMPLEMENTATION_ROADMAP.md
```

**Path B (TDD Learning):**
```bash
# Start here
cat QUICK_START.md
cat DEVELOPMENT_PLAN.md
```

---

## 📞 Need Help?

| Question | Answer |
|----------|--------|
| How do I integrate with backend? | `BACKEND_INTEGRATION.md` |
| What's the backend roadmap? | `IMPLEMENTATION_ROADMAP.md` |
| How does TDD work? | `QUICK_START.md` |
| What should I build first? | Week 1 in either roadmap |
| How do I track progress? | `CHECKLIST.md` |

---

**Let's build the Nexus Portal! 🚀**

*Choose your path and start coding!*
