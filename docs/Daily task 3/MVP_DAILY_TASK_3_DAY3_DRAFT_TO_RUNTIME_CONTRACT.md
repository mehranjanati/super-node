# MVP Daily Task 3 - Day 3: Draft To Runtime Contract

هدف روز:

- بستن mapping روشن بین draft agent در فرانت و contract اجرای واقعی
- جلوگیری از این‌که draft فقط یک فرم ذخیره‌شده بماند

## چرا این روز حیاتی است

بعد از Day 2 دیگر ابهامی درباره use case اصلی وجود ندارد: agent اول باید `Workflow Insight Agent` باشد.

اما تا وقتی بین draft ذخیره‌شده در `Foundry` و مسیر اجرای واقعی در `BFF/Go` یک contract کوتاه، صریح و قابل‌حمل بسته نشود:

- فرانت فقط state نگه می‌دارد اما execution معنی‌دار ندارد
- `GlobalChat` نمی‌فهمد کدام draft یا capability باید فعال شود
- `BFF` مجبور می‌شود فقط از روی متن آزاد کاربر intent را حدس بزند
- MVP دوباره به سمت تجربه مبهم و demo-like برمی‌گردد

## وضعیت شروع روز

- `Foundry` و `Projects` draft واقعی agent را نگه می‌دارند
- use case نهایی Day 2 روی `Workflow Insight Agent` بسته شده است
- `GlobalChat` هنوز context انتخاب draft را به backend منتقل نمی‌کند
- shape مشخصی برای capability، execution mode و result surface در draft agent وجود ندارد

## کارهای امروز

1. صریح‌کردن contract سمت فرانت:
   - تعریف typeهای capability و execution mode
   - جلوگیری از `config` مبهم و بدون shape مشخص
2. افزودن فیلدهای runtime-facing به `Foundry`:
   - `capability`
   - `executionMode`
   - `resultSurface`
   - `tools`
3. نمایش contract در UI:
   - preview JSON
   - summary قابل‌دیدن در خود `Foundry`
4. عبور context draft انتخاب‌شده به `GlobalChat`
5. آماده‌کردن `BFF` برای دریافت selected draft context

## وضعیت فعلی

- `completed`

## کارهای انجام‌شده

- در `portal1/src/lib/types/index.ts` typeهای زیر اضافه شدند:
  - `AgentCapability`
  - `AgentExecutionMode`
  - `AgentResultSurface`
  - `AgentConfig`
- `Agent.config` از یک `Record<string, unknown>` مبهم به contract صریح‌تر `AgentConfig` نزدیک شد
- `Foundry` حالا برای هر draft این فیلدها را نگه می‌دارد:
  - `capability`
  - `executionMode`
  - `resultSurface`
- logic مشتق‌سازی execution mode از capability تعریف شد تا contract اولیه برای agentهای read-only ساده و deterministic بماند
- JSON preview در `Foundry` حالا contract واقعی draft را نشان می‌دهد، نه فقط metadata نمایشی
- در خود `Foundry` یک summary مجزا برای execution contract اضافه شد تا کاربر ببیند draft او قرار است با چه mode و چه toolی اجرا شود
- `GlobalChat` selected draft را از store می‌خواند و همراه request به `BFF` می‌فرستد
- `BFF` حالا selected draft context را در prompt سیستمی دریافت می‌کند تا routing و تصمیم‌گیری روی capability واقعی استوار شود

## contract نهایی این روز

### input از فرانت

- `selectedAgent.id`
- `selectedAgent.name`
- `selectedAgent.type`
- `selectedAgent.config.systemPrompt`
- `selectedAgent.config.capability`
- `selectedAgent.config.executionMode`
- `selectedAgent.config.resultSurface`
- `selectedAgent.config.tools`
- `currentPath`
- `currentRoute`

### shape حداقلی draft config

```json
{
  "systemPrompt": "Summarize workflow state for the user",
  "targetAudience": "Ops / Founders",
  "language": "fa",
  "framework": "svelte",
  "runtime": "go",
  "capability": "workflow_insight",
  "executionMode": "read_only_workflow_insight",
  "resultSurface": "global_chat",
  "tools": ["workflow_insight"],
  "source": "foundry_draft"
}
```

### mapping اصلی

- `capability=workflow_insight`
  - `executionMode=read_only_workflow_insight`
  - `tools=["workflow_insight"]`
  - `resultSurface=global_chat`
- `capability=deploy_website`
  - فعلاً فقط برای compatibility نگه داشته می‌شود
  - use case اصلی MVP این سری نیست

## فایل‌های اصلی

- `portal1/src/lib/types/index.ts`
- `portal1/src/lib/components/foundry/Foundry.svelte`
- `portal1/src/lib/components/chat/GlobalChat.svelte`
- `BFF/index.ts`

## validation

```bash
cd /Users/elbaan/Documents/super\ node\ 1/portal1
npm run check

cd /Users/elbaan/Documents/super\ node\ 1/BFF
bunx tsc --noEmit index.ts --lib es2022,dom --module esnext --moduleResolution bundler --target es2022 --types bun
```

خروجی مورد انتظار:

- contract typeها بدون خطای TypeScript ثبت شوند
- `Foundry` فیلدهای capability و execution را ذخیره کند
- `GlobalChat` selected draft را همراه request بفرستد

## limitationهای فعلی

- persistence draft همچنان local-first است و backend-synced نیست
- contract فعلاً بیشتر برای `GlobalChat` مصرف می‌شود و هنوز surfaceهای دیگر consumer کامل آن نیستند
- این روز contract را واقعی می‌کند، اما هنوز execution path کامل و result نهایی را به‌تنهایی نمی‌بندد

## carry-over به روز بعد

- اجرای واقعی capability `workflow_insight`
- اتصال `BFF` به read-only endpointهای موجود `Go Gateway`
- نمایش result ساخت‌یافته و معنادار در UI
