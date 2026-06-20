> **⚠️ توجه (آپدیت استراتژی):** این سند ممکن است شامل جزئیات قدیمی مربوط به نسخه‌های اولیه فرانت‌اند باشد. سیستم اکنون از رویکرد **Template-Driven** و **GitOps کاربر‌محور** استفاده می‌کند. برای مشاهده استراتژی نهایی حتماً به [MASTER_PLAN.md](../../MASTER_PLAN.md) و [MVP_DEVELOPMENT_PLAN.md](../MVP_DEVELOPMENT_PLAN.md) مراجعه کنید.

# 🚀 Nexus Portal - Revised Implementation Plan (Backend-First)

**Project:** Nexus Portal (Super Node Interface)  
**Backend:** http://localhost:3001 (REAL)  
**Approach:** TDD + Backend Integration  
**Timeline:** 6 weeks to production

---

## 🎯 Project Overview

### What We're Building
Unified frontend for Nexus Super Node with three core modules:

1. **Crypto Trading** (REAL) - WebSocket alerts + Temporal workflow approvals
2. **App Builder** (SIMULATED) - AI-powered app generation with 2s delay
3. **Chat** (MOCKED_STREAM) - Streaming chat with typing effect

### Backend Integration
- **HTTP API:** `http://localhost:3001`
- **WebSocket:** `ws://localhost:3001/ws/market`
- **Auth:** SIWE (Sign-In with Ethereum)

---

## 📅 Revised Timeline (6 Weeks)

### Week 1: Foundation + Auth + WebSocket
- Testing infrastructure
- SIWE authentication
- WebSocket service
- Basic UI components

### Week 2: Trading Module (REAL)
- Real-time market alerts
- Trade approval flow
- Trading dashboard
- Notification system

### Week 3: App Builder (SIMULATED)
- Builder service integration
- Split-screen UI
- File tree visualization
- Code preview

### Week 4: Chat Module (MOCKED_STREAM)
- Streaming chat service
- Typing effect UI
- Message history
- Markdown rendering

### Week 5: Dashboard + Integration
- Unified dashboard
- Module integration
- Settings page
- User config management

### Week 6: Polish + Deploy
- Performance optimization
- Accessibility
- CI/CD
- Documentation

---

## 🔥 Week 1: Foundation + Auth + WebSocket

### Day 1-2: Testing Infrastructure (8 hours)

**Install Dependencies:**
```bash
# Testing
npm install -D vitest @vitest/ui @testing-library/svelte @testing-library/jest-dom @testing-library/user-event jsdom @playwright/test @vitest/coverage-v8

# Web3 / Auth
npm install ethers wagmi viem @rainbow-me/rainbowkit

# WebSocket
npm install ws
npm install -D @types/ws
```

**Setup:**
- [x] Create `vitest.config.ts`
- [x] Create `playwright.config.ts`
- [x] Create `tests/setup.ts`
- [x] Update `package.json` scripts
- [ ] Verify tests run

---

### Day 3-4: SIWE Authentication (12 hours)

**Implementation:**

1. **Auth Service** (`src/lib/services/auth.ts`)
   - [ ] Write tests for getNonce()
   - [ ] Write tests for signMessage()
   - [ ] Write tests for verify()
   - [ ] Implement SIWE flow
   - [ ] Test with real backend

2. **Auth Store** (`src/lib/stores/auth.ts`)
   - [ ] Write tests for signIn()
   - [ ] Write tests for signOut()
   - [ ] Write tests for getUser()
   - [ ] Implement store
   - [ ] Test localStorage persistence

3. **Login Component** (`src/lib/components/auth/LoginButton.svelte`)
   - [ ] Write component tests
   - [ ] Implement wallet connection
   - [ ] Implement sign-in flow
   - [ ] Add loading states
   - [ ] Add error handling

**E2E Test:**
```typescript
// tests/e2e/auth.spec.ts
test('user can sign in with wallet', async ({ page }) => {
  await page.goto('/');
  await page.click('[data-testid="connect-wallet"]');
  // Mock wallet interaction
  await expect(page.locator('[data-testid="user-address"]')).toBeVisible();
});
```

**Verification:**
- ✅ Can connect wallet
- ✅ Can sign message
- ✅ Backend verifies signature
- ✅ Token stored in localStorage
- ✅ Protected routes work

---

### Day 5-7: WebSocket Service (16 hours)

**Implementation:**

1. **Market WebSocket Service** (`src/lib/services/market-ws.ts`)
   - [ ] Write tests (mock WebSocket)
   - [ ] Implement connection logic
   - [ ] Implement reconnection logic
   - [ ] Handle market alerts
   - [ ] Trigger UI notifications
   - [ ] Test with real backend

2. **Trading Dashboard** (`src/lib/components/trading/TradingDashboard.svelte`)
   - [ ] Write component tests
   - [ ] Display connection status
   - [ ] Show real-time alerts
   - [ ] Implement alert cards
   - [ ] Add time formatting

3. **Trade Approval Modal** (`src/lib/components/trading/TradeApprovalModal.svelte`)
   - [ ] Write component tests
   - [ ] Show trade details
   - [ ] Implement approve button
   - [ ] Call approveTrade() (simulation)
   - [ ] Show success/error states

**E2E Test:**
```typescript
// tests/e2e/trading.spec.ts
test('receives market alert and can approve', async ({ page }) => {
  await page.goto('/trading');
  
  // Wait for WebSocket connection
  await expect(page.locator('text=connected')).toBeVisible();
  
  // Simulate alert (or wait for real one)
  // Click approve button
  await page.click('[data-testid="approve-trade"]');
  
  await expect(page.locator('text=Signal sent')).toBeVisible();
});
```

**Verification:**
- ✅ WebSocket connects to backend
- ✅ Receives real-time alerts
- ✅ Displays BUY/SELL signals
- ✅ Approval modal works
- ✅ Reconnects on disconnect

---

## 🔥 Week 2: Trading Module (REAL)

### Day 1-2: Trading UI Components (12 hours)

**Components to Build:**

1. **AlertCard** (`src/lib/components/trading/AlertCard.svelte`)
   - [ ] Write tests
   - [ ] Display symbol, price, strategy
   - [ ] Show confidence level
   - [ ] Add approve button (if BUY/SELL)
   - [ ] Format timestamp

2. **TradingHistory** (`src/lib/components/trading/TradingHistory.svelte`)
   - [ ] Write tests
   - [ ] Display past alerts
   - [ ] Filter by strategy (BUY/SELL/HOLD)
   - [ ] Search by symbol
   - [ ] Pagination

3. **ConnectionStatus** (`src/lib/components/trading/ConnectionStatus.svelte`)
   - [ ] Write tests
   - [ ] Show connected/disconnected
   - [ ] Show reconnection attempts
   - [ ] Manual reconnect button

---

### Day 3-4: Notification System (10 hours)

**Implementation:**

1. **Browser Notifications**
   - [ ] Request permission
   - [ ] Show notification on BUY/SELL
   - [ ] Click notification → open modal
   - [ ] Test on different browsers

2. **In-App Notifications**
   - [ ] Create notification store
   - [ ] Toast notifications
   - [ ] Notification center
   - [ ] Mark as read

---

### Day 5-7: Trading Dashboard Page (14 hours)

**Page:** `src/routes/(app)/trading/+page.svelte`

**Features:**
- [ ] Real-time alerts grid
- [ ] Connection status bar
- [ ] Filter controls (BUY/SELL/HOLD)
- [ ] Search by symbol
- [ ] Trading history sidebar
- [ ] Statistics (total alerts, success rate)

**E2E Test:**
```typescript
test('trading dashboard shows all features', async ({ page }) => {
  await page.goto('/trading');
  
  await expect(page.locator('[data-testid="connection-status"]')).toBeVisible();
  await expect(page.locator('[data-testid="alerts-grid"]')).toBeVisible();
  await expect(page.locator('[data-testid="filter-controls"]')).toBeVisible();
});
```

---

## 🔥 Week 3: App Builder (SIMULATED)

### Day 1-2: Builder Service (10 hours)

**Implementation:**

1. **Builder Service** (`src/lib/services/builder.ts`)
   - [ ] Write tests (mock fetch)
   - [ ] Implement generate()
   - [ ] Handle 2s delay
   - [ ] Parse file structure
   - [ ] Error handling

2. **Builder Store** (`src/lib/stores/builder.ts`)
   - [ ] Write tests
   - [ ] Store file structure
   - [ ] Track loading state
   - [ ] Store selected file
   - [ ] Error state

---

### Day 3-5: Builder UI (18 hours)

**Components:**

1. **PromptInput** (`src/lib/components/builder/PromptInput.svelte`)
   - [ ] Write tests
   - [ ] Textarea with placeholder
   - [ ] Character count
   - [ ] Generate button
   - [ ] Loading state (2s spinner)

2. **FileTree** (`src/lib/components/builder/FileTree.svelte`)
   - [ ] Write tests
   - [ ] Recursive tree rendering
   - [ ] Expand/collapse folders
   - [ ] File selection
   - [ ] Icons for file types

3. **CodePreview** (`src/lib/components/builder/CodePreview.svelte`)
   - [ ] Write tests
   - [ ] Syntax highlighting (Prism.js or Shiki)
   - [ ] Line numbers
   - [ ] Copy button
   - [ ] Download button

---

### Day 6-7: Builder Page (12 hours)

**Page:** `src/routes/(app)/builder/+page.svelte`

**Layout:**
- Left panel: Prompt input
- Right panel: File tree + Code preview

**Features:**
- [ ] Split-screen layout
- [ ] Responsive (stack on mobile)
- [ ] Loading state with progress
- [ ] Empty state
- [ ] Error state

**E2E Test:**
```typescript
test('can generate app and view code', async ({ page }) => {
  await page.goto('/builder');
  
  await page.fill('[data-testid="prompt-input"]', 'Create a todo app');
  await page.click('[data-testid="generate-button"]');
  
  // Wait for 2s loading
  await expect(page.locator('[data-testid="loading-spinner"]')).toBeVisible();
  
  // File structure appears
  await expect(page.locator('[data-testid="file-tree"]')).toBeVisible({ timeout: 3000 });
  
  // Click a file
  await page.click('[data-testid="file-node"]').first();
  
  // Code preview shows
  await expect(page.locator('[data-testid="code-preview"]')).toBeVisible();
});
```

---

## 🔥 Week 4: Chat Module (MOCKED_STREAM)

### Day 1-2: Chat Service (10 hours)

**Implementation:**

1. **Chat Service** (`src/lib/services/chat.ts`)
   - [ ] Write tests (mock streaming)
   - [ ] Implement sendMessage()
   - [ ] Handle chunked transfer
   - [ ] Parse chunks
   - [ ] Error handling

2. **Chat Store** (`src/lib/stores/chat.ts`)
   - [ ] Write tests
   - [ ] Store messages
   - [ ] Track streaming state
   - [ ] Handle typing effect

---

### Day 3-5: Chat UI (18 hours)

**Components:**

1. **MessageList** (`src/lib/components/chat/MessageList.svelte`)
   - [ ] Write tests
   - [ ] Display messages
   - [ ] User vs Assistant styling
   - [ ] Auto-scroll to bottom
   - [ ] Typing indicator

2. **MessageInput** (`src/lib/components/chat/MessageInput.svelte`)
   - [ ] Write tests
   - [ ] Textarea with auto-resize
   - [ ] Send button
   - [ ] Keyboard shortcuts (Enter to send)
   - [ ] Disabled while streaming

3. **TypingEffect** (`src/lib/components/chat/TypingEffect.svelte`)
   - [ ] Write tests
   - [ ] Render text character by character
   - [ ] Configurable speed
   - [ ] Cursor animation

---

### Day 6-7: Chat Page (12 hours)

**Page:** `src/routes/(app)/chat/+page.svelte`

**Features:**
- [ ] Message history
- [ ] Streaming responses
- [ ] Markdown rendering
- [ ] Code syntax highlighting
- [ ] Copy code blocks
- [ ] Clear history

**E2E Test:**
```typescript
test('can send message and see streaming response', async ({ page }) => {
  await page.goto('/chat');
  
  await page.fill('[data-testid="message-input"]', 'Hello!');
  await page.click('[data-testid="send-button"]');
  
  // User message appears
  await expect(page.locator('text=Hello!')).toBeVisible();
  
  // Typing indicator appears
  await expect(page.locator('[data-testid="typing-indicator"]')).toBeVisible();
  
  // Response streams in
  await expect(page.locator('[data-testid="assistant-message"]')).toBeVisible({ timeout: 5000 });
});
```

---

## 🔥 Week 5: Dashboard + Integration

### Day 1-3: Unified Dashboard (18 hours)

**Page:** `src/routes/(app)/+page.svelte`

**Sections:**

1. **Active Agents** (Top)
   - [ ] Show running agents
   - [ ] Agent status indicators
   - [ ] Quick actions

2. **Recent Market Alerts** (Middle)
   - [ ] Last 5 alerts
   - [ ] Quick approve buttons
   - [ ] Link to full trading page

3. **Quick Actions** (Bottom)
   - [ ] Generate App
   - [ ] Start Chat
   - [ ] View Trading History

---

### Day 4-5: Settings Page (12 hours)

**Page:** `src/routes/(app)/settings/+page.svelte`

**Features:**
- [ ] User profile (address, avatar)
- [ ] Notification preferences
- [ ] Trading settings
- [ ] Builder preferences
- [ ] Theme toggle (dark/light)

**Store in localStorage for now**

---

### Day 6-7: Module Integration (10 hours)

**Tasks:**
- [ ] Ensure all modules work together
- [ ] Shared state management
- [ ] Navigation between modules
- [ ] Consistent styling
- [ ] Error boundary

---

## 🔥 Week 6: Polish + Deploy

### Day 1-2: Performance (12 hours)
- [ ] Code splitting
- [ ] Lazy loading
- [ ] Image optimization
- [ ] Bundle analysis
- [ ] Lighthouse > 90

### Day 3-4: Accessibility (10 hours)
- [ ] Keyboard navigation
- [ ] Screen reader testing
- [ ] ARIA labels
- [ ] Color contrast
- [ ] WCAG AA compliance

### Day 5-6: CI/CD (10 hours)
- [ ] GitHub Actions setup
- [ ] Automated tests
- [ ] Deploy to production
- [ ] Environment variables

### Day 7: Documentation (8 hours)
- [ ] Update README
- [ ] API documentation
- [ ] User guide
- [ ] Deployment guide

---

## 📊 Testing Strategy

### Unit Tests (80%+ coverage)
- Services (auth, market-ws, builder, chat)
- Stores (auth, trading, builder, chat)
- Components (all UI components)

### Integration Tests
- Auth flow
- WebSocket reconnection
- Trade approval flow
- Builder generation
- Chat streaming

### E2E Tests (Critical paths)
- Complete user journey
- Auth → Trading → Approve
- Auth → Builder → Generate
- Auth → Chat → Send message

---

## 🎯 Success Criteria

### Functionality
- ✅ SIWE auth works with backend
- ✅ WebSocket receives real alerts
- ✅ Trade approval sends signal
- ✅ Builder generates file structure
- ✅ Chat streams responses

### Performance
- ✅ Lighthouse score > 90
- ✅ First Contentful Paint < 1.5s
- ✅ WebSocket reconnects automatically

### Quality
- ✅ 80%+ test coverage
- ✅ Zero TypeScript errors
- ✅ WCAG AA compliant
- ✅ Works on mobile

---

## 🚀 Quick Start

### 1. Install Dependencies
```bash
npm install -D vitest @vitest/ui @testing-library/svelte @testing-library/jest-dom @testing-library/user-event jsdom @playwright/test @vitest/coverage-v8
npm install ethers wagmi viem @rainbow-me/rainbowkit ws
npm install -D @types/ws
```

### 2. Start Backend
```bash
# In backend directory
npm run dev  # Should start on http://localhost:3001
```

### 3. Start Frontend
```bash
npm run dev  # Should start on http://localhost:5173
```

### 4. Run Tests
```bash
npm test              # Unit tests
npm run test:e2e      # E2E tests
npm run test:coverage # Coverage report
```

---

## 📚 Key Files

### Services
- `src/lib/services/auth.ts` - SIWE authentication
- `src/lib/services/market-ws.ts` - WebSocket for trading
- `src/lib/services/builder.ts` - App builder API
- `src/lib/services/chat.ts` - Streaming chat API

### Stores
- `src/lib/stores/auth.ts` - Auth state
- `src/lib/stores/trading.ts` - Trading alerts
- `src/lib/stores/builder.ts` - Builder state
- `src/lib/stores/chat.ts` - Chat messages

### Pages
- `src/routes/(app)/+page.svelte` - Dashboard
- `src/routes/(app)/trading/+page.svelte` - Trading
- `src/routes/(app)/builder/+page.svelte` - Builder
- `src/routes/(app)/chat/+page.svelte` - Chat
- `src/routes/(app)/settings/+page.svelte` - Settings

---

**This plan integrates your real backend with TDD best practices!** 🚀
