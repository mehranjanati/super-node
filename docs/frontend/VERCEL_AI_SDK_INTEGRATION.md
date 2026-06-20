> **⚠️ توجه (آپدیت استراتژی):** این سند ممکن است شامل جزئیات قدیمی مربوط به نسخه‌های اولیه فرانت‌اند باشد. سیستم اکنون از رویکرد **Template-Driven** و **GitOps کاربر‌محور** استفاده می‌کند. برای مشاهده استراتژی نهایی حتماً به [MASTER_PLAN.md](../../MASTER_PLAN.md) و [MVP_DEVELOPMENT_PLAN.md](../MVP_DEVELOPMENT_PLAN.md) مراجعه کنید.

# 🤖 Vercel AI SDK Integration Guide

## 📋 Overview

**Purpose:** Integrate Vercel AI SDK for AI-powered app generation  
**Alternative to:** Custom backend `/builder/generate` endpoint  
**Benefits:** 
- Streaming responses
- Better AI models (GPT-4, Claude, etc.)
- Built-in UI components
- Type-safe API

---

## 🚀 Installation

### Step 1: Install Dependencies

```bash
# Core AI SDK
npm install ai

# OpenAI provider (or use Anthropic, Google, etc.)
npm install @ai-sdk/openai

# Optional: UI components for React
npm install @ai-sdk/react

# Optional: Svelte integration
npm install @ai-sdk/svelte
```

### Step 2: Environment Variables

Create `.env.local`:
```bash
# OpenAI API Key
OPENAI_API_KEY=sk-...

# Or use other providers:
# ANTHROPIC_API_KEY=...
# GOOGLE_API_KEY=...

# Backend URL (if using custom backend)
VITE_BACKEND_URL=http://localhost:3001
```

---

## 🏗️ Integration Options

### Option 1: Use Vercel AI SDK (Recommended)
**Pros:**
- Streaming responses
- Better AI models
- Built-in UI components
- Type-safe

**Cons:**
- Requires API key
- Costs money (OpenAI/Anthropic)

### Option 2: Use Custom Backend
**Pros:**
- Free (if using local models)
- Full control
- Already implemented

**Cons:**
- Limited AI capabilities
- No streaming (unless implemented)

### Option 3: Hybrid Approach
**Pros:**
- Use AI SDK for generation
- Use custom backend for other features
- Best of both worlds

**Cons:**
- More complex setup

---

## 💻 Implementation: Vercel AI SDK

### Step 1: Create AI Service

**File:** `src/lib/services/ai.ts`

```typescript
import { generateText, streamText } from 'ai';
import { openai } from '@ai-sdk/openai';

/**
 * Generate app code using AI
 */
export async function generateApp(prompt: string, options: {
  appType: 'erp' | 'ecommerce' | 'blog' | 'dashboard' | 'social' | 'custom';
  framework: 'react' | 'vue' | 'svelte' | 'nextjs';
  theme: string;
  includeAuth?: boolean;
  includeDatabase?: boolean;
  includeAPI?: boolean;
  includeTests?: boolean;
}) {
  const systemPrompt = `You are an expert full-stack developer. Generate a complete, production-ready ${options.appType} application using ${options.framework}.

Requirements:
- App Type: ${options.appType}
- Framework: ${options.framework}
- Theme: ${options.theme}
- Include Authentication: ${options.includeAuth ? 'Yes' : 'No'}
- Include Database Schema: ${options.includeDatabase ? 'Yes' : 'No'}
- Include API Endpoints: ${options.includeAPI ? 'Yes' : 'No'}
- Include Unit Tests: ${options.includeTests ? 'Yes' : 'No'}

Generate a complete file structure with all necessary files. Return the response as a JSON object with this structure:
{
  "files": [
    {
      "path": "src/App.tsx",
      "content": "// file content here",
      "type": "file"
    },
    {
      "path": "src/components",
      "type": "directory",
      "children": [...]
    }
  ],
  "metadata": {
    "totalFiles": 50,
    "estimatedLines": 5000,
    "dependencies": ["react", "typescript", ...]
  }
}`;

  const { text } = await generateText({
    model: openai('gpt-4-turbo'),
    system: systemPrompt,
    prompt: prompt,
    temperature: 0.7,
    maxTokens: 4000,
  });

  // Parse the JSON response
  try {
    const result = JSON.parse(text);
    return result;
  } catch (error) {
    throw new Error('Failed to parse AI response');
  }
}

/**
 * Stream app generation (for real-time updates)
 */
export async function streamAppGeneration(
  prompt: string,
  options: any,
  onChunk: (chunk: string) => void
) {
  const systemPrompt = `Generate a ${options.appType} app...`;

  const { textStream } = await streamText({
    model: openai('gpt-4-turbo'),
    system: systemPrompt,
    prompt: prompt,
  });

  for await (const chunk of textStream) {
    onChunk(chunk);
  }
}
```

---

### Step 2: Create Builder Store (with AI SDK)

**File:** `src/lib/stores/builder.ts`

```typescript
import { writable } from 'svelte/store';
import { generateApp } from '$lib/services/ai';
import type { FileNode } from '$lib/types';

interface BuilderState {
  appType: 'erp' | 'ecommerce' | 'blog' | 'dashboard' | 'social' | 'custom' | null;
  framework: 'react' | 'vue' | 'svelte' | 'nextjs' | null;
  theme: string | null;
  prompt: string;
  options: {
    includeAuth: boolean;
    includeDatabase: boolean;
    includeAPI: boolean;
    includeTests: boolean;
    includeDocker: boolean;
    includeCI: boolean;
  };
  fileStructure: FileNode[];
  selectedFile: FileNode | null;
  isGenerating: boolean;
  error: string | null;
  currentStep: 1 | 2 | 3 | 4;
}

function createBuilderStore() {
  const initialState: BuilderState = {
    appType: null,
    framework: null,
    theme: null,
    prompt: '',
    options: {
      includeAuth: true,
      includeDatabase: true,
      includeAPI: true,
      includeTests: true,
      includeDocker: false,
      includeCI: false,
    },
    fileStructure: [],
    selectedFile: null,
    isGenerating: false,
    error: null,
    currentStep: 1,
  };

  const { subscribe, set, update } = writable<BuilderState>(initialState);

  return {
    subscribe,

    setAppType(appType: BuilderState['appType']) {
      update((state) => ({ ...state, appType, currentStep: 2 }));
    },

    setFramework(framework: BuilderState['framework']) {
      update((state) => ({ ...state, framework }));
    },

    setTheme(theme: string) {
      update((state) => ({ ...state, theme }));
    },

    nextStep() {
      update((state) => ({
        ...state,
        currentStep: Math.min(4, state.currentStep + 1) as 1 | 2 | 3 | 4,
      }));
    },

    prevStep() {
      update((state) => ({
        ...state,
        currentStep: Math.max(1, state.currentStep - 1) as 1 | 2 | 3 | 4,
      }));
    },

    setPrompt(prompt: string) {
      update((state) => ({ ...state, prompt }));
    },

    toggleOption(option: keyof BuilderState['options']) {
      update((state) => ({
        ...state,
        options: {
          ...state.options,
          [option]: !state.options[option],
        },
      }));
    },

    async generate() {
      update((state) => ({ ...state, isGenerating: true, error: null }));

      try {
        const state = get(builderStore);

        if (!state.appType || !state.framework || !state.theme) {
          throw new Error('Please complete all steps');
        }

        // Call AI SDK
        const result = await generateApp(state.prompt, {
          appType: state.appType,
          framework: state.framework,
          theme: state.theme,
          includeAuth: state.options.includeAuth,
          includeDatabase: state.options.includeDatabase,
          includeAPI: state.options.includeAPI,
          includeTests: state.options.includeTests,
        });

        update((state) => ({
          ...state,
          fileStructure: result.files,
          selectedFile: result.files[0],
          isGenerating: false,
        }));
      } catch (error) {
        update((state) => ({
          ...state,
          error: error instanceof Error ? error.message : 'Generation failed',
          isGenerating: false,
        }));
      }
    },

    selectFile(file: FileNode) {
      update((state) => ({ ...state, selectedFile: file }));
    },

    reset() {
      set(initialState);
    },
  };
}

export const builderStore = createBuilderStore();
```

---

### Step 3: Update Builder Page Component

**File:** `src/routes/(app)/builder/+page.svelte`

```svelte
<script lang="ts">
  import { builderStore } from '$lib/stores/builder';
  import Button from '$lib/components/ui/Button.svelte';
  import AppTypeSelector from '$lib/components/builder/AppTypeSelector.svelte';
  import ThemeSelector from '$lib/components/builder/ThemeSelector.svelte';
  import PromptInput from '$lib/components/builder/PromptInput.svelte';
  import ReviewStep from '$lib/components/builder/ReviewStep.svelte';
  import FileTree from '$lib/components/builder/FileTree.svelte';
  import CodePreview from '$lib/components/builder/CodePreview.svelte';

  let currentStep = $derived($builderStore.currentStep);
  let isGenerating = $derived($builderStore.isGenerating);
  let fileStructure = $derived($builderStore.fileStructure);
  let selectedFile = $derived($builderStore.selectedFile);
</script>

<div class="builder-page">
  <h1>🏗️ App Builder</h1>

  {#if currentStep === 1}
    <AppTypeSelector />
  {:else if currentStep === 2}
    <ThemeSelector />
  {:else if currentStep === 3}
    <PromptInput />
  {:else if currentStep === 4}
    <ReviewStep />
  {/if}

  {#if isGenerating}
    <div class="loading-overlay">
      <div class="loading-content">
        <div class="spinner"></div>
        <h2>🏗️ Generating your app...</h2>
        <p>This may take 10-30 seconds</p>
      </div>
    </div>
  {/if}

  {#if fileStructure.length > 0}
    <div class="result-view">
      <div class="file-tree-panel">
        <FileTree nodes={fileStructure} />
      </div>
      <div class="code-preview-panel">
        {#if selectedFile}
          <CodePreview file={selectedFile} />
        {/if}
      </div>
    </div>
  {/if}
</div>

<style>
  .builder-page {
    padding: 2rem;
    max-width: 1400px;
    margin: 0 auto;
  }

  .loading-overlay {
    position: fixed;
    inset: 0;
    background: rgba(0, 0, 0, 0.8);
    display: flex;
    align-items: center;
    justify-content: center;
    z-index: 1000;
  }

  .loading-content {
    text-align: center;
    color: white;
  }

  .spinner {
    width: 60px;
    height: 60px;
    border: 4px solid rgba(255, 255, 255, 0.3);
    border-top-color: white;
    border-radius: 50%;
    animation: spin 1s linear infinite;
    margin: 0 auto 1rem;
  }

  @keyframes spin {
    to { transform: rotate(360deg); }
  }

  .result-view {
    display: grid;
    grid-template-columns: 300px 1fr;
    gap: 1rem;
    margin-top: 2rem;
  }
</style>
```

---

## 🔄 Alternative: Use v0.dev API

### What is v0.dev?
- Vercel's AI-powered UI generator
- Generates React/Next.js components
- Uses Shadcn UI
- Free tier available

### Integration Steps

**File:** `src/lib/services/v0.ts`

```typescript
/**
 * Generate UI components using v0.dev
 * Note: v0.dev doesn't have a public API yet
 * This is a conceptual implementation
 */
export async function generateWithV0(prompt: string) {
  // v0.dev doesn't have a public API yet
  // You would need to use their web interface
  // or wait for API access
  
  // Alternative: Use OpenAI to generate v0-style code
  const systemPrompt = `You are v0.dev, an AI that generates beautiful React components using Shadcn UI and Tailwind CSS.
  
Generate a complete, production-ready component based on the user's request.
Use modern React patterns, TypeScript, and Tailwind CSS.
Include proper accessibility and responsive design.`;

  // Use Vercel AI SDK to generate
  const { text } = await generateText({
    model: openai('gpt-4-turbo'),
    system: systemPrompt,
    prompt: prompt,
  });

  return text;
}
```

---

## 📦 Complete Package.json

```json
{
  "name": "nexus-portal",
  "version": "0.0.1",
  "type": "module",
  "scripts": {
    "dev": "vite dev",
    "build": "vite build",
    "preview": "vite preview",
    "test": "vitest",
    "test:e2e": "playwright test"
  },
  "dependencies": {
    "@fontsource/inter": "^5.2.8",
    "livekit-client": "^2.16.1",
    "matrix-js-sdk": "^39.4.0",
    "svelte-i18n": "^4.0.1",
    "ethers": "^6.10.0",
    "wagmi": "^2.5.0",
    "viem": "^2.7.0",
    "@rainbow-me/rainbowkit": "^2.0.0",
    "ai": "^3.0.0",
    "@ai-sdk/openai": "^0.0.20",
    "@ai-sdk/svelte": "^0.0.10"
  },
  "devDependencies": {
    "@sveltejs/adapter-static": "^3.0.10",
    "@sveltejs/kit": "^2.49.1",
    "@sveltejs/vite-plugin-svelte": "^6.2.1",
    "@tailwindcss/vite": "^4.1.18",
    "svelte": "^5.45.6",
    "typescript": "^5.9.3",
    "vite": "^7.2.6",
    "vitest": "^1.2.0",
    "@playwright/test": "^1.41.0",
    "@testing-library/svelte": "^4.1.0"
  }
}
```

---

## 🎯 Usage Examples

### Example 1: Generate Mini ERP

```typescript
import { generateApp } from '$lib/services/ai';

const result = await generateApp(
  'Create a Mini ERP for manufacturing with inventory, purchase orders, CRM, and HR',
  {
    appType: 'erp',
    framework: 'react',
    theme: 'modern',
    includeAuth: true,
    includeDatabase: true,
    includeAPI: true,
    includeTests: true,
  }
);

console.log(result.files); // Array of FileNode
console.log(result.metadata); // { totalFiles, estimatedLines, dependencies }
```

### Example 2: Stream Generation

```typescript
import { streamAppGeneration } from '$lib/services/ai';

await streamAppGeneration(
  'Create an e-commerce store',
  { appType: 'ecommerce', framework: 'nextjs', theme: 'glassmorphic' },
  (chunk) => {
    console.log('Received chunk:', chunk);
    // Update UI with chunk
  }
);
```

---

## 🔐 Security Best Practices

### 1. Never Expose API Keys in Frontend

**❌ Wrong:**
```typescript
const apiKey = 'sk-...'; // Never do this!
```

**✅ Correct:**
```typescript
// Use environment variables
const apiKey = import.meta.env.VITE_OPENAI_API_KEY;

// Or better: Call backend API
const response = await fetch('/api/generate', {
  method: 'POST',
  body: JSON.stringify({ prompt }),
});
```

### 2. Create Backend API Route

**File:** `src/routes/api/generate/+server.ts`

```typescript
import { generateApp } from '$lib/services/ai';
import { json } from '@sveltejs/kit';
import type { RequestHandler } from './$types';

export const POST: RequestHandler = async ({ request }) => {
  const { prompt, options } = await request.json();

  try {
    const result = await generateApp(prompt, options);
    return json(result);
  } catch (error) {
    return json(
      { error: 'Generation failed' },
      { status: 500 }
    );
  }
};
```

---

## 📊 Cost Estimation

### OpenAI GPT-4 Turbo Pricing
- Input: $10 / 1M tokens
- Output: $30 / 1M tokens

### Typical App Generation
- Input: ~1,000 tokens (prompt + system)
- Output: ~3,000 tokens (code)
- **Cost per generation: ~$0.11**

### Monthly Estimates
- 100 generations/month: ~$11
- 1,000 generations/month: ~$110
- 10,000 generations/month: ~$1,100

---

## 🚀 Deployment

### Environment Variables (Production)

```bash
# .env.production
OPENAI_API_KEY=sk-...
VITE_BACKEND_URL=https://api.nexusportal.com
```

### Vercel Deployment

```bash
# Install Vercel CLI
npm install -g vercel

# Deploy
vercel

# Set environment variables
vercel env add OPENAI_API_KEY
```

---

## ✅ Integration Checklist

- [ ] Install `ai` package
- [ ] Install `@ai-sdk/openai` (or other provider)
- [ ] Install `@ai-sdk/svelte`
- [ ] Create `.env.local` with API key
- [ ] Create `src/lib/services/ai.ts`
- [ ] Update `src/lib/stores/builder.ts`
- [ ] Create backend API route (for security)
- [ ] Update Builder page component
- [ ] Test generation with sample prompts
- [ ] Add error handling
- [ ] Add loading states
- [ ] Add cost monitoring
- [ ] Deploy to production

---

## 🎉 Summary

**الان می‌تونید:**
1. ✅ Vercel AI SDK رو نصب کنید
2. ✅ با OpenAI/Anthropic/Google integrate کنید
3. ✅ App generation رو با AI انجام بدید
4. ✅ Streaming responses داشته باشید
5. ✅ Type-safe API داشته باشید

**فایل‌های مورد نیاز:**
- `src/lib/services/ai.ts` - AI service
- `src/lib/stores/builder.ts` - Builder store
- `src/routes/api/generate/+server.ts` - Backend API
- `.env.local` - API keys

**هزینه:**
- ~$0.11 per generation
- ~$11/month for 100 generations

**بعدی:**
```bash
npm install ai @ai-sdk/openai @ai-sdk/svelte
```

همه چیز آماده است! 🚀
