> **⚠️ توجه (آپدیت معماری):** این سند ممکن است شامل جزئیات قدیمی باشد. برای مشاهده معماری نهایی سیستم (شامل اضافه شدن Rivet، Matrix، OpenClaw، Redpanda و LiveKit) حتماً به [FULL_SYSTEM_ARCHITECTURE.md](./FULL_SYSTEM_ARCHITECTURE.md) و برای برنامه اجرایی به [MVP_DEVELOPMENT_PLAN.md](./MVP_DEVELOPMENT_PLAN.md) مراجعه کنید.

# 🛠️ پلن اجرایی توسعه BFF (SvelteKit + Bun + Redis)

این سند شامل مراحل عملیاتی و قدم‌به‌قدم برای پیاده‌سازی لایه BFF در پروژه `portal1` است. ما این مراحل را به ترتیب اجرا خواهیم کرد.

---

## فاز ۱: پیکربندی محیط (Bun و SvelteKit)
**هدف:** آماده‌سازی پروژه `portal1` برای اجرا روی موتور Bun.

1. **نصب آداپتور Bun:**
   اجرای دستور `bun add -d svelte-adapter-bun` در پوشه `portal1`.
2. **بروزرسانی `svelte.config.js`:**
   جایگزینی آداپتور پیش‌فرض (auto یا node) با `svelte-adapter-bun`.
3. **تست اجرای اولیه:**
   اجرای `bun run dev` برای اطمینان از صحت کارکرد پروژه.

---

## فاز ۲: راه‌اندازی کلاینت Redis در SvelteKit
**هدف:** ایجاد یک ماژول مرکزی برای ارتباط با Redis بومی Bun.

1. **ایجاد فایل `src/lib/server/redis.ts`:**
   نوشتن کلاینت سینگلتون (Singleton) برای اتصال به Redis.
   ```typescript
   import { Redis } from 'bun';
   export const redis = new Redis(process.env.REDIS_URL || "redis://localhost:6379");
   ```
2. **راه‌اندازی سرور لوکال Redis:**
   اطمینان از اجرای Redis روی سیستم (از طریق Docker یا نصب مستقیم) برای تست.

---

## فاز ۳: ایجاد هسته BFF (API Route چت)
**هدف:** ساخت نقطه اتصال (Endpoint) برای Vercel AI SDK.

1. **ایجاد مسیر API:**
   ساخت فایل `src/routes/api/chat/+server.ts`.
2. **پیاده‌سازی Rate Limiting:**
   استفاده از `redis.incr` برای محدود کردن درخواست‌های هر IP.
3. **پیاده‌سازی `streamText`:**
   اتصال به OpenAI (یا یک مدل Mock برای تست لوکال) و برگرداندن `toDataStreamResponse()`.
4. **تعریف یک ابزار تستی (Tool):**
   تعریف ابزار `deploy_website` در `streamText` که فعلاً یک لاگ ساده چاپ می‌کند (تا در فاز ۵ به Go وصل شود).

---

## فاز ۴: اتصال فرانت‌اند (ChatInterface)
**هدف:** جایگزینی کدهای Mock فعلی با هوک واقعی `useChat`.

1. **ویرایش `src/lib/components/chat/ChatInterface.svelte`:**
   حذف سرویس `superNode` قدیمی و وارد کردن `useChat` از `@ai-sdk/svelte`.
2. **اتصال فرم:**
   بایند کردن (Bind) ورودی کاربر به `$input` و فرم به `handleSubmit`.
3. **رندر پیام‌ها و ابزارها:**
   نمایش `$messages` و مدیریت وضعیت لودینگ ابزارها (`toolInvocations`) در رابط کاربری.

---

## فاز ۵: اتصال BFF به بک‌اند Go (Super Node)
**هدف:** تکمیل چرخه با ارسال درخواست اجرای ابزار از Bun به Go.

1. **ایجاد یک سرور ساده Go (Internal API):**
   نوشتن یک روت داخلی در Go (مثلاً `/internal/tools/execute`) که درخواست‌های JSON را از Bun دریافت می‌کند و یک پاسخ JSON ساده (بدون استریمینگ پیچیده) برمی‌گرداند.
2. **بروزرسانی ابزار در BFF:**
   تغییر ابزار `deploy_website` در `+server.ts` تا با `fetch` به سرور Go درخواست بفرستد، نتیجه JSON را بگیرد و آن را به صورت استریم (با پروتکل Vercel AI) به فرانت‌اند برگرداند.

---

## 🚀 وضعیت اجرا
- [ ] فاز ۱
- [ ] فاز ۲
- [ ] فاز ۳
- [ ] فاز ۴
- [ ] فاز ۵
