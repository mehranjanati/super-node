# MVP Daily Task 3 - Day 2: Real Agent Use Case Selection

هدف روز:

- انتخاب اولین use case واقعی برای agent MVP
- جلوگیری از بازگشت به scopeهای demo-centric یا orchestration-heavy

## چرا این روز حیاتی است

اگر `Daily task 3` بخواهد فقط روی "draft UI" بماند، دوباره به یک MVP نمایشی تبدیل می‌شود.

برای عبور از این وضعیت باید خیلی زود، دقیق و صریح مشخص شود که:

- agent اول دقیقاً چه کاری انجام می‌دهد
- این کار چقدر واقعی است
- ریسک آن برای MVP چقدر است
- خروجی آن در UI چگونه دیده می‌شود

تا وقتی این انتخاب روشن نشود:

- contract روزهای بعد مبهم می‌ماند
- فرانت معلوم نیست برای چه result surface باید آماده شود
- backend/BFF ممکن است دوباره به سمت use caseهای مبهم و پرریسک بروند

## وضعیت شروع روز

- `Foundry` و `Projects` حالا draft واقعی agent را نگه می‌دارند
- هنوز هیچ agentی به runtime واقعی وصل نشده است
- مسیر `deploy_website` وجود دارد، اما برای "agent واقعی و ساده" هنوز بیش از حد سنگین و misleading است
- استک read-only dataهای مفیدی مثل `manifest`, `workflows`, `logs` را در اختیار دارد

## اصل تصمیم

use case روز 2 فقط وقتی مجاز است که این شرایط را داشته باشد:

- read-only یا کم‌ریسک باشد
- نیازمند orchestration پیچیده یا multi-step deploy نباشد
- failure آن واضح و قابل‌نمایش باشد
- نتیجه آن برای کاربر بلافاصله قابل‌فهم باشد
- بتوان آن را از draft agent به runtime با contract محدود و روشن وصل کرد

## گزینه‌های مجاز

### گزینه A: Workflow Insight Agent

تعریف:

- agent فهرست workflowها را می‌خواند
- یک workflow را بر اساس ورودی کاربر انتخاب یا summarize می‌کند
- وضعیت، `planning_source`، step و آخرین log مهم را به کاربر نشان می‌دهد

مزیت‌ها:

- از data آماده backend استفاده می‌کند
- output آن textual/structured و کم‌هزینه است
- به surfaceهای موجود `Workflows` و `Logs` نزدیک است

ریسک‌ها:

- اگر workflow data هنوز ناقص یا stuck باشد، output ممکن است کم‌ارزش شود

### گزینه B: Tool Catalog / Manifest Agent

تعریف:

- agent manifest داخلی را می‌خواند
- toolهای موجود را توضیح می‌دهد
- برای use case کاربر مناسب‌ترین tool را پیشنهاد می‌دهد

مزیت‌ها:

- از contractهای فعلی استفاده می‌کند
- failure و output آن روشن است
- برای onboarding agentic platform ارزش دارد

ریسک‌ها:

- اگر فقط به "tool listing" محدود بماند، ممکن است همچنان کمی نمایشی به نظر برسد

### گزینه C: Log Summary Agent

تعریف:

- agent snapshot لاگ‌ها را می‌خواند
- خطاها، warningها یا آخرین وضعیت runtime را خلاصه می‌کند

مزیت‌ها:

- value عملی برای debugging و عملیات دارد
- read-only و low-risk است

ریسک‌ها:

- کیفیت خروجی به غنای لاگ‌ها وابسته است
- ممکن است برای کاربر product-facing کمتر از حد لازم جذاب باشد

## گزینه‌های غیرمجاز در این روز

- agentهای multi-step با تصمیم‌گیری پیچیده
- agentهایی که بدون guardrail عمل write/delete انجام می‌دهند
- agentهایی که دوباره روی deploy-centric orchestration بنا می‌شوند
- scopeهایی که نیازمند auth، billing، multi-user state یا external integration جدید هستند

## خروجی مورد انتظار از این روز

در پایان روز 2 باید این موارد ثبت شده باشند:

1. use case نهایی انتخاب‌شده
2. دلیل انتخاب
3. use caseهایی که عمداً رد شده‌اند
4. input/output دقیق agent
5. surfaceهای فرانت درگیر
6. endpointها و contractهای backend/BFF مورد نیاز

## وضعیت فعلی

- `completed`

## قالب تصمیم نهایی

در پایان روز باید این بلوک پر شود:

### انتخاب نهایی

- `Workflow Insight Agent`

### چرا انتخاب شد

- از endpointها و data modelهای موجود `workflows` و `logs` استفاده می‌کند و نیاز به integration جدید ندارد
- نسبت به `deploy_website` یک capability واقعی اما read-only و کم‌ریسک می‌دهد
- خروجی آن می‌تواند بلافاصله برای کاربر قابل‌فهم باشد: وضعیت workflow، step فعلی، planning source و آخرین log مهم
- به surfaceهای آماده‌ی `Workflows` و `Logs` نزدیک است و می‌تواند بدون orchestration جدید به `GlobalChat` وصل شود
- با وضعیت فعلی MVP که بعضی workflowها در `INIT` می‌مانند هم سازگار است، چون همین گیرکردن را به insight قابل‌مصرف تبدیل می‌کند

### inputهای اولیه

- `agent_id` یا draft انتخاب‌شده از `Foundry/Projects`
- `question` یا intent کاربر مثل:
  - `آخرین workflowهای من را خلاصه کن`
  - `بگو کدام workflow fail شده`
  - `وضعیت workflow <id> چیست؟`
- فیلتر اختیاری:
  - `workflow_id`
  - `status`
  - `limit`
- context مسیر فعلی فرانت برای کمک به result surface:
  - `currentPath`
  - `currentRoute`

### outputهای مورد انتظار

- summary متنی کوتاه و قابل‌فهم برای کاربر
- نتیجه ساخت‌یافته شامل:
  - `workflowId`
  - `name`
  - `status`
  - `currentStep`
  - `planningSource`
  - `lastLog`
  - `updatedAt`
- در صورت نیاز، فهرست 1 تا 5 workflow منتخب با reason کوتاه برای انتخاب یا highlight
- failure روشن و قابل‌نمایش:
  - `workflow_not_found`
  - `no_workflows_available`
  - `workflow_data_unavailable`
  - `logs_unavailable`

### surfaceهای هدف در UI

- `Foundry`
- `Projects`
- `GlobalChat` یا surface اختصاصی agent result

## گزینه‌های عمداً ردشده

### گزینه B: Tool Catalog / Manifest Agent

- با اینکه endpoint و contract آن ساده است، خروجی‌اش بیش از حد نزدیک به `tool listing` می‌ماند
- برای Day 2 ارزش product-facing کمتری نسبت به insight روی runtime واقعی دارد
- می‌تواند بعداً به‌عنوان capability دوم agent اضافه شود، اما انتخاب اول مناسبی برای بستن `agent that does real work` نیست

### گزینه C: Log Summary Agent

- از نظر safety مناسب است، اما quality آن بیش از حد به غنای snapshot لاگ‌ها وابسته است
- در استک فعلی، لاگ‌ها بیشتر برای debugging مناسب‌اند تا اولین experience قابل‌فهم agent برای کاربر
- بخشی از نیاز این گزینه با `Workflow Insight Agent` هم پوشش داده می‌شود، چون workflow insight می‌تواند آخرین log مهم را هم نشان بدهد

## contract پیشنهادی برای Day 3

- draft agent در `Foundry` باید یک execution capability روشن داشته باشد:
  - `execution_mode: read_only_workflow_insight`
  - `result_surface: global_chat`
  - `tools: ["workflow_insight"]`
- `BFF` باید یک tool جدید non-deploy تعریف کند که query کاربر را به retrieval ساده از workflow data تبدیل کند
- `Go Gateway` لازم نیست orchestration جدید بسازد؛ کافی است read modelهای موجود `GET /workflows`، `GET /workflows/:id` و `GET /logs` را برای BFF قابل‌مصرف نگه دارد
- contract پاسخ باید علاوه بر متن summary، payload ساخت‌یافته‌ی workflow insight را هم برگرداند تا فرانت بتواند result card یا expandable detail نشان بدهد

## endpointها و contractهای موردنیاز

- `BFF /api/chat`
  - افزودن tool جدید مثل `workflow_insight`
  - routing intentهای read-only workflow به‌جای `deploy_website`
- `Go Gateway /workflows`
  - لیست workflow executionها برای انتخاب یا summarize
- `Go Gateway /workflows/:id`
  - جزئیات workflow انتخاب‌شده
- `Go Gateway /logs`
  - گرفتن آخرین log مهم برای enrich کردن insight
- contract پاسخ اولیه:
  - `status`
  - `summary`
  - `selected_workflows`
  - `planning_source`
  - `last_log`
  - `error`

## نتیجه نهایی روز

- use case نهایی Day 2 برابر با `Workflow Insight Agent` است
- مسیر `deploy_website` عمداً از scope agent اول کنار گذاشته می‌شود
- Day 3 باید mapping بین draft agent و capability `workflow_insight` را رسمی کند
- Day 4 باید اجرای read-only این agent را در `GlobalChat` یا result surface نزدیک به آن سبز کند

## فایل‌های محتمل برای روز بعد

- `portal1/src/lib/components/foundry/Foundry.svelte`
- `portal1/src/lib/components/chat/GlobalChat.svelte`
- `portal1/src/lib/services/supernode.ts`
- `BFF/index.ts`
- `internal/adapters/gateway/echo.go`

## validation روز

این روز validation اجرایی کامل ندارد، اما خروجی تصمیم باید با این سوال validate شود:

> آیا اگر همین use case فردا پیاده‌سازی شود، کاربر حس می‌کند یک agent ساده واقعاً کاری برایش انجام داده است؟

اگر پاسخ مبهم باشد، انتخاب این روز درست نبوده است.

## معیار done

- یک use case واقعی، کم‌ریسک و non-demo انتخاب شده باشد
- output آن برای user قابل‌فهم باشد
- dependencyهای آن با استک فعلی هم‌خوان باشند
- scope آن برای Day 3 و Day 4 قابل‌اجرا باشد

## carry-over

اگر تصمیم نهایی در این روز بسته نشد، فقط این موارد به روز بعد منتقل می‌شوند:

- مقایسه محدود بین دو گزینه نهایی
- ثبت dependency یا blocker دقیق

هر نوع گسترش scope بیشتر از این، failure همین روز محسوب می‌شود.
