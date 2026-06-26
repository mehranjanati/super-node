# MVP Daily Task 2 - Day 7: Release Candidate And Handoff

هدف روز:

- بستن نسخه candidate برای MVP
- ثبت handoff روشن برای ادامه توسعه بعد از MVP

## وضعیت شروع روز

- مسیر اصلی `deploy_website` باید سبز و قابل تکرار باشد
- فرانت SPA برای routeهای اصلی باید پایدار شده باشد
- اگر tool دوم اضافه شده، باید کم‌ریسک و بدون regression باشد

## کارهای امروز

1. اجرای validation نهایی برای use case اصلی:
   - chat -> BFF -> Go -> workflow -> UI
2. اجرای validation نهایی برای artifact:
   - `workflow_id`
   - `planning_source`
   - `liveUrl` یا `previewUrl` یا `repoUrl`
3. اجرای sanity check برای failure mode:
   - backend unavailable
   - workflow status unavailable
   - tool invocation failure
4. ثبت checklist release candidate:
   - چه چیزهایی سبز هستند
   - چه چیزهایی عمداً defer شده‌اند
   - چه ریسک‌هایی پذیرفته شده‌اند
5. ثبت handoff برای بعد از MVP:
   - backlog کوتاه
   - debtهای باقی‌مانده
   - اولویت‌های post-MVP

## وضعیت فعلی

- `completed`

## validation نهایی انجام‌شده

### 1. health

```bash
curl http://localhost:3001/api/health
curl http://localhost:3000/internal/health
```

نتیجه:

- `BFF` روی `http://localhost:3001` سالم بود
- `Go gateway` و dependency مربوط به `voltagent-service` روی `http://localhost:3000` سالم بودند

### 2. trigger مستقیم use case اصلی

```bash
curl -X POST http://localhost:3000/internal/tools/deploy \
  -H 'Content-Type: application/json' \
  -d '{
    "project_name":"mvp-release-rc",
    "prompt":"Create a simple landing page for an AI studio with hero, pricing, and CTA sections.",
    "framework":"svelte",
    "theme":"modern",
    "template":"default"
  }'
```

نتیجه:

- backend پاسخ `status=started` برگرداند
- `workflow_id` و `planning_source=remote_voltagent` ثبت شدند

### 3. trigger از مسیر chat -> BFF -> Go

```bash
curl -N -X POST http://localhost:3001/api/chat \
  -H 'Content-Type: application/json' \
  -d '{
    "messages":[
      {
        "id":"msg-1",
        "role":"user",
        "parts":[
          {
            "type":"text",
            "text":"برای من یک سایت ساده برای استودیو هوش مصنوعی بساز. اسم پروژه mvp-release-chat باشد و با svelte اجرا شود."
          }
        ]
      }
    ],
    "data":{
      "currentPath":"builder",
      "currentRoute":"#/builder"
    }
  }'
```

نتیجه:

- `BFF` tool invocation از نوع `deploy_website` را استریم کرد
- خروجی tool شامل `workflow_id`, `planning_source`, `current_step`, `status` بود
- workflow از مسیر چت با موفقیت در read model قابل مشاهده شد

### 4. workflow و logs

workflowهای تازه‌ساخته‌شده:

- `deploy-site-mvp-release-rc-1782481385807`
- `deploy-site-mvp-release-chat-1782481389197`
- اجرای smoke script رسمی:
  - `scripts/mvp_release_candidate_smoke.sh`

نتیجه:

- `GET /workflows/:id` برای هر دو workflow پاسخ معتبر برگرداند
- `GET /logs` entryهای مرتبط با workflowهای جدید را برگرداند
- status فعلی در runtime موجود همچنان `RUNNING / INIT` بود
- artifactهای موجود فعلاً شامل `project_name`, `template`, `planning_source`, `message` هستند و در این اجرا `liveUrl` یا `previewUrl` یا `repoUrl` برنگشت

### 5. sanity check برای failure mode

```bash
curl -i http://localhost:3000/workflows/workflow-does-not-exist
curl -i -X POST http://localhost:3000/internal/tools/deploy \
  -H 'Content-Type: application/json' \
  -d '{}'
```

نتیجه:

- workflow ناموجود با `404 workflow_not_found` شفاف برگشت
- tool invocation نامعتبر با `400 missing_fields` شفاف برگشت
- failure mode مربوط به `backend unavailable` عمداً روی runtime زنده اجرا نشد تا سرویس‌های موجود مختل نشوند؛ این مورد در UI از نظر کد review شد ولی در این دور به‌صورت destructive تست نشد

## Release Candidate Checklist

### سبز

- trigger اصلی از `chat` و `Go direct deploy` کار می‌کند
- `workflow_id` و `planning_source` از contract اصلی برمی‌گردند
- `workflows` و `logs` در read model قابل مشاهده‌اند
- failureهای پایه مثل payload نامعتبر و workflow ناموجود شفاف هستند
- smoke script رسمی برای تکرار validation اضافه شد:
  - `scripts/mvp_release_candidate_smoke.sh`

### defer شده

- حذف کامل legacyهای غیر MVP مثل `main-router.svelte`
- productization کامل routeها و surfaceهای غیر MVP
- browser-level chaos check برای سناریوی `backend unavailable`

### ریسک‌های پذیرفته‌شده

- latest workflowها در runtime فعلی هنوز از `INIT` جلو نرفته‌اند و artifact URL معنادار برنگردانده‌اند
- نتیجه فعلی برای validation داخلی و ادامه توسعه مفید است، اما هنوز برای ادعای MVP رسمیِ کامل بدون ambiguity کافی نیست

## نتیجه Go/No-Go

- `No-Go` برای release رسمی MVP

دلیل:

- معیار Day 7 می‌خواست حداقل یک artifact معنی‌دار مثل `liveUrl`, `previewUrl` یا `repoUrl` دیده شود، اما در validation فعلی این خروجی‌ها حاضر نبودند
- workflowها ساخته می‌شوند و surfaceها قابل مشاهده‌اند، ولی runtime هنوز در همین اجرای نهایی evidence کافی برای "نتیجه نهایی قابل استفاده" نداده است

برداشت عملی:

- `Go` برای validation داخلی، demo فنی و ادامه hardening
- `No-Go` برای اعلام MVP رسمی نهایی تا وقتی artifact قابل مصرف و progression فراتر از `INIT` دیده نشود

## handoff بعد از MVP

backlog کوتاه:

- ریشه‌یابی این‌که چرا workflowهای جدید در read model روی `INIT` می‌مانند
- عبور artifactهای نهایی مثل `previewUrl`, `liveUrl` یا `repoUrl` از backend تا UI
- ثبت یک browser-level smoke یا e2e سبک برای `chat -> BFF -> Go -> workflow -> UI`

debtهای باقی‌مانده:

- cleanup فایل‌ها و routeهای legacy که دیگر consumer اصلی ندارند
- کاهش تکیه برخی surfaceها به polling و بهبود observability برای progression workflow

اولویت‌های post-MVP:

- سبز کردن artifact نهایی و status progression
- بستن sanity check واقعی برای `backend unavailable` در محیط کنترل‌شده
- تصمیم‌گیری دوباره درباره release رسمی بعد از یک validation تکرارپذیرِ کاملاً سبز

## Go/No-Go پیشنهادی

### Go

- deploy از chat یا builder قابل trigger باشد
- workflow در UI دیده شود
- حداقل یک artifact معنی‌دار دیده شود
- خطاهای اصلی به‌صورت واضح نمایش داده شوند
- smoke validation رسمی بدون ambiguity اجرا شود

### No-Go

- trigger اصلی فقط در demo mode کار کند
- artifactها fake یا غیرقابل تشخیص باشند
- routeهای اصلی SPA ناپایدار باشند
- failure surface هنوز misleading باشد

## معیار done

- تیم بتواند بگوید MVP دقیقاً در چه وضعیتی release-candidate شده است
- handoff برای کارهای بعد از MVP مستند و کوتاه باشد

## خروجی مورد انتظار

در پایان روز 7 یک پاسخ شفاف وجود داشته باشد:

> آیا این نسخه برای نمایش، validation داخلی یا ادامه توسعه به‌عنوان MVP رسمی قابل قبول است یا نه؟
