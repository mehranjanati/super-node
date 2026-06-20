> **⚠️ توجه (آپدیت استراتژی):** این سند ممکن است شامل جزئیات قدیمی مربوط به نسخه‌های اولیه فرانت‌اند باشد. سیستم اکنون از رویکرد **Template-Driven** و **GitOps کاربر‌محور** استفاده می‌کند. برای مشاهده استراتژی نهایی حتماً به [MASTER_PLAN.md](../../MASTER_PLAN.md) و [MVP_DEVELOPMENT_PLAN.md](../MVP_DEVELOPMENT_PLAN.md) مراجعه کنید.

# 🔌 Backend Integration Guide - Nexus Portal

## 📡 Backend Configuration

### Endpoints
- **HTTP API:** `http://localhost:3001`
- **WebSocket:** `ws://localhost:3001`
- **Auth Method:** SIWE (Sign-In with Ethereum)

---

## 🔐 Authentication Flow (SIWE)

### Implementation Steps

#### 1. Install Dependencies
```bash
npm install ethers wagmi viem @rainbow-me/rainbowkit
npm install -D @types/node
```

#### 2. Auth Service (`src/lib/services/auth.ts`)

```typescript
import { ethers } from 'ethers';

const API_URL = 'http://localhost:3001';

export interface AuthUser {
  address: string;
  nonce?: string;
  token?: string;
}

class AuthService {
  private user: AuthUser | null = null;

  /**
   * Step 1: Get nonce from backend
   */
  async getNonce(address: string): Promise<string> {
    const response = await fetch(`${API_URL}/auth/nonce`, {
      method: 'GET',
      headers: {
        'Content-Type': 'application/json',
      },
    });
    
    const data = await response.json();
    return data.nonce;
  }

  /**
   * Step 2: Sign message with wallet
   */
  async signMessage(address: string, nonce: string): Promise<string> {
    if (!window.ethereum) {
      throw new Error('No Ethereum wallet found');
    }

    const provider = new ethers.BrowserProvider(window.ethereum);
    const signer = await provider.getSigner();
    
    const message = `Sign in to Nexus Portal

Address: ${address}
Nonce: ${nonce}`;
    const signature = await signer.signMessage(message);
    
    return signature;
  }

  /**
   * Step 3: Verify signature with backend
   */
  async verify(address: string, signature: string, nonce: string): Promise<AuthUser> {
    const response = await fetch(`${API_URL}/auth/verify`, {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json',
      },
      body: JSON.stringify({ address, signature, nonce }),
    });

    if (!response.ok) {
      throw new Error('Authentication failed');
    }

    const data = await response.json();
    this.user = {
      address,
      token: data.token,
    };
    
    // Store in localStorage
    localStorage.setItem('nexus_auth', JSON.stringify(this.user));
    
    return this.user;
  }

  /**
   * Complete SIWE flow
   */
  async signIn(address: string): Promise<AuthUser> {
    // Step 1: Get nonce
    const nonce = await this.getNonce(address);
    
    // Step 2: Sign message
    const signature = await this.signMessage(address, nonce);
    
    // Step 3: Verify
    const user = await this.verify(address, signature, nonce);
    
    return user;
  }

  /**
   * Get current user
   */
  getUser(): AuthUser | null {
    if (this.user) return this.user;
    
    const stored = localStorage.getItem('nexus_auth');
    if (stored) {
      this.user = JSON.parse(stored);
    }
    
    return this.user;
  }

  /**
   * Sign out
   */
  signOut(): void {
    this.user = null;
    localStorage.removeItem('nexus_auth');
  }

  /**
   * Get auth token for API requests
   */
  getToken(): string | null {
    const user = this.getUser();
    return user?.token || null;
  }
}

export const authService = new AuthService();
```

#### 3. Auth Store (`src/lib/stores/auth.ts`)

```typescript
import { writable } from 'svelte/store';
import { authService, type AuthUser } from '$lib/services/auth';

interface AuthState {
  user: AuthUser | null;
  isAuthenticated: boolean;
  isLoading: boolean;
  error: string | null;
}

function createAuthStore() {
  const { subscribe, set, update } = writable<AuthState>({
    user: authService.getUser(),
    isAuthenticated: !!authService.getUser(),
    isLoading: false,
    error: null,
  });

  return {
    subscribe,

    async signIn(address: string) {
      update((state) => ({ ...state, isLoading: true, error: null }));

      try {
        const user = await authService.signIn(address);
        set({
          user,
          isAuthenticated: true,
          isLoading: false,
          error: null,
        });
      } catch (error) {
        update((state) => ({
          ...state,
          error: error instanceof Error ? error.message : 'Sign in failed',
          isLoading: false,
        }));
        throw error;
      }
    },

    signOut() {
      authService.signOut();
      set({
        user: null,
        isAuthenticated: false,
        isLoading: false,
        error: null,
      });
    },

    clearError() {
      update((state) => ({ ...state, error: null }));
    },
  };
}

export const authStore = createAuthStore();
```

---

## 📊 Crypto Trading Module (REAL)

### WebSocket Integration

#### 1. Market WebSocket Service (`src/lib/services/market-ws.ts`)

```typescript
import { writable } from 'svelte/store';

const WS_URL = 'ws://localhost:3001/ws/market';

export interface MarketAlert {
  timestamp: string;
  symbol: string;
  price: number;
  strategy: 'BUY' | 'SELL' | 'HOLD';
  confidence: number;
  reason?: string;
}

class MarketWebSocketService {
  private ws: WebSocket | null = null;
  private reconnectTimeout: number | null = null;
  private reconnectAttempts = 0;
  private maxReconnectAttempts = 5;

  // Store for alerts
  public alerts = writable<MarketAlert[]>([]);
  public connectionStatus = writable<'connected' | 'disconnected' | 'connecting'>('disconnected');

  connect() {
    if (this.ws?.readyState === WebSocket.OPEN) {
      console.log('WebSocket already connected');
      return;
    }

    this.connectionStatus.set('connecting');
    this.ws = new WebSocket(WS_URL);

    this.ws.onopen = () => {
      console.log('✅ Market WebSocket connected');
      this.connectionStatus.set('connected');
      this.reconnectAttempts = 0;
    };

    this.ws.onmessage = (event) => {
      try {
        const alert: MarketAlert = JSON.parse(event.data);
        console.log('📊 Market Alert:', alert);

        // Add to alerts store
        this.alerts.update((alerts) => [alert, ...alerts].slice(0, 100)); // Keep last 100

        // Trigger UI notification if BUY or SELL
        if (alert.strategy === 'BUY' || alert.strategy === 'SELL') {
          this.triggerTradeAlert(alert);
        }
      } catch (error) {
        console.error('Failed to parse market message:', error);
      }
    };

    this.ws.onerror = (error) => {
      console.error('❌ WebSocket error:', error);
    };

    this.ws.onclose = () => {
      console.log('🔌 Market WebSocket disconnected');
      this.connectionStatus.set('disconnected');
      this.attemptReconnect();
    };
  }

  disconnect() {
    if (this.reconnectTimeout) {
      clearTimeout(this.reconnectTimeout);
      this.reconnectTimeout = null;
    }
    
    if (this.ws) {
      this.ws.close();
      this.ws = null;
    }
    
    this.connectionStatus.set('disconnected');
  }

  private attemptReconnect() {
    if (this.reconnectAttempts >= this.maxReconnectAttempts) {
      console.error('Max reconnection attempts reached');
      return;
    }

    this.reconnectAttempts++;
    const delay = Math.min(1000 * Math.pow(2, this.reconnectAttempts), 30000);
    
    console.log(`Reconnecting in ${delay}ms (attempt ${this.reconnectAttempts})`);
    
    this.reconnectTimeout = window.setTimeout(() => {
      this.connect();
    }, delay);
  }

  private triggerTradeAlert(alert: MarketAlert) {
    // Dispatch custom event for UI to handle
    window.dispatchEvent(
      new CustomEvent('trade-alert', { detail: alert })
    );

    // Browser notification (if permitted)
    if ('Notification' in window && Notification.permission === 'granted') {
      new Notification(`${alert.strategy} Signal: ${alert.symbol}`, {
        body: `Price: $${alert.price} | Confidence: ${alert.confidence}%`,
        icon: '/favicon.png',
      });
    }
  }

  /**
   * Send approve trade signal to backend
   * TODO: Backend needs to implement this endpoint
   */
  async approveTrade(alert: MarketAlert): Promise<void> {
    console.log('📤 Sending approve_trade signal (Simulation):', alert);
    
    // TODO: Implement when backend endpoint is ready
    // const response = await fetch('http://localhost:3001/api/trading/approve', {
    //   method: 'POST',
    //   headers: {
    //     'Content-Type': 'application/json',
    //     'Authorization': `Bearer ${authService.getToken()}`,
    //   },
    //   body: JSON.stringify({
    //     symbol: alert.symbol,
    //     strategy: alert.strategy,
    //     timestamp: alert.timestamp,
    //   }),
    // });
    
    // For now, just log
    await new Promise((resolve) => setTimeout(resolve, 500));
    console.log('✅ Signal sent (Simulation)');
  }
}

export const marketWS = new MarketWebSocketService();
```

#### 2. Trading Dashboard Component

```svelte
<!-- src/lib/components/trading/TradingDashboard.svelte -->
<script lang="ts">
  import { onMount, onDestroy } from 'svelte';
  import { marketWS, type MarketAlert } from '$lib/services/market-ws';
  import Button from '$lib/components/ui/Button.svelte';
  import Card from '$lib/components/ui/Card.svelte';

  let alerts: MarketAlert[] = $state([]);
  let connectionStatus: 'connected' | 'disconnected' | 'connecting' = $state('disconnected');
  let showApprovalModal = $state(false);
  let selectedAlert: MarketAlert | null = $state(null);

  // Subscribe to stores
  const unsubAlerts = marketWS.alerts.subscribe((value) => {
    alerts = value;
  });

  const unsubStatus = marketWS.connectionStatus.subscribe((value) => {
    connectionStatus = value;
  });

  // Listen for trade alerts
  function handleTradeAlert(event: CustomEvent<MarketAlert>) {
    selectedAlert = event.detail;
    showApprovalModal = true;
  }

  async function approveTrade() {
    if (!selectedAlert) return;
    
    await marketWS.approveTrade(selectedAlert);
    showApprovalModal = false;
    selectedAlert = null;
  }

  onMount(() => {
    marketWS.connect();
    window.addEventListener('trade-alert', handleTradeAlert as EventListener);
  });

  onDestroy(() => {
    marketWS.disconnect();
    unsubAlerts();
    unsubStatus();
    window.removeEventListener('trade-alert', handleTradeAlert as EventListener);
  });
</script>

<div class="trading-dashboard">
  <!-- Connection Status -->
  <div class="status-bar">
    <span class="status-indicator" class:connected={connectionStatus === 'connected'}>
      {connectionStatus === 'connected' ? '🟢' : '🔴'}
    </span>
    <span>Market Feed: {connectionStatus}</span>
  </div>

  <!-- Recent Alerts -->
  <div class="alerts-grid">
    {#each alerts as alert}
      <Card>
        <div class="alert-card">
          <div class="alert-header">
            <h3>{alert.symbol}</h3>
            <span class="strategy" class:buy={alert.strategy === 'BUY'} class:sell={alert.strategy === 'SELL'}>
              {alert.strategy}
            </span>
          </div>
          <div class="alert-body">
            <p>Price: ${alert.price.toFixed(2)}</p>
            <p>Confidence: {alert.confidence}%</p>
            {#if alert.reason}
              <p class="reason">{alert.reason}</p>
            {/if}
          </div>
          <div class="alert-footer">
            <small>{new Date(alert.timestamp).toLocaleString()}</small>
            {#if alert.strategy !== 'HOLD'}
              <Button variant="primary" size="sm" onclick={() => {
                selectedAlert = alert;
                showApprovalModal = true;
              }}>
                Approve Trade
              </Button>
            {/if}
          </div>
        </div>
      </Card>
    {/each}
  </div>

  <!-- Approval Modal -->
  {#if showApprovalModal && selectedAlert}
    <div class="modal-overlay" onclick={() => showApprovalModal = false}>
      <div class="modal" onclick={(e) => e.stopPropagation()}>
        <h2>Approve Trade</h2>
        <div class="modal-content">
          <p><strong>Symbol:</strong> {selectedAlert.symbol}</p>
          <p><strong>Action:</strong> {selectedAlert.strategy}</p>
          <p><strong>Price:</strong> ${selectedAlert.price.toFixed(2)}</p>
          <p><strong>Confidence:</strong> {selectedAlert.confidence}%</p>
        </div>
        <div class="modal-actions">
          <Button variant="secondary" onclick={() => showApprovalModal = false}>
            Cancel
          </Button>
          <Button variant="primary" onclick={approveTrade}>
            Approve
          </Button>
        </div>
      </div>
    </div>
  {/if}
</div>

<style>
  .trading-dashboard {
    padding: 1rem;
  }

  .status-bar {
    display: flex;
    align-items: center;
    gap: 0.5rem;
    padding: 0.5rem 1rem;
    background: var(--color-surface);
    border-radius: 0.5rem;
    margin-bottom: 1rem;
  }

  .status-indicator {
    font-size: 0.75rem;
  }

  .alerts-grid {
    display: grid;
    grid-template-columns: repeat(auto-fill, minmax(300px, 1fr));
    gap: 1rem;
  }

  .alert-card {
    padding: 1rem;
  }

  .alert-header {
    display: flex;
    justify-content: space-between;
    align-items: center;
    margin-bottom: 0.5rem;
  }

  .strategy {
    padding: 0.25rem 0.75rem;
    border-radius: 1rem;
    font-weight: 600;
    font-size: 0.875rem;
  }

  .strategy.buy {
    background: #10b981;
    color: white;
  }

  .strategy.sell {
    background: #ef4444;
    color: white;
  }

  .alert-footer {
    display: flex;
    justify-content: space-between;
    align-items: center;
    margin-top: 1rem;
    padding-top: 1rem;
    border-top: 1px solid var(--color-border);
  }

  .modal-overlay {
    position: fixed;
    inset: 0;
    background: rgba(0, 0, 0, 0.5);
    display: flex;
    align-items: center;
    justify-content: center;
    z-index: 1000;
  }

  .modal {
    background: white;
    padding: 2rem;
    border-radius: 0.5rem;
    max-width: 500px;
    width: 90%;
  }

  .modal-actions {
    display: flex;
    gap: 1rem;
    margin-top: 1.5rem;
    justify-content: flex-end;
  }
</style>
```

---

## 🏗️ App Builder Module (SIMULATED)

### Implementation

#### 1. Builder Service (`src/lib/services/builder.ts`)

```typescript
const API_URL = 'http://localhost:3001';

export interface BuilderRequest {
  prompt: string;
  framework?: string;
}

export interface FileNode {
  name: string;
  type: 'file' | 'directory';
  content?: string;
  children?: FileNode[];
}

export interface BuilderResponse {
  fileStructure: FileNode[];
  generatedAt: string;
}

class BuilderService {
  async generate(request: BuilderRequest): Promise<BuilderResponse> {
    const response = await fetch(`${API_URL}/builder/generate`, {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json',
      },
      body: JSON.stringify(request),
    });

    if (!response.ok) {
      throw new Error('Failed to generate app');
    }

    const data = await response.json();
    return data;
  }
}

export const builderService = new BuilderService();
```

#### 2. Builder Page Component

```svelte
<!-- src/routes/(app)/builder/+page.svelte -->
<script lang="ts">
  import { builderService, type FileNode } from '$lib/services/builder';
  import Button from '$lib/components/ui/Button.svelte';

  let prompt = $state('');
  let isLoading = $state(false);
  let fileStructure: FileNode[] = $state([]);
  let selectedFile: FileNode | null = $state(null);
  let error = $state<string | null>(null);

  async function handleGenerate() {
    if (!prompt.trim()) return;

    isLoading = true;
    error = null;

    try {
      const result = await builderService.generate({ prompt });
      fileStructure = result.fileStructure;
      
      // Auto-select first file
      if (fileStructure.length > 0) {
        selectedFile = findFirstFile(fileStructure);
      }
    } catch (err) {
      error = err instanceof Error ? err.message : 'Failed to generate';
    } finally {
      isLoading = false;
    }
  }

  function findFirstFile(nodes: FileNode[]): FileNode | null {
    for (const node of nodes) {
      if (node.type === 'file') return node;
      if (node.children) {
        const found = findFirstFile(node.children);
        if (found) return found;
      }
    }
    return null;
  }

  function selectFile(file: FileNode) {
    if (file.type === 'file') {
      selectedFile = file;
    }
  }
</script>

<div class="builder-page">
  <!-- Left: Prompt Input -->
  <div class="prompt-panel">
    <h2>App Builder</h2>
    <textarea
      bind:value={prompt}
      placeholder="Describe the app you want to build..."
      rows="10"
    ></textarea>
    
    <Button
      variant="primary"
      loading={isLoading}
      onclick={handleGenerate}
      disabled={!prompt.trim() || isLoading}
    >
      {isLoading ? 'Generating...' : 'Generate App'}
    </Button>

    {#if error}
      <div class="error">{error}</div>
    {/if}

    {#if isLoading}
      <div class="loading-state">
        <div class="spinner"></div>
        <p>Generating your app structure...</p>
        <p class="hint">This takes about 2 seconds</p>
      </div>
    {/if}
  </div>

  <!-- Right: Code Preview -->
  <div class="preview-panel">
    {#if fileStructure.length > 0}
      <div class="file-explorer">
        <h3>File Structure</h3>
        <FileTree nodes={fileStructure} onSelect={selectFile} />
      </div>

      <div class="code-editor">
        {#if selectedFile}
          <div class="editor-header">
            <span>{selectedFile.name}</span>
          </div>
          <pre><code>{selectedFile.content || '// No content'}</code></pre>
        {:else}
          <div class="empty-state">
            Select a file to view its content
          </div>
        {/if}
      </div>
    {:else}
      <div class="empty-state">
        <p>Enter a prompt and click "Generate App" to see the file structure</p>
      </div>
    {/if}
  </div>
</div>

<style>
  .builder-page {
    display: grid;
    grid-template-columns: 1fr 1fr;
    gap: 1rem;
    height: calc(100vh - 4rem);
    padding: 1rem;
  }

  .prompt-panel,
  .preview-panel {
    border: 1px solid var(--color-border);
    border-radius: 0.5rem;
    padding: 1.5rem;
    overflow: auto;
  }

  textarea {
    width: 100%;
    padding: 0.75rem;
    border: 1px solid var(--color-border);
    border-radius: 0.5rem;
    font-family: monospace;
    margin: 1rem 0;
  }

  .loading-state {
    text-align: center;
    padding: 2rem;
  }

  .spinner {
    width: 40px;
    height: 40px;
    border: 4px solid #f3f3f3;
    border-top: 4px solid #3498db;
    border-radius: 50%;
    animation: spin 1s linear infinite;
    margin: 0 auto 1rem;
  }

  @keyframes spin {
    0% { transform: rotate(0deg); }
    100% { transform: rotate(360deg); }
  }

  .code-editor {
    background: #1e1e1e;
    color: #d4d4d4;
    padding: 1rem;
    border-radius: 0.5rem;
    margin-top: 1rem;
  }

  .code-editor pre {
    margin: 0;
    overflow: auto;
  }

  .empty-state {
    display: flex;
    align-items: center;
    justify-content: center;
    height: 100%;
    color: var(--color-text-secondary);
  }
</style>
```

---

## 💬 Chat Module (MOCKED_STREAM)

### Implementation

```typescript
// src/lib/services/chat.ts
const API_URL = 'http://localhost:3001';

export interface ChatMessage {
  role: 'user' | 'assistant';
  content: string;
  timestamp: Date;
}

class ChatService {
  async sendMessage(
    message: string,
    onChunk: (chunk: string) => void
  ): Promise<void> {
    const response = await fetch(`${API_URL}/api/chat`, {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json',
      },
      body: JSON.stringify({ message }),
    });

    if (!response.ok) {
      throw new Error('Chat request failed');
    }

    const reader = response.body?.getReader();
    const decoder = new TextDecoder();

    if (!reader) {
      throw new Error('No response body');
    }

    while (true) {
      const { done, value } = await reader.read();
      if (done) break;

      const chunk = decoder.decode(value, { stream: true });
      onChunk(chunk);
    }
  }
}

export const chatService = new ChatService();
```

---

## 📋 Updated Implementation Priority

### Phase 1: Authentication & WebSocket (Week 1)
1. ✅ Install Web3 dependencies
2. ✅ Implement SIWE auth flow
3. ✅ Create auth store
4. ✅ Build login component
5. ✅ Implement WebSocket service
6. ✅ Create trading dashboard
7. ✅ Test real-time alerts

### Phase 2: Trading Module (Week 2)
1. ✅ Build trade approval modal
2. ✅ Implement approve_trade signal
3. ✅ Add browser notifications
4. ✅ Create trading history view
5. ✅ Add filters and search

### Phase 3: App Builder (Week 3)
1. ✅ Implement builder service
2. ✅ Create split-screen UI
3. ✅ Build file tree component
4. ✅ Add code syntax highlighting
5. ✅ Implement loading states

### Phase 4: Chat (Week 4)
1. ✅ Implement streaming chat
2. ✅ Build chat UI with typing effect
3. ✅ Add message history
4. ✅ Implement markdown rendering

---

## 🧪 Testing Strategy

### Unit Tests
- Auth service (SIWE flow)
- WebSocket service (mock WebSocket)
- Builder service (mock fetch)
- Chat service (streaming)

### Integration Tests
- Auth flow end-to-end
- WebSocket reconnection
- Trade approval flow
- Builder generation flow

### E2E Tests
- Complete user journey
- Real WebSocket connection
- Trade alert → approval
- Builder prompt → preview

---

This integration plan connects your existing TDD infrastructure with the real backend!
