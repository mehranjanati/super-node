# پلن توسعه روزانه

این سند مرجع اجرایی روزبه‌روز برای توسعه MVP است. مبنای اولویت‌ها:

- مرجع اصلی محصول: `MVP_DEVELOPMENT_PLAN.md`
- مرجع اجرایی BFF: `BFF_IMPLEMENTATION_PLAN.md`
- مرجع اتصال Go: `GO_INTEGRATION_PLAN.md`

---

## اصل اجرا

قانون این پلن:

1. هر روز فقط یک Vertical Slice کوچک اما کامل جلو می‌رود.
2. هر گام باید خروجی قابل تست داشته باشد.
3. اول `curl` و تست دستی، بعد UI، بعد تست خودکار اگر ارزشمند بود.
4. تا قبل از کامل شدن زنجیره `SPA -> BFF -> Go` وارد LiveKit، Matrix، Redpanda و بخش‌های سنگین‌تر نمی‌شویم.

---

## تعریف Done برای هر روز

هر روز وقتی تمام است که این 4 مورد سبز باشند:

- کد build/check شود
- مسیر جدید با `curl` تست شود
- رفتار از UI یا لایه مصرف‌کننده قابل مشاهده باشد
- خطای واضح و fallback مناسب داشته باشد

---

## روز 1: Health Contract

هدف:

- شفاف شدن وضعیت سلامت BFF و Go
- فراهم شدن پایه تست برای روزهای بعد

کارها:

- افزودن `GET /api/health` در BFF
- افزودن یا یکسان‌سازی `GET /internal/health` در Go
- برگرداندن JSON ساخت‌یافته به جای متن ساده
- تست با `curl`

تست:

```bash
curl http://localhost:3001/api/health
curl http://localhost:3000/internal/health
```

خروجی مورد انتظار:

```json
{
  "status": "ok"
}
```

---

## روز 2: Contract بین BFF و Go

هدف:

- تثبیت قرارداد request/response برای اولین tool

کارها:

- تعریف payload استاندارد برای `deploy_website`
- تعیین endpoint داخلی Go
- تعیین ساختار خطا
- تست مستقیم Go با `curl`

تست:

```bash
curl -X POST http://localhost:3000/internal/tools/deploy \
  -H 'Content-Type: application/json' \
  -d '{"project_name":"demo-site","template":"svelte"}'
```

---

## روز 3: Tool Calling واقعی

هدف:

- جایگزینی mock فعلی در BFF با درخواست واقعی به Go

کارها:

- حذف mock در `BFF/index.ts`
- `fetch` واقعی از BFF به Go
- مدیریت timeout و خطای Go
- نمایش نتیجه tool در چت

تست:

- `curl` مستقیم به BFF
- تست از UI chat با prompt ثابت

---

## روز 4: تست مسیر کامل

هدف:

- اثبات زنجیره کامل از UI تا Go

کارها:

- تعریف prompt پایدار برای trigger ابزار
- بررسی render نتیجه tool در UI
- تست خطای Go down
- تست invalid input

تست:

- ارسال prompt از UI
- توقف Go و مشاهده fallback

---

## روز 5: Temporal Dummy Workflow

هدف:

- اجرای tool از طریق workflow ساده به‌جای اجرای مستقیم

کارها:

- ایجاد workflow کوتاه و deterministic
- trigger از Go
- ثبت status اولیه

تست:

- `curl` به Go
- بررسی لاگ یا status برگشتی

---

## روز 6: Persistence و Read Model

هدف:

- ثبت نتیجه اجرای tool در دیتابیس

کارها:

- ثبت execution record
- در دسترس قرار دادن از طریق Hasura یا endpoint خواندن
- آماده‌سازی برای داشبورد

تست:

- write از chat
- read از DB/Hasura

---

## روز 7: پایداری و تمیزکاری

هدف:

- پایدار شدن تجربه توسعه

کارها:

- بهبود پیام‌های خطا
- بررسی env ها
- چک health ها قبل از اجرا
- مستندسازی curlهای اصلی

تست:

- smoke test کل مسیر

---

## ترتیب اجرای واقعی

ترتیب واقعی همین حالا:

1. روز 1 را اجرا می‌کنیم
2. بعد از سبز شدن healthها می‌رویم روز 2
3. سپس mock tool را حذف می‌کنیم و روز 3 را جلو می‌بریم

---

## فایل‌های درگیر اولیه

- `BFF/index.ts`
- `internal/adapters/gateway/echo.go`
- `cmd/nexus-super-node/main.go`
- `portal1/src/lib/components/chat/GlobalChat.svelte`

---

## وضعیت

- [x] روز 1: Health Contract
- [x] روز 2: Contract بین BFF و Go
- [x] روز 3: Tool Calling واقعی
- [x] روز 4: تست مسیر کامل
- [ ] روز 5: Temporal Dummy Workflow
- [ ] روز 6: Persistence و Read Model
- [ ] روز 7: پایداری و تمیزکاری
