> **⚠️ توجه (آپدیت استراتژی):** این سند ممکن است شامل جزئیات قدیمی مربوط به نسخه‌های اولیه فرانت‌اند باشد. سیستم اکنون از رویکرد **Template-Driven** و **GitOps کاربر‌محور** استفاده می‌کند. برای مشاهده استراتژی نهایی حتماً به [MASTER_PLAN.md](../../MASTER_PLAN.md) و [MVP_DEVELOPMENT_PLAN.md](../MVP_DEVELOPMENT_PLAN.md) مراجعه کنید.

# 🌐 پلن توسعه فرانت‌اند: یکپارچه‌سازی Vercel AI SDK با بک‌اند Go

این سند شامل مراحل دقیق و عملی برای اتصال فرانت‌اند (SvelteKit) به بک‌اند (Go) با استفاده از Vercel AI SDK است.

## ۱. بررسی وضعیت فعلی Vercel AI SDK در `portal1`
در حال حاضر در `portal1` پکیج‌های `@ai-sdk/svelte` و `ai` نصب شده‌اند، اما در فایل `ChatInterface.svelte` از یک سرویس Mock به نام `superNode` استفاده می‌شود. هدف ما جایگزینی این Mock با قابلیت‌های واقعی Vercel AI SDK است.

## ۲. استراتژی ارتباطی (Client-Side)
Vercel AI SDK در سمت کلاینت (Svelte) هوک قدرتمند `useChat` را ارائه می‌دهد. این هوک به صورت خودکار مدیریت وضعیت پیام‌ها، استریمینگ (Streaming) و فراخوانی ابزارها (Tool Calling) را انجام می‌دهد.

### مراحل اجرایی فرانت‌اند:

#### گام ۱: راه‌اندازی `useChat` در `ChatInterface.svelte`
باید سرویس `superNode` را حذف کرده و از `useChat` استفاده کنیم.

```svelte
<script lang="ts">
  import { useChat } from '@ai-sdk/svelte';
  
  const { messages, input, handleSubmit, isLoading } = useChat({
    api: 'http://localhost:8080/api/chat', // آدرس بک‌اند Go
    initialMessages: [
      { id: '1', role: 'assistant', content: 'سلام! من نماینده Nexus شما هستم.' }
    ],
    // مدیریت فراخوانی ابزارها در صورت نیاز
    onToolCall: ({ toolCall }) => {
      console.log('Tool called:', toolCall);
    }
  });
</script>
```

#### گام ۲: مدیریت فرم و ارسال پیام
فرم چت باید به جای توابع دستی، از `handleSubmit` و `input` که توسط `useChat` ارائه می‌شوند استفاده کند.

```svelte
<form on:submit={handleSubmit} class="flex items-center gap-2">
  <input 
    bind:value={$input} 
    placeholder="پیام خود را بنویسید..." 
    disabled={$isLoading}
  />
  <button type="submit" disabled={$isLoading}>ارسال</button>
</form>
```

#### گام ۳: رندر کردن پیام‌ها و ابزارها (Tool Invocations)
Vercel AI SDK از قابلیت Tool Calling پشتیبانی می‌کند. اگر بک‌اند Go یک ابزار (مثلاً `deploy_website`) را فراخوانی کند، در آرایه `$messages` قرار می‌گیرد.

```svelte
{#each $messages as message}
  <div class="message {message.role}">
    {message.content}
    
    <!-- نمایش وضعیت ابزارها -->
    {#if message.toolInvocations}
      {#each message.toolInvocations as tool}
        <div class="tool-call">
          در حال اجرای: {tool.toolName}
          {#if tool.state === 'result'}
            نتیجه: {JSON.stringify(tool.result)}
          {/if}
        </div>
      {/each}
    {/if}
  </div>
{/each}
```

## ۳. مدیریت CORS و Proxy
از آنجایی که فرانت‌اند روی پورت متفاوتی (مثلاً 5173) و بک‌اند Go روی پورت دیگری (مثلاً 8080) اجرا می‌شود، باید در `vite.config.ts` یک پروکسی تنظیم کنیم تا درخواست‌های `/api/chat` به بک‌اند Go هدایت شوند و مشکل CORS پیش نیاید.

```typescript
// vite.config.ts
export default defineConfig({
  server: {
    proxy: {
      '/api': {
        target: 'http://localhost:8080',
        changeOrigin: true,
      }
    }
  }
});
```

## ۴. فازهای توسعه فرانت‌اند
1. **فاز اول:** پیاده‌سازی `useChat` با یک پیام متنی ساده (بدون ابزار) برای تست اتصال استریمینگ به Go.
2. **فاز دوم:** اضافه کردن پشتیبانی از رندر کردن `toolInvocations` در UI (مانند نمایش لودینگ برای دیپلوی سایت).
3. **فاز سوم:** مدیریت خطاهای شبکه و قطع شدن استریم.
