# VoltAgent Service Day 1 Baseline

این سند baseline تاریخی روز اول migration است. این نسخه sync شده تا:

- تصمیم‌ها و snapshot روز اول حفظ شوند
- با وضعیت فعلی repo تناقض نداشته باشد
- روشن شود کدام blockerهای روز اول بعدا resolve شده‌اند و کدام بخش‌ها هنوز باز هستند

## نقش این سند

این فایل دیگر نباید به‌عنوان «وضعیت زنده فعلی» خوانده شود.

کاربرد درست آن:

- ثبت تصمیم‌های foundational روز اول
- ثبت inventory اولیه dependencyهای embedded
- ثبت scope رسمی phase 1
- مقایسه‌ی day 1 با وضعیت فعلی migration

برای وضعیت زنده فعلی، سندهای مرجع این‌ها هستند:

- `docs/VOLTAGENT_SERVICE_DAILY_PLAN.md`
- `docs/VOLTAGENT_STANDALONE_ARCHITECTURE_DECISION.md`
- `docs/VOLTAGENT_SERVICE_CONTRACT.md`

## وضعیت تاییدشده در repo فعلی

- [x] baseline روز اول در فایل جدا نگه داشته شده و دیگر به‌عنوان وضعیت زنده خوانده نمی‌شود
- [x] scope رسمی phase 1 و use case اصلی `deploy_website` داخل سند ثبت شده‌اند
- [x] مرز مسئولیت `BFF`, `Go`, `VoltAgent`, `Temporal` در سند روشن شده است
- [x] blockerهای resolved و blockerهای باز به‌صورت جداگانه ثبت شده‌اند
- [x] sync note انتهای سند با وضعیت فعلی repo هم‌راستا است

## تصمیم‌های روز اول که هنوز معتبرند

در فاز اول:

- `voltagent-service` فقط نقش reasoning/planning service را می‌گیرد
- `Go` orchestrator اصلی باقی می‌ماند
- فقط یک use case واقعی migrate می‌شود: `deploy_website`
- execution واقعی همچنان در `Go + Temporal` انجام می‌شود
- fallback موقت به implementation embedded تا زمان پایدار شدن integration حفظ می‌شود

این تصمیم‌ها هنوز مبنای migration هستند و تغییر نکرده‌اند.

## snapshot روز اول

### 1. Gateway به bridge implementation متکی بود

در روز اول، `EchoGateway` مستقیما به `*voltagent.VoltAgentService` وصل بود و routeهای runtime هنوز service embedded را صدا می‌زدند.

فایل‌های کلیدی:

- `internal/adapters/gateway/echo.go`
- `internal/core/services/voltagent/types.go`
- `internal/core/services/voltagent/fx.go`

### 2. routeهای legacy داخل Go ثبت بودند

در روز اول، این routeها هنوز در gateway Go وجود داشتند:

- `GET /voltagent/manifest`
- `POST /voltagent/execute`

و هنوز مستقیم این متدها را صدا می‌زدند:

- `g.voltSvc.GetManifest()`
- `g.voltSvc.ExecuteTool(...)`

### 3. chat path خارج از scope فاز اول بود

در روز اول:

- `POST /api/chat`

هنوز به:

- `g.voltSvc.StreamChat(ctx, body.Messages)`

وصل بود و این بخش عمدا خارج از scope phase 1 باقی ماند.

### 4. planning و execution هنوز از هم جدا نشده بودند

در روز اول، منطق زیر هنوز در `internal/core/services/voltagent/voltagent.go` متمرکز بود:

- `GetManifest()`
- `ExecuteTool()`
- workflow trigger
- status lookup

این snapshot برای ثبت نقطه شروع migration مهم است.

## scope رسمی phase 1

### داخل scope

- تعریف contract داخلی `Go -> voltagent-service`
- ساخت endpointهای `GET /health` و `POST /plan`
- ساخت `VoltAgentClient` در Go
- migrate کردن use case `deploy_website`
- حفظ execution در `Go + Temporal`
- حفظ fallback embedded

### خارج از scope phase 1

- migrate کردن کامل chat path
- migrate کردن کامل manifest/tool discovery
- حذف فوری `internal/core/services/voltagent`
- تماس مستقیم `BFF -> voltagent-service`
- حذف کامل fallback

## use case رسمی phase 1

use case انتخاب‌شده:

- `deploy_website`

دلیل انتخاب در روز اول:

- implementation embedded آن از قبل وجود داشت
- مستقیم به workflowهای Go/Temporal وصل بود
- boundary بین planning و execution را خوب نشان می‌داد
- برای validation end-to-end مناسب بود

## مرز مسئولیت لایه‌ها

### BFF

- UX
- session-facing behavior
- streaming به کاربر
- تماس فقط با Go

### Go

- ورودی اصلی backend
- auth
- orchestration
- workflow trigger
- persistence
- read model
- fallback policy
- observability مرکزی

### voltagent-service

- reasoning
- planning
- tool selection
- intent normalization
- plan generation

### Temporal

- execution durable
- retries
- long-running workflows

## owner پیشنهادی

### Node / VoltAgent Service Owner

مسئول:

- API `health/plan`
- validation
- prompt/model config
- structured logs
- response contract stability

### Go Backend Owner

مسئول:

- `VoltAgentClient`
- config/env
- fallback policy
- routing integration
- workflow execution mapping

### Infra / Compose Owner

مسئول:

- docker compose
- health checks
- service discovery
- env injection
- deploy docs

### Workflow / Persistence Owner

مسئول:

- plan-to-workflow mapping
- workflow execution records
- read model visibility
- artifacts persistence

## blockerهای روز اول و وضعیت فعلی آن‌ها

### blockerهای resolved شده

موارد زیر در روز اول blocker بودند اما الان resolve شده‌اند:

1. contract رسمی `POST /plan` وجود نداشت
2. `Go -> voltagent-service` client هنوز ساخته نشده بود
3. config رسمی برای URL/timeout/fallback هنوز وجود نداشت
4. `voltagent-service` هنوز endpoint canonical برای `GET /health` و `POST /plan` نداشت
5. compose اصلی هنوز این سرویس را به‌صورت رسمی load نمی‌کرد

resolved state فعلی:

- `docs/VOLTAGENT_SERVICE_CONTRACT.md` الان source of truth است
- `internal/adapters/voltagentclient/client.go` ساخته شده است
- `internal/config/config.go` و `config/config.yaml` config رسمی `voltagent` را دارند
- `voltagent-service/src/index.ts` الان `GET /health` و `POST /plan` دارد
- `docker-compose.yml` و `docker-compose.app.yml` الان `voltagent-service` را دارند

### blockerهایی که هنوز به‌صورت کامل resolve نشده‌اند

موارد زیر هنوز به‌صورت کامل بسته نشده‌اند:

1. migration کامل chat path
2. migration کامل manifest/tool discovery
3. observability کامل بین `Go` و `voltagent-service`
4. health matrix تجمیعی در backend اصلی
5. hardening کامل deployment و internal-only exposure
6. تکمیل تست end-to-end برای کل مسیر phase 1

## sync note بر اساس repo فعلی

در وضعیت فعلی repo:

- `EchoGateway` هنوز `*voltagent.VoltAgentService` را inject می‌کند
- routeهای `GET /voltagent/manifest` و `POST /voltagent/execute` هنوز وجود دارند
- `POST /api/chat` هنوز از service embedded عبور می‌کند
- `voltagent-service` در composeهای اصلی stack حضور دارد
- `VoltAgentClient` در Go وجود دارد
- wiring واقعی `Go -> voltagent-service` برای use case `deploy_website` با fallback embedded شروع شده است
- bridge implementation هنوز حذف نشده و بخشی از سیستم را نگه می‌دارد

این یعنی:

- baseline روز اول از نظر historical درست است
- اما دیگر نباید جملات آن به‌صورت «وضعیت فعلی زنده» خوانده شوند

## فایل‌هایی که در روزهای بعدی لمس شدند یا هنوز مهم‌اند

### Node side

- `voltagent-service/src/index.ts`
- `voltagent-service/package.json`
- `voltagent-service/Dockerfile`

### Go side

- `internal/config/config.go`
- `internal/adapters/gateway/echo.go`
- `internal/adapters/voltagentclient/client.go`
- `internal/core/services/voltagent/voltagent.go`
- `internal/core/services/voltagent/types.go`
- `internal/core/services/voltagent/fx.go`

### Infra side

- `docker-compose.yml`
- `docker-compose.app.yml`
- `config/config.yaml`

## خروجی واقعی روز اول

روز اول با این خروجی‌ها done محسوب می‌شد:

- inventory مسیرهای embedded ثبت شد
- use case فاز اول قفل شد: `deploy_website`
- boundary بین reasoning و execution روشن شد
- ownerهای پیشنهادی ثبت شدند
- blockerهای اولیه مشخص شدند

## نتیجه این baseline بعد از sync

این سند بعد از sync باید این‌طور خوانده شود:

- decision baseline: هنوز معتبر
- implementation snapshot روز اول: تاریخی
- blocker list: بخشی resolved شده، بخشی هنوز باز است
- source of truth فعلی migration: در سندهای روزانه، contract و architecture decision قرار دارد
