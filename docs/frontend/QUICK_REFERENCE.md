> **⚠️ توجه (آپدیت استراتژی):** این سند ممکن است شامل جزئیات قدیمی مربوط به نسخه‌های اولیه فرانت‌اند باشد. سیستم اکنون از رویکرد **Template-Driven** و **GitOps کاربر‌محور** استفاده می‌کند. برای مشاهده استراتژی نهایی حتماً به [MASTER_PLAN.md](../../MASTER_PLAN.md) و [MVP_DEVELOPMENT_PLAN.md](../MVP_DEVELOPMENT_PLAN.md) مراجعه کنید.

# 📋 Nexus Portal - Quick Reference Guide

## 🗺️ Page Structure Overview

| # | Page | Route | Auth | Purpose | Key Features |
|---|------|-------|------|---------|--------------|
| 1 | **Login** | `/` | No | Authentication | SIWE wallet connection |
| 2 | **Dashboard** | `/dashboard` | Yes | Overview | Agents, Alerts, Quick Actions |
| 3 | **Trading** | `/trading` | Yes | Market alerts | Real-time WebSocket, Trade approval |
| 4 | **Builder** | `/builder` | Yes | App generation | AI-powered, File tree, Code preview |
| 5 | **Chat** | `/chat` | Yes | AI assistant | Streaming responses, Markdown |
| 6 | **Agents** | `/agents` | Yes | Agent management | CRUD operations, Performance metrics |
| 7 | **Settings** | `/settings` | Yes | Configuration | Profile, Notifications, Preferences |

---

## 🎨 Component Inventory

### Navigation (1 component)
- **NavBar** - Persistent top navigation with logo, links, user avatar

### Buttons (5 variants)
- **Primary** - Main actions (blue)
- **Secondary** - Alternative actions (gray)
- **Ghost** - Subtle actions (transparent)
- **Danger** - Destructive actions (red)
- **Icon** - Icon-only buttons

### Cards (6 types)
- **Agent Card** - Shows agent info, status, actions
- **Alert Card** - Market alert with approve button
- **Widget Card** - Dashboard widgets
- **File Card** - File tree nodes
- **Message Card** - Chat messages
- **Empty State Card** - No content placeholder

### Inputs (8 types)
- **Text Input** - Single-line text
- **Textarea** - Multi-line text
- **Dropdown** - Select from options
- **Toggle** - On/off switch
- **Slider** - Range selection
- **File Upload** - Image/file picker
- **Search** - Search with icon
- **Date Picker** - Date range selection

### Modals (3 types)
- **Trade Approval Modal** - Confirm trade
- **Agent Details Modal** - View agent info
- **Create Agent Modal** - 3-step wizard

### Status Indicators (4 types)
- **Connection Status** - WebSocket status (🟢/🔴)
- **Agent Status** - Active/Paused/Error (🟢/⏸️/🔴)
- **Loading Spinner** - Processing indicator
- **Typing Indicator** - Chat typing (● ● ●)

### Specialized Components (10+)
- **File Tree** - Recursive folder/file display
- **Code Editor** - Syntax-highlighted code
- **Chat Message** - User/Assistant bubbles
- **Alert Grid** - Responsive alert cards
- **Performance Chart** - ROI/metrics visualization
- **Progress Bar** - Multi-step progress
- **Toast Notification** - Success/error messages
- **Breadcrumbs** - Navigation path
- **Tabs** - Tabbed content sections
- **Sidebar** - Collapsible side panel

---

## 🔄 User Flows

### Flow 1: Sign In
```
Landing (/) 
  → Click "Connect Wallet"
  → Wallet popup appears
  → User approves
  → Sign message
  → Backend verifies
  → Redirect to /dashboard
```

### Flow 2: Approve Trade
```
Dashboard or Trading page
  → WebSocket receives BUY/SELL alert
  → Alert card appears
  → User clicks "Approve Trade"
  → Modal opens with details
  → User clicks "Confirm"
  → Signal sent to backend
  → Success toast appears
  → Card updates to "Approved"
```

### Flow 3: Generate App
```
Builder page (/builder)
  → User enters prompt
  → User selects framework
  → User clicks "Generate App"
  → Loading state (2s)
  → File tree appears
  → First file auto-selected
  → Code preview shows content
  → User can copy/download
```

### Flow 4: Chat with AI
```
Chat page (/chat)
  → User types message
  → User clicks Send (or Enter)
  → User message appears
  → Typing indicator shows
  → Response streams in
  → Text appears character-by-character
  → Typing indicator disappears
  → Message complete
```

### Flow 5: Create Agent
```
Agents page (/agents)
  → Click "+ Create Agent"
  → Modal opens (Step 1/3)
  → Select agent type
  → Click "Next"
  → Step 2: Configure (name, strategy, risk)
  → Click "Next"
  → Step 3: Review
  → Click "Deploy Agent"
  → Loading state
  → Success → Modal closes
  → New agent card appears
```

---

## 🎯 Interaction Patterns

### Click Interactions
| Element | Action | Feedback |
|---------|--------|----------|
| Button | Execute action | Loading state or navigation |
| Card | Expand or navigate | Highlight or modal |
| Link | Navigate | Page transition |
| Toggle | Switch state | Animated flip |
| Dropdown | Show options | Menu appears |

### Hover Interactions
| Element | Effect |
|---------|--------|
| Button | Darker shade |
| Card | Subtle lift (shadow) |
| Link | Underline |
| Icon | Tooltip appears |

### Keyboard Interactions
| Key | Action |
|-----|--------|
| Tab | Navigate between elements |
| Enter | Submit form / Activate button |
| Escape | Close modal |
| Shift+Enter | New line in textarea |
| Ctrl/Cmd+C | Copy (in code blocks) |

---

## 📱 Responsive Breakpoints

| Breakpoint | Width | Layout Changes |
|------------|-------|----------------|
| **Mobile** | < 768px | Single column, stacked widgets, hamburger menu |
| **Tablet** | 768px - 1024px | 2-column grid, some widgets side-by-side |
| **Desktop** | > 1024px | 3-column grid, full sidebar, split screens |

### Responsive Behaviors

**Dashboard:**
- Desktop: 3 columns (Agents | Alerts | Actions)
- Tablet: 2 columns (Agents+Alerts | Actions)
- Mobile: 1 column (stacked)

**Trading:**
- Desktop: Grid + Sidebar
- Tablet: Grid only (sidebar hidden)
- Mobile: List view

**Builder:**
- Desktop: 50/50 split
- Tablet: 40/60 split
- Mobile: Stacked (prompt top, preview bottom)

**Chat:**
- All devices: Full-width messages
- Mobile: Smaller font, narrower bubbles

---

## 🎨 Design Tokens

### Colors
```css
/* Primary */
--color-primary: #3B82F6;      /* Blue */
--color-secondary: #6B7280;    /* Gray */
--color-success: #10B981;      /* Green */
--color-warning: #F59E0B;      /* Yellow */
--color-error: #EF4444;        /* Red */

/* Status */
--color-active: #10B981;       /* Green */
--color-paused: #F59E0B;       /* Yellow */
--color-error: #EF4444;        /* Red */

/* Background */
--color-bg-primary: #FFFFFF;   /* White */
--color-bg-secondary: #F9FAFB; /* Light gray */
--color-bg-dark: #1F2937;      /* Dark gray */

/* Text */
--color-text-primary: #111827;   /* Dark */
--color-text-secondary: #6B7280; /* Gray */
--color-text-inverse: #FFFFFF;   /* White */
```

### Typography
```css
/* Font Families */
--font-display: 'Inter', sans-serif;
--font-body: 'Inter', sans-serif;
--font-mono: 'Fira Code', monospace;

/* Font Sizes */
--text-xs: 0.75rem;    /* 12px */
--text-sm: 0.875rem;   /* 14px */
--text-base: 1rem;     /* 16px */
--text-lg: 1.125rem;   /* 18px */
--text-xl: 1.25rem;    /* 20px */
--text-2xl: 1.5rem;    /* 24px */
--text-3xl: 1.875rem;  /* 30px */
```

### Spacing
```css
/* 8px base grid */
--space-1: 0.25rem;  /* 4px */
--space-2: 0.5rem;   /* 8px */
--space-3: 0.75rem;  /* 12px */
--space-4: 1rem;     /* 16px */
--space-6: 1.5rem;   /* 24px */
--space-8: 2rem;     /* 32px */
--space-12: 3rem;    /* 48px */
```

### Shadows
```css
--shadow-sm: 0 1px 2px rgba(0, 0, 0, 0.05);
--shadow-md: 0 4px 6px rgba(0, 0, 0, 0.1);
--shadow-lg: 0 10px 15px rgba(0, 0, 0, 0.1);
--shadow-xl: 0 20px 25px rgba(0, 0, 0, 0.15);
```

### Border Radius
```css
--radius-sm: 0.25rem;  /* 4px */
--radius-md: 0.5rem;   /* 8px */
--radius-lg: 1rem;     /* 16px */
--radius-full: 9999px; /* Pill shape */
```

---

## ⚡ Performance Targets

| Metric | Target | Measurement |
|--------|--------|-------------|
| **First Contentful Paint** | < 1.5s | Lighthouse |
| **Time to Interactive** | < 3s | Lighthouse |
| **Lighthouse Score** | > 90 | Lighthouse |
| **Bundle Size** | < 500KB | Gzipped |
| **WebSocket Reconnect** | < 2s | Manual test |
| **API Response** | < 500ms | Network tab |

---

## ♿ Accessibility Checklist

### WCAG 2.1 AA Compliance

- [ ] **Color Contrast:** 4.5:1 for normal text, 3:1 for large text
- [ ] **Keyboard Navigation:** All interactive elements accessible via Tab
- [ ] **Focus Indicators:** Visible focus outline on all elements
- [ ] **ARIA Labels:** All buttons, inputs, and regions labeled
- [ ] **Screen Reader:** Tested with NVDA/JAWS
- [ ] **Alt Text:** All images have descriptive alt text
- [ ] **Heading Hierarchy:** Proper H1 → H2 → H3 structure
- [ ] **Form Labels:** All inputs have associated labels
- [ ] **Error Messages:** Clear, descriptive error messages
- [ ] **Skip Links:** "Skip to main content" link

---

## 🧪 Testing Checklist

### Unit Tests
- [ ] All components render correctly
- [ ] Props work as expected
- [ ] Events trigger correctly
- [ ] Edge cases handled

### Integration Tests
- [ ] Auth flow works end-to-end
- [ ] WebSocket connects and receives messages
- [ ] Forms submit correctly
- [ ] Modals open and close

### E2E Tests
- [ ] User can sign in
- [ ] User can approve trade
- [ ] User can generate app
- [ ] User can chat with AI
- [ ] User can create agent
- [ ] User can update settings

### Cross-Browser Tests
- [ ] Chrome
- [ ] Firefox
- [ ] Safari
- [ ] Edge

### Device Tests
- [ ] Desktop (1920x1080)
- [ ] Laptop (1366x768)
- [ ] Tablet (768x1024)
- [ ] Mobile (375x667)

---

## 📚 Documentation Index

| Document | Purpose | When to Use |
|----------|---------|-------------|
| **START_HERE.md** | Overview of everything | First time setup |
| **UX_FLOW_SPECIFICATION.md** | Complete page specs | Building UI |
| **BACKEND_INTEGRATION.md** | API integration guide | Connecting to backend |
| **IMPLEMENTATION_ROADMAP.md** | 6-week plan | Project planning |
| **DEVELOPMENT_PLAN.md** | 8-week TDD plan | Learning TDD |
| **QUICK_START.md** | TDD workflow | Daily development |
| **CHECKLIST.md** | Task tracker | Progress tracking |

---

## 🎯 Quick Stats

- **Total Pages:** 7
- **Total Modals:** 3
- **Total Components:** 45+
- **Total Interactions:** 100+
- **Total User Flows:** 10+
- **Lines of Documentation:** 5000+

---

**Everything you need to build the Nexus Portal is documented!** 🚀

Use `UX_FLOW_SPECIFICATION.md` for detailed page layouts and interactions.
