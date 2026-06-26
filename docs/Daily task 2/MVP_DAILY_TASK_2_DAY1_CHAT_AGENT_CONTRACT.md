# MVP Daily Task 2 - Day 1: Chat Agent Contract

هدف روز:

- تبدیل `GlobalChat` به اولین entry رسمی MVP
- بستن مسیر `chat -> BFF -> Go -> workflow -> visible status`
- حذف drift بین contract چت و contract واقعی `deploy_website`

## وضعیت شروع روز

موارد done قبل از شروع:

- `portal1` دارای `GlobalChat` با پشتیبانی از tool invocation و polling workflow است
- `BFF` دارای route `POST /api/chat` و tool calling برای `deploy_website` است
- `Go` دارای route `POST /internal/tools/deploy` و read model برای workflowها و logها است
- `voltagent-service` و fallback backend برای `deploy_website` در Phase 1 اعتبارسنجی شده‌اند

gapهای باز:

- chat path فقط `projectName` و `template` را حمل می‌کرد و از payload کامل deploy عقب بود
- mapping بین زبان طبیعی کاربر و args نهایی deploy هنوز محدود است
- UI chat در failure mode هنوز تعریف product-grade و بدون ambiguity ندارد
- success path از chat تا artifact نهایی هنوز رسمی و تکرارپذیر ثبت نشده است

## پیشرفت انجام‌شده

- schema ابزار `deploy_website` در `BFF/index.ts` گسترش یافت و حالا `prompt`, `projectName`, `framework`, `theme`, `template` را می‌پذیرد
- `BFF` حالا در صورت نبود `projectName`، یک slug پایدار از متن درخواست تولید می‌کند
- system prompt در `BFF` صریح‌تر شد تا مدل برای درخواست‌های ساخت/دیپلوی سایت از tool مناسب استفاده کند
- fallbackهای mock در `portal1/src/lib/services/supernode.ts` به mode صریح `VITE_ENABLE_DEMO_MODE=true` محدود شدند
- `GlobalChat` حالا tool resultهای خطادار را به‌عنوان failure نمایش می‌دهد و `planning_source` را هم کنار status نشان می‌دهد

## کارهای امروز

1. audit کردن contract چت:
   - `portal1/src/lib/components/chat/GlobalChat.svelte`
   - `BFF/index.ts`
   - `internal/adapters/gateway/echo.go`
2. یکسان‌سازی payload `deploy_website`:
   - `project_name`
   - `prompt`
   - `framework`
   - `theme`
   - `template`
3. شفاف‌کردن response tool result:
   - `workflow_id`
   - `planning_source`
   - `current_step`
   - `message`
   - artifactهای قابل نمایش
4. تعیین failure surface:
   - timeout در BFF
   - خطای backend
   - unreachable بودن workflow status
5. تعریف validation رسمی:
   - prompt ثابت برای trigger deploy
   - مشاهده tool invocation در chat
   - مشاهده workflow status در همان UI

## فایل‌های اصلی

- `portal1/src/lib/components/chat/GlobalChat.svelte`
- `portal1/src/lib/services/supernode.ts`
- `BFF/index.ts`
- `internal/adapters/gateway/echo.go`
- `internal/core/services/voltagent/voltagent.go`

## ریسک‌های روز

- mismatch بین schema مورد انتظار مدل در BFF و payload واقعی backend
- fallbackهای نمایشی که خطای runtime را پنهان می‌کنند
- نبود artifact معتبر در success path که باعث شود تجربه همچنان demo به نظر برسد
- تفاوت رفتار بین dev mode و production mode در مسیر `/api/chat`

## معیار done

- کاربر بتواند از chat یک درخواست deploy بدهد
- `BFF` payload کامل را به route جدید backend بفرستد
- `workflow_id` و status در chat قابل مشاهده باشند
- اگر backend fail شد، UI failure واقعی را نشان بدهد و به mock خاموش برنگردد
- دستورهای validation و expected result در همین فایل ثبت شده باشند

## وضعیت فعلی

- `completed`

کارهای انجام‌شده:
- validation دستی انجام شد و contractها به درستی کار می‌کنند.
- خطاها به‌طور صریح در UI نمایش داده می‌شوند و fallbacks مخفی حذف شده‌اند.
- payload بین chat, BFF, و Go یکسان‌سازی شد.

## Validation پیشنهادی

### تست 1: Health

```bash
curl http://localhost:3001/api/health
curl http://localhost:3000/internal/health
```

### تست 2: Deploy Contract

```bash
curl -X POST http://localhost:3000/internal/tools/deploy \
  -H 'Content-Type: application/json' \
  -d '{
    "project_name": "mvp-chat-demo",
    "prompt": "Create a simple landing page for an AI studio",
    "framework": "svelte",
    "theme": "modern",
    "template": "default"
  }'
```

### تست 3: Chat Trigger

prompt پیشنهادی در `GlobalChat`:

```text
برای من یک سایت ساده برای استودیو هوش مصنوعی بساز. اسم پروژه mvp-chat-demo باشد و با svelte اجرا شود.
```

خروجی مورد انتظار:

- tool invocation در chat دیده شود
- `workflow_id` برگردد
- status اولیه `RUNNING` یا `started` دیده شود
- در `#/workflows` همان رکورد وجود داشته باشد
- در `#/logs` حداقل log اولیه workflow دیده شود

## carry-over احتمالی

اگر امروز کامل نشد، فقط این موارد به روز بعد منتقل می‌شوند:

- حذف fallbackهای demo در surfaceهای باقیمانده
- ثبت artifact واقعی در chat result
- نوشتن smoke validation رسمی برای use case اصلی
