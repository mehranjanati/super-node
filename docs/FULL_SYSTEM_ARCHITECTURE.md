# 🏛️ معماری جامع سیستم (Nexus Super Node)
## نقشه استک، سرویس‌ها و مدل ارتباطی

این سند نمای کلی و جامع از کل سیستم، سرویس‌های درگیر، تکنولوژی‌های استفاده شده و نحوه ارتباط آن‌ها با یکدیگر را پس از بازبینی نهایی (شامل Rivet، Matrix، OpenClaw و Postgres) نشان می‌دهد.

---

## 🗺️ نقشه معماری (Architecture Diagram)

در اینجا نقشه ارتباطی سیستم با استفاده از نمودار Mermaid رسم شده است (می‌توانید این کد را در ابزارهایی مثل `mermaid.live` یا گیت‌هاب مشاهده کنید):

```mermaid
graph TD
    %% کاربران
    User((کاربر نهایی / کلاینت))

    %% لایه ۱: رابط کاربری و دسترسی
    subgraph Tier1 [لایه ۱: UI & Access Layer]
        Svelte[SvelteKit Portal<br/>UI & Chat]
        LiveKit[LiveKit<br/>Real-time Audio/Video]
        Hasura[Hasura GraphQL<br/>Data Access Layer]
        OpenClaw[OpenClaw<br/>Multi-Channel Gateway]
    end

    %% لایه ۲: ارتباطات و مسیریابی
    subgraph Tier2 [لایه ۲: Routing & Comms]
        Bun[Bun BFF<br/>AI SDK & Gateway]
        Redis[(Redis<br/>Pub/Sub & Cache)]
        Matrix[Matrix Conduit<br/>Decentralized Chat]
    end

    %% سرویس‌های هوش مصنوعی خارجی
    LLM((OpenAI / Anthropic))

    %% لایه ۳: هسته اجرایی و منطق
    subgraph Tier3 [لایه ۳: Core Execution]
        GoServer[Golang Core API<br/>Nexus Super Node]
        Temporal[Temporal.io<br/>Workflow Orchestration]
        Podman[Podman Engine<br/>Container Management]
        VoltAgent[VoltAgent<br/>Reasoning Engine]
        Rivet[Rivet Service<br/>Visual Logic Graph]
        Wasm[Wazero / TinyGo Agent<br/>Edge Runtime]
        TemplateRegistry[Template Registry<br/>CMS & E-Commerce]
    end

    %% لایه ۴: داده و رویداد
    subgraph Tier4 [لایه ۴: Data & Event Fabric]
        Redpanda[Redpanda<br/>Event Streaming]
        Connect[Redpanda Connect<br/>Data Pipelines]
        TiDB[(TiDB<br/>Main HTAP DB)]
        Postgres[(PostgreSQL<br/>Hasura & Temporal State)]
        MinIO[MinIO<br/>Object Storage]
        IPFS[IPFS<br/>Decentralized Storage]
    end

    %% محیط اکانت کاربر (User Environment)
    subgraph UserEnv [محیط اختصاصی کاربر (User GitHub)]
        GHAction[GitHub Actions<br/>CI/CD Pipeline]
        GHPages[GitHub Pages / Vercel<br/>Frontend Hosting]
    end

    %% ارتباطات لایه ۱ و ۲
    User <-->|HTTP/WS| Svelte
    User <-->|WebRTC| LiveKit
    User <-->|External Clients| OpenClaw
    Svelte <-->|Chat API| Bun
    Svelte <-->|GraphQL| Hasura
    Bun <-->|LLM API| LLM
    Bun <-->|Cache| Redis

    %% ارتباطات Matrix و OpenClaw
    OpenClaw <-->|Matrix Protocol| Matrix
    OpenClaw <-->|Stream| Redpanda

    %% ارتباطات لایه ۳
    Bun <-->|gRPC/REST| GoServer
    LiveKit <-->|Webhooks| GoServer
    GoServer <-->|Manage| Podman
    GoServer <-->|Workflows| Temporal
    GoServer <-->|Reasoning| VoltAgent
    GoServer <-->|gRPC Logic| Rivet
    GoServer <-->|Clone/Deploy| TemplateRegistry
    Temporal <-->|Execute| Wasm

    %% جریان GitOps کاربر
    GoServer -.->|OAuth Push| GHAction
    GHAction -.->|Deploy SPA/SSR| GHPages
    GHAction -.->|Release Wasm| Wasm

    %% ارتباطات لایه ۴
    GoServer <-->|Pub/Sub| Redpanda
    Hasura <-->|Query| Postgres
    Temporal <-->|State| Postgres
    GoServer <-->|Persist| TiDB
    Hasura -.->|Optional Query| TiDB
    Redpanda <-->|ETL| Connect
    Connect <-->|Store| MinIO
```

---

## 📚 تشریح لایه‌ها و استک تکنولوژی

### ۱. لایه رابط کاربری و دسترسی (UI & Access Layer)
*   **تکنولوژی‌ها:** Svelte 5, SvelteKit, Vercel AI SDK, LiveKit, Hasura, OpenClaw.
*   **نقش:** این لایه نقطه ورود تمام کاربران و کلاینت‌های خارجی به سیستم است. شامل پرتال کاربری وب (Svelte)، ارتباطات صوتی و تصویری آنی (LiveKit)، دسترسی سریع و منعطف به داده‌ها (Hasura) و دروازه ورود کانال‌های ارتباطی مختلف (OpenClaw) است.
*   **تغییرات جدید:** اضافه شدن `OpenClaw` به عنوان درگاه چندکاناله (Multi-Channel Gateway) که به کلاینت‌های خارجی اجازه اتصال می‌دهد.

### ۲. لایه ارتباطات و مسیریابی (Routing & Comms)
*   **تکنولوژی‌ها:** Bun Runtime, Redis, Matrix (Conduit).
*   **نقش:** مدیریت ترافیک، مسیریابی پیام‌ها و کش کردن اطلاعات. `Bun` به عنوان BFF و دروازه هوش مصنوعی عمل می‌کند، `Redis` وظیفه مدیریت State‌های لحظه‌ای و محدودیت درخواست‌ها را دارد. `Matrix Conduit` نیز بستر پیام‌رسانی غیرمتمرکز (Decentralized Chat) برای ارتباط امن را فراهم می‌کند.
*   **تغییرات جدید:** `Matrix Conduit` به عنوان یک سرور ارتباطی قدرتمند اضافه شده است که به همراه OpenClaw بستر ارتباطی غیرمتمرکز DPIN را شکل می‌دهد.

### ۳. لایه هسته اجرایی و منطق (Core Execution - Super Node)
*   **تکنولوژی‌ها:** Golang, Temporal, Podman, VoltAgent, Rivet Service, Wazero / Wasm Agent.
*   **نقش:** موتور اصلی و مغز متفکر سیستم. `Go Server` به عنوان هسته اصلی، ورک‌فلوها را به `Temporal` می‌سپارد، کانتینرها را با `Podman` مدیریت می‌کند و استدلال‌های پیچیده را به `VoltAgent` می‌دهد. 
*   **سرویس Rivet:** سرویس `Rivet` (مبتنی بر Node.js و gRPC) یک ماشین مجازی اجرای منطق گراف بصری است که به Go Server اجازه می‌دهد جریان‌های منطقی عامل‌ها را اجرا کند.
*   **سرویس Wasm:** `Wazero` و عامل‌های نوشته شده با `TinyGo` برای اجرای کدهای سبک و لبه‌ای (Edge) به شکل ایزوله استفاده می‌شوند. انتخاب TinyGo یکپارچگی زبان (Go) و سهولت در تولید کد توسط هوش مصنوعی را تضمین می‌کند.
*   **سیستم Template Registry:** برای تسریع توسعه و جلوگیری از تولید کد صفر (Code Gen)، سیستم از قالب‌های آماده (مثل `svelte-commerce` و `sveltia-cms`) استفاده کرده و فقط آن‌ها را کانفیگ و مستقر می‌کند.

### ۴. لایه رویداد و داده (Data & Event Fabric)
*   **تکنولوژی‌ها:** Redpanda, Redpanda Connect, TiDB, Postgres, MinIO, IPFS.
*   **نقش:** ستون فقرات ذخیره‌سازی و جریان داده‌ها. `Redpanda` (جایگزین Kafka) برای استریم رویدادها، و `Redpanda Connect` برای پردازش خط‌لوله‌های داده استفاده می‌شود. 
*   **پایگاه‌های داده:** `Postgres` به عنوان دیتابیس رابطه‌ای پایه برای سرویس‌هایی مثل Hasura و Temporal ایفای نقش می‌کند. `TiDB` دیتابیس اصلی HTAP برای کوئری‌های تحلیلی و برداری است. `MinIO` و `IPFS` نیز برای ذخیره‌سازی فایل‌ها و داده‌های خام به کار می‌روند.

### ۵. لایه میزبانی و استقرار کاربر (User-Centric GitOps Layer)
*   **تکنولوژی‌ها:** GitHub Actions, GitHub Pages, Vercel/Netlify.
*   **نقش:** تمام عملیات کامپایل Wasm و بیلد/میزبانی فرانت‌اند روی اکانت GitHub کاربر انجام می‌شود. این کار بار پردازشی را از روی Super Node برداشته و مالکیت ۱۰۰٪ کد را به کاربر می‌دهد (Zero Compute Cost).

---

## 🔄 مدل ارتباطی (Communication Flow)

برای درک بهتر، بیایید چند سناریوی واقعی را بررسی کنیم:

### سناریو ۱: چت و دیپلوی سایت
**"کاربر در چت می‌نویسد: یک سایت شخصی برای من دیپلوی کن"**
1.  **کاربر -> فرانت‌اند:** پیام از طریق Svelte ارسال می‌شود.
2.  **فرانت‌اند -> BFF:** سرور Bun پیام را گرفته و پس از بررسی Redis (Rate Limit)، به LLM می‌فرستد.
3.  **BFF -> Go Super Node:** LLM تشخیص اجرای ابزار می‌دهد و Bun درخواست را به Go Server می‌فرستد.
4.  **Go Server -> Podman/TiDB:** سرور Go از طریق Podman کانتینر را می‌سازد و نتیجه را در TiDB ثبت می‌کند.

### سناریو ۲: منطق بصری با Rivet
**"یک عامل نیاز به پردازش منطق چندمرحله‌ای (Graph) دارد"**
1.  **Go Server -> Rivet Service:** سرور Go از طریق gRPC یک گراف منطقی (فایل `.rivet-project`) را به `Rivet Service` می‌فرستد.
2.  **Rivet Service -> LLM/Tools:** ماشین مجازی Rivet گراف را اجرا کرده و در صورت نیاز با مدل‌ها یا ابزارها تعامل می‌کند.
3.  **Rivet Service -> Go Server:** نتیجه نهایی پردازش گراف به سرور Go بازگردانده می‌شود تا در Temporal یا TiDB ثبت شود.

### سناریو ۳: ارتباطات از طریق OpenClaw و Matrix
**"دریافت پیام از یک کلاینت خارجی (مثل تلگرام/Discord از طریق OpenClaw)"**
1.  **External Client -> OpenClaw:** پیام به دروازه OpenClaw می‌رسد.
2.  **OpenClaw -> Matrix / Redpanda:** پیام در شبکه Matrix برای همگام‌سازی توزیع شده و همزمان به عنوان یک رویداد در Redpanda منتشر می‌شود.
3.  **Redpanda -> Go Server:** عامل‌های مستقر در Go Server رویداد را می‌خوانند، پردازش می‌کنند (با کمک VoltAgent یا Rivet) و پاسخ را تولید می‌کنند.

---

## 🛡️ مزایای این معماری (چرا این بهترین روش است؟)

1.  **جداسازی وظایف:** هر لایه مسئولیت مشخصی دارد (رابط کاربری، مسیریابی، منطق اجرایی، ذخیره‌سازی).
2.  **توسعه‌پذیری منطق:** وجود `Rivet Service` به توسعه‌دهندگان اجازه می‌دهد تا منطق عامل‌ها را به صورت بصری و بدون نیاز به کامپایل مجدد Go Server ویرایش کنند.
3.  **ارتباطات غیرمتمرکز:** ادغام `Matrix` و `OpenClaw` سیستم را برای پشتیبانی از معماری DPIN (شبکه‌های زیرساخت فیزیکی غیرمتمرکز) و کلاینت‌های متنوع آماده می‌کند.
4.  **پایداری و مقیاس‌پذیری:** تفکیک دیتابیس‌ها (Postgres برای Stateهای داخلی و TiDB برای داده‌های اصلی) و استفاده از Redpanda باعث می‌شود سیستم توان عملیاتی بسیار بالایی داشته باشد.
