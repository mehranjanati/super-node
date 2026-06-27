# MVP Daily Task 3 - Day 8: Backend Orchestration Consolidation

هدف روز:

- انتقال orchestration مرتبط با agent MVP از `BFF` به لایه backend
- حفظ `BFF` به‌عنوان bridge قراردادمحور، نه host اصلی orchestration

## چرا این روز حالا لازم است

تا پایان Day 7، اگر مسیر `design -> save -> ask -> visible result` برای `Workflow Insight Agent` سبز شده باشد، یک سوال معماری مهم باقی می‌ماند:

> owner واقعی orchestration این مسیر کدام لایه است؟

اگر orchestration بین `BFF` و `Go Gateway` پخش بماند:

- منطق runtime به‌مرور در دو لایه تکثیر می‌شود
- contract draft-to-runtime در مرزهای مختلف تفسیر متفاوت پیدا می‌کند
- hardening بعدی برای capabilityهای جدید پرهزینه‌تر می‌شود
- تصمیم معماری ثبت‌شده پروژه که `Go` را owner orchestration می‌داند، ضعیف می‌شود

این روز قرار نیست use case جدید بسازد؛ فقط باید ownership orchestration را روشن و متمرکز کند.

## وضعیت شروع روز

- `Workflow Insight Agent` به‌عنوان capability اصلی MVP ثبت شده است
- `GlobalChat` selected draft را به `BFF` می‌فرستد
- `BFF` فعلاً هم routing prompt و هم orchestration ابزار `workflow_insight` را نگه می‌دارد
- `Go Gateway` read modelهای لازم مثل `GET /workflows` و `GET /logs` را دارد
- تصمیم‌های معماری پروژه ترجیح می‌دهند orchestration در `Go` متمرکز بماند

## اصل تصمیم

در این روز:

- `Go Gateway` باید owner منطق اجرای capability `workflow_insight` شود
- `BFF` باید request contract-aware را forward کند و از نگه‌داشتن orchestration domain-specific سبک‌تر شود

این روز مجاز است:

- endpoint جدید backend برای execution سطح agent اضافه کند
- منطق retrieval/filtering/summarization مرتبط با `workflow_insight` را به backend منتقل کند
- response ساخت‌یافته فعلی را حفظ یا رسمی‌تر کند
- mapping capability و execution mode را به contract روشن‌تری نزدیک کند

این روز مجاز نیست:

- capability جدید اضافه کند
- multi-agent orchestration بسازد
- stateful planner جدید یا supervisor جدید وارد کند
- UX جدید خارج از نیازهای حفظ compatibility بسازد

## کارهای امروز

1. تعریف contract جدید backend-facing:
   - ورودی:
     - `question`
     - `selectedAgent`
     - `currentPath`
     - `currentRoute`
   - خروجی:
     - `status`
     - `capability`
     - `summary`
     - `selected_workflows`
     - `filters`
     - `error`
2. افزودن endpoint جدید در `Go Gateway` برای agent execution:
   - مثال: `POST /api/agents/execute`
   - یا مسیر داخلی معادل که ownership orchestration را روشن کند
3. انتقال منطق `workflow_insight` از `BFF` به backend:
   - retrieval از `workflows`
   - retrieval از `logs`
   - filtering
   - summarization
   - failure modeling
4. سبک‌کردن `BFF`:
   - حذف منطق domain-heavy مربوط به `workflow_insight`
   - تبدیل آن به adapter برای forwarding request و streaming پاسخ
5. حفظ compatibility:
   - `GlobalChat` نباید نیاز به contract جدیدی در فرانت پیدا کند
   - مسیر `deploy_website` نباید regress شود

## قرارداد هدف این روز

### request از `BFF` به backend

```json
{
  "question": "کدام workflow fail شده است؟",
  "selectedAgent": {
    "id": "agent-1",
    "name": "workflow-insight-agent",
    "type": "analytics",
    "config": {
      "capability": "workflow_insight",
      "executionMode": "read_only_workflow_insight",
      "resultSurface": "global_chat",
      "tools": ["workflow_insight"],
      "systemPrompt": "Summarize workflow state for the user"
    }
  },
  "currentPath": "global_chat",
  "currentRoute": "#/dashboard"
}
```

### response از backend

```json
{
  "status": "success",
  "capability": "workflow_insight",
  "summary": "Found 1 failed workflow. Latest signal is available.",
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
  "filters": {
    "status": "FAILED",
    "limit": 3
  }
}
```

## فایل‌های محتمل

- `internal/adapters/gateway/echo.go`
- `internal/...` سرویس یا use case جدید مرتبط با agent execution
- `BFF/index.ts`
- `portal1/src/lib/components/chat/GlobalChat.svelte`
- `docs/Daily task 3/MVP_DAILY_TASK_3_DAY8_BACKEND_ORCHESTRATION_CONSOLIDATION.md`

## وضعیت نهایی (Status)
**وضعیت:** انجام شد (Done) - موفقیت‌آمیز ✅
- مالکیت Orchestration به درستی به Go Gateway در قالب مسیر `POST /api/agents/execute` منتقل شد.
- توابع سنگین از `BFF/index.ts` حذف شدند و BFF اکنون تنها به عنوان یک Bridge با آگاهی از قرارداد (Contract) عمل می‌کند.
- اسکریپت Smoke Test (`mvp_real_agent_smoke.sh`) روی معماری جدید با موفقیت پاس شد.

## validation پیشنهادی

### 1. فرانت

```bash
cd /Users/elbaan/Documents/super\ node\ 1/portal1
npm run check
```

### 2. BFF

```bash
cd /Users/elbaan/Documents/super\ node\ 1/BFF
bunx tsc --noEmit index.ts --lib es2022,dom --module esnext --moduleResolution bundler --target es2022 --types bun
```

### 3. backend

```bash
cd /Users/elbaan/Documents/super\ node\ 1
go test ./...
```

### 4. request مستقیم

```bash
curl -X POST http://localhost:3000/api/agents/execute \
  -H 'Content-Type: application/json' \
  -d '{
    "question":"آخرین workflowها را خلاصه کن",
    "selectedAgent":{
      "id":"agent-smoke",
      "name":"workflow-insight-agent",
      "type":"analytics",
      "config":{
        "capability":"workflow_insight",
        "executionMode":"read_only_workflow_insight",
        "resultSurface":"global_chat",
        "tools":["workflow_insight"],
        "systemPrompt":"Summarize workflow state for the user"
      }
    },
    "currentPath":"global_chat",
    "currentRoute":"#/dashboard"
  }'
```

انتظار:

- ownership orchestration روی backend روشن باشد
- response ساخت‌یافته `workflow_insight` از backend برگردد
- `BFF` فقط نقش bridge/transport را نگه دارد

## معیار done

- `workflow_insight` دیگر orchestration اصلی خود را در `BFF` نگه ندارد
- backend owner روشن execution path برای capability فعلی باشد
- contract فرانت تغییر شکسته ایجاد نکند
- مسیر `deploy_website` regress نکند
- توزیع مسئولیت بین `portal1`, `BFF` و `Go` از نظر معماری واضح‌تر از قبل شود

## carry-over

- اگر این جابه‌جایی برای MVP فعلی زیادی پرریسک بود، فقط contract endpoint و plan مهاجرت ثبت می‌شود
- هر نوع گسترش scope به capability دوم یا orchestration چندمرحله‌ای failure این روز محسوب می‌شود
