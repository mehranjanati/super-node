> **⚠️ توجه (آپدیت معماری):** این سند ممکن است شامل جزئیات قدیمی باشد. برای مشاهده معماری نهایی سیستم (شامل اضافه شدن Rivet، Matrix، OpenClaw، Redpanda و LiveKit) حتماً به [FULL_SYSTEM_ARCHITECTURE.md](./FULL_SYSTEM_ARCHITECTURE.md) و برای برنامه اجرایی به [MVP_DEVELOPMENT_PLAN.md](./MVP_DEVELOPMENT_PLAN.md) مراجعه کنید.

# 🚀 معماری پیشرفته سه‌لایه (Advanced 3-Tier Architecture)
## SvelteKit (Bun + Redis) به عنوان BFF و Go به عنوان Super Node

این سند معماری پیشرفته و نهایی سیستم را برای ادغام **Vercel AI SDK** با بک‌اند **Go** از طریق یک لایه میانی (BFF) مبتنی بر **Bun** و **Redis** شرح می‌دهد.

---

## 🌟 چرا این معماری؟ (مزایای استراتژیک)
با اضافه شدن Bun و Redis به عنوان لایه BFF (Backend for Frontend) در دل SvelteKit، ما پیچیدگی پیاده‌سازی دستی پروتکل‌های استریمینگ در Go را حذف می‌کنیم. 
* **Go (Super Node):** فقط روی منطق تجاری سنگین (دیپلوی، بلاکچین، پردازش) تمرکز می‌کند.
* **Bun (BFF):** با سرعت فوق‌العاده، ارتباط با LLM (OpenAI/Anthropic) و استریمینگ به کلاینت را مدیریت می‌کند.
* **Redis (بومی Bun):** کش کردن، مدیریت Rate Limit و ارتباط Pub/Sub بین Go و Bun را با تاخیر صفر (Zero-latency) انجام می‌دهد.

---

## 🏗️ ساختار لایه‌ها

### لایه ۱: کلاینت (SvelteKit Frontend)
* **وظیفه:** رندر کردن UI، مدیریت فرم چت، نمایش وضعیت ابزارها (Generative UI).
* **تکنولوژی:** Svelte 5, Tailwind, `@ai-sdk/svelte`
* **نحوه کار:** فقط از هوک `useChat` استفاده می‌کند و به مسیر `/api/chat` (که در لایه ۲ است) متصل می‌شود.

### لایه ۲: لایه میانی هوشمند (SvelteKit Server / BFF)
* **وظیفه:** دریافت درخواست کلاینت، ارتباط با LLM، مدیریت Tool Calling، کش کردن پاسخ‌ها.
* **تکنولوژی:** Bun Runtime, SvelteKit Server Routes (`+server.ts`), Bun Native Redis, Vercel AI SDK Core.
* **نحوه کار:** 
  1. درخواست را می‌گیرد.
  2. Redis را چک می‌کند (Rate Limit و Cache).
  3. با LLM ارتباط برقرار می‌کند (`streamText`).
  4. اگر LLM نیاز به ابزاری داشت، یک درخواست HTTP/REST به لایه ۳ (Go) می‌فرستد.

### لایه ۳: هسته پردازشی (Go Super Node)
* **وظیفه:** اجرای دستورات واقعی (دیپلوی کانتینر، تراکنش‌های مالی، کوئری‌های دیتابیس).
* **تکنولوژی:** Golang, Docker/Podman SDK, Temporal (در صورت نیاز).
* **نحوه کار:** APIهای ساده و سریع (REST یا gRPC) ارائه می‌دهد که فقط توسط لایه ۲ (BFF) فراخوانی می‌شوند.

---

## 🛠️ پلن اجرایی و پیاده‌سازی (مرحله به مرحله)

### گام ۱: تنظیم SvelteKit برای اجرا روی Bun
ابتدا باید پروژه فرانت‌اند (`portal1`) را برای استفاده از Bun پیکربندی کنیم.
```bash
# در پوشه portal1
bun add -d svelte-adapter-bun
```
سپس در `svelte.config.js`:
```javascript
import adapter from 'svelte-adapter-bun';
// ...
kit: {
  adapter: adapter()
}
```

### گام ۲: پیاده‌سازی مسیر `/api/chat` در BFF با ادغام Redis
ایجاد فایل `src/routes/api/chat/+server.ts`:

```typescript
import { streamText, tool } from 'ai';
import { openai } from '@ai-sdk/openai';
import { z } from 'zod';
import { Redis } from 'bun';

// اتصال فوق سریع و بومی Bun به Redis
const redis = new Redis("redis://localhost:6379");

export async function POST({ request, getClientAddress }) {
  const ip = getClientAddress();
  
  // ۱. Rate Limiting با Redis
  const requests = await redis.incr(`rate_limit:${ip}`);
  if (requests === 1) await redis.expire(`rate_limit:${ip}`, 60);
  if (requests > 20) return new Response('Too Many Requests', { status: 429 });

  const { messages } = await request.json();

  // ۲. ارتباط با LLM و تعریف ابزارها
  const result = streamText({
    model: openai('gpt-4o'),
    messages,
    tools: {
      deploy_website: tool({
        description: 'Deploy a new website or container',
        parameters: z.object({
          projectName: z.string(),
          template: z.string(),
        }),
        execute: async ({ projectName, template }) => {
          // ۳. ارتباط BFF با بک‌اند Go
          // به جای درگیری با استریم، یک درخواست ساده به Go می‌زنیم
          const goResponse = await fetch('http://localhost:8080/internal/tools/deploy', {
            method: 'POST',
            headers: { 'Content-Type': 'application/json' },
            body: JSON.stringify({ project_name: projectName, template })
          });
          
          if (!goResponse.ok) throw new Error("Go Backend Failed");
          return await goResponse.json(); // نتیجه به LLM برمی‌گردد
        },
      }),
    },
  });

  // ۴. بازگرداندن استریم استاندارد به کلاینت
  return result.toDataStreamResponse();
}
```

### گام ۳: ساده‌سازی بک‌اند Go
بک‌اند Go دیگر نیازی به هندل کردن `0:` و `9:` ندارد. فقط یک API ساده برای اجرای ابزارها می‌سازیم:

```go
// در بک‌اند Go (Super Node)
func DeployToolHandler(w http.ResponseWriter, r *http.Request) {
    var req struct {
        ProjectName string `json:"project_name"`
        Template    string `json:"template"`
    }
    json.NewDecoder(r.Body).Decode(&req)

    // منطق سنگین دیپلوی...
    deployURL := runDeployment(req.ProjectName)

    // بازگرداندن یک JSON ساده به BFF
    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(map[string]string{
        "status": "success",
        "url":    deployURL,
    })
}
```

### گام ۴ (پیشرفته): استفاده از Redis Pub/Sub برای کارهای طولانی (Long-running Tasks)
اگر دیپلوی در Go چند دقیقه طول می‌کشد، HTTP Request تایم‌اوت می‌شود. راه حل:
1. **BFF (Bun):** درخواست دیپلوی را به Go می‌فرستد و بلافاصله به یک کانال Redis سابسکرایب می‌کند.
2. **Go:** پروسه را در بک‌گراند شروع می‌کند و وضعیت‌ها را در Redis پابلیش می‌کند (`PUBLISH deploy_status "building..."`).
3. **BFF (Bun):** پیام‌های Redis را دریافت کرده و از طریق استریم Vercel AI (یا SSE مجزا) به فرانت‌اند می‌فرستد تا کاربر نوار پیشرفت را ببیند.

---

## 🎯 نتیجه‌گیری
این معماری **Advance**، بهترین‌های هر سه دنیا را ترکیب می‌کند:
1. **SvelteKit + Vercel AI:** بهترین تجربه کاربری (UX) و توسعه‌دهنده (DX) برای چت هوش مصنوعی.
2. **Bun + Redis:** سریع‌ترین لایه میانی ممکن برای مدیریت ترافیک، کش و استریمینگ.
3. **Go:** قدرتمندترین موتور برای پردازش‌های سیستمی، کانتینرها و بلاکچین.
