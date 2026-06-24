# VoltAgent Standalone Architecture Decision

این سند تصمیم معماری رسمی برای جایگاه `VoltAgent` در استک `Nexus Super Node` است.

هدف این سند این است که ابهام فعلی بین:

- `VoltAgent` به عنوان یک مفهوم در معماری
- `VoltAgent` به عنوان یک ماژول embedded داخل Go
- `VoltAgent` به عنوان یک سرویس مستقل در استک

را به یک تصمیم اجرایی روشن تبدیل کند.

---

## وضعیت تصمیم

**تصمیم نهایی:**

`VoltAgent` از این نقطه به بعد باید به عنوان **سرویس مستقل (Standalone Service)** در استک توسعه داده، deploy و integrate شود.

این تصمیم جایگزین مدل فعلی `embedded-only` است.

---

## مسئله

در وضعیت فعلی پروژه، بین اسناد، کد و مدل ذهنی محصول یک ناهماهنگی وجود دارد:

1. در اسناد معماری، `VoltAgent` به عنوان مغز تصمیم‌گیری و reasoning engine معرفی شده است.
2. در کد فعلی، `VoltAgent` به شکل یک module داخلی داخل بک‌اند Go پیاده شده است.
3. در runtime فعلی، سرویس مستقل `voltagent-service` در استک بالا می‌آید، اما هنوز مسیر اصلی backend و edge به‌صورت end-to-end به آن migrate نشده‌اند.

نتیجه این ناهماهنگی:

- مرز مسئولیت بین `Go` و `VoltAgent` مبهم است.
- توسعه agentic behavior در مسیرهای مختلف پخش شده است.
- اضافه کردن supervisor logic، sub-agent orchestration و plan execution شفاف نیست.
- migration به معماری نهایی سخت‌تر می‌شود، چون implementation فعلی جهت‌گیری نهایی را منعکس نمی‌کند.

---

## چرا این تصمیم گرفته شد

ما `VoltAgent` را سرویس مستقل می‌کنیم چون:

1. **Reasoning باید از orchestration جدا شود.**
   بک‌اند Go نباید هم‌زمان API gateway، workflow orchestrator، event consumer و reasoning engine باشد.

2. **مرز معماری واضح‌تر می‌شود.**
   وقتی `VoltAgent` سرویس مستقل باشد، مشخص می‌شود تصمیم‌گیری کجا انجام می‌شود و اجرا کجا.

3. **توسعه multi-agent واقعی ممکن می‌شود.**
   supervisor/sub-agent patterns، tool planning و task decomposition در یک سرویس مستقل قابل رشدتر از یک package embedded هستند.

4. **استک به vision اصلی پروژه نزدیک‌تر می‌شود.**
   در معماری هدف، `VoltAgent` یک capability اصلی است، نه فقط helper code داخل Go.

5. **مقیاس‌پذیری آینده بهتر می‌شود.**
   در آینده می‌توان `VoltAgent` را مستقل scale، observe و harden کرد، بدون اینکه کل Go backend را دست بزنیم.

---

## چیزی که این تصمیم نیست

این تصمیم به معنی این نیست که:

- `VoltAgent` جای `Go` را بگیرد
- `Go` فقط به proxy ساده تبدیل شود
- همه منطق پروژه فورا از Go به `VoltAgent` migrate شود
- کل استک به یک rewrite بزرگ وارد شود

این سند فقط می‌گوید:

- reasoning layer باید مستقل باشد
- `VoltAgent` باید وارد استک به عنوان سرویس واقعی شود
- migration باید مرحله‌ای و کم‌ریسک باشد

---

## جایگاه نهایی در معماری

### نقش `BFF`

- مدیریت chat UX
- streaming پاسخ‌ها
- session-facing behavior
- تماس با `Go` به عنوان backend اصلی

### نقش `Go Super Node`

- API gateway داخلی سیستم
- مدیریت persistence و read model
- trigger کردن workflowها در `Temporal`
- ارتباط با `Redpanda`, `Hasura`, `Rivet`, `OpenClaw`, `LiveKit`
- اجرای system-level side effects

### نقش `VoltAgent`

- reasoning
- planning
- tool selection
- task decomposition
- supervisor/sub-agent coordination
- تبدیل intent یا context به execution plan ساخت‌یافته

### نقش `Temporal`

- اجرای durable و قابل‌اتکا برای workflowها
- retry، statefulness و long-running tasks

### نقش `Rivet / Wasm`

- execution engineهای تخصصی
- نه لایه تصمیم‌گیری مرکزی

---

## اصل مرزی مهم

`VoltAgent` باید **تصمیم بگیرد**، اما `Go` باید **اجرا و ثبت کند**.

به زبان ساده:

- `VoltAgent` می‌گوید: "برای این درخواست چه planی مناسب است؟"
- `Go` می‌گوید: "این plan را با Temporal / MCP / Rivet / DB چگونه اجرا کنم؟"

این تفکیک، مرز مرکزی معماری جدید است.

---

## مدل ارتباطی جدید

مدل پایه ارتباطی به شکل زیر خواهد بود:

```text
User -> SPA -> BFF -> Go -> VoltAgent
                          -> Temporal
                          -> Redpanda
                          -> Rivet
                          -> Persistence
```

در این مدل:

- `BFF` مستقیم با `VoltAgent` حرف نمی‌زند
- `Go` تنها ورودی سیستم برای orchestration باقی می‌ماند
- `VoltAgent` به عنوان decision service توسط `Go` فراخوانی می‌شود

---

## تصمیم مهم: BFF مستقیم با VoltAgent حرف نمی‌زند

این تصمیم آگاهانه است.

دلیل:

1. اگر `BFF` مستقیم با `VoltAgent` حرف بزند، orchestration بین دو لایه پخش می‌شود.
2. `Go` باید نقطه کنترل مرکزی برای auth، workflow، persistence و observability باقی بماند.
3. `VoltAgent` باید پشت contract داخلی محافظت شود، نه اینکه مستقیما در edge exposed شود.

پس در فاز فعلی:

- `BFF -> Go`
- `Go -> VoltAgent`

و نه:

- `BFF -> VoltAgent`

---

## API اولیه VoltAgent

در فاز اول، `VoltAgent` فقط یک API کوچک و روشن خواهد داشت.

contract canonical این API در این سند نگهداری می‌شود:

- `docs/VOLTAGENT_SERVICE_CONTRACT.md`

نمونه‌های این بخش فقط جهت توضیح معماری هستند. source of truth برای schema، versioning و taxonomy خطاها این سند است:

- `docs/VOLTAGENT_SERVICE_CONTRACT.md`

### 1. Health

```http
GET /health
```

نمونه پاسخ:

```json
{
  "status": "ok",
  "service": "voltagent-service",
  "contract_version": "v1alpha1",
  "request_id": "0f7c95a0-0d8c-4d0f-a8f3-9e245bfa7a61",
  "checks": {
    "api": "ok",
    "planner": "ok"
  }
}
```

### 2. Plan

```http
POST /plan
Content-Type: application/json
```

نمونه درخواست:

```json
{
  "contract_version": "v1alpha1",
  "intent": "deploy_website",
  "input": {
    "project_name": "demo-site",
    "prompt": "A minimal landing page for a crypto product",
    "framework": "svelte",
    "theme": "minimal"
  },
  "context": {
    "user_id": "usr_123",
    "session_id": "sess_456",
    "source": "go-gateway"
  }
}
```

نمونه پاسخ:

```json
{
  "status": "ok",
  "contract_version": "v1alpha1",
  "request_id": "0f7c95a0-0d8c-4d0f-a8f3-9e245bfa7a61",
  "plan": {
    "intent": "deploy_website",
    "kind": "workflow",
    "execution_target": "go-temporal",
    "workflow": {
      "name": "website-deployment-v1",
      "action": "start_dynamic_pipeline",
      "input": {
        "project_name": "demo-site",
        "prompt": "A minimal landing page for a crypto product",
        "framework": "svelte",
        "theme": "minimal"
      }
    },
    "artifacts": {
      "project_name": "demo-site"
    },
    "warnings": []
  }
}
```

### 3. Execute

در فاز اول، `POST /execute` جزو contract canonical نیست و source of truth برای API داخلی فقط این دو endpoint هستند:

- `GET /health`
- `POST /plan`

اگر در آینده endpointی مثل `POST /execute` اضافه شود:

- باید lightweight و non-durable باشد
- نباید جای execution اصلی در `Go + Temporal` را بگیرد
- باید با versioning و contract جداگانه مستند شود

---

## مسئولیت‌هایی که از Go خارج می‌شوند

به مرور زمان این بخش‌ها از Go به `VoltAgent` منتقل می‌شوند:

- انتخاب tool بر اساس intent
- plan generation
- supervisor logic
- تصمیم‌گیری درباره اینکه workflow لازم است یا execution سبک کافی است
- orchestration ذهنی بین sub-agentها

---

## مسئولیت‌هایی که در Go باقی می‌مانند

این موارد باید در `Go` بمانند:

- routeهای اصلی سیستم
- DB writes
- workflow triggering
- integration با Temporal
- integration با Redpanda
- integration با Rivet
- integration با Hasura/read model
- system audit trail
- fallback behavior

---

## وضعیت implementation فعلی

در implementation فعلی:

- سرویس مستقل `voltagent-service` با endpointهای `GET /health` و `POST /plan` وجود دارد
- `VoltAgentClient` در Go ایجاد شده است
- `VoltAgentService` embedded هنوز داخل کد Go وجود دارد
- routeهای `/voltagent/manifest` و `/voltagent/execute` هنوز داخل gateway Go ثبت شده‌اند
- chat path فعلی هنوز از service embedded عبور می‌کند
- wiring runtime اصلی هنوز به‌صورت کامل به `voltagent-service` منتقل نشده است

این implementation فعلی باید به عنوان:

**bridge implementation**

در نظر گرفته شود، نه architecture نهایی.

---

## تصمیم migration

ما migration را به صورت مرحله‌ای انجام می‌دهیم، نه big-bang.

### فاز 1: Standalone Skeleton

- ایجاد سرویس جدید `voltagent-service`
- اضافه شدن آن به `docker-compose`
- health endpoint
- config و env مستقل
- logging مستقل

وضعیت فعلی:

- این فاز عملا انجام شده است

### فاز 2: Client در Go

- ایجاد `VoltAgentClient` در Go
- تماس `Go -> VoltAgent /plan`
- حفظ implementation فعلی embedded به عنوان fallback موقت

وضعیت فعلی:

- client ایجاد شده است
- wiring runtime هنوز کامل نشده است

### فاز 3: اولین use case واقعی

اولین use case پیشنهادی:

- `deploy_website`

جریان:

1. `BFF` درخواست را به `Go` می‌فرستد
2. `Go` درخواست planning را به `VoltAgent` می‌فرستد
3. `VoltAgent` execution plan برمی‌گرداند
4. `Go` plan را به workflow یا tool execution ترجمه و اجرا می‌کند
5. نتیجه در persistence/read model ثبت می‌شود

### فاز 4: حذف منطق embedded از مسیر اصلی

- routeهای قدیمی `voltagent` در Go deprecated می‌شوند
- manifest/tool logic به سرویس مستقل منتقل می‌شود
- embedded module فقط تا پایان migration باقی می‌ماند

### فاز 5: Removal

بعد از stable شدن:

- `internal/core/services/voltagent` از نقش اصلی خارج می‌شود
- یا حذف می‌شود
- یا فقط به عنوان compatibility shim باقی می‌ماند

---

## preconditions قبل از شروع migration

برای اینکه migration کم‌ریسک بماند، این شرط‌ها باید رعایت شوند:

1. contract JSON request/response باید قبل از implementation نهایی شود
2. `Go` باید client داخلی واضح برای `VoltAgent` داشته باشد
3. fallback behavior باید تعریف شود
4. observability باید از ابتدا در نظر گرفته شود
5. فقط یک use case در فاز اول migrate شود

وضعیت فعلی این preconditionها:

- contract رسمی وجود دارد
- client رسمی در Go وجود دارد
- fallback در config تعریف شده اما هنوز در runtime enforce نشده است
- observability بین‌سرویسی هنوز کامل نشده است
- use case فاز اول همچنان `deploy_website` است

---

## fallback policy

اگر `VoltAgent` down باشد:

- `Go` باید خطای واضح برگرداند
- در فاز migration، fallback محدود به embedded implementation مجاز است
- بعد از نهایی شدن migration، fallback باید یا degraded mode باشد یا error مشخص

نمونه خطا:

```json
{
  "status": "error",
  "code": "VOLTAGENT_UNAVAILABLE",
  "message": "VoltAgent service is unavailable"
}
```

---

## observability

از روز اول، این موارد باید برای `VoltAgent` در نظر گرفته شوند:

- request id
- correlation id با `Go`
- structured logs
- latency metrics
- error classification
- trace boundary بین `Go` و `VoltAgent`

اگر این بخش از ابتدا گذاشته نشود، debugging کل stack خیلی سخت می‌شود.

---

## امنیت و دسترسی

در فاز اول:

- `VoltAgent` فقط در شبکه داخلی stack قابل دسترسی باشد
- public exposure نداشته باشد
- فقط `Go` client رسمی آن باشد

وضعیت فعلی:

- composeهای فعلی هنوز پورت host برای `voltagent-service` publish می‌کنند، بنابراین hardening کامل این بخش هنوز done نیست

در آینده:

- mTLS یا shared internal token
- rate limit داخلی
- auth بین سرویس‌ها

قابل اضافه شدن است.

---

## تغییرات لازم در استک

این تصمیم مستقیما این تغییرات را لازم می‌کند:

1. اضافه شدن `voltagent-service` به `docker-compose.yml`
2. اضافه شدن config/env برای `VoltAgent`
3. اضافه شدن client در Go
4. جدا شدن planning logic از `internal/core/services/voltagent`
5. اضافه شدن health check جدید به ماتریس سلامت سیستم

وضعیت فعلی:

- موردهای 1 تا 3 انجام شده‌اند
- موردهای 4 و 5 هنوز ناقص هستند

---

## چیزهایی که فعلا خارج از scope هستند

این سند فعلا درباره این‌ها تصمیم اجرایی نمی‌دهد:

- multi-region deployment برای `VoltAgent`
- auto-scaling policy
- memory subsystem کامل داخل `VoltAgent`
- migration کامل به sub-agent swarm
- جایگزینی `VoltAgent` با `Ruflo` یا frameworkهای دیگر
- تماس مستقیم `BFF` با `VoltAgent`

---

## چرا فعلا Ruflo جایگزین انتخاب نشد

`Ruflo` به عنوان reference جالب است، اما فعلا جایگزین مستقیم انتخاب نمی‌شود چون:

1. استک فعلی پروژه بر پایه `Go + Temporal + Redpanda + Hasura + Rivet` طراحی شده است.
2. `Ruflo` بیشتر یک ecosystem مستقل برای multi-agent orchestration و MCP-centric workflows است.
3. جایگزینی مستقیم در این مرحله complexity migration را بالا می‌برد.

در نتیجه:

- `Ruflo` می‌تواند منبع الهام یا benchmark باشد
- اما decision فعلی بر پایه `VoltAgent standalone` است

---

## معیار Done برای این تصمیم

این تصمیم زمانی از حالت سند به معماری واقعی تبدیل می‌شود که:

- سرویس `voltagent-service` در compose بالا بیاید
- `GET /health` پاسخ درست بدهد
- `Go` با client رسمی به آن متصل شود
- یک use case واقعی از آن عبور کند
- نتیجه از مسیر `BFF -> Go -> VoltAgent -> Go -> Temporal` قابل تست باشد

---

## جمع‌بندی

تصمیم نهایی این سند:

- `VoltAgent` باید **سرویس مستقل** باشد
- `Go` باید **orchestrator اصلی** باقی بماند
- مرز بین reasoning و execution باید **شفاف و enforceable** باشد
- migration باید **مرحله‌ای، contract-first و کم‌ریسک** انجام شود

این سند مبنای توسعه بعدی `VoltAgent` در پروژه است.
