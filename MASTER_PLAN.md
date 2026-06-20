# 🗺️ Master Plan: Nexus Super Node

این سند نقشه راه (Roadmap) نهایی، استراتژیک و یکپارچه کل سیستم Nexus Super Node است. این سند نشان می‌دهد که سیستم در سطح کلان چه کاری انجام می‌دهد و اجزای مختلف چگونه در کنار هم قرار گرفته‌اند.

---

## ۱. چشم‌انداز (Vision)
تبدیل شدن به یک **"AI Code Generator & Orchestrator"** کاملاً خودکار که می‌تواند وب‌سایت‌ها، فروشگاه‌ها، عامل‌های هوشمند (Agents) و خط‌لوله‌های داده را تنها با یک پرامپت ساده (Natural Language) تولید، پیکربندی و دیپلوی کند، در حالی که زیرساخت نهایی و مالکیت سورس‌کد **صد در صد در اختیار کاربر** است.

---

## ۲. استراتژی‌های کلیدی (Core Strategies)

این سیستم بر پایه ۴ تصمیم استراتژیک بنا شده است:

### الف) Zero Compute Cost for Frontend & Build
برای جلوگیری از بار اضافی روی سرورهای Nexus، فرآیند Build و Hosting فرانت‌اند به اکوسیستم خود کاربر منتقل شده است:
- **Frontend Hosting:** از طریق GitHub Pages (برای SSG/SPA) یا Vercel (برای SSR) دیپلوی می‌شود.
- **Wasm Compilation:** کامپایل کدهای Go به فایل‌های `.wasm` توسط GitHub Actions در ریپازیتوری کاربر انجام می‌شود.
> *رجوع شود به: `docs/GITOPS_WORKFLOW.md`*

### ب) Template-Driven Development (به جای تولید کد از صفر)
هوش مصنوعی مستقیماً کدهای پیچیده (مثل سبد خرید فروشگاه) را خط به خط نمی‌نویسد (جلوگیری از باگ و هذیان).
- **رویکرد:** سیستم از **Template Registry** شامل قالب‌های آماده، امن و متن‌باز (مانند `svelte-commerce` یا `sveltia-cms`) استفاده می‌کند.
- **نقش AI:** هوش مصنوعی صرفاً این قالب‌ها را بر اساس سلیقه کاربر کانفیگ می‌کند (تنظیم تم، محتوا، متغیرها).
> *رجوع شود به: `docs/ECOMMERCE_CMS_STRATEGY.md`*

### ج) یکپارچگی استک (Unified Go Stack)
- کل هسته بک‌اند (Super Node) با **Go** نوشته شده است.
- زبان توسعه عامل‌های لبه‌ای (Wasm Agents) نیز از Rust به **TinyGo** تغییر یافت تا تولید کد توسط AI بسیار دقیق‌تر، بدون خطای کامپایل (Borrow Checker) و کاملاً یکپارچه با هسته سیستم باشد.
> *رجوع شود به: `docs/WASM_TINYGO_STRATEGY.md`*

### د) معماری غیرمتمرکز و ارتباطی (DPIN)
- سیستم دارای یک دروازه چندکاناله به نام **OpenClaw** است.
- ارتباطات متنی امن و پایدار بر بستر **Matrix Conduit** بنا شده‌اند.
- تماس‌های صوتی و تصویری Real-time با استفاده از **LiveKit** مدیریت می‌شوند.
> *رجوع شود به: `docs/FULL_SYSTEM_ARCHITECTURE.md`*

---

## ۳. فازهای توسعه (MVP to Production)

توسعه سیستم به صورت چابک و در ۴ فاز کلیدی پیش می‌رود:

*   **فاز ۱ (هفته ۱) - Foundation:**
    راه‌اندازی دیتابیس‌ها (TiDB/Postgres)، سرور BFF (Bun) و اتصال آن به SvelteKit.
*   **فاز ۲ (هفته ۲) - AI & Execution:**
    اضافه کردن Vercel AI SDK، راه‌اندازی Temporal، ساخت Template Registry و سرویس اجرای گراف منطقی (Rivet).
*   **فاز ۳ (هفته ۳) - Events & GitOps:**
    راه‌اندازی Redpanda برای استریم رویدادها، فعال‌سازی جریان کاری GitHub Actions (GitOps) برای استقرار کد در اکانت کاربران.
*   **فاز ۴ (هفته ۴) - Advanced Comms:**
    استقرار درگاه‌های ارتباطی (Matrix، OpenClaw، LiveKit) و یکپارچگی نهایی.

> *برای جزئیات دقیق‌تر مراحل اجرایی به `docs/MVP_DEVELOPMENT_PLAN.md` مراجعه کنید.*

---

## ۴. ساختار دایرکتوری مستندات (Docs Index)

برای درک هر بخش از سیستم، به مستندات اختصاصی آن مراجعه کنید:

1.  **معماری کلان:** `docs/FULL_SYSTEM_ARCHITECTURE.md`
2.  **پلن توسعه گام‌به‌گام (MVP):** `docs/MVP_DEVELOPMENT_PLAN.md`
3.  **نقشه راه اجرایی و چک‌لیست توسعه (Detailed Roadmap):** `docs/DETAILED_DEVELOPMENT_ROADMAP.md`
4.  **استراتژی گیت‌آپس و استقرار:** `docs/GITOPS_WORKFLOW.md`
4.  **استراتژی استفاده از قالب‌های آماده:** `docs/ECOMMERCE_CMS_STRATEGY.md`
5.  **استراتژی توسعه Wasm با TinyGo:** `docs/WASM_TINYGO_STRATEGY.md`
6.  **الگوهای چندعاملی (Multi-Agent) با VoltAgent:** `docs/dev/sub_agents_and_supervisors.md`
7.  **مدیریت حافظه و سیستم RAG:** `docs/dev/memory_and_rag.md`
8.  **توسعه ابزارها و پروتکل MCP:** `docs/dev/building_mcp_tools.md`
9.  **ارزیابی‌ها و نرده‌های محافظ (Evals & Guardrails):** `docs/dev/evals_and_guardrails.md`
10. **مشاهده‌پذیری و دیباگ هوش مصنوعی (Observability):** `docs/dev/observability_and_tracing.md`

---
*آخرین بروزرسانی: همگام با معماری نهایی Template-Driven، GitOps و TinyGo.*