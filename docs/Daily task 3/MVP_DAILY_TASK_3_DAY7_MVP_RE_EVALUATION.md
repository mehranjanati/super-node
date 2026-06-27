# MVP Daily Task 3 - Day 7: MVP Re-Evaluation

هدف روز:

- تصمیم دوباره برای Go/No-Go بر اساس agent MVP واقعی
- ثبت handoff روشن برای ادامه مسیر بعد از این سری

## چرا این روز حیاتی است

Day 7 در این سری نباید فقط تکرار Day 7 از `Daily task 2` باشد.

در `Daily task 2` معیار اصلی بیشتر حول deploy path و artifactهای آن بود.

اما در `Daily task 3` سوال نهایی عوض شده است:

> آیا کاربر حالا می‌تواند یک agent ساده، کم‌ریسک و واقعی را تعریف کند، ذخیره کند، اجرا کند و نتیجه قابل‌فهم ببیند؟

اگر این سوال هنوز پاسخ شفاف نداشته باشد:

- سری Daily Task 3 هنوز MVP کاربردی را نبسته است
- حتی اگر بخش‌هایی از زیرساخت یا UI خوب کار کنند، claim نهایی محصول مبهم می‌ماند

## وضعیت شروع روز

- Day 1 باید draft flow واقعی را بسته باشد
- Day 2 باید use case واقعی را بسته باشد
- Day 3 باید contract بین draft و runtime را بسته باشد
- Day 4 باید اجرای read-only agent را سبز کرده باشد
- Day 5 باید result surface را harden کرده باشد
- Day 6 باید smoke path رسمی و evidence آن را ثبت کرده باشد

## سوال نهایی این روز

در پایان این روز فقط باید به این سوال پاسخ داده شود:

> آیا مسیر `design -> save -> ask -> visible result` برای agent MVP به اندازه کافی واقعی، کم‌ریسک و قابل‌فهم شده است که مبنای MVP کاربردی قرار بگیرد؟

## کارهای امروز

1. جمع‌آوری evidence روزهای 1 تا 6
2. ثبت checklist نهایی:
   - چه چیزهایی سبز هستند
   - چه چیزهایی defer شده‌اند
   - چه ریسک‌هایی پذیرفته شده‌اند
3. تصمیم Go/No-Go برای این سری
4. ثبت handoff کوتاه برای بعد از این سری

## checklist تصمیم

### سبز لازم

- کاربر می‌تواند draft agent بسازد
- draft agent ذخیره و دوباره باز می‌شود
- capability agent مشخص و non-demo است
- `GlobalChat` یا result surface معادل، result قابل‌فهم نشان می‌دهد
- failureهای پایه misleading نیستند
- smoke path رسمی ثبت و قابل تکرار است

### defer قابل‌قبول

- backend sync برای draftها
- multi-user state
- capability دوم مثل `manifest/tool catalog`
- surface اختصاصی و کامل برای result history
- evalها و telemetry عمیق‌تر

### No-Go triggers

- execution فقط با توضیح شفاهی قابل درک باشد
- result همچنان بیشتر شبیه raw tool payload باشد تا agent output
- failure mode اصلی کاربر را گمراه کند
- smoke path رسمی ناقص یا غیرقابل تکرار باشد
- use case اصلی دوباره عملاً به deploy-centric path برگشته باشد

## Go/No-Go پیشنهادی این سری

### Go

- اگر `Workflow Insight Agent` واقعاً از draft انتخاب‌شده اجرا شود
- اگر user نتیجه‌ای ببیند که شامل summary یا structured result قابل‌فهم باشد
- اگر مسیر `design -> save -> ask -> result` بدون ambiguity مستند شده باشد
- اگر failureهای اصلی روشن و low-risk باشند

### No-Go

- اگر contract فقط در کد وجود داشته باشد ولی در تجربه کاربر معلوم نباشد
- اگر output فقط برای developer قابل‌فهم باشد
- اگر validation رسمی این سری بسته نشده باشد
- اگر regression یا scope creep باعث شود agent MVP دوباره نمایشی شود

## نتیجه ارزیابی نهایی (Evaluation Result)

### وضعیت نهایی: Go ✅

بر اساس شواهد (Evidence) ثبت شده در روزهای گذشته، به خصوص موفقیت اسکریپت `smoke validation` در **Day 6**:

1. **Workflow Insight Agent** به عنوان یک use-case کم‌ریسک و کاملاً کاربردی از روی draft انتخاب‌شده اجرا می‌شود.
2. قرارداد (Contract) قوی بین `Foundry` (ساخت)، `Projects` (نگهداری)، و `GlobalChat` (اجرا و نمایش) برقرار شده است.
3. خروجی در چت صرفاً یک raw payload نیست، بلکه یک **Structured Result Card** و یک Summary متنی قابل فهم است.
4. مسیر خطاها (مانند `workflow_not_found`) به درستی شناسایی و به کاربر در UI بازگردانده می‌شود.
5. تیم می‌تواند به راحتی با اجرای `./scripts/mvp_real_agent_smoke.sh` کل این مسیر را به شکل ایزوله و مطمئن تست کند.

---

## handoff بعد از این سری

### backlog کوتاه

- اضافه‌کردن capability دوم بعد از تثبیت `workflow_insight`
- ثبت history یا session summary برای agent resultها
- همگام‌سازی draftها با backend در صورت نیاز product
- افزودن smoke script رسمی قابل اجرا در CI سبک یا محیط dev

### debtهای باقی‌مانده

- local-first بودن persistence draft
- تکیه result surface به `GlobalChat`
- نیاز به polish بیشتر برای `Projects` و surfaceهای غیر chat

### اولویت‌های post-series

- hardening نهایی UX
- validation تکرارپذیر در محیط نزدیک‌تر به release
- تصمیم‌گیری برای گسترش capabilityها فقط بعد از تثبیت path اول

## معیار done

- یک پاسخ شفاف Go/No-Go برای agent MVP وجود داشته باشد
- evidence کافی برای آن پاسخ ثبت شده باشد
- handoff برای کارهای بعدی کوتاه، دقیق و غیرمبهم باشد

## خروجی مورد انتظار

در پایان Day 7 باید بتوان با اطمینان گفت:

> MVP این سری دیگر صرفاً "agent draft با ظاهر خوب" نیست، بلکه یک agent ساده و واقعی دارد که در scope محدود خود کار قابل‌فهم انجام می‌دهد؛ یا اگر هنوز این معیار بسته نشده، دلیل دقیق No-Go ثبت شده است.
