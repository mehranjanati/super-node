# MVP Daily Task 3 - Day 6: Smoke Validation For Real Agent MVP

هدف روز:

- ثبت smoke path رسمی برای `design -> save -> ask -> result`
- تولید evidence تکرارپذیر برای این‌که agent MVP فقط در حد implementation پراکنده نمانده است

## چرا این روز حیاتی است

اگر Day 6 ثبت نشود، کل Daily Task 3 هنوز روی "فکر می‌کنیم کار می‌کند" می‌ماند.

برای MVP کاربردی، implementation کافی نیست. باید بتوان یک مسیر روشن و تکرارپذیر نشان داد که بگوید:

- draft ساخته شد
- draft ذخیره شد
- draft انتخاب شد
- agent واقعاً اجرا شد
- نتیجه meaningful برگشت
- failureهای پایه هم قابل‌فهم بودند

## وضعیت شروع روز

- `Foundry` و `Projects` draft واقعی را نگه می‌دارند
- contract بین draft و runtime path ثبت شده است
- `workflow_insight` در `BFF` و `GlobalChat` فعال است
- result surface اولیه وجود دارد
- هنوز smoke path رسمی این سری ثبت نشده است

## scope این روز

Day 6 فقط روی validation متمرکز است.

مجاز:

- ساخت اسکریپت smoke سبک
- تعریف checklist دقیق
- ثبت evidence و expected output
- sanity check برای failure modeهای اصلی

غیرمجاز:

- refactor سنگین
- افزودن capability جدید
- تغییر architecture
- بازکردن scope به deploy artifactهای خارج از use case اصلی این سری

## کارهای امروز

1. تعریف smoke path رسمی:
   - `Foundry -> Save Draft`
   - `Projects -> Reopen Draft`
   - `GlobalChat -> Ask Workflow Insight`
   - `Visible Result`
2. ساخت validation فنی سبک:
   - check فرانت
   - check تایپ BFF
   - request مستقیم به `BFF /api/chat`
3. ثبت validation برای failure modeها:
   - draft انتخاب‌نشده
   - workflow ناموجود
   - data unavailable
4. ثبت اسکریپت یا commandهای تکرارپذیر
5. ثبت نتیجه نهایی smoke

## smoke path رسمی

### مسیر اصلی

1. در `Foundry` یک draft با این مشخصات بساز:
   - `capability=workflow_insight`
   - `executionMode=read_only_workflow_insight`
   - `resultSurface=global_chat`
2. draft را ذخیره کن
3. از `Projects` همان draft را دوباره باز کن
4. `GlobalChat` را باز کن
5. یکی از queryهای زیر را بفرست:
   - `آخرین workflowها را خلاصه کن`
   - `کدام workflowها fail شده‌اند؟`
6. بررسی کن:
   - selected draft در چت دیده شود
   - tool صحیح invoke شود
   - summary متنی برگردد
   - structured result card نمایش داده شود

### sanity check برای failure

- سوال با `workflowId` ناموجود:
  - انتظار: `workflow_not_found`
- نبود داده:
  - انتظار: `no_workflows_available` یا `workflow_data_unavailable`
- selected draft نداشتن:
  - انتظار: fallback رفتاری قابل‌فهم و non-misleading

## validation فنی پیشنهادی

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

### 3. request مستقیم به BFF

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
            "text":"آخرین workflowها را خلاصه کن"
          }
        ]
      }
    ],
    "data":{
      "currentPath":"global_chat",
      "currentRoute":"#/dashboard",
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
      }
    }
  }'
```

انتظار:

- tool invocation از نوع `workflow_insight`
- خروجی شامل `summary`
- خروجی شامل `selected_workflows`

## فایل‌های محتمل

- `scripts/mvp_real_agent_smoke.sh`
- `portal1/src/lib/components/chat/GlobalChat.svelte`
- `BFF/index.ts`
- `docs/Daily task 3/MVP_DAILY_TASK_3_DAY6_SMOKE_VALIDATION_FOR_REAL_AGENT_MVP.md`

## معیار done

- smoke path رسمی این سری روی use case واقعی agent ثبت شده باشد
- commands و expected outputها قابل تکرار باشند
- failure modeهای پایه مستند شده باشند
- تیم بتواند بدون اتکا به توضیح شفاهی، مسیر `design -> save -> ask -> result` را validate کند

## carry-over

- اگر smoke سبز نشد، فقط blocker دقیق و evidence شکست به Day 7 منتقل می‌شود
- هر نوع گسترش scope یا feature work جدید در این روز failure محسوب می‌شود
