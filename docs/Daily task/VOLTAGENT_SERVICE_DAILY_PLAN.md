# Daily Plan For VoltAgent Service Integration

این فایل مرجع اجرایی روزبه‌روز migration مربوط به `voltagent-service` است و نسخه canonical روزانه باید در همین فولدر بماند:

- `docs/Daily task/`

قاعده اجرا:

- baseline روز اول در فایل جدا نگه داشته می‌شود
- هر روز بعدی فایل جداگانه خودش را دارد
- این فایل index اصلی روزها، وضعیت کلی و کارهای باقی‌مانده است

## فایل‌های روزانه

- `docs/Daily task/VOLTAGENT_SERVICE_DAY1_BASELINE.md`
- `docs/Daily task/VOLTAGENT_SERVICE_DAY2_BACKEND_RUNTIME.md`
- `docs/Daily task/VOLTAGENT_SERVICE_DAY3_BFF_PORTAL_ALIGNMENT.md`
- `docs/Daily task/VOLTAGENT_SERVICE_DAY4_HEALTH_AND_OBSERVABILITY.md`
- `docs/Daily task/VOLTAGENT_SERVICE_DAY5_INFRA_HARDENING.md`
- `docs/Daily task/VOLTAGENT_SERVICE_DAY6_TESTS_AND_VALIDATION.md`
- `docs/Daily task/VOLTAGENT_SERVICE_DAY7_LEGACY_CLEANUP_AND_RELEASE.md`

## وضعیت فعلی migration

### انجام‌شده

- `voltagent-service` ساخته شده و در composeهای اصلی حضور دارد
- `GET /health` و `POST /plan` وجود دارند
- contract رسمی phase 1 نوشته شده است
- `VoltAgentClient` در Go ساخته شده است
- wiring اولیه `Go -> voltagent-service` با fallback برای `deploy_website` اضافه شده است
- baseline روز اول و سند architecture decision sync شده‌اند

### هنوز باز

- migrate شدن همه call siteهای runtime
- health aggregation در Go
- observability و metricهای بین‌سرویسی
- alignment لایه BFF و portal با مسیر جدید
- hardening deployment و internal-only exposure
- تست‌های unit, integration و end-to-end
- cleanup و deprecation مسیرهای legacy

## اصل اجرای روزانه

برای هر روز:

1. فقط یک slice مشخص جلو می‌رود
2. خروجی باید قابل تست باشد
3. docs همان روز هم‌زمان update می‌شود
4. اگر کاری done نشده، باید به روز بعد منتقل و explicit ثبت شود

## روزها و وضعیت

### روز 1: Baseline و scope

وضعیت:

- `done`

خروجی:

- snapshot تاریخی migration
- scope phase 1
- ownerها و blockerهای اولیه

مرجع:

- `docs/Daily task/VOLTAGENT_SERVICE_DAY1_BASELINE.md`

### روز 2: Backend runtime wiring

وضعیت:

- `in_progress`

کارهای done:

- wiring اولیه `Go -> voltagent-service` برای `deploy_website`
- fallback embedded برای failure remote
- عبور metadataهای `request_id`, `correlation_id`, `source`
- ثبت `planning_source` در response و read model/persistence

کارهای باقی‌مانده همین روز:

- شفاف‌کردن facade داخلی برای remote/fallback تا از `VoltAgentService` embedded جدا شود
- تکمیل status/response contract برای مسیر جدید
- افزودن تست اختصاصی برای `voltagentclient`

مرجع:

- `docs/Daily task/VOLTAGENT_SERVICE_DAY2_BACKEND_RUNTIME.md`

### روز 3: BFF و portal alignment

وضعیت:

- `in_progress`

کارها:

- هماهنگ‌کردن BFF با مسیر backend جدید
- تصمیم‌گیری روشن برای portal: موقت روی legacy بماند یا migrate شود
- کاهش وابستگی مستقیم frontend به `/voltagent/manifest` و `/voltagent/execute`
- migrate شدن deploy flow در `portal1` به `/internal/tools/deploy`
- migrate شدن manifest و non-deploy tool execution در `portal1` به `/internal/tools/*`
- routeهای `/voltagent/*` فقط برای compatibility در backend باقی می‌مانند

مرجع:

- `docs/Daily task/VOLTAGENT_SERVICE_DAY3_BFF_PORTAL_ALIGNMENT.md`

### روز 4: Health و observability

وضعیت:

- `in_progress`

کارها:

- افزودن health تجمیعی در Go
- latency/error classification برای `Go -> voltagent-service`
- تعریف metricهای پایه
- مشخص‌کردن logهای remote در برابر fallback

مرجع:

- `docs/Daily task/VOLTAGENT_SERVICE_DAY4_HEALTH_AND_OBSERVABILITY.md`

### روز 5: Infra hardening

وضعیت:

- `in_progress`

کارها:

- تکمیل env matrix برای `NEXUS_VOLTAGENT_*`
- بررسی publish بودن host port و internal-only exposure
- runbook بالا آوردن stack
- hardening dependency و readiness

مرجع:

- `docs/Daily task/VOLTAGENT_SERVICE_DAY5_INFRA_HARDENING.md`

### روز 6: Tests و validation

وضعیت:

- `in_progress`

کارها:

- unit test برای `voltagentclient`
- integration test برای `POST /plan`
- smoke test برای compose
- end-to-end test برای `deploy_website`

مرجع:

- `docs/Daily task/VOLTAGENT_SERVICE_DAY6_TESTS_AND_VALIDATION.md`

### روز 7: Legacy cleanup و release readiness

وضعیت:

- `pending`

کارها:

- deprecate کردن مسیرهای embedded از مسیر اصلی
- sync نهایی docs
- release checklist
- incident/runbook کوتاه

مرجع:

- `docs/Daily task/VOLTAGENT_SERVICE_DAY7_LEGACY_CLEANUP_AND_RELEASE.md`

## لیست دقیق کارهای باقی‌مانده

### Backend

- تفکیک روشن facade runtime از bridge embedded فعلی
- migrate کردن call siteهای باقی‌مانده به مسیر جدید
- بازنگری contract response برای routeهای جدید
- افزودن تست اختصاصی برای `voltagentclient`

### Runtime

- عبور BFF از مسیر جدید به‌صورت شفاف
- migrate شدن portal به contract داخلی `/internal/tools/*`
- محدودشدن routeهای legacy به compatibility backend

### Infra

- اضافه‌کردن envهای explicit برای Go
- health/readiness aggregation
- internal-only policy برای `voltagent-service`
- runbook و hardening deployment

### Docs

- نگه‌داشتن تمام day fileها داخل `docs/Daily task`
- sync کردن هر روز بعد از هر تغییر واقعی در repo
- جلوگیری از باقی‌ماندن bullets قدیمی در اسناد روزانه

### Validation

- unit test
- integration test
- e2e test
- smoke test compose

## قانون به‌روزرسانی

از این نقطه به بعد:

- هر کار جدید یا باقی‌مانده باید اول در همین index یا فایل روز مربوطه ثبت شود
- اگر یک روز کامل نشد، remaining taskها باید به روز بعد carry شوند
- اگر فایل روز جدید لازم باشد، باید داخل همین فولدر ساخته شود

## هدف نهایی

این migration وقتی done است که:

- `Go` در مسیر اصلی planning از `voltagent-service` استفاده کند
- fallback observable باشد
- BFF و portal به مسیر نهایی align شده باشند
- health و observability کامل شده باشند
- docs روزانه و docs معماری با repo sync باشند
- مسیر `BFF -> Go -> voltagent-service -> Go -> Temporal` تست‌پذیر و پایدار باشد
