# Daily Task 2: MVP Agent Path

این فایل مرجع اجرایی مرحله بعد از اتمام Phase 1 مربوط به `voltagent-service` است و نسخه canonical این سری باید در همین فولدر بماند:

- `docs/Daily task 2/`

هدف این سری جدید:

- رساندن استک به یک MVP واقعی
- بستن یک مسیر قابل مشاهده که در آن کاربر بتواند با یک agent ساده یک کار واقعی انجام دهد
- کاهش فاصله بین demo surface و runtime واقعی

مرجع‌های اصلی این مرحله:

- `docs/MVP_DEVELOPMENT_PLAN.md`
- `docs/Daily task/VOLTAGENT_SERVICE_DAILY_PLAN.md`
- `docs/Daily task/VOLTAGENT_SERVICE_DAY7_LEGACY_CLEANUP_AND_RELEASE.md`

## جمع‌بندی وضعیت فعلی

آنچه الان واقعاً در repo وجود دارد:

- Phase 1 مربوط به `voltagent-service` عملاً بسته شده و مسیرهای جدید `/internal/tools/*` تثبیت شده‌اند
- `portal1` برای `manifest`، `execute` و `deploy` به contract داخلی backend وصل شده است
- `portal1` در لایه UI هم `Builder`، `Workflows`، `Logs` و `GlobalChat` را دارد
- `BFF` برای chat از tool calling استفاده می‌کند و `deploy_website` را به Go می‌فرستد
- backend برای `deploy_website`، `workflow_id` و `planning_source` را برمی‌گرداند و در read model نگه می‌دارد

اما gapهای مهم برای MVP هنوز باز هستند:

- chat path هنوز contract کامل `deploy_website` را مثل builder حمل نمی‌کند و فقط `projectName` و `template` می‌فرستاد
- `portal1/src/lib/services/supernode.ts` هنوز در چند نقطه fallbackهای mock و simulation داشت
- بعضی surfaceها هنوز failure واقعی را نشان نمی‌دادند و در حالت خطا به رفتار demo برمی‌گشتند
- use case اصلی `deploy_website` از نظر orchestration جلو رفته، اما بخشی از activityها هنوز placeholder یا کم‌اعتماد هستند
- هنوز یک smoke path رسمی نداریم که بگوید: "از chat یا builder شروع شد، workflow ساخته شد، status دیده شد، artifact معنی‌دار برگشت"

## اصل اولویت

در این مرحله هر کاری که انجام می‌دهیم باید به این سوال پاسخ بدهد:

> آیا الان کاربر می‌تواند با یک agent ساده، یک درخواست واقعی بدهد و نتیجه آن را در UI ببیند؟

اگر پاسخ منفی باشد، آن کار هنوز جزو MVP ضروری نیست.

## تعریف Done این سری

این Daily Task 2 وقتی موفق محسوب می‌شود که این موارد سبز باشند:

- `portal1` از یک entry واقعی مثل chat یا builder بتواند یک tool واقعی را trigger کند
- `BFF` و `Go` contract یکسان و بدون drift برای همان use case داشته باشند
- failure به‌صورت واضح نمایش داده شود و با mock پنهان نشود
- workflow و log و artifact از UI قابل مشاهده باشند
- حداقل یک validation تکرارپذیر برای مسیر اصلی ثبت شده باشد

## اولویت‌های واقعی MVP

### Priority 0: بستن حلقه Chat Agent

هدف:

- کاربر از `GlobalChat` یک درخواست طبیعی بدهد
- `BFF` tool مناسب را صدا بزند
- backend workflow را شروع کند
- نتیجه و status در همان chat و surfaceهای workflow دیده شود

چرا اولویت اول است:

- چون این نزدیک‌ترین تعریف به "یک agent ساده که یک کار انجام می‌دهد" است
- الان UI chat و BFF هر دو وجود دارند، اما contract بینشان هنوز کامل و product-grade نشده است

کارهای ضروری:

- یکسان‌سازی payload `deploy_website` بین `GlobalChat`، `BFF` و `Go`
- عبور فیلدهای `prompt`, `framework`, `theme`, `template` از chat path
- استانداردسازی response برای tool result
- ثبت و نمایش errorهای واقعی بدون fallback نمایشی

### Priority 0: حذف Demo Mode پنهان

هدف:

- اگر backend یا workflow در دسترس نبود، UI باید failure واقعی را نشان بدهد
- mock فقط در صورت mode صریح dev/demo فعال باشد

چرا اولویت بالا است:

- تا وقتی UI در سکوت به `mock-wf-*` و simulation برمی‌گردد، نمی‌شود گفت MVP واقعاً کار می‌کند

کارهای ضروری:

- حذف fallback خاموش در `portal1/src/lib/services/supernode.ts`
- تبدیل mockها به feature flag یا `demo_mode`
- نمایش پیام خطای واضح در `Builder`, `Workflows`, `Logs`, `GlobalChat`

### Priority 1: واقعی‌کردن Artifact خروجی

هدف:

- نتیجه workflow فقط یک `workflow_id` نباشد
- حداقل یکی از این خروجی‌ها معنی‌دار و قابل استفاده باشد:
  - `repo_url`
  - `pr_url`
  - `liveUrl`
  - `previewUrl`

چرا مهم است:

- MVP بدون artifact قابل استفاده هنوز بیشتر شبیه orchestration demo است
- کاربر باید بفهمد خروجی agent دقیقاً چه چیزی بوده است

کارهای ضروری:

- نرمال‌سازی artifactها در backend read model
- نمایش artifact معتبر در UI
- مشخص‌کردن اینکه کدام activityها placeholder هستند و کدام‌ها آماده استفاده‌اند

### Priority 1: Smoke Validation رسمی مسیر اصلی

هدف:

- یک check تکرارپذیر وجود داشته باشد که کل مسیر MVP را تایید کند

کارهای ضروری:

- تعریف سناریوی رسمی `chat -> BFF -> Go -> workflow -> UI`
- ثبت commandها و expected output
- اگر لازم بود افزودن یک تست سبک e2e یا smoke script

### Priority 2: افزودن دومین کار ساده برای Agent

این مورد فقط وقتی شروع می‌شود که Priority 0 و 1 سبز شده باشند.

گزینه‌های ممکن:

- یک non-deploy tool ساده از manifest
- یک query/read-only tool با ریسک کمتر
- یا نسخه پایدارتر از `crypto_analysis` فقط در حد trigger و visibility

قانون:

- قبل از سبز شدن مسیر `deploy_website`، نباید scope را با tool دوم باز کنیم

## ترتیب اجرای روزها

### روز 1: Chat Agent Contract

وضعیت:

- `completed`

کارهای انجام‌شده:

- contract اولیه `BFF` برای deploy از حالت حداقلی خارج شد و حالا `prompt`, `framework`, `theme`, `template` را هم پشتیبانی می‌کند
- برای chat path، اگر `projectName` صریح نباشد، نام پروژه در BFF derive می‌شود
- fallbackهای silent در `portal1/src/lib/services/supernode.ts` به `VITE_ENABLE_DEMO_MODE=true` محدود شدند
- `GlobalChat` حالا نتیجه‌های error را به‌عنوان failure واقعی نمایش می‌دهد و `planning_source` را هم نشان می‌دهد

کارهای باز:

- validation واقعی مسیر end-to-end با runtime روشن
- تصمیم نهایی برای artifact اولیه قابل نمایش در chat
- ثبت smoke command رسمی بعد از اجرای دستی

مرجع:

- `docs/Daily task 2/MVP_DAILY_TASK_2_DAY1_CHAT_AGENT_CONTRACT.md`

### روز 2: حذف Mock و Failure Surface

وضعیت:

- `completed`

هدف:

- حذف fallbackهای خاموش و واقعی‌کردن تجربه خطا در فرانت

مرجع:

- `docs/Daily task 2/MVP_DAILY_TASK_2_DAY2_REMOVE_MOCKS_AND_FAILURE_SURFACE.md`

### روز 3: Workflow Artifact و Result Surface

وضعیت:

- `completed`

هدف:

- نشان‌دادن artifact معنی‌دار در `Builder`, `Workflows`, `Logs`, `GlobalChat`

مرجع:

- `docs/Daily task 2/MVP_DAILY_TASK_2_DAY3_WORKFLOW_ARTIFACTS_AND_RESULT_SURFACE.md`

### روز 4: Smoke Validation و MVP Gate

وضعیت:

- `completed`

هدف:

- تعریف Go/No-Go برای MVP use case اصلی

مرجع:

- `docs/Daily task 2/MVP_DAILY_TASK_2_DAY4_SMOKE_VALIDATION_AND_MVP_GATE.md`

### روز 5: Tool دوم یا Polish هدفمند

وضعیت:

- `completed`

هدف:

- فقط بعد از سبز شدن مسیر اصلی، یا polish کم‌ریسک انجام شود یا یک tool دوم read-only و کم‌ریسک اضافه شود

مرجع:

- `docs/Daily task 2/MVP_DAILY_TASK_2_DAY5_SECOND_TOOL_OR_TARGETED_POLISH.md`

### روز 6: Frontend Stability و MVP Hardening

وضعیت:

- `completed`

هدف:

- کاهش خطاهای فرانت که روی اعتماد به MVP و routeهای اصلی SPA اثر مستقیم دارند

کارهای انجام‌شده:

- نرمال‌سازی hash routeها در `portal1` یکپارچه شد تا `#/`, `#/dashboard`, `#/builder`, `#/workflows`, `#/logs` drift نداشته باشند
- `GlobalChat` حالا context صفحه واقعی را از hash route می‌گیرد و به `BFF` می‌فرستد؛ در نتیجه system prompt دیگر روی `"/"` ثابت نمی‌ماند
- `Topbar`, `Sidebar` و route اصلی SPA از helper مشترک برای page title و active state استفاده می‌کنند
- `portal1` بعد از hardening دوباره با `npm run check` و `npm run build` سبز شد

مرجع:

- `docs/Daily task 2/MVP_DAILY_TASK_2_DAY6_FRONTEND_STABILITY_AND_MVP_HARDENING.md`

### روز 7: Release Candidate و Handoff

وضعیت:

- `completed`

هدف:

- بستن release candidate برای MVP و ثبت handoff شفاف برای post-MVP

کارهای انجام‌شده:

- validation نهایی health برای `BFF` و `Go gateway` با runtime روشن انجام شد
- مسیر `chat -> BFF -> Go -> workflow -> UI` با اجرای مستقیم `POST /api/chat` و `POST /internal/tools/deploy` دوباره validate شد
- یک smoke script رسمی در `scripts/mvp_release_candidate_smoke.sh` ثبت شد تا health, deploy, chat trigger و failure sanity تکرارپذیر باشند
- sanity check برای `missing_fields` و `workflow_not_found` سبز شد، اما سناریوی `backend unavailable` روی runtime زنده عمداً destructive اجرا نشد
- نتیجه release candidate به‌صورت شفاف ثبت شد: مسیر trigger و visibility سبز است، اما چون artifact نهایی مثل `liveUrl`, `previewUrl` یا `repoUrl` در validation نهایی حاضر نبود، وضعیت فعلی برای release رسمی MVP هنوز `No-Go` است

مرجع:

- `docs/Daily task 2/MVP_DAILY_TASK_2_DAY7_RELEASE_CANDIDATE_AND_HANDOFF.md`

## قانون carry-over

اگر هر روز کامل نشد:

- taskهای باز باید explicit به روز بعد منتقل شوند
- دلیل باز ماندن task باید ثبت شود
- اگر blocker محیطی یا dependency بیرونی وجود داشت، باید همان روز نوشته شود

## خروجی نهایی مورد انتظار

در پایان این سری جدید، باید بتوان این جمله را با اطمینان گفت:

> در `portal1` یک agent ساده وجود دارد که از chat یا builder یک درخواست واقعی می‌گیرد، از مسیر اصلی استک اجرا می‌شود، workflow می‌سازد، و نتیجه‌اش را به‌صورت قابل مشاهده به کاربر برمی‌گرداند.
