# برنامه توسعه فازبندی شده (MVP Development Plan)

این سند مرجع محصول برای MVP است، اما باید با تسک‌های روزانه داخل `docs/Daily task/` هم‌راستا بماند. هدف این نسخه این است که بین «پیشروی سریع معماری» و «دیدن تدریجی نتیجه در فرانت و تست‌پذیری» تعادل ایجاد کند؛ یعنی هر گام کوچک، هم خروجی قابل مشاهده داشته باشد، هم validation روشن، هم docs قابل اتکا.

---

## هدف MVP

نسخه MVP باید فقط این هسته را به‌صورت پایدار و قابل تکرار تحویل بدهد:

1. یک مسیر اصلی کارا و قابل مشاهده از `SPA/BFF -> Go -> voltagent-service -> Go -> Temporal`
2. امکان اجرای حداقل یک use case واقعی مثل `deploy_website`
3. امکان دیدن اثر اجرا در یک مصرف‌کننده واقعی مثل chat، builder، workflows یا logs
4. امکان تست دستی و تست هدفمند برای همین مسیر اصلی
5. حداقل visibility برای health، error و fallback

موارد زیر در این MVP جزو اولویت اول نیستند:

- تماس‌های صوتی/تصویری LiveKit
- شبکه‌های چندکاناله کامل با Matrix و OpenClaw
- کلاسترهای سنگین Wasm
- داشبوردهای بسیار کامل و production-grade قبل از تثبیت مسیر اصلی
- code generation از صفر برای سایت‌های پیچیده

توجه استراتژیک:

- رویکرد اصلی همچنان **Template-Driven** است، نه تولید کد پیچیده از صفر
- day fileها مرجع اجرای migration هستند و این سند باید فقط جهت‌گیری محصول و ترتیب sliceها را روشن کند

---

## تعادل بین MVP و Daily Task

برای جلوگیری از drift بین برنامه MVP و کارهای روزانه، این قانون‌ها باید رعایت شوند:

1. `docs/MVP_DEVELOPMENT_PLAN.md` مرجع اولویت محصول و ترتیب sliceها است.
2. فایل‌های داخل `docs/Daily task/` مرجع اجرای روزبه‌روز migration و carry-over هستند.
3. هر slice جدید فقط وقتی ارزش دارد که:
   - با `curl` یا تست دستی قابل تایید باشد
   - از یک لایه مصرف‌کننده مثل SPA یا BFF قابل مشاهده باشد
   - اگر ریسک regression دارد، یک test هدفمند هم برایش اضافه شود
4. قبل از کامل‌شدن زنجیره اصلی، وارد featureهای سنگین‌تر مثل LiveKit، Matrix، Redpanda UI و polishهای کم‌ارزش نمی‌شویم.
5. هر تغییری که فقط backend را بهتر کند ولی از سمت مصرف‌کننده قابل دیدن نباشد، باید هم‌زمان یک surface مشاهده یا plan تست روشن هم داشته باشد.

---

## تعریف Done برای هر Slice

هر slice MVP وقتی done محسوب می‌شود که این موارد سبز باشند:

- کد مربوطه `build/check` یا حداقل compile pass داشته باشد
- مسیر جدید با `curl` یا تست دستی قابل تکرار باشد
- نتیجه از یک consumer واقعی مثل SPA، BFF chat، workflows یا logs دیده شود
- خطا، timeout یا fallback رفتار صریح و قابل تشخیص داشته باشد
- docs همان روز و در صورت نیاز این فایل MVP به‌روز شده باشند

وقتی ارزش افزوده واضح دارد، test خودکار هم اضافه می‌شود:

- unit test برای contract یا mapping
- integration test برای endpoint اصلی
- e2e یا smoke test فقط برای مسیرهای حیاتی و کم‌نویز

---

## وضعیت فعلی نسبت به MVP

با توجه به repo فعلی و day fileهای migration، پروژه الان در میانه MVP قرار دارد:

- foundation اولیه فرانت، BFF و Go وجود دارد
- مسیر `deploy_website` تا حد خوبی به routeهای جدید align شده است
- `portal1` در لایه service برای manifest و execute و deploy به contract داخلی `/internal/tools/*` نزدیک شده است
- workflow/log surface در فرانت وجود دارد، اما همه صفحه‌ها هنوز به backend واقعی وصل نیستند
- health، observability و validation هنوز نیاز به تکمیل هدفمند دارند

نتیجه عملی:

- تمرکز فعلی باید روی کامل‌کردن مسیر اصلی و قابل‌دیدن‌کردن آن در SPA باشد
- نه روی بازکردن featureهای جدید خارج از MVP

---

## فاز ۱: Foundation و Visibility پایه

هدف:

- تثبیت health contract
- آماده‌کردن پایه‌ای که از همان ابتدا در فرانت یا BFF قابل مشاهده باشد

کارها:

- بالا آوردن `BFF` با `GET /api/health`
- بالا آوردن `Go` با `GET /internal/health`
- یکسان‌سازی envelopeهای ساده برای health و error
- تثبیت پروژه `SvelteKit` و کامپوننت‌های پایه
- داشتن حداقل یک screen یا status card که سلامت سرویس‌ها را نشان بدهد

خروجی قابل مشاهده:

- health endpointها با پاسخ JSON ساخت‌یافته
- امکان دیدن وضعیت سلامت از طریق فرانت یا حداقل از یک صفحه/پنل debug

validation:

```bash
curl http://localhost:3001/api/health
curl http://localhost:3000/internal/health
```

هم‌راستایی با day fileها:

- `docs/Daily task/DEVELOP_ROOZANEH.md`
- `docs/Daily task/VOLTAGENT_SERVICE_DAY4_HEALTH_AND_OBSERVABILITY.md`

---

## فاز ۲: Contract داخلی و First Real Slice

هدف:

- بستن اولین vertical slice واقعی برای `deploy_website`
- حذف ambiguity بین callerهای edge و backend

کارها:

- تثبیت contract برای `POST /internal/tools/deploy`
- شفاف‌کردن contractهای `GET /internal/tools/manifest` و `POST /internal/tools/execute`
- هم‌راستا کردن `BFF` و `portal1` با routeهای جدید
- نگه‌داشتن routeهای legacy فقط برای compatibility backend
- ثبت `workflow_id`, `planning_source`, `current_step` و خطاهای اصلی

خروجی قابل مشاهده:

- امکان trigger کردن deploy از یک consumer واقعی
- دیده‌شدن workflow یا log یا status در فرانت
- قابل مشاهده بودن failureهای request یا timeout

validation:

```bash
curl -X POST http://localhost:3000/internal/tools/deploy \
  -H 'Content-Type: application/json' \
  -d '{"project_name":"demo-site","template":"svelte"}'
```

تست‌های ارزشمند این فاز:

- unit test برای serviceهای فرانت که route درست را صدا می‌زنند
- test هدفمند برای gateway یا contract mapping
- تست دستی از UI برای اطمینان از visible بودن نتیجه

هم‌راستایی با day fileها:

- `docs/Daily task/VOLTAGENT_SERVICE_DAY2_BACKEND_RUNTIME.md`
- `docs/Daily task/VOLTAGENT_SERVICE_DAY3_BFF_PORTAL_ALIGNMENT.md`

---

## فاز ۳: Workflow، Read Model و Surface قابل مشاهده

هدف:

- تبدیل execution ساده به جریان قابل پیگیری
- نشان‌دادن status و log در یک surface واقعی در SPA

کارها:

- trigger شدن workflow از Go به Temporal
- persist شدن execution record و status
- فراهم‌شدن read model برای workflowها و logها
- وصل شدن صفحه‌هایی مثل `Workflows` و `Logs` به داده واقعی
- اگر builder صفحه فعال محصول است، mount شدن route آن در SPA فعال

خروجی قابل مشاهده:

- ایجاد workflow بعد از `deploy_website`
- دیده‌شدن execution در `#/workflows` یا surface مشابه
- دیده‌شدن logها یا status به‌صورت قابل فهم برای توسعه‌دهنده

validation:

- trigger از chat یا builder یا endpoint مستقیم
- مشاهده workflow در UI
- بررسی log/status از endpoint یا consumer

تست‌های ارزشمند این فاز:

- test برای persistence/read model
- test برای mapping workflow status
- smoke test سبک برای use case رسمی phase 1

هم‌راستایی با day fileها:

- `docs/Daily task/DEVELOP_ROOZANEH.md`
- `docs/Daily task/VOLTAGENT_SERVICE_DAY6_TESTS_AND_VALIDATION.md`

---

## فاز ۴: Observability، Hardening و Validation تکرارپذیر

هدف:

- بالا بردن confidence مسیر اصلی بدون باز کردن scopeهای جدید

کارها:

- health تجمیعی `Go -> voltagent-service`
- latency و error classification
- logهای روشن برای `remote_voltagent` و `embedded_fallback`
- env matrix و readiness
- smoke test برای stack
- validation رسمی مسیر `deploy_website` در حالت remote و fallback

خروجی قابل مشاهده:

- health معنادارتر در endpoint و پنل‌های فرانت
- امکان تشخیص واضح failure mode
- confidence بیشتر برای ادامه توسعه UI و اضافه‌کردن تست‌ها

تست‌های ارزشمند این فاز:

- unit test برای `voltagentclient`
- integration test برای `POST /plan`
- smoke test برای compose
- end-to-end validation برای `deploy_website`

هم‌راستایی با day fileها:

- `docs/Daily task/VOLTAGENT_SERVICE_DAY4_HEALTH_AND_OBSERVABILITY.md`
- `docs/Daily task/VOLTAGENT_SERVICE_DAY5_INFRA_HARDENING.md`
- `docs/Daily task/VOLTAGENT_SERVICE_DAY6_TESTS_AND_VALIDATION.md`

---

## فازهای بعد از MVP

فقط بعد از سبز شدن مسیر اصلی و validation آن وارد این بخش‌ها می‌شویم:

- Redpanda event flows کامل و dashboardهای real-time پیشرفته
- agent orchestrationهای چندمرحله‌ای فراتر از use case اصلی
- LiveKit، Matrix و OpenClaw به‌عنوان featureهای قابل عرضه
- polish کامل فرانت، coverage گسترده و release hardening پیشرفته

---

## استراتژی ادغام

برای جلوگیری از سردرگمی، قانون زیر باید رعایت شود:

1. هدف نهایی محصول این است که عملیات write/chat از لایه edge از طریق `BFF` عبور کند.
2. در فاز migration فعلی، استفاده صریح `portal1` از contract داخلی `/internal/tools/*` قابل قبول است، چون هدفش حذف وابستگی به routeهای legacy و کم‌کردن drift است.
3. تمام readهای پایدار و محصولی ترجیحاً از طریق Hasura یا read modelهای مشخص انجام می‌شوند، نه از endpointهای موقت و ad-hoc.
4. `Go` فقط با پیام‌های ساخت‌یافته و contract-aware کار می‌کند و وظیفه‌اش اجرای منطق، workflow و orchestration است.
5. هر مسیر موقت که برای migration نگه داشته می‌شود باید explicit در day fileها ثبت شود و plan حذف داشته باشد.

---

## قاعده اولویت از این نقطه به بعد

اگر بین چند کار مردد بودیم، اولویت با این ترتیب است:

1. چیزی که زنجیره اصلی `SPA/BFF -> Go -> voltagent-service -> Go -> Temporal` را کامل‌تر می‌کند
2. چیزی که نتیجه را زودتر در فرانت یا consumer قابل مشاهده می‌کند
3. چیزی که confidence را با test یا smoke validation بالا می‌برد
4. چیزی که observability و failure mode را روشن‌تر می‌کند
5. featureهای جدید خارج از مسیر اصلی
