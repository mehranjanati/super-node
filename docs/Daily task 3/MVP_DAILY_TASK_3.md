# Daily Task 3: Real Agent MVP Path

این سری جدید بعد از `Daily task 2` شروع می‌شود و پاسخش به یک واقعیت ساده است:

> MVP فعلی هنوز بیش از حد deploy/demo-centric است و کاربر هنوز نمی‌تواند یک agent ساده را واقعاً طراحی، ذخیره، بازبینی و به اجرای قابل‌فهم نزدیک کند.

نسخه canonical این سری باید در همین فولدر بماند:

- `docs/Daily task 3/`

مرجع‌های اصلی:

- `docs/Daily task 2/MVP_DAILY_TASK_2.md`
- `docs/Daily task 2/MVP_DAILY_TASK_2_DAY7_RELEASE_CANDIDATE_AND_HANDOFF.md`
- `docs/dev/agents.md`

## چرا این سری لازم است

آنچه از `Daily task 2` باقی ماند:

- مسیر `chat -> BFF -> Go -> workflow -> UI` برای `deploy_website` وجود دارد
- health و trigger و read model تا حدی کار می‌کنند
- فرانت از نظر route و failure surface پایدارتر شده است

اما gapهای حیاتی هنوز باز هستند:

- `Foundry` و `Projects` تا همین لحظه بیشتر surface نمایشی بودند تا entry واقعی برای طراحی agent
- کاربر نمی‌توانست یک draft ساده از agent خودش را بسازد و بعد دوباره به آن برگردد
- use case اصلی هنوز بیش از حد روی deploy متمرکز است و هنوز به "agent that does real work" نرسیده
- بین draft طراحی agent و runtime واقعی هنوز contract روشن و کوتاه وجود ندارد
- هنوز یک مسیر product-grade نداریم که بگوید:
  - کاربر agent را تعریف کرد
  - config او ذخیره شد
  - agent به یک capability واقعی وصل شد
  - نتیجه meaningful در UI دیده شد

## اصل اولویت این سری

در این سری هر کاری باید به این سوال پاسخ بدهد:

> آیا کاربر الان می‌تواند یک agent ساده بسازد که فقط ظاهر نمایشی نداشته باشد و حداقل یک کار واقعی، قابل‌فهم و قابل‌دیدن انجام بدهد؟

اگر پاسخ منفی باشد، آن کار هنوز MVP واقعی نیست.

## تعریف Done این سری

این Daily Task 3 وقتی موفق است که این موارد سبز باشند:

- `Foundry` یا surface مشابه بتواند یک draft واقعی از agent را بسازد و نگه دارد
- draft agent به یک contract مشخص برای runtime یا tool execution وصل شود
- حداقل یک agent ساده، non-demo و کم‌ریسک، یک کار واقعی انجام دهد
- خروجی agent فقط `workflow_id` یا shell status نباشد و برای کاربر معنی‌دار باشد
- validation رسمی برای مسیر `design -> save -> run -> visible result` ثبت شده باشد

## اولویت‌های واقعی

### Priority 0: واقعی‌کردن Agent Design Surface

هدف:

- `Foundry` و `Projects` از حالت mock خارج شوند
- user بتواند draft agent را بسازد، ذخیره کند، دوباره باز کند و تغییر بدهد

چرا اولویت اول است:

- چون بدون این مرحله، حتی "طراحی یک agent ساده" هم در محصول واقعی نیست
- این نزدیک‌ترین gap بین UI موجود و product واقعی است

### Priority 0: بستن یک Use Case واقعی برای Agent

هدف:

- یک agent ساده، کم‌ریسک و قابل‌فهم انتخاب شود که کار واقعی انجام دهد

گزینه‌های مناسب:

- query/read-only روی logs یا workflows
- retrieval ساده از manifest و tool catalog
- summary یا analysis سبک روی داده‌های read-only موجود

قانون:

- تا وقتی agent دومن‌محور واقعی سبز نشده، scope نباید دوباره روی deploy demo باز شود

### Priority 1: Contract بین Draft و Runtime

هدف:

- draft agent در UI فقط یک فرم ذخیره‌شده نباشد
- بین config فرانت و contract اجرای backend/BFF یک mapping روشن تعریف شود

نمونه فیلدهای مورد نیاز:

- `name`
- `system_prompt`
- `targetAudience`
- `tools`
- `execution_mode`
- `result_surface`

### Priority 1: Result Surface واقعی

هدف:

- output agent در UI به‌صورت قابل‌مصرف دیده شود

نمونه outputهای قابل‌قبول:

- summary
- structured result
- retrieved records
- selected workflow/log insight

### Priority 2: Validation رسمی برای Agent MVP

هدف:

- یک smoke path رسمی برای `design -> save -> run -> result` تعریف شود

## لنگر اجرایی این سری

برای جلوگیری از بازگشت به scopeهای demo-centric، این سری فعلاً حول یک capability مشخص می‌چرخد:

- `Workflow Insight Agent`

برداشت عملی:

- معیار این سری دیگر `deploy artifact` نیست
- معیار این سری این است که draft agent انتخاب‌شده بتواند با contract روشن، روی data واقعی read-only اجرا شود و summary یا structured result قابل‌فهم برگرداند

## ترتیب روزها

### روز 1: Agent Draft Foundation

وضعیت:

- `completed`

هدف:

- تبدیل `Foundry` و `Projects` از surface نمایشی به draft flow واقعی

مرجع:

- `docs/Daily task 3/MVP_DAILY_TASK_3_DAY1_AGENT_DRAFT_FOUNDATION.md`

### روز 2: Real Agent Use Case Selection

وضعیت:

- `completed`

هدف:

- انتخاب یک کار واقعی، کم‌ریسک و non-demo برای agent اول

خروجی:

- `Workflow Insight Agent` به‌عنوان use case اول انتخاب شد
- دلیل انتخاب، input/output و contract اولیه در
  `docs/Daily task 3/MVP_DAILY_TASK_3_DAY2_REAL_AGENT_USE_CASE_SELECTION.md`
  ثبت شد

### روز 3: Draft To Runtime Contract

وضعیت:

- `completed`

هدف:

- تعریف mapping روشن بین draft agent و contract اجرا در BFF/Go

خروجی:

- contract صریح برای draft agent شامل `capability`, `executionMode`, `resultSurface` و `tools` تعریف شد
- selected draft context از فرانت به `GlobalChat` و `BFF` عبور داده شد

مرجع:

- `docs/Daily task 3/MVP_DAILY_TASK_3_DAY3_DRAFT_TO_RUNTIME_CONTRACT.md`

### روز 4: Read-Only Agent Execution

وضعیت:

- `completed`

هدف:

- اجرای اولین agent ساده با capability واقعی و safe

خروجی:

- tool جدید `workflow_insight` در `BFF` اضافه شد
- agent MVP حالا از endpointهای read-only موجود `Go Gateway` استفاده می‌کند
- `GlobalChat` summary و structured result مربوط به workflow insight را نمایش می‌دهد

مرجع:

- `docs/Daily task 3/MVP_DAILY_TASK_3_DAY4_READ_ONLY_AGENT_EXECUTION.md`

### روز 5: Result Surface And UX Hardening

وضعیت:

- `completed`

هدف:

- نمایش output معنادار agent و بستن failure surface مرتبط

مرجع:

- `docs/Daily task 3/MVP_DAILY_TASK_3_DAY5_RESULT_SURFACE_AND_UX_HARDENING.md`

### روز 6: Smoke Validation For Real Agent MVP

وضعیت:

- `pending`

هدف:

- ثبت smoke path رسمی برای `design -> save -> run -> result`

مرجع:

- `docs/Daily task 3/MVP_DAILY_TASK_3_DAY6_SMOKE_VALIDATION_FOR_REAL_AGENT_MVP.md`

### روز 7: MVP Re-Evaluation

وضعیت:

- `pending`

هدف:

- تصمیم دوباره برای Go/No-Go بر اساس agent MVP واقعی، نه فقط deploy orchestration

مرجع:

- `docs/Daily task 3/MVP_DAILY_TASK_3_DAY7_MVP_RE_EVALUATION.md`

### روز 8: Backend Orchestration Consolidation

وضعیت:

- `pending`

هدف:

- متمرکزکردن ownership orchestration capability فعلی در backend
- سبک‌کردن `BFF` و نزدیک‌کردن آن به نقش bridge قراردادمحور

خروجی مورد انتظار:

- endpoint روشن backend-facing برای agent execution
- انتقال orchestration اصلی `workflow_insight` از `BFF` به backend
- حفظ compatibility فرانت و جلوگیری از regression در مسیرهای موجود

مرجع:

- `docs/Daily task 3/MVP_DAILY_TASK_3_DAY8_BACKEND_ORCHESTRATION_CONSOLIDATION.md`

## خروجی مورد انتظار

در پایان این سری باید بتوان با اطمینان گفت:

> در `portal1` کاربر می‌تواند یک agent ساده را تعریف کند، آن را از حالت draft نمایشی خارج کند، به یک capability واقعی وصل کند و نتیجه قابل‌فهم آن را در UI ببیند.
