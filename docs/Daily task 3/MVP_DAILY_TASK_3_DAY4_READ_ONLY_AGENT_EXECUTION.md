# MVP Daily Task 3 - Day 4: Read-Only Agent Execution

هدف روز:

- اجرای اولین agent ساده با capability واقعی و safe
- عبور از contract صرف به execution قابل‌مصرف در UI

## چرا این روز حیاتی است

تا پایان Day 3، draft agent دیگر فقط یک فرم نبود و contract روشن‌تری پیدا کرد.

اما هنوز سوال اصلی Daily Task 3 باز می‌ماند:

> آیا کاربر واقعاً می‌تواند یک agent ساده را اجرا کند و نتیجه قابل‌فهم ببیند؟

اگر Day 4 سبز نشود:

- تمام پیشرفت Day 1 تا Day 3 هنوز در مرز design باقی می‌ماند
- MVP همچنان از نظر کاربر نهایی یک surface زیبا ولی بدون action واقعی خواهد بود
- risk بازگشت به use caseهای deploy-centric دوباره بالا می‌رود

## وضعیت شروع روز

- `Workflow Insight Agent` در Day 2 انتخاب شده است
- contract draft به runtime در Day 3 روشن شده است
- endpointهای read-only موجود در `Go Gateway` شامل:
  - `GET /workflows`
  - `GET /workflows/:id`
  - `GET /logs`
- `GlobalChat` هنوز باید result را به شکلی product-usable نمایش دهد

## کارهای امروز

1. افزودن tool جدید در `BFF`:
   - `workflow_insight`
2. استفاده از endpointهای read-only موجود در `Go Gateway`
3. تبدیل data خام workflow/log به output قابل‌مصرف:
   - summary
   - selected workflows
   - latest log
   - status و step
4. نمایش result ساخت‌یافته در `GlobalChat`
5. حفظ مسیر `deploy_website` بدون regression، ولی خارج از use case اصلی MVP

## وضعیت فعلی

- `completed`

## کارهای انجام‌شده

- در `BFF/index.ts` tool جدید `workflow_insight` اضافه شد
- این tool query کاربر را به retrieval ساده از `workflows` و `logs` تبدیل می‌کند
- `BFF` از endpointهای موجود `Go Gateway` استفاده می‌کند و orchestration جدیدی اضافه نشده است
- response این tool فقط یک message آزاد نیست و payload ساخت‌یافته شامل این فیلدها را برمی‌گرداند:
  - `status`
  - `capability`
  - `summary`
  - `selected_workflows`
  - `total_workflows`
  - `filters`
  - `error`
- `GlobalChat` selected draft را همراه request می‌فرستد و حالا اگر نتیجه از نوع `workflow_insight` باشد:
  - summary را نمایش می‌دهد
  - workflowهای منتخب را card-style نشان می‌دهد
  - `status`, `currentStep`, `planningSource`, `lastLog`, `updatedAt` را render می‌کند
- initial assistant copy در چت از deploy-centric بودن فاصله گرفت و به workflow/runtime insight نزدیک شد
- failure modeهای read-only به‌صورت روشن‌تر model شدند:
  - `question_required`
  - `workflow_not_found`
  - `no_workflows_available`
  - `workflow_data_unavailable`

## contract اجرای این روز

### input

- `question`
- `workflowId` اختیاری
- `status` اختیاری
- `limit` اختیاری
- `selectedAgent`
- `currentPath`
- `currentRoute`

### output

```json
{
  "status": "success",
  "capability": "workflow_insight",
  "summary": "Found 2 workflows. Running: 1, completed: 0, failed: 1.",
  "selected_workflows": [
    {
      "workflowId": "wf-123",
      "name": "deploy-site-demo",
      "status": "FAILED",
      "currentStep": "INIT",
      "planning_source": "remote_voltagent",
      "last_log": "Planner returned no executable artifact.",
      "updated_at": "2026-06-26T17:00:00Z",
      "matched_by": "status"
    }
  ],
  "total_workflows": 1,
  "filters": {
    "status": "FAILED",
    "limit": 3
  }
}
```

## فایل‌های اصلی

- `BFF/index.ts`
- `portal1/src/lib/components/chat/GlobalChat.svelte`
- `internal/adapters/gateway/echo.go` (به‌عنوان provider endpointهای موجود)

## validation

```bash
cd /Users/elbaan/Documents/super\ node\ 1/portal1
npm run check

cd /Users/elbaan/Documents/super\ node\ 1/BFF
bunx tsc --noEmit index.ts --lib es2022,dom --module esnext --moduleResolution bundler --target es2022 --types bun
```

خروجی مورد انتظار:

- tool `workflow_insight` بدون خطای type اجرا شود
- `GlobalChat` نتیجه را به‌صورت summary و structured result نمایش دهد
- مسیر قبلی `deploy_website` از کار نیفتد

## limitationهای فعلی

- execution path هنوز درون `GlobalChat` embed شده و surface اختصاصی result ندارد
- validation رسمی end-to-end برای مسیر `design -> save -> ask -> result` هنوز مستند نشده است
- UX خطاها و حالت empty هنوز نیاز به hardening دارد
- `Projects` هنوز contract/capability را روی کارت‌ها به‌وضوح نشان نمی‌دهد

## carry-over به روز بعد

- hardening result surface
- شفاف‌تر کردن failure states و empty states
- کم‌کردن ambiguity تجربه کاربر در `Foundry`, `Projects` و `GlobalChat`
